package types

import (
	"errors"
)

var (
	// ErrorWrongSentence 语法错误
	ErrorWrongSentence = errors.New("语法错误")
	// ErrorWrongVarOperation 语法错误,变量操作不合法
	ErrorWrongVarOperation = errors.New("语法错误,变量操作不合法")
	// ErrorNotFoundVariable 变量不存在，或者你的定义非法
	ErrorNotFoundVariable = errors.New("变量不存在，或者你的定义非法")
	// ErrorNotFoundFunction 找不到函数
	ErrorNotFoundFunction = errors.New("找不到函数")
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
	// ErrorHandleUnsupported 不支持的操作，类型不符
	ErrorHandleUnsupported = errors.New("不支持的操作，类型不符")
)

// 专门针对if语句声明如下错误提示
// if 语句需要做的检查特别多
var (
	// ErrorIfExpression if 语句定义不合法
	ErrorIfExpression = errors.New("if 语句定义不合法")
	// ErrorLogicOfIfExpression if 语句逻辑顺序有误
	ErrorLogicOfIfExpression = errors.New("if 语句逻辑顺序有误")
	// ErrorElseExpression else 语句定义不合法，不应该带有条件表达式
	ErrorElseExpression = errors.New("else 语句定义不合法，不应该带有条件表达式")
	// ErrorIfExpressionWithoutConditions if 语句定义不合法，缺少条件定义
	ErrorIfExpressionWithoutConditions = errors.New("if 语句定义不合法，缺少条件定义")
	// ErrorIlligleIfExpressionOfElse 不合理的if语句, 在一个完整的if句子内,else 关键词最多只应该出现一次
	ErrorIlligleIfExpressionOfElse = errors.New("if 语句定义不合法, 在一个完整的if句子内, else关键词最多只应该出现一次")
	// ErrorIlligleIfExpressionOfIf 不合理的if语句, 在一个完整的if句子内,else 关键词最多只应该出现一次
	ErrorIlligleIfExpressionOfIf = errors.New("不合理的if语句, 在一个完整的if句子内, 单独的if关键词最多只应该出现一次")
)

// for 流程控制语句错误处理
var (
	ErrorForExpression = errors.New("for 语法错误")
	ErrorForContinue   = errors.New("continue 关键字应在循环语句中出现")
	ErrorForBreak      = errors.New("break 关键字应在循环语句中出现")
)

func WithError(pos, err string) error {
	return errors.New(pos + err)
}
