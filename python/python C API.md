### 头文件
* Python.h里面包含了所有需要的CPython中的内容，并且包含stdio.h, string.h, errno.h, limits.h, assert.h, stdlib.h(如果Platform有这些头文件的话)
* CPython中的所有Object或者函数，宏等都是以Py开头，\_Py开头的一般是CPython内部使用的
### Objects, Types and Reference Counts
* Python中一切解释Object，有一个抽象的类型PyObject
* steal reference 的概念, PyTuple\_SetItem, PyList\_SetItem
``` c++
PyObject *t;
t = PyTuple_New(3);
PyTuple_SetItem(t, 0, PyInt_FromLong(1L));
PyTuple_SetItem(t, 1, PyInt_FromLong(2L));
PyTuple_SetItem(t, 2, PyString_FromString("three"));
```
* 相比于PyXXX_New 方法，一般使用Py\_BuildValue方法来快速创建
```C++
PyObject *tuple, *list;
tuple = Py_BuildValue("(iis)", 1, 2, "three");
list = Py_BuildValue("[iis]", 1, 2, "three");
```
* 注意PyXXX\_GetItem和PySequence\_GetItem的区别，使用PyXXX\_GetItem返回一个Object实际上并不会对这个Object增加引用，但是PySequence\_GetItem则会，所以使用完之后需要主动调用Py_DECREF来释放
```C++
long
sum_list(PyObject *list)
{
    int i, n;
    long total = 0;
    PyObject *item;

    n = PyList_Size(list);
    if (n < 0)
        return -1; /* Not a list */
    for (i = 0; i < n; i++) {
        item = PyList_GetItem(list, i); /* Can't fail */
        if (!PyInt_Check(item)) continue; /* Skip non-integers */
        total += PyInt_AsLong(item);
    }
    return total;
}
long
sum_sequence(PyObject *sequence)
{
    int i, n;
    long total = 0;
    PyObject *item;
    n = PySequence_Length(sequence);
    if (n < 0)
        return -1; /* Has no length */
    for (i = 0; i < n; i++) {
        item = PySequence_GetItem(sequence, i);
        if (item == NULL)
            return -1; /* Not a sequence, or other failure */
        if (PyInt_Check(item))
            total += PyInt_AsLong(item);
        Py_DECREF(item); /* Discard reference ownership */
    }
    return total;
}
```
### Exceptions
* 一个例子, Python代码对应的C代码如下
```Python
def incr_item(dict, key):
    try:
        item = dict[key]
    except KeyError:
        item = 0
    dict[key] = item + 1
```
```C++
int
incr_item(PyObject *dict, PyObject *key)
{
    /* Objects all initialized to NULL for Py_XDECREF */
    PyObject *item = NULL, *const_one = NULL, *incremented_item = NULL;
    int rv = -1; /* Return value initialized to -1 (failure) */

    item = PyObject_GetItem(dict, key);
    if (item == NULL) {
        /* Handle KeyError only: */
        if (!PyErr_ExceptionMatches(PyExc_KeyError))
            goto error;

        /* Clear the error and use zero: */
        PyErr_Clear();
        item = PyInt_FromLong(0L);
        if (item == NULL)
            goto error;
    }
    const_one = PyInt_FromLong(1L);
    if (const_one == NULL)
        goto error;

    incremented_item = PyNumber_Add(item, const_one);
    if (incremented_item == NULL)
        goto error;

    if (PyObject_SetItem(dict, key, incremented_item) < 0)
        goto error;
    rv = 0; /* Success */
    /* Continue with cleanup code */

 error:
    /* Cleanup code, shared by success and failure path */

    /* Use Py_XDECREF() to ignore NULL references */
    Py_XDECREF(item);
    Py_XDECREF(const_one);
    Py_XDECREF(incremented_item);

    return rv; /* -1 for error, 0 for success */
}
```
### Embedding Python
* Py\_Initialize 做初始化，实际上做了这几件事：初始化加载的模块、创建\_\_builtin\_\_, \_\_main\_\_, sys 和exceptions，初始化模块的module的搜索路径（sys.path）
* PySys\_SetArgvEx(argc, argv, updatepath),用来设置sys.argv
* 自己写的embedding代码可能会面临找不到python模块的问题，可以使用PySys\_SetPath来设置搜索路径，或者使用PyRun\_SimpleString("import sys");PyRun\_SimpleString("sys.path.append('.')");
* 一个Embedding 的例子
```Python
# -*- coding:utf-8 -*-
def init():
    print 'hello, world!'
```
```C++
#include <Python.h>
#include <stdio.h>

int main()
{
    printf("before\n");
    Py_Initialize();
    PyRun_SimpleString("import sys");
    PyRun_SimpleString("sys.path.append('..')");
    PyObject *pModule = NULL;
    PyObject *pFunc = NULL;
    pModule = PyImport_ImportModule("hello");
    if (pModule == NULL)
    {
        return -1;
    }
    pFunc = PyObject_GetAttrString(pModule, "init");
    if (pFunc == NULL)
    {
        return -1;
    }
    PyEval_CallObject(pFunc, NULL);
    Py_Finalize();
    return 0;
}
```
### 给Python添加C模块
* 一个简单的例子
```Python
import spam
status = spam.system("ls -l")
```
```C++
static PyObject* spam_system(PyObject *self, PyObject *args)
{
    const char *command;
    int sts;
    if (!PyArg_ParseTuple(args, "s", &command))
        return NULL;
    sts = system(command);
    return Py_BuildValue("i",sts);
}
static PyMethodDef SpamMethods[] = {
    ...
    {"system",  spam_system, METH_VARARGS,
     "Execute a shell command."},
    ...
    {NULL, NULL, 0, NULL}        /* Sentinel */
};
int
main(int argc, char *argv[])
{
    /* Pass argv[0] to the Python interpreter */
    Py_SetProgramName(argv[0]);

    /* Initialize the Python interpreter.  Required. */
    Py_Initialize();

    /* Add a static module */
    Py_InitModule("spam", SpamMethods);
    ...

```
* 在这里C的函数默认需要两个参数self和args，如果是模块方法，self为NULL或者Py_InitModule4传过来的，如果是成员函数，则这里的self就是一个对象
* args是一个Python Tuple的指针，包含了参数，Tuple中的每一个参数都是一个Python Object
* 使用PyArg_ParseTuple把Python的Tuple转化为C中的类型
### 脚本引擎绑定
* 无非就是写了一套wrapper，脚本调用Application的方法的时候，首先调用到暴露给脚本的wrapper方法，wrapper在转发调用Application对应的接口
* 像Cocos的脚本绑定，实际上是使用工具自动导出这一套wrapper
* 如果我需要给Cocos添加一个Python脚本引擎的话，思路大概也是这样，创建一个cocos工程，写一个cc模块，用来暴露给Python层使用
* 在脚本层想要通过cc.Director这样的方式来获取Director对象，也就是说需要给模块增加属性
* 上面已经有给模块添加方法的实例了，那么还需要给模块添加新的Object类型
* 使用 PyModule_Addobject方法来添加PythonTypeObject










