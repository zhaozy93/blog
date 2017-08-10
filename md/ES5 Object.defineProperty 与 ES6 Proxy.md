# ES5 Object.defineProperty 与 ES6 Proxy

先来各自看一下基础知识

## Object.defineProperty
`Object.defineProperty`在[MDN](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/defineProperty)中这样介绍它
> Object.defineProperty() 方法会直接在一个对象上定义一个新属性，或者修改一个对象的现有属性， 并返回这个对象。
说了和没说一样， 了解一下就知道`Object.defineProperty`帮我们解决的事getter、setter的劫持问题，何为劫持呢？举例

``` js
let obj = {};  // 我们定义了一个obj对象
obj.name = 'hello'; // 我们调用setter为obj设置了name属性，且值为hello
console.log(obj.name); // 我们调用getter来获取obj的name属性对应的值
```
我们劫持了getter、setter，言外之意就是我们可以获得用户设置或者读取某个属性的事件(暂且理解为这是一种事件)， 我们就可以在这个事件中作自己爱做的事情了。

``` js
let obj = {};
Object.defineProperty(obj, 'name', {
  get(){
    return this._name || 'empty_name';
  },
  set(newValue){
    if(newValue === this._name){
      return;
    }
    this._name = newValue;
  }
});
obj.name;
obj.name = 'hello';
```
![img1](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/defineProperty-1.jpeg)
看chrome下的控制台可以看到，obj对象上面的name属性不是一个普通的字符串，当你点击三个点时，他就又变成字符串了，这个很神奇，可以认为通过defineProperty劫持后的属性是个懒属性吧，需要主动去触发。

关于`Object.defineProperty`参数尤其是第三个description参数的问题可以参考MDN、与其余几片文章。
- [解析神奇的 Object.defineProperty](http://blog.csdn.net/fengyinchao/article/details/55805966)
- [Object.defineproperty实现数据和视图的联动](http://www.cnblogs.com/oceanxing/p/3938443.html)
- [vue.js关于Object.defineProperty的利用原理](http://www.jianshu.com/p/07ba2b0c8fca)
- [MVVM双向绑定实现之Object.defineProperty](http://www.mamicode.com/info-detail-1159006.html)
- [Vue.js双向绑定实现原理详解](http://www.jb51.net/article/100815.htm)
## Proxy
前面讲到了`Object.defineProperty`大概做的是劫持的事情，劫持对象的某个属性的读与写事件。那么我们就可以认为从ES5到ES6，JS的功力又大增了一次，潜心修炼N年。

`Proxy`可以劫持更多的事件了，也可以简单认为`Proxy`是`Object.defineProperty`的升级加强版。

[es6 javascript的Proxy实例的方法](http://blog.csdn.net/qq_30100043/article/details/53443017)这篇文档对基本的`Proxy`用法进行了讲解。

还是拿劫持setter、getter来讲，在`Proxy`中，升级到一次劫持整个对象的getter、setter而不再是某个属性了。
```js
var person = {  
    name: " 张三 "  
};
var proxy = new Proxy(person, {
  get: function(target, property) {
    if(property in target) {
      return target[property];
    } else {
      throw new ReferenceError("Property \"" + property + "\" does not exist.");
    }
  }
});
proxy.name // " 张三 "  
proxy.age //  抛出一个错误 
```
![img2](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/proxy-1.jpeg)
这一次我们发现getter的劫持不再局限于某个属性，而是一次劫持整个对象。 同时在`Object.defineProperty`有一个缺陷就是 如果我们劫持了name这个属性的getter，那么我们在get方法内部时不能使用`obj.name`,聪明的你肯定想到了`this.name`,但是也不可以。 所以上面的例子才采用_name来真正存储数据。但是这一次在Proxy方法下，我们不再需要了_name这种迂回战术了，因为我们其实不是对target本身操作，而是对他的劫持器(Proxy)操作。

Proxy不仅能劫持最普通的getter、setter甚至还可以劫持函数方法的调用、new指令、删除属性，甚甚甚甚Proxy可以监控对象的`defineProperty`.....这就非常尴尬了，大哥直接监控了小弟。

## Vue
在介绍 `Object.defineProperty` 的时候给了几个blog地址，发现它们的标题都和vue有关系，毕竟那么火， 原因就在于Vue作为一个MVVM框架，它最基本双向绑定就是通过`Object.defineProperty`来实现的。
``` html
<div>
  <p>你好，<span id='nickName'></span></p>
  <div id="introduce"></div>
</div>
```
``` js
//视图控制器
var userInfo = {};
Object.defineProperty(userInfo, "nickName", {
  get: function(){
    return document.getElementById('nickName').innerHTML;
  },
  set: function(nick){
    document.getElementById('nickName').innerHTML = nick;
  }
});
Object.defineProperty(userInfo, "introduce", {
  get: function(){
    return document.getElementById('introduce').innerHTML;
  },
  set: function(introduce){
    document.getElementById('introduce').innerHTML = introduce;
  }
})
```
通过上面两个最基本的`Object.defineProperty` 我们就实现了 只关心数据 不关心DOM元素。当我们更新数据的时候，其实我们时劫持了更新事件，然后我们用DOM元素来作为真正的数据容器，但是呢， 这好像不是双向绑定。 那怎么实现双向绑定呢？

``` html
<input id='input'/>
<span id='text'></span>
```
``` js
let userInput = {};
Object.defineProperty(userInput, 'text', {
  get(){
    return document.getElementById('text').innerHTML
  },
  set(newValue){
    document.getElementById('text').innerHTML = newValue;
    document.getElementById('input').value = newValue
  }
});
document.getElementById('input').addEventListener('input', function(E){
  userInput.text = E.target.value;
})
```
然后神奇的双向绑定就出现了。 当然Vue做的不仅仅是这些，毕竟日常工作中基本不接触vue，全铺在react上，这篇文档也算是对新事物的了解，以及为接下来react的分析开个头。。。 xx__xx

## 参考
[ES6新特性：Proxy代理器](http://www.cnblogs.com/diligenceday/p/5467634.html)