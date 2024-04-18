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
	if e.msg == "" {
		e.msg = e.format
	}
	if e.reason == "" {
		return fmt.Sprintf("code:%s,\nmsg:%s", e.code, e.msg)
	}
	return fmt.Sprintf("code:%s,\nmsg: %s,\nreason:%s", e.code, e.msg, e.reason)
}

func (e ErrCode) Sprintf(msg ...any) ErrCode {
	e.msg = fmt.Sprintf(e.format, msg...)
	return e
}

func (e ErrCode) Code() string {
	return e.code
}

func ValidFormat(str string) bool {
	pattern := `^Err_[0-8]{9}[1-9x]$`
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return matched
}
