# 高性能JavaScript

## 加载和执行
1. 合并文件 减少http的开销
2. script放在最后 减少对页面渲染的阻塞
3. 使用延迟脚本 如async、defer属性
    - async 加载完立即执行
    - defer 页面构建完成后才会执行 只对带有src属性的标签有作用
4. 无阻塞脚本 即在window触发load事件后才执行

## 数据存取
1. 字面量和局部变量  比 数组和对象 读取快很多
2. 作用域链访问不宜过长。  可以将作用域链过长的变量暂时存为局部变量
3. with不推荐使用 并且它会新建一个对象作为作用域并且推入整个作用域链的顶部。
4. 尽量不使用动态作用域: eval、with等， 浏览器有静态分析

## DOM编程
1. JS与DOM是两个独立的引擎，因此天生就慢
2. clone DOM可能比createElement快一点
3. HTML集合是一种假定时时态，即当底层文档对象更新时，它也会自动更新。
4. 遍历数组比遍历HTML集合快很多，但是将HTML转为数组也会有一步损耗
5. 低版本IE中nextElementSibling比childNode快很多
6. 使用直接获得元素的API而不是哪些不区分元素节点、注释和文本节点的API
7. 接6 如nextElementSibling替代nextSibling
8. CSS选择器返回的并不是HTML集合，不是假定时时态，会快很多。querySelectorAll
9. DOM树表示页面结构、渲染树表示DOM节点如何显示
10. DOM树种每一个需要现实的节点在渲染树中至少存在一个对应的节点。
11. 当DOM树中的元素发生了几何属性更改会导致渲染树部分失效，然后进行重排reflow。 重排之后，浏览器会对受到影响的部分进行重绘repaint
12. 重排
    - 添加删除元素
    - 元素位置改变
    - 元素尺寸改变
    - 内容更改 文字或图片等
    - 页面渲染器初始化
    - 浏览器尺寸改变
13. 浏览器的重排和重绘以队列形式执行。但某些操作会导致其立即执行，因为这些命令都需要返回最新的布局信息
    - offsetTop、 offsetLeft、offsetRight、offsetBottom
    - scrollTop、 scrollLeft、scrollRight、scrollBottom
    - clientTop、 clientLeft、clientRight、clientBottom
    - getComputedStyle(）、 currentStyle()
14. 多个css更改时可以考虑用class的更改来代替
15. 使用文档片段createDocumentFragment来进行多个dom的修改
16. 动画元素可以脱离文档流 绝对定位
17. 使用事件委托代理

## 算法和流程控制
1. for in的执行速度明显慢于其他三个
2. 条件数量越大 越倾向于switch 而不是if else
3. 优化if else的条件设置以及条件的先后位置
4. 有时候简单的查找表比switch和ifelse都要快
5. 查找表r=[r0, r1, r2, r3]。 直接r[4]
6. 尝试优化递归 如 尾递归

## 字符串和正则表达式
1. 正则有可能造成回溯失控，要注意正则的编写
2. 正则并不意味着高效，具体情况具体分析 不要一味追求正则
3. 正则的工作原理是回溯
4. 正则回溯： 正则匹配字符串时，从左到右，看能否找到匹配项。 在遇到量词和分支时，需要决策下一步如何处理。 同时有必要的话会记录其他选择，以备返回时使用。 如果匹配时失败，会回到上一个决策点，再进行匹配。直到匹配成功或者全部匹配完毕失败。
5. 提高正则匹配效率
    - 关注如何让匹配更快的失败
    - 以简单、必需的字开始
    - 减少分支
    - 使用非捕获组
    - 只捕获感兴趣的文本

## AJAX
- GET产生一个TCP数据包；POST产生两个TCP数据包。 但还是更具需要使用 毕竟现在网络哪有那么差
- get、post区别[url](http://blog.csdn.net/happy_xiahuixiax/article/details/72859762)
- XMLHttpRequest 有一个3的状态 有的浏览器支持，对于文件流有很大的改善帮助

## 编程实践
- 不使用String作为定时器的参数
- 避免重复判断，如判断是否是IE 判断后缓存一个变量 不要老是去判断
- 延迟加载某些函数 在函数内部进行判断，然后更改函数为对应的版本
- 尝试使用位操作
- 尽可能使用原生方法

## 构建并部署高性能JS应用
1. 合并多个JS文件
2. 设置缓存