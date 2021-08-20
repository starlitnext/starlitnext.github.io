
* [在 VSCode 下用 Markdown Preview Enhanced 愉快地写文档](https://zhuanlan.zhihu.com/p/56699805)

### 语法增强
* 添加目录，在文件开始添加下面代码
```
<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=3 orderedList=false} -->
```
* `==`来==高亮==
* Emoji: L^A^T~E~X :smile:
* L^A^T~E~X 公式： `$$\int_{-\infty}^{\infty} e^{-x^2} = \sqrt{\pi}$$`
$$\int_{-\infty}^{\infty} e^{-x^2} = \sqrt{\pi}$$
* 图表：ditaa
``` ditaa {cmd=true args="-E" hide=true}
+----------------------+
|        c1FF          |
|        SKY           |
+----------------------+
|                      |
|                      |
|        cGRE          |
|                      |
|        PLAY          |
|       GROUND         |
|                      |
|                      |
|                      |
|                      |
+----------------------+
```
* 在文档里跑代码
  - 需要打开 `Enable Script Execution`
  - `Ctrl + Shift + p` 运行 `Run Code Trunk`
``` lua {cmd=true hide=true}
print("test lua")
```
