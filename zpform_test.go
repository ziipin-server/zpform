package zpform

import (
	"testing"
)

type TestForm struct {
	Name string `
				label:"姓名"
				max:"3"
				required
			`
	Age int `label:"年龄"`
}

func TestGenerateHTML(t *testing.T) {
	f := TestForm{Name: "123"}
	ParseStruct(f)
}
