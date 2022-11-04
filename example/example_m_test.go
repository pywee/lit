package example

import (
	"testing"

	goExpr "github.com/pywee/lit"
)

// TestSetVal 变量声明及调用
// 执行算术表达式
func TestSetVal(t *testing.T) {
	exprs := []byte(`
		a = 10;
		b = a + 20 / (10-(2+1*(3*a))-1+2);
		print(b);
	`)
	_, err := goExpr.NewExpr(exprs)
	if err != nil {
		panic(err)
	}
}

// TestRunInnerFunction 内置函数调用
func TestRunInnerFunction(t *testing.T) {
	exprs := []byte(`
		a = "hello 123";
		a = replace(a, "123", "456...", -1);
		print(trim(a));
	`)
	if _, err := goExpr.NewExpr(exprs); err != nil {
		panic(err)
	}
}

// TestRunCustomFunction 自定义函数声明及调用
func TestRunCustomFunction(t *testing.T) {
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

// TestRunCustomBool 弱类型处理
func TestRunCustomBool(t *testing.T) {
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
