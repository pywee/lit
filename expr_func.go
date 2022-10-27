package lit

import (
	fn "github.com/pywee/lit/function"
	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

func (r *Expression) execFunc(funcName string, expr []*global.Structure, pos string, innerVariable map[string]*global.Structure) (*global.Structure, error) {
	fArgs := fn.CheckFunctionName(funcName)
	if fArgs == nil {
		return nil, types.ErrorNotFoundFunction
	}

	// 函数内参数检查
	// 获取传入执行函数的具体参数
	// 并将它们的结果值递归解析出来
	args := fn.GetFunctionArgList(expr)
	argsLen := len(args)
	if fArgs.MustAmount > argsLen {
		return nil, types.ErrorArgsNotEnough
	}
	if fArgs.MaxAmount != -1 && fArgs.MaxAmount < argsLen {
		return nil, types.ErrorTooManyArgs
	}

	// FIXME
	// get params after parsing
	// 汇总解析成功之后的实参数据
	// 传入回调函数进行实际执行
	// 当前仅支持内置函数
	var paramsList []*global.Structure
	for k, varg := range args {
		// FIXME
		// 函数中的实参表达式 实参可以是函数、变量、算术表达式等等
		// global.Output(varg)
		rv, err := r.parse(varg, pos, innerVariable)
		if err != nil {
			return nil, err
		}

		if fArgs.MaxAmount != -1 {
			// 检查最终解析出来的参数值类型是否与函数要求的形参类型一致
			if fa := fArgs.Args[k]; fa.Type != types.INTERFACE && fa.Type != rv.Tok {
				// TODO
				// 参数[弱类型]支持
				// 参数[提前在形参中设置默认值]支持
				// fmt.Println(fa.Type, rv.Tok)
				return nil, types.ErrorArgsNotSuitable
			}
		}
		paramsList = append(paramsList, rv)
	}
	fRet, err := fArgs.FN(pos, paramsList...)
	if err != nil {
		return nil, err
	}
	return fRet, err
}

// execCustomFunc 解析并执行自定义函数
func (r *Expression) execCustomFunc(expr []*global.Structure, pos string) error {
	var (
		exprSingular  = make([]*global.Structure, 0, 3)
		innerVariable = make(map[string]*global.Structure)
	)
	for _, v := range expr {
		if v.Tok == ";" {
			cfnParsed := parseExprInnerFunc(exprSingular)
			if cfnParsed == nil {
				return types.ErrorWrongSentence
			}

			// 变量赋值
			if cfnParsed.typ == varStatemented {
				rv, err := r.parse(cfnParsed.varExpr, pos, nil)
				if err != nil {
					return err
				}
				innerVariable[cfnParsed.vName] = rv
			} else if cfnParsed.typ == funcImplemented {
				_, err := r.parse(cfnParsed.varExpr, pos, innerVariable)
				if err != nil {
					return err
				}
			}

			// global.Output(exprSingular)
			exprSingular = nil
			continue
		}
		exprSingular = append(exprSingular, v)
	}
	return nil
}

const (
	// varStatemented 变量声明或者赋值
	varStatemented = 1
	// funcImplemented 函数调用
	funcImplemented = 2
)

type innerFuncExpr struct {
	// typ 类型 1-表示变量赋值 2-函数调用
	typ int8
	// vName 变量名称 如果有的话
	vName string
	// tok 变量操作符 有可能是赋值时用的 = 或者是i++之类
	tok string
	// varExpr 表达式
	varExpr []*global.Structure
	// varParsed 表达式 varExpr 解析后的最终值
	// varParsed *global.Structure
}

func parseExprInnerFunc(expr []*global.Structure) *innerFuncExpr {
	sLen := len(expr)
	if sLen > 2 {
		expr0 := expr[0]
		expr1 := expr[1]

		// 变量声明或者境外赋值
		if expr1.Tok == "=" && expr[0].Tok == "IDENT" {
			return &innerFuncExpr{
				typ:     varStatemented,
				vName:   expr0.Lit,
				tok:     expr1.Tok,
				varExpr: expr[2:],
			}
		}

		// 函数调用
		if sLen >= 3 && expr0.Tok == "IDENT" && global.IsVariableOrFunction(expr0) && expr[sLen-1].Tok == ")" {
			return &innerFuncExpr{
				typ:     funcImplemented,
				vName:   expr0.Lit,
				varExpr: expr,
			}
		}
	}

	return nil
}
