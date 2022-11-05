package lit

import (
	"github.com/pywee/lit/global"
)

// 2022.11.02 备注
// 针对 if 句子的解析, 例如
// if a... else if b... else
// 其实最终的逻辑可转换为 if a || b || c 的逻辑
// if 语句通常在作用域内 所以需要维护 innerVal 上下文

// ExpressionIfStructure 保存 if 句子信息
type ExpressionIfStructure struct {
	List []*ExIf
}

type ExIf struct {
	// condition if条件
	condition []*global.Structure
	// body if句子内数据
	// 此处仍会出现if 需要通过递归层层解析
	body []*global.Structure
	// bodyLen 数据体长度
	bodyLen int
}

// parseIf 解析 if 语句
func (r *Expression) parseIf(expr []*global.Structure, pos string, innerVariable map[string]*global.Structure) (*global.Structure, error) {
	var (
		// foundIF 发现if语句标记
		foundIF bool
		// foundElse
		// foundElse bool
		// curlyBracket 标记大括号
		curlyBracket int
		// thisIfBody if 数据体
		thisIfBody = make([]*global.Structure, 0, 10)
		// thisIfConditions 单个 if 句子内的条件
		thisIfConditions = make([]*global.Structure, 0, 10)
		// expressionIF if 表达式列表
		// 此处为平级标记 仅标记一次最外部发现的if 然后对内部嵌套的 if 进行递归处理
		// if.. else if ... else
		expressionIF = make([]*ExIf, 0, 5)
	)

	for _, v := range expr {
		// 在此处解析 if 语句
		if !foundIF && curlyBracket == 0 && v.Tok == "if" {
			foundIF = true
			continue
		}
		if curlyBracket == 0 && v.Tok == "else" && !foundIF {
			foundIF = true
			continue
		}
		if v.Tok == "{" {
			if curlyBracket == 0 {
				foundIF = false
			}
			curlyBracket++
		}
		if v.Tok == "}" {
			curlyBracket--
			if curlyBracket == 0 {
				thisIfBody = append(thisIfBody, v)
			}
		}
		if curlyBracket > 0 {
			thisIfBody = append(thisIfBody, v)
			continue
		}
		if !foundIF && curlyBracket == 0 {
			expressionIF = append(expressionIF, &ExIf{condition: thisIfConditions, body: thisIfBody, bodyLen: len(thisIfBody)})
			thisIfBody = nil
			thisIfConditions = nil
		}
		if foundIF {
			thisIfConditions = append(thisIfConditions, v)
			continue
		}

		if len(expressionIF) > 0 {
			var (
				err error
				ret *global.Structure
			)

			for _, vv := range expressionIF {
				if vv.bodyLen <= 0 {
					continue
				}

				// 处理 else
				// 此时是没有条件的
				if len(vv.condition) == 0 {
					if _, err = r.parse(vv.body[:vv.bodyLen-1], pos, innerVariable); err != nil {
						return nil, err
					}
					break
				}

				// 处理 if
				if ret, err = r.parse(vv.condition, pos, innerVariable); err != nil {
					return nil, err
				}
				if ret.Lit == "true" {
					// 截取if内部数据
					if vv.bodyLen > 2 && vv.body[0].Tok == "{" && vv.body[vv.bodyLen-1].Tok == "}" {
						vv.body = vv.body[1 : vv.bodyLen-1]
						r.parseIfBody(vv.body, pos, innerVariable)
					}
					// if _, err = r.parse(vv.body, pos, innerVariable); err != nil {
					// 	return nil, err
					// }
					break
				}
			}
			expressionIF = nil
		}
	}
	return nil, nil
}

// parseIfBody 解析if句子内的代码块
func (r *Expression) parseIfBody(expr []*global.Structure, pos string, innerVariable map[string]*global.Structure) (*global.Structure, error) {
	var (
		list    = make([]*global.Structure, 0, 5)
		foundIF bool
	)

	// global.Output(expr)
	for _, v := range expr {
		if v.Tok == "if" && !foundIF {
			foundIF = true
		}
		if !foundIF && v.Tok == ";" {
			if _, err := r.parse(list, pos, innerVariable); err != nil {
				return nil, err
			}
			list = nil
		}
		list = append(list, v)
	}

	if rlen := len(list); rlen > 0 {
		if list[rlen-1].Tok == ";" && list[rlen-1].Lit == "\n" {
			list = list[:rlen-1]
		}
		if _, err := r.parse(list, pos, innerVariable); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
