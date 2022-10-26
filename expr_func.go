package goExpr

import (
	fn "github.com/pywee/goExpr/function"
	"github.com/pywee/goExpr/global"
	"github.com/pywee/goExpr/types"
)

func (r *Expression) execFunc(funcName string, expr []*global.Structure, pos string) (*global.Structure, error) {
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
		rv, err := r.parse(varg, pos)
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
