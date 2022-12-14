package lit

import (
	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// FIXME 未针对for语句的合法性做充分检查
// parseIdentedFOR 解析for语句
func (r *expression) parseIdentedFOR(expr []*global.Structure, blocks []*global.Block, innerVar global.InnerVar, i int) ([]*global.Block, int, error) {
	var (
		// rlen 数据长度
		rlen = len(expr)
		// foundCurlyBracket 花括号标记
		foundCurlyBracket bool
		// curlyBracketCount 计算花括号
		curlyBracketCount uint8
		// conditions for 条件
		conditions = make([]*global.Structure, 0, 12)
		// curlyBracketCode 花括号内的代码块
		curlyBracketCode = make([]*global.Structure, 0, 10)
	)

	for j := i; j < rlen; j++ {
		exprJ := expr[j]
		if exprJ.Tok == "{" {
			curlyBracketCount++
			foundCurlyBracket = true
		} else if exprJ.Tok == "}" {
			curlyBracketCount--
		}
		if foundCurlyBracket {
			curlyBracketCode = append(curlyBracketCode, exprJ)
			if curlyBracketCount == 0 {
				i = j
				break
			}
		} else {
			conditions = append(conditions, exprJ)
		}
	}

	blocks = append(blocks, &global.Block{
		Name:   "FOR",
		ForExt: &global.ForExpression{Type: types.TypeForExpressionIteration, Conditions: conditions, Code: curlyBracketCode[1 : len(curlyBracketCode)-1]},
		Type:   types.CodeTypeIdentFOR,
	})
	return blocks, i, nil
}

// execFORType1 解析以下形式的 for 流程控制:
// n = 0; n < y; n ++
func (r *expression) execFORType1(forExpr *global.ForExpression, innerVar global.InnerVar) error {
	var (
		lf         int
		cd1        = make([]*global.Structure, 0, 5)
		cd2        = make([]*global.Structure, 0, 5)
		conditions = forExpr.Conditions
	)

	for i := 0; i < len(conditions); i++ {
		if conditions[i].Tok == ";" {
			if len(cd1) == 0 {
				cd1 = conditions[1 : i+1]
			} else if len(cd2) == 0 {
				cd2 = conditions[len(cd1)+1 : i]
				lf = i
				break
			}
		}
	}

	// n = 0
	if _, err := r.initExpr(cd1, innerVar, nil); err != nil {
		return err
	}

	// n ++
	cd3 := conditions[lf+1:]
	cLen := len(cd3)
	if cLen < 2 {
		return types.ErrorForExpression
	}

	if cLen == 2 {
		if tok := cd3[1].Tok; tok != "++" && tok != "--" {
			return types.ErrorForExpression
		}
	}

	// FIXME 暂未基于此类表达式做通用处理
	if cLen >= 3 {
		if tok := cd3[1].Tok; tok != "+=" && tok != "-=" {
			return types.ErrorForExpression
		}
	}

	cd3 = append(cd3, &global.Structure{Tok: ";", Lit: ";"})
	var result *global.Structure
	for {
		// n < y
		rv, err := r.parse(cd2, innerVar)
		if err != nil {
			return err
		}
		if !global.TransformAllToBool(rv) {
			break
		}

		if result, err = r.initExpr(forExpr.Code, innerVar, &parsing{isInLoop: true}); err != nil {
			return err
		}

		if _, err = r.initExpr(cd3, innerVar, &parsing{isInLoop: true}); err != nil {
			return err
		}

		// 发现 continue 则跳出当前循环
		if result != nil && result.Tok == "continue" {
			continue
		}

		// 发现 break 则跳出当前循环
		if result != nil && result.Tok == "break" {
			break
		}
	}
	return nil
}
