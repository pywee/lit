package lit

import (
	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

type forExpression struct {
	i      int
	blocks []*global.Block
}

type forArgs struct {
	i        int
	rlen     int
	blocks   []*global.Block
	expr     []*global.Structure
	innerVar map[string]*global.Structure
}

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
		ForExt: &global.ForExpression{Type: 1, Conditions: conditions, Code: curlyBracketCode[1 : len(curlyBracketCode)-1]},
		Type:   types.CodeTypeIdentFOR,
	})
	return blocks, i, nil
}

// execFORType1 解析以下形式的 for 流程控制:
// n = 0; n < y; n ++
func (r *expression) execFORType1(forExpr *global.ForExpression, innerVar global.InnerVar) (*global.Structure, error) {
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
	if _, err := r.initExpr(cd1, innerVar); err != nil {
		return nil, err
	}

	// n ++
	cd3 := conditions[lf+1:]
	if len(cd3) < 2 {
		return nil, types.ErrorForExpression
	}
	if tok := cd3[1].Tok; tok != "++" && tok != "--" && tok != "+=" && tok != "-=" {
		return nil, types.ErrorForExpression
	}
	cd3 = append(cd3, &global.Structure{Tok: ";", Lit: ";"})

	for {
		// n < y
		rv, err := r.parse(cd2, innerVar)
		if err != nil {
			return nil, err
		}
		if !global.ChangeToBool(rv) {
			break
		}

		if _, err = r.initExpr(forExpr.Code, innerVar); err != nil {
			return nil, err
		}

		if _, err = r.initExpr(cd3, innerVar); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
