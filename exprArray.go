package lit

import "github.com/pywee/lit/global"

// parseIdentARRAY 解析数组定义
// 递归操作 得出树状数组
// 此处仅根据字义的结构保存好树状的数据 在被调用时才会对树状数组里的各个元素进行解析
func parseIdentARRAY(expr []*global.Structure) *global.Array {
	var (
		key       int
		rlen      = len(expr)
		arr       = new(global.Array)
		arrays    = make([]*global.ArrayIdent, 0, 3)
		innerExpr = make([]*global.Structure, 0, 3)
	)

	for i := 1; i < rlen; i++ {
		if expr[i].Tok == "," {
			if ilen := len(innerExpr); ilen > 2 && innerExpr[0].Tok == "[" && innerExpr[ilen-1].Tok == "]" {
				if sonArr := parseIdentARRAY(innerExpr); sonArr != nil {
					arrays = append(arrays, &global.ArrayIdent{Child: sonArr, InnerKey: key})
				}
			} else {
				arrays = append(arrays, &global.ArrayIdent{Values: innerExpr, InnerKey: key})
			}
			key++
			innerExpr = nil
			continue
		}
		if i+1 < rlen {
			innerExpr = append(innerExpr, expr[i])
		}
	}

	if ilen := len(innerExpr); ilen > 2 && innerExpr[0].Tok == "[" && innerExpr[ilen-1].Tok == "]" {
		if sonArr := parseIdentARRAY(innerExpr); sonArr != nil {
			arrays = append(arrays, &global.ArrayIdent{Child: sonArr, InnerKey: key})
		}
	} else {
		arrays = append(arrays, &global.ArrayIdent{Values: innerExpr, InnerKey: key})
	}
	arr.List = arrays
	return arr
}
