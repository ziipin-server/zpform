package zpform

import (
	"time"
)

type DateTS uint32

func (f *DateTS) FromString(value string) error {
	t, err := time.Parse("2006-01-02", value)
	if nil != err {
		return err
	}
	*f = DateTS(uint32(t.Unix()))
	return nil
}
