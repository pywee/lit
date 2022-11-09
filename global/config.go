package global

type Block struct {
	// Type [0-未知, 非法定义; 1.变量赋值; 2.函数声明; 3.函数调用; 4.if语句; 5.for; 6.return]
	Type int8
	// 关键词 如果类型是函数 则显示函数名
	Name string
	// Code 代码块
	Code []*Structure
	// if 语句
	IfExt []*ExIf
}

type ExIf struct {
	// Tok 标识
	Tok string
	// Condition if条件
	Condition []*Structure
	// ConditionLen 条件句子长度
	ConditionLen int
	// Body if句子内数据
	// 此处仍会出现if 需要通过递归层层解析
	Body []*Structure
	// BodyLen 数据体长度
	BodyLen int
}

// Structure 基础数据解析 每一个符号代码
type Structure struct {
	Position string
	Tok      string
	Lit      string
}

func (s *Structure) IsBoolTrue() bool {
	return s.Lit == "true"
}

func (s *Structure) IsBoolFalse() bool {
	return s.Lit == "false"
}

func (s *Structure) Val() string {
	return s.Lit
}

func (s *Structure) Type() string {
	return s.Tok
}
