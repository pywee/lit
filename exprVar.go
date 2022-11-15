package lit

import (
	"strconv"

	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// parseIdentedVAR 解析变量声明
func parseIdentedVAR(blocks []*global.Block, expr []*global.Structure, i int, rlen int) ([]*global.Block, int) {
	code := make([]*global.Structure, 0, 5)
	for j := i; j < rlen; j++ {
		exprJ := expr[j]
		if exprJ.Tok == ";" {
			blocks = append(blocks, &global.Block{Type: types.CodeTypeIdentVAR, Code: code})
			i = j
			break
		}
		code = append(code, exprJ)
	}
	return blocks, i
}

// parseIdentedVarPLUS 解析变量自增
func parseIdentedVarPLUS(blocks []*global.Block, expr []*global.Structure, i int, rlen int) ([]*global.Block, int) {
	vplus := make([]*global.Structure, 0, 3)
	for j := i; j < rlen; j++ {
		exprJ := expr[j]
		if exprJ.Tok == ";" {
			blocks = append(blocks, &global.Block{
				Name: expr[i].Lit,
				Type: types.CodeTypeVariablePlus,
				Code: vplus,
			})
			i = j
			break
		}
		vplus = append(vplus, exprJ)
	}
	return blocks, i
}

// parseIdentedVarREDUCE 解析变量自减
func parseIdentedVarREDUCE(blocks []*global.Block, expr []*global.Structure, i int, rlen int) ([]*global.Block, int) {
	vreduce := make([]*global.Structure, 0, 3)
	for j := i; j < rlen; j++ {
		exprJ := expr[j]
		if exprJ.Tok == ";" {
			blocks = append(blocks, &global.Block{
				Name: expr[i].Lit,
				Type: types.CodeTypeVariableReduce,
				Code: vreduce,
			})
			i = j
			break
		}
		vreduce = append(vreduce, exprJ)
	}
	return blocks, i
}

// execVarPlusReduce 执行变量递增、递减操作
func execVarPlusReduce(block *global.Block, innerVar map[string]*global.Structure, plus bool) (*global.Structure, error) {
	blockName := block.Name
	v1, ok := innerVar[blockName]
	if !ok {
		return nil, types.ErrorNotFoundVariable
	}

	vTok := v1.Tok
	if vTok != "INT" && vTok != "FLOAT" && vTok != "STRING" {
		return nil, types.ErrorHandleUnsupported
	}

	var rv *global.Structure
	if vTok == "INT" {
		retParsed, err := strconv.ParseInt(v1.Lit, 10, 64)
		if err != nil {
			return nil, err
		}
		if plus {
			retParsed++
		} else {
			retParsed--
		}
		rv = &global.Structure{Tok: "INT", Lit: strconv.FormatInt(retParsed, 10)}
	} else if vTok == "FLOAT" {
		retParsed, err := strconv.ParseFloat(v1.Lit, 64)
		if err != nil {
			return nil, err
		}
		if plus {
			retParsed++
		} else {
			retParsed--
		}
		// rv = &global.Structure{Tok: "FLOAT", Lit: strconv.FormatFloat(retParsed, 'e', -1, 64)}
		rv = &global.Structure{Tok: "FLOAT", Lit: strconv.FormatFloat(retParsed, 'f', -1, 64)}
	} else if vTok == "STRING" {
		if ok, _ := global.IsInt(v1.Lit); ok {
			retParsed, err := strconv.ParseInt(v1.Lit, 10, 64)
			if err != nil {
				return nil, err
			}
			if plus {
				retParsed++
			} else {
				retParsed--
			}
			rv = &global.Structure{Tok: "INT", Lit: strconv.FormatInt(retParsed, 10)}
		} else if ok, _ := global.IsFloat(v1.Lit); ok {
			retParsed, err := strconv.ParseFloat(v1.Lit, 64)
			if err != nil {
				return nil, err
			}
			if plus {
				retParsed++
			} else {
				retParsed--
			}
			rv = &global.Structure{Tok: "FLOAT", Lit: strconv.FormatFloat(retParsed, 'f', -1, 64)}
		}
	}
	if rv != nil {
		innerVar[blockName] = rv
	}
	return rv, nil
}
