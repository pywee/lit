package lit

import (
	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// parseIdentedVAR 解析变量声明
func parseIdentedVAR(blocks []*global.Block, expr []*global.Structure, i int, rlen int) ([]*global.Block, int) {
	block := &global.Block{Type: types.CodeTypeIdentVAR, Code: make([]*global.Structure, 0, 5)}
	expr[i].Tok = "VAR"
	for j := i; j < rlen; j++ {
		exprJ := expr[j]
		if exprJ.Tok == ";" {
			blocks = append(blocks, block)
			block = nil
			i = j
			break
		}
		block.Code = append(block.Code, exprJ)
	}
	return blocks, i
}
