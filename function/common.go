package function

import (
	global "github.com/pywee/lit/global"
)

// 内置函数列表
var functions = make([]*FunctionInfo, 0, 20)

type FunctionInfo struct {
	// StructName 所属的结构体名称
	StructName string
	// FunctionName 名称
	FunctionName string
	// MustAmount 必填参数数量
	MustAmount int
	// 最大可传入参数 -1表示无限制
	MaxAmount int
	// List 形参信息定义
	Args []*functionArgs
	// FN 函数体
	FN func(string, ...*global.Structure) (*global.Structure, error)
	// CustFN 自定义函数体
	CustFN []*global.Structure
}

type functionArgs struct {
	// Must 是否必须
	Must bool
	// Type 参数类型
	Type string
	// Name 参数名
	Name string
	// Value 参数默认值
	Value string
}

func init() {
	// 输入内置函数列表
	fns := [][]*FunctionInfo{strFunctions, numberFunctions, baseFunctions}
	for _, v := range fns {
		functions = append(functions, v...)
	}
}

func IsExprFunction(expr []*global.Structure, rlen int) bool {
	if rlen < 3 {
		return false
	}
	if expr[0] == nil {
		return false
	}
	return expr[0].Tok == "IDENT" && expr[1].Tok == "(" && ((expr[rlen-1].Tok == ";" && expr[rlen-2].Tok == ")") || expr[rlen-1].Tok == ")")
}

// GetInnerIdentedFunc 查找内部函数
func GetInnerIdentedFunc(name string) *FunctionInfo {
	for _, v := range functions {
		if v.FunctionName == name {
			return v
		}
	}
	return nil
}

// 获取传入执行函数的具体参数
func GetFunctionArgList(expr []*global.Structure) [][]*global.Structure {
	if len(expr) == 0 {
		return [][]*global.Structure{}
	}

	var (
		seenK int8
		list  = make([][]*global.Structure, 0, 3)
		arg   = make([]*global.Structure, 0, 5)
	)
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

func BoolToInt(src *global.Structure) {}
