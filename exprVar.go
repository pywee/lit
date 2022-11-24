package lit

import (
	"strconv"

	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

type parseVar struct {
	rlen   int
	tok    string
	r      *expression
	blocks []*global.Block
	expr   []*global.Structure
}

// parseIdentedVAR 解析变量声明
func parseIdentedVAR(arg *parseVar, innerVal global.InnerVar, i int) (int, error) {
	var (
		err     error
		rv      *global.Structure
		tok     = arg.tok
		expr    = arg.expr
		rlen    = arg.rlen
		thisLit = expr[0].Lit
		thisTok = expr[0].Tok
		code    = make([]*global.Structure, 0, 5)
	)

	for j := i; j < rlen; j++ {
		exprJ := expr[j]
		if exprJ.Tok != ";" {
			code = append(code, exprJ)
			continue
		}

		if len(code) < 3 {
			return 0, types.ErrorWrongVarOperation
		}

		code2 := code[2:]
		if tok == "=" {
			if rv, err = arg.r.parse(code2, innerVal); err != nil {
				return 0, err
			}
		} else if tok == "+=" {
			code2 = append(code2, &global.Structure{Tok: ")"})
			nexpr := append([]*global.Structure{{Tok: thisTok, Lit: thisLit}, {Tok: "+"}, {Tok: "("}}, code2...)
			if rv, err = arg.r.parse(nexpr, innerVal); err != nil {
				return 0, err
			}
		} else if tok == "-=" {
			code2 = append(code2, &global.Structure{Tok: ")"})
			nexpr := append([]*global.Structure{{Tok: thisTok, Lit: thisLit}, {Tok: "-"}, {Tok: "("}}, code2...)
			if rv, err = arg.r.parse(nexpr, innerVal); err != nil {
				return 0, err
			}
		} else if tok == "*=" {
			code2 = append(code2, &global.Structure{Tok: ")"})
			nexpr := append([]*global.Structure{{Tok: thisTok, Lit: thisLit}, {Tok: "*"}, {Tok: "("}}, code2...)
			if rv, err = arg.r.parse(nexpr, innerVal); err != nil {
				return 0, err
			}
		} else if tok == "/=" {
			code2 = append(code2, &global.Structure{Tok: ")"})
			nexpr := append([]*global.Structure{{Tok: thisTok, Lit: thisLit}, {Tok: "/"}, {Tok: "("}}, code2...)
			if rv, err = arg.r.parse(nexpr, innerVal); err != nil {
				return 0, err
			}
		} else if tok == "%=" {
			code2 = append(code2, &global.Structure{Tok: ")"})
			nexpr := append([]*global.Structure{{Tok: thisTok, Lit: thisLit}, {Tok: "%"}, {Tok: "("}}, code2...)
			if rv, err = arg.r.parse(nexpr, innerVal); err != nil {
				return 0, err
			}
		} else if tok == "&=" {
			code2 = append(code2, &global.Structure{Tok: ")"})
			nexpr := append([]*global.Structure{{Tok: thisTok, Lit: thisLit}, {Tok: "&"}, {Tok: "("}}, code2...)
			if rv, err = arg.r.parse(nexpr, innerVal); err != nil {
				return 0, err
			}
		} else if tok == "|=" {
			code2 = append(code2, &global.Structure{Tok: ")"})
			nexpr := append([]*global.Structure{{Tok: thisTok, Lit: thisLit}, {Tok: "|"}, {Tok: "("}}, code2...)
			if rv, err = arg.r.parse(nexpr, innerVal); err != nil {
				return 0, err
			}
		} else if tok == "^=" {
			code2 = append(code2, &global.Structure{Tok: ")"})
			nexpr := append([]*global.Structure{{Tok: thisTok, Lit: thisLit}, {Tok: "^"}, {Tok: "("}}, code2...)
			if rv, err = arg.r.parse(nexpr, innerVal); err != nil {
				return 0, err
			}
		}
		i = j
		innerVal[thisLit] = rv
		break
	}

	return i, nil
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
