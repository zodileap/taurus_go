package err

import (
	"fmt"
	"regexp"
)

func New(code string, format string, reason string) ErrCode {
	return ErrCode{
		code:   code,
		format: format,
		reason: reason,
	}
}

type ErrCode struct {
	code   string
	format string
	msg    string
	reason string
}

func (e ErrCode) Error() string {
	if e.reason == "" {
		return fmt.Sprintf("code: %s, msg: %s", e.code, e.msg)
	}
	if e.msg == "" {
		e.msg = e.format
	}
	return fmt.Sprintf("code: %s, msg: %s, reason: %s", e.code, e.msg, e.reason)
}

func (e ErrCode) Sprintf(msg ...any) ErrCode {
	e.msg = fmt.Sprintf(e.format, msg...)
	return e
}

func ValidFormat(str string) bool {
	pattern := `^Err_[0-9]{10}$`
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return matched
}
