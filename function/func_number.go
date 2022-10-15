package function

import "github.com/pywee/goExpr/global"

// numberFunctions
// 支持的内置函数: 数字处理函数
const (
	FUNCTION_INT        = "int"
	FUNCTION_FLOOR      = "floor"
	FUNCTION_STRING     = "string"
	FUNCTION_ISNUMBERIC = "isNumberic"
	FUNCTION_ISINT      = "isInt"
	FUNCTION_ISFLOAT    = "isFloat"
)

// numberFunctions
// 支持的内置函数: 数字处理函数
var numberFunctions = []*functionInfo{
	{
		// 检查当前输入的值是否为数字（包括整型和浮点型型）
		FunctionName: FUNCTION_ISNUMBERIC,
		MustAmount:   1,
		MaxAmount:    1,
		Args: []*functionArgAttr{
			{Type: TYPE_INTERFACE, Must: true},
		},
		FN: func(args ...*global.Structure) (*global.Structure, error) {
			match, err := global.IsNumber(args[0].Lit)
			if err != nil {
				return nil, err
			}
			if match {
				return &global.Structure{Tok: TYPE_BOOL, Lit: "true"}, nil
			}
			return &global.Structure{Tok: TYPE_BOOL, Lit: "false"}, nil
		},
	},
	{
		FunctionName: FUNCTION_ISINT,
		MustAmount:   1,
		MaxAmount:    1,
		Args: []*functionArgAttr{
			{Type: TYPE_INTERFACE, Must: true},
		},
		FN: func(args ...*global.Structure) (*global.Structure, error) {
			match, err := global.IsInt(args[0].Lit)
			if err != nil {
				return nil, err
			}
			if match {
				return &global.Structure{Tok: TYPE_BOOL, Lit: "true"}, nil
			}
			return &global.Structure{Tok: TYPE_BOOL, Lit: "false"}, nil
		},
	},
	{
		FunctionName: FUNCTION_ISFLOAT,
		MustAmount:   1,
		MaxAmount:    1,
		Args: []*functionArgAttr{
			{Type: TYPE_INTERFACE, Must: true},
		},
		FN: func(args ...*global.Structure) (*global.Structure, error) {
			match, err := global.IsFloat(args[0].Lit)
			if err != nil {
				return nil, err
			}
			if match {
				return &global.Structure{Tok: TYPE_BOOL, Lit: "true"}, nil
			}
			return &global.Structure{Tok: TYPE_BOOL, Lit: "false"}, nil
		},
	},
}
