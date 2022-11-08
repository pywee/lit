package lit

import (
	"github.com/pywee/lit/global"
)

func (r *Expression) parseAnd(expr []*global.Structure, k int, pos string, innerVariable map[string]*global.Structure) (*global.Structure, error) {
	var (
		err error
		rv  *global.Structure
	)
	if rv, err = r.parse(expr[:k], pos, innerVariable); err != nil {
		return nil, err
	}
	// if rv.Tok != "IDENT" {
	// }
	if !global.ChangeToBool(rv) {
		return rv, nil
	}
	if rv, err = r.parse(expr[k+1:], pos, innerVariable); err != nil {
		return nil, err
	}
	global.ChangeToBool(rv)
	return rv, nil
}

func (r *Expression) parseOr(expr []*global.Structure, k int, pos string, innerVariable map[string]*global.Structure) (*global.Structure, error) {
	var (
		err error
		rv  *global.Structure
	)
	if rv, err = r.parse(expr[:k], pos, innerVariable); err != nil {
		return nil, err
	}
	if global.ChangeToBool(rv) {
		return rv, nil
	}
	if rv, err = r.parse(expr[k+1:], pos, innerVariable); err != nil {
		return nil, err
	}

	global.ChangeToBool(rv)
	return rv, nil
}
