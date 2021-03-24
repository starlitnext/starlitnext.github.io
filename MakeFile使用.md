
## 参考文档
* [GNU make](https://www.gnu.org/software/make/manual/html_node/index.html)
* [makefile tutorial](https://makefiletutorial.com)
* [通过实例学Makefile](https://zhuanlan.zhihu.com/p/317716664)

## make & makefile
* make 是一个工具，可以自动识别哪些程序修改了需要重新编译
* makefile 是一个文件，告诉make如何编译和链接一个程序

## makefile基础语法
``` makefile
targets: prerequisites
    command
    command
    command
```
* 一个makefile由上述的规则（rules）组成
* rules:
    * *targets*: 文件名，一个ruleyige
    * *commands*: 一些列的步骤，用来生成target
    * *prerequisites*: 一系列用空格分隔的文件名，command运行需要的文件列表

## Beginner Examples
``` makefile
blah: blah.o
    cc blah.o -o blah # Runs third

blah.o: blah.c
    cc -c blah.c -o blah.o # Runs second

blah.c:
    echo "int main() { return 0; }" > blah.c # Runs first
```
* 运行 make blah 时，会生成一个blah可执行程序，步骤如下：
    * 先搜索blah作为target
    * blah需要blah.o，搜索生成blah.o的规则
    * blah.o 需要blah.c，搜索生成blah.c的规则
    * blah.c 没有依赖，运行echo command
    * 运行 cc -c 生成 blah.o
    * 运行 cc 生成 blah
* 默认的target是第一个target（写在最前面的那个）

## 变量
* 可以把重复的文件赋值给一个变量
* 使用 ${} 或者 $() 引用变量
* 注意 **=** 和 **:=** 的区别
    * **=** 递归展开变量 （使用的时候才展开）
    * **:=** 简单展开变量（立马展开）
``` makefile
files = file1 file2
some_file: $(files)
    echo "Look at this variable: " $(files)
    touch some_file

file1:
    touch file1
file2:
    touch file2

clean:
    rm -f file1 file2 some_file
```

## all target
* 使用 all target 来生成多个target
``` makefile
all: one two three

one:
    touch one
two:
    touch two
three:
    touch three

clean:
    rm -f one two three
```

## .PHONY
``` makefile
objects := main.o kbd.o command.o display.o \
          insert.o search.o files.o utils.o

edit : $(objects)
    cc -o edit $(objects)

.PHONY : clean
clean :
    rm edit $(objects)
```
* 一个"虚假"的目标，即它并不是一个文件名，但它可以作为一个目标，当我们输入make clean时就可以执行它。如果你的目录内真有一个clean的文件，make判断clean文件存在，而它又没有依赖其他的目标，下面的rm...就不会执行到。为了告诉make不要把它当成一个文件判断，这里使用了.PHONY把clean指定为虚假目标，这样make就一定会执行这个则。

## Automatic Variables
* [文档](https://www.gnu.org/software/make/manual/html_node/Automatic-Variables.html)
* $@    target的文件名
* $<    The name of the first prerequisite
* $?    The names of all the prerequisites that are newer than the target, with spaces between them
* $^    The names of all the prerequisites, with spaces between them 
``` makefile
hey: one two
    # Outputs "hey", since this is the first target
    echo $@

    # Outputs all prerequisites newer than the target
    echo $?

    # Outputs all prerequisites
    echo $^

    touch hey

one:
    touch one

two:
    touch two

clean:
    rm -f hey one two
```

## Wildcard
* \* 和 % 符号
* \* 
    * 从文件系统中匹配文件名
    * 可以在 target、prerequisite 和 wildcard 函数中使用
    * 变量定义中不能直接使用，可以用 wildcard
    * 如果匹配不到文件名，则会保持写的样子，除非在wildcard中使用
``` makefile
thing_wrong := *.o # Don't do this! '*' will not get expanded
thing_right := $(wildcard *.o)

all: one two three four

# Fails, because $(thing_wrong) is the string "*.o"
one: $(thing_wrong)

# Stays as *.o if there are no files that match this pattern :(
two: *.o 

# Works as you would expect! In this case, it does nothing.
three: $(thing_right)

# Same as rule three
four: $(wildcard *.o)
```
* %
    * 替换引用，把.c文件替换成.o，$(var:a=b)
``` makefile
foo := a.c b.c c.c
bar := $(foo:%.c=dir/%.c)
default:
    @echo $(bar)
# 输出：dir/a.c dir/b.c dir/c.c
```
```makefile
SRCS := $(wildcard *.c) # 使用内置函数自动匹配源文件
OBJS := $(SRCS:.c=.o)   # 将源文件转成目标文件

PROG := edit
CC := gcc -std=c99
CFLAGS := -Wall -O2
LDFLAGS :=
LIBS :=

$(PROG) : $(OBJS)
    $(CC) -o $@ $^ $(LDFLAGS) $(LIBS)

.PHONY : clean
clean :
    rm $(PROG) $(OBJS)
```

## 生成头文件依赖
* 使用 gcc 的 -MM 来生成目标文件依赖的源文件和头文件, 把这些依赖关系拷贝到Make文件里就行，比较麻烦
* 自动生成依赖规则，使用GNU make 的 remake 特性可以做到
``` makefile
SRCS := $(wildcard *.c)
OBJS := $(SRCS:.c=.o)
DEPS := $(SRCS:.c=.d)

PROG := edit
CC := gcc -std=c99
CFLAGS := -Wall -O2
LDFLAGS :=
LIBS := -lm

$(PROG) : $(OBJS)
    $(CC) -o $@ $^ $(LDFLAGS) $(LIBS)

#############################################
# 自动生成依赖文件
ifneq ($(MAKECMDGOALS),clean)
-include $(DEPS)
endif

%.d: %.c
    @echo "make depend: $@"
    @set -e; rm -f $@; \
    $(CC) $(CFLAGS) -MM $< | sed -E 's,($*).o[: ]*,\1.o $@: ,g' > $@
#############################################

.PHONY : clean
clean :
    rm -f $(PROG) $(OBJS) $(DEPS)
```
