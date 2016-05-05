package zpform

import (
	"fmt"
	"reflect"
)

func ParseStruct(form interface{}) {
	formType := reflect.TypeOf(form)
	for i := 0; i < formType.NumField(); i++ {
		fieldObj := formType.Field(i)
		fmt.Println(fieldObj.Name)
		fmt.Println("label: ", fieldObj.Tag.Get("label"))
		fmt.Println("max: ", fieldObj.Tag.Get("required"))
	}
}
