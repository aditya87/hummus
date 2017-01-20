package hummus

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/aditya87/hummus/tree"
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
	parseTree := tree.NewTree()

	for i := 0; i < t.NumField(); i++ {
		err := parseTree.Insert(string(t.Field(i).Tag), v.Field(i).Interface(), isEmptyValue(v.Field(i)))
		if err != nil {
			return nil, err
		}
	}

	return parseTree.BuildJSON(), nil
}

func parseHummusTag(tag reflect.StructTag) (hummusTag, error) {
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
