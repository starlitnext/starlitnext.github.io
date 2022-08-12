## lua reload 方案

### 概述
* 目前市面上用的lua reload方案基本是适用云风提供的那个实现方案，[源码](https://github.com/cloudwu/luareload)，[如何让 lua 做尽量正确的热更新](https://blog.codingnow.com/2016/11/lua_update.html)
* 其基本原理是用一个沙盒把新的模块require进来，然后对面新旧两个模块进行对比，把有修改的内容引用成新的，而在替换函数的实现上，他的方案是把lua的VM整个进行遍历，把引用老函数的地方替换成新的
* 遍历VM的做法性能很差，在开发期可以使用，但是无法在线上使用，因此需要对替换函数的实现进行优化，比如在python中可以直接替换func_code 来实现，不需要去遍历VM替换函数对象，只是把字节码对象改了就ok了，那么在lua中没有func_code这样的结构提供给脚本，因此需要对lua的函数实现源码进行分析。
* 因为无法知道两个函数是否不一致，因此遍历VM的时候等于模块的所有函数都要进行替换，这个消耗是巨大的

### lua的函数实现分析
* lua中存在两种函数对象，用C实现的函数 `CClosure` 和 脚本中实现的函数 `LClosure`，这个结构体中都有一个 Struct Proto 的指针，实际上就类似于python中的funcode，只不过脚本中拿不到这个对象的引用
```
typedef struct LClosure {
  ClosureHeader;
  struct Proto *p;
  UpVal *upvals[1];  /* list of upvalues */
} LClosure;

typedef struct Proto {
  CommonHeader;
  lu_byte numparams;  /* number of fixed (named) parameters */
  lu_byte is_vararg;
  lu_byte maxstacksize;  /* number of registers needed by this function */
  int sizeupvalues;  /* size of 'upvalues' */
  int sizek;  /* size of 'k' */
  int sizecode;
  int sizelineinfo;
  int sizep;  /* size of 'p' */
  int sizelocvars;
  int sizeabslineinfo;  /* size of 'abslineinfo' */
  int linedefined;  /* debug information  */
  int lastlinedefined;  /* debug information  */
  TValue *k;  /* constants used by the function */
  Instruction *code;  /* opcodes */
  struct Proto **p;  /* functions defined inside the function */
  Upvaldesc *upvalues;  /* upvalue information */
  ls_byte *lineinfo;  /* information about source lines (debug information) */
  AbsLineInfo *abslineinfo;  /* idem */
  LocVar *locvars;  /* information about local variables (debug information) */
  TString  *source;  /* used for debug information */
  GCObject *gclist;
} Proto;
```
* 在第一次require的时候，lua底层会创建对应的Proto对象，里面包含了运行时的一些信息，以及upvalue的一些信息，lua中的函数对象就这一个实现
* lua中的upvalue是指函数中使用的上层定义域中的变量，在LClosure中才保存真正的引用对象，Proto中保存的是Upvaluedesc，包含了upvalue的名字等信息

### lua函数替换原理
* 线上的bug修复，一般函数修改会比较小，因此大部分情况下，upvalue是不会变的（函数都没有修改），这种情况下可以直接替换Proto对象，并处理一下upvalue的顺序就可以完成，这种方式速度非常快
* 如果upvalue有变化的话，则只能替换LClosure对象了，但是我也不希望通过遍历VM的方式来找出所有引用
    - 在lua的VM实现中，对象都会用GCObject放到一个全局的gclist中，用于实现lua的gc，那么可以利用这个数据结构来找出某个LClosure对象的引用
    - gclist中遍历时需要注意下，table只需要遍历第一层就行了，如果table中如果引用了另一个table，那么这个被引用的table一定也会存在gclist中，后面也可以遍历到

### lua reload 实现
* 同样和云风的方案一样，用一个沙盒require新的模块，两者对比，进行深度遍历对比，找出需要替换的函数，并判断一下upvalue是否改变。
    - 基础的值数据类型直接替换
    - table的话继续深度遍历
* 按照上述lua函数替换原理对函数进行替换
* 对于有些模块，可能不希望做替换，外面引用也不管了，可以提供一个标记，目前我的实现是在模块开头内定义一个local __reload_all = true，这种模块的reload只需要从package.loaded 中移除，再require一下就行了
* 做一个lua的c模块实现以下基础接口
    - replace_proto
    - replace_function
    - is_reload_all