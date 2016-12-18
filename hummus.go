package hummus

import (
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

type gabsTag struct {
	tagName   string
	omitEmpty bool
}

func Marshal(input interface{}) ([]byte, error) {
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)

	jsonObj, err := marshalReflect(t, v)
	if err != nil {
		panic(err)
	}

	return []byte(jsonObj.String()), nil
}

func marshalReflect(t reflect.Type, v reflect.Value) (*gabs.Container, error) {
	jsonObj := gabs.New()

	for i := 0; i < t.NumField(); {
		gt := parseGabsTag(t.Field(i).Tag)

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
				childTag := fmt.Sprintf("gabs:%q", at.childPath)

				childFields = append(childFields, reflect.StructField{
					Name:   t.Field(i).Name,
					Type:   t.Field(i).Type,
					Tag:    reflect.StructTag(childTag),
					Offset: 0,
				})

				childValues = append(childValues, v.Field(i).Interface())
				for j < t.NumField() {
					gtn := parseGabsTag(t.Field(j).Tag)
					if atn, ok := parseArrayTag(gtn.tagName); ok &&
						atn.arrayPath == at.arrayPath && atn.arrayIndex == at.arrayIndex {
						nextChildTag := fmt.Sprintf("gabs:%q", atn.childPath)
						childFields = append(childFields, reflect.StructField{
							Name:   t.Field(j).Name,
							Type:   t.Field(j).Type,
							Tag:    reflect.StructTag(nextChildTag),
							Offset: 0,
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
					panic(err)
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
		i = j
	}

	return jsonObj, nil
}

func parseGabsTag(tag reflect.StructTag) gabsTag {
	gabsTagString := tag.Get("gabs")
	tagFields := strings.Split(gabsTagString, ",")
	var omitEmpty bool

	if len(tagFields) > 2 {
		panic("hello")
	}

	if len(tagFields) == 2 && tagFields[1] == "omitempty" {
		omitEmpty = true
	}

	return gabsTag{
		tagName:   tagFields[0],
		omitEmpty: omitEmpty,
	}
}

func parseArrayTag(tag string) (arrayTag, bool) {
	if strings.Contains(tag, "[") {
		arrayRegex := regexp.MustCompile("(.*)\\[(\\d+)\\]\\.*(.*)")
		matches := arrayRegex.FindStringSubmatch(tag)
		if len(matches) != 4 {
			panic("hi")
		}

		arrayIndex, err := strconv.Atoi(matches[2])
		if err != nil {
			panic(err)
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
