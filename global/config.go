package global

type Structure struct {
	Position string
	Tok      string
	Lit      string
}

func (s *Structure) Val() string {
	return s.Lit
}

func (s *Structure) Type() string {
	return s.Tok
}
