# JS框架设计
首先声明 这本书的大部分内容可能对于一个有经验的FEer来说都是温故，但是可能也是目前市面上难得的少见的以书写前端框架为思路写书的。书内有很多知识点都是值得学习，借鉴的。 内容分别用jQuery、prototype、avalon等框架来做对比用于解释某些内容，虽然书籍跟不上时代了，内部很多东西都不再是日后工作和应用的重点，甚至伴随着浏览器的升级与标准化，很多兼容性内容也早已退出舞台，但是了解历史、学习他人思路是很重要的一个成长过程。

# 节点模块
1. innerHTML: 在ie某些read-only标签是不支持的如 Col Head HTML...   [stackOverflow](https://stackoverflow.com/questions/16234410/detecting-whether-innerhtml-is-readonly)

2. innerHTML: 在ie中直接使用某些标签会无效， 需要在前面加上一些东西 比如文字或者标签
    - innerHTML = '<meta content="IE=9"/>'  // 失败
    - innerHTML = 'X<meta content="IE=9"/>'    // 生效

3. Range对象 很神奇 [MDN](https://developer.mozilla.org/en-US/docs/Web/API/Document/createRange)

4. MutationObserver: window.MutationObserver可以用来监听dom元素的变化， 兼容性有待提高
    - 类似事件， 但是异步触发
    - 一个for循环插入100个节点 由于是异步 会进行一次数据封装，只触发一次MutationObserver
    - 但是会导致文本内容变得支离破碎，由原本1个文本节点变为N个，这样会导致MVVM框架对比扫描节点失败

5. registerElement: document.registerElement 可以注册自定义标签元素， 在定义的时候可以传输指定原型，在指定原型上有生命钩子函数，如createdCallback、attachedCallback、detachedCallback、attributeChangedCallback， 这也相当于是对dom变化的监听

6. - innerText 触发reflow  textContent不会
   - innerText返回值会格式化

7.  iframe
    - 隐藏边框 iframe.frameBorder = 0
    - 去掉滚动条 iframe.scrolling = 'no'
    - 获取window iframe.contentWindow
    - 获取document iframe.contentDocument