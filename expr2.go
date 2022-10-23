package goExpr

// func test(expr []*global.Structure, pos string) {
// 	i := 0
// 	foundK := -1
// 	kList := make([]*global.Structure, 0, 10)
// 	for k, v := range expr {
// 		if v.Tok == "(" {
// 			if foundK == -1 {
// 				foundK = k
// 			}
// 			i++
// 		}
// 		if v.Tok == ")" {
// 			i--
// 		}
// 		if foundK >= 0 {
// 			kList = append(kList, v)
// 		}
// 		if i == 0 {
// 			global.Output(kList)
// 		}
// 	}
// }

// parse 解析器主入口
// func (r *Expression) parseB(expr []*global.Structure, pos string, foundAndOr bool) (*global.Structure, error) {
// 	var err error
// 	rLen := len(expr)
// 	if rLen == 0 {
// 		return nil, nil
// 	}

// 	e0 := expr[0]
// 	if rLen == 1 && e0.Tok != "IDENT" {
// 		return e0, nil
// 	}

// 	// 只找括号不换函数
// 	i := 0
// 	foundK := -1
// 	kList := make([]*global.Structure, 0, 10)

// 	for k, v := range expr {
// 		if v.Tok == "(" {
// 			if foundK == -1 {
// 				foundK = k
// 			}
// 			i++
// 		}
// 		if v.Tok == ")" {
// 			i--
// 		}
// 		if foundK >= 0 {
// 			kList = append(kList, v)
// 		}

// 		// 括号 非函数
// 		if foundK >= 0 && i == 0 {
// 			if foundK == 0 || expr[foundK-1].Tok != "IDENT" {
// 				// Fixed
// 				// VarDump((1)+(2)) 此时会进入死循环 因为每次都取到 (1）
// 				var rv *global.Structure
// 				if kLen := len(kList); kLen > 1 && kList[0].Tok == "(" && kList[kLen-1].Tok == ")" {
// 					rv, err = r.parseB(kList[1:kLen-1], pos, foundAndOr)
// 					if err != nil {
// 						return nil, err
// 					}
// 				} else if rv, err = r.parseB(kList, pos, foundAndOr); err != nil {
// 					return nil, err
// 				}

// 				k1 := expr[k+1:]
// 				expr = append(expr[:foundK], rv)
// 				expr = append(expr, k1...)
// 				foundK = -1
// 				return r.parseB(expr, pos, foundAndOr)
// 			}

// 			// 在括号的结尾发现函数
// 			if kLen := len(kList); kLen > 1 && kList[0].Tok == "(" && kList[kLen-1].Tok == ")" {
// 				if rv, _ := r.doFunc(expr, pos, foundAndOr); rv != nil {
// 					return rv, nil
// 				}
// 			}
// 		}
// 	}

// 	// FIXME
// 	// 临时有符号的整型和浮点型的处理逻辑
// 	if rLen == 2 && (e0.Tok == "-" || e0.Tok == "+") {
// 		// TODO
// 		// 考虑弱类型支持 如 "-100" 和 -100 是否等义
// 		if e1 := expr[1]; e1.Tok == "INT" || e1.Tok == "FLOAT" {
// 			if e0.Tok == "+" {
// 				return &global.Structure{Tok: e1.Tok, Lit: e1.Lit}, nil
// 			}
// 			return &global.Structure{Tok: e1.Tok, Lit: e0.Tok + e1.Lit}, nil
// 		}
// 	}

// 	// FIXME
// 	// 针对有 && 和 || 符号的表达式才进入这个逻辑
// 	// 后期可以进一步优化
// 	if foundAndOr {
// 		rvAfter, err := r.parseAndOr(expr, pos, foundAndOr)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if rvAfter != nil {
// 			return rvAfter, nil
// 		}
// 	}

// 	// 函数判断
// 	if rLen >= 3 {
// 		rv, err := r.doFunc(expr, pos, foundAndOr)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if rv != nil {
// 			return rv, nil
// 		}
// 	}

// 	// 只找函数不找括号
// 	count := 0
// 	found := false
// 	foundFNKey := -1
// 	fns := []*global.Structure{}
// 	for k, v := range expr {
// 		if v.Tok == "(" {
// 			if foundFNKey == -1 && k > 0 && expr[k-1].Tok == "IDENT" {
// 				foundFNKey = k
// 				fns = append(fns, expr[k-1])
// 			}
// 			found = true
// 			count++
// 		} else if v.Tok == ")" {
// 			count--
// 		}
// 		if count > 0 {
// 			fns = append(fns, v)
// 		}
// 		if count == 0 && found {
// 			if foundFNKey != -1 {
// 				// 补全最后一个括号
// 				fns = append(fns, v)
// 				rv, err := r.parseB(fns, pos, foundAndOr)
// 				if err != nil {
// 					return nil, err
// 				}
// 				right := []*global.Structure{}
// 				if k+1 < len(expr) {
// 					right = expr[k+1:]
// 				}
// 				if foundFNKey > 1 {
// 					expr = append(expr[:foundFNKey-1], rv)
// 				} else {
// 					expr = []*global.Structure{rv}
// 				}
// 				expr = append(expr, right...)
// 				return r.parseB(expr, pos, foundAndOr)
// 			}
// 		}
// 	}

// 	// 找出剩余的未完成的函数
// 	// 如 true+isInt(1)+isFloat(1.1)
// 	// 如果不在此处进行检查 那么当前函数只会解析前面的 isInt
// 	// 而后面的 isFloat 会丢失
// 	// 这与 isInt(1)+isFloat(1.1)+true 的处理逻辑不一样
// 	// 所以这里还需要一次递归处理
// 	// foundLastFuncExpr := false
// 	for k, v := range expr {
// 		if v.Tok == "IDENT" && k+1 < len(expr) && expr[k+1].Tok == "(" {
// 			rv, err := r.parseB(expr, pos, foundAndOr)
// 			if err != nil {
// 				return nil, err
// 			}
// 			expr = append(expr[:k], rv)
// 		}
// 	}

// 	// FIXME
// 	if len(expr) == 0 {
// 		println("表达式可能有误！！！")
// 		return nil, nil
// 	}

// 	rv, err := r.parseExpr(expr, pos)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return rv[0], nil
// }

// func (r *Expression) doFunc(expr []*global.Structure, pos string, foundAndOr bool) (*global.Structure, error) {
// 	// 判断是否表达式为函数
// 	// 如果表达式是 replace("1", "2", "", 1) 则可生效
// 	// FIXME 如果表达式是 replace("1", "2", "", 1) + "xxxx" 则不生效 fixed
// 	rLen := len(expr)
// 	if fn.IsExprFunction(expr, rLen) {
// 		if global.IsVariableOrFunction(expr[0]) {
// 			funcName := expr[0]
// 			fArgs, err := fn.CheckFunctionName(funcName.Lit)
// 			if err != nil {
// 				return nil, err
// 			}

// 			// 函数内参数检查
// 			// 获取传入执行函数的具体参数
// 			// 并将它们的结果值递归解析出来
// 			args := fn.GetFunctionArgList(expr[2 : rLen-1])
// 			argsLen := len(args)
// 			if fArgs.MustAmount > argsLen {
// 				return nil, types.ErrorArgsNotEnough
// 			}
// 			if fArgs.MaxAmount != -1 && fArgs.MaxAmount < argsLen {
// 				return nil, types.ErrorTooManyArgs
// 			}

// 			// FIXME
// 			// get params after parsing
// 			// 汇总解析成功之后的实参数据
// 			// 传入回调函数进行实际执行
// 			// 当前仅支持内置函数
// 			var paramsList []*global.Structure
// 			for k, varg := range args {
// 				// FIXME
// 				// 函数中的实参表达式 实参可以是函数、变量、算术表达式等等
// 				rv, err := r.parseB(varg, pos, foundAndOr)
// 				if err != nil {
// 					return nil, err
// 				}

// 				// 检查最终解析出来的参数值类型是否与函数要求的形参类型一致
// 				if fa := fArgs.Args[k]; fa.Type != types.INTERFACE && fa.Type != rv.Tok {
// 					// TODO
// 					// 参数[弱类型]支持
// 					// 参数[提前在形参中设置默认值]支持
// 					// fmt.Println(fa.Type, rv.Tok)
// 					return nil, types.ErrorArgsNotSuitable
// 				}
// 				paramsList = append(paramsList, rv)
// 			}
// 			fRet, err := fArgs.FN(pos, paramsList...)
// 			return fRet, err
// 		}
// 	}
// 	return nil, nil
// }

// parseAndOr
// func (r *Expression) parseAndOr(expr []*global.Structure, pos string, foundAndOr bool) (*global.Structure, error) {
// 	// FIXME 针对 && 符号的解析
// 	// 优先处理括号
// 	// 1.针对已经声明的布尔值没有处理正确
// 	// example false && 12345;
// 	// 2.使用函数的时候 在带有 && 符号语句中没有解析出正确结果
// 	foundK := false
// 	for k, v := range expr {
// 		if v.Tok == "(" {
// 			foundK = true
// 		}
// 		if v.Tok == "&&" && len(expr) >= 3 && k > 0 {
// 			if foundK {
// 				return nil, nil
// 			}
// 			rvLeft, err := r.parseB(expr[:k], pos, foundAndOr)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if fn.ChangeBool(rvLeft).IsBoolFalse() {
// 				return rvLeft, nil
// 			}
// 			rvRight, err := r.parseB(expr[k+1:], pos, foundAndOr)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return fn.ChangeBool(rvRight), nil
// 		}

// 		if v.Tok == "||" && len(expr) >= 3 && k > 0 {
// 			rvLeft, err := r.parseB(expr[:k], pos, foundAndOr)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if fn.ChangeBool(rvLeft).IsBoolTrue() {
// 				return rvLeft, nil
// 			}
// 			rvRight, err := r.parseB(expr[k+1:], pos, foundAndOr)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return fn.ChangeBool(rvRight), nil
// 		}
// 	}
// 	return nil, nil
// }

// func (r *Expression) FindFunction(expr []*global.Structure, pos string, foundAndOr bool) ([]*global.Structure, error) {
// 	count := 0
// 	found := false
// 	foundFNKey := -1
// 	fns := []*global.Structure{}
// 	for k, v := range expr {
// 		if v.Tok == "(" {
// 			if foundFNKey == -1 && k > 0 && expr[k-1].Tok == "IDENT" {
// 				foundFNKey = k
// 				fns = append(fns, expr[k-1])
// 			}
// 			found = true
// 			count++
// 		} else if v.Tok == ")" {
// 			fns = append(fns, v)
// 			count--
// 		}
// 		if count > 0 {
// 			fns = append(fns, v)
// 		}
// 		if count == 0 && found {
// 			if foundFNKey != -1 {
// 				rv, err := r.parseB(fns, pos, foundAndOr)
// 				if err != nil {
// 					return nil, err
// 				}
// 				if foundFNKey > 1 {
// 					expr = append(expr[:foundFNKey-1], rv)
// 				}
// 				if k+1 < len(expr) {
// 					expr = append(expr, expr[k+1:]...)
// 				}
// 			}
// 			break
// 		}
// 	}
// 	return expr, nil
// }
