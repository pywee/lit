package goExpr

import (
	"github.com/pywee/goExpr/global"
)

func test(expr []*global.Structure, pos string) {
	i := 0
	foundK := -1
	kList := make([]*global.Structure, 0, 10)
	for k, v := range expr {
		if v.Tok == "(" {
			if foundK == -1 {
				foundK = k
			}
			i++
		}
		if v.Tok == ")" {
			i--
		}
		if foundK >= 0 {
			kList = append(kList, v)
		}
		if i == 0 {
			global.Output(kList)
		}
	}
}
