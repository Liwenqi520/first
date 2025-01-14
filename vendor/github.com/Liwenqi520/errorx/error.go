package errorx
import (
	"fmt"
	"runtime"
	"github.com/pkg/errors"
)
var _ error = (*Error)(nil)
// ErrEmpty 空错误
var ErrEmpty = errors.New("")
// CodePrefix 返回码前缀
var CodePrefix int64
// Error 自定义错误信息
type Error struct {
	Code     int64
	Msg      string //自定义翻译信息
	Position string
	Cause    error //存储错误信息堆栈
}
func (e *Error) Error() string {
	var cause string
	if e.Cause != ErrEmpty {
		cause = fmt.Sprintf(";err: %s", e.Cause.Error())
	}
	return fmt.Sprintf("code: %d ;position: %s %s", e.Code, e.Position, cause)
}
//New 创建异常附带自定义错误msg
func New(code int64, msg ...string) error {
	code = CodePrefix*1000000 + code
	e := &Error{
		Cause: ErrEmpty,
	}
	e.Code = code
	if len(msg) > 0 {
		e.Cause = errors.New(msg[0])
	}
	e.caller(1)
	return e
}

//NewWithMsg 返回自定义翻译msg
func NewWithMsg(code int64, msg string) error {
	code = CodePrefix*1000000 + code
	e := &Error{
		Cause: ErrEmpty,
	}
	e.Code = code
	e.Msg = msg
	e.caller(1)
	return e
}
//Caller 获取异常抛出的位置
func (e *Error) caller(skip int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		file = "???"
		line = 0
	}
	e.Position = fmt.Sprintf("%s:%d", file, line)
}