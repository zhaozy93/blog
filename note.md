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