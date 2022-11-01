package function

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"

	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// strFunctions
// 支持的内置函数: 字符串处理函数
var strFunctions = []*FunctionInfo{
	{
		FunctionName: FUNCTION_REPLACE,
		MustAmount:   4,
		MaxAmount:    4,
		Args: []*functionArgs{
			{Type: types.INTERFACE, Must: true},
			{Type: types.STRING, Must: true},
			{Type: types.STRING, Must: true},
			{Type: types.INT, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			a0 := args[0].Lit
			a1 := args[1].Lit
			a2 := args[2].Lit
			a3, _ := strconv.Atoi(args[3].Lit)
			rx := strings.Replace(a0, a1, a2, a3)
			return &global.Structure{Tok: types.STRING, Lit: rx}, nil
		},
	},
	{
		// 必选参数1个
		// 可选参数1个 当不输入可选参数时 默认为空格 " "
		FunctionName: FUNCTION_TRIM,
		MustAmount:   1,
		MaxAmount:    2,
		Args: []*functionArgs{
			{Type: types.INTERFACE, Must: true},
			{Type: types.STRING, Must: false},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			if len(args) > 1 {
				return &global.Structure{Tok: types.STRING, Lit: strings.Trim(args[0].Lit, args[1].Lit)}, nil
			}
			return &global.Structure{Tok: types.STRING, Lit: strings.Trim(args[0].Lit, " ")}, nil
		},
	},
	{
		// 必选参数1个
		// 可选参数1个 当不输入可选参数时 默认为空格 " "
		FunctionName: FUNCTION_LTRIM,
		MustAmount:   1,
		MaxAmount:    2,
		Args: []*functionArgs{
			{Type: types.INTERFACE, Must: true},
			{Type: types.STRING, Must: false},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			if len(args) > 1 {
				return &global.Structure{Tok: types.STRING, Lit: strings.TrimLeft(args[0].Lit, args[1].Lit)}, nil
			}
			return &global.Structure{Tok: types.STRING, Lit: strings.TrimLeft(args[0].Lit, " ")}, nil
		},
	},
	{
		// 必选参数1个
		// 可选参数1个 当不输入可选参数时 默认为空格 " "
		FunctionName: FUNCTION_RTRIM,
		MustAmount:   1,
		MaxAmount:    2,
		Args: []*functionArgs{
			{Type: types.INTERFACE, Must: true},
			{Type: types.STRING, Must: false},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			if len(args) > 1 {
				return &global.Structure{Tok: types.STRING, Lit: strings.TrimRight(args[0].Lit, args[1].Lit)}, nil
			}
			return &global.Structure{Tok: types.STRING, Lit: strings.TrimRight(args[0].Lit, " ")}, nil
		},
	},
	{
		// 必选参数1个
		FunctionName: FUNCTION_TRIMSPACE,
		MustAmount:   1,
		MaxAmount:    1,
		Args:         []*functionArgs{{Type: types.INTERFACE, Must: true}},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			return &global.Structure{Tok: types.STRING, Lit: strings.TrimSpace(args[0].Lit)}, nil
		},
	},
	{
		FunctionName: FUNCTION_LEN,
		MustAmount:   1,
		MaxAmount:    1,
		Args:         []*functionArgs{{Type: types.INTERFACE, Must: true}},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			return &global.Structure{Tok: "INT", Lit: fmt.Sprintf("%d", len(args[0].Lit))}, nil
		},
	},
	{
		FunctionName: FUNCTION_UTF8LEN,
		MustAmount:   1,
		MaxAmount:    1,
		Args:         []*functionArgs{{Type: types.INTERFACE, Must: true}},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			return &global.Structure{Tok: "INT", Lit: fmt.Sprintf("%d", strings.Count(args[0].Lit, "")-1)}, nil
		},
	},
	{
		FunctionName: FUNCTION_SUBSTR,
		MustAmount:   3,
		MaxAmount:    3,
		Args: []*functionArgs{
			{Type: types.INTERFACE, Must: true},
			{Type: types.INT, Must: true},
			{Type: types.INT, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			print(strings.Count(args[0].Lit, "") - 1)
			return nil, nil
		},
	},
	{
		FunctionName: FUNCTION_MD5,
		MustAmount:   1,
		MaxAmount:    1,
		Args:         []*functionArgs{{Type: types.INTERFACE, Must: true}},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			return &global.Structure{
				Tok: "STRING",
				Lit: fmt.Sprintf("%x", md5.Sum([]byte(args[0].Lit))),
			}, nil
		},
	},
	{
		FunctionName: FUNCTION_CONTAINS,
		MustAmount:   2,
		MaxAmount:    2,
		Args: []*functionArgs{
			{Type: types.STRING, Must: true},
			{Type: types.STRING, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			return &global.Structure{
				Tok: "BOOL",
				Lit: fmt.Sprintf("%v", strings.Contains(args[0].Lit, args[1].Lit)),
			}, nil
		},
	},
	{
		FunctionName: FUNCTION_INDEX,
		MustAmount:   2,
		MaxAmount:    2,
		Args: []*functionArgs{
			{Type: types.STRING, Must: true},
			{Type: types.STRING, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			return &global.Structure{
				Tok: "INT",
				Lit: fmt.Sprintf("%d", strings.Index(args[0].Lit, args[1].Lit)),
			}, nil
		},
	},
	{
		FunctionName: FUNCTION_LASTINDEX,
		MustAmount:   2,
		MaxAmount:    2,
		Args: []*functionArgs{
			{Type: types.STRING, Must: true},
			{Type: types.STRING, Must: true},
		},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			return &global.Structure{
				Tok: "INT",
				Lit: fmt.Sprintf("%d", strings.LastIndex(args[0].Lit, args[1].Lit)),
			}, nil
		},
	},
	{
		FunctionName: FUNCTION_TOLOWER,
		MustAmount:   1,
		MaxAmount:    1,
		Args:         []*functionArgs{{Type: types.STRING, Must: true}},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			return &global.Structure{
				Tok: "STRING",
				Lit: strings.ToLower(args[0].Lit),
			}, nil
		},
	},
	{
		FunctionName: FUNCTION_TOUPPER,
		MustAmount:   1,
		MaxAmount:    1,
		Args:         []*functionArgs{{Type: types.STRING, Must: true}},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			return &global.Structure{
				Tok: "STRING",
				Lit: strings.ToUpper(args[0].Lit),
			}, nil
		},
	},
	{
		FunctionName: FUNCTION_TOTITLE,
		MustAmount:   1,
		MaxAmount:    1,
		Args:         []*functionArgs{{Type: types.STRING, Must: true}},
		FN: func(pos string, args ...*global.Structure) (*global.Structure, error) {
			return &global.Structure{
				Tok: "STRING",
				Lit: strings.ToTitle(args[0].Lit),
			}, nil
		},
	},
}
