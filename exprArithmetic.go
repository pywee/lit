package lit

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// parseExpr 解析算术表达式入口
// 1+11|23-12;
// 1+2+(3%10+(11*22+1.1-10)+41+(5+9)+10)+1+(2-(10*2%2));
// 2+1/(20-1*(3+10)-1)%2^1;
func (r *expression) parseExpr(expr []*global.Structure, pos string) ([]*global.Structure, error) {
	v, err := r.findExprK(expr, pos)
	if err != nil {
		return nil, err
	}
	return parsePlusReduceMulDivB(v, pos)
}

func (r *expression) findExprK(expr []*global.Structure, pos string) ([]*global.Structure, error) {
	exprLen := len(expr)

	// 如果返回的是最终值 则不再需要进一步解析了
	// 加快处理速度
	if exprLen == 1 {
		if e0 := expr[0]; inArray(e0.Tok, []string{"INT", "STRING", "FLOAT", "BOOL"}) != "" {
			return expr, nil
		}
	}

	// 错误表达式处理
	if exprLen > 1 {
		if e0Tok := expr[0].Tok; (e0Tok == "FLOAT" || e0Tok == "INT") && expr[1].Tok == "IDENT" {
			return nil, types.ErrorWrongSentence
		}
	}

	for k, v := range expr {
		if global.IsVariableOrFunction(v) {
			ret, err := r.Get(v.Lit)
			if err != nil {
				return nil, err
			}
			expr[k] = ret
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
func parsePlusReduceMulDivB(arr []*global.Structure, pos string) ([]*global.Structure, error) {
	result, err := parsePlusReduceMulDiv([]string{"*", "/", "%"}, arr)
	if err != nil {
		return nil, err
	}
	if len(result) == 1 {
		return result, nil
	}
	if result, err = parsePlusReduceMulDiv([]string{"&"}, result); err != nil {
		return nil, err
	}
	if len(result) == 1 {
		return result, nil
	}
	if len(result) > 0 {
		if result, err = parsePlusReduceMulDiv([]string{"+", "-", "|", "^"}, result); err != nil {
			return nil, err
		}
	}
	// if len(result) != 1 {
	// 	return nil, types.ErrorWrongSentence
	// }
	return result, nil
}

// parsePlusReduceMulDiv 递归解析表达式 最内层
func parsePlusReduceMulDiv(arr []string, expr []*global.Structure) ([]*global.Structure, error) {
	if len(expr) == 0 {
		return expr, nil
	}

	var result []*global.Structure
	for k, v := range expr {
		if inArray(v.Tok, arr) != "" && k-1 > -1 && k+2 <= len(expr) {
			r, err := findExprBetweenSymbool(expr[k-1], expr[k], expr[k+1])
			if err != nil {
				return nil, err
			}
			if r.Type == "STRING" {
				result = append([]*global.Structure{{Tok: "STRING", Lit: r.Value.(string)}}, expr[k+2:]...)
				return parsePlusReduceMulDiv(arr, result)
			}

			middle := &global.Structure{}
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
func findExprBetweenSymbool(l, m, r *global.Structure) (*exprResult, error) {
	var (
		exprTyp string
		left    interface{}
		right   interface{}
		lTok    = l.Tok
		mTok    = m.Tok
		rTok    = r.Tok
	)

	// 弱类型处理
	// 布尔值转整型
	if l.Tok == "BOOL" {
		if l.Lit == "true" {
			lTok = "INT"
			l.Tok = "INT"
			l.Lit = "1"
		} else if l.Lit == "false" {
			lTok = "INT"
			l.Tok = "INT"
			l.Lit = "0"
		}
	}
	if r.Tok == "BOOL" {
		if r.Lit == "true" {
			rTok = "INT"
			r.Tok = "INT"
			r.Lit = "1"
		} else if r.Lit == "false" {
			rTok = "INT"
			r.Tok = "INT"
			r.Lit = "0"
		}
	}

	// 字符串拼接及弱类型处理的算术计算
	if (lTok == "STRING" || lTok == "INTERFACE") && (rTok == "STRING" || rTok == "INTERFACE") {
		// 弱类型处理 如果左右两边都是字符串数字则允许进行算术计算
		isLeftNumeric, err := global.IsNumber(l.Lit)
		if err != nil {
			log.Print(err)
		}
		isRightNumeric, err := global.IsNumber(r.Lit)
		if err != nil {
			log.Print(err)
		}

		// 不是完全由数字构成的字符串则直接拼接
		// 否则转换类型进行算术计算
		// example A
		// a = "1";
		// b = "2";
		// print(a+b) 输出 (int)3

		// example B
		// a = "1";
		// b = "2a";
		// print(a+b) 输出 (string)12a

		if !isLeftNumeric || !isRightNumeric {
			if mTok == "+" {
				return &exprResult{Type: "STRING", Value: l.Lit + r.Lit}, nil
			}
			return nil, types.ErrorNonNumberic
		}
	}

	// 弱类型处理
	// 数字字符串转为整型或浮点型
	if lTok == "STRING" || lTok == "CHAR" {
		lit := l.Lit
		isFloat, err := global.IsFloat(lit)
		if err != nil {
			return nil, types.ErrorWrongSentence
		}
		if isFloat {
			if _, err = strconv.ParseFloat(lit, 64); err != nil {
				return nil, types.ErrorWrongSentence
			}
			lTok = "FLOAT"
			l.Lit = lit
		}

		isInt, err := global.IsInt(lit)
		if err != nil {
			return nil, types.ErrorWrongSentence
		}
		if isInt {
			if _, err = strconv.ParseInt(lit, 10, 64); err != nil {
				return nil, types.ErrorWrongSentence
			}
			lTok = "INT"
			l.Lit = lit
		}
		if !isFloat && !isInt {
			return nil, types.ErrorWrongSentence
		}
	}
	if rTok == "STRING" || rTok == "CHAR" {
		_, err := strconv.ParseInt(r.Lit, 10, 64)
		if err != nil {
			return nil, types.ErrorWrongSentence
		}
		rTok = "INT"
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
	return nil, types.ErrorWrongSentence
}

// findStrInfrontSymbool 查找指定符号之前的数据
func findStrInfrontSymbool(sym string, src []*global.Structure) ([]*global.Structure, int) {
	for k, v := range src {
		if sym == v.Tok {
			return src[:k], k
		}
	}
	return nil, -1
}
