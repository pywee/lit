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
func (r *Expression) parseIdentedFOR(arg *forArgs) ([]*global.Block, int, error) {
	var (
		// i 解析代码块时的游标
		i = arg.i
		// rlen 代码块长度
		rlen = arg.rlen
		// expr 原始代码
		expr   = arg.expr
		blocks = arg.blocks
		// semicolonCount 计算分号
		semicolonCount uint8
		// conditions
		conditions = make([][]*global.Structure, 0, 3)
		// forExpr 通过;分割后的表达式
		forExpr = make([]*global.Structure, 0, 5)
		// foundCurlyBracket 花括号标记
		foundCurlyBracket bool
		// curlyBracketCount 计算花括号
		curlyBracketCount uint8
		// curlyBracketCode 花括号内的代码块
		curlyBracketCode = make([]*global.Structure, 0, 10)
	)

	// global.Output(expr)
	// block := new(global.Block)
	for j := i + 1; j < rlen; j++ {
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
			continue
		}

		forExpr = append(forExpr, exprJ)
		if exprJ.Tok == ";" {
			semicolonCount++
			if semicolonCount == 1 {
				if len(forExpr) < 4 || forExpr[1].Tok != "=" {
					return nil, 0, types.ErrorForExpression
				}
				i = j
				conditions = append(conditions, forExpr)
				forExpr = nil
			} else if semicolonCount == 2 {
				fLen := len(forExpr)
				if fLen == 0 {
					return nil, 0, types.ErrorForExpression
				}
				i = j
				conditions = append(conditions, forExpr)
				forExpr = nil
			}
		}
	}

	fLen := len(forExpr)
	if fLen < 2 || len(conditions) < 2 {
		return nil, 0, types.ErrorForExpression
	}
	forExpr = append(forExpr, &global.Structure{Tok: ";", Lit: ";"})
	conditions = append(conditions, forExpr)
	blocks = append(blocks, &global.Block{
		// Name:   "FOR",
		ForExt: &global.ForExpression{Type: 1, Conditions: conditions, Code: curlyBracketCode},
		Type:   types.CodeTypeIdentFOR,
	})
	return blocks, i, nil
}
