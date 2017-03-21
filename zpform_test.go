package zpform

import (
	"net/http"
	"testing"
)

func TestParseSimple(t *testing.T) {
	req, err := http.NewRequest("GET", "http://some.domain/?a=1&b=2", nil)
	if err != nil {
		t.Error(err.Error())
	}
	var form struct {
		A int
		B string
	}
	if err := ReadStructForm(req, &form); err != nil {
		t.Error(err.Error())
	}
	if form.A != 1 {
		t.Error("parse A fail")
	}
	if form.B != "2" {
		t.Error("parse B fail")
	}
}

func TestParseSlice(t *testing.T) {
	req, err := http.NewRequest("GET", "http://some.domain/?a=1&a=2&b[]=3&b[]=4", nil)
	if err != nil {
		t.Error(err.Error())
	}
	var form struct {
		A []int
		B []int
	}
	if err := ReadStructForm(req, &form); err != nil {
		t.Error(err.Error())
	}
	if len(form.A) != 2 || form.A[0] != 1 || form.A[1] != 2 {
		t.Error("parse A fail")
	}
	if len(form.B) != 2 || form.B[0] != 3 || form.B[1] != 4 {
		t.Error("parse B fail")
	}
}
