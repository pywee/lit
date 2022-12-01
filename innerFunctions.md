### 内置函数支持


```
len 返回字符串长度或数组长度 (支持半角和全角长度判断)
len(str, [rune bool]) int 

trim 返回去掉两边指定字符后的字符串 默认去掉空格
trim(str, [charlist=' ' string...]) string

trimLeft 返回去掉左边指定字符后的字符串 默认去掉空格
trimLeft(str, [charlist='', string...]) string

trimRight 返回去掉右边指定字符后的字符串 默认去掉空格
trimRight(str, [charlist='', string...]) string

trimSpace 返回将 str 两端所有空白都去掉之后的字符串 包含空格 \n \r\n \t
trimSpace(str string) string

rand 返回指定范围内随机数
rand(min int, max int) int

split 根据 sep 分割字符串，返回数组
split(str string, sep string) array

查询 sep 是否为 arr 的值
inArray(arr array, sep string|int) bool

查询 key 是否在数组中
keyInArray(arr array, key string|int) bool

返回数组中所有的 key
arrayKeys(arr array) array

返回数组中所有的 value (不保留原有的下标名称，所有下标将重新以整型递增返回)
arrayValues(arr array) array

删除数组指定的下标 返回删除了下标后的数组 
如果没有找到要删除的下标 将数组原样返回
delete(arr array, index int|string) array

base64Encode(str string) string
base64Decode(str string) string

urlEncode(str string) string
urlDecode(str string) string
```


