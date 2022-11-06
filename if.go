package lit

import (
	fn "github.com/pywee/lit/function"
	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
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

type ifCondition struct {
	// countIF 计算if出现次数
	countIF int8
	// countELSE 计算else出现次数
	countELSE int8
	// countELSEIF 计算elseif出现次数
	countELSEIF int8
}

type ExIf struct {
	// tok 标识
	tok string
	// condition if条件
	condition []*global.Structure
	// conditionLen 条件句子长度
	conditionLen int
	// body if句子内数据
	// 此处仍会出现if 需要通过递归层层解析
	body []*global.Structure
	// bodyLen 数据体长度
	bodyLen int
}

// parseIf 解析 if 语句
func (r *Expression) parseIf(expr []*global.Structure, pos string, innerVariable map[string]*global.Structure) error {
	var (
		// err 错误处理
		err error
		// foundIF 发现if语句标记
		foundIF bool
		// foundElse
		// foundElse bool
		// curlyBracket 标记大括号
		curlyBracket int
		// thisIfTok 标识名称 if/ else if/ else
		thisIfTok string
		// thisIfBody if 数据体
		thisIfBody = make([]*global.Structure, 0, 10)
		// thisIfConditions 单个 if 句子内的条件
		thisIfConditions = make([]*global.Structure, 0, 10)
		// expressionIF if 表达式列表
		// 此处为平级标记 仅标记一次最外部发现的if 然后对内部嵌套的 if 进行递归处理
		// if.. else if ... else
		expressionIF = make([]*ExIf, 0, 5)
		// conditionChecking 检查if语句合法性
		conditionChecking = new(ifCondition)
	)

	if innerVariable == nil {
		innerVariable = make(map[string]*global.Structure, 5)
	}

	for k, v := range expr {
		// 在此处解析 if 语句
		if !foundIF && curlyBracket == 0 && v.Tok == "if" {
			foundIF = true
			thisIfTok = "if"
			continue
		}
		if curlyBracket == 0 && v.Tok == "else" && !foundIF {
			foundIF = true
			thisIfTok = "else"
			continue
		}
		if k > 0 && v.Tok == "if" && expr[k-1].Tok == "else" {
			foundIF = true
			thisIfTok = "elseif"
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
				if err = checkIfConditions(thisIfTok, thisIfConditions, conditionChecking); err != nil {
					return err
				}
				thisIfBody = append(thisIfBody, v)
				expressionIF = append(expressionIF, &ExIf{
					tok:          thisIfTok,
					condition:    thisIfConditions,
					conditionLen: len(thisIfConditions),
					body:         thisIfBody,
					bodyLen:      len(thisIfBody),
				})
				thisIfBody = nil
				thisIfConditions = nil
			}

			// FIXME 此处有问题
			// continue
		}

		if curlyBracket > 0 {
			thisIfBody = append(thisIfBody, v)
			continue
		}
		if !foundIF && curlyBracket == 0 {
			if len(thisIfBody) > 0 {
				if err = checkIfConditions(thisIfTok, thisIfConditions, conditionChecking); err != nil {
					return err
				}
				expressionIF = append(expressionIF, &ExIf{
					tok:          thisIfTok,
					condition:    thisIfConditions,
					conditionLen: len(thisIfConditions),
					body:         thisIfBody,
					bodyLen:      len(thisIfBody),
				})
			}
			thisIfBody = nil
			thisIfConditions = nil
		}
		if foundIF {
			thisIfConditions = append(thisIfConditions, v)
			continue
		}

	}

	// 解析if代码块
	if len(expressionIF) > 0 {
		if conditionChecking.countELSE > 1 {
			return types.ErrorIlligleIfExpressionOfElse
		}
		if conditionChecking.countIF > 1 || conditionChecking.countIF <= 0 {
			return types.ErrorIlligleIfExpressionOfIf
		}
		if err := r.parseIfConditionList(expressionIF, pos, innerVariable); err != nil {
			return err
		}
	}
	return nil
}

// checkIfConditions 检查if语句定义
func checkIfConditions(tok string, conditions []*global.Structure, i *ifCondition) error {
	if tok == "if" {
		i.countIF++
		if len(conditions) == 0 {
			return types.ErrorIfExpressionWithoutConditions
		}
	} else if tok == "elseif" {
		i.countELSEIF++
		if len(conditions) == 0 {
			return types.ErrorIfExpressionWithoutConditions
		}
	} else if tok == "else" {
		i.countELSE++
	}
	return nil
}

func (r *Expression) parseIfConditionList(expressionIF []*ExIf, pos string, innerVariable map[string]*global.Structure) error {
	var (
		err  error
		ELSE bool
		ret  *global.Structure
	)

	for _, vv := range expressionIF {
		// 处理 else
		// 此时是没有条件的
		if vv.tok == "else" && len(vv.condition) == 0 {
			if ELSE {
				return types.ErrorIfExpression
			}
			ELSE = true
			if vv.bodyLen > 2 && vv.body[0].Tok == "{" && vv.body[vv.bodyLen-1].Tok == "}" {
				if vv.tok == "elseif" && vv.conditionLen <= 0 {
					return types.ErrorIfExpression
				}
				if _, err = r.parseIfBody(vv.body[1:vv.bodyLen-1], pos, innerVariable); err != nil {
					return err
				}
			}
			break
		}

		// 处理 if
		if ret, err = r.parse(vv.condition, pos, innerVariable); err != nil {
			return err
		}
		if ret == nil {
			return types.ErrorIfExpression
		}

		// if ret.Lit == "true" {
		if _, ok := fn.ChangeToBool(ret); ok {
			// 截取if内部数据 向内部截取一层
			if vv.bodyLen > 2 && vv.body[0].Tok == "{" && vv.body[vv.bodyLen-1].Tok == "}" {
				vv.body = vv.body[1 : vv.bodyLen-1]
				if _, err = r.parseIfBody(vv.body, pos, innerVariable); err != nil {
					return err
				}
			}
			// FIXME
			// 当前句子为真就直接完成这个函数了 需进一步测试
			break
		}
	}
	return nil
}

// parseIfBody 解析if句子内的代码块
func (r *Expression) parseIfBody(expr []*global.Structure, pos string, innerVariable map[string]*global.Structure) (*global.Structure, error) {
	var (
		list    = make([]*global.Structure, 0, 5)
		foundIF bool
	)

	for _, v := range expr {
		if v.Tok == "if" && !foundIF {
			foundIF = true
		}
		if !foundIF && v.Tok == ";" {
			// 此处需要再次分割 = 符号为变量赋值
			var (
				err   error
				vName string
				rv    *global.Structure
			)
			if vleft, vLeftListEndIdx := findStrInfrontSymbool("=", list); vLeftListEndIdx != -1 {
				if vLeftListEndIdx == 1 {
					vName = vleft[0].Lit
					list = list[vLeftListEndIdx+1:]
				}
			}

			if rv, err = r.parse(list, pos, innerVariable); err != nil {
				return nil, err
			}
			if len(vName) > 0 && rv != nil {
				innerVariable[vName] = rv
			}
			list = nil
			continue
		}
		list = append(list, v)
	}

	// 花括号闭合处
	if rlen := len(list); rlen > 0 {
		if list[rlen-1].Tok == ";" && list[rlen-1].Lit == "\n" {
			list = list[:rlen-1]
		}
		if _, err := r.parse(list, pos, innerVariable); err != nil {
			return nil, err
		}
	}

	// global.Output(list)
	return nil, nil
}
