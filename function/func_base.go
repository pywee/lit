package function

import (
	"github.com/pywee/goExpr/global"
	"github.com/pywee/goExpr/types"
)

// baseFunctions
// 支持的内置函数: 通用处理函数
var baseFunctions = []*functionInfo{
	{
		// print
		FunctionName: FUNCTION_PRINT,
		MustAmount:   1,
		MaxAmount:    -1,
		Args: []*functionArgAttr{
			{Type: types.INTERFACE, Must: true},
			{Type: types.INTERFACE, Must: true},
		},
		FN: func(args ...*global.Structure) (*global.Structure, error) {
			for _, v := range args {
				if v != nil {
					print(v.Lit, " ")
				}
			}
			print("\n")
			return nil, nil
		},
	},
}
