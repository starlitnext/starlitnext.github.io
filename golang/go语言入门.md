<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=3 orderedList=false} -->

<!-- code_chunk_output -->

- [环境安装](#环境安装)
- [语言结构](#语言结构)
- [语法基础](#语法基础)
- [数据类型](#数据类型)
- [变量](#变量)
- [常量](#常量)
- [运算符](#运算符)
- [条件语句](#条件语句)
- [循环语句](#循环语句)
- [函数](#函数)
- [数组](#数组)
- [指针](#指针)
- [结构体](#结构体)
- [切片](#切片)
- [Map](#map)
- [接口（interface）](#接口interface)
- [并发](#并发)

<!-- /code_chunk_output -->
## 环境安装
* 安装包下载地址 [golang](https://golang.org/dl/)

#### windows 下安装
* 下载msi安装文件
* 将bin目录添加到path

#### linux 下安装
* 下载linux版本安装文件（xxx.tar.gz）
* 解压至`/usr/local`: `tar -C /usr/local -xzf xxx.tar.gz`
* 将`usr/local/go/bin` 添加到PATH : `export PATH=$PATH:/usr/local/go/bin`

#### helloworld 测试
* 注意不要命名为xxx_test.go，golang默认`*_test.go`为测试文件，会报`go run: cannot run *_test.go files`
``` go
// 01_demo.go
package main
import "fmt"
func main() {
  fmt.Println("Hello, World!")
}
// >>> go run test.go
// Hello, World!
```

## 语言结构
* `package` 定义包名，`package main` 表示一个可独立执行的程序，每个Go应用程序都包含一个名为`main`的包
* `import` 引入包，可以用 `import ()`同时引入多个包
* func 定义函数，`main`函数时每个可执行程序必须包含的
* 注释 `/*...*/` 或 `//`
* `fmt.Println` 和 `fmt.Print`
* 当标识符以**大写**字母开头，就可以被外部包代码使用，称为**导出**，小写字母开头则对包外是不可见的
* `go run *.go` 运行
* `go build *.go` 生产二进制文件
* `{` 不能单独放在一行，`func main() {` 这里 `{`不能放到下一行

## 语法基础
* 行分隔符：不需要`;`，如果多个语句写在同一行，则需要`;`分隔，不推荐
* 标识符：和 `c语言` 一样
* 字符串拼接 `+`，如 `fmt.Println("Hello," + "World!")`
* 关键字：25个关键字或保留字，36个预定义标识符
* 变量声明 `var identifier type`
* 格式化字符串：`fmt.Sprintf(fmtstr, var1, var2)`

## 数据类型
* bool：`var b bool = true`
* 基于架构的类型 
  * `uint8 uint16 uint32 uint64 int8 int16 int32 int64`
  * `float32 float64 complex64 complex128`
* 其它数字类型
  * `byte`: 类似`uint8`
  * `rune`: 类似`int32`
  * `uint`: 32位或64位
  * `int`: 与`uint`大小一样
  * `uintptr`: 无符号整型，用于存放一个指针
* 类型转换 `type_name(expression)`

## 变量
* `var v_name v_type` `v_name = value` 声明和赋值分开，系统默认值
* `var v_name = value` 自行判定变量类型
* `v_name := value` 和上面一样，不过直能在函数体中出现，不能放在全局作用域下
* 多变量声明 
  * `var vname1, vname2, vname3 type`
  * 因式分解关键字写法，一般用于声明全局变量
    ``` go
    var (
      vname1 vtype1
      vname2 vtype2
    )
    ```
* 取内存地址 `&i`
* 多变量赋值
  ``` go
  var a, b int
  var c string
  a, b, c = 5, 6, "abc"
  ```
* 交换两个变量的值：`a, b = b, a`，两个变量的类型必须相同
* `_`用于抛弃值，如 `_, b = 5, 7`

## 常量
* 程序运行时，不会被修改的量
* `const identifier [type] = value`
* 用作枚举
  ``` go
  const (
    Unknow = 0
    Female = 1
    Male = 2
  )
  ```
* 常量可以用`len, cap, unsafe.Sizeof` 函数计算表达式的值，常量表达式中函数必须是内置函数
``` go
package main
import (
  "unsafe"
  "fmt"
)
const (
  a = "abc"
  b = len(a)
  c = unsafe.Sizeof(a)
)
func main() {
  fmt.Println(a, b, c) // abc 3 16
}
```
#### iota
* 一个特殊的常量，一个可以被编译器修改的常量
* `iota` 在`const`关键字出现时被重置位0，`const`中每新增一行`常量声明`，`iota`计数一次
* `iota`可以理解为`const`语句块中的`行索引`
``` go
package main
import "fmt"
func main() {
  const (
    a = itoa    // 0
    b           // 1
    c           // 2
    d = "ha"    // 独立值，iota += 1
    e           // "ha" iota += 1
    f = 100     // 100 iota += 1
    g           // 100 iota += 1
    h = itoa    // 7 恢复计数
    i           // 8
  )
  fmt.Println(a,b,c,d,e,f,g,h,i) // 0 1 2 ha ha 100 100 7 8
}
```
``` go
package main
import "fmt"
const {
  i = 1<<iota
  j = 3<<iota
  k 
  l
}
func main() {
  fmt.Println("i=", i)  // i= 1
  fmt.Println("j=", j)  // j= 6
  fmt.Println("k=", k)  // k= 12
  fmt.Println("l=", l)  // l= 24
}
```

## 运算符
* 和 `C语言` 基本完全一致

## 条件语句
* `if` `if ... else if ... else` `switch` `select`
* `select` 随机执行一个可运行的case，如果没有case可运行，它将阻塞

## 循环语句
* `for init; condition; post {}`
* `for condition {}` 和 `C语言`的 `while`一样
* `for {}` 和 `C语言`的`for(;;)`一样
* `for ke, value := range oldMap {}`
``` go
package main
import "fmt"
func main() {
  strings := []string{"google", "runoob"}
  for i, s := range strings {
    fmt.Println(i, s)
  }
  numbers := [6]int{1, 2, 3, 5};
  for i, x := range numbers {
    fmt.Println("第%d位x的值=%d\n", i, x)
  }
}
```

## 函数
``` go
func function_name([parameter list]) [return_types] {
  // function body
}
```
* 闭包
``` go
package main
import "fmt"
func getSequence() func() int {
  i := 0
  return func() int{
    i += 1
    return i
  }
}
func main() {
  nextNumber := getSequence()
  fmt.Println(nextNumber()) // 1
  fmt.Println(nextNumber()) // 2
  fmt.Println(nextNumber()) // 3

  nextNumber1 := getSequence()
  fmt.Println(nextNumber1()) // 1
  fmt.Println(nextNumber1()) // 2
}
```
* 方法
``` go
func (variable_name variable_data_type) function_name() [return_type] {
  // method body
}
```

## 数组
* 一维数组 `var variable_name [SIZE] variable_type`
* 多维数组 `var variable_name [SIZE1][SIZE2]...[SIZEN] variable_type`
* 初始化数组
  * `var balance = [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}`
  * `balance := [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}`
  * 上面 `[5]` 如果写成 `[...]` 编译器会自动根据**元素个数**推断数组长度
  * 根据**下标**来初始化元素 `balance := [5]float32{1:2.0, 3:7.0}`
  * 下标从 `0` 开始

## 指针
`var var_name *var_type`
* 空指针 `nil`
* 指向指针的指针 `var ptr **int`

## 结构体
``` go
type struct_variable_type struct {
  member definition
  member definition
  ...
  member definition
}
variable_name := struct_variable_type {value1, value2 ... valuen}
variable_name := struct_variable_type {key1: value1, key2: value2 ... keyn: valuen}
```
* 结构体方法 见 `函数` 章节

## 切片
* 对数组的抽象，实现动态数组
* 定义一个位指定大小的数组来定义切片 `var identifier []type`
* 或使用 `make` 函数来创建切片 `var slice1 []type = make([]type, len, capacity)`, `capacity`可选

## Map
``` go
var map_variable map[key_data_type]value_data_type
map_variable := make(map[key_data_type]value_data_type)
```
* 可以用 range 遍历 `for key := range myMap {}`
* 查看元素是否存在 `value, ok := myMap[key]`
* `delete` 函数删除集合的元素

## 接口（interface）
``` go
type interface_name interface {
  method_name1 [return_type]
  method_name2 [return_type]
  method_name3 [return_type]
  ...
  method_namen [return_type]
}
type struct_name struct {
  // variables
}
func (struct_name_variable struct_name) method_name1() [return_type] {
}
func (struct_name_variable struct_name) method_namen() [return_type] {
}
```
* 实例，实现多态
``` go
package main
import "fmt"
type Phone interface {
  call()
}
type NokiaPhone struct {
}
func (nokiaPhane NokiaPhone) call() {
  fmt.Println("I am Nokia, I can call you!")
}
type IPhone struct {
}
func (iPhone IPhone) call() {
  fmt.Println("I am iPhone, I can call you!")
}
func main() {
  var phone Phone
  phone = new(NokiaPhone)
  phone.call()

  phone = new(IPhone)
  phone.call()
}
```

## 并发
* 语言支持并发，通过 `go` 关键字来开启 `goroutine`
* `goroutine` 是轻量级线程，由 `Golang runtime` 进行调度管理
* `go func_name([params])` 语法格式
* 同一个程序中的所有 `goroutine` 共享同一个地址空间
* 实例
``` go
package main
import (
  "fmt"
  "time"
)
func say(s string) {
  for i := 0; i < 5; i++ {
    time.Sleep(100 * time.Millisecond)
    fmt.Println(s)
  }
}
func main() {
  go say("world")
  say("hello")
}
```

#### 通道
* 用于两个 `goroutine` 之间传递一个指定类型的值类同步运行和通讯
* `<-` 用于指定通道的方向，发送或接收。如未指定方向，则为**双向通道**
```
ch := make(chan int)  // 创建通道，用 `chan` 关键字
ch <- v               // 把v发送到通道 ch
v := <- ch            // 从ch接收数据，并赋值给v
```
* 默认情况下，通道不带缓冲区，发送端发送数据，同时必须有接收端接收相应的数据
* 实例，通过两个`goroutine` 来计算数字之和
``` go
package main
import "fmt"
func sum(s []int, c chan int) {
  sum := 0
  for _, v := range s {
    sum += v
  }
  c <- sum    // 把sum 发送给通道 c
}
func main() {
  s := []int{7, 2, 8, -9, 4, 0}
  c := make(chan int)
  go sum(s[:len(s)/2], c)
  go sum(s[len(s)/2:], c)
  x, y := <-c, <-c  // 从通道c接收
  fmt.Println(x, y, x+y) // -5, 17, 12
}
```
* 如果通道缓冲区满了。则发送方会阻塞，直到写入成功
* 设置缓冲区大小 `ch := make(chan int, 100)`
* 发送端使用close关闭通道，接收端使用`range`读取可以停止，否则会在发送端不发送以后阻塞接收端
``` go
package main
import "fmt"
func fibonacci(n int, c chan int) {
  x, y := 0, 1
  for i := 0; i < n; i++ {
    c <- x
    x, y = y, x+y
  }
  close(c)
}
func main() {
  c := make(chan int, 10)
  go fibonacci(cap(c), c)
  // 如果上面不close，则range不会结束，导致阻塞
  for i := range c {
    fmt.Println(i)
  }
}
```

@import "echarts.min.js"
``` javascript {cmd=true element="<div id='showechart' style='width: 500px; height: 300px; margin: 0 auto'></div>" hide}
var myChart = echarts.init(document.getElementById('showechart'));
// 指定图表的配置项和数据
var option = {
    xAxis: {
        type: 'category',
        data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
    },
    yAxis: {
        type: 'value'
    },
    series: [{
        data: [150, 230, 224, 218, 135, 147, 260],
        type: 'line'
    }]
};
// 使用刚指定的配置项和数据显示图表。
myChart.setOption(option);
```
