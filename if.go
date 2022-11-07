package lit

// 2022.11.02 备注
// 针对 if 句子的解析, 例如
// if a... else if b... else
// 其实最终的逻辑可转换为 if a || b || c 的逻辑
// if 语句通常在作用域内 所以需要维护 innerVal 上下文

import (
	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

type parsedIf struct {
	i      int
	blocks []*global.Block
}

// parseIdentedIF 解析if
func parseIdentedIF(blocks []*global.Block, expr []*global.Structure, i int, rlen int) (*parsedIf, error) {
	var (
		count      int16
		parsed     = new(parsedIf)
		body       = make([]*global.Structure, 0, 10)
		conditions = make([]*global.Structure, 0, 10)
	)
	block := &global.Block{
		Type:  types.CodeTypeIdentIF,
		Code:  make([]*global.Structure, 0, 5),
		IfExt: make([]*global.ExIf, 0, 5),
	}

	for j := i; i < rlen; j++ {
		exprJ := expr[j]
		block.Code = append(block.Code, exprJ)
		if exprJ.Tok == "{" {
			count++
			// global.Output(conditions)
		} else if exprJ.Tok == "}" {
			count--
			if count < 0 {
				return nil, types.ErrorIfExpression
			}
			if count == 0 {
				body = append(body, exprJ)
				block.IfExt = append(block.IfExt, &global.ExIf{
					Tok:          "if",
					Body:         body,
					BodyLen:      len(body),
					ConditionLen: len(conditions),
					Condition:    conditions,
				})
				parsed.i = j
				parsed.blocks = append(blocks, block)
				break
			}
		}
		if count > 0 {
			body = append(body, exprJ)
		}
		if j > i && count == 0 {
			conditions = append(conditions, exprJ)
		}
	}

	// TODO 错误处理
	return parsed, nil
}

type pasedElse struct {
	i         int
	foundElse bool
	blocks    []*global.Block
}

// parseIdentELSE 解析if句子else和elseif部分
func parseIdentELSE(blocks []*global.Block, expr []*global.Structure, i int, rlen int) (*pasedElse, error) {
	var (
		count      int16
		foundElse  bool
		parsed     = new(pasedElse)
		elseIF     = "else"
		code       = make([]*global.Structure, 0, 5)
		body       = make([]*global.Structure, 0, 10)
		conditions = make([]*global.Structure, 0, 10)
	)
	if i+1 < rlen && expr[i+1].Tok == "if" {
		elseIF = "elseif"
	}

	foundElse = true
	for j := i; j < rlen; j++ {
		exprJ := expr[j]
		code = append(code, exprJ)
		if exprJ.Tok == "{" {
			count++
		} else if exprJ.Tok == "}" {
			count--
			if count < 0 {
				return nil, types.ErrorIfExpression
			}
			if count == 0 {
				body = append(body, exprJ)
				if blen := len(blocks); blen > 0 {
					if elseIF == "elseif" {
						if len(conditions) <= 2 {
							return nil, types.ErrorIfExpression
						}
						conditions = conditions[2:]
					} else {
						conditions = conditions[1:]
					}
					blocksEnd := blocks[blen-1]
					bsCode := blocksEnd.Code
					blocks[blen-1].IfExt = append(blocksEnd.IfExt, &global.ExIf{
						Tok:          elseIF,
						Body:         body,
						Condition:    conditions,
						BodyLen:      len(body),
						ConditionLen: len(conditions),
					})
					blocks[blen-1].Code = append(bsCode, code...)
					if j+1 < rlen && expr[j+1].Tok == ";" && expr[j+1].Lit == "\n" {
						foundElse = false
					}
					parsed.i = j
					parsed.blocks = blocks
					parsed.foundElse = foundElse
					break
				}
				return nil, types.ErrorIfExpression
			}
		}
		if count == 0 {
			conditions = append(conditions, exprJ)
		} else if count > 0 {
			body = append(body, exprJ)
		}
	}
	return parsed, nil
}
