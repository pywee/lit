# goExpr

Scientific computing in golang.
本代码实现了在 Golang 里面直接引用后可以对算术表达式的文本进行解释，科学地计算出结果。包括加、减、乘、除、与、或、非。

对于需要进行解释型动态开发的需求可以使用本代码包。

### 使用方法


```
go get github.com/pywee/goExpr
```

示例：

```
import "github.com/pywee/goExpr"

func main() {

    // Golang 语言原生计算 
    // output 
    // +1.600000e+001
    // 125
	println((2 + 100 ^ 2 - (10*1.1 - 22 + (22 | 11))) / 10 * 2)
    println(12/333+31+(5/10)-6|100)


    // 使用 goExpr 计算文本中的数据
    // 表达式文本
	exprs := []byte(`
        a = (2 + 100^2-(10*1.1-22+(22|11)))/10*2;
        b = a+11203-(11*10&2/16);
        c = 12/333+31+(5/10)-6|100;
    `)

	exec, err := goExpr.NewExpr(exprs)
    valueA := exec.Get("a")
	valueB := goExpr.Get("b")
	valueC := goExpr.Get("c")
    println(valueA)
    println(valueB)
    println(valueC)

    // output
    // &{ FLOAT 16} <nil>
    // &{ FLOAT 11219} <nil>
    // &{ INT 125} <nil>
}
```

---

####  请注意，goExpr 优先使用 Golang 的算术符号优先级进行数据计算。
**每个语言对算术符号的优先级处理都有一定区别，如，针对以下表达式进行计算时：**

``` golang
// 2 + 100 ^ 2 - (10*1.1 - 22 + (22 | 11)) / 10 * 2
// php 输出 -104
// node.js 输出 -104
// golang 输出 96
// goExpr 输出 96
```

---

##### goExpr 算术符号优先级
第一级  ``` * / % ```
第二级  ``` + - | & ^ ```
