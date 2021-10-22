---
html:
    toc:true:
---
<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [前言](#前言)
    - [A First Example](#a-first-example)
    - [The Stack](#the-stack)
- [Extending Your Application](#extending-your-application)
    - [基础实例](#基础实例)
    - [table 操作](#table-操作)

<!-- /code_chunk_output -->

## 前言
* `embeddable language` 和 `extensible language`
    * 嵌入到application中来扩展功能
    * 使用C语言给lua提供额外功能
* 两种交互模式：
    * C 作为application code，lua作为library
    * lua has the control，C作为library
* C API: C code 和 Lua 交互的 functions，constants and types 的集合
    * read and write lua global variables
    * call lua functions
    * run pieces of lua code
    * register c functions that can be called by lua code
* lua 自带的一些代码可作为参考：
    * appliction code : lua.c
    * standard libriries: lmathlib.c、 lstrlib.c

#### A First Example
* A bare-bones stand-alone Lua interpreter
``` lua
#include <stdio.h>
#include <string.h>
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>

int main(int argc, char* argv[]) {
    char buff[256];
    int error;
    lua_State *L = luaL_newstate();
    luaL_openlibs(L);

    while (fgets(buff, sizeof(buff), stdin) != NULL) {
        error = luaL_loadstring(L, buff) || lua_pcall(L, 0, 0, 0);
        if (error) {
            fprintf(stderr, "%s\n", lua_tostring(L, -1));
            lua_pop(L, 1); // pop error message from the stack
        }
    }
    
    lua_close(L);
    return 0;
}
```
* `luaL_loadstring` : 如果成功，返回0，并把结果压入栈顶
* `lua_pcall` : 从栈顶pop出函数然后在保护模式下执行，如果没报错则返回0
* 上述两个函数如果出错，则会在栈顶压入错误信息

#### The Stack
* Lua和C之间交互通过一个`virtual stack` 来实现，几乎所有的C API都会对这个stack进行操作
* stack中的每个slot都可以放任意的lua数据类型
* 每一个 `lua type`都有一个对应的 `push function`来把一个lua的值push到这个stack中
``` cpp
void lua_pushnil(lua_State *L);
void lua_pushboolean(lua_State *L, int bool);
...
```
所有这些方法都可以在`lua.h`中找到声明
* 使用 `luaL_checkstack` 或 `lua_ckeckstack` 方法来确定stack是否有足够的空间
* API 使用 `indicies`来引用stack中的元素
    - push的第一个元素索引为1，第二个为2，也就是栈底元素为1，然后往上累加
    - 栈顶元素的索引为-1，下一个为-2，然后往下递减
``` ditaa {cmd=true args=["-E"] hide=true}
栈顶
/--------\
| c0A0   |  -1   3 <-lua_gettop()
+--------+
| c0A0   |  -2   2
+--------+
| c0A0   |  -3   1
\--------/
栈底
```
* `lua_is*` 判断栈中某个元素是否为lua的某个类型
* `lua_type` 返回栈中某个元素的类型：`LUA_TNIL,LUA_TBOOLEAN,...`
* `lua_to*` 把栈中某个元素转化为某种C的类型
* 一个例子 Dumping the stack
``` c

```
* 对stack进行操作的api
``` lua
int lua_gettop(lua_State *L);
void lua_settop(lua_State *L, int index);
void lua_pushvalue(lua_State *L, int index);
void lua_rotate(lua_State *L, int index, int n);
void lua_remove(lua_State *L, int index);
void lua_insert(lua_State *L, int index);
void lua_replace(lua_State *L, int index);
void lua_copy(lua_State *L, int fromidx, int toidx);
```

## Extending Your Application

#### 基础实例
* 这个例子中，我们将把lua作为应用的配置脚本使用
* 假设我们有个应用程序，需要从配置文件中读取窗口的大小，lua的配置如下
``` lua
-- define window size
width = 200
height = 300
```
* 在应用程序中读取到配置的数据
    - lua_loadfile 把文件中的代码块加载到stack中，栈顶为`compiled chunk`
    - lua_pcall 运行 `compiled chunk`
    - lua_getglobal 从global中获取value，push到栈顶
``` c
int getglobint(lua_State *L, const char *var) {
    int isnum, result;
    lua_getglobal(L, var); // push value to stack
    result = (int)lua_tointegerx(L, -1, &isnum);
    if (!isnum)
        error(L, "'%s' should be a number\n", var);
    lua_pop(L, 1);  // remove result from stack
    return result;
}

void load(lua_State *L, const char *fname, int *w, int *h) {
    if (luaL_loadfile(L, fname) || lua_pcall(L, 0, 0, 0))
        error(L, "cannot run config. file: %s", lua_tostring(L, -1));
    *w = getglobint(L, "width");
    *h = getglobint(L, "height");
}
```

#### table 操作
* 假设我们需要在lua配置文件中获取背景颜色，那么在lua中配置如下
``` lua
BLUE = {red = 0, green = 0, blue = 1.0}
-- background = {red = 0.30, green = 0.10, blue = 0}
background = BLUE
```
* 在c中读取lua table 的数据
    - lua_istable 判断栈中某个元素是否为table
    - lua_gettable pop栈顶元素作为key，把对应index处的table的值push到栈顶
``` c
// assume that table is on the top of the stack
int getcolorfield(lua_State *L, const char *key) {
    int result, isnum;
    lua_pushstring(L, key); // push key
    lua_gettable(L, -2); // pop key, push background[key]
    result = (int)(lua_tonumberx(L, -1, &isnum) * MAX_COLOR);
    if (!isnum)
        error(L, "invalid component '%s' in color", key);
    lua_pop(L, 1); // remove number
    return result;
}

...

lua_getglobal(L, "background");
    if (!lua_istable(L, -1))
        error(L, "'background' is not a table");
    int red, green, blue;
    red = getcolorfield(L, "red");
    green = getcolorfield(L, "green");
    blue = getcolorfield(L, "blue");
    printf("background<rgb>: <%d, %d, %d>\n", red, green, blue);
```