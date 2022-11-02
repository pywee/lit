package example

import (
	goExpr "github.com/pywee/lit"
)

// SetVal 变量声明及调用
// 执行算术表达式
func SetVal() {
	exprs := []byte(`
		a = 10;
		b = a + 20 / (10-(2+1*(3*a))-1+2);
		print(b);
	`)
	if _, err := goExpr.NewExpr(exprs); err != nil {
		panic(err)
	}
}

// RunInnerFunction 内置函数调用
func RunInnerFunction() {
	exprs := []byte(`
		a = "hello 123";
		a = replace(a, "123", "456...", -1);
		print(trim(a, "."));
	`)
	if _, err := goExpr.NewExpr(exprs); err != nil {
		panic(err)
	}
}

// RunCustomFunction 自定义函数声明及调用
func RunCustomFunction() {
	exprs := []byte(`
		func demo(a, b = 10) {
			c = 0;
			return a + 20 + b > 10 + c;
		}
		a = demo(10);
		print(a);
	`)
	if _, err := goExpr.NewExpr(exprs); err != nil {
		panic(err)
	}
}

// RunCustomBool 弱类型处理
func RunCustomBool() {
	exprs := []byte(`
		a = false;
		b = true;
		c = "123";
		d = 456;
		print(a + b);
		print(a >= b);
		print(a + b + c + d);
	`)
	if _, err := goExpr.NewExpr(exprs); err != nil {
		panic(err)
	}
}
