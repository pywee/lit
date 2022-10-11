package goExpr

import (
	"fmt"
	"strconv"
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
	FN func(...*structure) (*structure, error)
}

type functionArgAttr struct {
	// Must 是否必须
	Must bool
	// Type 参数类型
	Type string
}

const (
	TYPE_INTERFACE = "INTERFACE"
	TYPE_STRING    = "STRING"
	TYPE_INT       = "INT"
	TYPE_BOOL      = "BOOL"
)

const (
	FUNCTION_PRINT   = "print"
	FUNCTION_REPLACE = "replace"
)

// privateFunctions 内置函数
var privateFunctions = []*functionInfo{
	{
		FunctionName: FUNCTION_PRINT,
		MustAmount:   1,
		MaxAmount:    -1,
		Args: []*functionArgAttr{
			{Type: TYPE_INTERFACE, Must: true},
			{Type: TYPE_INTERFACE, Must: true},
		},
		FN: func(args ...*structure) (*structure, error) {
			for _, v := range args {
				fmt.Print(v.Lit)
			}
			fmt.Print("\n")
			return nil, nil
		},
	},
	{
		FunctionName: FUNCTION_REPLACE,
		MustAmount:   1,
		MaxAmount:    4,
		Args: []*functionArgAttr{
			{Type: TYPE_INTERFACE, Must: true},
			{Type: TYPE_STRING, Must: true},
			{Type: TYPE_STRING, Must: true},
			{Type: TYPE_INT, Must: true},
		},
		FN: func(args ...*structure) (*structure, error) {
			a0 := args[0].Lit
			a1 := args[1].Lit
			a2 := args[2].Lit
			a3, _ := strconv.Atoi(args[3].Lit)
			rx := strings.Replace(a0, a1, a2, a3)

			// FIXME
			return &structure{Tok: "IDENT", Lit: rx}, nil
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
