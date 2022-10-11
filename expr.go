package goExpr

import (
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
	"regexp"
	"strconv"
	"strings"
)

type Expression struct {
	publicVariable map[string]*structure
}

func NewExpr(src []byte) (*Expression, error) {
	var result = &Expression{
		publicVariable: make(map[string]*structure, 10),
	}
	var s scanner.Scanner
	var list = make([]*structure, 0, 100)

	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		if tok.String() == ";" {
			var vName string
			if vleft, vLeftListEndIdx := findStrInfrontSymbool("=", list); vLeftListEndIdx != -1 {
				vName = vleft[0].Lit
				list = list[vLeftListEndIdx+1:]
			}

			// 递归解析表达式
			rv, err := result.parse(list, fset.Position(pos))
			if err != nil {
				return nil, errors.New(fset.Position(pos).String() + " error: " + err.Error())
			}

			// 变量赋值
			if vName != "" {
				// fmt.Printf("set %s to %v\n", vName, rv)
				result.publicVariable[vName] = rv
			}
			list = nil

			continue
		}

		tokString := tok.String()
		if tok.String() == "CHAR" {
			tokString = "STRING"
		}
		if tokString == "STRING" {
			lit = formatString(lit)
		}

		list = append(list, &structure{
			Position: fset.Position(pos).String(),
			Tok:      tokString,
			Lit:      lit,
		})

		// fmt.Printf("[ %s ]\t[ %s ]\t [ %s ] \n", fset.Position(pos).String(), tok, lit)
	}

	return result, nil
}

// parse 解析器主入口
func (r *Expression) parse(expr []*structure, pos token.Position) (*structure, error) {
	// 执行函数
	// print(a);
	rLen := len(expr)
	if rLen >= 3 && expr[0].Tok == "IDENT" && expr[1].Tok == "(" && expr[rLen-1].Tok == ")" {
		funcName := expr[0]
		if r.IsVariableOrFunction(funcName.Lit) {
			fArgs, err := checkFunctionName(funcName.Lit)
			if err != nil {
				return nil, err
			}

			// 获取传入执行函数的具体参数
			// 并将它们的结果值递归解析出来
			args := getFunctionArgList(expr[2 : rLen-1])
			argsLen := len(args)
			if fArgs.MustAmount > argsLen {
				return nil, ErrorArgsNotEnough
			}
			if fArgs.MaxAmount != -1 && fArgs.MaxAmount < argsLen {
				return nil, ErrorTooManyArgs
			}

			// FIXME
			// get params after parsing
			// 汇总解析成功之后的实参数据
			// 传入回调函数进行实际执行
			// 当前仅支持内置函数
			var paramsList []*structure
			for k, varg := range args {
				// FIXME
				// 函数中的实参表达式 实参可以是函数、变量、算术表达式等等
				rv, err := r.parse(varg, pos)
				if err != nil {
					return nil, err
				}

				// 检查最终解析出来的参数值类型是否与函数要求的形参类型一致
				if fa := fArgs.Args[k]; fa.Type != TYPE_INTERFACE && fa.Type != rv.Tok {
					// TODO
					// 参数弱类型支持
					// fmt.Println(fa.Type, rv.Tok)
					return nil, ErrorArgsNotSuitable
				}

				paramsList = append(paramsList, rv)
			}

			// output(paramsList)

			fRet, err := fArgs.FN(paramsList...)
			return fRet, err
		}
		return nil, nil
	}

	rv, err := r.parseExpr(expr, pos.String())
	if err != nil {
		return nil, err
	}

	rv0 := rv[0]
	return rv0, nil
}

// parseExpr 解析算术表达式入口
// 1+11|23-12;
// 1+2+(3%10+(11*22+1.1-10)+41+(5+9)+10)+1+(2-(10*2%2));
// 2+1/(20-1*(3+10)-1)%2^1;
func (r *Expression) parseExpr(src []*structure, pos string) ([]*structure, error) {
	v, err := r.findExprK(src, pos)
	if err != nil {
		return nil, err
	}
	return parsePlusReduceMulDivB(v, src[0].Position)
}

func (r *Expression) findExprK(expr []*structure, pos string) ([]*structure, error) {
	if len(expr) > 1 && (expr[0].Tok == "FLOAT" || expr[0].Tok == "INT") && expr[1].Tok == "IDENT" {
		return nil, ErrorWrongSentence
	}

	for k, v := range expr {
		if v.Tok == "IDENT" {
			if b := strings.ToLower(v.Lit); b == "true" || b == "false" {
				v.Tok = "BOOL"
				v.Lit = strings.ToLower(v.Lit)
				expr[k] = v
			} else if r.IsVariableOrFunction(v.Lit) {
				ret, err := r.Get(v.Lit)
				if err != nil {
					return nil, err
				}
				expr[k] = ret
			}
		}
	}

	var (
		startI   int
		endI     int
		startIdx = -1
		endIdx   = -1
	)
	for k, v := range expr {
		if v.Tok == "(" {
			startI++
			if startIdx == -1 {
				startIdx = k
			}
		}
		if v.Tok == ")" {
			endI++
			endIdx = k
		}
		if endI == startI && startIdx != -1 && endIdx != -1 {
			temp, err := r.findExprK(expr[startIdx+1:endIdx], pos)
			if err != nil {
				return nil, err
			}
			result := expr[:startIdx]
			result = append(result, temp...)
			return r.findExprK(append(result, expr[endIdx+1:]...), pos)
		}
	}

	return parsePlusReduceMulDivB(expr, pos)
}

// parsePlusReduceMulDivB 处理加减乘除
func parsePlusReduceMulDivB(arr []*structure, pos string) ([]*structure, error) {
	result, err := parsePlusReduceMulDiv([]string{"*", "/", "%"}, arr)
	if err != nil {
		return nil, err
	}

	// ?????
	if len(result) == 1 {
		// fmt.Println(result[0], arr[0])
		return result, nil
	}

	if len(result) > 0 {
		if result, err = parsePlusReduceMulDiv([]string{"+", "-", "|", "&", "^"}, result); err != nil {
			return nil, err
		}
	}

	if len(result) != 1 {
		return nil, ErrorWrongSentence
	}
	return result, nil
}

// parsePlusReduceMulDiv 递归解析表达式 最内层
func parsePlusReduceMulDiv(arr []string, expr []*structure) ([]*structure, error) {
	if len(expr) == 0 {
		return expr, nil
	}

	var result []*structure
	for k, v := range expr {
		if inArray(v.Tok, arr) && k-1 > -1 && k+2 <= len(expr) {
			r, err := findExprBetweenSymbool(expr[k-1], expr[k], expr[k+1])
			if err != nil {
				return nil, err
			}
			if r.Type == "STRING" {
				return []*structure{{Lit: r.Value.(string), Tok: r.Type}}, nil
			}

			middle := &structure{}
			if r.Type == "FLOAT" {

				// FIXME
				// 待观察是否对其他表达式有影响
				// example:
				// a = 2
				// b = 3
				// c = a + 100/10*2.0
				// output FLOAT 22.
				// 在使用 strings.TrimRight(lit, ".") 之后
				// 输出 FLOAT(22)

				middle.Tok = "FLOAT"
				middle.Lit = strings.TrimRight(fmt.Sprintf("%f", r.Value), "0")
				middle.Lit = strings.TrimRight(middle.Lit, ".")
			} else if r.Type == "INT" {
				middle.Tok = "INT"
				middle.Lit = fmt.Sprintf("%d", r.Value)
			}

			result = append(expr[:k-1], middle)
			result = append(result, expr[k+2:]...)
			return parsePlusReduceMulDiv(arr, result)
		}
	}
	return expr, nil
}

// findExprBetweenSymbool 查找某运算符数据
// 如: 1+2-3 输出查找 + 号 则返回 1+2
func findExprBetweenSymbool(l, m, r *structure) (*exprResult, error) {
	var (
		exprTyp string
		left    interface{}
		right   interface{}
		lTok    = l.Tok
		mTok    = m.Tok
		rTok    = r.Tok
	)

	// 字符串拼接操作 非算术计算
	if lTok == "STRING" && rTok == "STRING" {
		if mTok == "+" {
			return &exprResult{Type: "STRING", Value: l.Lit + r.Lit}, nil
		}
	}

	// 弱类型处理
	// 数字字符串转为整型或浮点型
	if lTok == "STRING" || lTok == "CHAR" {
		lit := l.Lit
		isFloat, err := regexp.MatchString(`^[0-9]+[.]+[0-9]*$`, lit)
		if err != nil {
			return nil, ErrorWrongSentence
		}
		if isFloat {
			if _, err = strconv.ParseFloat(lit, 64); err != nil {
				return nil, ErrorWrongSentence
			}
			lTok = "FLOAT"
			l.Lit = lit
		}

		isInt, err := regexp.MatchString(`^[0-9]+$`, lit)
		if err != nil {
			return nil, ErrorWrongSentence
		}
		if isInt {
			if _, err = strconv.ParseInt(lit, 10, 64); err != nil {
				return nil, ErrorWrongSentence
			}
			lTok = "INT"
			l.Lit = lit
		}
		if !isFloat && !isInt {
			return nil, ErrorWrongSentence
		}
	}
	if rTok == "STRING" || rTok == "CHAR" {
		_, err := strconv.ParseInt(r.Lit, 10, 64)
		if err != nil {
			return nil, ErrorWrongSentence
		}
		rTok = "INT"
	}

	// 弱类型处理
	// 布尔值转整型
	if l.Lit == "true" {
		lTok = "INT"
		l.Lit = "1"
	}
	if l.Lit == "false" {
		lTok = "INT"
		l.Lit = "0"
	}
	if r.Lit == "true" {
		rTok = "INT"
		r.Lit = "1"
	}
	if r.Lit == "false" {
		rTok = "INT"
		r.Lit = "0"
	}

	if lTok == "INT" && rTok == "INT" {
		exprTyp = "INT"
		left, _ = strconv.ParseInt(l.Lit, 10, 64)
		right, _ = strconv.ParseInt(r.Lit, 10, 64)
	}
	if lTok == "FLOAT" || rTok == "FLOAT" {
		exprTyp = "FLOAT"
		left, _ = strconv.ParseFloat(l.Lit, 64)
		right, _ = strconv.ParseFloat(r.Lit, 64)
	}

	if mTok == "*" {
		if exprTyp == "INT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) * right.(int64)}, nil
		}
		if exprTyp == "FLOAT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) * right.(float64)}, nil
		}
	}
	if mTok == "/" {
		if exprTyp == "INT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) / right.(int64)}, nil
		}
		if exprTyp == "FLOAT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) / right.(float64)}, nil
		}
	}
	if mTok == "+" {
		if exprTyp == "INT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) + right.(int64)}, nil
		}
		if exprTyp == "FLOAT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) + right.(float64)}, nil
		}
	}
	if mTok == "-" {
		if exprTyp == "INT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) - right.(int64)}, nil
		}
		if exprTyp == "FLOAT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) - right.(float64)}, nil
		}
	}
	if mTok == "%" && exprTyp == "INT" {
		return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) % right.(int64)}, nil
	}
	if mTok == "^" && exprTyp == "INT" {
		return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) ^ right.(int64)}, nil
	}
	if mTok == "|" && exprTyp == "INT" {
		return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) | right.(int64)}, nil
	}
	if mTok == "&" && exprTyp == "INT" {
		return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) & right.(int64)}, nil
	}
	return nil, ErrorWrongSentence
}

// findStrInfrontSymbool 查找指定符号之前的数据
func findStrInfrontSymbool(sym string, src []*structure) ([]*structure, int) {
	for k, v := range src {
		if sym == v.Tok {
			return src[:k], k
		}
	}
	return nil, -1
}

func inArray(sep string, arr []string) bool {
	for _, v := range arr {
		if sep == v {
			return true
		}
	}
	return false
}

func output(expr []*structure) {
	for _, v := range expr {
		fmt.Println("output:", v.Tok, v.Lit)
	}
	println("")
}
