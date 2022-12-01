package function

import (
	"fmt"
	"strings"

	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// baseFunctions
// 支持的内置函数: 通用处理函数
var baseFunctions = []*FunctionInfo{
	{
		// print
		FunctionName: FUNCTION_PRINT,
		MustAmount:   1,
		MaxAmount:    -1,
		Args: []*functionArgs{
			{Type: types.INTERFACE, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			for _, v := range args {
				if v != nil {
					print(v.Lit)
				}
			}
			return nil, nil
		},
	},
	{
		// println
		FunctionName: FUNCTION_PRINTLN,
		MustAmount:   1,
		MaxAmount:    -1,
		Args: []*functionArgs{
			{Type: types.INTERFACE, Must: true},
			{Type: types.INTERFACE, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			for _, v := range args {
				if v == nil {
					continue
				}
				if len(v.Arr) > 0 {
					print("Array ")
				} else {
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
		},
	},
	{
		FunctionName: FUNCTION_ISINT,
		MustAmount:   1,
		MaxAmount:    1,
		Args:         []*functionArgs{{Type: types.INTERFACE, Must: true}},
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
		Args:         []*functionArgs{{Type: types.INTERFACE, Must: true}},
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
	{
		// 检查当前输入的值是否为数字（包括整型和浮点型型）
		FunctionName: FUNCTION_ISNUMERIC,
		MustAmount:   1,
		MaxAmount:    1,
		Args: []*functionArgs{
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
		FunctionName: FUNCTION_LEN,
		MustAmount:   1,
		MaxAmount:    2,
		Args: []*functionArgs{
			{Type: types.INTERFACE, Must: true},
			{Type: types.INTERFACE, Must: false},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			arg0 := args[0]
			if len(args) > 1 {
				if arg1 := args[1]; arg1.Tok == "BOOL" || arg1.Tok == "INT" {
					if arg1.Lit == "true" || arg1.Lit != "0" {
						return &global.Structure{Tok: "INT", Lit: fmt.Sprintf("%d", strings.Count(arg0.Lit, "")-1)}, nil
					}
				}
				return nil, types.ErrorFunctionArgsNotSuitable
			}
			return &global.Structure{Tok: "INT", Lit: fmt.Sprintf("%d", len(arg0.Lit))}, nil
		},
	},
}
