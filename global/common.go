package global

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pywee/lit/types"
)

func IsNumber(s string) (bool, error) {
	return regexp.MatchString(`^[-0-9]+[.]{0,1}[0-9]*$`, s)
}

func IsFloat(s string) (bool, error) {
	return regexp.MatchString(`^[0-9]+[.]{1}[0-9]+$`, s)
}

func IsInt(s string) (bool, error) {
	return regexp.MatchString(`^[-0-9]+$`, s)
}

// IsVariableOrFunction 判断是否为标准变量和函数名称
func IsVariableOrFunction(expr *Structure) bool {
	if expr != nil && expr.Tok == "IDENT" {
		match, _ := regexp.MatchString(`^[a-zA-Z_]{1}[a-zA-Z0-9_]*$`, expr.Lit)
		return match
	}
	return false
}

func InArrayString(str string, arr []string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}
	return false
}

// FIXME
// TransformAllToBool 将当前的输入转换为布尔值
func TransformAllToBool(src *Structure) bool {
	if src == nil {
		return false
	}

	sTok := src.Tok
	sLit := src.Lit
	if sTok == "BOOL" {
		if sLit != "" && sLit != "false" {
			src.Lit = "true"
			return true
		}
		src.Lit = "false"
		return false
	}

	var returnBool bool
	if sTok == "STRING" && sLit != "" && sLit != "0" {
		src.Lit = "true"
		returnBool = true
	} else if sTok == "INT" && sLit != "0" {
		src.Lit = "true"
		returnBool = true
	} else if sTok == "FLOAT" && sLit != "0" {
		src.Lit = "true"
		returnBool = true
	} else {
		src.Lit = "false"
	}
	src.Tok = "BOOL"
	return returnBool
}

// TransformBoolToInt 将布尔值转换为整型
func TransformBoolToInt(src *Structure) error {
	src.Tok = "INT"
	if src.Lit == "false" {
		src.Lit = "0"
		return nil
	}
	if src.Lit == "true" {
		src.Lit = "1"
		return nil
	}
	return types.ErrorIdentType
}

// TODO
func TransformAllToInt(src *Structure) error {
	if src.Tok == "BOOL" {
		return TransformBoolToInt(src)
	}
	return nil
}

// TransformTokTypeStringToTypeIntOrFloat 将字符串数字标记为整型
func TransformTokTypeStringToTypeIntOrFloat(src *Structure) error {
	var (
		ok  bool
		err error
	)

	if src.Tok == "STRING" {
		if src.Lit == "" || src.Lit == "false" {
			src.Tok = "INT"
			src.Lit = "0"
			return nil
		}
		if src.Lit == "true" {
			src.Tok = "INT"
			src.Lit = "1"
			return nil
		}
	}
	if ok, err = IsInt(src.Lit); err != nil {
		return err
	}
	if ok {
		src.Tok = "INT"
		return nil
	}
	if ok, err = IsFloat(src.Lit); err != nil {
		return err
	}
	if ok {
		src.Tok = "FLOAT"
		return nil
	}
	return types.ErrorStringIntCompared
}

func Output(expr interface{}, x ...interface{}) {
	if fmt.Sprintf("%T", expr) == "[]*global.ExIf" {
		for _, v := range expr.([]*ExIf) {
			fmt.Println("output from []arr:", v.Condition, v.Body)
		}
	} else if fmt.Sprintf("%T", expr) == "[][]*global.Structure" {
		for _, v := range expr.([][]*Structure) {
			for _, vv := range v {
				fmt.Println("output from [][]arr:", vv)
			}
			println()
		}
	} else if fmt.Sprintf("%T", expr) == "[]*global.Structure" {
		for _, v := range expr.([]*Structure) {
			fmt.Println("output from []arr:", v)
		}
	} else if fmt.Sprintf("%T", expr) == "*global.Structure" {
		fmt.Println("result value output:", expr)
	} else {
		fmt.Println(expr)
	}
	if len(x) > 0 {
		Output(x[0])
	}
}

func Output2(expr [][]*Structure, k int) {
	for _, x := range expr {
		for _, v := range x {
			fmt.Println(k, "output:", v.Tok, v.Lit)
		}
	}
	println("")
}

type CodeInfomation struct {
	Name  string
	Type  string
	Value interface{}
}

func FormatString(s string) string {
	var (
		slen = len(s)
	)

	if s[0] == 34 && s[slen-1] == 34 {
		s = s[1 : slen-1]
		s = strings.Replace(s, `\"`, `"`, -1)
		s = strings.Replace(s, `\\`, `\`, -1)
	}
	return s

	// if s[0] == 39 && s[slen-1] == 39 { // 引号 '
	// 	// FIXME
	// 	lit = strings.TrimRight(s[1:], "'")
	// } else if s[0] == 34 && s[slen-1] == 34 { // 引号 "
	// 	lit = strings.TrimRight(s[1:], `"`)
	// 	lit = strings.TrimRight(lit, "\"")
	// 	lit = strings.Replace(lit, `\"`, `"`, -1)
	// 	lit = strings.Replace(lit, `\\`, `\`, -1)
	// }
	// return lit
}
