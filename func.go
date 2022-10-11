package goExpr

import (
	"fmt"
	"strings"
)

type functionInfo struct {
	// FunctionName 名称
	FunctionName string
	// MustAmount 必填参数数量
	MustAmount int
	// 最大可传入参数 -1表示无限制
	MaxAmount int
	// List 形参信息定义
	Args []*functionArgAttr
	// FN 函数体
	FN func(...interface{}) (interface{}, error)
}

type functionArgAttr struct {
	// Must 是否必须
	Must bool
	// TypeName 参数类型
	TypeName string
}

const (
	TYPE_INTERFACE = "interface"
	TYPE_STRING    = "string"
	TYPE_INT       = "int"
	TYPE_BOOL      = "bool"
)

const (
	FUNCTION_PRINT   = "print"
	FUNCTION_REPLACE = "replace"
)

// privateFunctions 内置函数
var privateFunctions = []*functionInfo{
	{
		FunctionName: FUNCTION_PRINT,
		MustAmount:   2,
		MaxAmount:    -1,
		Args: []*functionArgAttr{
			{TypeName: TYPE_INTERFACE, Must: true},
			{TypeName: TYPE_INTERFACE, Must: true},
		},
		FN: func(args ...interface{}) (interface{}, error) {
			fmt.Println(args...)
			return nil, nil
		},
	},
	{
		FunctionName: FUNCTION_REPLACE,
		MustAmount:   4,
		MaxAmount:    4,
		Args: []*functionArgAttr{
			{TypeName: TYPE_INTERFACE, Must: true},
			{TypeName: TYPE_STRING, Must: true},
			{TypeName: TYPE_STRING, Must: true},
			{TypeName: TYPE_INT, Must: true},
		},
		FN: func(args ...interface{}) (interface{}, error) {
			a1 := args[0].(string)
			a2 := args[1].(string)
			a3 := args[2].(string)
			a4 := args[3].(int)
			return strings.Replace(a1, a2, a3, a4), nil
		},
	},
}

func checkFunctionName(name string) (*functionInfo, error) {
	for _, v := range privateFunctions {
		if v.FunctionName == name {
			return v, nil
		}
	}
	return nil, ErrorNotFoundFunction
}

// 获取传入执行函数的具体参数
func getFunctionArgList(expr []*structure) [][]*structure {
	if len(expr) == 0 {
		return [][]*structure{}
	}

	var list = make([][]*structure, 0, 3)
	var arg = make([]*structure, 0, 5)
	for _, v := range expr {
		if v.Tok == "," {
			if len(arg) > 0 {
				list = append(list, arg)
				arg = nil
			}
			continue
		}
		arg = append(arg, v)
	}
	list = append(list, arg)
	return list
}
