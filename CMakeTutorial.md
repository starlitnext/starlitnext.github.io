---
html:
    toc:true:
---
<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=5 orderedList=false} -->

<!-- code_chunk_output -->

- [1. 开始](#1-开始)
    - [添加可执行文件](#添加可执行文件)
    - [构建和运行](#构建和运行)
    - [添加版本号和可配置头文件](#添加版本号和可配置头文件)
    - [设置C++ 标准](#设置c-标准)
    - [Rebuild](#rebuild)
- [2. 添加 Library](#2-添加-library)
    - [添加一个子文件夹](#添加一个子文件夹)
    - [添加一个选项来决定是否使用我们自己的库](#添加一个选项来决定是否使用我们自己的库)
- [3.添加使用我这个库的依赖](#3添加使用我这个库的依赖)
- [4.安装](#4安装)
    - [安装规则](#安装规则)
    - [安装根路径](#安装根路径)
- [5.添加系统检查](#5添加系统检查)

<!-- /code_chunk_output -->

## 1. 开始
#### 添加可执行文件
最基础的project是从源代码构建一个可执行文件
``` cmake
cmake_minimum_required(VERSION 3.10)
# set the project name
project(Tutorial)
# add then executable
add_executable(Tutorial tutorial.cxx)
```
#### 构建和运行
``` shell
mkdir build
cd build
# 配置project并生成一个native build system（如make）
cmake ..
# 然构建系统真正的compile/link工程
cmake --build .
```
最后会生成名为`Tutorial`的可执行文件

#### 添加版本号和可配置头文件
给项目添加一个版本号，可以直接把版本号写在代码中，但使用cmake可以提供更大的灵活性
* 修改CMakeLists.txt 在`project()`中设置工程名和版本号
``` cmake
# set the project name and version
project(Tutorial VERSION 1.0)
```
* 从header file传递版本号给源代码
``` cmake
configure_file(TutorialConfig.h.in TutorialConfig.h)
```
* TutorialConfig.h 会生成到build目录下，所以我们需要把目录添加到头文件搜索目录
``` cmake
target_include_directories(Tutorial PUBLIC
                           "${PROJECT_BINARY_DIR}"
                           )
```
* 在源码目录中创建TutorialConfig.h
``` cpp
// the configured options and settings for Tutorial
#define Tutorial_VERSION_MAJOR @Tutorial_VERSION_MAJOR@
#define Tutorial_VERSION_MINOR @Tutorial_VERSION_MINOR@
```
当CMake在处理这个头文件的时候，`@Tutorial_VERSION_MAJOR@` 和 `@Tutorial_VERSION_MINOR@`将会替换成正确的值
* 修改 tutorial.cxx 并包含 TutorialConfig.h
``` cpp
if (argc < 2) {
    std::cout << argv[0] << " Version " << Tutorial_VERSION_MAJOR << "."
              << Tutorial_VERSION_MINOR << std::endl;
    std::cout << "Usage: " << argv[0] << " number" << std::endl;
    return 1;
}
```

#### 设置C++ 标准
* 在tutorial.cxx中使用c++11的feature
``` cpp
const double inputValue = std::stdod(argv[1]);
```
* 在 CMakeLists.txt 中指定使用c++11标准
``` cmake
# specify the c++ standard
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED True)
```

#### Rebuild
``` cmake
cd build
cmake --build .
```

## 2. 添加 Library
#### 添加一个子文件夹
* 我们添加子文件夹MathFunctions，来生成我们的库，子文件夹包含MathFunctions.h和mysqrt.cxx
* 添加 MathFunctions/CMakeLists.txt
``` cmake
add_library(MathFunctions mysqrt.cxx)
```
* 在顶层的CMakeLists.txt中添加
``` cmake
# add the MathFunctions library
add_subdirectory(MathFunctions)

# add the executable
add_executable(Tutorial tutorial.cxx)

target_link_libraries(Tutorial PUBLIC MathFunctions)

# add the binary tree to the search path for include files
# so that we will find TutorialConfig.h
target_include_directories(Tutorial PUBLIC
    "${PROJECT_BINARY_DIR}"
    "${PROJECT_SOURCE_DIR}/MathFunctions"
)
```
#### 添加一个选项来决定是否使用我们自己的库
* 顶层CMakeLists.txt 添加一个 `option`，默认开启
``` cmake
option(USE_MYMATH "Use tutorial provided math implementation" ON)
```
* 使用 `if` 来决定是否链接我们的library
``` cmake
if (USE_MYMATH)
    add_subdirectory(MathFunctions)
    list(APPEND EXTRA_LIBS MathFunctions)
    list(APPEND EXTRA_INCLUDES "${PROJECT_SOURCE_DIR}/MathFunctions")
endif()

add_executable(Tutorial tutorial.cxx)

target_link_libraries(Tutorial PUBLIC ${EXTRA_LIBS})

target_include_directories(Tutorial PUBLIC
    "${PROJECT_BINARY_DIR}",
    ${EXTRA_INCLUDES}
)
```
* tutorial.cxx 中也需要相应的处理
``` cpp
#ifdef USE_MYMATH
#   include "MathFunctions.h"
#endif

#ifdef USE_MYMATH
    const double outputValue = mysqrt(inputValue);
#else
    const double outputValue = sqrt(inputValue);
#endif
```
* 由于代码中需要使用 `USE_MYMATH`，我们可以添加到 `TutorialConfig.h.in` 中
``` cpp
#cmakedefine USE_MYMATH
```
* 最后可以在生成build的时候修改选项
``` shell
cmake .. -DUSE_MYMATH=OFF
```

## 3.添加使用我这个库的依赖
* 上面我们使用一个库的时候，需要添加`EXTRA_INCLUDES`头文件搜索目录。实际上可以在MathFunctions中就声明，使用我这个库的地方需要包含我的头文件目录
* 使用`INTERFACE`表示使用者需要，我自己不需要，在 `MathFunctions/CMakeLists.txt`中添加以下内容
``` cmake
target_include_directories(MathFunctions
    INTERFACE ${CMAKE_CURRENT_SOURCE_DIR}
)
```
* 接下来就可以删除 `EXTRA_INCLUDES`了
* 这里好处再于，由提供者来声明使用我的人需要包含哪些头文件路径，不容易出问题

## 4.安装
#### 安装规则
* MathFunctions安装需要头文件和library文件，在MathFunctions/CMakeLists.txt中添加
``` cmake
install(TARGETS MathFunctions DESTINATION lib)
install(FILES MathFunctions.h DESTINATION include)
```
* 顶层的CMakeLists.txt中添加
``` cmake
install(TARGETS Tutorial DESTINATION bin)
install(FILES "${PROJECT_BINARY_DIR}/TutorialConfig.h"
    DESTINATION include
)
```
#### 安装根路径
* `CMAKE_INSTALL_PREFIX` 决定安装的根路径
* 使用 `cmake --instal . --prefix "/home/myuser/installdir"` 可以配置安装根路径
* 在vscode 中 settings.json 中添加
``` json
"cmake.configureArgs":["-DCMAKE_INSTALL_PREFIX=${workspaceFolder}/install"]
```
然后重新生成缓存

## 5.添加系统检查
* 使用`CheckSymbolExists`来检查函数是否存在
* 检查log和exp函数的例子，如果检查不到，则链接m库来查找，MathFunctions/CMakeLists.txt中添加
``` cmake
target_include_directories(MathFunctions
    INTERFACE ${CMAKE_CURRENT_SOURCE_DIR}
)

# does this system provide the log and exp functions?
include(CheckSymbolExists)
check_symbol_exists(log "math.h" HAV_LOG)
check_symbol_exists(exp "math.h" HAV_EXP)
if(NOT (HAV_LOG AND HAV_EXP))
    unset(HAV_LOG CACHE)
    unset(HAV_EXP CACHE)
    set(CMAKE_REQUIRED_LIBRARIES "m")
    check_symbol_exists(log "math.h" HAV_LOG)
    check_symbol_exists(exp "math.h" HAV_EXP)
    if (HAV_LOG AND HAV_EXP)
        target_link_libraries(MathFunctions PRIVATE m)
    endif()
endif()

# if available，specify `HAV_LOG` and `HAV_EXP` as `PRIVATE` compile definitions
if (HAV_LOG AND HAV_EXP)
    target_compile_definitions(MathFunctions
                               PRIVATE "HAV_LOG" "HAV_EXP")
endif()
```
* mysqrt.cxx 中使用
``` cpp
#if defined(HAV_LOG) && defined(HAV_EXP)
    double result = exp(log(x) * 0.5);
#else
    double result = x;
#endif
```
