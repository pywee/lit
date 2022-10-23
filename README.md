# Lit 轻量级解释型弱类型语言 (迭代中...)


goExpr Lit 是一个轻量级解释型语言的锥形，目前仍在持续开发中。它基于 Golang 开发，还没有<b>正式的名字</b>，暂且叫作 <b>Lit</b> 好了。原本我只打算用于实现对文本中包含的算术表达式进行计算，后来逐步地扩充了更多的特性。索性将它定位为一个解释型编程语言好了。

目前我已经实现了一些编程语言必须具备的基础功能，尤其是比较重要的算术表达式、弱类型转换、内置函数、变量声明等。仍然需要进一步完善。<u>因此你还不能直接将其用于日常开发中。</u>

我对 Lit (goExpr) 的定位是一个解释型弱类型语言，希望它能用于日常的 Web 开发，这是我的目标。至于能否顺利完成，还需要很长时间才能下定论。

关于目前已实现的特性，请看如下文档，我将定期更新：

当前文档更新日期是 2022.10.19

---

#### 使用方法


```
go get github.com/pywee/goExpr
```



**一、支持变量声明**
```golang
import "github.com/pywee/goExpr"

func main() {
    // 执行以下句子，最终会输出
    src := []byte(`
        a = 123;
        b = a + 456;
        Print(b); // 579
    `)
    _, err := goExpr.NewExpr(src)
}

```

---

**二、算术表达式的计算。算术符号的优先级保持与 Go 语言相同。请看示例：**

```golang
    // 不同的语言符号优先级是不完全一样的，Lit 的算术符号优先级保持与 Golang 一致
    // 首先我们执行 Go 语言原生函数进行数学表达式计算
    // 以下句子最终会输出 
    // +1.600000e+001
    // 125
    println((2 + 100 ^ 2 - (10*1.1 - 22 + (22 | 11))) / 10 * 2)
    println(12/333+31+(5/10)-6|100)

    // 使用 Lit 计算文本中的数据
    // 表达式文本
    // 执行下面的句子 最终会输出
    exprs := []byte(`
        a = (2 + 100 ^ 2 - (10*1.1 - 22 + (22 | 11))) / 10 * 2;
        b = 12 / 333 + 31 + (5 / 10) - 6 | 100;
        Print(a); // 16
        Print(b); // 125
    `)
    _, err = goExpr.NewExpr(exprs)

    // *** 同样的表达式放在 PHP 中，会输出 -24 ***

```

---

**三、当前已支持部分常用内置函数(测试阶段)，更多的内置函数我将在接下来继续完成**

```golang

    // 下面的句子调用了两个函数 
    // IsInt(arg) 用来检查 arg 是否为整型 
    // Replace(arg1, arg2, arg3, arg4) 用来做字符串替换

    // 执行下面语句 最终会输出
     exprs := []byte(`
        a = Replace("hello word111", "1", "", 2-IsInt((1+(1 + IsInt(123+(1+2)))-1)+2)-2);
        VarDump(a); // STRING hello word
    `)
    _, err = goExpr.NewExpr(exprs)

```
---
**四、弱类型转换，弱类型的这一特性我将它设计为与 PHP 基本一样**
```golang
    // 当布尔值参与运算时，底层会将 true 转为 1, false 转为 0
    // 执行下面句子 将输出
    src := []byte(`
        a = true - 1;
        b = IsInt(1);
        c = IsFloat(1);
        d = false == 0.0;
        e = "false" == 0.0;
        Print(a); // 0
        Print(b); // true
        Print(c); // true
        Print(d); // true
        Print(e); // true
    `)
    _, err := goExpr.NewExpr(src)


    // 与其他弱类型语言一样
    // 字符串数字与整型相操作，在 Lit 的底层会将字符串数字转换为整型
    // 执行下面句子 将输出
    src := []byte(`
        a = "1" - 1;
        b =  0.0 >= false+1 || (1<=21 && 1==1);
        Print(a); // 0
        Print(b); // true
    `)
    _, err := goExpr.NewExpr(src)


    // 与其他弱类型语言一样
    // 字符串数字与整型相操作，在 Lit 的底层会将字符串数字转换为整型
    // 执行下面句子 将输出
    src := []byte(`
        a = "1" - 1;
        Print(a); // 0
    `)
    _, err := goExpr.NewExpr(src)


    // 字符串与字符串相加时 将进行字符串的拼接
    // 执行以下句子，将会输出
    src := []byte(`
        a = "abc" + "def";
        Print(a); // abcdef
    `)
    _, err := goExpr.NewExpr(src)


    // 但如果当两个字符串都为数字时 对他们进行相加 则会被底层转换为数字
    // 执行如下句子，将会输出
     src := []byte(`
        a = "123" + "456";
        Print(a); // 579
    `)
    _, err := goExpr.NewExpr(src)


    // 其他字符串+整型将会报错
    // 执行如下句子
    src := []byte(`
    	a = "abcwwww1230"+0.01;
    	Print(a); // 报错
    `)
    _, err := goExpr.NewExpr(src)
```


---


**五、"并且" 与 "或者" 符号处理**
```golang
    // 执行如下句子，将会输出
    src := []byte(`
        a = IsInt(1) && 72+(11-2) || 1-false;
        VarDump(a); // BOOL true
    `)
    _, err := goExpr.NewExpr(src)

```


---
**请注意，Lit 的算术符号优先级向 Golang 看齐。每个语言对算术符号的优先级处理都有一定区别，如，针对以下表达式进行计算时：**

``` golang
// 2 + 100 ^ 2 - (10*1.1 - 22 + (22 | 11)) / 10 * 2
// PHP 输出 -104
// Node.js 输出 -104
// Golang 输出 96
// Lit (goExpr) 输出 96
```

---

**Lit 算术符号优先级**
第一级  ``` () && || ```

第一级  ``` > < >= <= == != ===```

第三级  ``` * / % ```

第四级 ```| &``` 

第五级  ``` + - ^ ```

---

**当前支持的内置函数有如下，更多函数将会在逐步补充**
```golang
    // 通用处理函数
    Print
    VarDump

    // 字符串处理函数
    // 函数的命名基本参考了 Go 语言
    // 除了个别函数有差别，如 
    // UTF8Len 用于检测字符串字数的函数
    // IsNumeric 用于判断当前输入是否为数字
    Trim
    TrimLeft
    TrimRight
    TrimSpace
    Len
    UTF8Len
    MD5
    Replace
    Contains
    Index
    LastIndex
    ToLower
    ToUpper
    ToTitle
    Repeat

    // 其他函数
    IsNumeric
    IsBool
    IsInt
    IsFloat
```


##### 后期将会实现更多特性，敬请期待...

