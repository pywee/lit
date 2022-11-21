package function

import (
	"fmt"

	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// baseFunctions
// 支持的内置函数: 通用处理函数
var baseFunctions = []*FunctionInfo{
	{
		// println
		FunctionName: FUNCTION_PRINT,
		MustAmount:   1,
		MaxAmount:    -1,
		Args: []*functionArgs{
			{Type: types.INTERFACE, Must: true},
			{Type: types.INTERFACE, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			for _, v := range args {
				if v != nil {
					print(v.Lit, " ")
				}
			}
			println()
			return nil, nil
		},
	},
	{
		// varDump
		FunctionName: FUNCTION_VARDUMP,
		MustAmount:   1,
		MaxAmount:    -1,
		Args: []*functionArgs{
			{Type: types.INTERFACE, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			for _, v := range args {
				if v != nil {
					fmt.Printf("%s %s ", v.Tok, v.Lit)
				}
			}
			print("\n")
			return nil, nil
		},
	},
	{
		FunctionName: FUNCTION_ISBOOL,
		MustAmount:   1,
		MaxAmount:    1,
		Args:         []*functionArgs{{Type: types.INTERFACE, Must: true}},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			if args[0].Tok == "BOOL" {
				return &global.Structure{Tok: "BOOL", Lit: "true"}, nil
			}
			return &global.Structure{Tok: "BOOL", Lit: "false"}, nil

			// if s := args[0]; s.Tok == "BOOL" && (s.Lit == "true" || s.Lit == "false") {
			// 	return &global.Structure{Tok: "BOOL", Lit: "true"}, nil
			// }
			// return &global.Structure{Tok: "BOOL", Lit: "false"}, nil
		},
	},
}
