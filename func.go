package lit

import (
	"errors"
	"go/scanner"
	"go/token"
	"strings"

	fn "github.com/pywee/lit/function"
	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// parseExecFUNC 解析自定义函数声明
func parseExecFUNC(blocks []*global.Block, expr []*global.Structure, i int, rlen int) ([]*global.Block, int) {
	block := &global.Block{Type: types.CodeTypeFunctionExec, Code: make([]*global.Structure, 0, 10)}
	for j := i; j < rlen; j++ {
		if expr[j].Tok == ";" {
			blocks = append(blocks, block)
			i = j
			break
		}
		block.Code = append(block.Code, expr[j])
	}
	return blocks, i
}

// parseIdentFUNC 解析函数定义
func parseIdentFUNC(funcBlocks []*fn.FunctionInfo, expr []*global.Structure, i int, rlen int) ([]*fn.FunctionInfo, int, error) {
	var (
		bracket uint16
		block   = &global.Block{Type: types.CodeTypeIdentFN, Code: make([]*global.Structure, 0, 30)}
	)
	for j := i; j < rlen; j++ {
		block.Code = append(block.Code, expr[j])
		if expr[j].Tok == "{" {
			bracket++
		} else if expr[j].Tok == "}" {
			bracket--
			if bracket == 0 {
				if len(block.Code) < 7 {
					return nil, 0, errors.New(expr[i].Position + types.ErrorFunctionIlligle.Error())
				}
				funcsParsed, err := cfn.ParseCutFunc(block.Code, expr[i].Position)
				if err != nil {
					return nil, 0, errors.New(expr[i].Position + err.Error())
				}
				funcBlocks = append(funcBlocks, funcsParsed)
				i = j
				break
			}
		}
	}
	return funcBlocks, i, nil
}

// execCustomFunc 执行自定义函数
// 当自定义的函数被调用时才会调用此方法
// realArgValues 为函数被调用时得到的实参
func (r *Expression) execCustomFunc(fni *fn.FunctionInfo, realArgValues []*global.Structure, pos string, innerVarInFuncParams map[string]*global.Structure) (*global.Structure, error) {
	// innerVar 函数体内的变量声明
	var innerVar = make(map[string]*global.Structure)

	// 以下场景需要在维护上下文 innerVarInFuncParams 局部变量
	for k, v := range innerVarInFuncParams {
		innerVar[k] = v
	}

	// 为形参赋值
	// 即: 将调用函数时传入的实参赋值给形参
	if err := r.setInnerVal(fni, realArgValues, pos, innerVar); err != nil {
		return nil, err
	}

	// 函数体代码解析
	// 递归回去 从头开始操作
	// FIXME.需要进一步检查参数上下文传递问题
	return r.createExpr(fni.CustFN, innerVar)
}

// execFunc 执行内置函数
func (r *Expression) execFunc(funcName string, expr []*global.Structure, pos string, innerVar map[string]*global.Structure) (*global.Structure, error) {
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
		rv, err := r.parse(varg, pos, innerVar)
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
	// 变量声明或者赋值
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
func (r *Expression) setInnerVal(fni *fn.FunctionInfo, realArgValues []*global.Structure, pos string, innerVar map[string]*global.Structure) error {
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
		realArgValueParsed, err := r.parse(thisArg, pos, innerVar)
		if err != nil {
			return err
		}
		nArg.Value = realArgValueParsed.Lit
		innerVar[nArg.Name] = realArgValueParsed
	}
	return nil
}

// getFunctionDefined 获取文本内所有自定义函数
func getFunctionDefined(s scanner.Scanner, file *token.File, fset *token.FileSet) error {
	var (
		funcKuo          int8
		foundCustomeFunc bool
		funcList         = make([]*global.Structure, 0, 10)
	)

	cfn = fn.NewCustomFunctions()
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		stok := tok.String()
		posString := fset.Position(pos).String()
		posLine := "第" + strings.Split(posString, ":")[0] + "行, "
		if stok == "func" {
			foundCustomeFunc = true
		}
		if foundCustomeFunc {
			if stok == "CHAR" || stok == "STRING" {
				lit = formatString(lit)
			}
			if sLit := strings.ToLower(lit); stok != "STRING" && (sLit == "false" || sLit == "true") {
				lit = sLit
				stok = "BOOL"
			}

			// FIXME 此处去掉符号 ; 可能导致其他逻辑有问题
			// 需要进一步观察测试
			// 去掉 go 语言解析包多余的分割标识符
			// if stok == ";" && lit == "\n" {
			// lit = ""
			// continue
			// }

			funcList = append(funcList, &global.Structure{
				Position: fset.Position(pos).String(),
				Tok:      stok,
				Lit:      lit,
			})

			if stok == "{" {
				funcKuo++
			} else if stok == "}" {
				funcKuo--
				if funcKuo == 0 {
					if len(funcList) < 7 {
						return errors.New(posLine + types.ErrorFunctionIlligle.Error())
					}
					funcsParsed, err := cfn.ParseCutFunc(funcList, posString)
					if err != nil {
						return errors.New(posLine + err.Error())
					}
					funcList = nil
					foundCustomeFunc = false
					cfn.AddFunc("", funcsParsed)
				}
			}

			// if tok.String() == ";" && lit == "\n" {
			// 	if len(funcList) < 7 {
			// 		return nil, errors.New(posLine + types.ErrorFunctionIlligle.Error())
			// 	}
			// 	funcsParsed, err := cfn.ParseCutFunc(funcList, posString)
			// 	if err != nil {
			// 		return nil, errors.New(posLine + err.Error())
			// 	}
			// 	funcList = nil
			// 	foundCustomeFunc = false
			// 	cfn.AddFunc("", funcsParsed)
			// }

			continue
		}
	}
	return nil
}
