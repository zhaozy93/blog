# chapter6 浏览器测试

浏览器兼容永远是前端说不完的痛，尤其对于2C的客户，需要面对各种浏览器，尤其是低版本的IE，尽管现在大家基本已经放弃了IE67，但IE8仍然不是一个标准的W3C浏览器。仍然一堆兼容性需要处理，例如最基本的document.attachEvent

## 测试策略
浏览器功能性检测大概可以分为两类：
- 检测navigator.userAgent，用户代理检测法
- 检测浏览器的功能特性，即功能特性检测法

1. 用户代理检测方法主要依靠依据的是大家的工作经验，比如大家知道IE绑定事件是document.attachEvent，因此我们可以这样写

```
if ( IE ){
    document.attachEvent(.....)
} else {
    document.addEventListener
}
```
但是这样的兼容性我们可能会忽略或者不全面。 比如一个小众的浏览器，我们可能不会那么准确的通过userAgent来判断他是不是IE还是W3C标准浏览器。

关于userAgent，这里有一片不错的[文章](http://www.cnblogs.com/ifantastic/p/3481231.html)介绍它的乱七八糟的混乱，但是userAgent仍是后台来简单区分用户浏览器的一个不错的选择。

2. 功能特征检测法则依据浏览器是否支持某项特定的功能特性，来决定程序的执行分支。还是刚刚的例子，这一次的判断条件不再是IE。
```
if (document.attachEvent ){
    document.attachEvent(.....)
} else if ( document.addEventListener ){
    document.addEventListener(...)
} else {
    return new Error()
}
```
这一次我们不在关心页面运行在哪个浏览器，也不关心浏览器版本，我们知道有两个方法都可以完成我们的任务， 目标浏览器支持哪个就运行哪个分支，这样确保程序可以正常运行，如果发现两个方法都不被目标浏览器支持，那么不好意思了，目前没有可以替代的第三种方案，那只能报错了。。。

综合上面两种方案，使用功能特性检测法可能更符合我们的希望目标，毕竟有一点很重要，我们不可能穷尽所有厂商各版本的浏览器特征，所以我们可以只检测我们用到的特性方法。jQuery也是采用的第二种功能特性检测发来实现的浏览器测试部分。

## 源码
这一次直接将全部代码一次性作为整体，不再分割成多块。

```
jQuery.support = (function() {
  // 定义一堆变量以后使用， 新建一个div元素作为容器，同时获取html元素的引用
  var support,
    all,
    a,
    select,
    opt,
    input,
    marginDiv,
    fragment,
    tds,
    events,
    eventName,
    i,
    isSupported,
    div = document.createElement( "div" ),
    documentElement = document.documentElement;

  // 初步检测
  div.setAttribute("className", "t");
  div.innerHTML = "   <link/><table></table><a href='/a' style='top:1px;float:left;opacity:.55;'>a</a><input type='checkbox'/>";

  all = div.getElementsByTagName( "*" );
  a = div.getElementsByTagName( "a" )[ 0 ];

  // 如果连innerHTML、getElementsByTagName、getElementsByTagName任意一个都失败，那么support就没意义了
  if ( !all || !all.length || !a ) {
    return {};
  }

  // 新建select和option以及获取input引用
  select = document.createElement( "select" );
  opt = select.appendChild( document.createElement("option") );
  input = div.getElementsByTagName( "input" )[ 0 ];

  support = {
    // 刚才innerHTML那段字符串是以空格开头的， 测试浏览器是否会保留innerHTML的前导空格
    // 如果保留空格符的话 那么第一个子元素将是文本node， nodetype为3， 普通dom元素nodetype是1
    // 保留签到空格为true 
    leadingWhitespace: ( div.firstChild.nodeType === 3 ),

    // 测试浏览器会不会自动为空table插入tbody元素， 通过测试刚才那段innerHTML里有没有自动生成tbody即可
    // 不自动插入为true
    tbody: !div.getElementsByTagName("tbody").length,

    // 测试能不能正确的序列化link标签
    // 低版本ie如果设置innerHTML='<link/>'会直接把link吃掉，
    // 解决办法是加一个'div<div><link/></div>'，然后再一层层取出来link元素  不知道为什么。
    // 能保留link则为true
    htmlSerialize: !!div.getElementsByTagName("link").length,

    // 测试能不能通过getAttribute获取内连style信息
		// ie678需要使用a.style.cssText来获取内连style
    // getAttribute成功则为true
    style: /top/.test( a.getAttribute("style") ),

    // 测试a标签的URL会不会被改为全路径
    // 如果不更改则为true，
    hrefNormalized: ( a.getAttribute("href") === "/a" ),

    // Make sure that element opacity exists
    // (IE uses filter instead)
    // Use a regex to work around a WebKit issue. See #5145
    // 测试浏览器支持opacity属性与否
    // 在ie8中实测会返回.55，这样测试就不会通过，因为支持opacity的浏览器会解析.55为0.55
    // 但需要注意的是正则里面的.不是小数点的点，而是匹配通配符
    opacity: /^0.55/.test( a.style.opacity ),

    // 测试是否支持cssFloat属性，如果不支持 后续需要使用styleFloat来修正
    // 支持cssFloat则为true
    cssFloat: !!a.style.cssFloat,

    // 测试默认的<input type='checkbox'/>的value值
    // 不过实测safari、chrome、ie8全为on，可能是老版本的safari会是默认空字符串
    checkOn: ( input.value === "on" ),

    //   select = document.createElement( "select" );
    // opt = select.appendChild( document.createElement("option") );
    // 测试option默认是否选中，ie8是false， safari、chrome、firefox为true
    optSelected: opt.selected,

    // div.setAttribute("className", "t");
    // 测试能否在setAttribute需要传入的是dom属性名还是html属性名
    // 同样的方法还有getAttribute和removeAttribute
    // class这个属性是html属性， className是dom属性
    // 需要传入html属性是true， dom属性是false
    getSetAttribute: div.className !== "t",

    // enctype 属性规定在发送到服务器之前应该如何对表单数据进行编码
    // 直接测试form元素有没有这个属性即可
    enctype: !!document.createElement("form").enctype,

    // 测试浏览器能否正确的复制html5的元素
    // 以nav标签为测试例子
    // 这里如果outerHTML为undefined也是认为是可以正确复制html5元素的
    // 能正确复制为true
    html5Clone: document.createElement("nav").cloneNode( true ).outerHTML !== "<:nav></:nav>",

    // 后续的几项在后面再进行测试
    submitBubbles: true,
    changeBubbles: true,
    focusinBubbles: false,
    deleteExpando: true,
    noCloneEvent: true,
    inlineBlockNeedsLayout: false,
    shrinkWrapBlocks: false,
    reliableMarginRight: true
  };

  // 测试能不能正确的复制input的checked值
  input.checked = true;
  support.noCloneChecked = input.cloneNode( true ).checked;

  // 当select设为disabled时，子元素option能否也自动设为disabled
  select.disabled = true;
  support.optDisabled = !opt.disabled;

  // 测试浏览器是否允许删除dom元素的属性
  // 这里使用了try因为直接报错容易崩溃
  // ie不允许删除dom元素的属性
  try {
    delete div.test;
  } catch( e ) {
    support.deleteExpando = false;
  }

	// ie在复制元素的时候会连同事件一起复制了，要测试是否存在这个兼容问题
	// 在没有addEventListener但有attachEvent的浏览器，即ie上测试
	// 先为div添加一个相应事件，on开通哦。 
	// 复制元素并且执行onclike， 如果执行了，那么noCloneEvent就为false了。
	// 如果没有复制事件，那么也就不会执行，就会是默认true
  if ( !div.addEventListener && div.attachEvent && div.fireEvent ) {
    div.attachEvent( "onclick", function() {
      // Cloning a node shouldn't copy over any
      // bound event handlers (IE does this)
      support.noCloneEvent = false;
    });
    div.cloneNode( true ).fireEvent( "onclick" );
  }

  // 测试radio类型会不会让input的value丢失掉
  input = document.createElement("input");
  input.value = "t";
  input.setAttribute("type", "radio");
  support.radioValue = input.value === "t";


  // 测试浏览器能否在文档片段中正确的复制checked状态值
  input.setAttribute("checked", "checked");
  div.appendChild( input );
  fragment = document.createDocumentFragment();
  fragment.appendChild( div.lastChild );
  support.checkClone = fragment.cloneNode( true ).cloneNode( true ).lastChild.checked;

  // 测试将一个input插入到文档片段中，check属性是否保留
  support.appendChecked = input.checked;

  fragment.removeChild( input );
  fragment.appendChild( div );



  div.innerHTML = "";


  // 这是一个webkit的小bug https://bugs.webkit.org/show_bug.cgi?id=13343
  // 不会对普通block element返回margin-right值
  // 只要得出的结果不是0那么就是存在这个bug
  // 正确计算marginright值为true
  if ( window.getComputedStyle ) {
    marginDiv = document.createElement( "div" );
    marginDiv.style.width = "0";
    marginDiv.style.marginRight = "0";
    div.style.width = "2px";
    div.appendChild( marginDiv );
    support.reliableMarginRight =
      ( parseInt( ( window.getComputedStyle( marginDiv, null ) || { marginRight: 0 } ).marginRight, 10 ) || 0 ) === 0;
  }

	// 只在ie上测试
	// 对submit、change、focusin循环处理
	// 如果div上面没有这个认为不支持
	// 如果不支持 再试试setAttribute能不能正确设置
	// 最后设置最终的是否支持
	// 通过名字来看是测试能否冒泡，但是没看见有冒泡的检测啊？
	// 原因在于 div本身不应该有submit、change、focus这类事件，如果有则表明支持冒泡
  if ( div.attachEvent ) {
    for( i in {
      submit: 1,
      change: 1,
      focusin: 1
    }) {
      eventName = "on" + i;
      isSupported = ( eventName in div );
      if ( !isSupported ) {
        div.setAttribute( eventName, "return;" );
        isSupported = ( typeof div[ eventName ] === "function" );
      }
      support[ i + "Bubbles" ] = isSupported;
    }
  }

  fragment.removeChild( div );

  // Null elements to avoid leaks in IE
  fragment = select = opt = marginDiv = div = input = null;

  // 后面的测试需要等到文档加载完，需要body元素
  jQuery(function() {
    var container, outer, inner, table, td, offsetSupport,
      conMarginTop, ptlm, vb, style, html,
      body = document.getElementsByTagName("body")[0];

    // 没有body 直接return掉
    if ( !body ) {
      return;
    }

    conMarginTop = 1;
    ptlm = "position:absolute;top:0;left:0;width:1px;height:1px;margin:0;";
    vb = "visibility:hidden;border:0;";
    style = "style='" + ptlm + "border:5px solid #000;padding:0;'";
    html = "<div " + style + "><div></div></div>" +
      "<table " + style + " cellpadding='0' cellspacing='0'>" +
      "<tr><td></td></tr></table>";

    container = document.createElement("div");
    container.style.cssText = vb + "width:0;height:0;position:static;top:0;margin-top:" + conMarginTop + "px";
    body.insertBefore( container, body.firstChild );

    // Construct the test element
    div = document.createElement("div");
    container.appendChild( div );

    // https://github.com/jquery/jquery/pull/35
		// 构建了一个table表格，一行内第一个td无内容，第二个td有内容
		// 测试这种情况下第一个td无论display是否为none都应该高度为0
		// 空td的offsetheight为0 则true
    div.innerHTML = "<table><tr><td style='padding:0;border:0;display:none'></td><td>t</td></tr></table>";
    tds = div.getElementsByTagName( "td" );
    isSupported = ( tds[ 0 ].offsetHeight === 0 );

    tds[ 0 ].style.display = "";
    tds[ 1 ].style.display = "none";

    support.reliableHiddenOffsets = isSupported && ( tds[ 0 ].offsetHeight === 0 );

    // 通过检测offsetWidth来判断是不是支持标准盒模型
    div.innerHTML = "";
    div.style.width = div.style.paddingLeft = "1px";
    jQuery.boxModel = support.boxModel = div.offsetWidth === 2;


    if ( typeof div.style.zoom !== "undefined" ) {
			// 测试对于inline的元素是不是正常显示
			// 通过设置div display: inline zoom: 1
			// 再去查看它的offsetWidth，正常应该为0
			// 如果inline之后还有offsetWidth则 inlineBlockNeedsLayout为true
      div.style.display = "inline";
			// 设置zoom主要是在ie环境下为元素出发hasLayout属性
      div.style.zoom = 1;
      support.inlineBlockNeedsLayout = ( div.offsetWidth === 2 );

			// 构建宽度大于父元素的子元素，然后测试父元素的可见宽度
			// 如果父元素可见宽度增大 那么就是shrinkWrapBlock为true
      div.style.display = "";
      div.innerHTML = "<div style='width:4px;'></div>";
      support.shrinkWrapBlocks = ( div.offsetWidth !== 2 );
    }

		// ptlm = "position:absolute;top:0;left:0;width:1px;height:1px;margin:0;";
    vb = "visibility:hidden;border:0;";
    style = "style='" + ptlm + "border:5px solid #000;padding:0;'";
    html = "<div " + style + "><div></div></div>" +
      "<table " + style + " cellpadding='0' cellspacing='0'>" +
      "<tr><td></td></tr></table>";

    div.style.cssText = ptlm + vb;
    div.innerHTML = html;

    outer = div.firstChild;
    inner = outer.firstChild;
    td = outer.nextSibling.firstChild.firstChild;

    offsetSupport = {
			// ie8下offsetTop会包含父元素的上border宽度
      doesNotAddBorder: ( inner.offsetTop !== 5 ),
			// IE FireFox中td距离tr的offsetTop会包含table的上边框宽度
      doesAddBorderForTableAndCells: ( td.offsetTop === 5 )
    };

    inner.style.position = "fixed";
    inner.style.top = "20px";

    // safari subtracts parent border width here which is 5px
		// 判断能否正确的计算fixed元素的位置
    offsetSupport.fixedPosition = ( inner.offsetTop === 20 || inner.offsetTop === 15 );
    inner.style.position = inner.style.top = "";

    outer.style.overflow = "hidden";
    outer.style.position = "relative";
		// hidden的父元素，子元素在计算offseTop时会不会减去父元素的变宽
		// 主流浏览器都是false
    offsetSupport.subtractsBorderForOverflowNotVisible = ( inner.offsetTop === -5 );
    // 测试body元素的offsetTop包含不包含body的外边距
		// 可是没看懂为什么和一个固定值比较
		offsetSupport.doesNotIncludeMarginInBodyOffset = ( body.offsetTop !== conMarginTop );

    body.removeChild( container );
    div  = container = null;

    jQuery.extend( support, offsetSupport );
  });

  return support;
})();
```

## 总结
其实写到这里我都不知道自己写了些什么，而且这里面很多兼容性测试都是第一次听说，甚至连`onfocusin`都是第一次见，可能工作时间短、再加上不是to c的原因吧，对于兼容性的知识确实了解不多。这就正需要补习这方面知识了，前端创造者为了提供更完美的效果也应该使页面尽可能的兼容更多的平台。 最后再总体总结一下这里面的兼容性测试项目。至于兼容性测试之后对应的解决方案，会在之后用到的时候再分析。


| 测试项目        | 含义           | 
| ------------- |:-------------:|
|dom测试|
| leadingWhitespace      | innerHTML会不会保留前导空白符 |
| tbody     | 会不会为空的table增加tbody元素 |
| htmlSerialize | innerHTML能否正确序列化link标签 |
| hrefNormalized | a标签的href会不会被补为绝对路径 |
| checkOn | checkbox的value是否为on |
| getSetAttribute | 测试getAttribute、setAttribute能否正常工作 |
| enctype | 表单是否支持enctype |
| html5Clone | 浏览器能否正确的复制html5元素 |
| deleteExpando | 能否删除DOM元素上的属性 |
| optSelected | option默认是否选中 |
| noCloneChecked | 复制元素时、是否会复制check状态 |
| optDisabled | disable select标签时、内部option会不会一起被禁 |
| radioValue | radio的input会不会丢失value |
| checkClone | 文档片段中能否正确复制check状态 |
| appendChecked | 已选中的check在添加到文档时，能否正确的保留状态 |
|css测试|
| style | inline style能不能直接element.style访问 |
| opacity | 浏览器是否支持opacity |
| cssFloat | 能否支持style.cssFloat |
|盒模型测试|
| inlineBlockNeedsLayout | display:inline; zoom: 1之后能否按照inline-block来显示 |
| shrinkWrapBlocks | ie中、拥有haslayout的固定width、height元素会不会被子元素撑开 |
| reliableMarginRight | 能否返回正确的右边距，仅chrome某个小版本有问题 |
| reliableHiddenOffsets | 空单元格可见高度是否为0 |
| boxModel | 是不是标准W3C盒子模型 |
| doesNotAddBorder | 子元素距父元素的offsetTop是否包含父元素的border-top-width |
| doesAddBorderForTableAndCells | td距tr的offset是否包含table的border-top-width |
| fixedPosition | 正确返回fixed元素的窗口坐标 |
| subtractsBorderForOverflowNotVisible | 父元素overflow:hidden，子元素的offsetTop减去父元素的border-width |
| doesNotIncludeMarginInBodyOffset | body距html边框的距离是否包括body的margin |
|事件测试|
| noCloneEvent | 复制元素时是否复制事件 |
| submitBubbles | submit事件能正常冒泡 |
| changeBubbles | change事件能正常冒泡 |
| focusinBubbles | focusin事件能正常冒泡 |

未来还有很长的路要走，兼容是不可回避的问题