package zpform

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

func panicErr(err error) {
	if nil != err {
		panic(err)
	}
}

func toSnake(camelStr string) string {
	newstr := make([]rune, 2*len(camelStr))
	idx := 0
	for _, ch := range camelStr {
		if ch >= 'A' && ch <= 'Z' {
			if idx > 0 {
				newstr[idx] = '_'
				idx++
			}
			ch += 'a' - 'A'
		}
		newstr[idx] = ch
		idx++
	}
	return string(newstr[:idx])
}

type _FieldMeta struct {
	FieldName   string
	Label       string
	Widget      string
	FieldFormat string
	Validators  []ValidateFunc
}

func getFieldName(fieldType *reflect.StructField) string {
	fname := fieldType.Tag.Get("zpf_name")
	if "" == fname {
		fname = toSnake(fieldType.Name)
	}
	return fname
}

func getValidators(tag *reflect.StructTag) []ValidateFunc {
	vf := make([]ValidateFunc, 0)
	if required := tag.Get("zpf_reqd"); required != "" && required != "false" {
		vf = append(vf, Required())
	}
	if regexp := tag.Get("zpf_re"); "" != regexp {
		vf = append(vf, Regexp(regexp))
	}
	if minLen := tag.Get("zpf_minlen"); "" != minLen {
		minInt, err := strconv.ParseInt(minLen, 10, 32)
		panicErr(err)
		vf = append(vf, LengthGT(int(minInt)))
	}
	if maxLen := tag.Get("zpf_maxlen"); "" != maxLen {
		maxInt, err := strconv.ParseInt(maxLen, 10, 32)
		panicErr(err)
		vf = append(vf, LengthLT(int(maxInt)))
	}
	if lenRange := tag.Get("zpf_len"); "" != lenRange {
		var min, max int
		n, err := fmt.Sscanf(lenRange, "%d %d", &min, &max)
		if 2 != n || nil != err {
			panic("invali zpf_len param: " + lenRange)
		}
		vf = append(vf, LengthRange(min, max))
	}
	if minNum := tag.Get("zpf_minnum"); "" != minNum {
		minInt, err := strconv.ParseInt(minNum, 10, 64)
		panicErr(err)
		vf = append(vf, NumberGT(minInt))
	}
	if maxNum := tag.Get("zpf_maxnum"); "" != maxNum {
		maxInt, err := strconv.ParseInt(maxNum, 10, 64)
		panicErr(err)
		vf = append(vf, NumberLT(maxInt))
	}
	if numRange := tag.Get("zpf_num"); "" != numRange {
		var min, max int64
		n, err := fmt.Sscanf(numRange, "%d %d", &min, &max)
		if 2 != n || nil != err {
			panic("invali zpf_num param: " + numRange)
		}
		vf = append(vf, NumberRange(min, max))
	}
	return vf
}

func getWidget(fieldType *reflect.StructField) string {
	if widget := fieldType.Tag.Get("zpf_widget"); "" != widget {
		return widget
	}
	return "textbox"
}

var fieldMetaMap map[string]*_FieldMeta
var fieldMetaLock sync.RWMutex

func getFieldMeta(formType reflect.Type, fieldType *reflect.StructField) (meta *_FieldMeta) {
	// metaKey := formType.PkgPath() + "/" + formType.Name() + "." + fieldType.Name
	metaKey := formType.String() + "|" + fieldType.Name
	fieldMetaLock.RLock()
	meta, exists := fieldMetaMap[metaKey]
	fieldMetaLock.RUnlock()
	if exists {
		return meta
	}
	meta = &_FieldMeta{}
	tag := &fieldType.Tag
	meta.FieldName = getFieldName(fieldType)
	meta.FieldFormat = tag.Get("zpf_format")
	meta.Validators = getValidators(tag)
	if label := tag.Get("zpf_label"); "" != label {
		meta.Label = label
	} else {
		meta.Label = fieldType.Name
	}
	meta.Widget = getWidget(fieldType)
	fieldMetaLock.Lock()
	fieldMetaMap[metaKey] = meta
	fieldMetaLock.Unlock()
	return meta
}

func getBindingByField(field reflect.Value, formType reflect.Type, fieldType *reflect.StructField) *_F {
	meta := getFieldMeta(formType, fieldType)
	if "-" == meta.FieldName {
		return nil
	}
	fieldVar := field.Addr().Interface()
	return NewF(fieldVar, meta.FieldName, meta.Label, meta.FieldFormat, meta.Validators...)
}

func GetBindings(form interface{}) []*_F {
	return GetReflectedBindings(reflect.ValueOf(form))
}

func GetReflectedBindings(formVal reflect.Value) []*_F {
	if reflect.Ptr != formVal.Kind() {
		panic("form must be a pointer")
	}
	formObj := formVal.Elem()
	if reflect.Struct != formObj.Kind() {
		panic("form must point to a struct")
	}
	formType := formObj.Type()
	fieldCnt := formType.NumField()

	bindings := make([]*_F, 0)
	for fieldIdx := 0; fieldIdx < fieldCnt; fieldIdx++ {
		field := formObj.Field(fieldIdx)
		fieldType := formType.Field(fieldIdx)
		binding := getBindingByField(field, formType, &fieldType)
		if nil != binding {
			bindings = append(bindings, binding)
		}
	}
	return bindings
}

func init() {
	fieldMetaMap = make(map[string]*_FieldMeta)
}
