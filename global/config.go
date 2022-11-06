package global

type Logic struct {
	// Type [0-未知, 非法定义; 1.变量赋值; 2.函数声明; 3.函数调用; 4.if语句; 5.for; 6.return]
	Type int8
	// Code 代码块
	Code []*Structure
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
