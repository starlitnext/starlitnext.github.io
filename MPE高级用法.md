* [Markdown Preview Enhanced 官方文档](https://shd101wyy.github.io/markdown-preview-enhanced/#/)
* [在 VSCode 下用 Markdown Preview Enhanced 愉快地写文档](https://zhuanlan.zhihu.com/p/56699805)

### 语法增强
* 添加目录，在文件开始添加下面代码
```
<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=3 orderedList=false} -->
```
* 转化成html以后目录放到侧边栏，需要在文件顶部加上：
``` 
---
html:
    toc:true
---
```
* vscode 添加代码片段，`ctrl+shift+p` 打开settings.json添加
```
"[markdown]":  {
    "editor.quickSuggestions": true
  }
```
`文件`->`首选项`->`用户片段`打开markdown.json添加
```
"Add TOC": {
		"prefix": "toc",
		"body": [
			"---",
			"html:",
    		"    toc:true:",
			"---",
			"<!-- @import \"[TOC]\" {cmd=\"toc\" depthFrom=2 depthTo=3 orderedList=false} -->"
		],
		"description": "add toc"
	}
```
* `==`来==高亮==
* Emoji: L^A^T~E~X :smile:
* L^A^T~E~X 公式： `$$\int_{-\infty}^{\infty} e^{-x^2} = \sqrt{\pi}$$`
$$\int_{-\infty}^{\infty} e^{-x^2} = \sqrt{\pi}$$
* 画图：[ditaa](https://github.com/stathissideris/ditaa)
``` ditaa {cmd=true args=["-E"] hide=true}
+----------------------\
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
                                                   MAX
                                                    v
/------------+--------------------+-------+-------+-\
| c700       |   c990             | c0A0  | cA0A  | |
\------------+--------------------+-------+-------+-/
```
* 在文档里跑代码
  - 需要打开 `Enable Script Execution`
  - `Ctrl + Shift + p` 运行 `Run Code Trunk` or `Run All Code Trunk`
``` lua {cmd=true hide=true}
print("test lua")
```
* 代码输出可以是`HTML`或者`Markdown`，无缝成为文档的一部分，实现动态文档
  - 输出html的例子
``` lua {cmd=true hide=true output="html"}
local lv = 10
print("<h2 align=center>等级为"..lv.."</h2>")
```
  - 输出markdown的例子
``` lua {cmd=true hide=true output="markdown"}
print("**标题**")
print("|等级|攻击力|伤害|")
print("|--|--|--|")
for i = 1, 5 do
  print("|"..i.."|"..(i*10).."|"..(i*20).."|")
end
```
* 使用`@import`
@import "https://cloud.githubusercontent.com/assets/1908863/22716507/f352a4b6-ed5b-11e6-9bac-88837f111de0.gif" {width="500px" title="my title" alt="my alt"}
  - 插入图片
    `@import "apue/环境表.png"`
  - 插入很长的代码
    `@import "test.lua" {cmd=lua hide}`
  - 插入其它`markdown`
    `@import "sec1.md"`
  - 插入表格
    `@import "test.csv"`