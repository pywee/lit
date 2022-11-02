package lit

import (
	"strings"

	fn "github.com/pywee/lit/function"
	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// execFunc 执行内置函数
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

// execCustomFunc 执行自定义函数
// 当自定义的函数被调用时才会调用此方法
// realArgValues 为函数被调用时得到的实参
func (r *Expression) execCustomFunc(fni *fn.FunctionInfo, realArgValues []*global.Structure, pos string) (*global.Structure, error) {
	var (
		// exprSingularLine 函数体内的每一句表达式
		exprSingularLine = make([]*global.Structure, 0, 3)
		// innerVariable 函数体内的变量声明
		innerVariable = make(map[string]*global.Structure)
	)

	// 函数体代码解析
	for _, v := range fni.CustFN {
		if v.Tok == ";" {
			// 获得当前代码行的类型
			cfnInnerLineParsed := parseExprInnerFunc(exprSingularLine)
			if cfnInnerLineParsed == nil {
				return nil, types.ErrorWrongSentence
			}

			// 函数体内 return 语句
			if cfnInnerLineParsed.typ == returnIdent {
				if err := r.setInnerVal(fni, realArgValues, pos, innerVariable); err != nil {
					return nil, err
				}
				return r.parse(cfnInnerLineParsed.varExpr, pos, innerVariable)
			}

			// 变量赋值
			if cfnInnerLineParsed.typ == varStatemented {
				rv, err := r.parse(cfnInnerLineParsed.varExpr, pos, innerVariable)
				if err != nil {
					return nil, err
				}
				innerVariable[cfnInnerLineParsed.vName] = rv
			} else if cfnInnerLineParsed.typ == funcImplemented {

				// 对比函数形参和实参
				// realArgList := fn.GetFunctionArgList(realArgValues)
				// rLen := len(realArgList)
				// for k, nArg := range fni.Args {
				// 	var thisArg = make([]*global.Structure, 0, 3)
				// 	if k+1 > rLen {
				// 		if nArg.Must {
				// 			return nil, types.ErrorArgsNotEnough
				// 		}
				// 		thisArg = []*global.Structure{
				// 			{Tok: nArg.Type, Lit: nArg.Value},
				// 		}
				// 	} else {
				// 		thisArg = realArgList[k]
				// 	}

				// 	// 解析传入的实参 因为实参可能也是函数
				// 	realArgValueParsed, err := r.parse(thisArg, pos, innerVariable)
				// 	if err != nil {
				// 		return nil, err
				// 	}
				// 	nArg.Value = realArgValueParsed.Lit
				// 	innerVariable[nArg.Name] = realArgValueParsed
				// }

				var err error
				if err = r.setInnerVal(fni, realArgValues, pos, innerVariable); err != nil {
					return nil, err
				}
				if _, err = r.parse(cfnInnerLineParsed.varExpr, pos, innerVariable); err != nil {
					return nil, err
				}
			}
			exprSingularLine = nil
			continue
		}
		exprSingularLine = append(exprSingularLine, v)
	}
	return &global.Structure{Tok: "INT", Lit: "1"}, nil
}

const (
	// varStatemented 变量声明或者赋值
	varStatemented = 1
	// funcImplemented 函数调用
	funcImplemented = 2
	// returnIdent return 语句
	returnIdent = 3
)

type innerFuncExpr struct {
	// typ 类型 1-表示变量赋值 2-函数调用 3-return句子
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

// parseExprInnerFunc 解析当前函数体内的某一行代码
func parseExprInnerFunc(expr []*global.Structure) *innerFuncExpr {
	sLen := len(expr)
	if sLen < 2 {
		return nil
	}

	expr0 := expr[0]
	expr1 := expr[1]
	if strings.ToLower(expr0.Lit) == "return" {
		return &innerFuncExpr{
			typ:     returnIdent,
			vName:   "return",
			tok:     "return",
			varExpr: expr[1:],
		}
	}

	// sLen > 2
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

	return nil
}

// setInnerVal 解析函数体内的变量声明句子
// for example:
// a = 1;
// b = abc();
// c = a + b;
// return a+b+c
func (r *Expression) setInnerVal(fni *fn.FunctionInfo, realArgValues []*global.Structure, pos string, innerVariable map[string]*global.Structure) error {
	// 对比函数形参和实参
	realArgList := fn.GetFunctionArgList(realArgValues)
	rLen := len(realArgList)
	for k, nArg := range fni.Args {
		var thisArg = make([]*global.Structure, 0, 3)
		if k+1 > rLen {
			if nArg.Must {
				return types.ErrorArgsNotEnough
			}
			thisArg = []*global.Structure{{Tok: nArg.Type, Lit: nArg.Value}}
		} else {
			thisArg = realArgList[k]
		}

		// 解析传入的实参 因为实参可能也是函数
		realArgValueParsed, err := r.parse(thisArg, pos, innerVariable)
		if err != nil {
			return err
		}
		nArg.Value = realArgValueParsed.Lit
		innerVariable[nArg.Name] = realArgValueParsed
	}
	return nil
}
