# chapter2
在[chapter1](https://github.com/zhaozy93/blog/issues/2)中对jquery整体结构和最基本的部分有一个简单的认识了。 现在来看jquery是如何帮我们查找dom元素的。

## Sizzle
Sizzle既是一个单独的用js实现的css选择器包，也是jquery不可分割的一部分。在低级浏览器中jquery正是靠着Sizzle才能简单的稳定的帮我们获取到想要的dom元素集合。

我们都知道css的查询是从右向左查即div > p 先查所有的p再筛选出满足父元素为div的p。 关于css左查和右查的规则以及优劣可以看[blog](http://www.cnblogs.com/zhaodongyu/p/3341080.html)、[segmentfault](https://segmentfault.com/q/1010000000713509)。

在Sizzle中基本也是与css规则相同即从右向左查询，进行逐步筛选，逐步缩小的方法。除了有伪类存在时是从左向右查。原因是 div button:first 是查找div下面的第一个button，位置伪类过滤的是 div button而不是button。关于这一点之后慢慢分析。

## Sizzle in jQuery
在jQuery中的Sizzle位置大概在 3859-5302 行之间，大约1500行左右。 由一个立即执行函数包裹，在内部最后通过把一些方法挂载到jQuery变量上面来实现让jQuery的访问与使用，可以看到jQuery.find其实就是Sizzle的主入口。 立即执行函数的好处就不再叙述了，与jQuery的整体思想一致。 内部实现方法、定义变量，最终将接口抛出来：jQuery挂在window上，Sizzle挂在jQuery上。
```
(function(){
........
Sizzle.attr = jQuery.attr;
Sizzle.selectors.attrMap = {};
jQuery.find = Sizzle;
jQuery.expr = Sizzle.selectors;
jQuery.expr[":"] = jQuery.expr.filters;
jQuery.unique = Sizzle.uniqueSort;
jQuery.text = Sizzle.getText;
jQuery.isXMLDoc = Sizzle.isXML;
jQuery.contains = Sizzle.contains;
})();
```

## 源码

### chunker
在Sizzle整个代码的第一行是一个非常长的正则表达式--chunker, 此行正则表达式竟然有150个字符，简直是一个逆天的正则。那就从它开始看。
```
var chunker = /((?:\((?:\([^()]+\)|[^()]+)+\)|\[(?:\[[^\[\]]*\]|['"][^'"]*['"]|[^\[\]'"]+)+\]|\\.|[^ >+~,(\[\\]+)+|[>+~])(\s*,\s*)?((?:.|\r|\n)*)/g,
```
首先拆分一下正则，利用()分为三部分
```
1: ((?:\((?:\([^()]+\)|[^()]+)+\)|\[(?:\[[^\[\]]*\]|['"][^'"]*['"]|[^\[\]'"]+)+\]|\\.|[^ >+~,(\[\\]+)+|[>+~])
2: (\s*,\s*)?   // \s表示一个单个的空白符， 因此匹配的就是单个,(逗号)，甚至可以不匹配 css选择器表达式中逗号,表示并列选择器
3: ((?:.|\r|\n)*)    // ,后面的并列选择器的文本
```
第一部分还是太长 继续拆分，利用 |

```
1—1: (?:\((?:\([^()]+\)|[^()]+)+\)|\[(?:\[[^\[\]]*\]|['"][^'"]*['"]|[^\[\]'"]+)+\]|\\.|[^ >+~,\[\\]+)+
1-2: [>+~]
```
还是太长 继续拆分利用 | 和 ()
```
1_1_1: \((?:\([^()]+\)|[^()]+)+\)     // 匹配 ((tag))或者(tag)  并且tag必须存在
1_1_2: \[(?:\[[^\[\]]*\]|['"][^'"]*['"]|[^\[\]'"]+)+\]  // 匹配 [[]] 或者[""]、[''] 或者 [tag]
1_1_3: \\.
1_1_4: [^ >+~,(\[\\]+  // 不包含关系符号、属性符号[、伪类(和转义符\，可以匹配.class、 #id
1_1_5: [>+~]  // 匹配层级选择器 > + ~
```

因此整个chunker的第一部分就是选择器表达式，第二部分表示并列选择器表达式的flag，第三部分表示并列表达式。第一部分的内部又分为伪类()部分、属性[]部分、引号包裹部分、简单选择器部分。 当然这样简单的处理理解肯定是不足的。具体的之后再继续深入，先知道chunker大概。

### Sizzle入口

#### parts.1
在立即执行函数内定义一个Sizzle函数变量，最后挂在jQuery上，作为Sizzle的入口函数。

Sizzle接受4个参数，selector选择器字符串、context上下文、results接收结果的容器、seed结果将从该元素中进行过滤
```
var Sizzle = function( selector, context, results, seed ) {
  // 常规的默认参数修正
  results = results || [];
  context = context || document;
  // 记录原上下文
  var origContext = context;

  if ( context.nodeType !== 1 && context.nodeType !== 9 ) {
    return [];
  }
  
  if ( !selector || typeof selector !== "string" ) {
    return results;
  }

  var m, set, checkSet, extra, ret, cur, pop, i,
    prune = true,
    contextXML = Sizzle.isXML( context ),
    parts = [],
    soFar = selector;
  // 采用do-while的形式循环匹配出所有的 选择表达式
  do {
    // 执行一次正则，但永远不会得到匹配，这样的作用是chunker.lastIndex = 0
    // 因为带g的正则都有lastIndex属性，每次匹配的时候都是从这个位置开始，
    // 关于lastIndex可以看[blog](http://www.cnblogs.com/sniper007/archive/2011/12/20/2295044.html),但这篇文章Sizzle解析的部分个人感觉是错误的
    chunker.exec( "" );
    // 真正执行一次匹配
    m = chunker.exec( soFar );
    

    // 如果解析成功，那么将这次匹配到的内容压入parts数组，并且如果有,这个flag则更改将多余的部分付给extra
    if ( m ) {
      soFar = m[3];
    
      parts.push( m[1] );
      // 在什么情况下m[2]会存在，什么情况下会不存在呢？
      // m[2]存在比较好理解，就是包含 , 的字符串
      // m[2]不存在的情况就是 div > p这样提取出了m[1] = div, m[2]=undeifned， m[3]=' > p'
      if ( m[2] ) {
        extra = m[3];
        break;
      }
    }
  } while ( m );
  .....
```

对这段代码单独提取一下并且测试一下会发现,并且输出parts数组， 这段代码直接帮我们把一个`选择器表达式`进行拆解了， 同时关系也比较明确
```
exec('div:children > div , div')     // ["div:children", ">", "div"]
exec('div div')   // ["div", "div"]
```
#### parts.2
继续看代码

这一次稍长一点 是一个if-else，条件是parts长度大于1并且有伪类存在选择器中。
origPos = /:(nth|eq|gt|lt|first|last|even|odd)(?:\((\d*)\))?(?=[^\-]|$)/
因为前面也说过了伪类选择器是Sizzle决定左查还是右查的一个重要指标
```
  // 根据parts长度和是否有伪类选择器决定左查或者右查
  if ( parts.length > 1 && origPOS.exec( selector ) ) {
    // 满足条件则是从左向右查找
    // 如果parts长度为2并且第一个为关系符 ? > ~ 则直接调用posProcess方法
    // 说明这是一个比较简单的 其实只有一个选择器的查询
    if ( parts.length === 2 && Expr.relative[ parts[0] ] ) {
      set = posProcess( parts[0] + parts[1], context, seed );
    } else {
      // 否则的话再判断第一位是不是关系符决定set的初始值
      set = Expr.relative[ parts[0] ] ?
        [ context ] :
        Sizzle( parts.shift(), context );
      // 依旧用while循环来依次取出parts的首位并且调用posProcess方法，直到parts长度为0
      while ( parts.length ) {
        selector = parts.shift();

        if ( Expr.relative[ selector ] ) {
          selector += parts.shift();
        }
        
        set = posProcess( selector, set, seed );
      }
    }
    // 在这个if中可以看出来posProcess才是真正干活的， 每次传进去一个选择器selector，然后再传入set，set就是上一次得出的结果， 可以理解为上下文吧，然后每次进行缩小操作， 从左向右嘛，一直在做减法
  } else {
    // 上面是从左向右， 这里就是从右向左了
    // 满足一系列条件的话就修正context和初始化ret
    // parts长度大于1， 首个元素是id， 最后一个不是id
    if ( !seed && parts.length > 1 && context.nodeType === 9 && !contextXML &&
        Expr.match.ID.test(parts[0]) && !Expr.match.ID.test(parts[parts.length - 1]) ) {
      ret = Sizzle.find( parts.shift(), context, contextXML );
      context = ret.expr ?
        Sizzle.filter( ret.expr, ret.set )[0] :
        ret.set[0];
    }

    // 如果前一个if不满足条件的话context在这里一定会有值
    // 只有前一个if执行了，但是根据第一个id选择器没有找到元素，修正context失败了， 那么也意味着没有符合条件的元素，之内else返回[]
    if ( context ) {
      // 初始化参数， 如果有seed则直接食用seed作为候选集，否则调用最后一个parts进行候选集初始化
      ret = seed ?
        { expr: parts.pop(), set: makeArray(seed) } :
        Sizzle.find( parts.pop(), parts.length === 1 && (parts[0] === "~" || parts[0] === "+") && context.parentNode ? context.parentNode : context, contextXML );
      // 元素的候选集，上下文，从右向左是一直做减法
      set = ret.expr ?
        Sizzle.filter( ret.expr, ret.set ) :
        ret.set;

      // 如果parts还有元素就对set做个备份， 否则的话prune设为false
      if ( parts.length > 0 ) {
        checkSet = makeArray( set );
      } else {
        prune = false;
      }

      // 相同的把戏，这次是pop即parts从右向左依次推出，来慢慢缩小查找范围
      while ( parts.length ) {
        // 取出最后一个元素
        cur = parts.pop();
        pop = cur;
        
        // 如果最后一个元素不是关系符, 则把cur设为 ''， pop则重新在parts中取出最后一个
        // 这样确保pop为选择器， cur为选择符号
        if ( !Expr.relative[ cur ] ) {
          cur = "";
        } else {
          pop = parts.pop();
        }
        // 当发现pop为null时即上一步pop时已经是首位了，[0]之前没有了， 则直接把context设为pop
        if ( pop == null ) {
          pop = context;
        }
        // 然后在利用关系符、checkSet、pop来缩小查找范围，结果通过checkSet来传递
        Expr.relative[ cur ]( checkSet, pop, contextXML );
      }

    } else {
      checkSet = parts = [];
    }
  }
```

#### parts.3
最后进行统一的处理
```
  // 如果没有checkSet，则把set值赋给checkSet，因为从左向右查没有set，从右向左查并且只有一个表达式时也没有checkSet
  if ( !checkSet ) {
    checkSet = set;
  }

  if ( !checkSet ) {
    Sizzle.error( cur || selector );
  }
  // 如果checkSet也就是前面的结果是数组 再决定要不要筛选，最终将结果打入results中，最后返回results
  if ( toString.call(checkSet) === "[object Array]" ) {
    if ( !prune ) {
      results.push.apply( results, checkSet );

    } else if ( context && context.nodeType === 1 ) {
      for ( i = 0; checkSet[i] != null; i++ ) {
        if ( checkSet[i] && (checkSet[i] === true || checkSet[i].nodeType === 1 && Sizzle.contains(context, checkSet[i])) ) {
          results.push( set[i] );
        }
      }

    } else {
      for ( i = 0; checkSet[i] != null; i++ ) {
        if ( checkSet[i] && checkSet[i].nodeType === 1 ) {
          results.push( set[i] );
        }
      }
    }

  } else {
    makeArray( checkSet, results );
  }

  // 如果还有extra，则继续调用Sizzle入口函数，并将results和OrigContext传入，然后再筛选结果。
  if ( extra ) {
    Sizzle( extra, origContext, results, seed );
    Sizzle.uniqueSort( results );
  }

  return results;
```

通过对Sizzle入口函数的解析，过程可归结如下
1、先对选择器整体的字符串进行拆解
2、无论是从左到右还是从右到左，一步步调用方法来不断减少范围
3、最后对两种方向得出的结果进行统一，最后返回结果。

里面真正做事情的函数不是特别多，大概有三个主要出现的。 `posProcess`, `Sizzle.find`,`Expr.relative[ cur ]`


### Sizzle.find()
用于查找与表达式匹配的集合。

在Sizzle函数的解析中我们发现真正做事情的内容不多， Sizzle.find算一个。
```
Sizzle.find = function( expr, context, isXML ) {
  var set, i, len, match, type, left;
  if ( !expr ) {
    return [];
  }
  // Expr.order = [ "ID", "NAME", "TAG" ];
  // 依次检测id， name和tag
  for ( i = 0, len = Expr.order.length; i < len; i++ ) {
    type = Expr.order[i];
    // 利用leftMatch正则来依次检测这个选择器是否符合要求
    if ( (match = Expr.leftMatch[ type ].exec( expr )) ) {
      // 进入if则表示正则匹配成功了
      left = match[1];
      match.splice( 1, 1 );
      // 这里优点不明白的是为一个字符串执行substr(length-1)之后应该永远不可能等于'\\'，只可能等于一个长度
      // 总之就是成功了就继续Expr.find来获取元素集合
      if ( left.substr( left.length - 1 ) !== "\\" ) {
        match[1] = (match[1] || "").replace( rBackslash, "" );
        set = Expr.find[ type ]( match, context, isXML );
        if ( set != null ) {
          // 如果元素集合获取成功了则把已经用过的expr部分删除掉
          expr = expr.replace( Expr.match[ type ], "" );
          break;
        }
      }
    }
  }

  if ( !set ) {
    // 走到这里说明没有匹配到任何元素，相对应expr也没有被删掉任何一部分，这样就获取当前元素的所有子元素作为集合，expr再还回去
    set = typeof context.getElementsByTagName !== "undefined" ?
      context.getElementsByTagName( "*" ) :
      [];
  }
  // set表示找到的所有元素集合， expr表示剩余的选择器部分。
  return { set: set, expr: expr };
};
```

这里面发现真正干活的是`Expr.find`


### Sizzle.filter()
利用块表达式来过滤元素集合
```
Sizzle.filter = function( expr, set, inplace, not ) {
  var match, anyFound,
    type, found, item, filter, left,
    i, pass,
    old = expr,
    result = [],
    curLoop = set,
    isXMLFilter = set && set[0] && Sizzle.isXML( set[0] );

  while ( expr && set.length ) {
    for ( type in Expr.filter ) {
      if ( (match = Expr.leftMatch[ type ].exec( expr )) != null && match[2] ) {
        filter = Expr.filter[ type ];
        left = match[1];

        anyFound = false;

        match.splice(1,1);

        if ( left.substr( left.length - 1 ) === "\\" ) {
          continue;
        }

        if ( curLoop === result ) {
          result = [];
        }

        if ( Expr.preFilter[ type ] ) {
          match = Expr.preFilter[ type ]( match, curLoop, inplace, result, not, isXMLFilter );

          if ( !match ) {
            anyFound = found = true;

          } else if ( match === true ) {
            continue;
          }
        }

        if ( match ) {
          for ( i = 0; (item = curLoop[i]) != null; i++ ) {
            if ( item ) {
              found = filter( item, match, i, curLoop );
              pass = not ^ found;

              if ( inplace && found != null ) {
                if ( pass ) {
                  anyFound = true;

                } else {
                  curLoop[i] = false;
                }

              } else if ( pass ) {
                result.push( item );
                anyFound = true;
              }
            }
          }
        }

        if ( found !== undefined ) {
          if ( !inplace ) {
            curLoop = result;
          }

          expr = expr.replace( Expr.match[ type ], "" );

          if ( !anyFound ) {
            return [];
          }

          break;
        }
      }
    }

    // Improper expression
    if ( expr === old ) {
      if ( anyFound == null ) {
        Sizzle.error( expr );

      } else {
        break;
      }
    }

    old = expr;
  }

  return curLoop;
};
```



---
未完待续
Sizzle暂时不继续进行分析， 对js帮助并不是特别大。
先着重进行jQuery内容的分析

