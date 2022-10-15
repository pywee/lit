package function

import (
	"fmt"

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
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			for _, v := range args {
				if v != nil {
					print(v.Lit, " ")
				}
			}
			print("\n")
			return nil, nil
		},
	},
	{
		// varDump
		FunctionName: FUNCTION_VARDUMP,
		MustAmount:   1,
		MaxAmount:    -1,
		Args: []*functionArgAttr{
			{Type: types.INTERFACE, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			print("varDump: ")
			for _, v := range args {
				if v != nil {
					fmt.Print(v, " ")
				}
			}
			print("\n")
			return nil, nil
		},
	},
}
