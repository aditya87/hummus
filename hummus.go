package hummus

import (
	"reflect"
	"strings"

	"github.com/jeffail/gabs"
)

func Marshal(input interface{}) ([]byte, error) {
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)

	jsonObj := gabs.New()

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("gabs")
		tagFields := strings.Split(tag, ",")

		if len(tagFields) > 2 {
			panic("hello")
		}

		if len(tagFields) == 2 && tagFields[1] == "omitempty" && isEmptyValue(v.Field(i)) {
			continue
		}

		jsonObj.SetP(v.Field(i).Interface(), tagFields[0])
	}

	return []byte(jsonObj.String()), nil
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
