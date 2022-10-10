package goExpr

import (
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
			if err := result.parse(list, fset.Position(pos)); err != nil {
				// fmt.Printf("%s", err.Error())
				return nil, err
			}
			list = nil
			continue
		}

		list = append(list, &structure{
			Position: fset.Position(pos).String(),
			Tok:      tok.String(),
			Lit:      lit,
		})

		// fmt.Printf("[ %s ]\t[ %s ]\t [ %s ] \n", fset.Position(pos).String(), tok, lit)
	}

	// fmt.Println("reeee", result)

	return result, nil
}

// parse 解析器主入口
func (r *Expression) parse(expr []*structure, pos token.Position) error {
	var vName string
	// var result []*structure
	if vleft, vLeftListEndIdx := findStrInfrontSymbool("=", expr); vLeftListEndIdx != -1 {
		// result = expr[vLeftListEndIdx+1:]
		vName = vleft[0].Lit
	} else {
		// result = expr
	}

	// 执行函数
	// print(a);
	rLen := len(expr)
	if rLen >= 3 && expr[0].Tok == "IDENT" && expr[1].Tok == "(" && expr[rLen-1].Tok == ")" {
		funcName := expr[0]
		if r.IsVariableOrFunction(funcName.Lit) {
			exprT := expr[2 : rLen-1]
			fmt.Println(r.parseExpr(exprT, pos.String()))
			// fmt.Println(funcName.Lit)

			// if r.IsVariableOrFunction(varT.Tok) {
			// 	ret, ok := r.publicVariable[varT.Lit]
			// 	if !ok {
			// 		return ErrorNotFoundVariable
			// 	}
			// 	fmt.Println("变量被函数调用", ret)
			// }
			// output("<<<<", result[2:rLen-1])
		}
	} else {
		// 解析变量
		// output(":", expr)

		// FIXME 仅针对等于号右边是表达式的情况
		// 其余情况尚未处理
		rv, err := r.parseExpr(expr, pos.String())
		if err != nil {
			return err
		}
		if vName != "" {
			// fmt.Printf("set %s with value %v\n", vName, rv[0])
			r.publicVariable[vName] = rv[0]
		}
	}

	// temp(result)

	// FIXME

	return nil
}

// parseExpr 解析算术表达式入口
// 1+11|23-12;
// 1+2+(3%10+(11*22+1.1-10)+41+(5+9)+10)+1+(2-(10*2%2));
// 2+1/(20-1*(3+10)-1)%2^1;
func (r *Expression) parseExpr(src []*structure, pos string) ([]*structure, error) {
	var (
		// vLeftListEndIdx 等于号左边的数据结束位置
		vLeftListEndIdx int
		// vLeftList 等于号左边的数据 通常为变量的名称
		vLeftList []*structure
	)

	vLeftList, vLeftListEndIdx = findStrInfrontSymbool("=", src)
	if vLeftListEndIdx == -1 {
		return nil, ErrorWrongSentence
	}

	vName := vLeftList[0]
	if vName.Tok != "IDENT" {
		return nil, ErrorWrongSentence
	}

	v, err := r.findExprK(src[vLeftListEndIdx+1:], pos)
	if err != nil {
		return nil, err
	}
	return parsePlusReduceMulDivB(v, src[0].Position)
}

func (r *Expression) findExprK(expr []*structure, pos string) ([]*structure, error) {
	if len(expr) > 1 && (expr[0].Tok == "FLOAT" || expr[0].Tok == "INT") && expr[1].Tok == "IDENT" {
		return nil, ErrorWrongSentence
	}

	// for _, v := range expr {
	// 	fmt.Println(",,,,,", v)
	// }

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
	result := parsePlusReduceMulDiv([]string{"*", "/", "%"}, arr)
	if len(result) > 0 {
		result = parsePlusReduceMulDiv([]string{"+", "-", "|", "&", "^"}, result)
	}

	if len(result) != 1 {
		return nil, WithError(1002, fmt.Sprintf("%s wrong sentence", pos))
	}
	return result, nil
}

// parsePlusReduceMulDiv 递归解析表达式 最内层
func parsePlusReduceMulDiv(arr []string, expr []*structure) []*structure {
	if len(expr) == 0 {
		return expr
	}

	var result []*structure
	for k, v := range expr {
		if inArray(v.Tok, arr) && k-1 > -1 && k+2 <= len(expr) {
			if ok, r := findExprBetweenSymbool(expr[k-1], expr[k], expr[k+1]); ok {
				middle := &structure{}
				if r.Type == "FLOAT" {

					// FIXME 待观察是否对其他表达式有影响
					// example:
					// a = 2
					// b = 3
					// c = a + 100/10*2.0
					// output FLOAT 22.
					// after using expression strings.TrimRight(lit, ".")
					// then output FLOAT 22

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
			return nil
		}
	}

	return expr
}

// findExprBetweenSymbool 查找某运算符数据
// 如: 1+2-3 输出查找 + 号 则返回 1+2
func findExprBetweenSymbool(l, m, r *structure) (bool, *exprResult) {
	var (
		exprTyp string
		left    interface{}
		right   interface{}
		lTok    = l.Tok
		mTok    = m.Tok
		rTok    = r.Tok
	)

	// 弱类型处理
	// 数字字符串转为整型或浮点型
	if lTok == "STRING" || lTok == "CHAR" {
		lit := formatString(l.Lit)
		isFloat, err := regexp.MatchString(`^[0-9]+[.]+[0-9]*$`, lit)
		if err != nil {
			return false, nil
		}
		if isFloat {
			if _, err = strconv.ParseFloat(lit, 64); err != nil {
				return false, nil
			}
			lTok = "FLOAT"
			l.Lit = lit
		}

		isInt, err := regexp.MatchString(`^[0-9]+$`, lit)
		if err != nil {
			return false, nil
		}
		if isInt {
			if _, err = strconv.ParseInt(lit, 10, 64); err != nil {
				return false, nil
			}
			lTok = "INT"
			l.Lit = lit
		}
		if !isFloat && !isInt {
			return false, nil
		}
	}
	if rTok == "STRING" || rTok == "CHAR" {
		_, err := strconv.ParseInt(formatString(r.Lit), 10, 64)
		if err != nil {
			return false, nil
		}
		rTok = "INT"
		r.Lit = formatString(r.Lit)
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
			return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) * right.(int64)}
		}
		if exprTyp == "FLOAT" {
			return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) * right.(float64)}
		}
	}
	if mTok == "/" {
		if exprTyp == "INT" {
			return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) / right.(int64)}
		}
		if exprTyp == "FLOAT" {
			return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) / right.(float64)}
		}
	}
	if mTok == "+" {
		if exprTyp == "INT" {
			return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) + right.(int64)}
		}
		if exprTyp == "FLOAT" {
			return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) + right.(float64)}
		}
	}
	if mTok == "-" {
		if exprTyp == "INT" {
			return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) - right.(int64)}
		}
		if exprTyp == "FLOAT" {
			return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) - right.(float64)}
		}
	}
	if mTok == "%" && exprTyp == "INT" {
		return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) % right.(int64)}
	}
	if mTok == "^" && exprTyp == "INT" {
		return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) ^ right.(int64)}
	}
	if mTok == "|" && exprTyp == "INT" {
		return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) | right.(int64)}
	}
	if mTok == "&" && exprTyp == "INT" {
		return true, &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) & right.(int64)}
	}
	return false, nil
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

func output(sign string, expr []*structure) {
	for _, v := range expr {
		fmt.Println(sign, v.Tok, v.Lit)
	}
	println("")
}
