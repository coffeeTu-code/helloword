``# 原文链接：[JS中的柯里化(currying)](http://www.zhangxinxu.com/wordpress/2013/02/js-currying/)



# 函数式编程 - 柯里化(currying)



## 定义

柯里化（Currying），又称部分求值（Partial Evaluation），是把接受多个参数的函数变换成接受一个单一参数（最初函数的第一个参数）的函数，并且返回接受余下的参数而且返回结果的新函数的技术。

> 柯里化有3个常见作用：
> 1. 提前返回；
> 2. 参数复用；
> 3. 延迟计算/运行。

## 柯里化函数举例

柯里化函数（javascript） 举例一：
```javascript
var addEvent = (function(){
    if (window.addEventListener) {
        return function(el, sType, fn, capture) {
            el.addEventListener(sType, function(e) {
                fn.call(el, e);
            }, (capture));
        };
    } else if (window.attachEvent) {
        return function(el, sType, fn, capture) {
            el.attachEvent("on" + sType, function(e) {
                fn.call(el, e);
            });
        };
    }
})();
```

柯里化函数（go） 举例二：
```go
package main

import (
	"fmt"
	"testing"
)

func say(self string) func(word string) (talk string) {
	return func(word string) (talk string) {
		return "<" + self + ">: " + word
	}
}

var xiaomingSay = say("小明")
var xiaohongSay = say("小红")

func TestSay(t *testing.T) {
	fmt.Println(xiaomingSay("美丽的女士，可以请求你跳一支舞吗？"))
	fmt.Println(xiaohongSay("我的荣幸!"))
	//Output:
	//<小明>: 美丽的女士，可以请求你跳一支舞吗？
	//<小红>: 我的荣幸!
}
```

柯里化函数（javascript） 举例三：
```javascript
var curryWeight = function(fn) {
    var _fishWeight = [];
    return function() {
        if (arguments.length === 0) {
            return fn.apply(null, _fishWeight);
        } else {
            _fishWeight = _fishWeight.concat([].slice.call(arguments));
        }
    }
};
var fishWeight = 0;
var addWeight = curryWeight(function() {
    var i=0; len = arguments.length;
    for (i; i<len; i+=1) {
        fishWeight += arguments[i];
    }
});

addWeight(2.3);
addWeight(6.5);
addWeight(1.2);
addWeight(2.5);
addWeight();    //  这里才计算

console.log(fishWeight);    // 12.5
```

