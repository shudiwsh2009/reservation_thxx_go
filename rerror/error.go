package rerror

import (
	"fmt"
	"github.com/mijia/sweb/log"
)

type RError struct {
	message string
	err     error
	code    int
	context []interface{}
}

func (e RError) Code() int {
	return e.code
}

func (e RError) Context() []interface{} {
	return e.context
}

func NewRError(m string, err error) *RError {
	ret := &RError{
		message: m,
		err:     err,
		code:    ErrorUnknown,
	}
	log.Error(ret.Error())
	return ret
}

func NewRErrorCode(m string, err error, code int) *RError {
	ret := &RError{
		message: m,
		err:     err,
		code:    code,
	}
	log.Error(ret.Error())
	return ret
}

func NewRErrorCodeContext(m string, err error, code int, context ...interface{}) *RError {
	ret := &RError{
		message: m,
		err:     err,
		code:    code,
		context: context,
	}
	log.Error(ret.Error())
	return ret
}

func (e *RError) Error() string {
	if e.err == nil {
		return fmt.Sprintf("%s", e.message)
	} else {
		return fmt.Sprintf("%s: %v", e.message, e.err)
	}
}

func (e *RError) DisplayMessage() string {
	if e.context == nil {
		return ReturnMessage(e.code)
	} else {
		return ReturnMessage(e.code, e.context...)
	}
}
