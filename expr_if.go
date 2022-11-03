package lit

import "github.com/pywee/lit/global"

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
}
