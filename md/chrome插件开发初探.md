# Chrome插件开发初探

> 凡是可以用 JavaScript 来写的应用，最终都会用 JavaScript 来写。 
>
> ——Atwood定律（Jeff Atwood在2007年提出）

如今chrome基本是所有程序员的标配，每个人或多或少的都会使用到chrome的插件吧。 比如我常用的就有前5个。
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/chromeExtensionmd01.jpeg)
- [stash](https://chrome.google.com/webstore/detail/stash/bnhjedgfogckebfhnlicnkbdjlmpibck) 用于关闭并暂存当前所有tab的， 很实用， 方便在不同项目之间切换(毕竟API搬运工，每个项目要搬不同的API)
- React、Redux 前端开发不解释
- [Open In IINA](https://chrome.google.com/webstore/detail/open-in-iina/pdnojahnhpgmdhjdhgphgdcecehkbhfo) 又是一个很好用的视频播放插件，当页面有flash时可以点击这个 调用本地的播放器 而不让flash造成mac发热蒸鸡蛋的现象
- [掘金](https://chrome.google.com/webstore/detail/%E6%8E%98%E9%87%91/lecdifefmmfjnjjinhaennhdlmcaeeeb) 掘金的一个首页插件，每次打开新tab时就会自动加载掘金网站的热门帖子， 掘金PM想法不错

- 还有最后一个 [Time recorder(时间记录器)](https://github.com/zhaozy93/blog/blob/master/chromExtension/timeRecorder.crx])，虽然没有发布在chrome商店，但是你可以自己去下载啊 -_-。

先上两张图片，由于家里电脑插件也是刚安装，要本还不是特别多。
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/chromeExtensionmd02.jpeg)
![image](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/chromeExtensionmd03.jpeg)

## 想法
每天电脑切换次数、打开时间最长的就是chrome > iterm > vscode， 因此就想看看自己每天都浏览什么网站。 历史记录虽然有但是不直观啊， 作为懒人怎么可能去翻历史记录，万一翻到不该看得怎么办。

## 于是有了Chrome插件之旅

搜了搜之后发现360做了一件良心之事，把chrome extension[开发文档](http://open.chrome.360.cn/extension_dev/overview.html)翻译了！！！ 我的天，尽管网页丑的要死，但还是不错的。但还是看着怪怪的，于是乖乖回去看原生的了。 [Chrome extension develop docs](https://developer.chrome.com/extensions)。 chrome团队文档写的很好，还配有示例。

## Chrome extension编写
chrome extension就是按照特定要求编写的网站，依旧有HTML、CSS、JS + Chrome APIs组成， Chrome继续利用V8引擎来运行我们编写的插件。 举个例子就是 chrome启动之后就在后台运行了我们编写的Web程序。 当点击icon时把我们需要的页面展示给用户就好了。

在插件中最重要的文件(或者称为入口文件, 本人取得名字 -_-) 就是manifest.json， 他就像我们的package.json一样， 一个配置文件，格式要求严格， 内部字段都是预先定义的。 所有字段与预设值都可以[查到](https://developer.chrome.com/extensions/manifest)

来看一下time recorder的json文件来作为示例
```json
{
  "manifest_version": 2,
  // 对整个插件最基本的描述
  "name": "Time recorder",
  "description": "This extension allows the user to record counts and time of each page you opened .",
  "version": "1.0",
  // 这就是我们浏览器右上角的那个小图标啦
  // 用户点击时就打开popup.html
  // 鼠标悬停显示 See Detail!
  "browser_action": {
    "default_icon": "icon.png",
    "default_popup": "popup.html",
    "default_title": "See Detail!"
  },

  // 我们的程序是否需要后台运行呢？ 必然是的
  // 是否需要 persistent 持久运行呢？ 是的
  // 后台运行的脚本有吗？ 有啊background.js
  "background": {
    "persistent": true,
    "scripts": ["background.js"]
  },
  // 这个主要是图标描述文件
  "icons": { 
    "16": "icon.png",
    "48": "icon.png",
   "128": "icon.png" 
  },

  // 插件需要的权限
  // 比如 需要存储、需要tab的信息、 需要后台运行等等等
  "permissions": [
    "activeTab",
    "storage",
    "tabs",
    "background"
  ]
}
```

现在我们就大概明白了chrome extension整体运行的逻辑了。 也能大概想明白time recorder插件的执行逻辑， chrome启动就执行插件，每次打开一个新网址就记录，关闭的时候记录时间，同时记录每个网址打开的次数。 但是记录具体的网址不太现实，也没啥作用，大部分网址可能只打开一次就算了，因此实际记录的是域名。

还特地编写了从url-->host的处理逻辑，可能有时候还是不太严谨，但目前够用了。
```js
function getHost(url) {
  if(!url){ return null; }
  var host = url;
  var regex = /.*\:\/\/([^\/]*).*/;
  var match = url.match(regex);
  if(typeof match !== "undefined" && null !== match) {
    host = match[1];
  }
  return host;
}
```

如何知道有新网址被打开呢， 主要是依赖于Chrome APIs， 这里面就有tab状态更新的事件，我们只需要在background.js脚本中监听这个事件，就能得知新tab打开，tab网址更换，tab被关闭等事件，与前端JS如出一辙，只是API变了一下而已。

经过重重处理，最后将数据存储在 Chrome.storage， 一个与localstorage非常类似的接口，但它是异步的。 

在用户点击浏览器右上角插件icon时，弹出popup.html文件，其实也是被作为一个新的tab打开的。 就像普通的绘图一样，把我们之前记录的数据绘制出来即可。

### Chrome APIS坑
chrome提供的大部分API都是异步的，向Chrome.storage这样的存储都是异步的，因此需要注意程序书写时代码的位置

## 不足
当前插件记录的时网页打开时间，但是我有个坏习惯，tab打开之后在tab不是特别多的时候经常不及时关闭，而且晚上下班浏览器还一直开着，那这些时间其实不应该被算进去。我们关心的应该是focused的网页，所以当前大逻辑还有待更改。

## 总结
其实我们能做事情还很多，有更多的领域等待去探索。