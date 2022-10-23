package goExpr

import (
	"strconv"

	fn "github.com/pywee/goExpr/function"
	"github.com/pywee/goExpr/global"
)

// parseCompare 比较运算
// == >= <= != > <
func (r *Expression) parseCompare(expr []*global.Structure, pos string) (*global.Structure, error) {
	var (
		err     error
		rvLeft  *global.Structure
		rvRight *global.Structure
	)
	for k, v := range expr {
		if v.Tok == "==" {
			if rvLeft, err = r.parse(expr[:k], pos); err != nil {
				return nil, err
			}
			if rvRight, err = r.parse(expr[k+1:], pos); err != nil {
				return nil, err
			}
			if compareEqual(rvLeft, rvRight) {
				return &global.Structure{Tok: "BOOL", Lit: "true"}, nil
			}
			return &global.Structure{Tok: "BOOL", Lit: "false"}, nil
		}
	}

	return nil, nil
}

// compareEqual 比较符: 等于
func compareEqual(left, right *global.Structure) bool {
	var (
		lTok = left.Tok
		rTok = right.Tok
		lLit = left.Lit
		rLit = right.Lit
	)

	// 两边都是布尔值
	if lTok == "BOOL" && lTok == rTok {
		return lLit == rLit
	}
	// 两边都是字符串
	if lTok == "STRING" && lTok == rTok {
		return lLit == rLit
	}
	// 两边都是整型
	if lTok == "INT" && lTok == rTok {
		return lLit == rLit
	}

	l, err := formatValueTypeToCompare(left)
	if err != nil {
		return false
	}
	r, err := formatValueTypeToCompare(right)
	if err != nil {
		return false
	}

	// if left.Tok == "FLOAT" || right.Tok == "FLOAT" {
	// 	return l.(float64) == r.(float64)
	// }
	// return l.(int64) == r.(int64)

	return l == r
}

// formatValueTypeToCompare 进行弱类型转换用于比较运算
// 如 "1.1">1.1 正常返回
// "你好" > 1.1 将返回错误 不能将非数字类的字符串进行比较运算
func formatValueTypeToCompare(src *global.Structure) (interface{}, error) {
	var (
		v   interface{}
		err error
	)

	if src.Tok == "STRING" {
		if err = fn.ChangeTokTypeStringToTypeIntOrFloat(src); err != nil {
			return nil, err
		}
	} else if src.Tok == "BOOL" {
		if err = fn.ChangeBoolToInt(src); err != nil {
			return nil, err
		}
	}

	if src.Tok == "FLOAT" {
		if v, err = strconv.ParseFloat(src.Lit, 64); err != nil {
			return nil, err
		}
	}
	if src.Tok == "INT" {
		if v, err = strconv.ParseInt(src.Lit, 10, 64); err != nil {
			return nil, err
		}
	}
	return v, nil
}

// // compareNotEqual 比较符: 不等于
// func compareNotEqual(syn string, left, right *global.Structure) bool {
// 	var (
// 		lTok = left.Tok
// 		rTok = right.Tok
// 		lLit = left.Lit
// 		rLit = right.Lit
// 	)

// 	if left.Tok == right.Tok {
// 		return left.Lit != right.Lit
// 	}
// 	if lTok == "INT" && rTok == "FLOAT" {
// 		return lLit+".0" != rLit
// 	}
// 	if rTok == "INT" && lTok == "FLOAT" {
// 		return rTok+".0" != lLit
// 	}
// 	_, rvLeftBool := fn.ChangeToBool(left)
// 	_, rvRightBool := fn.ChangeToBool(right)

// 	return rvLeftBool != rvRightBool
// }

// // compareGT 比较符: 大于
// func compareGT(left, right *global.Structure) bool {
// 	var (
// 		lTok = left.Tok
// 		rTok = right.Tok
// 		lLit = left.Lit
// 		rLit = right.Lit
// 	)

// 	if lTok == "BOOL" {
// 		lTok = "INT"
// 		if lLit == "false" {
// 			lLit = "0"
// 		} else if lLit == "true" {
// 			lLit = "1"
// 		}
// 	}
// 	if rTok == "BOOL" {
// 		rTok = "INT"
// 		if rLit == "false" {
// 			rLit = "0"
// 		} else if rLit == "true" {
// 			rLit = "1"
// 		}
// 	}

// 	if left.Tok == right.Tok {
// 		return left.Lit != right.Lit
// 	}
// 	if lTok == "INT" && rTok == "FLOAT" {
// 		return lLit+".0" != rLit
// 	}
// 	if rTok == "INT" && lTok == "FLOAT" {
// 		return rTok+".0" != lLit
// 	}

// 	_, rvLeftBool := fn.ChangeToBool(left)
// 	_, rvRightBool := fn.ChangeToBool(right)
// 	return rvLeftBool != rvRightBool
// }
