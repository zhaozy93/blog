# 各类基础知识点

## 进程与线程(浅显解释)

计算机核心是CPU，现代计算机都是多核CPU， 一个内核只能运行一个任务(进程)。如果任务(进程)过多就要排队等候。同时也就造成了进程与进程之间不共享内存。

一个进程则包含许多线程，线程们之间协助完成任务，那么他们就必须共享内存，但是资源(内存)的共享也是有限制的， 有些内存只能同时被一个线程访问，那么其余的就要排队(互斥锁)。还有一种资源(内存)可以被多人同时访问，那么超过数量限制的线程就也要排队等待(信号量)。

因此对于操作系统是多进程形式，允许多个任务同时运行。

对于程序是多线程，允许把单个任务分割成多个部分去完成。

通常把进程作为分配资源的基本单位，而把线程作为独立运行和独立调度的基本单位

---- 分割线

- 单进程单线程：一个人在一个桌子上吃菜。
- 单进程多线程：多个人在同一个桌子上一起吃菜。
- 多进程单线程：多个人每个人在自己的桌子上吃菜。

多线程的问题是多个人同时吃一道菜的时候容易发生争抢，例如两个人同时夹一个菜，一个人刚伸出筷子，结果伸到的时候已经被夹走菜了。。。此时就必须等一个人夹一口之后，在还给另外一个人夹菜，也就是说资源共享就会发生冲突争抢。

对于 Windows 系统来说，【开桌子】的开销很大(创建进程的时间开销很大)，因此 Windows 鼓励大家在一个桌子上吃菜。因此 Windows 多线程学习重点是要大量面对资源争抢与同步方面的问题。

对于 Linux 系统来说，【开桌子】的开销很小，因此 Linux 鼓励大家尽量每个人都开自己的桌子吃菜。这带来新的问题是：坐在两张不同的桌子上，说话不方便。因此，Linux 下的学习重点大家要学习进程间通讯的方法。

--- 分割线

一个任务主流执行方式
- 单线程串行执行 一个环节卡住了整体都卡住了， io与计算任务没区分开
- 多线程并行执行 开销主要在于创建线程以及线程执行期间上下文切换还有就是数据共享问题， 演变而来就是锁、状态同步等问题

reference: http://www.ruanyifeng.com/blog/2013/04/processes_and_threads.html

## JS的单线程

浏览器是单进程还是多进程呢？ 浏览器的实现方式不同。主要表现在多标签页时。

多进程浏览器：
- IE 10
- chrome
- opera 15或以上

多线程浏览器：
- IE 6
- oprea 12或以上
- Firefox

采用多进程的优点：
- 安全：现代系统都有进程的安全机制，单个进程有自己独立的内存空间
- 稳定：不会因为一个线程的崩溃导致整个应用的崩溃
- 性能

采用多进程的缺点：
- 内存占用大
- 进程间通讯的成本大
- 进程启动的开销大

JS异步机制
- JS的单线程是指一个浏览器进程中只有一个JS的执行线程
- 而异步机制是浏览器的两个或以上常驻线程共同完成的，例如异步请求是由两个常驻线程：JS执行线程和事件触发线程共同完成的，JS的执行线程发起异步请求（这时浏览器会开一条新的HTTP请求线程来执行请求，这时JS的任务已完成，继续执行线程队列中剩下的其他任务），然后在未来的某一时刻事件触发线程监视到之前的发起的HTTP请求已完成，它就会把完成事件插入到JS执行队列的尾部等待JS处理。

## 并发、并行与串行

- 并发是两个队列交替使用一台咖啡机，
- 并行是两个队列同时使用两台咖啡机，
- 串行，一个队列使用一台咖啡机，那么哪怕前面那个人便秘了去厕所呆半天，后面的人也只能死等着他回来才能去接咖啡，这效率无疑是最低的。
区别
- 并发的关键是你有处理多个任务的能力，不一定要同时。
- 并行的关键是你有同时处理多个任务的能力。
- 并发和并行都可以是很多个线程，就看这些线程能不能同时被（多个）cpu执行，如果可以就说明是并行，而并发是多个线程被（一个）cpu 轮流切换着执行。

## Node.js异步IO
其实异步IO的本质与 前面提到的  *JS的单线程* 标题下面*JS异步机制* 类似
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs01.jpeg)

## 阻塞IO与非阻塞IO、异步IO与同步IO

1、 异步IO* 与 *阻塞IO* 听起来很类似， 目的上也确实达到了我们前面图片描述的那样 并行执行IO的目标。 但本质是不一样的
*系统内核只有 *阻塞IO与非阻塞IO* 两种
  - 阻塞IO的意思就是应用程序要等待系统完成所有的IO操作调用才算结束，程序才能继续执行。 读取文件为例就是磁盘寻道-->读取数据-->复制数据到内存-->调用结束， 结果就是CPU等待IO，浪费CPU性能
  - 非阻塞IO为调用之后系统立即返回，但是不包含真实数据，只有一个文件描述符。CPU可以继续做其他事情

2、 在非阻塞IO模式下程序如何得知IO完毕了呢？
  
为了获取完整数据，CPU需要不断的去重复调用IO来确认完成状态， 叫做轮询(非常类似于前端的轮询)
- read。它是最原始、性能最低的， 通过一直重复调用I/O的状态来完成数据的读取。在真正数据返回之前CPU一直浪费在调动、等待上面，效率极低
- select。 对read的一种改进， 通过文件描述符上的事件状态来判断， 不在进行无意义的重复read，但是有个限制 对多对1024个文件描述符同时检查
- poll。 对select的改进， 突破了1024的限制，但是文件多时效率还是极低
- epoll。 Linux下效率最高的IO方式。在轮询过程中如果没有真正数据就会休眠，直到IO结束后利用完成事件将它唤醒。 真正做到了时间通知、执行回调的方式，而不是遍历查询，节省CPU性能。

![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs02.jpeg)

轮询技术确实保证了非阻塞IO获取数据完整的要求，但对于应用程序而言仍然还是一种同步，因为要么处于遍历文件描述符、要么处于休眠状态。实际效果还是不够好

3、 现实的异步IO(主程序单线程模式下)

利用系统多开一个或多个线程执行阻塞或非阻塞IO， 在IO完成后通过线程之间的通信来告诉主线程 IO完毕

## Node.js 异步IO
1、Node.js自身执行模型---事件循环

 通俗讲Node.js最外层是一个While(true)的循环， 循环内部被称为Tick。
 ![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs03.jpeg)
 注意：这里的退出仅仅是Tick的退出 也就是进入下一次while循环

2、观察者  Tick中向谁询问是否还有事件

在Node中事件主要来源是网络请求和IO读取，每一类事件都有对应的观察者：I/O观察者、网络请求I/O观察者。 同时实现了观察者对事件的分类

事件循环是一个典型的*生产者与观察者*模型。 
- 异步IO、网络请求是生产者：为node产生源源不断地事件
- 事件被传递到对应分类的观察者那里
- 事件循环不断去询问观察者

在Windows下，这个事件循环基于IOCP创建， 在*nix下循环基于多线程创建。

注: 观察者中 *idle观察者* 优先于 *IO观察者* 优先于 *check观察者* 

3、请求对象

下图是node中执行fs.open()经历的过程
- js调用函数
- 调用核心模块的c++库
- 內建模块根据平台区分调用对应的不同平台的具体实现
- 在uv\_fs\_open()的调用过程中，我们创建了一个FSReqWrap请求对象， js调用时的参数都被传入给了uv\_fs\_open方法。至此对象封装完毕。
- window平台下，调用QueueUserWorkItem()将这个FSReqWrap对象推入线程池中。至此JS的调用结束
- JS线程继续执行后续的操作、计算等等
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs04.jpeg)

4、回调函数

还是上面的例子，讲到过所有参数都传给了请求对象， 那么回调函数也赋给了FSReqWrap请求对象oncomplete_sym事件上面， 当文件打开后这个事件就被触发了。

IO线程在结束后会调用PostQueuedCompletionStatus(作用是通知IOCP状态有更改，并归还线程至线程池)通知IOCP 这里有一个事件结束了。

下一次Tick时候，会调用GetQueuedCompletionStatus()从观察者那里查询事件。 查到事件并且有回调函数那么就执行回调函数

![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs05.jpeg)

注意 这里有一点说明。
- 组装好IO请求对象并推入线程池这是IO第一部分
- 回调通知是IO的第二部分(包含执行回调函数)。
- 也就是说我们的回调函数不是由我们自己的js调用，所以很多时候try catch是无效的

## 非IO得异步

1、setTimeout、setInterval

两者它们的实现与IO异步类似，只是不需要I/O线程参与。 用setTimeout()或者者setInterval()创建的定时器会插入到定时器观察者内部的一个红黑树中。

之后每次执行Tick询问定时器观察者，如果超过时间了就形成一个事件，然后把回调函数执行。 

所以定时器并不是严格按照时间执行的

2、process.nextTick

如果只是想在下一个循环中执行一个回调的话 可以直接使用process.nextTick, 直接把任务插入到下一次事件循环，比setTimeout((),0 )效率高很多。

3、setImmediate

效果与process.nextTick是一致的，但是有细微差别。

process.next优先级更高，主要是因为process.nextTick属于idle观察者， setImmediate属于check观察者。

## 异步编程缺点与难点
- 异常处理： 上面也提到过，很多场合使用简单的try catch是不能捕获到cb里面的异常
- 回调嵌套： 前端也存在这个问题， 万年坑
- 阻塞代码： 一种语言竟然没有sleep这样的函数
- 多线程编程： 无法充分利用多核CPU性能

## 异步编程解决方案
- 事件发布／订阅模式
- Promise／Deferred模式
- 流程控制库

### 事件发布／订阅模式
把回调函数当作一种事件来处理。 有事件的生产者发布消息，由事件监听者订阅消息并执行对应的回调函数。

事件发布与订阅并不是严格意义的异步解决方案，但是它的完成是依赖于事件循环体系，所以也广泛应用于解决异步回调问题。
```js
// 订阅
emitter.on("event1", function (message) {
    console.log(message); 
});
// 发布
emitter.emit('event1', "I am message!");
```
事件除了普通的监听还有特殊的
```js
event.once("event1", callback);  // 只可能被执行一次
proxy.all("template", "data", "resources", callback);  // 多事件之间协同作战
```

### Promise／Deferred模式

Promise/A抽象定义
- Promise操作只有三种状态：未完成态、完成态、失败态
- Promise状态只能从未完成向完成或失败转换，且不能逆反。完成态和失败态不能互相转换
- Promise状态一旦改变，将不能再次被更改

events模块来模拟Promise实现
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs06.jpeg)

接着来看Deferred， Deferred其实是用来触发Promise的
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs07.jpeg)
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs08.jpeg)

Promise／Deferred之间的差别就在于Deferred主要是用于内部，维护异步模型的状态，而Promise主要用于外部，把then暴露给外界添加自定义逻辑。
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs09.jpeg)

### 流程控制库
这里没有一个统一的规范，各路神仙大显神通。
#### 尾触发与Next
非常适合中间件的使用
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs10.jpeg)
#### async
略过

## 异步并发控制
这里的并发不是表面线程、进程层面上的并发， 而是异步概念的并发
```js
for(let i =0; i< 100; i++){
  async();
}
```
尽管可以很简单的同时打开100个异步请求去，类似于读取磁盘文件，但是磁盘也是有过载保护的。 所以还是需要控制异步的数量

给出了一个包`bagpipe`， 这个包就是一个池子，每次像池子里面`push`，当出水口不满时，就持续增加水流(加大执行个数)，当出水口满了，多余的水则待在池子里等着。

## V8内存
一般服务器端的后台开发都不会有内存限制，但是由于Node是构建于V8之上，因此V8对内存的限制也就限制了Node，同时Node对JS对象的管理都是通过V8自己的方式来进行分配和管理的。(64位一般约为1.4G， 32位约为0.7G)

V8中所有JS对象都是通过堆来进行分配的，Node中提供了查看内存使用量的方法`process.memoryUsage()`，其中heapTotal是V8申请到的内存，heapUsed是使用的内存量。当然也可以`--max-old-space-size(单位为MB) 或 --max-new-space-size(单位为KB)来更改内存限制`

V8为何要做内存限制呢？

深层原因是V8的垃圾回收机制，官方测试1.5G的堆内存垃圾回收为例，V8做一次小的垃圾回收需要50ms，做一次非增量式的垃圾回收甚至需要1秒以上，这段时间是垃圾回收引起JS线程暂停的时间，应用性能和响应能力都直线下降，这是前端或者后端都不能接受的。

V8内存主要分为新生代和老生代两代。 也就是前面提到的`new-space`与`old-space`。新生代内存中对象为存活时间较短的对象，老生代内存中对象为存活时间较长或者常驻内存的对象。
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs11.jpeg)
老生代的大小：在Node源码中可以看到，在64位系统下为1400MB，在32位系统下为700MB，也就是前面说的约。

新生代内存： 由两个reserved\_semispace\_size\_构成。reserved\_semispace\_size\_在64位系统32位系统上分别为16MB和8MB， 然后double一下。

## V8垃圾回收机制
没有一种垃圾回收机制可以完美应对各种环境，因此V8其实采用了多种解决方案，将内存分为新老生代就是其中一种办法。

### 新生代
刚才提到新生代其实是由两个reserved\_semispace\_size\_构成，新生代垃圾回收主要采用的是Scavenge思想的Cheney实现算法。 将堆内存一分为二，称为semispace，两个中一个处于活动状态，一个处于闲置状态。处于使用状态的称为From，闲置状态的称为To。分配对象时分配给From，垃圾回收时将活动的对象复制到To，垃圾回收结束后将From和To对调，然后释放掉To(对调前的From)。

Scavenge思想是一种典型的空间换时间的思路，因此无法应用在大规模的垃圾回收中。但是比较适合新生代内存，空间小，生命周期短的特性。

在由From向To转换的过程中，需要进行特定的检查，经过一定检查的对象需要转移到老生代内存中，完成晋升。
- 检查一： 是否经历过Scavenge检测
- 检查二： To空间内存占比超过一定比例(比例为25%， 因为将To转换为From之后还要负责接收新的对象内存开销)
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs12.jpeg)

### 老生代
Mark-Sweep & Mark-Compact共同完成老生代的垃圾回收。 老生代特点是存活对象比例大且存活时间长。

Mark-Sweep分为标记与清除阶段。先将所有存活对象进行标记，其后再将所有未标记对象清除。 但是造成的后果就是，清除后得到的都是一堆不连续的内存地址，如果下次需要分配比较大的空间可能无法完成，然后提前触发一次不必要的垃圾回收。

![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs13.jpeg)

为了解决清除后内存地址不连贯的问题，提出了Mark-Compact思路，在清除阶段将所有存活对象向内存一端移动，之后清除掉存活对象内存边界外的即可。

![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs14.jpeg)

因为需要对象的移动，因此Mark-Compact并不是完全优于Mark-Sweep的，两者在工作中也是互补的。 大部分时间进行Mark-Sweep逻辑，当空间不足时使用Mark-Compact。

![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs15.jpeg)

## Incremental Marking
垃圾回收机制一般都需要将程序逻辑执行暂时暂停，以确保JavaScript应用程序与垃圾回收器一致，这种行为一般称为全停顿(stop-the-world)。新生代由于体积小，即使全停顿影响也不大，但是老生代由于体积较大，全停顿带来的影响还是比较明显的。因此V8推出了增量标记(Incremental Marking)的方式来将原本一次标记清除拆分为多次执行，没做完一小步就让JS执行一会，再标记，再让JS执行，知道完成所有标记。
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs16.jpeg)

因此要提高整体程序执行效率就要减少老生代的垃圾回收次数。

## V8内存--堆外内存
在调用process.momoryUsage()有一个属性为rss，全称为resident set size，进程的常驻内存部分。heapTotal是堆申请到的总内存。 因此与之对应的还有堆外内存。 rss减去heapTotal就是堆外内存所占用的大小。 _`V8对堆外内存没有限制，因此可以使用堆外内存来突破普通内存的限制`_。直接的方法就是Buffer对象，在Node中为了满足流处理，简单的字符串操作不能满足性能需求，因此Buffer对象由Node进行操作，不由V8管理。


## Buffer
### 字节（Byte）
字节是通过网络传输信息（或在硬盘或内存中存储信息）的单位。字节是计算机信息技术用于计量存储容量和传输容量的一种计量单位，1个字节等于8位二进制，它是一个8位的二进制数，是一个很具体的存储空间。

### Buffer结构
Buffer模块是JS与C++结合的模块，性能部分由C++实现，非性能部分由JS实现。 同时Buffer属于堆外内存。
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs17.jpeg)

### Buffer对象介绍
buffer对象类似于数组，每个元素是由两位16进制组成，即0-255数值。

但是下面的例子可以看得出来，不同编码格式、不同字符占用的buffer的元素个数不一样。 utf-8编码中汉子占3个位置。
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs18.jpeg)

### Buffer内存
Buffer性能部分由C++实现，当然内存部分肯定由C++实现。 是由C++申请内存，JS分配内存的策略。 但是C++频繁申请内存也会造成较大系统压力。

Node采用slab(动态内存管理机制)分配机制，C++每次申请一个特定大小(大小Buffer对象的区分线)的内存，如果Buffer对象小于这个这个界限，那么创建一个slab单元(内存大小就是区分值)，并且分配足够的内存给Buffer对象，同时Buffer对象的parent指向这个slab单元，并且记录下从slab哪里开始到哪里结束。 当再次新建小对象时检测当前的slab空间是否足够，不足的话再申请一个slab空间，同时上一次slab剩余的空间就被浪费了。 如果申请一个足够大的Buffer对象，那么直接申请一个足够大的slab对象，并且这个slab对象被独占。

### Buffer与字符串转换
`new Buffer(str, [encoding]);`可以实现字符串向buffer转换，但是只能指定一种编码格式。

`buffer.write(str, [offset], [length], [enconding])` 可以实现同一buffer指定不同编码格式。 但是转换为字符串时需要特别注意。

`buf.toString([encoding], [start], [end])`buffer转字符串。

### Buffer读取
前面提到过在utf8编码下汉字占3个字符，但是如果我们在读取的时候恰好最后2个字符，那么正好读出来就会是乱码。 
- readable.setEncoding(encoding)
- chunks = []; chunks.push(chunk); size += chunk.length;  ==> Buffer.concat(chunks, size)


##网络编程
OSI模型7层构成
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs19.jpeg)

### http
利用curl -v来看 http组成
- tcp三次握手
- 请求报头
- 响应内容： 响应头+响应体
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs20.jpeg)

### http与tcp
http协议是基于tcp协议的，node中http服务继承TCP服务器(net模块)，由于采用事件驱动(请求视为一个事件),不会创建额外的线程或进程，因此保持低内存，便于高并发。但是开启keepalive之后，一个TCP回话可以用于多次请求和响应。TCP以connection为单位，http以request为单位。 http模块是connection到request的封装。
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs21.jpeg)

### websocket握手
请求头中包含Sec-WebSocket-Key字段，服务器将其与”258EAFA5-E914-47DA-95CA-C5AB0DC85B11″这个字符串进行拼接，然后对拼接后的字符串进行sha-1运算，再进行base64编码，最后以“Sec-WebSocket-Accept”字段形式返回给客户端。 客户端验证通过则握手成功。

### 网络传输加密与解密
在客户端与服务器建立安全传输之前，两者需要互换公钥。
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/nodejs22.jpeg)

