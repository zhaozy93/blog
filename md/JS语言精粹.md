# JS语言精粹
只记录重点

## 第四章 函数
- 调用一个函数会暂停当前函数的执行，传递控制权和参数给新函数。
- 每个函数会接受两个附加参数: this、arguments
- 如果函数调用时在前面加上new, 且返回值不是一个对象，则返回this(该新对象)
- 尾递归是一种在函数的最后执行递归调用语句的特殊形式的递归。[ruanyifeng尾调用解释](http://es6.ruanyifeng.com/#docs/function#尾调用优化)

## 第六章 数组
- 设置更大的length不会给数组分配更多的空间
- 设置小的length将导致所有下标大于或等于新length的属性被删除

## 第七章 正则表达式
- 正则表达式分组
  - 捕获型: ()形式， 可以再在正则表达式中使用\1、 \2、 \3等形式来表示捕获的内容
  - 非捕获型: (?:) 表示当前的这个分组不需要被记录，带来微弱的性能优势，并且不会干扰捕获型分组的编号
  - 向前正向匹配: (?=) 表示在匹配这个组的内容后，文本会倒回它开始的地方，实际上并不匹配任何内容。 不好的特性
  - 向前负向匹配: (?!) 表示在匹配失败时才继续向前进行匹配。 不好的特性

## 第八章 方法
- Array
  - array.join() 在IE8之前使用join方法连接大量字符串具有性能优势，但之后浏览器都对字符串的+运算做了优化，更推荐使用+来做性能追求。
  - array.sort() 是默认对字符串进行排序，因此如果元素不是字符串会先转换 造成[11, 12, 13, 2, 22, 3, 4]的问题
- Number

toExponential、toFixed、toPrecision的结果都是字符串哦
toFixed、toPrecision区别在于toPrecision对指数形式的也有效，但toFixed对指数形式数字无效

  - number.toExponential(fractionDigits) 把number转换成一个指数形式的字符串， fractionDigits控制其小数点后的数字位数
  - number.toFixed(fractionDigits)  把number转换成一个十进制数形式的字符串， fractionDigits控制其小数点后的数字位数
  - number.toPrecision(fractionDigits)  把number转换成一个十进制数形式的字符串， fractionDigits控制其小数点后的数字位数
- Regex
  - 记得regex的g模式有lastIndex属性哦
- String
  - string.localCompare
  - string.search vs string.indexOf: search可以接受正则表达式，indexOf只接受字符串

## 附录A 糟粕
- 数字采用标准 IEEE 754 浮点数的计算有bug
- JS没有真正的数组，但JS的数组非常好用，无需设置维度， 也没有越界的错误。 但性能和真正的数组比可能比较糟糕了