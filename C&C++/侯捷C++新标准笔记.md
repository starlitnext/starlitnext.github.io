# 视频课程地址
* [b站](https://www.bilibili.com/video/BV1Rp4y1n7fb?from=search&seid=13773620021362443107)

# 演进、环境与资源
* C++标准演进
  - C++98 (1.0)
  - C++03 (TR1, Technical Report 1)
  - C++11 (2.0)
  - C++14 对11进行了一些小修复和改进
* C++2.0 新特性包括`语言`和`标准库`两个层面，后者以header files形式呈现
  - C++标准库的header files不带后缀（.h），如 `#include<vector>`
  - 新式 C header files 不带.h，如 `#include<cstdio>`
  - 旧式 C header files 任可用，如 `#include<stdio.h>`
* C++2.0 新增的header files
``` c++
#include <type_traits>
#include <unordered_set>
#include <forward_list>
#include <array>
#include <tuple>
#include <reges>
#include <regex>
#include <thread>
using namespace std;
```
* TR1 之前`namespace std::tr1`的内容都放到`namespace std`了

# c++2.0 语言部分

## Variadic Templates
* 数量不定的模板参数，接收任意个参数，每个参数为任意类型
``` c++
void print()
{
}
template<typename T, typename... Types>
void print(const T& firstArg, const Types&... args)
{
  cout << firstArg << endl; // print first argument
  print(args...);           // call print() for remaining arguments
}

print(7.5, "hello", bitset<16>(377), 42);
```
* `sizeof...(args)` 为参数的个数
* 特例化共存,考虑与上面代码共存时执行哪个？可以共存，会优先调用上面这个，因为`比较特化`，下面这个永远不会被调用
``` c++
template<typename... Types>
void print(const Types&... args)
{
  cout << sizeof...(args) << endl;
}
```
* 方便的完成 recursive function call，如实现hash function
``` c++
class CustomerHash {
public:
  std::size_t operator()(const CustomerHash& c) {
    return hash_val(c.fname, c.lno, c.no); // call 1
  }
};

// 1
template <typename... Types>
inline size_t hash_val(const Types&... args) {
  size_t seed = 0;
  hash_val(seed, args...); // call 2
  return seed;
}

// 2
template<typename T, typename... Types>
inline void hash_val(size_t& seed, const T&val, const Types&... args) {
  hash_combine(seed, val); // call 4
  hash_val(seed, args...); // recursive call 2，final call 3
}

// 3
template<typename T>
inline void hash_val(size_t& seed, const T& val) {
  hash_combine(seed, val); // call 4
}

// call 4
template<typename T>
inline void hash_combine(size_t& seed, const T& val){
  seed ^= std::hash<T>()(val) + 0x9e3779b9 + (seed << 6) + (seed << 2);
}
```
* 做 recursive inheritance
``` c++
template<typename... Values>class tuple;
template<> class tuple<> {};
template<typename Head, typename... Tail>
class tuple<Head, Tail...>
  : private tuple<Tail...>
{
  typedef tuple<Tail...> inherited;
public:
  tuple() {}
  tuple(Head v, Tail... vtail)
    : m_head(v), inherited(vtail...) {}
  
  typename Head::type head() { return m_head; }
  inherited& tail() { return *this; }
  
protected:
  Head m_head;
};
```

## 一些小特性
* `vector<list<int> >`可以不用这个空格了，写成 `vector<list<int>>`
* 增加`nullptr`，nullptr是关键字，nullptr 的类型为 `std::nulptr_t`
``` c++
// stddef.h
#if defined(__cplusplus) && __cplusplus >= 201103L
#ifndef _GXX_NULLPTR_T
#define _GXX_NULLPTR_T
  typedef decltype(nullptr) nullptr_t;
#endif
#endif /* C++11.  */

void f(int);
void f(void*);
f(0);       // calls f(int)
f(NULL);    // calls f(int) if NULL is 0, ambiguous otherwise
f(nullptr); // calls f(void*)
```
* 增加`auto` 关键字，自动类型推导

## Uniform Initialization
* 统一使用大括号初始化
``` c++
int values[] {1, 2, 3};
vector<int> v {2, 3, 5, 7};
vector<string> cities {
  "Berlin", "New York", "London"
};
complex<double> c{4.0, 3.0};
```
* 编译器看到{t1, t2, ...} 便做出一个 initializer_list\<T>, 它关联至一个array\<T, n>, 调用函数（如 ctor）时，该array内的元素可被编译器分解逐一传给函数
* 但若函数参数是个initializer_list\<T>，调用者却不能给予数个T参数，然后自动转为一个initializer_t\<T>传入

## Initializer_list
* 初始化
``` c++
int i;    // i has undefined value
int j{};  // j is initialized by 0
int* p;   // p has undefined value
int* q{}; // q is initialized by nullptr
```
* narrowing initializations
``` c++
int x1(5.3);        // OK, but OUCH: x1 becomes 5
int x2 = 5.3;       // OK, but OUCH: x2 becomes 5
int x3{5.0};        // ERROR: narrowing, GCC [warning]
int x4 = {5.3};     // ERROR: narrowing
char c1{7};         // OK, even though 7 is an int, this is not narrowing
char c2{99999};     // ERROR: narrowing (if 9999 doesn't fit into a char)
std::vector<int> v1 {1, 2, 3, 4, 5};  // OK
std::vector<int> v2 {1, 2.3, 4, 5.6}; // ERROR: narrowing
```
* 做可变函数参数
``` c++
void print(std::initializer_list<int> vals)
{
}
print({1,2,3,4,5});
```
* 作为构造函数
  - 如果没有构造函数2，则q和s调用1，r报错
``` c++
class P
{
public:
  // 1
  P(int a, int b)
  {
    cout << "P(int, int), a=" << a << ", b=" << b << endl;
  }

  // 2
  P(initializer_list<int> initlist)
  {
    cout << "P(initializer_list<int>), value= ";
    for (auto i : initlist)
      cout << i << ' ';
    cout << endl;
  }
};

P p(77, 5);         // P(int, int), a=77, b=5
P q{77, 5};         // P(initializer_list<int>), value= 77 5
P r{77, 5, 42};         // P(initializer_list<int>), value= 77 5 42
P s={77, 5};        // P(initializer_list<int>), value= 77 5
```

## Explicit for ctors taking more than one augument
* c++1.0 explicit 使用场景
``` c++
struct Complex
{
  int real, imag;
  Complex(int re, int im=0) : real(re), imag(im)
  { }
  Complex operator+(const Complex& x)
  { return Complex(real + x.real, imag + x.imag); }
};
Complex c1(12, 5);
Complex c2 = c1 + 5;  // 5 隐式转换

struct Complex
{
  int real, imag;
  explicit
  Complex(int re, int im=0) : real(re), imag(im)
  { }
  Complex operator+(const Complex& x)
  { return Complex(real + x.real, imag + x.imag); }
};
Complex c1(12, 5);
Complex c2 = c1 + 5;  
// [ERROR] no match for 'operator+'(operand types are 'Complex' and 'int')
```
* c++2.0 任意参数个数的ctors都可以指定explicit

## range-based `for` statement
``` c++
for ( decl : coll ) {
  statement
}
```

## =default, =delete
* 如果你自己定义了一个ctor，那么编译器就不会再给你一个default ctor
* 如果你强制加上=default，就可以重新获得并使用default ctor
``` c++
class Zoo
{
public:
  Zoo(int i1, int i2) : d1(i1), d2(i2) { }
  Zoo(const Zoo&) = delete;
  Zoo(Zoo&&) = default;
  Zoo& operator=(const Zoo&) = default;
  Zoo& operator=(const Zoo&&) = delete;
  virtual ~Zoo() { }
private:
  int d1, d2;
};
```
* 什么情况下需要写Big-Three
  - 带有指针类型的成员
* default 的拷贝，按bit拷贝对象
* 2.0 Big-Five
* No-Copy and Private-Copy
``` c++
struct NoCopy {
  NoCopy() = default;
  NoCopy(const NoCopy&) = delete;
  NoCopy &operator=(const NoCopy&) = delete;
  ~NoCopy() = default;
};
class PrivateCopy {
private:
  PrivateCopy(const PrivateCopy&);
  PrivateCopy &operator=(const Private&);
public:
  PrivateCopy() = default;
  ~PrivateCopy() = default;
};
```

## Alias Template (template typedef)
``` c++
template <typename T>
using Vec = std::vector<T, MyAlloc<T>>;
Vec<int> coll;
// is equivalent to
std::vector<int, MyAlloc<int>> coll;
```
* 使用 macro 和 typedef 实现不了上面这种用法，typedef接收参数
* It is `not possible` to partially or explicity specialize an alias template

## Template template parameter
``` c++
template<typename T,
         template <class>
         class Container
        >
class XCls
{
private:
  Container<T> c;
public:
  XCls() {
    for (long i=0; i < SIZE; ++i)
      c.insert(c.end(), T());
    output_static_data(T());
    Container<T> c1(c);
    Contianer<T> c2(std::move(c));
    c1.swap(c2);
  }
};

template<typename T>
using Vec = vector<T, allocator<T>>;
template<typename T>
using Lst = list<T, allocator<T>>;

XCls<MyString, Vec> c1;
```

## Type Alias, noexcept, override, final
* Type Alias, similar to typedef
``` c++
// typedef void (*func)(int, int);
using func = void(*)(int, int);
void example(int, int) {}
func fn = example;
```
* using 的应用场景
``` c++
// using-directives for namespaces and using-declarations for namespace members
using namespace std;
using std::cout;
// using-declarations for class members
protected:
  using _Base::MyAlloc;
// type alias and alias template
using func = void(*)(int, int);
template<typename T>
struct Container {
  using value_type = T;
};
template<class CharT> using mystring = std::basic_string<CharT, std::char_traits<CharT>>;
```
* noexcept
* override
* final

## decltype
``` c++
map<string, float> coll;
decltype(coll)::value_type elem;

template<typename T1, typename T2>
decltype(x+y) add(T1 x, T2 y);

template<typename T1, typename T2>
auto add(T1 x, T2 y) -> decltype(x+y);

auto cmp = [](const Person& p1, const Person& p2) {
  return p1.lastname() < p2.lastname()
};
std::set<Person, decltype(cmp)> coll(cmp);
```
* defines a type equivalent to `the type of an expression`
* 使用场景
  * 声明return type
  * metaprogramming
  * pass the type of a lambda

## Lambdas
``` c++
[...](...)mutable_opt throwSpec_opt -> retType_opt {...}
```
* `inline` functionality
* can be used as a `parameter` or a `local object`
* 如果没有参数，可以不写 `()`
* `[]`可以传入外部变量
  - [=] 按值传递
  - [&] 按引用传递
  - [=, &y] y按引用传递，其余按值传递
  - 没有加特殊符号时，按值传递
* `mutable` 按值传递时，可以修改这个值，但不影响外部的值
``` c++
int id = 0;
auto f = [id]() mutable {
  cout << "id: " << id << endl;
  ++id; // 如果没有mutable：[Error] increment of read-obly variable 'id'
}; // 可以理解为一个class，里面有一个int id 的成员，这个id由外部传入
id = 43;
f();
f();
f();
cout << id << endl;
// id: 0
// id: 1
// id: 2
// 42

int id = 0;
auto f = [&id](int param) {
  cout << "id: " << id << endl;
  ++id; ++param; // OK
};
id = 42;
f(7);
f(7);
f(7);
cout << id << endl;
// id: 42
// id: 43
// id: 44
// 45
```
* lambda 类型没有默认构造函数和赋值运算符
``` c++
auto cmp = [](const Person& p1, const Person& p2) {
  return p1.lastname() < p2.lastname();
}
std::set<Person, decltype(cmp)> coll(cmp);

template<class Key,
         class Compare = less<Key>
         class Alloc = alloc>
class set {
...
public:
  set() : t(Compare()) {} // 如果std::set<Person, decltype(cmp)> coll;调用的话，那么这里需要调用delctype(cmp) 类型的默认构造函数，导致编译报错
  explicit set(const Compare& comp) : t(comp) {}
};
```
* 和使用仿函数的对比
``` c++
vector<int> vi { 5, 28, 50, 83, 70, 590, 245, 59, 24 };
int x = 30;
int y = 100;
vi.erase( remove_if(vi.begin(),
                    vi.end(),
                    [x, y](int n) { return x < n && n < y;}
                    ),
          vi.end()
        ); // remove_if 把满足条件的元素移到容器后面，并返回第一个满足条件的位置
for (auto i : vi)
  cout << i << " "; // 5, 28, 590, 245, 24
cout << endl;

class LambdaFunctor {
public:
  LambdaFunctor(int a, int b) : m_a(a), m_b(b) { }
  bool operator()(int n) const { return m_a < n && n < m_b; }
private:
  int m_a;
  int m_b;
};
vi.erase( remove_if(v.begin(), v.end(),
                    LambdaFunctor(x, y)),
          v.end()
        );
```

## Variadic Templates
* 变化的是template parameters
  - 参数个数 variable number
  - 参数类型 different type
* 递归处理模板参数
``` c++
void func() { /*...*/ }
template<typename T, typename... Types>
void func(const T& firstArg, const Types&... args)
{
  // 处理 first Arg ...
  func(args...);
}
```
* 重写printf
``` c++
void printf(const char* s) {
  while (*s) {
    if (*s == '%' && *(++s) != '%')
      throw std::runtime_error("invalid format string: missing arguments");
    std::cout << *s++;
  }
}
template<typename T, typename... Args>
void printf(const char* s, T value, Args... args) {
  while (*s) {
    if (*s == '%' && *(++s) != '%') {
      std::cout << value;
      printf(++s, args...);
      return;
    }
    std::cout << *s++;
  }
  throw std::logic_error("extra arguments provided to printf");
}

int* pi = new int;
printf(
  "%d %s %p %f\n",
  15,
  "This is Ace.",
  pi,
  3.141592653
);
// 15 This is Ace.0x3e4ab8 3.14159
```
* 实现max
``` c++
// 使用initializer_list
max({ 56, 48, 60, 100, 20, 18 }); // 需要加上 {}

// 使用 variadic templates
int maximum(int n)
{
  return n;
}
template<typename... Arg>
int maximum(int n, Args... args)
{
  return std::max(n, maximum(args...)); // 最后一次调用max(a, b)这样的形式
}
maxmimum(56, 48, 60, 100, 20, 18); // 不需要加上 {}
```
* 打印 [7.5,hello,0000011111,42]，即头尾的处理方式不同，需要加上[]
``` c++
cout << make_tupe(7.5, string("hello"), bitset<16>(377), 42);

template<typename... Args>
ostream& operator<<(ostream& os, const tuple<Args...>& t) {
  os << "[";
  PRINT_TUPLE<0, sizeof...(Args), Args...>::print(os, t);
  return << "]";
}
template<int IDX, int MAX, typename... Args>
struct PRINT_TUPLE {
  static void print(ostream& os, const tuple<Args...>& t) {
    os << get<IDX>(t) << (IDX+1 == MAX ? "" : ","); // 非最后一个元素打印 `,`
    PRINT_TUPLE<IDX+1, MAX, Args...>::print(os, t); // 最终调用到 PRINT_TUPLE<MAX, MAX, Args...>
  }
};
// partial specialization to end the recursion
template<int MAX, typename... Args>
struct PRINT_TUPLE<MAX, MAX, Args...> {
  static void print(ostream& os, const tuple<Args...>& t) { }
};
```
* 实现递归继承，recursive inheritance，`tuple的实现`
* 实现复合，recursive composition
``` c++
template<typename... Values> class tup;
template<> class tup<> { };
template<typename Head, typename... Tail>
class tup<Head, Tail...> {
  typedef tup<Tail...> composited;
protected:
  composited m_tail;
  Head m_head;
public:
  tup() { }
  tup(Head v, Tail... vtail) : m_tail(vtail...), m_head(v) { }
  Head head() { return m_head; }
  composited& tail() { return m_tail; }
}
```

# 标准库部分

## 右值引用 Rvalue reference
* Lvalue: 可以出现在 operator= 左侧
* Rvalue: 只能出现在 operator= 右侧
* Unperfect Forwarding
``` c++
void process(int& i) { cout << "process(int&):" << i << endl; }
void process(int&& i) { cout << "process(int&&):" <<i << endl; }
void forward(int&& i) {
  cout << "forward(int&&):" << i << ",";
  process(i);
}

int a = 0;
process(a);                     // process(int&):0  变量被视为lvalue处理
process(1);                     // process(int&&):1 temp object 被视为rvalue处理
process(move(a));               // process(int&&):0 强制将a由lvalue改为rvalue
forward(2);                     // forward(int&&):2,process(int&):2 rvalue 经由 forward() 传给另一个函数却变成了lvalue
forward(move(a));               // forward(int&&):0,process(int&):0
forward(a);                     // [Error] cannot bind 'int' lvalue to 'int&&'
const int& b = 1;
process(b);                     // [Error] no matching function for call to 'process(const int&)'
process(move(b));               // [Error] no matching function for call to 'process(std::remove_reference<const int&>::type)'
int& x(5);                      // [Error] invalid initialization of non-const reference of type 'int&' from an rvalue of type 'int
```
* Perfect forwarding
``` c++
template<typename T1, template T2>
void functionA(T1&& t1, T2&& t2) {
  functionB(std::forward<T1>(t1),
            std::forward<T2>(t2));
}
```
* 写一个move aware class
``` c++
class MyString {
public:
  static size_t DCtor;    // default ctor 调用次数
  static size_t Ctor;     // ctor 调用次数
  static size_t CCtor;    // copy ctor 调用次数
  static size_t CAsgn;    // copy assignment 调用次数
  static size_t MCtor;    // move ctor 调用次数
  static size_t MAsgn;    // move assignment 调用次数
  static size_t Dtor;     // dtor 调用次数
private:
  char* _data;
  size_t _len;
  void _init_data(const char* s) {
    _data = new char[_len+1];
    memcpy(_data, s, _len);
    _data[len] = '\0';
  }
public:
  MyString() : _data(NULL), _len(0) { ++DCtor; }
  MyString(const char* p) : _len(strlen(p)) { ++Ctor; _init_data(p); }
  MyString(const MyString& str) : _len(str.len) { ++CCtor; _init_data(str._data); }
  MyString(MyString&& str) noexcept 
    : _data(str._data), _len(str._len) {
    ++MCtor;
    str._len = 0;
    str._data = NULL; // 重要
  }
  MyString& operator=(const MyString& str) {
    ++CAsgn;
    if (this != &str) {
      if (_data) delete _data;
      _len = str._len;
      _init_data(str._data);
    } else {
    }
    return *this;
  }
  MyString& operator=(MyString&& str) noexcept {
    ++MAsgn;
    if (this != &str) {
      if (_data) delete _data;
      _len = str._len;
      _data = str._data;
      str._len = 0;
      str._data = NULL;
    }
    return *this;
  }
  virtual ~MyString() {
    ++Dtor;
    if (_data) delete _data;
  }
};
```