package lit

import (
	"errors"

	fn "github.com/pywee/lit/function"
	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// parseFuncExec 解析自定义函数声明
func parseFuncExec(blocks []*global.Block, expr []*global.Structure, i int, rlen int) ([]*global.Block, int) {
	block := &global.Block{Type: types.CodeTypeFunctionExec, Code: make([]*global.Structure, 0, 10)}
	for j := i; j < rlen; j++ {
		if expr[j].Tok == ";" {
			blocks = append(blocks, block)
			i = j
			break
		}
		block.Code = append(block.Code, expr[j])
	}
	return blocks, i
}

// parseIdentFUNC 解析函数定义
func parseIdentFUNC(funcBlocks []*fn.FunctionInfo, expr []*global.Structure, i int, rlen int) ([]*fn.FunctionInfo, int, error) {
	var (
		bracket uint16
		block   = &global.Block{Type: types.CodeTypeIdentFN, Code: make([]*global.Structure, 0, 30)}
	)
	for j := i; j < rlen; j++ {
		block.Code = append(block.Code, expr[j])
		if expr[j].Tok == "{" {
			bracket++
		} else if expr[j].Tok == "}" {
			bracket--
			if bracket == 0 {
				if len(block.Code) < 7 {
					return nil, 0, errors.New(expr[i].Position + types.ErrorFunctionIlligle.Error())
				}
				funcsParsed, err := cfn.ParseCutFunc(block.Code, expr[i].Position)
				if err != nil {
					return nil, 0, errors.New(expr[i].Position + err.Error())
				}
				funcBlocks = append(funcBlocks, funcsParsed)
				i = j
				break
			}
		}
	}
	return funcBlocks, i, nil
}
