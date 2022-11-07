package lit

import (
	"strings"

	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// var (
// functions 内部全局方法
// functions = []string{"print", "printf", "sprintf"}
// privateVariable 自定义局部变量
// privateVariable = make(map[string]*structure, 10)
// publicVariable 自定义全局变量
// publicVariable = make(map[string]*structure, 10)
// )

type CodeInfomation struct {
	Name  string
	Type  string
	Value interface{}
}

type exprResult struct {
	Type  string
	Tok   string
	Value interface{}
}

func (r *Expression) Get(vName string) (*global.Structure, error) {
	ret, ok := r.publicVariable[vName]
	if !ok {
		return nil, types.ErrorNotFoundVariable
	}
	// fmt.Printf("get variable by name: %s, value: %v\n", vName, ret)
	return ret, nil
}

func (r *Expression) GetVal(vName string) interface{} {
	ret, ok := r.publicVariable[vName]
	if !ok {
		return types.ErrorNotFoundVariable
	}
	if len(ret.Lit) > 1 && (ret.Tok == "STRING" || ret.Tok == "CHAR") {
		return formatString(ret.Lit)
	}
	return ret.Lit
}

func formatString(s string) string {
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
