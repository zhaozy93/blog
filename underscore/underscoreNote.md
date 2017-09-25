# Note

## intro
 - 详细源码分析放在 https://github.com/zhaozy93/blog/blob/master/underscore/underscore_source_code.md
 - 此处仅部分知识点


## void 0
underscore中使用void 0 来代替 undefined
- 代码压缩上面字节更少
- undefined 在低版本浏览器中是可以被更改的
- void一直都属于关键字，不会被更改

## apply 与 call
apply与call用法上的的区别就不再重复了，很基础。

分析一下为什么underscore不直接return要先进行switch
- call的性能比apply快很多倍
- apply要先对参数进行检查和深拷贝
- 可以参考ecma的官方定义 http://www.ecma-international.org/ecma-262/5.1/#sec-15.3.4.3 http://www.ecma-international.org/ecma-262/5.1/#sec-15.3.4.4

```js
  var optimizeCb = function(func, context, argCount) {
    // 如果环境是undefined 则直接返回方法
    if (context === void 0) return func;
    // 其余的后面用到再来讲 现在用不到 理解不了
    // 其实主要也就是把返回的函数封装一层，增加几个参数 方便后面方法传参调用
    switch (argCount) {
      case 1: return function(value) {
        return func.call(context, value);
      };
      // The 2-parameter case has been omitted only because no current consumers
      // made use of it.
      case null:
      case 3: return function(value, index, collection) {
        return func.call(context, value, index, collection);
      };
      case 4: return function(accumulator, value, index, collection) {
        return func.call(context, accumulator, value, index, collection);
      };
    }
    // 这里我们发现一个奇怪的点， 下面这句使用apply直接return似乎完全可以替代上面的switch方法啊
    // 为什么还要switch
    // 原因在于apply的性能比call的性能低太多
    return function() {
      return func.apply(context, arguments);
    };
  };
```
## 表达式不能以 { 或者 function开头
很多时候我们会 [].length， 但是直接执行 {}.valueOf() 会报错， 原因就在于表达式不能以{ 或者 function开头， 但是在下面这行语句中不会报错 因为以！开头

`var hasEnumBug = !{toString: null}.propertyIsEnumerable('toString');`

## 类型判断
在underscore中通过这样的方法来批量产生类型检测方法。
```js
_.each(['Arguments', 'Function', 'String', 'Number', 'Date', 'RegExp', 'Error'], function(name) {
  _['is' + name] = function(obj) {
    return toString.call(obj) === '[object ' + name + ']';
  };
});
```
但后面还有很多额外的补充措施, 例如 Array、NaN、Finite等
```js
_.isArray = Array.isArray || function(obj) {
  return toString.call(obj) === '[object Array]';
};
```

## 可枚举属性 bug
IE9以下且重写了`['valueOf', 'isPrototypeOf', 'toString','propertyIsEnumerable', 'hasOwnProperty', 'toLocaleString']`这些属性，那么这些属性是不会被for...in...正确返回的。

检测方法 `var hasEnumBug = !{toString: null}.propertyIsEnumerable('toString');`

如果有bug的话 那么在for...in...中计算keys就要单独处理了

原本的逻辑判断太复杂了， 简化一下其实就是判断直接读取到的和原型链上面的是不是相等的，如果相等那就是一样的，没有改写 那也就没必要算入keys里面了
```js
// proto = (isFunction(obj.constructor) && obj.constructor.prototype) || Object.prototype
if (obj[prop] !== proto[prop] && !_.contains(keys, prop)) {
      keys.push(prop);
    }
```