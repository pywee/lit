package lit

import "github.com/pywee/lit/global"

// parseAnd 解析逻辑运算符 "&&"
func (r *expression) parseAnd(expr, nExpr []*global.Structure, innerVar global.InnerVar, i int) (*global.Structure, error) {
	var (
		err   error
		left  *global.Structure
		right *global.Structure
	)
	if left, err = r.parse(nExpr, innerVar); err != nil {
		return nil, err
	}
	if !global.ChangeToBool(left) {
		return left, nil
	}

	if right, err = r.parse(expr[i+1:], innerVar); err != nil {
		return nil, err
	}
	global.ChangeToBool(right)

	return right, nil
}

// parseOr 解析逻辑运算符 "||"
func (r *expression) parseOr(expr, nExpr []*global.Structure, innerVar global.InnerVar, i int) (*global.Structure, error) {
	var (
		err   error
		left  *global.Structure
		right *global.Structure
	)

	if left, err = r.parse(nExpr, innerVar); err != nil {
		return nil, err
	}
	if global.ChangeToBool(left) {
		return left, nil
	}

	if right, err = r.parse(expr[i+1:], innerVar); err != nil {
		return nil, err
	}
	global.ChangeToBool(right)
	return right, nil
}
