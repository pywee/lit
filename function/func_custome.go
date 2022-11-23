package function

import (
	global "github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

type CustomFunctions struct {
	cusFunctions map[string]*FunctionInfo
}

func NewCustomFunctions() *CustomFunctions {
	return &CustomFunctions{
		cusFunctions: make(map[string]*FunctionInfo, 0),
	}
}

// ParseCutFunc 解析函数数据
func (f *CustomFunctions) ParseCutFunc(expr []*global.Structure, pos string) (*FunctionInfo, error) {
	var (
		bracket    = 0
		bigBracket = 0
		expr1      = expr[1]
		exprLen    = len(expr)
		arg        = make([]*global.Structure, 0, 5)
		args       = make([][]*global.Structure, 0, 10)
		code       = make([]*global.Structure, 0, 20)
	)

	for i := 0; i < exprLen; i++ {
		if expr[i].Tok == "(" {
			bracket++
			for j := i + 1; j < exprLen; j++ {
				exprJ := expr[j]
				if exprJ.Tok == "," {
					args = append(args, arg)
					arg = nil
				}
				if exprJ.Tok == "(" {
					return nil, types.ErrorFunctionArgsIrregular
				}
				if exprJ.Tok == ")" {
					bracket--
					if bracket == 0 {
						i = j
						break
					}
				}
				if exprJ.Tok != "," {
					arg = append(arg, exprJ)
				}
			}
			continue
		}

		if expr[i].Tok == "{" {
			bigBracket++
			for j := i + 1; j < exprLen; j++ {
				exprJ := expr[j]
				if exprJ.Tok == "{" {
					bigBracket++
				} else if exprJ.Tok == "}" {
					bigBracket--
					if bigBracket == 0 {
						i = j
						break
					}
				}
				code = append(code, exprJ)
			}
		}
	}

	if len(arg) > 0 {
		args = append(args, arg)
	}

	// 检查参数定义合法性
	var (
		mustArgsAmount int
		optionalArgs   bool
		maxArgsAmount  = len(args)
		resultArgList  = make([]*functionArgs, 0, maxArgsAmount)
	)

	for _, v := range args {
		v0 := v[0]
		if v0.Tok != "IDENT" {
			return nil, types.ErrorFunctionArgsNotSuitable
		}

		vLen := len(v)
		if vLen != 1 && vLen != 3 {
			return nil, types.ErrorFunctionArgsIrregular
		}

		if vLen >= 3 {
			if v[1].Tok != "=" {
				return nil, types.ErrorFunctionArgsIrregular
			}
			v2 := v[2]
			if v2.Tok == "IDENT" {
				return nil, types.ErrorIrregularOfFuncArgValue
			}
			if !optionalArgs {
				optionalArgs = true
			}
			resultArgList = append(resultArgList, &functionArgs{
				Name:  v0.Lit,
				Type:  v2.Tok,
				Value: v2.Lit,
			})
		}

		if vLen == 1 && optionalArgs {
			return nil, types.ErrorWrongFuncArgsIdented
		}
		if !optionalArgs {
			mustArgsAmount++
			resultArgList = append(resultArgList, &functionArgs{
				Must: true,
				Type: v0.Tok,
				Name: v0.Lit,
			})
		}
	}

	return &FunctionInfo{
		FunctionName: expr1.Lit,
		CustFN:       code,
		MustAmount:   mustArgsAmount,
		MaxAmount:    maxArgsAmount,
		Args:         resultArgList,
	}, nil
}

func (f *CustomFunctions) AddFunc(structName string, fni *FunctionInfo) {
	f.cusFunctions[fni.FunctionName] = fni
}

func (f *CustomFunctions) GetCustomeFunc(name string) *FunctionInfo {
	if v, ok := f.cusFunctions[name]; ok {
		return v
	}
	return nil
}
