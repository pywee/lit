package goExpr

import (
	"fmt"
)

var (
	// ErrorNotFoundVariable 找不到变量
	ErrorNotFoundVariable = WithError(1001, "not found variable")
	// ErrorWrongSentence 语法错误
	ErrorWrongSentence = WithError(1002, "wrong sentence")
)

func WithError(code int, err string) error {
	return fmt.Errorf("found error (code %d), notice: %s\n", code, err)
}
