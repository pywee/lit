package function

import (
	"github.com/pywee/goExpr/global"
	"github.com/pywee/goExpr/types"
)

// 内置函数列表
var functions = make([]*functionInfo, 0, 100)

func init() {
	// 输入内置函数列表
	fns := [][]*functionInfo{strFunctions, numberFunctions, baseFunctions}
	for _, v := range fns {
		functions = append(functions, v...)
	}
}

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
	FN func(...*global.Structure) (*global.Structure, error)
}

type functionArgAttr struct {
	// Must 是否必须
	Must bool
	// Type 参数类型
	Type string
}

func IsExprFunction(expr []*global.Structure, rlen int) bool {
	if rlen < 3 {
		return false
	}
	if expr[0] == nil {
		return false
	}
	return expr[0].Tok == "IDENT" && expr[1].Tok == "(" && expr[rlen-1].Tok == ")"
}

func CheckFunctionName(name string) (*functionInfo, error) {
	for _, v := range functions {
		if v.FunctionName == name {
			return v, nil
		}
	}
	return nil, types.ErrorNotFoundFunction
}

// 获取传入执行函数的具体参数
func GetFunctionArgList(expr []*global.Structure) [][]*global.Structure {
	if len(expr) == 0 {
		return [][]*global.Structure{}
	}

	var list = make([][]*global.Structure, 0, 3)
	var arg = make([]*global.Structure, 0, 5)
	var seenK int8
	for _, v := range expr {
		if v.Tok == "(" {
			seenK++
		} else if v.Tok == ")" {
			seenK--
		}
		if v.Tok == "," {
			if len(arg) > 0 && seenK == 0 {
				list = append(list, arg)
				arg = nil
			} else {
				arg = append(arg, v)
			}
			continue
		}
		arg = append(arg, v)
	}
	list = append(list, arg)
	return list
}
