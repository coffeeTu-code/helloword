# Markdown 编写规范 - 基于有道云
[>_>注释]:<# 一级标题>
[>_>注释]:<## 二级标题>
[>_>注释]:<### 三级标题>
[>_>注释]:<#### 四级标题>
[>_>注释]:<##### 五级标题>
[>_>注释]:<###### 六级标题>


[TOC]  

[>_>注释]:<在段落中填写 [TOC] 以显示全文内容的目录结构。>  

## 文字设置  
### 字体  
正常  *斜体*  _斜体_  **加粗**  ***斜体加粗***  ~~删除线~~  

### 特殊符号
\\  反斜线  \*  星号  \_  下划线  \{ ... \} \[ ... \] \( ... \)  括号  \#  井号  \+  加号  \-  减号  \.  句号  \!  感叹号  
[特殊字符Unicode码 查找链接](https://unicode-table.com/cn/)  
[emoji表情查找链接](https://www.webfx.com/tools/emoji-cheat-sheet/)  

### 缩进
你好  
&#160;你好&nbsp;你好，缩进1/4中文  
&#8194;你好&ensp;你好，缩进1/2中文  
&#8195;你好&emsp;你好，缩进1个中文（建议选择）  

### 引用
> 第一段文字[^1]  
> 
>> 二级引用  
>   
> 第二段文字[^2]  

[^1]:注脚1，显示在页面的末端  
[^2]:注脚2，显示在页面的末端  

### 字体颜色
详细的字体颜色设置需要查看附录。  
<font face="黑体">我是黑体字</font>  
<font face="微软雅黑">我是微软雅黑</font>  
<font face="STCAIYUN">我是华文彩云</font>  
<font color=#0099ff size=3 face="黑体">我是黑体字</font>  
<font color=#00fffff size=5>我是黑体字</font>  
<font color=gray size=7>我是黑体字</font>  
<table><tr><td bgcolor=orange>设置背景色</td></tr></table>  

## 数据设置
### 插入列表
在使用列表时，只要是数字后面加上英文的点，就会无意间产生列表，比如2017.12.30。解决方式：在每个点前面加上'\\'就可以了。  
无序列表  
* 第一种使用方式  
+ 第二种使用方式  
- 第三种使用方式  

有序列表  
1. 一  
2. 二  
3. 三  

### 插入表格
表格对齐方式：我们可以指定表格单元格的对齐方式，冒号在左边表示左对齐，右边表示有对齐，两边都有表示居中。
学号    |姓名   |分数  
:-:     |:-:    |:-:  
小明    |男     |75  
小红    |女     |79  
小强    |男     |89

### 插入公式
行内公式`$\sqrt{x}$`   
行间公式
```math
\sqrt{x}
``` 
[LaTex数学公式使用参考链接](https://blog.csdn.net/testcs_dn/article/details/44229085)  

### 插入链接
![本地图片实例](G:/有道云笔记/我的文件/pic/b.jpg "标题：本地图片")
![网络图片实例](http://p.ssl.qhimg.com/t018ccca46939d3c0bf.jpg "标题：网络图片")
[示例：百度链接](https://www.baidu.com/)
<https://www.baidu.com/>  

### 插入代码块

第一个示例：`printf("Hello world !\n")`  

第二个示例：
```c++
#include <stdio.h>
int main(void)
{
    printf("代码高亮 !\n");    
}
```  

第三个示例：
```html
<table>
 <tr>
  <th rowspan="2">值班人员</th>
  <th>星期一</th>
  <th>星期二</th>
  <th>星期三</th>
 </tr>
 <tr>
  <th>李强</th>
  <th>张明</th>
  <th>王平</th>
 </tr>
</table>
```

第四个示例：
```python
@requires_authorization
def myfunc(param1='', param2=0):
    '''A doc string'''
    if param1 > param2:#interesting
        print 'Greater'
    return (param2 - param1 + 1)or None
class SomeClass:
    pass
>>> message='''interpreter
...prompt'''
```

### 插入流程图
[>_>注释]:< graph 确定流程图的类型，分别为TB、BT、RL、LR> 
[>_>注释]:< start 定义流程图的开始>  
[>_>注释]:< operation 定义一个矩形流程框>  
[>_>注释]:< condition 定义一个判断>  
[>_>注释]:< end 定义流程图的结束>  

```
graph LR
    start[开始] --> operation(流程一)
    operation --> condition{Yes or No?}
    condition --> |no| operation(流程一)
    condition --> |yes| op2((流程二))
    op2 --> ends[结束]
```
[查看更详细的流程图语法，点击这里](http://knsv.github.io/mermaid/#flowcharts-basic-syntax)


### 插入待办事项
- [ ] 已处理的事情1
- [ ] 已处理的事情2
- [x] 未处理的事情
  - [x] 未处理的事情1
  - [x] 未处理的事情2

### 插入序列图
```
sequenceDiagram

Alice->Bob: Hello Bob,how are you?  
Note right of Bob: Bob thinks  
Bob-->Alice: I am good thanks!  
```

### 插入甘特图
```
gantt
dateFormat YYYY-MM-DD
title 项目开发流程  

section 项目确定  
 需求分析:      2016-06-22,3d  
 可行性报告:    5d 
 概念验证:      5d  
section 项目实施  
 概要设计:      2016-07-05,5d  
 详细设计:      2016-07-08,10d  
 编码:          2016-07-15,10d  
 测试:          2016-07-22,5d  
section 发布验收  
 发布:          2d  
 验收:          3d  
```
[查看更详细的甘特图语法，点击这里](http://knsv.github.io/mermaid/#styling39)

## 快捷键
[>_>注释]:<换行需要在行尾添加两个‘空格’>  
加粗	Ctrl + B1  
斜体	Ctrl + I  
引用	Ctrl + Q  
插入链接	Ctrl + L  
插入代码	Ctrl + K  
插入图片	Ctrl + G  
提升标题	Ctrl + H  
有序列表	Ctrl + O  
无序列表	Ctrl + U  
横线	Ctrl + R  
撤销	Ctrl + Z  
重做	Ctrl + Y  



## 附录
HTML字体颜色名列表  
![HTML字体颜色名列表](G:/有道云笔记/我的文件/pic/color.png)

