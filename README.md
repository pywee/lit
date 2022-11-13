# 基于 Golang 编写的解释型语言 (迭代中)

这是一个基于纯 Golang 开发的解释型编程语言，目前仍在持续开发中，还没有<b>正式的名字</b>，暂且叫作 <b>Lit</b>。

最终我想把它实现为一个解释型的弱类型编程语言，希望它能用于日常的 Web 开发。

关于目前已实现的特性，请看如下文档示例，我将会定期更新。

当前文档更新日期是: 2022.11.13

---

### 已支持的特性

**一、变量**

**二、算术表达式**

**三、弱类型**

**四、关系运算**

**五、逻辑运算**

**六、内置函数**

**七、自定义函数**

**八、流程控制 if**

**九、流程控制 for (实验阶段，暂未包含 continue, break关键字)**

---



### 使用方法 非常简单

##### 第一步：下载 lit 项目包

```
go get github.com/pywee/lit
```
---

##### 第二步：在当前目录新建 main.go 文件，内容如下

```golang
import (
	"fmt"
	"os"
	"github.com/pywee/lit"
)

func main() {
	exprs, err := os.ReadFile("./src.lit")
	if err != nil {
		panic(err)
	}
	if _, err = goExpr.NewExpr(exprs); err != nil {
		fmt.Println(err)
		return
	}
}

```

---

##### 第三步：在当前目录新建 src.lit 文件

```go

// src.lit 文件内容

// 例1
// 变量声明、进行算术运算、执行内置函数
a = 5;
b = a + 6;
c = a + b * (1 + 2 * (3 + 4) - 5 / b);
print(a + b + c); // 输出 186


// 例2
// 较为复杂的算术运算，优先级的处理逻辑与 golang 保持一致
// * 以下代码如果是在 PHP 中执行会输出 -24，两种语言针对优先级的处理略微不一样，其他语言也是如此
a = (2 + 100 ^ 2 - (10*1.1 - 22 + (22 | 11))) / 10 * 2;
b = 12 / 333 + 31 + (5 / 10) - 6 | 100;
print(a); // 输出 16
print(b); // 输出 125


// 例3
// 自定义函数以及弱类型处理
func demo(arg) {
    arg++;
    return arg-false*3/10+(20-1);
}
print(demo(10)+70); // 100


// 例4
// 下面的句子调用了两个函数 
// isInt(arg) 用来检查 arg 是否为整型 
// replace(arg1, arg2, arg3, arg4) 用来做字符串替换
a = replace("Hello Word111!", "1", "", 2-isInt((1+(1 + isInt(123+(1+2)))-1)+2)-2);
print(a); // 输出 Hello word!


// 例5
// 更多关于弱类型的处理，基本上参考了 PHP
a = true - 1;
b = isInt(1);
c = isFloat(1.0);
d = false == 0.0;
e = "false" == 0.0;
print(a); // 0
print(b); // true
print(c); // true
print(d); // true
print(e); // true

a = false;
b = true;
c = "123";
d = 456;
print(a + b); // 1
print(a >= b); // false
print(a + b + c + d); // 580


// 例6
// 流程控制语句 if 多if嵌套示例
a = 100;
func functionDemo() {
    return 1;
}
if a == 101 - functionDemo() {
    if 1+1 == 2 && 3-1 == 20 {
        print("yes");
    } else {
        if 1+a == 102{
            print("no");
        } else {
            if false {
                print(0)
            } else {
                print("Hello world!")
            }
        }
    }
}


// 例7
// 流程控制语句 for
a = 5;
if a == 6-demo() {
    for i = 0; i < demo()+a; i ++ {
        print(a+i); // 输出 5 6 7 8 9 10
    }
}
func demo() {
    return 1;
}

// * for 功能仍在迭代中，仅实现了一种
// * 我在完成对数组的支持后会继续迭代 for 功能，实现 range 操作

```
```golang
// 最后即可在当前目录执行 go run . 输出结果
// go run .
```

---
**请注意，Lit 的算术符号优先级向 Golang 看齐。每个语言对算术符号的优先级处理都有一定区别，如，针对以下表达式进行计算时：**

``` golang
// 2 + 100 ^ 2 - (10*1.1 - 22 + (22 | 11)) / 10 * 2
// PHP 输出 -104
// Node.js 输出 -104
// Golang 输出 96
// Lit 输出 96
```

---

**Lit 算术符号优先级**

第一级  ``` () && || ```

第二级  ``` > < >= <= == != ===```

第三级  ``` * / % ```

第四级 ```| &``` 

第五级  ``` + - ^ ```

---

**当前支持的内置函数有如下，更多函数将会在逐步补充 (当前仍然存在bug）**
通用处理函数
```golang
print
varDump
```
---
字符串处理函数
函数的命名基本参考了 Go 语言
除了个别函数有差别，如 
utf8Len 用于检测字符串字数的函数
isNumeric 用于判断当前输入是否为数字

```
trim
trimLeft
trimRight
trimSpace
len
utf8Len
md5
replace
contains
index
lastIndex
toLower
toUpper
toTitle
repeat
```
---
其他函数
```
isNumeric
isBool
isInt
isFloat
```

---


#### 计划后期支持的特性

##### 数组、对象或结构体 (从中选择其一实现)、递归、基础的动态库 (如 tcp/http/mysql等)、异常处理

由于它将是一个解释型语言，出于效率考虑，我不会让它支持任何的语法糖，不支持多余的、锦上添花的特性，避免在运行时过多影响效率。

得益于 Golang 原生支持的一些高级特性，例如 管道、协程 这些可轻松实现异步操作，Lit 会有更多扩展空间，它会比 php 更灵活、更小巧。

敬请期待...
