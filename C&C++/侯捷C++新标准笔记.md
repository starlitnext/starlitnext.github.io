## 视频课程地址
* [b站](https://www.bilibili.com/video/BV1Rp4y1n7fb?from=search&seid=13773620021362443107)

## 演进、环境与资源
* C++标准演进
  - C++98 (1.0)
  - C++03 (TR1, Technical Report 1)
  - C++11 (2.0)
  - C++14 对11进行了一些小修复和改进
* C++2.0 新特性包括`语音`和`标准库`两个层面，后者以header files形式呈现
  - C++标准库的header files不带后缀（.h），如 `#include\<vector>`
  - 新式 C header files 不带.h，如 `#include\<cstdio>`
  - 旧式 C header files 任可用，如 `#include\<stdio.h>`
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
* 特例化共存,考虑与上面代码共存时执行哪个？
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
* 什么情况下需要写Bit-Three
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
