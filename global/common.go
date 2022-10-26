package global

import (
	"fmt"
	"regexp"
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
		fmt.Println(" ", x)
	}
	println("")
}

func Output2(expr [][]*Structure, k int) {
	for _, x := range expr {
		for _, v := range x {
			fmt.Println(k, "output:", v.Tok, v.Lit)
		}
	}
	println("")
}
