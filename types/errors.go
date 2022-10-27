package types

import (
	"errors"
)

var (
	// ErrorNotFoundVariable 找不到变量
	ErrorNotFoundVariable = errors.New("找不到变量")
	// ErrorWrongSentence 语法错误
	ErrorWrongSentence = errors.New("语法错误")
	// ErrorNotFoundFunction 找不到调用的函数
	ErrorNotFoundFunction = errors.New("找不到声明的函数")
	// ErrorArgsNotEnough 参数不足
	ErrorArgsNotEnough = errors.New("函数参数值数量不足")
	// ErrorTooManyArgs 参数过多
	ErrorTooManyArgs = errors.New("函数接收的参数过多")
	// ErrorArgsNotSuitable 参数类型不符
	ErrorArgsNotSuitable = errors.New("参数类型不符")
	// ErrorNonNumberic 非法字符参与数学计算
	ErrorNonNumberic = errors.New("a non-numeric value encountered has found")
	// ErrorFuncArgsAmountNotOK 函数参数不数量不符
	ErrorFuncArgsAmountNotOK = errors.New("函数参数不数量不符")
	// ErrorIdentType 源类型错误，当前类型不能转换为整型
	ErrorIdentType = errors.New("源类型错误，当前类型不能转换为整型")
	// ErrorStringIntCompared 不能将非数字类的字符串用于比较运算
	ErrorStringIntCompared = errors.New("不能将非数字类的字符串用于比较运算")
	// ErrorNotAllowToCompare 无法将这两种类型进行比较
	ErrorNotAllowToCompare = errors.New("无法将这两种类型进行比较")
	// ErrorFunctionIlligle 函数定义语法错误
	ErrorFunctionIlligle = errors.New("函数定义语法错误")
	// ErrorFunctionNameIrregular 函数名称不符合规范
	ErrorFunctionNameIrregular = errors.New("函数名称不符合规范")
	// ErrorFunctionArgsIrregular 函数参数定义不符合规范
	ErrorFunctionArgsIrregular = errors.New("函数参数定义不符合规范")
)

func WithError(pos, err string) error {
	return errors.New(pos + err)
}
