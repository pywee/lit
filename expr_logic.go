package goExpr

import (
	fn "github.com/pywee/goExpr/function"
	"github.com/pywee/goExpr/global"
)

func (r *Expression) parseAnd(expr []*global.Structure, k int, pos string) (*global.Structure, error) {
	var (
		err error
		rv  *global.Structure
	)
	if rv, err = r.parse(expr[:k], pos); err != nil {
		return nil, err
	}
	// if rv.Tok != "IDENT" {
	// }
	if _, ok := fn.ChangeBool(rv); !ok {
		return rv, nil
	}
	if rv, err = r.parse(expr[k+1:], pos); err != nil {
		return nil, err
	}
	fn.ChangeBool(rv)
	return rv, nil
}

func (r *Expression) parseOr(expr []*global.Structure, k int, pos string) (*global.Structure, error) {
	var (
		err error
		rv  *global.Structure
	)
	if rv, err = r.parse(expr[:k], pos); err != nil {
		return nil, err
	}
	if _, ok := fn.ChangeBool(rv); ok {
		return rv, nil
	}
	if rv, err = r.parse(expr[k+1:], pos); err != nil {
		return nil, err
	}
	fn.ChangeBool(rv)
	return rv, nil
}

// func logic() (*global.Structure, error) {
// 	var (
// 		err error
// 		rv  *global.Structure
// 	)

// 	if v.Tok == "&&" && firstKey == -1 {
// 		if rv, err = r.parse(expr[:k], pos); err != nil {
// 			return nil, err
// 		}
// 		// if rv.Tok != "IDENT" {
// 		// }
// 		if fn.ChangeBool(rv).IsBoolFalse() {
// 			return rv, nil
// 		}
// 		if rv, err = r.parse(expr[k+1:], pos); err != nil {
// 			return nil, err
// 		}
// 		return fn.ChangeBool(rv), nil
// 	} else if v.Tok == "||" && firstKey == -1 {
// 		if rv, err = r.parse(expr[:k], pos); err != nil {
// 			return nil, err
// 		}
// 		if fn.ChangeBool(rv).IsBoolTrue() {
// 			return rv, nil
// 		}
// 		if rv, err = r.parse(expr[k+1:], pos); err != nil {
// 			return nil, err
// 		}
// 		return fn.ChangeBool(rv), nil
// 	}
// }
