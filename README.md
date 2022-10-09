# goExpr

Scientific computing in golang.

本代码实现了在 Golang 里面直接引用后可以对算术表达式的文本进行解释，科学地计算出结果。包括加、减、乘、除、与、或、非。

对于需要进行解释型动态开发的需求可以使用本代码包。

### 使用方法


```
go get github.com/pywee/goExpr
```

示例：

```golang
import "github.com/pywee/goExpr"

func main() {

    // Golang 语言原生计算 
    println((2 + 100 ^ 2 - (10*1.1 - 22 + (22 | 11))) / 10 * 2)
    println(12/333+31+(5/10)-6|100)
    // output 
    // +1.600000e+001
    // 125


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

**支持简单的变量操作**
```golang
	src := []byte(`
        a = 123;
        b = a + 456
    `)
    exec, err := goExpr.NewExpr(src)
    fmt.Println(expr.GetVal("b"))

    // output
    // &{ INT 579} <nil>
```
---

**支持针对字符串类型的整型或浮点型参与算术运算**

```golang
    // 浮点型字符串+整型
    // 最终输出结果的底层类型将变为浮点型
    // example:
    // a = '333.910'+0.01
	src := []byte(`a = '333.910'+0.01`)
    exec, err := goExpr.NewExpr(src)
    fmt.Println(expr.GetVal("a"))
    // output
    // &{ FLOAT 333.92} <nil>

    --------------------------------------

    // 其他字符串+整型将会报错
    // example
    // a = "abcwwww1230"+0.01
	src := []byte(`a = "abcwwww1230"+0.01`)
    exec, err := goExpr.NewExpr(src)
    fmt.Println(expr.GetVal("a"))
    // output 
    // found error (code 1002), notice: xx:xx wrong sentence

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
第一级  ``` () ```
第二级  ``` * / % ```
第三级  ``` + - | & ^ ```
