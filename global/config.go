package global

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
