0、 基本数据类型
  Undefined、Null、Boolean、Number、String
1、parseInt()、 Number()、parseFloat()
三个方法都属于window
Number():
  - Boolean 转为 1、0
  - null转为0
  - undefined 转为 NaN
  - 任何不符合要求都转为NaN
  - 空字符串转为 0
  - 无参数也转为 0
parseInt(string, radix):
  - 空字符串为NaN
  - 无参数转为 NaN
  - 尽可能的去转换， 尽可能的不报NaN  parseInt('23xd') --> 23
  - 使用时尽可能的传入radix
parseFloat()
  - 只解析十进制
  - 16进制始终转换为 0

2、 Number.prototype.toString(radix)

  尽量传入radix以确保正确的十进制

3、String

  字符串是不可变的，也就是说字符串一旦创建，值就不能改变
  `a = 'a'; a= a + 'b'`
  先创建更大的新字符串，完成拼接工作，再销毁掉原本'a'。

4、Object.isPrototypeOf()

  用于检测传入的对象是否是传入对象的原型
  ``` js
  let a = {}, b = Object.create(a);
  a.isPrototypeOf(b)   // true
  ```
5、Object.propertyIsEnumerable(propertyName)

  检查属性是否支持枚举、是否能for-in出来

6、 ++、--操作转换

  有效的数值字符串 --> 数字

  无效的数值字符串('23rdf') --> NaN

  false --> 0

  true --> 1

  object依次调用valueOf() --> toString()

7、位运算
  ES规定所有的数值都是以IEEE-754 64位格式存储的，但真正操作中后台会转为32位的整数，然后执行操作，最后再转换回64位，因此开发人员是无法感知64位存储格式的。

  其中前31位是整数的值， 32位表示符号位， 0是整数，1是负数

  正数以纯二进制形式存储
  18:   10010

  负数以补码的形式储存
  -18: 1111 1111 1111 1111 1111 1111 1110 1110

  补码的计算步骤
  1) 绝对值的二进制码
  2) 二进制的反码，即0和1互换
  3) 反码 + 1

  正是由于64位与32位之间存在转换关系，导致一个bug，即NaN与Infinity在2进制计算过程中会被当作0(至于为啥不知道...)
  ～NaN === ～Infinity === -1

  ～按位非 求反码(变换符号再减1)
   ～25 === -26
   ～-25 === 24

  & 按位与 

  | 按位或  

  ^ 按位异或

  << 左移 将数值的所有位向左移动指定位数，右侧用0来补齐，左移不会影响符号位
  
  \>\> 右移 与左移刚好相反，右移出来的空位用0来补充，也不影响符号位
  
  \>\>\> 无符号右移 不保留符号位的右移
    导致正数的 \>\>\> 与 \>\> 效果一致
    负数会变为正数，因为符号位被0填充，同时会非常巨大，因为负数原本以补码形式存储

8、 *转换规则
  - 超出范围为 Infinity 或 -Infinity
  - 有一方为NaN， 结果为NaN
  - Infinity * 0， 结果为NaN
  - Infinity * 非0， 结果为Infinity加符号位计算
  - Infinity * Infinity， 结果Infinity
  - 非数值，调用Number()再计算

9、 ／转换规则
  - 超出范围为 Infinity 或 -Infinity
  - 有一方为NaN， 结果为NaN
  - 0 / 0, 结果NaN
  - 0 / Infinity, 结果0加符号位
  - 非0 / 0， 结果为Infinity加符号位
  - 非0 / Infinity, 结果为Infinity加符号位
  - Infinity / Infinity， 结果NaN
  - 非数值，调用Number()再计算

10、% 取余转换规则
  - Infinity / 有限数， 结果为 NaN
  - 有限数 /  Infinity， 结果为有限数
  - Infinity ／ Infinity， 结果位 NaN
  - any / 0, 结果位NaN
  - 非数值，调用Number()再计算

11、 +号规则

   都是数值:
  - 一方为NaN，结果为NaN
  - Infinity + -Infinity  = NaN
  - +0 + -0 = +0

  有字符串时，双方都转为字符串再做拼接

  对象、数组、布尔调用toString()， undefined、null调用String()方法获取对应字符串

  但是！！！

  [] + {} + undefined + null  -->  "[object Object]undefinednull"
  
  undefined + null --> NaN

12、 - 规则
  - 一方为NaN，结果为NaN
  - Infinity - -Infinity  = NaN
  - -Infinity - Infinity  = NaN
  - +0 - -0 = -0
  - -0 - -0 = +0
  字符串、布尔、null、undefined调用Number()
  对象一次调用valueOf() toString()

13、比较符号 > < >= <=
  - 双方数值 直接比较
  - 双方字符串， 对比字符编码值
  - 一方数值，另一方转为数值再比较
  - 任何数值和NaN比较都会得到false
  ``` js
  'Brick' < 'alphabet'  // true B: 66  a:97
  '23' < '3' // true 
  ```

14、 ==
  == 在转换时会尽量向数值上面靠，
  但是在向数值转换前有几个要求
  - null == undefined
  - null 和 undefined 不应被转换为任何数值
  - NaN 出现则返回false
  - 对象则判断是不是引用同一个实体

15、var与 no var

  es5中 不使用var会将变量定义到全局环境下，即window... 容易搞事情

16、label配合for循环 ！！！废弃 不推荐使用

  类似于c语言的go语法
  ``` js
  outermost:
    for(var i=0 ;i< 10; i++){
      for(var j=0 ;j< 10; j++){
        if( j == 5 && i == 5){
          break outermost;
        }
        console.log( '' + i + j)
      }
    }
  ```
  break outermost;起到一次性跳出两个循环的作用，最后只会打印到54
  continue outermost; 起到一次性跳出两个循环的作用，但是还会进去的。。。。 也就是说打印完54，下一次就是60了

17、 with(object) ！！！废弃 不推荐使用, 严格模式直接报错
  ```js
  with(obj){
    console.log(key);
  }
  ```
  在with语句块内，变量会先寻找局部变量，找不到在寻找obj上面是否有同名属性。

  是不是有点全局变量window的意思，但是不可以使用，性能差，而且不支持严格模式。

18、 原始数据类型和引用数据类型 内存区别
  栈：原始数据类型（Undefined，Null，Boolean，Number、String）
  堆：引用数据类型（对象、数组和函数）

 两种类型的区别是：存储位置不同；
 原始数据类型直接存储在栈(stack)中的简单数据段，占据空间小、大小固定，属于被频繁使用数据，所以放入栈中存储；
 引用数据类型存储在堆(heap)中的对象,占据空间大、大小不固定。如果存储在栈中，将会影响程序运行的性能；引用数据类型在栈中存储了指针，该指针指向堆中该实体的起始地址。当解释器寻找引用值时，会首先检索其在栈中的地址，取得地址后从堆中获得实体

 19、没有块级作用域，只有作用域(执行环境)(ES5)
  ``` js
  if(true){
    var color = 'red'
  }
  alert(color) // 正确
  for(var i=0; i< 10; i++){var m = 100}
  alert(i, m) // 正确
  ```

20、 IE、opera可以主动触发垃圾回收  ！！！不建议

  IE： window.CollectGarbage()
  
  Opera: window.opera.collect()

21、new Array()

使用`new Array()` 与 `Array()`效果是一样的。

`new Array(arg)` 参数长度为1时要注意：
```js
  new Array(20) // []长度为20，用undefined填充
  new Array('20') // [ '20' ]  长度为1，内容为字符串20
  [,,,]  // IE8及其以下为 [undefined, undefined, undefined, undefined]
          // 其余浏览器为 [undefined, undefined, undefined ]
          // 也就是说最后一位尽量不要留逗号
```

22、数组方法

  - 改变原数组：reverse、sort、
  - 产生新数组：slice、concat、splice
  - 其余方法：indexOf、lastIndexOf、reduce、reduceRight、every、filter、forEach、map、some、

23、Array.prototype.indexOf()

 可以接受第二个参数表示查找的开始位置 arr.indexOf(searchElement[, fromIndex])

24、Date()

```js
// 两者返回的是距离1970、1、1零时的毫秒数
Date.parse() // 需要看浏览器心情，差异性很大
Date.UTC(2005, 4, 1, 1, 59, 59)  // 年、月、日、时、分、秒    月0-11 省略日默认1， 其余省略默认0

Date.now() // 静态方法获取当前时间毫秒数
+ new Date() // 在不支持的浏览器上兼容

```

25、 函数复制

函数复制(赋值)是真的复制。看例子

```js
// example 1
function a(){ console.log('a') }
b = a;
a();   // ok
b();   // ok
a = null;
a();   //error
b()    // ok
// example 2
var a = function(){ console.log('a') }
b = a;
a();   // ok
b();   // ok
a = null;
a();   //error
b()    // ok
```

26、arguments.callee  ！！！ 不建议使用，不允许使用

一个神奇的方法，可以获得参数的调用函数自身，在递归中可以使用，方便获取自身。但是！！！

警告：在严格模式下，第5版 ECMAScript (ES5) 禁止使用 arguments.callee()。当一个函数必须调用自身的时候, 避免使用 arguments.callee(), 通过要么给函数表达式一个名字,要么使用一个函数声明.

27、Function.caller   ！！！ 不建议使用，不允许使用

caller是javascript函数的一个属性，它指向调用当前函数的函数，如果函数是在全局范围内调用的话，那么caller的值为null。

28、arguments.caller  ！！！ 不建议使用，不允许使用

arguments.caller 这是我们遇到的第二个caller，没啥用，在严格模式下无法访问，非严格模式下值也为undefined，而且貌似被废弃了

29、String.prototype.slice、String.prototype.substring.prototype.substr

只传入一个参数，且为正数，则结果一致，表示从某一位开始截取到最后

传入两个正参数，slice、与 substring结果一致，表示从第几位到第几位， substr表示从第几位开始截取几位长度的字符串

slice中负数表示倒着数就可以了

substr只能识别第一位为负数，第二位出现负数会将第二位转为0

substring会将所有负数转为0
```js
var a = 'hello world';
a.slice(4); //"o world"
a.slice(4,7);  //"o w"
a.slice(-7, -4); //"o w"

a.substring(4); //"o world"
a.substring(4,7);  //"o w"
a.substring(-4);  //"hello world"
a.substring(-4, -7);  //""

a.substr(4);  //"o world"
a.substr(4, 1);  //"o"
a.substr(-1); //"d"
a.substr(-1, 0); //""
```

30、String.prototype.replace

如果参数1传入的是字符串，那么只会替换一次，传入带有g符号的正则才会全局替换

31、String.prototype.localCompare

`str.localCompare(str2)`按字符一个一个比较两者之间在字符标的前后顺序，如果str在str2参数之前返回-1， 如果在其之后返回1， 相等返回0

32、String.fromCharCode()

所做之事与charCodeAt()刚好相反

`String.fromCharCode(104, 101, 108, 108, 111)` --> 'hello'

33、Math
```js
Math.ceil() // 向上取舍
Math.floor() // 向下取舍
Math.round() // 标准取舍
```

34、Object.key 数据属性、访问器属性

对象的每一个属性(property)都包含数据属性和访问器属性。
数据属性包含`[[Configurable]] [[Enumerable]]  [[Writable]] [[Value]]`
访问器属性包含`[[Configurable]] [[Enumerable]]  [[Get]] [[Set]]`

可以看出两者都包含`[[Configurable]] [[Enumerable]]`属性。 那什么时候会用到呢
```js
let book = {};
Object.defineProperty(book, 'author', {});
// 在这里第三个属性就是描述属性， 但是只能传入数据属性或者访问器属性，不能串了，否则报错
```

35、Object.getPrototypeOf(obj)

获取一个对象的原型， 也是chrome私有实现的`obj.__proto__`的规范实现

36、for in、Object.keys()、Object.getWonPropertyNames()

- for..in.. 是对某个对象的可枚举属性进行遍历，不保证属性出现顺序，不保证属性属于对象本身还是原型链
- Object.keys() 是某个对象的所有可枚举属性名组成的数组，数组内容与for..in..一致，甚至顺序与for..in..一致。
- Object.getWonPropertyNames() 同样返回一个数组，但是返回的是对象自身的所有属性名组成的数组

37、Object.create(proto, properties)

Object.create(proto) 作为一种常见的实现继承的方式已经不难理解，但它还接受第二个参数，并且第二个参数与Object.definieProperties()的第二个参数格式一致，将作为新对象的熟悉。

38、location 

location即属于window对象又属于document，并且两者是对同一对象的引用 `window.location === document.location`。
```js
location.assign(new_url); //location.href、window.location赋值也是触发此函数
location.replace(new_url);
location.reload();
```

39、检测浏览器插件

浏览器插件在很多时候是很多网站工作的必备元素，例如很多网站需要flash， 可以提前检测flash，如果不存在就让提示用户去下载。

在标准浏览器中呢
```js
navigator.plugins  // 返回一个数组，包含当前开启的所有插件
plugin = {
  description : "Enables Widevine licenses for playback of HTML audio/video content. (version: 1.4.8.1000)",
  filename : "widevinecdmadapter.plugin",
  length : 1,
  name : "Widevine Content Decryption Module"
}
```
IE
```js
try{
  new ActiveXObject(COMname)
  return true;
}catch(e){
  return false;
}
```

40、HTMLCollection

`HTMLCollection[0]、HTMLCollection.item(0)、HTMLCollection.namedItem('name')`
当我们通过document.getElementsByXX获取到的是一个HTMLCollection类数组对象，它可以直接通过下标[0]来获取第一个元素，也可以通过.item(0)来获取第一个元素，同时还可以HTMLCollection.namedItem('name')来根据name属性获取对应元素

41、NodeList、NamedNodeMap、HTMLCollection实时动态性

这三个神奇的亲属关系的元素类型有一个神奇的特点。那就是特们总是保持最新的页面元素。换句话讲什么意思呢，代码表示
```js
var a = document.getElementsByTagName('div'); // a.length = 249
document.body.append(document.createElement('div'));
a.length // 250
```
看到没，明明对a(HTMLCollection)没有操作,但是a的内容却根据页面的变化而实际变化。 所以对这三者的操作要注意。

- NodeList由`Node.childNodes 和 document.querySelectorAll`产生
- HTMLCollection由`Element.getElementsBy`产生
(JS高级编程10.2.4标出，但实测chrome下document.querySelectorAll产生的NodeList并不具有像getElementsByTagName获取到的HTMLCollection那样具有实时动态性)
(但是实测a = document.body.childNodes产生的NodeList具有实时动态性)

42、compatMode 渲染模式

渲染模式一般分为标准模式和混杂模式。 document.compatMode可以获取当前的渲染模式，值为 CSS1Compat 或者 BackCompate。 那种是标准应该非常容易分辨喽

43、Element.insertAdjacentElement( position, element);

方法将一个给定的元素节点插入到相对于被调用的元素的给定的一个位置。
这个方法的兼容性简直爆炸，can.i.use显示基本全部兼容(firefox47以前不兼容)，但好像平时很少使用它。

```js
/*
'beforebegin': 在该元素本身的前面. 作为同级元素
'afterbegin':只在该元素当中, 在该元素第一个子孩子前面.
'beforeend':只在该元素当中, 在该元素最后一个子孩子后面.
'afterend': 在该元素本身的后面. 作为同级元素
*/
document.body.('afterbegin', document.createElement('div'))
```

44、contains

如果判断一个元素是不是某个元素的子元素，可以直接通过`parentElem.contains(childElem)`的结果true or false来判断，而且兼容性还是很好的，除了在非常低的safari下(但是apple的升级率很高啊，哪还有低版本safari)

45、scrollIntoView、scrollIntoViewIfNeeded

两者都是作用在元素上的， 
- elem.scrollIntoViewIfNeeded会判断当前元素是不是在显示区域内，如果不在就滚动页面让其显示，
- elem.scrollIntoView会滚动页面让元素显示在最上方
scrollIntoView兼容性不错，scrollIntoViewIfNeeded则惨不忍睹

46、DOM2增加了xhtml命名空间

html本身是不支持命名空间的，但由于xml支持，所以xhtml也支持了命名空间。

与命名空间相关的方法在方法名字中都有`NS`或者以`NS`结尾。

47、document.defaultView、document.parentWindow

使用document.defaultView是获取window的一个途径，有时是在elem.ownerDocument.defaultView, 但此方法在IE9以下不被支持，替代方法 document.parentWindow

48、访问iframe内容

iframeElem.contentDocument || iframeElem.contentWindow.document

49、style操作

`window.getComputedStyle(elem)` 返回的样式是一个实时的 CSSStyleDeclaration 对象，当元素的样式更改时，它会自动更新本身;
```js
CSSStyleDeclaration.cssText;
CSSStyleDeclaration.length;
CSSStyleDeclaration.parentRule;
CSSStyleDeclaration.getPropertyPriority(propertyName);
CSSStyleDeclaration.getPropertyValue(propertyName);
CSSStyleDeclaration.setProperty(propertyName, value, priority)
CSSStyleDeclaration.removeProperty(propertyName)
```

50、offset、client
  - 偏移量 offsetHeight、offsetWidth不包含margin、包含border、padding
  - 客户区 clientHeight、clientWidth只包含padding

51、document.createNodeIterator(elem)、document.createTreeWalker(elem)

 - createNodeIterator：返回Iterator遍历对象，可以用来nextNode()一直循环便利elem的所有元素
 - createTreeWalker: 返回一个方法更佳丰富的遍历对象

52、document.createRange()

范围，可以选择文档的一个区域，IE单独实现的是document.createTextRange()。 range表示包含节点和部分文本节点的文档片段。  一般在处理选择文本是常用到。