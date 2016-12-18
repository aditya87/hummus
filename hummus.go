package hummus

import (
	"reflect"

	"github.com/jeffail/gabs"
)

func Marshal(input interface{}) ([]byte, error) {
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)

	jsonObj := gabs.New()

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("gabs")
		jsonObj.SetP(v.Field(i).Interface(), tag)
	}

	return []byte(jsonObj.String()), nil
}
