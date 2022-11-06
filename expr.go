package lit

import (
	"fmt"
	"go/scanner"
	"go/token"
	"log"
	"strconv"
	"strings"

	fn "github.com/pywee/lit/function"
	"github.com/pywee/lit/global"
	"github.com/pywee/lit/types"
)

// FIXME.函数内如果表达式没有用分号分开 程序会报错

// cfn 代码体内的自定义函数
var cfn *fn.CustomFunctions

type Expression struct {
	funcBlocks     []*fn.FunctionInfo
	codeBlocks     []*global.Block
	publicVariable map[string]*global.Structure
}

func NewExpr(src []byte) (*Expression, error) {
	var (
		s      scanner.Scanner
		fset   = token.NewFileSet()
		expr   = make([]*global.Structure, 0, 100)
		result = &Expression{
			publicVariable: make(map[string]*global.Structure, 10),
		}
	)

	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		stok := tok.String()
		posString := fset.Position(pos).String()
		posLine := strings.Split(posString, ":")[0]
		if sLit := strings.ToLower(lit); stok != "STRING" && (sLit == "false" || sLit == "true") {
			stok = "BOOL"
		}
		expr = append(expr, &global.Structure{
			Tok:      stok,
			Lit:      lit,
			Position: posLine,
		})
	}

	_, err := result.parseExprs(expr, nil)
	return nil, err
}

// parseExprs 解析代码块
func (r *Expression) parseExprs(expr []*global.Structure, innerVariable map[string]*global.Structure) (*global.Structure, error) {
	var (
		err        error
		foundElse  bool
		rlen       = len(expr)
		blocks     = make([]*global.Block, 0, 20)
		funcBlocks = make([]*fn.FunctionInfo, 0, 20)
	)

	for i := 0; i < rlen; i++ {
		thisExpr := expr[i]
		if thisExpr.Tok == ";" && thisExpr.Lit == "\n" {
			continue
		}

		// 变量声明
		if expr[i].Tok == "IDENT" && i < rlen && expr[i+1].Tok == "=" {
			blocks, i = parseIdentedVAR(blocks, expr, i, rlen)
			continue
		}

		if thisExpr.Tok == "return" {
			var returnExpr = make([]*global.Structure, 0, 5)
			for j := i + 1; j < rlen; j++ {
				if expr[j].Tok == ";" {
					// blocks = append(blocks, &global.Block{
					// 	Type: types.CodeTypeIdentRETURN,
					// 	Code: returnExpr,
					// })
					// i = j
					// break
					return r.parse(returnExpr, "", innerVariable)
				}
				returnExpr = append(returnExpr, expr[j])
			}
			continue
		}

		// 函数定义
		if expr[i].Tok == "func" {
			if funcBlocks, i, err = parseIdentFUNC(funcBlocks, expr, i, rlen); err != nil {
				return nil, err
			}
			continue
		}

		// 函数调用
		if expr[i].Tok == "IDENT" && i < rlen && expr[i+1].Tok == "(" {
			blocks, i = parseFuncExec(blocks, expr, i, rlen)
			continue
		}

		// if 句子
		if expr[i].Tok == "if" && !foundElse {
			parsed, err := parseIdentedIF(blocks, expr, i, rlen)
			if err != nil {
				return nil, err
			}
			i = parsed.i
			blocks = parsed.blocks
			continue
		}

		// else 处理
		if expr[i].Tok == "else" {
			parsed, err := parseIdentELSE(blocks, expr, i, rlen)
			if err != nil {
				return nil, err
			}
			i = parsed.i
			blocks = parsed.blocks
			foundElse = parsed.foundElse
		}
	}

	r.codeBlocks = blocks
	r.funcBlocks = funcBlocks

	global.Output(blocks)

	for _, block := range blocks {
		if block.Type == types.CodeTypeFunctionExec {
			_, err := r.parse(block.Code, "", innerVariable)
			if err != nil {
				return nil, err
			}
		} else if block.Type == types.CodeTypeIdentIF {
			for _, v := range block.IfExt {
				global.Output(v.Condition)
				println("---")
			}
		} else if block.Type == types.CodeTypeIdentVAR {
			var vName string
			code := block.Code
			if vleft, vLeftListEndIdx := findStrInfrontSymbool("=", code); vLeftListEndIdx != -1 {
				if vLeftListEndIdx == 1 {
					vName = vleft[0].Lit
					code = code[vLeftListEndIdx+1:]
				}
			}
			if vName != "" {
				rv, err := r.parse(code, "", innerVariable)
				if err != nil {
					return nil, err
				}
				if innerVariable != nil {
					innerVariable[vName] = rv
				} else {
					r.publicVariable[vName] = rv
				}
			}
		}
	}
	return nil, nil
}

func (r *Expression) parse(expr []*global.Structure, pos string, innerVariable map[string]*global.Structure) (*global.Structure, error) {
	rLen := len(expr)
	if rLen == 0 {
		return nil, nil
	}

	// 变量解析
	if rLen == 1 {
		rv := expr[0]
		if rv != nil && rv.Tok == "IDENT" && global.IsVariableOrFunction(rv) {
			// 先寻找作用域变量 再找全局变量
			if innerRv, ok := innerVariable[rv.Lit]; ok {
				return innerRv, nil
			}
			if pubRv, ok := r.publicVariable[rv.Lit]; ok {
				return pubRv, nil
			}
			return nil, types.ErrorNotFoundVariable
		}
		return rv, nil
	}

	var (
		// err error
		err error
		// foundBracket 发现括号标记
		foundBracket bool
		// bracketCount 小括号计数器
		bracketCount int8
		// firstKey 收集前后代码使用的标记
		firstKey = -1
		// firstIdent 首个出现的定义代码
		firstIdent *global.Structure
	)

	for k, v := range expr {
		if v == nil {
			return nil, types.ErrorWrongSentence
		}
		if v.Tok == "if" {
			return nil, r.parseIf(expr, pos, innerVariable)
		}
		if v.Tok == "||" && firstKey == -1 {
			return r.parseOr(expr, k, pos, innerVariable)
		}
		if v.Tok == "&&" && firstKey == -1 {
			return r.parseAnd(expr, k, pos, innerVariable)
		}

		if v.Tok == "(" {
			bracketCount++
			if firstKey == -1 {
				firstKey = k
			}
			if firstIdent == nil && k > 0 && !foundBracket {
				firstIdent = expr[k-1]
			}
			foundBracket = true
		} else if v.Tok == ")" {
			bracketCount--
		}

		// first + middle + end
		// 主要递归 middle
		// a = (("wwww") + 222))
		// a = IsInt(("wwww")+222);
		// a = "你"+Replace("你好", "2", "3", 4)+"xxx";
		// FIXME a = (1+IsInt(1+(3))+(11+2)); fixed
		if bracketCount == 0 && foundBracket && firstKey != -1 {
			first := []*global.Structure{}
			if firstKey > 0 {
				// 发现函数
				// 执行 IsInt(1) 和 1+Isint(1) 时
				// 得到的 first 不一样
				first = expr[:firstKey]
				if firstIdent != nil && firstIdent.Tok == "IDENT" && expr[firstKey-1] == firstIdent {
					first = expr[:firstKey-1]
				}
			}

			// end 必须放在 middle 前面
			// 避免变量 expr 被修改了长度
			end := []*global.Structure{}
			if k < len(expr) && len(expr[k+1:]) > 0 {
				end = expr[k+1:]
			}

			// 发现中间表达式为函数执行调用
			var middle *global.Structure
			if global.IsVariableOrFunction(firstIdent) {
				funcName := firstIdent.Lit

				// 此判断在前面则可实现对内置函数的重写
				// global.Output(funcName)
				// if fni := cfn.GetCustomeFunc(funcName); fni != nil {
				if fni := getFuncIdented(r.funcBlocks, funcName); fni != nil {
					// 执行自定义函数
					// global.Output(funcName)
					// global.Output(expr[firstKey+1 : k])
					// fmt.Println(funcName, innerVariable, expr[firstKey+1:k])
					if middle, err = r.execCustomFunc(fni, expr[firstKey+1:k], pos, innerVariable); err != nil {
						return nil, err
					}
				}

				// 查找是否有内置函数
				if getFunc := fn.CheckFunctionName(funcName); getFunc != nil {
					// expr[firstKey+1 : k] 为实参
					// global.Output(expr[firstKey+1 : k])
					if middle, err = r.execFunc(funcName, expr[firstKey+1:k], pos, innerVariable); err != nil {
						return nil, err
					}
				}
			} else if middle, err = r.parse(expr[firstKey+1:k], pos, innerVariable); err != nil {
				return nil, err
			}

			// 左中右结合再次递归
			expr = append(first, middle)
			expr = append(expr, end...)

			return r.parse(expr, pos, innerVariable)
		}
	}

	// 进入这里的已经是最小粒度了 --------
	// 寻找变量值
	var (
		exists bool
		value  *global.Structure
	)
	for k, v := range expr {
		if v.Tok == "IDENT" && global.IsVariableOrFunction(v) {
			if value, exists = innerVariable[v.Lit]; exists {
				expr[k] = value
				continue
			}
			if value, exists = r.publicVariable[v.Lit]; exists {
				expr[k] = value
				continue
			}
			return nil, types.ErrorNotFoundVariable
		}
	}

	// 比较运算符处理
	// 如果有 则进入该逻辑
	if rLen = len(expr); rLen > 1 {
		var rv *global.Structure
		if rv, err = r.parseCompare(expr, pos); err != nil {
			return nil, err
		}
		if rv != nil {
			return rv, nil
		}
	}

	// global.Output(expr)

	// 最小粒度进入到算术表达式中计算
	rv, err := r.parseExpr(expr, pos)
	if err != nil {
		return nil, err
	}
	return rv[0], nil
}

// getFuncIdented 获取定义好的自定义函数代码块
func getFuncIdented(blocks []*fn.FunctionInfo, fname string) *fn.FunctionInfo {
	for _, v := range blocks {
		if v.FunctionName == fname {
			return v
		}
	}
	return nil
}

// parseExpr 解析算术表达式入口
// 1+11|23-12;
// 1+2+(3%10+(11*22+1.1-10)+41+(5+9)+10)+1+(2-(10*2%2));
// 2+1/(20-1*(3+10)-1)%2^1;
func (r *Expression) parseExpr(expr []*global.Structure, pos string) ([]*global.Structure, error) {
	v, err := r.findExprK(expr, pos)
	if err != nil {
		return nil, err
	}
	return parsePlusReduceMulDivB(v, pos)
}

func (r *Expression) findExprK(expr []*global.Structure, pos string) ([]*global.Structure, error) {
	exprLen := len(expr)

	// 如果返回的是最终值 则不再需要进一步解析了
	// 加快处理速度
	if exprLen == 1 {
		if e0 := expr[0]; inArray(e0.Tok, []string{"INT", "STRING", "FLOAT", "BOOL"}) {
			return expr, nil
		}
	}

	// 错误表达式处理
	if exprLen > 1 {
		if e0Tok := expr[0].Tok; (e0Tok == "FLOAT" || e0Tok == "INT") && expr[1].Tok == "IDENT" {
			return nil, types.ErrorWrongSentence
		}
	}

	for k, v := range expr {
		if global.IsVariableOrFunction(v) {
			ret, err := r.Get(v.Lit)
			if err != nil {
				return nil, err
			}
			expr[k] = ret
		}
	}

	var (
		startI   int
		endI     int
		startIdx = -1
		endIdx   = -1
	)
	for k, v := range expr {
		if v.Tok == "(" {
			startI++
			if startIdx == -1 {
				startIdx = k
			}
		}
		if v.Tok == ")" {
			endI++
			endIdx = k
		}
		if endI == startI && startIdx != -1 && endIdx != -1 {
			temp, err := r.findExprK(expr[startIdx+1:endIdx], pos)
			if err != nil {
				return nil, err
			}
			result := expr[:startIdx]
			result = append(result, temp...)
			return r.findExprK(append(result, expr[endIdx+1:]...), pos)
		}
	}

	return parsePlusReduceMulDivB(expr, pos)
}

// parsePlusReduceMulDivB 处理加减乘除
func parsePlusReduceMulDivB(arr []*global.Structure, pos string) ([]*global.Structure, error) {
	result, err := parsePlusReduceMulDiv([]string{"*", "/", "%"}, arr)
	if err != nil {
		return nil, err
	}
	if len(result) == 1 {
		return result, nil
	}

	if result, err = parsePlusReduceMulDiv([]string{"&"}, result); err != nil {
		return nil, err
	}
	if len(result) == 1 {
		return result, nil
	}

	if len(result) > 0 {
		if result, err = parsePlusReduceMulDiv([]string{"+", "-", "|", "^"}, result); err != nil {
			return nil, err
		}
	}

	if len(result) != 1 {
		return nil, types.ErrorWrongSentence
	}
	return result, nil
}

// parsePlusReduceMulDiv 递归解析表达式 最内层
func parsePlusReduceMulDiv(arr []string, expr []*global.Structure) ([]*global.Structure, error) {
	if len(expr) == 0 {
		return expr, nil
	}

	var result []*global.Structure
	for k, v := range expr {
		if inArray(v.Tok, arr) && k-1 > -1 && k+2 <= len(expr) {
			r, err := findExprBetweenSymbool(expr[k-1], expr[k], expr[k+1])
			if err != nil {
				return nil, err
			}
			if r.Type == "STRING" {
				result = append([]*global.Structure{{Tok: "STRING", Lit: r.Value.(string)}}, expr[k+2:]...)
				return parsePlusReduceMulDiv(arr, result)
			}

			middle := &global.Structure{}
			if r.Type == "FLOAT" {
				// FIXME
				// 待观察是否对其他表达式有影响
				// example:
				// a = 2
				// b = 3
				// c = a + 100/10*2.0
				// output FLOAT 22.
				// 在使用 strings.TrimRight(lit, ".") 之后
				// 输出 FLOAT(22)

				middle.Tok = "FLOAT"
				middle.Lit = strings.TrimRight(fmt.Sprintf("%f", r.Value), "0")
				middle.Lit = strings.TrimRight(middle.Lit, ".")
			} else if r.Type == "INT" {
				middle.Tok = "INT"
				middle.Lit = fmt.Sprintf("%d", r.Value)
			}

			result = append(expr[:k-1], middle)
			result = append(result, expr[k+2:]...)
			return parsePlusReduceMulDiv(arr, result)
		}
	}
	return expr, nil
}

// findExprBetweenSymbool 查找某运算符数据
// 如: 1+2-3 输出查找 + 号 则返回 1+2
func findExprBetweenSymbool(l, m, r *global.Structure) (*exprResult, error) {
	var (
		exprTyp string
		left    interface{}
		right   interface{}
		lTok    = l.Tok
		mTok    = m.Tok
		rTok    = r.Tok
	)

	// 弱类型处理
	// 布尔值转整型
	if l.Tok == "BOOL" {
		if l.Lit == "true" {
			lTok = "INT"
			l.Tok = "INT"
			l.Lit = "1"
		} else if l.Lit == "false" {
			lTok = "INT"
			l.Tok = "INT"
			l.Lit = "0"
		}
	}
	if r.Tok == "BOOL" {
		if r.Lit == "true" {
			rTok = "INT"
			r.Tok = "INT"
			r.Lit = "1"
		} else if r.Lit == "false" {
			rTok = "INT"
			r.Tok = "INT"
			r.Lit = "0"
		}
	}

	// 字符串拼接及弱类型处理的算术计算
	if (lTok == "STRING" || lTok == "INTERFACE") && (rTok == "STRING" || rTok == "INTERFACE") {
		// 弱类型处理 如果左右两边都是字符串数字则允许进行算术计算
		isLeftNumeric, err := global.IsNumber(l.Lit)
		if err != nil {
			log.Print(err)
		}
		isRightNumeric, err := global.IsNumber(r.Lit)
		if err != nil {
			log.Print(err)
		}

		// 不是完全由数字构成的字符串则直接拼接
		// 否则转换类型进行算术计算
		// example A
		// a = "1";
		// b = "2";
		// print(a+b) 输出 (int)3

		// example B
		// a = "1";
		// b = "2a";
		// print(a+b) 输出 (string)12a

		if !isLeftNumeric || !isRightNumeric {
			if mTok == "+" {
				return &exprResult{Type: "STRING", Value: l.Lit + r.Lit}, nil
			}
			return nil, types.ErrorNonNumberic
		}
	}

	// 弱类型处理
	// 数字字符串转为整型或浮点型
	if lTok == "STRING" || lTok == "CHAR" {
		lit := l.Lit
		isFloat, err := global.IsFloat(lit)
		if err != nil {
			return nil, types.ErrorWrongSentence
		}
		if isFloat {
			if _, err = strconv.ParseFloat(lit, 64); err != nil {
				return nil, types.ErrorWrongSentence
			}
			lTok = "FLOAT"
			l.Lit = lit
		}

		isInt, err := global.IsInt(lit)
		if err != nil {
			return nil, types.ErrorWrongSentence
		}
		if isInt {
			if _, err = strconv.ParseInt(lit, 10, 64); err != nil {
				return nil, types.ErrorWrongSentence
			}
			lTok = "INT"
			l.Lit = lit
		}
		if !isFloat && !isInt {
			print(10)
			return nil, types.ErrorWrongSentence
		}
	}
	if rTok == "STRING" || rTok == "CHAR" {
		_, err := strconv.ParseInt(r.Lit, 10, 64)
		if err != nil {
			return nil, types.ErrorWrongSentence
		}
		rTok = "INT"
	}

	if lTok == "INT" && rTok == "INT" {
		exprTyp = "INT"
		left, _ = strconv.ParseInt(l.Lit, 10, 64)
		right, _ = strconv.ParseInt(r.Lit, 10, 64)
	}
	if lTok == "FLOAT" || rTok == "FLOAT" {
		exprTyp = "FLOAT"
		left, _ = strconv.ParseFloat(l.Lit, 64)
		right, _ = strconv.ParseFloat(r.Lit, 64)
	}

	if mTok == "*" {
		if exprTyp == "INT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) * right.(int64)}, nil
		}
		if exprTyp == "FLOAT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) * right.(float64)}, nil
		}
	}
	if mTok == "/" {
		if exprTyp == "INT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) / right.(int64)}, nil
		}
		if exprTyp == "FLOAT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) / right.(float64)}, nil
		}
	}
	if mTok == "+" {
		if exprTyp == "INT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) + right.(int64)}, nil
		}
		if exprTyp == "FLOAT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) + right.(float64)}, nil
		}
	}
	if mTok == "-" {
		if exprTyp == "INT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) - right.(int64)}, nil
		}
		if exprTyp == "FLOAT" {
			return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(float64) - right.(float64)}, nil
		}
	}
	if mTok == "%" && exprTyp == "INT" {
		return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) % right.(int64)}, nil
	}
	if mTok == "^" && exprTyp == "INT" {
		return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) ^ right.(int64)}, nil
	}
	if mTok == "|" && exprTyp == "INT" {
		return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) | right.(int64)}, nil
	}
	if mTok == "&" && exprTyp == "INT" {
		return &exprResult{Type: exprTyp, Tok: mTok, Value: left.(int64) & right.(int64)}, nil
	}
	return nil, types.ErrorWrongSentence
}

// findStrInfrontSymbool 查找指定符号之前的数据
func findStrInfrontSymbool(sym string, src []*global.Structure) ([]*global.Structure, int) {
	for k, v := range src {
		if sym == v.Tok {
			return src[:k], k
		}
	}
	return nil, -1
}

func inArray(sep string, arr []string) bool {
	for _, v := range arr {
		if sep == v {
			return true
		}
	}
	return false
}

// isSingularExpr 判断当前是否不需要再进一步解析
// func isSingularExpr(expr *global.Structure) bool {
// 	t := expr.Tok
// 	return t == "STRING" || t == "INT" || t == "BOOL" || t == "FLOAT"
// }

// func NewExpr2(src []byte) (*Expression, error) {
// 	var (
// 		funcKuo       int8
// 		foundCustFunc bool // 发现自定义函数时 保存其形式文本
// 		s             scanner.Scanner
// 		fset          = token.NewFileSet()
// 		list          = make([]*global.Structure, 0, 20)
// 		result        = &Expression{publicVariable: make(map[string]*global.Structure, 10)}
// 	)

// 	file := fset.AddFile("", fset.Base(), len(src))
// 	s.Init(file, src, nil, scanner.ScanComments)

// 	// TODO
// 	// 提前获取结构体定义部分
// 	// 提前获取常量定义部分

// 	// 提前获取函数定义部分
// 	if err := getFunctionDefined(s, file, fset); err != nil {
// 		return nil, err
// 	}

// 	var (
// 		foundIf        bool
// 		ifCurlyBracket int8
// 	)
// 	s.Init(file, src, nil, scanner.ScanComments)
// 	for {
// 		pos, tok, lit := s.Scan()
// 		if tok == token.EOF {
// 			break
// 		}

// 		stok := tok.String()
// 		posString := fset.Position(pos).String()
// 		// posInt, _ := strconv.ParseInt(strings.Split(posString, ":")[1], 10, 64)
// 		posLine := strings.Split(posString, ":")[0]

// 		// 跳过文本内所有自定义函数
// 		if stok == "func" {
// 			foundCustFunc = true
// 		}
// 		if stok == "if" && !foundIf && !foundCustFunc {
// 			foundIf = true
// 		}
// 		if stok == "else" && !foundIf && !foundCustFunc {
// 			foundIf = true
// 		}

// 		// fmt.Println(tok, lit, foundIf)

// 		// 标记 if 语句代码块
// 		if foundIf {
// 			if stok == "{" {
// 				ifCurlyBracket++
// 			} else if stok == "}" {
// 				ifCurlyBracket--
// 				if ifCurlyBracket == 0 {
// 					foundIf = false
// 				}
// 			}
// 		}

// 		// 标记函数定义代码块
// 		if foundCustFunc {
// 			if stok == "{" {
// 				funcKuo++
// 			} else if stok == "}" {
// 				funcKuo--
// 				if funcKuo == 0 {
// 					foundCustFunc = false
// 				}
// 			}
// 			continue
// 		}

// 		if !foundIf && stok == ";" {
// 			var vName string
// 			if vleft, vLeftListEndIdx := findStrInfrontSymbool("=", list); vLeftListEndIdx != -1 {
// 				if vLeftListEndIdx == 1 {
// 					vName = vleft[0].Lit
// 					list = list[vLeftListEndIdx+1:]
// 				}
// 			}

// 			// 递归解析表达式
// 			for _, v := range list {
// 				if sLit := strings.ToLower(v.Lit); stok != "STRING" && (sLit == "false" || sLit == "true") {
// 					v.Tok = "BOOL"
// 				}
// 			}

// 			rv, err := result.parse(list, posLine, nil)
// 			if err != nil {
// 				return nil, errors.New("第" + posLine + "行: " + err.Error())
// 			}
// 			// 变量赋值
// 			if vName != "" {
// 				// fmt.Printf("set %s to %v\n", vName, rv)
// 				result.publicVariable[vName] = rv
// 			}
// 			list = nil
// 			continue
// 		}

// 		// fmt.Println(tok, lit, foundIf)

// 		if stok == "CHAR" {
// 			stok = "STRING"
// 		}
// 		if stok == "STRING" {
// 			lit = formatString(lit)
// 		}
// 		list = append(list, &global.Structure{
// 			Position: posString,
// 			Tok:      stok,
// 			Lit:      lit,
// 		})
// 	}
// 	return result, nil
// }
