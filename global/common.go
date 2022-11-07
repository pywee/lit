package global

import (
	"fmt"
	"regexp"

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

func Output(expr interface{}, x ...interface{}) {
	if fmt.Sprintf("%T", expr) == "[][]*global.Structure" {
		for _, v := range expr.([][]*Structure) {
			for _, vv := range v {
				fmt.Println("output from [][]arr:", vv)
			}
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

// FIXME
// ChangeToBool 将当前的输入转换为布尔值
func ChangeToBool(src *Structure) (*Structure, bool) {
	if src.Tok == "BOOL" {
		if src.Lit != "" && src.Lit != "false" {
			src.Lit = "true"
			return src, true
		}
		src.Lit = "false"
		return src, false
	}

	var returnBool bool
	if src.Tok == "STRING" && src.Lit != "" && src.Lit != "0" {
		src.Lit = "true"
		returnBool = true
	} else if src.Tok == "INT" && src.Lit != "0" {
		src.Lit = "true"
		returnBool = true
	} else if src.Tok == "FLOAT" && src.Lit != "0" {
		src.Lit = "true"
		returnBool = true
	} else {
		src.Lit = "false"
	}
	src.Tok = "BOOL"
	return src, returnBool
}

// ChangeBoolToInt 将布尔值转换为整型
func ChangeBoolToInt(src *Structure) error {
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

// ChangeTokTypeStringToTypeIntOrFloat 将字符串数字标记为整型
func ChangeTokTypeStringToTypeIntOrFloat(src *Structure) error {
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
