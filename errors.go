package goExpr

import (
	"errors"
	"fmt"
)

var (
	// ErrorNotFoundVariable 找不到变量
	ErrorNotFoundVariable = errors.New("找不到变量")
	// ErrorWrongSentence 语法错误
	ErrorWrongSentence = errors.New("语法错误")
	// ErrorNotFoundFunction 找不到调用的函数
	ErrorNotFoundFunction = errors.New("call of undefined function")
	// ErrorArgsNotEnough 参数不足
	ErrorArgsNotEnough = errors.New("not enough args to input")
	// ErrorTooManyArgs 参数过多
	ErrorTooManyArgs = errors.New("too many args to input")
	// ErrorArgsNotSuitable 参数类型不符
	ErrorArgsNotSuitable = errors.New("参数类型不符")
	// ErrorNonNumberic 非法字符参与数学计算
	ErrorNonNumberic = errors.New("a non-numeric value encountered has found")
)

func WithError(pos, err string) error {
	return errors.New(pos + err)
}

func temp(arr []*structure) {
	for _, v := range arr {
		fmt.Println(v)
	}
}
