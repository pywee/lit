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
		// conditions
		conditions = make([]*global.Structure, 0, 10)
		// foundCurlyBracket 花括号标记
		foundCurlyBracket bool
		// curlyBracketCount 计算花括号
		curlyBracketCount uint8
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

	global.Output(curlyBracketCode)

	blocks = append(blocks, &global.Block{
		Name:   "FOR",
		ForExt: &global.ForExpression{Type: 1, Conditions: conditions, Code: curlyBracketCode},
		Type:   types.CodeTypeIdentFOR,
	})
	return blocks, i, nil
}

// execFORType1 解析以下形式的 for 流程控制:
// n = 0; n < y; n ++
func (r *expression) execFORType1(forExpr *global.ForExpression, innerVar global.InnerVar) (*global.Structure, error) {
	return nil, nil
}
