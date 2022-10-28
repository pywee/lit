package function

import (
	"github.com/pywee/lit/global"
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
	// 判断是否为类方法或普通函数
	funcName := expr[1].Tok
	// 普通函数处理
	if funcName != "IDENT" {
		return nil, types.ErrorNotFoundFunction
	}

	if !global.IsVariableOrFunction(expr[1]) {
		return nil, types.ErrorFunctionNameIrregular
	}

	argsDefinitions, err := getFunctionArgs(expr, pos)
	if err != nil {
		return nil, err
	}

	return &FunctionInfo{
		FunctionName: expr[1].Lit,
		CustFN:       argsDefinitions.fnBody,
		MustAmount:   argsDefinitions.needArgsAmount,
		MaxAmount:    argsDefinitions.maxArgsAmount,
		Args:         argsDefinitions.list,
	}, nil
}

type functionArgsInfo struct {
	// needArgsAmount 解析后得到必传参数的数量
	needArgsAmount int
	// maxArgsAmount 解析后得最大可传参数的数量
	maxArgsAmount int
	// list 解析后的函数形参定义数据
	list []*functionArgs
	// fnBody 函数体
	fnBody []*global.Structure
}

// getFunctionArgs 获取函数内的参数定义信息
func getFunctionArgs(expr []*global.Structure, pos string) (*functionArgsInfo, error) {
	var (
		err error
		// functionArgsInfo 解析后的函数形参定义数据
		argsInfo = new(functionArgsInfo)
		// argDefinition 解析后的函数形参定义数据
		argDefinition *functionArgs
		// foundBracket 标记发现了括号
		foundBracket = 0
		// foundBracketOnce 标记发现了括号
		foundBracketOnce = false
		// curlyBracket 标记发现了花括号
		curlyBracket = 0
		// arg 收集到的函数 expr 形式
		arg = make([]*global.Structure, 0, 5)
		// funcBody 函数体数据保存
		funcBody = make([]*global.Structure, 0, 20)
	)

	for _, v := range expr {
		if v.Tok == "(" && !foundBracketOnce {
			foundBracket++
			if foundBracket == 1 {
				continue
			}
		}
		if v.Tok == ")" {
			foundBracket--
		}

		// 参数定义位置结束
		if v.Tok == "{" {
			foundBracketOnce = true
			curlyBracket++
			if curlyBracket == 1 {
				continue
			}
		}
		if v.Tok == "}" {
			curlyBracket--
			if foundBracketOnce && curlyBracket == 0 {
				foundBracketOnce = false
				argsInfo.fnBody = funcBody
				break
			}
		}

		if curlyBracket > 0 {
			funcBody = append(funcBody, v)
			continue
		}

		if foundBracket > 0 {
			if foundBracket == 1 && v.Tok == "," {
				// 形参数据
				// a, a = 1, b = false
				if aLen := len(arg); aLen > 0 {
					if argDefinition, err = checkArguments(arg, aLen); err != nil {
						return nil, err
					}
					if argDefinition.Must {
						argsInfo.needArgsAmount++
					}
					argsInfo.maxArgsAmount++
					argsInfo.list = append(argsInfo.list, argDefinition)
					arg = nil
					continue
				}
			}
			arg = append(arg, v)
		}
	}

	if aLen := len(arg); aLen > 0 {
		if argDefinition, err = checkArguments(arg, aLen); err != nil {
			return nil, err
		}
		if argDefinition != nil {
			if argDefinition.Must {
				argsInfo.needArgsAmount++
			}
			argsInfo.maxArgsAmount++
			argsInfo.list = append(argsInfo.list, argDefinition)
		}
	}

	if curlyBracket != 0 {
		return nil, types.ErrorWrongSentence
	}
	return argsInfo, nil
}

// checkArguments 检查定义的函数的形参定义信息并返回合法数据
func checkArguments(arg []*global.Structure, argLen int) (*functionArgs, error) {
	if argLen == 1 {
		return &functionArgs{Type: types.INTERFACE, Must: true, Name: arg[0].Lit}, nil
	}
	if argLen == 3 {
		arg1 := arg[1]
		if arg1.Tok != "=" {
			return nil, types.ErrorFunctionArgsIrregular
		}
		arg2 := arg[2]
		return &functionArgs{
			Type:  types.INTERFACE,
			Must:  false,
			Name:  arg[0].Lit,
			Value: arg2.Lit,
		}, nil
	}
	if argLen > 0 {
		return nil, types.ErrorFunctionArgsIrregular
	}
	return nil, nil
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
