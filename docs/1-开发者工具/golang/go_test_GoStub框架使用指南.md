# GoStub框架使用指南

[TOC]

## 序言

在软件开发中，产品代码的正确性通过测试代码来保证，而测试代码的正确性谁来保证？答案是毫无争议的，肯定是程序员自己。这就要求测试代码必须足够简单且表达力强，让错误无处藏身。我们要有一个好鼻子，能够嗅出测试的坏味道，及时的进行测试重构，从而让测试代码易于维护。

要写出好的测试代码，必须精通相关的测试框架。对于Golang的程序员来说，至少需要掌握下面四个测试框架：

- GoConvey
- GoStub
- GoMock
- Monkey

本文将主要介绍**GoStub框架**的基本使用方法，从而指导读者更好的进行测试实践，最终写出简单优雅的测试代码。

## GoStub简介

通过上一篇文章《GoConvey框架使用指南》的学习，大家熟悉了GoConvey框架的基本使用方法，虽然已经可以写出简单优雅的测试代码，但是如果在被测函数中调用了底层操作函数，比如调用了os包的Stat函数，则需要在测试函数中先对该底层操作函数打桩。那么，该如何对函数高效的打桩呢？

本文给大家介绍一款轻量级的GoStub框架，接口友好，可以对全局变量、函数或过程打桩，我们一起来体验一下。

## 安装

在命令行运行下面的命令：

```
go get github.com/prashantv/gostub
```

运行完后你会发现，在$GOPATH/src/github.com目录下，新增了prashantv/gostub子目录，这就是本文的主角。

## 基本使用方法

GoStub框架的使用场景很多，依次为：

1. 基本场景：为一个全局变量打桩
2. 基本场景：为一个函数打桩
3. 基本场景：为一个过程打桩
4. 复合场景：由任意相同或不同的基本场景组合而成

### 为一个全局变量打桩

假设num为被测函数中使用的一个全局整型变量，当前测试用例中假定num的值大于100，比如为150，则打桩的代码如下：

```go
stubs := Stub(&num, 150)
defer stubs.Reset()
```

stubs是GoStub框架的函数接口Stub返回的对象，该对象有Reset操作，即将全局变量的值恢复为原值。

### 为一个函数打桩

假设我们产品的既有代码中有下面的函数定义：

```go
func Exec(cmd string, args ...string) (string, error) {
    ...
}
```

则Exec函数是不能通过GoStub框架打桩的。

若要想对Exec函数通过GoStub框架打桩，则仅需对该函数声明做很小的重构，即将Exec函数定义为匿名函数，同时将它赋值给Exec变量，重构后的代码如下：

```
var Exec = func(cmd string, args ...string) (string, error) {
    ...
}
```

说明：对于新增函数，请按上面的方式定义

当Exec函数重构成Exec变量后，丝毫不影响既有代码中对Exec函数的调用。由于Exec变量是函数变量，所以我们一般将这类变量也叫做函数。

现在我们可以对Exec函数打桩了，代码如下所示：

```
stubs := Stub(&Exec, func(cmd string, args ...string) (string, error) {
            return "xxx-vethName100-yyy", nil
})
defer stubs.Reset()
```

其实GoStub框架专门提供了StubFunc函数用于函数打桩，我们重构打桩代码：

```go
stubs := StubFunc(&Exec,"xxx-vethName100-yyy", nil)
defer stubs.Reset()
```

产品代码中很多函数都会调用Golang的库函数或第三方的库函数，我们又不能重构这些函数，那么该如何对这些库函数打桩？

答案很简单，即在适配层中定义库函数的变量，然后在产品代码中使用该变量。

定义库函数的变量：

```
package adapter

var Stat = os.Stat
var Marshal = json.Marshal
var UnMarshal = json.Unmarshal
...
```

使用UnMarshal的代码：

```go
bytes, err := adapter.Marshal(&student)
if err != nil {
    ...
    return err
}
...
...
```

我们现在可以对库函数进行打桩了。假设当前使用的库函数为Marshal，因为Marshal函数有成功或失败两种情况，所以它有两个桩函数，但对于每一个测试用例来说Unmarshal只有一个桩函数。

序列化成功时的打桩代码为：

```
var liLei = `{"name":"LiLei", "age":"21"}`
stubs := StubFunc(&adapter.Marshal, []byte(liLei), nil)
defer stubs.Reset()
```

序列化失败时的打桩代码为：

```
stubs := StubFunc(&adapter.Marshal, nil, ErrAny)
defer stubs.Reset()
```

### 为一个过程打桩

当一个函数没有返回值时，该函数我们一般称为过程。很多时候，我们将资源清理类函数定义为过程。

我们对过程DestroyResource的打桩代码为：

```go
stubs := StubFunc(&DestroyResource)
defer stubs.Reset()
```

### 任意相同或不同的原子场景的组合

不论是调用Stub函数还是StubFunc函数，都会生成一个stubs对象，该对象仍然有Stub方法和StubFunc方法，所以在一个测试用例中可以同时对多个全局变量、函数或过程打桩。这些全局变量、函数或过程会将初始值存在一个map中，并在延迟语句中通过Reset方法统一做回滚处理。

假设Sf为Stub或StubFunc函数的调用，Sm为Stub或StubFunc方法的调用，则在一个测试用例中使用GoStub框架的打桩代码为：

```
stubs := Sf
defer stubs.Reset()
stubs.Sm1
...
stubs.SmN
```

> 不推荐将打桩代码写成下面的形式：
```
stubs := Sf
defer stubs.Sm1.(...).SmN.Reset()
```

## TestFuncDemo

笔者在上一篇文章《GoConvey框架使用指南》中推荐读者使用Convey语句的嵌套，即一个函数有一个测试函数，测试函数中嵌套两级Convey语句，第一级Convey语句对应测试函数，第二级Convey语句对应测试用例。在第二级的每个Convey函数中都会产生一个stubs对象，彼此独立，互不影响。

我们看一个针对GoStub框架使用的较为完整的测试函数Demo：

```go
func TestFuncDemo(t *testing.T) {
    Convey("TestFuncDemo", t, func() {
        Convey("for succ", func() {
            stubs := Stub(&num, 150)
            defer stubs.Reset()
            stubs.StubFunc(&Exec,"xxx-vethName100-yyy", nil)
            var liLei = `{"name":"LiLei", "age":"21"}`
            stubs.StubFunc(&adapter.Marshal, []byte(liLei), nil)
            stubs.StubFunc(&DestroyResource)
            //several So assert
        })

        Convey("for fail when num is too small", func() {
            stubs := Stub(&num, 50)
            defer stubs.Reset()
            //several So assert
        })

        Convey("for fail when Exec error", func() {
            stubs := Stub(&num, 150)
            defer stubs.Reset()
            stubs.StubFunc(&Exec, "", ErrAny)
            //several So assert
        })

        Convey("for fail when Marshal error", func() {
            stubs := Stub(&num, 150)
            defer stubs.Reset()
            stubs.StubFunc(&Exec,"xxx-vethName100-yyy", nil)
            stubs.StubFunc(&adapter.Marshal, nil, ErrAny)
            //several So assert
        })

    })
}
```

## 不适用的复杂情况

尽管GoStub框架已经可以优雅的解决很多场景的函数打桩问题，但对于一些复杂的情况，却只能干瞪眼：

1. 被测函数中多次调用了数据库读操作函数接口 ReadDb，并且数据库为key-value型。被测函数先是 ReadDb 了一个父目录的值，然后在 for 循环中读了若干个子目录的值。在多个测试用例中都有将ReadDb打桩为在多次调用中呈现不同行为的需求，即父目录的值不同于子目录的值，并且子目录的值也互不相等

2. 被测函数中有一个循环，用于一个批量操作，当某一次操作失败，则返回失败，并进行错误处理。假设该操作为Apply，则在异常的测试用例中有将Apply打桩为在多次调用中呈现不同行为的需求，即Apply的前几次调用返回成功但最后一次调用却返回失败

3. 被测函数中多次调用了同一底层操作函数，比如 exec.Command，函数参数既有命令也有命令参数。被测函数先是创建了一个对象，然后查询对象的状态，在对象状态达不到期望时还要删除对象，其中查询对象是一个重要的操作，一般会进行多次重试。在多个测试用例中都有将 exec.Command 打桩为多次调用中呈现不同行为的需求，即创建对象、查询对象状态和删除对象对返回值的期望都不一样

## 小结

GoStub是一款轻量级的测试框架，接口友好，可以对全局变量、函数或过程打桩。本文详细阐述了GoStub框架的使用场景，并给出了一个较为完整的测试函数Demo，希望读者能够掌握GoStub框架的基本使用方法，提高单元测试水平，交付高质量的软件。