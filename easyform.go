package zpform

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ValidateFunc func(string) (bool, string)

type _F struct {
	FieldVar     *reflect.Value
	FieldName    string
	FieldLabel   string
	FieldFormat  string
	ValidateFunc []ValidateFunc
}

func NewF(fieldVar interface{}, fieldName, fieldLabel, fieldFormat string, validateFunc ...ValidateFunc) *_F {
	fieldVal := reflect.ValueOf(fieldVar)
	if reflect.Ptr != fieldVal.Kind() {
		panic("fieldVar must be a pointer")
	}
	validators := make([]ValidateFunc, len(validateFunc))
	for i := 0; i < len(validateFunc); i++ {
		validators[i] = validateFunc[i]
	}
	return &_F{
		FieldVar:     &fieldVal,
		FieldName:    fieldName,
		FieldLabel:   fieldLabel,
		FieldFormat:  fieldFormat,
		ValidateFunc: validators,
	}
}

func NewReF(fieldVar interface{}, fieldName, re string) *_F {
	return NewF(fieldVar, fieldName, fieldName, "", Regexp(re))
}

func setFieldValue(el reflect.Value, strVal, format string) error {
	if strVal == "" {
		return nil
	}
	switch el.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(strVal, 10, 64)
		if nil != err {
			return err
		}
		el.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(strVal, 10, 64)
		if nil != err {
			return err
		}
		el.SetUint(uintVal)
	case reflect.String:
		el.SetString(strVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(strVal, 64)
		if nil != err {
			return err
		}
		el.SetFloat(floatVal)
	case reflect.Bool:
		strVal = strings.ToLower(strVal)
		boolVal := strVal == "true" || strVal == "是" || strVal == "yes" || strVal == "y" || strVal == "on"
		el.SetBool(boolVal)
	case reflect.Struct, reflect.Ptr:
		switch el.Interface().(type) {
		case time.Time:
			t, err := time.Parse(format, strVal)
			if err != nil {
				return err
			}
			el.Set(reflect.ValueOf(t))
		default:
			fname := "Parse"
			if el.Kind() == reflect.Ptr {
				el.Set(reflect.New(el.Type().Elem()))
			}
			if el.Kind() == reflect.Struct {
				el = el.Addr()
			}
			f := el.MethodByName(fname)
			if !f.IsValid() {
				return errors.New(fmt.Sprintf("Struct has no method '%s'", fname))
			}
			if errVal := f.Call([]reflect.Value{
				reflect.ValueOf(format),
				reflect.ValueOf(strVal),
			}); len(errVal) > 0 && !errVal[0].IsNil() {
				return errVal[0].Interface().(error)
			}
		}
	default:
		return errors.New("unaccept fieldvar type")
	}
	return nil
}

func getErrorMsg(fieldName, fieldLabel, errMsg string) error {
	if errMsg == "" {
		errMsg = "not valid"
	}
	var name string
	if fieldLabel == "" {
		name = fieldName
	} else {
		name = fieldLabel
	}
	return errors.New(fmt.Sprintf("%v（%v）",
		name, errMsg,
	))
}

func ReadForm(req *http.Request, f ...*_F) error {
	req.ParseForm()
	for i := 0; i < len(f); i++ {
		finfo := f[i]
		el := finfo.FieldVar.Elem()
		if reflect.Slice == el.Kind() {
			postValues, exists := req.Form[finfo.FieldName]
			if !exists {
				postValues, exists = req.Form[finfo.FieldName+"[]"]
				if !exists {
					continue
				}
			}
			elemType := el.Type().Elem()
			for _, postValue := range postValues {
				for _, validator := range finfo.ValidateFunc {
					if ok, msg := validator(postValue); !ok {
						return getErrorMsg(
							finfo.FieldName, finfo.FieldLabel, msg,
						)
					}
				}
				newElem := reflect.New(elemType)
				if err := setFieldValue(newElem.Elem(), postValue, finfo.FieldFormat); err != nil {
					return err
				}
				el.Set(reflect.Append(el, newElem.Elem()))
			}
		} else {
			postValue := req.FormValue(finfo.FieldName)
			for _, validator := range finfo.ValidateFunc {
				if ok, msg := validator(postValue); !ok {
					return getErrorMsg(
						finfo.FieldName, finfo.FieldLabel, msg,
					)
				}
			}
			if err := setFieldValue(el, postValue, finfo.FieldFormat); err != nil {
				return err
			}
		}
	}
	return nil
}

func ReadStructForm(req *http.Request, form interface{}) error {
	return ReadForm(req, GetBindings(form)...)
}

func ReadReflectedStructForm(req *http.Request, formValue reflect.Value) error {
	return ReadForm(req, GetReflectedBindings(formValue)...)
}

func ReadFileForm(req *http.Request, form interface{}) error {
	f := GetBindings(form)
	req.ParseMultipartForm(32 << 20)
	if req.MultipartForm == nil || req.MultipartForm.File == nil {
		return errors.New("no file upload")
	}
	for i := 0; i < len(f); i++ {
		finfo := f[i]
		el := finfo.FieldVar.Elem()
		if reflect.Slice == el.Kind() {
			postValues, exists := req.MultipartForm.File[finfo.FieldName]
			if !exists {
				postValues, exists = req.MultipartForm.File[finfo.FieldName+"[]"]
				if !exists {
					continue
				}
			}
			elemType := el.Type().Elem()
			for _, postValue := range postValues {
				newElem := reflect.New(elemType)
				nel := newElem.Elem()
				nel.Set(reflect.ValueOf(postValue))
				el.Set(reflect.Append(el, newElem.Elem()))
			}
		} else {
			panic("invalid field: " + finfo.FieldName)
		}
	}
	return nil
}
