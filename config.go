package goExpr

import (
	"regexp"
)

var (
	// functions 内部全局方法
	// functions = []string{"print", "printf", "sprintf"}
	// privateVariable 自定义局部变量
	// privateVariable = make(map[string]*structure, 10)
	// publicVariable 自定义全局变量
	publicVariable = make(map[string]*structure, 10)
)

type CodeInfomation struct {
	Name  string
	Type  string
	Value interface{}
}

type structure struct {
	Position string
	Tok      string
	Lit      string
}

type exprResult struct {
	Type  string
	Tok   string
	Value interface{}
}

// IsVariableOrFunction 判断是否为标准变量和函数名称
func (r *Expression) IsVariableOrFunction(name string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z_]{1}[a-zA-Z0-9_]{0,}$`, name)
	return match
}

func (r *Expression) Get(vName string) (*structure, error) {
	ret, ok := r.publicVariable[vName]
	if !ok {
		return nil, ErrorNotFoundVariable
	}
	// fmt.Printf("get variable by name: %s, value: %v\n", vName, ret)
	return ret, nil
}
