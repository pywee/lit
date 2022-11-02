package lit

import (
	fn "github.com/pywee/lit/function"
	"github.com/pywee/lit/global"
)

func (r *Expression) parseAnd(expr []*global.Structure, k int, pos string) (*global.Structure, error) {
	var (
		err error
		rv  *global.Structure
	)
	if rv, err = r.parse(expr[:k], pos, nil); err != nil {
		return nil, err
	}
	// if rv.Tok != "IDENT" {
	// }
	if _, ok := fn.ChangeToBool(rv); !ok {
		return rv, nil
	}
	if rv, err = r.parse(expr[k+1:], pos, nil); err != nil {
		return nil, err
	}
	fn.ChangeToBool(rv)
	return rv, nil
}

func (r *Expression) parseOr(expr []*global.Structure, k int, pos string) (*global.Structure, error) {
	var (
		err error
		rv  *global.Structure
	)
	if rv, err = r.parse(expr[:k], pos, nil); err != nil {
		return nil, err
	}
	if _, ok := fn.ChangeToBool(rv); ok {
		return rv, nil
	}
	if rv, err = r.parse(expr[k+1:], pos, nil); err != nil {
		return nil, err
	}

	fn.ChangeToBool(rv)
	return rv, nil
}
