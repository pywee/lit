package global

import (
	"fmt"
	"regexp"
)

func IsNumber(s string) (bool, error) {
	return regexp.MatchString(`^[0-9]+[.]{0,1}[0-9]*$`, s)
}

func IsFloat(s string) (bool, error) {
	return regexp.MatchString(`^[0-9]+[.]{1}[0-9]+$`, s)
}

func IsInt(s string) (bool, error) {
	return regexp.MatchString(`^[0-9]+$`, s)
}

// IsVariableOrFunction 判断是否为标准变量和函数名称
func IsVariableOrFunction(name string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z_]{1}[a-zA-Z0-9_]{0,}$`, name)
	return match
}

func InArrayString(str string, arr []string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}
	return false
}

func Output(expr []*Structure) {
	for _, v := range expr {
		fmt.Println("output:", v)
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
