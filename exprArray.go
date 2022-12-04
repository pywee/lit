package lit

import (
	"github.com/pywee/lit/global"
)

// parseIdentARRAY 解析数组定义
// 递归操作 得出树状数组
// 此处仅根据字义的结构保存好树状的数据 在被调用时才会对树状数组里的各个元素进行解析

/* 请注意
 * 出于对执行效率和规范性考虑 Lit 的数组的下标必须是顺序的
 * 暂不支持非整型下标的访问和声明，但可以支持整型数字的字符串访问
 * 请看如下列子：
 * 声明数组 a = ['x', 'y', 'z'] 此时数组的下标分别为 0 1 2
 * 合法访问 a[0], a[1], a[2]
 * 合法访问 a["0"], a["1"], a["2"]
 * 非法访问 a["0.1"], a["你好"], a["hello"]

 * 数组不支持新增不存在的下标，这与 php 语言在很大程度上"完全不同"
 * 请看如下列子
 * 这样的声明是不合法的
 * a = ['1' => 'x', 'hello' => 'y', 'world' => 'z']
 * a['xxx'] = 1

 * 要叠加数组 必须通过 append 函数
 * a = []
 * b = []
 * a = append(a, 1, 2, b, '你好')
 * 循环当前数组时 输出: 1 2 Array 你好
 */
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
