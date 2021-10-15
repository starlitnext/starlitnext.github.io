
### 基础
* 安装 cmake 和 cmake tools 插件
* 在c_cpp_properties.json 中添加头文件搜索路径
``` json
"configurations": [
        {
            "includePath": [
                "${workspaceFolder}/**",
            ],
            "configurationProvider": "vector-of-bool.cmake-tools"
        }
    ],
```

### 添加cmake的参数
* 如想要实现 cmake .. -DUSE_MYMATH=OFF
* 在settings.json 中添加
``` json
"cmake.configureArgs":["-DUSE_MYMATH=OFF"]
```
* 修改参数以后需要重新生产一下缓存 `ctrl+shift+p` 运行 `CMake: Delete Cache and Reconfig`
* `CMakeLists.txt` 每次`ctrl+s`会重新生产缓存，但修改了参数以后必须使用上述方式重新生成

### 调试
* 生产CMake缓存以后，左侧CMake项目大纲里面每个构建目标会有一个对应的项目
* 选中某个目标，点击鼠标右键，点击`调试`