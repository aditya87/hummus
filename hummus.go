package hummus

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeffail/gabs"
)

type arrayTag struct {
	arrayPath  string
	arrayIndex int
	childPath  string
}

type hummusTag struct {
	tagName   string
	omitEmpty bool
}

func Marshal(input interface{}) ([]byte, error) {
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)

	jsonObj, err := marshalReflect(t, v)
	if err != nil {
		return []byte{}, err
	}

	return []byte(jsonObj.String()), nil
}

func marshalReflect(t reflect.Type, v reflect.Value) (*gabs.Container, error) {
	jsonObj := gabs.New()

	for i := 0; i < t.NumField(); {
		gt, err := parseGabsTag(t.Field(i).Tag)
		if err != nil && strings.Contains(err.Error(), "invalid struct tag") {
			i++
			continue
		} else if err != nil {
			return nil, err
		}

		if gt.omitEmpty && isEmptyValue(v.Field(i)) {
			i++
			continue
		}

		path := gt.tagName
		array := false
		object := v.Field(i).Interface()
		j := i + 1

		if at, ok := parseArrayTag(path); ok {
			if !jsonObj.ExistsP(at.arrayPath) {
				jsonObj.ArrayP(at.arrayPath)
			}

			if at.childPath != "" {
				var childFields []reflect.StructField
				var childValues []interface{}
				childTag := fmt.Sprintf("hummus:%q", at.childPath)
				childFields = append(childFields, reflect.StructField{
					Name: t.Field(i).Name,
					Type: t.Field(i).Type,
					Tag:  reflect.StructTag(childTag),
				})

				childValues = append(childValues, v.Field(i).Interface())
				for j < t.NumField() {
					gtn, err := parseGabsTag(t.Field(j).Tag)
					if err != nil {
						return nil, err
					}

					if atn, ok := parseArrayTag(gtn.tagName); ok &&
						atn.arrayPath == at.arrayPath && atn.arrayIndex == at.arrayIndex {
						if gtn.omitEmpty && isEmptyValue(v.Field(j)) {
							j++
							continue
						}
						nextChildTag := fmt.Sprintf("hummus:%q", atn.childPath)
						childFields = append(childFields, reflect.StructField{
							Name: t.Field(j).Name,
							Type: t.Field(j).Type,
							Tag:  reflect.StructTag(nextChildTag),
						})
						childValues = append(childValues, v.Field(j).Interface())
						j++
					} else {
						break
					}
				}

				var err error
				var childObject *gabs.Container
				childType := reflect.StructOf(childFields)
				childValue := reflect.New(childType).Elem()
				for k, v := range childValues {
					childValue.Field(k).Set(reflect.ValueOf(v))
				}
				childObject, err = marshalReflect(childType, childValue)
				if err != nil {
					return nil, err
				}

				object = childObject.Data()
			}

			array = true
			path = at.arrayPath
		}

		if array {
			jsonObj.ArrayAppendP(object, path)
		} else {
			jsonObj.SetP(object, path)
		}
		replaceHashTags(jsonObj, path)
		i = j
	}

	return jsonObj, nil
}

func replaceHashTags(obj *gabs.Container, path string) {
	keys := strings.Split(path, ".")
	curKey := keys[0]
	curSubTree := obj.Path(curKey)
	if strings.Contains(curKey, "#") {
		obj.DeleteP(curKey)
		curKey = strings.Replace(curKey, "#", ".", -1)
		existingSubTree := obj.S(curKey)
		subTreeData := existingSubTree.Data()
		subTreeData = mergeObjects(subTreeData, curSubTree.Data())
		obj.Set(subTreeData, curKey)
	}
	if len(keys) != 1 {
		replaceHashTags(curSubTree, strings.Join(keys[1:], "."))
	}
}

func mergeObjects(dst interface{}, src interface{}) interface{} {
	dstMap, dmok := dst.(map[string]interface{})
	srcMap, smok := src.(map[string]interface{})
	dstArr, daok := dst.([]interface{})
	srcArr, saok := src.([]interface{})
	if dmok && smok {
		for k, v := range srcMap {
			dstMap[k] = mergeObjects(dstMap[k], v)
		}
		dst = dstMap
	} else if daok && saok {
		for _, v := range srcArr {
			dstArr = append(dstArr, v)
		}
		dst = dstArr
	} else {
		dst = src
	}

	return dst
}

func parseGabsTag(tag reflect.StructTag) (hummusTag, error) {
	hummusTagString := tag.Get("hummus")
	if hummusTagString == "" {
		return hummusTag{}, fmt.Errorf("error: invalid struct tag %s", tag)
	}

	tagFields := strings.Split(hummusTagString, ",")
	var omitEmpty bool

	if len(tagFields) > 2 {
		return hummusTag{}, errors.New("error: invalid number of struct tag fields")
	}

	if len(tagFields) == 2 && tagFields[1] == "omitempty" {
		omitEmpty = true
	}

	return hummusTag{
		tagName:   tagFields[0],
		omitEmpty: omitEmpty,
	}, nil
}

func parseArrayTag(tag string) (arrayTag, bool) {
	if strings.Contains(tag, "[") {
		//regex needs to be non-greedy in order to catch the parent array path first
		//(in case of arrays inside arrays)
		arrayRegex := regexp.MustCompile("(.*?)\\[(\\d+)\\]\\.*(.*)")

		matches := arrayRegex.FindStringSubmatch(tag)
		if len(matches) != 4 {
			return arrayTag{}, false
		}

		arrayIndex, err := strconv.Atoi(matches[2])
		if err != nil {
			return arrayTag{}, false
		}

		return arrayTag{
			arrayPath:  matches[1],
			arrayIndex: arrayIndex,
			childPath:  matches[3],
		}, true
	}
	return arrayTag{}, false
}

// straight-up stole this from encoding/json
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
