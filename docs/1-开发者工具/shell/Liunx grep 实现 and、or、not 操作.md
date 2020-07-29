# grep -v、-e、-E


> [出处](https://www.cnblogs.com/franjia/p/4384362.html)


[TOC]



## 概述

在Linux的grep命令中如何使用or，and，not操作符呢？

其实，在grep命令中，有or和not操作符的等价选项，但是并没有grep and这种操作符。不过呢，可以使用patterns来模拟and操作的。下面会举一些例子来说明在Linux的grep命令中如何使用or，and，not。

在下面的例子中，会用到这个employee.txt文件，如下：

```
$ cat employee.txt  
100  Thomas  Manager    Sales       $5,000  
200  Jason   Developer  Technology  $5,500  
300  Raj     Sysadmin   Technology  $7,000  
400  Nisha   Manager    Marketing   $9,500  
500  Randy   Manager    Sales       $6,000  
```


## grep or 操作符


以下四种方法均能实现grep OR的操作。个人推荐方法3。

### 使用 `\|`

如果不使用grep命令的任何选项，可以通过使用 '|' 来分割多个pattern，以此实现OR的操作。

```
grep 'pattern1\|pattern2' filename  
```

例子如下：

```
$ grep 'Tech\|Sales' employee.txt  
100  Thomas  Manager    Sales       $5,000  
200  Jason   Developer  Technology  $5,500  
300  Raj     Sysadmin   Technology  $7,000  
500  Randy   Manager    Sales       $6,000  
```

### 使用选项 `-E`

grep -E 选项可以用来扩展选项为正则表达式。 如果使用了grep 命令的选项-E，则应该使用 | 来分割多个pattern，以此实现OR操作。

```
grep -E 'pattern1|pattern2' filename
```

例子如下：

```
$ grep -E 'Tech|Sales' employee.txt  
100  Thomas  Manager    Sales       $5,000  
200  Jason   Developer  Technology  $5,500  
300  Raj     Sysadmin   Technology  $7,000  
500  Randy   Manager    Sales       $6,000 
```

### 使用 `egrep`

egrep 命令等同于‘grep -E’。因此，使用egrep (不带任何选项)命令，以此根据分割的多个Pattern来实现OR操作。

```
egrep 'pattern1|pattern2' filename  
```

例子如下：

```
$ egrep 'Tech|Sales' employee.txt  
100  Thomas  Manager    Sales       $5,000  
200  Jason   Developer  Technology  $5,500  
300  Raj     Sysadmin   Technology  $7,000  
500  Randy   Manager    Sales       $6,000  
```

### 使用选项 `-e`

使用grep -e 选项，只能传递一个参数。在单条命令中使用多个 -e 选项，得到多个pattern，以此实现OR操作。

```
grep -e pattern1 -e pattern2 filename
```

例子如下：

```
$ grep -e Tech -e Sales employee.txt  
100  Thomas  Manager    Sales       $5,000  
200  Jason   Developer  Technology  $5,500  
300  Raj     Sysadmin   Technology  $7,000  
500  Randy   Manager    Sales       $6,000  
```


## grep and 操作


### 使用 `-E 'pattern1.*pattern2'`

grep命令本身不提供AND功能。但是，使用 -E 选项可以实现AND操作。

```
grep -E 'pattern1.*pattern2' filename  
grep -E 'pattern1.*pattern2|pattern2.*pattern1' filename 
```

第一个例子如下：（其中两个pattern的顺序是指定的）

```
$ grep -E 'Dev.*Tech' employee.txt  
200  Jason   Developer  Technology  $5,500 
```

第二个例子：（两个pattern的顺序不是固定的，可以是乱序的）

```
$ grep -E 'Manager.*Sales|Sales.*Manager' employee.txt  
```

### 使用多个grep命令

可以使用多个 grep 命令 ，由管道符分割，以此来实现 AND 语义。

```
grep -E 'pattern1' filename | grep -E 'pattern2'  
```

例子如下：

```
$ grep Manager employee.txt | grep Sales  
100  Thomas  Manager    Sales       $5,000  
500  Randy   Manager    Sales       $6,000  
```


## grep not 操作


### 使用选项 `grep -v`

使用 grep -v 可以实现 NOT 操作。 -v 选项用来实现反选匹配的（ invert match）。如，可匹配得到除下指定pattern外的所有lines。

```
grep -v 'pattern1' filename
```

例子如下：

```
$ grep -v Sales employee.txt  
200  Jason   Developer  Technology  $5,500  
300  Raj     Sysadmin   Technology  $7,000  
400  Nisha   Manager    Marketing   $9,500  
```

可以将NOT操作与其他操作联合起来，以此实现更强大的功能 组合。

如，这个例子将得到：‘Manager或者Developer，但是不是Sales’的结果：

```
$ egrep 'Manager|Developer' employee.txt | grep -v Sales  
200  Jason   Developer  Technology  $5,500  
400  Nisha   Manager    Marketing   $9,500 
```