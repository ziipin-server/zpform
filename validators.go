package zpform

import (
	"fmt"
	"regexp"
	"strconv"
)

func Required() ValidateFunc {
	return func(val string) (bool, string) {
		if len(val) > 0 {
			return true, ""
		}
		return false, "必填"
	}
}

func LengthRange(min, max int) ValidateFunc {
	return func(val string) (bool, string) {
		if len(val) >= min && len(val) <= max {
			return true, ""
		}
		return false, fmt.Sprintf("文字数量超过允许范围[%v, %v]", min, max)
	}
}

func LengthLT(max int) ValidateFunc {
	return func(val string) (bool, string) {
		if len(val) <= max {
			return true, ""
		}
		return false, fmt.Sprintf("文字过多（最多%v个字）", max)
	}
}

func LengthGT(min int) ValidateFunc {
	return func(val string) (bool, string) {
		if len(val) >= min {
			return true, ""
		}
		return false, fmt.Sprintf("文字过少（最少%v个字）", min)
	}
}

func NumberRange(min, max int64) ValidateFunc {
	return func(val string) (bool, string) {
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return false, "不是数字"
		}
		if intVal >= int64(min) && intVal <= int64(max) {
			return true, ""
		}
		return false, fmt.Sprintf("数值超过允许范围[%s, %s]", min, max)
	}
}

func NumberGT(min int64) ValidateFunc {
	return func(val string) (bool, string) {
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return false, "不是数字"
		}
		if intVal >= min {
			return true, ""
		}
		return false, fmt.Sprintf("数值过小（最小%v）", min)
	}
}

func NumberLT(max int64) ValidateFunc {
	return func(val string) (bool, string) {
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return false, "不是数字"
		}
		if intVal <= max {
			return true, ""
		}
		return false, fmt.Sprintf("数值过大（最大%v）", max)
	}
}

func Regexp(re string) ValidateFunc {
	return func(val string) (bool, string) {
		matched, err := regexp.MatchString(re, val)
		if err != nil {
			return false, "校验失败"
		}
		if matched {
			return true, ""
		}
		return false, "格式不符合要求"
	}
}
