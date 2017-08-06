# 终
  jQuery源码分析从[开始](https://github.com/zhaozy93/blog/issues/2)(17-05-30)到[现在](https://github.com/zhaozy93/blog/issues/20)(17-08-05)整整两月零五日，中间也因为工作暂停过，也高兴过头刷夜过，一路走来断断续续，好在没有放弃。回头看看写过的那些无意义无价值的注释恰恰是成长的过程。 单看某一段代码，某一段注释可能并不复杂，但当把整个项目系统的看过一遍以后，才发现原来一环扣一环，它的风格从一而终、它的写法自古不变，里面基本没有复杂的逻辑，有的只是更全更广的基础知识点。
  
  - 没经历过你不知道cssFloat、StyleFLoat的区别
  - 没经历过你不知道Regex.source是什么
  - 没经历过你不知道runtimeStyle又是什么

  只有经历过才知道知识的匮乏，外面的世界更大，自己知道的更少，认识到更渺小的自己。
  读过jQuery可能要比阅读一本生硬的设计模式会让你了解的更多，更重要的是你知道别人是如何设计，一天天一段段代码也会让你潜移默化的在日后给你这好像曾经在jQuerry中看到过、有过类似代码的感受。

  虽然并不是真正的对jQuery全部分析，但核心core部分确实基本全部走过，感谢前辈们的各种经验、blog、问题回答让后人在遇到问题时会有更清楚的认识，特别感谢[高云(nuysoft)](https://github.com/nuysoft)老师的jQuery技术内幕这本书，本人基本是按照书本的流程和帮助才能走完这晦涩无趣生硬的源码。[mock.js](http://mockjs.com/)是nuysoft老师的又一力作，在发现这个事实后也去搂了一眼mock的源码，发现里面的注释依旧那么清楚， 也许这就是习惯的力量吧。

  虽然标题是终，但前端的路途、学习的路途不可能终止，即使是jQuery的学习也不可能终止。这只是一小段旅途的休息站，日后的万里长城更等待着、静候着我去探索。

## 如何检测dom元素
  target.nodeType 

## 类数组对象
  下标0，1，2，3数字且拥有length属性(arguments对象是一个很好的例子)
  a = { 0: 1, 1: 2, 2: 3, length: 3};

## 类数组对象转换真正的数组
  利用slice方法，但原对象的key一定是0，1，2，3数字， 且必须有length属性。
  a = { 0: 1, 1: 2, 2: 3, length: 3};
  b = Array.prototype.slice.call(a, 0);

## apply和call
  虽然都是改变方法执行的this值，第一个参数都是this，但是apply只接受一个多余参数，call可以接受多个。
  apply会把第二个参入如果是数组扁平化一次，因此可以利用这个特性来实现二维数组转一位数组。
  Array.prototype.push.apply({}, [1,2,3])  --> { 0:1, 1:2, 2:3, length:3}
  Array.prototype.push.call({}, [1,2,3,4,5])  --> { 0:[1,2,3], length:1}

## 数字与字符串快速转变
  i = +i;
  i = i + ''

## 数字类型的判断方式
  function isNumeric( obj ) {
    return !isNaN( parseFloat(obj) ) && isFinite( obj );
  }

## 数字字符串的正则匹配
  /-?\d*(?:\.\d*)?(?:[eE][+\-]?\d+)?/g
  -?号可有可无  \d数字不限(可以直接.5)  (?:\.\d*)? 小数部分可有可无  (?:[eE][+\-]?\d+)? 科学计数部分可有可无

## 函数的四种创建方式
  function func () {}  // 普通方式
  let func = function(){}   // 匿名函数
  new Function ( string )  // 利用字符串构建函数
  还有一种类似函数的方法
  eval('(' + string + ')')  // eval利用

  new Function可以创建可复用的函数， eval仅仅是立即执行传入的字符串内容，
  stackoverflow 有一个不错的对比问题 [website](https://stackoverflow.com/questions/4599857/are-eval-and-new-function-the-same-thing)

## 函数执行上下文的一段解释， 翻译自 [stackoverflow](https://stackoverflow.com/questions/6796521/what-is-the-context-of-an-anonymous-function)
  Functions don't inherit the context they are called in. Every function has it's own this and what this refers to is determined by how the function is called. It basically boils down to:
  函数的执行上下文不会继承自调用他们的环境。每个函数都有自己的作用上下文，他们执行上下文指向取决于函数如何被调用，基本规则如下：
  "Standalone" ( func()) or as immediate function ( (function(){...}()) ) : this will refer to the global object (window in browsers)
  独立的函数或者立即执行函数： this将指向全局对象(在浏览器中即window)
  As property of an object ( obj.func()): this will refer to obj
  作为属性的函数(obj.func)： this指向obj
  With the new [docs] keyword ( new Func() ): this refers to an empty object that inherits from Func.prototype
  使用new Func()关键字创建的：this指向一个空对象，空对象继承自Func.prototype
  apply [docs] and call [docs] methods ( func.apply(someObj) ): this refers to someObj
  apply和call方法调用的函数：this执行传入的对象someObj

## String.prototype.replace
  str.replace(regexp|substr, newSubStr|function)
  可以看出第一个参数 正则|子串 很好理解，
  第二个参数 新子串|函数， 新子串也比较好理解， 函数好像不经常用到
  function(substr, [group1, group2, group3...], startIndex, str)
  函数的默认传入参数是这样构成的，很奇怪。 第一个参数是匹配到的子串， 倒数第二个参数是匹配到的子串开始index，最后一个参数是整个字符串， 中间部分可能没有，也可能有很多个。 主要是由正则中group的个数决定
i
## call、apply、bind
  1、三者都可以更改函数执行的上下文
  2、call、apply会立即执行函数，bind返回一个新函数
  3、apply只接受2个参数、call参数不限制， 同时apply会把第二个参数(数组格式)进行一个扁平化
  4、bind也不限制参数个数，但新函数中bind传入的参数会排在最前面，原函数的参数依旧向后移

## Regex.source
  日常中可能很少碰到需要进行正则拼接，或者修改某个已经定义的正则。
  但如果碰到 可以使用Regex.source 来获取正则的字符串形式。
  下面是个例子， 修改一个正则再生产一个新的正则
  /[\w]?/.source  --> '[\w]?'
  /[\w]?/.source.replace(/\?/, '*')  --> '[\w]*'
  new RegExp(/[\w]?/.source.replace(/\?/, '*'))  --> /[\w]*/

## String split 正则
  如果需要用 `空格` 来分割字符串， 第一反应可能是 str.split( ' ' )
  但str如果不太标准，例如有 `'  '`两个空格连在一起的情况, 是不是考虑用正则呢 str.split( /\s+/ )
  'qwer           asdf'.split(' ') --> ["qwer", "", "", "", "", "", "", "", "", "", "", "asdf"]
  'qwer           asdf'.split(/\s+/) --> ["qwer", "asdf"]
  尽管这也不是万能的，比如 '       '.split(/\s+/) --> ["", ""]， 但总是好一点点，
  换种思路再来一遍！

## innerHTML、innerText、textContent、createTextNode
  首先所有innerHTML与outerHTML类似，在读取时返回的字符串中包含子代标签信息，outter包含自身标签
  innerText和outerText都不会包含任何标签信息
  在设置时innerHTML、outerHTML会解析字符串内的转义字符、标签信息等
  innerText、outerText不会解析内容，直接字符串设置
  textContent行为与innerText基本一致
  createTextNode是创建文本标签，也不会对内部内容进行解析
  当然兼容性不在讨论范围内

## position
  一个元素只设置position是没有意义的，必须要有top或者left， 否则容易出现bug

## getAttributeNode、setAttributeNode与getAttribute、setAttribute
  前两者是获取属性节点，是一个对象，后两者是直接设置属性或读取属性，虽然目的是相同的，但后两者在IE67有很大问题，前两者可以很好的解决这个问题。

## float属性
  float是css中一个很好的、很常用的属性，但是在获取设置float值的时候却不是那么简单的
  在低版本ie中要读取float的姿势是 `elem.style.styleFloat`
  ie9及其以后 `elem.style.styleFloat` 和 `elem.style.cssFloat`都有
  在chrome、ff中正确姿势是 `elem.style.cssFloat`(实测高版本chrome`elem.style.float`也可以)
  但是呢？ 上面两种方法都是基于 `elem.style`的，大家都知道这只能获取内连style
  在获取计算float值时呢，getComputedStyle的作用就体现出来了
  `window.getComputedStyle(elem).cssFloat`就可以取到计算后的样式了，同时高版本的chrome也可以`window.getComputedStyle(elem).float`, 甚至可以`window.getComputedStyle(elem).getPropertyValue('float')`, 如果使`用getPropertyValue`就只能使用`float`，因为要求传入的是css属性名， 同时传入的应该是'-'分隔的属性名而非驼峰规则。
  但但是？ IE没有`window.getComputedStyle`,只有`elem.currentStyle('float')`。或者呢使用`elem.currentStyle.getAttribute('float')`,使用getAttribute也是不需要考虑怪异名称的，只是需要驼峰属性名即可。
  因此float属性读取和设置的时候还是要做if else判断的

## opacity属性
  低版本IE8及其以下是通过filter滤镜属性来实现opacity的， 同时ie中让filter生效的另一个因素就是要有layout属性， 触发layout的方法就是设置zoom值， `zoom = 1`
  `filter: alpha(opacity=50);` 
  当然现在基本不用考虑IE8及其以下了，但是filter的功能却不知事opacity透明度那么简单， 毕竟filter在奇遇浏览器也有这个属性！！！
  filter有14中滤镜
  ```
  Alpha     让HTML元件呈现出透明的渐进效果
  Blur     让HTML元件产生风吹模糊的效果
  Chroma     让图像中的某一颜色变成透明色
  DropShadow     让HTML元件有一个下落式的阴影
  FlipH     让HTML元件水平翻转
  FlipV     让HTML元件垂直翻转
  Glow     在元件的周围产生光晕而模糊的效果
  Gray     把一个彩色的图片变成黑白色
  Invert     产生图片的照片底片的效果
  Light     在HTML元件上放置一个光影
  Mask     利用另一个HTML元件在另一个元件上产生图像的遮罩
  Shadow     产生一个比较立体的阴影
  Wave     让HTML元件产生水平或是垂直方向上的波浪变形
  XRay     产生HTML元件的轮廓，就像是照X光一样
  ```
  Alpha滤镜又有多种参数
  ```
  Alpha 滤镜参数详解
  参数名     说明     取值说明 
  Opacity     不透明的程度，百分比。    从0到100，0表是完全透明，100表示完全不透明。
  FinishOpacity     这是一个同Opacity一起使用的选择性的参数，当同时Opacity和FinishOpacity时，可以制作出透明渐进的效果，比较酷。    从0到100，0表是完全透明，100表示完全不透明。
  Style     当同时设定了Opacity和finishOpacity产生透明渐进时，它主要是用赤指定渐进的显示形状。    0：没有渐进；1：直线渐进；2：圆形渐进；3：矩形辐射。
  StartX     渐进开始的 X 坐标值    
  StartY     渐进开始的 Y 坐标值    
  FinishX     渐进结束的 X 坐标值    
  FinishY     渐进结束的 Y 坐标值    
  ```
## runtimeStyle
  在IE中style对象有三类
  elem.style: 获取内连样式属性
  elem.currentStyle: 类似于getComputedStyle获取当前元素最终的属性
  elem.runtimeStyle: 顾名思义，运行时样式。
  前两者容易理解，那么什么是运行时样式呢
  [关于HTML Object中三个Style实例的区别](http://birdshome.cnblogs.com/archive/2005/01/16/92491.html)
  [关于使用runtimeStyle属性问题讨论](http://www.cnblogs.com/birdshome/archive/2006/04/24/runtimeStyle.html)
  这里恰巧有位m$的老哥的blog，讲解了runtimeStyle的个人理解。

  runtimeStyle的修改不会立即同步到style，也不会立即展现到页面元素。 在jQuery中就是利用这个特性来计算 em\%单位向px单位转换的。（pixelLeft属性帮助计算px）
  上面文章也提到了 
  - runtimeStyle属性一定要配对使用，即element.runtimeStyle.xxx = 'attribue';和element.runtimeStyle.xxx = '';配对


## document对象
  document.defaultView 返回window对象
  document.documentElement 返回html元素
  docuemnt.body 返回body元素
  document.ownerDocument   null
  document.getRootNode  ===  document
  document.compatMode 表明当前文档的渲染模式是混杂模式还是"标准规范模式".
                  BackCompat"代表"混杂模式", "CSS1Compat"代表"标准规范模式".
  document.designMode 整个文档是否可编辑

  elem.ownerDocument 返回document对象
  elem.getRootNode   返回document对象
  elem.getRootNode() == elem.ownerDocument
