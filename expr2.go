package goExpr

import (
	fn "github.com/pywee/goExpr/function"
	"github.com/pywee/goExpr/global"
	"github.com/pywee/goExpr/types"
)

func test(expr []*global.Structure, pos string) {
	i := 0
	foundK := -1
	kList := make([]*global.Structure, 0, 10)
	for k, v := range expr {
		if v.Tok == "(" {
			if foundK == -1 {
				foundK = k
			}
			i++
		}
		if v.Tok == ")" {
			i--
		}
		if foundK >= 0 {
			kList = append(kList, v)
		}
		if i == 0 {
			global.Output(kList)
		}
	}
}

func (r *Expression) doFunc(expr []*global.Structure, pos string, foundAndOr bool) (*global.Structure, error) {
	// 判断是否表达式为函数
	// 如果表达式是 replace("1", "2", "", 1) 则可生效
	// FIXME 如果表达式是 replace("1", "2", "", 1) + "xxxx" 则不生效 fixed
	rLen := len(expr)
	if fn.IsExprFunction(expr, rLen) {
		if global.IsVariableOrFunction(expr[0]) {
			funcName := expr[0]
			fArgs, err := fn.CheckFunctionName(funcName.Lit)
			if err != nil {
				return nil, err
			}

			// 函数内参数检查
			// 获取传入执行函数的具体参数
			// 并将它们的结果值递归解析出来
			args := fn.GetFunctionArgList(expr[2 : rLen-1])
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
				rv, err := r.parseB(varg, pos, foundAndOr)
				if err != nil {
					return nil, err
				}

				// 检查最终解析出来的参数值类型是否与函数要求的形参类型一致
				if fa := fArgs.Args[k]; fa.Type != types.INTERFACE && fa.Type != rv.Tok {
					// TODO
					// 参数[弱类型]支持
					// 参数[提前在形参中设置默认值]支持
					// fmt.Println(fa.Type, rv.Tok)
					return nil, types.ErrorArgsNotSuitable
				}
				paramsList = append(paramsList, rv)
			}
			fRet, err := fArgs.FN(pos, paramsList...)
			return fRet, err
		}
	}
	return nil, nil
}
