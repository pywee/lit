package function

import (
	"github.com/pywee/goExpr/global"
	"github.com/pywee/goExpr/types"
)

type CustomFunctions struct {
	cusFunctions map[string][]*global.Structure
}

func NewCustomFunctions() *CustomFunctions {
	return new(CustomFunctions)
}

// ParseCutFunc 解析函数数据
func (f *CustomFunctions) ParseCutFunc(expr []*global.Structure, pos string) error {
	// 判断是否为类方法或普通函数
	funcName := expr[1].Tok
	// 普通函数处理
	if funcName == "IDENT" {
		if !global.IsVariableOrFunction(expr[1]) {
			return types.ErrorFunctionNameIrregular
		}
		argsDefinitions, err := getFunctionArgsDefinitions(expr, pos)
		if err != nil {
			return err
		}

		funcDefinition := &functionInfo{
			FunctionName: funcName,
			MustAmount:   4,
			MaxAmount:    4,
			Args:         argsDefinitions,
		}

		global.Output(funcDefinition)
	}

	return nil
}

// getFunctionArgsDefinitions 获取函数内的参数定义信息
func getFunctionArgsDefinitions(expr []*global.Structure, pos string) ([]*functionArgs, error) {
	var (
		err error
		// argDefinition 解析后的函数形参定义数据
		argDefinition *functionArgs
		// argDefinitionList 解析后的函数形参定义数据
		argDefinitionList = make([]*functionArgs, 0, 5)
		// 检查当前函数是否为类方法
		foundBracket = 0
		// arg 收集到的函数 expr 形式
		arg = make([]*global.Structure, 0, 5)
		// needArgsAmount 解析后得到必传参数的数量
		needArgsAmount = 0
		// maxArgsAmount 解析后得最大可传参数的数量
		maxArgsAmount = 0
	)
	for _, v := range expr {
		if v.Tok == "(" {
			foundBracket++
			if foundBracket == 1 {
				continue
			}
		}
		if v.Tok == ")" {
			foundBracket--
		}
		if foundBracket > 0 {
			if foundBracket == 1 && v.Tok == "," {
				// 形参数据
				// a, a = 1, b = false
				aLen := len(arg)
				if aLen > 0 {
					if argDefinition, err = checkArguments(arg, aLen); err != nil {
						return nil, err
					}
					if argDefinition.Must {
						needArgsAmount++
					}
					maxArgsAmount++
					argDefinitionList = append(argDefinitionList, argDefinition)
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
	}
	if argDefinition != nil {
		argDefinitionList = append(argDefinitionList, argDefinition)
	}
	return argDefinitionList, nil
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

func (f *CustomFunctions) AddFunc(structName, funcName string, expr []*global.Structure) {
	f.cusFunctions[funcName] = expr
}

func (f *CustomFunctions) GetCustomeFunc(name string) []*global.Structure {
	if v, ok := f.cusFunctions[name]; ok {
		return v
	}
	return nil
}
