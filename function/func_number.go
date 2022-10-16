package function

import (
	"github.com/pywee/goExpr/global"
	"github.com/pywee/goExpr/types"
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
			{Type: types.INTERFACE, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			match, err := global.IsNumber(args[0].Lit)
			if err != nil {
				return nil, err
			}
			if match {
				return &global.Structure{Tok: types.BOOL, Lit: "true"}, nil
			}
			return &global.Structure{Tok: types.BOOL, Lit: "false"}, nil
		},
	},
	{
		FunctionName: FUNCTION_ISINT,
		MustAmount:   1,
		MaxAmount:    1,
		Args: []*functionArgAttr{
			{Type: types.INTERFACE, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			if len(args) == 0 {
				return nil, nil
			}
			match, err := global.IsInt(args[0].Lit)
			if err != nil {
				return nil, err
			}
			if match {
				return &global.Structure{Tok: types.BOOL, Lit: "true"}, nil
			}
			return &global.Structure{Tok: types.BOOL, Lit: "false"}, nil
		},
	},
	{
		FunctionName: FUNCTION_ISFLOAT,
		MustAmount:   1,
		MaxAmount:    1,
		Args: []*functionArgAttr{
			{Type: types.INTERFACE, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			match, err := global.IsFloat(args[0].Lit)
			if err != nil {
				return nil, err
			}
			if match {
				return &global.Structure{Tok: types.BOOL, Lit: "true"}, nil
			}
			return &global.Structure{Tok: types.BOOL, Lit: "false"}, nil
		},
	},
}
