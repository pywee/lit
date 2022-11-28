package lit

import (
	"go/scanner"
	"go/token"
	"strings"

	fn "github.com/pywee/lit/function"
	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// 2022.11.13
// 重写核心代码 加快解析速度
// 完成 for 循环第一种循环形式
// 修复一些小 bug
// 加入新字段标记自定义函数是否有返回值

// cfn 代码体内的自定义函数
var cfn *fn.CustomFunctions

type expression struct {
	funcBlocks     []*fn.FunctionInfo
	codeBlocks     []*global.Block
	publicVariable map[string]*global.Structure
}

type exprResult struct {
	Type  string
	Tok   string
	Value interface{}
}

func NewExpr(src []byte) (*expression, error) {
	var (
		s      scanner.Scanner
		fset   = token.NewFileSet()
		expr   = make([]*global.Structure, 0, 100)
		result = new(expression)
	)

	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	negative := ""
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		stok := tok.String()
		if stok == "COMMENT" {
			continue
		}

		// 负数处理
		if stok == "-" {
			if eLen := len(expr); eLen > 0 && expr[eLen-1].Lit == "" {
				negative = "-"
				continue
			}
		}

		if negative == "-" {
			negative = ""
			lit = "-" + lit
		}

		posString := fset.Position(pos).String()
		posLine := strings.Split(posString, ":")[0]
		if sLit := strings.ToLower(lit); sLit == "false" || sLit == "true" {
			stok = "BOOL"
		}
		if stok == "STRING" || stok == "CHAR" {
			lit = global.FormatString(lit)
		}
		expr = append(expr, &global.Structure{
			Tok:      stok,
			Lit:      lit,
			Position: posLine,
		})
	}

	innerVar := make(map[string]*global.Structure, 5)
	_, err := result.initExpr(expr, innerVar, nil)
	// global.Output(innerVar["a"])
	return nil, err
}

type parsing struct {
	isInLoop bool
	isInFunc bool
}

// parseExprs 解析代码块
func (r *expression) parseExprs(expr []*global.Structure, innerVar global.InnerVar, runtime *parsing) (*expression, error) {
	var (
		err        error
		foundElse  bool
		rlen       = len(expr)
		blocks     = make([]*global.Block, 0, 20)
		funcBlocks = make([]*fn.FunctionInfo, 0, 20)
	)
	for i := 0; i < rlen; i++ {
		thisExpr := expr[i]
		if thisExpr.Tok == ";" && thisExpr.Lit == "\n" {
			continue
		}

		// FIXME.
		// 负数处理 bug

		// TODO
		// For
		// 数组
		// 对象
		// 内置函数完善
		// 抛出错误时 明确返回出现错误的行数

		// for 流程控制语句
		// FIXME 未针对for语句的合法性做充分检查
		if thisExpr.Tok == "for" {
			blocks, i, err = r.parseIdentedFOR(expr, blocks, innerVar, i)
			if err != nil {
				return nil, err
			}
			continue
		}

		// continue 语句
		if thisExpr.Tok == "continue" {
			if runtime == nil || !runtime.isInLoop {
				return nil, types.ErrorForContinue
			}
			blocks = append(blocks, &global.Block{Type: types.CodeTypeContinue})
			return &expression{codeBlocks: blocks, funcBlocks: funcBlocks}, nil
		}

		// break 语句
		if thisExpr.Tok == "break" {
			if runtime == nil || !runtime.isInLoop {
				return nil, types.ErrorForBreak
			}
			blocks = append(blocks, &global.Block{Type: types.CodeTypeBreak})
			return &expression{codeBlocks: blocks, funcBlocks: funcBlocks}, nil
		}

		// 变量声明
		if thisExpr.Tok == "IDENT" && i < rlen {
			tok := expr[i+1].Tok
			if global.InArrayString(tok, global.MathSym) {
				if blocks, i = parseIdentedVAR(r, blocks, expr, innerVar, tok, rlen, i); i == -1 {
					return nil, types.ErrorWrongVarOperation
				}
				continue
			}
		}

		// return 语句
		if thisExpr.Tok == "return" {
			var returnExpr = make([]*global.Structure, 0, 5)
			for j := i + 1; j < rlen; j++ {
				exprJ := expr[j]
				if exprJ.Tok == ";" {
					blocks = append(blocks, &global.Block{Type: types.CodeTypeIdentRETURN, Code: returnExpr})
					return &expression{funcBlocks: funcBlocks, codeBlocks: blocks}, nil
				}
				returnExpr = append(returnExpr, exprJ)
			}
			continue
		}

		// 函数定义
		if expr[i].Tok == "func" {
			if !global.IsVariableOrFunction(expr[i+1]) {
				return nil, types.ErrorFunctionNameIrregular
			}
			if funcBlocks, i, err = parseIdentFUNC(funcBlocks, expr, i, rlen); err != nil {
				return nil, err
			}
			continue
		}

		// 函数调用
		if expr[i].Tok == "IDENT" && i < rlen && expr[i+1].Tok == "(" {
			blocks, i = parseExecFUNC(blocks, expr, i, rlen)
			continue
		}

		// if 流程控制语句
		if expr[i].Tok == "if" && !foundElse {
			parsed, err := parseIdentedIF(blocks, expr, i, rlen)
			if err != nil {
				return nil, err
			}
			i = parsed.i
			blocks = parsed.blocks
			continue
		}

		// else 处理
		if expr[i].Tok == "else" {
			parsed, err := parseIdentELSE(blocks, expr, i, rlen)
			if err != nil {
				return nil, err
			}
			i = parsed.i
			blocks = parsed.blocks
			foundElse = parsed.foundElse

			// FIXME 不确定此处 continue 是否存在副作用
			continue
		}

		// 变量自增操作
		if thisExpr.Tok == "IDENT" && i < rlen && expr[i+1].Tok == "++" {
			blocks, i = parseIdentedVarPLUS(blocks, expr, i, rlen)
			continue
		}

		// 变量自减操作
		if thisExpr.Tok == "IDENT" && i < rlen && expr[i+1].Tok == "--" {
			blocks, i = parseIdentedVarREDUCE(blocks, expr, i, rlen)
			continue
		}

		// 此处如果返回 则 if 语句中针对 else 的解析会出问题
		return nil, types.ErrorWrongSentence
	}

	return &expression{codeBlocks: blocks, funcBlocks: funcBlocks}, nil
}

// initExpr 全局表达式入口
// 代码块中如果带有 if 等复杂语句 则需要从这里进入递归
func (r *expression) initExpr(expr []*global.Structure, innerVar global.InnerVar, runtime *parsing) (*global.Structure, error) {
	bs, err := r.parseExprs(expr, innerVar, runtime)
	if err != nil {
		return nil, err
	}

	if len(r.funcBlocks) == 0 {
		r.funcBlocks = bs.funcBlocks
	}

	for _, block := range bs.codeBlocks {
		if block.Type == types.CodeTypeContinue {
			return &global.Structure{Tok: "continue", Lit: "continue"}, nil
		}

		if block.Type == types.CodeTypeBreak {
			return &global.Structure{Tok: "break", Lit: "break"}, nil
		}

		if block.Type == types.CodeTypeIdentRETURN {
			return r.parse(block.Code, innerVar)
		}

		if block.Type == types.CodeTypeIdentVAR {
			if block.Name == "" {
				return nil, types.ErrorWrongSentence
			}
			rv, err := r.parse(block.Code, innerVar)
			if err != nil {
				return nil, err
			}
			if rv == nil {
				rv = &global.Structure{Tok: "NULL", Lit: "NULL"}
			}
			innerVar[block.Name] = rv
			// return &global.Structure{Tok: "BOOL", Lit: "true"}, nil
		} else if block.Type == types.CodeTypeFunctionExec {
			rv, err := r.parse(block.Code, innerVar)
			if err != nil {
				return nil, err
			}
			if rv != nil {
				return rv, nil
			}
			// return &global.Structure{Tok: "NULL", Lit: "NULL"}, nil
		} else if block.Type == types.CodeTypeIdentIF {
			// 检查if语句合法性
			if err = checkLegitimateIF(block.IfExt); err != nil {
				return nil, err
			}

			var conditionResult bool
			for _, v := range block.IfExt {
				if v.ConditionLen == 0 {
					conditionResult = true
				} else if v.ConditionLen > 0 {
					rv, err := r.parse(v.Condition, innerVar)
					if err != nil {
						return nil, err
					}
					conditionResult = global.ChangeToBool(rv)
				}
				if !conditionResult {
					continue
				}

				blen := len(v.Body)
				if blen < 2 {
					return nil, types.ErrorIfExpression
				}
				ret, err := r.initExpr(v.Body[1:blen-1], innerVar, runtime)
				if err != nil {
					return nil, err
				}
				if ret != nil {
					return ret, nil
				}
				break
			}
		} else if block.Type == types.CodeTypeVariablePlus {
			if len(block.Code) != 2 {
				return nil, types.ErrorWrongSentence
			}
			if _, err = execVarPlusReduce(block, innerVar, true); err != nil {
				return nil, err
			}
		} else if block.Type == types.CodeTypeVariableReduce {
			if len(block.Code) != 2 {
				return nil, types.ErrorWrongSentence
			}
			if _, err = execVarPlusReduce(block, innerVar, false); err != nil {
				return nil, err
			}
		} else if block.Type == types.CodeTypeIdentFOR {
			forExpr := block.ForExt
			if forExpr.Type == 1 {
				if err = r.execFORType1(forExpr, innerVar); err != nil {
					return nil, err
				}
			}
		}
	}
	return nil, nil
}

// parse 解析算术表达式、逻辑运算等
// 用于最小粒度处理
func (r *expression) parse(expr []*global.Structure, innerVar global.InnerVar) (*global.Structure, error) {
	exLen := len(expr)
	if exLen == 1 {
		// 变量解析
		rv := expr[0]
		if rv != nil && rv.Tok == "IDENT" {
			if innerRv, ok := innerVar[rv.Lit]; ok {
				return innerRv, nil
			}
			return nil, types.ErrorNotFoundVariable
		}
		return rv, nil
	}

	var innerExpr = make([]*global.Structure, 0, 10)
	for i := 0; i < exLen; i++ {
		if expr[i].Tok == "IDENT" && i+1 < exLen && expr[i+1].Tok == "(" {
			var (
				brCount  uint8
				funcExec = make([]*global.Structure, 0, 5)
			)
			for j := i; j < exLen; j++ {
				funcExec = append(funcExec, expr[j])
				if expr[j].Tok == "(" {
					brCount++
				}
				if expr[j].Tok == ")" {
					brCount--
					if brCount == 0 {
						rv, err := r.execFUNC(funcExec, expr[i+2:j], innerVar)
						if err != nil {
							return nil, err
						}
						if rv != nil {
							innerExpr = append(innerExpr, rv)
						}
						i = j
						break
					}
				}
			}
			continue
		}

		if expr[i].Tok == "(" {
			var (
				bracketCount uint8
				bracketExprs = make([]*global.Structure, 0, 10)
			)
			for j := i; j < exLen; j++ {
				exprJ := expr[j]
				bracketExprs = append(bracketExprs, exprJ)
				if exprJ.Tok == "(" {
					bracketCount++
				} else if exprJ.Tok == ")" {
					bracketCount--
					if bracketCount == 0 {
						rv, err := r.parse(bracketExprs[1:len(bracketExprs)-1], innerVar)
						if err != nil {
							return nil, err
						}
						i = j
						innerExpr = append(innerExpr, rv)
						break
					}
				}
			}
			continue
		}

		// 获取变量的值
		if expr[i].Tok == "IDENT" {
			if value, exists := innerVar[expr[i].Lit]; exists {
				innerExpr = append(innerExpr, value)
				continue
			}
			return nil, types.ErrorNotFoundVariable
		}

		// 逻辑运算 ||
		if expr[i].Tok == "||" {
			return r.parseOr(expr, innerExpr, innerVar, i)
		}

		// 逻辑运算 &&
		if expr[i].Tok == "&&" {
			return r.parseAnd(expr, innerExpr, innerVar, i)
		}
		innerExpr = append(innerExpr, expr[i])
	}

	if iLen := len(innerExpr); iLen >= 3 {
		for n := 0; n < iLen; n++ {
			// 比较运算
			// 考虑后期再支持 ===
			if tok := inArray(innerExpr[n].Tok, []string{"==", "!=", ">", "<", ">=", "<="}); tok != "" {
				rv, err := r.parseComparison(innerExpr[:n], innerExpr[n+1:], tok, innerVar)
				if err != nil {
					return nil, err
				}
				return rv, nil
			}
		}
	}

	if len(innerExpr) == 1 {
		return innerExpr[0], nil
	}

	// 数学计算
	innerExpr, err := r.parseExpr(innerExpr, "")
	if err != nil {
		return nil, err
	}
	if len(innerExpr) == 1 {
		return innerExpr[0], nil
	}
	return nil, nil
}

// getIdentedCustomFunc 获取定义好的自定义函数代码块
func (r *expression) getIdentedCustomFunc(fname string) *fn.FunctionInfo {
	for _, v := range r.funcBlocks {
		if v.FunctionName == fname {
			return v
		}
	}
	return nil
}

func inArray(sep string, arr []string) string {
	for _, v := range arr {
		if sep == v {
			return v
		}
	}
	return ""
}

func (r *expression) Get(vName string) (*global.Structure, error) {
	ret, ok := r.publicVariable[vName]
	if !ok {
		return nil, types.ErrorNotFoundVariable
	}
	return ret, nil
}

func (r *expression) GetVal(vName string) interface{} {
	ret, ok := r.publicVariable[vName]
	if !ok {
		return types.ErrorNotFoundVariable
	}
	if len(ret.Lit) > 1 && (ret.Tok == "STRING" || ret.Tok == "CHAR") {
		return global.FormatString(ret.Lit)
	}
	return ret.Lit
}
