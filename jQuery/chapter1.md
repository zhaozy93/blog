# chapter1
* 解析jquery版本为1.7.1 没有解读最新版是因为水平有限，希望先借助当前已有资料进行练习
* 规定jquery通常表示jquery这个文件或这个包。
* jQuery表示内部的变量等代码层面的内容。
## jquery 整体
jquery整个文件是一个立即执行函数。 这样做的好处呢就是不会污染全局变量，同时全局的其余包或者变量也不会影响jquery的内部工作。 
不污染全局且不被外部影响是一个包被广泛使用的基本要求之一。
关于立即执行函数就不在解释。可以查看  
[github](https://github.com/AllFE/blog/issues/3)

  整个立即执行函数传入两个对象，window和undefined：
* 1、内部大量用到window和undefined，作为内部局域变量可以减少作用域链的查找，这样加快速度
* 2、方便代码压缩， a=window  b=undefine


整个代码的最后 window.jQuery = window.$ = jQuery 将jQuery对象抛出，并挂载在window上，这样可以全局任意位置随时调用。 日常的调用习惯$(params)可以推算出jQuery应该属于一个可执行的函数。

借助IDE可以轻松找到 window.jQuery = window.$ = jQuery 中jQuery(变量1)的定义位置， 23-958行。那么文件中起于的代码全部都是为这个jQuery变量服务的。恰巧这个jQuery(变量1)变量指向一个立即执行函数。这个函数最终返回的是一个内部变量也叫jQuery(变量2)。
换句话说， jquery文件是一个立即执行函数， 里面定义一个jQuery(变量1)变量，这个变量也指向一个立即执行函数。 这个函数最终返回的是一个指向具体方法的变量(变量2)。 也就是说整个jquery在调用的时候 其实是在调用一个具体的方法。这个方法也就是整个jquery的入口函数。因此下一步是要搞清楚`new jQuery.fn.init( selector, context, rootjQuery )` 做了哪些工作。

这里可能有点迷惑的就是定义了太多的jQuery，其实自习看一下作用域，删减掉部分代码 就比较能清晰的看出各个jQuery的级别或者作用域，然后就知道每次return的是哪个jQuery。
  ```
  // 下面定义的是jQuery(变量2)
  var jQuery = function( selector, context ) {
      // The jQuery object is actually just the init constructor 'enhanced'
      return new jQuery.fn.init( selector, context, rootjQuery );
    }
  ```
  ```
  // 下面定义的是jQuery(变量1)
  var jQuery = (function() {
    // 下面一行是定义jQuery(变量2)哦。
    var jQuery = function( selector, context ) {
      return new jQuery.fn.init( selector, context, rootjQuery );
    }
    // ....
    // ...
    // ...
    return jQuery  // 因此 这里return的是jQuery(变量2)
  })();
  window.jQuery = window.$ = jQuery   // 因此 这里window上挂载的是jQuery(变量1)
  ```
在jQuery(变量2)的var指令时，同时定义了高达40个左右的变量，尽管里面有注释，但有非常多的正则不太容易理解。因此我们不必现在理解他们，用到了再回头看。或者后期专门一节来分析所有正则。 这些变量中还有很多变量是用于缩写的， 可能这也是一个好习惯吧。

## jQuery.fn.init

又发现一个很有意思的地方就是 jQuery.fn = jQuery.prototype = {}，简而言之，jQuery.fn 又是一个缩写，它指向的是jQuery这个函数的原型链。 我们上一节讲到关心的是jQuery.fn.init这个函数。那么这个函数其实就是jQuery.prototype原型链上面定义的一个方法而已。为什么定义在原型链上，也就是原型链的好处可以待会再讲。

继续看代码 jQuery.fn = jQuery.prototype下面是 `jQuery.fn.init.prototype = jQuery.fn;`，这句话代表了什么呢。等价于
```
jQuery.fn = jQuery.prototype; 
jQuery.fn.init.prototype = jQuery.prototype;
```
这里就可以解释一下原型链的好处了。 原型链最大的好处也就是实现了类似于继承的概念。当一个属性在对象本身上面找不到的时候，他就会去寻找她的原型链，直到寻找到最顶端。 因为无论jQuery还是$调用的入口方法返回的是一个`new jQuery.fn.init`实例， 恰巧jquery还实现了链式写法，链式写法最基本的就是每个方法调用后要返回本身this，才能保证下一个方法可以继续执行。在这里就可以知道虽然所有的方法不会挂载在init上，应该属于jQuery，但是通过将两者的原型链共享，那么也可以让init随意调用jQurry的方法了。 注意：jquery链式操作返回的this应该是属于fn.init的。

看一下fn.init的源码, 这个正是我们使用$()时的调用方法。 然后接下来一行行分析这个init函数。
约定：
* 源码中所带的英文注释全部被删掉以方便查看。
* selector译为选择器
* context译为上下文
* rootjQuery译为根元素
```
 init: function( selector, context, rootjQuery ) {
      var match, elem, ret, doc;

      // 处理 $(""), $(null), or $(undefined)三种情况， 当选择器为假时实例一个空jQuery留作以后使用
      if ( !selector ) {
        return this;
      }

      // 处理 $(DOMElement) 选择器为原生dom元素时， 将上下文也设为该dom元素，并且将该元素压入数组， 同时可以得知其实this是一个数组
      if ( selector.nodeType ) {
        this.context = this[0] = selector;
        this.length = 1;
        return this;
      }

      // 处理选择器为body字符串并且没有上下文的情况
      // 上下文设为document，并且将body压入this数组
      if ( selector === "body" && !context && document.body ) {
        this.context = document;
        this[0] = document.body;
        this.selector = selector;
        this.length = 1;
        return this;
      }

      // 处理选择器是字符串的情况 
      if ( typeof selector === "string" ) {
        // 需要区分是普通字符串 还是 #id类型
        if ( selector.charAt(0) === "<" && selector.charAt( selector.length - 1 ) === ">" && selector.length >= 3 ) {
          // 假设字符串 为<开头>结尾并且长度大于3 就为字符串dom元素
          match = [ null, selector, null ];

        } else {
          // 否则用正则去判断string
          match = quickExpr.exec( selector );
        }
        // A simple way to check for HTML strings or ID strings
        // Prioritize #id over <tag> to avoid XSS via location.hash (#9521)
        // quickExpr = /^(?:[^#<]*(<[\w\W]+>)[^>]*$|#([\w\-]*)$)/;
        // 这个正则整体就是一个大括号， 大括号内部又分成两个括号
        // 但是大括号的开头时 ?: 也就是说这个()匹配的内容是不会被单独记录的。
        // 正则的两个小()被 | 分割也就是说可以匹配两个中任意一个即可
        // 第二部分： #([\w\-]*)$ 比较好理解 匹配#div这样的形式 当然后面id的名字也可以包含-
        // 第一部分： [^#<]*(<[\w\W]+>)[^>]*$ 
        // [^#<]* 表示开头不是#或<，长度为0或n
        // (<[\w\W]+>) <开始 中间任意字符串长度大于1 >结尾 
        // [^>]*$ 不以>结尾
        // 这样直接就把所有#开头的交给了第二个括号去匹配。因为如果#开头则表示 [^#<]* 匹配到0次(也不能算错)，但是也不是<，则匹配失败了。
        // 如果是以<开头的 则仍表示[^#<]* 匹配到0次， 然后正常匹配 <> 并且以贪婪模式保留<>中间的内容作为()的内容输出
        // 如果是以其余字符串形式开头表示[^#<]* 匹配到多次， 但是会被抛弃，只去寻找<>
        // 匹配的结果是一个数组， 长度为3
        // [ selector, htmlString, id]

        // 正则匹配成功， 并且 （ 匹配到了htmlString 或者 没有提供上下文）
        if ( match && (match[1] || !context) ) {

          // 匹配到了htmlString
          if ( match[1] ) {
            // 处理上下文以及正确的document
            context = context instanceof jQuery ? context[0] : context;
            doc = ( context ? context.ownerDocument || context : document );

            // rsingleTag = /^<(\w+)\s*\/?>(?:<\/\1>)?$/
            // 又是一个正则
            // ^< 以<开头
            // (\w+) 括号包裹，内部至少一个字符
            // \s* 0或多个空白符
            // \/? 0或1个/  反斜杠
            // (?:<\/\1>)? ?: 表示括号内的内容不用被记录输出， 最后一个？表示0或1次 ()内部其实是</加第一个括号内容> 其实就是html的闭合标签
            // <div></div> 通过
            // <div></div1> 不通过 
            // <div>gdfgdf</div1> 不通过
            // <div /></div> 也通过
            // 其实rsingleTag 就是用来检测 htmlString是不是单个简单标签，而且不含innerHTML和任何text
            ret = rsingleTag.exec( selector );

            if ( ret ) {
              // jQuery.isPlainObject是一个检测传入参数是不是 纯对象的方法， 只有当参数为纯javascript且没有Ownproperty才返回true
              // 这里当htmlString是简单标签时又根据上下文分了两种处理方式
              if ( jQuery.isPlainObject( context ) ) {
                selector = [ document.createElement( ret[1] ) ];
                jQuery.fn.attr.call( selector, context, true );

              } else {
                selector = [ doc.createElement( ret[1] ) ];
              }

            } else {
              // htmlString为复杂标签的处理方式
              ret = jQuery.buildFragment( [ match[1] ], [ doc ] );
              selector = ( ret.cacheable ? jQuery.clone(ret.fragment) : ret.fragment ).childNodes;
            }

            return jQuery.merge( this, selector );

          // else处理 选择器为字符串， 并且为id选择器的情况
          } else {
            elem = document.getElementById( match[2] );

            // Check parentNode to catch when Blackberry 4.6 returns
            // nodes that are no longer in the document #6963
            if ( elem && elem.parentNode ) {
              // Handle the case where IE and Opera return items
              // by name instead of ID
              if ( elem.id !== match[2] ) {
                return rootjQuery.find( selector );
              }

              // Otherwise, we inject the element directly into the jQuery object
              this.length = 1;
              this[0] = elem;
            }

            this.context = document;
            this.selector = selector;
            return this;
          }

        // 处理上下文不存在或者上下文为jquery对象时
        } else if ( !context || context.jquery ) {
          return ( context || rootjQuery ).find( selector );

        // 处理上下文存在的情况
        } else {
          return this.constructor( context ).find( selector );
        }

      // 处理选择器为函数的情况
      } else if ( jQuery.isFunction( selector ) ) {
        return rootjQuery.ready( selector );
      }

      // 如果选择器有selector属性， 则认为它是一个jquery对象，将其selector和context赋值给当前this
      if ( selector.selector !== undefined ) {
        this.selector = selector.selector;
        this.context = selector.context;
      }
      // 应对其余情况， 同时当选择器为jquery对象时也会被执行
      return jQuery.makeArray( selector, this );
    },
```
至此， 我们调用$()无论传入任何参数， 内部方法已经实现了分流。 分成几种不同情况去处理，有字符串，有domElement，还有jquery，当然还有假值。 当然现在也只能看到分流这一步， 分流后各小项是如何处理的当前还不得而知， 同时jQuery.fn也仅仅分析了init一个属性。 下一节 会继续分析jQuery.fn的其余属性。

## buildFragment
在前一节看到 当选择器selector为复杂string的时候，会调用buildFragment方法来产生相对应的dom元素。
```
 // htmlString为复杂标签的处理方式
ret = jQuery.buildFragment( [ match[1] ], [ doc ] );
```
这里match[1]为字符串 doc为矫正后的document
```
args 带转换为dom元素的字符串或html代码
nodes 用于修正创建文档片段DocumentFragment的文档对象
script 存放html代码中script元素
jQuery.buildFragment = function( args, nodes, scripts ) {
	var fragment, cacheable, cacheresults, doc,
	first = args[ 0 ];
  // 继续矫正document 确保后面不出现异常
	if ( nodes && nodes[0] ) {
		doc = nodes[0].ownerDocument || nodes[0];
	}
	if ( !doc.createDocumentFragment ) {
		doc = document;
	}
  // 一系列检测 决定是否可以缓存这段 待转换文本 
  if ( args.length === 1 && typeof first === "string" && first.length < 512 && doc === document &&
		first.charAt(0) === "<" && !rnocache.test( first ) &&
		(jQuery.support.checkClone || !rchecked.test( first )) &&
		(jQuery.support.html5Clone || !rnoshimcache.test( first )) ) {

		cacheable = true;
    // 如果符合缓存要求， 尝试从缓存仓库fragments中读取缓存的内容，
    // 可以理解为dom转换比较费时费力 不希望重复做两次同样的事情
    // jQuery.fragments = {};
		cacheresults = jQuery.fragments[ first ];
		if ( cacheresults && cacheresults !== 1 ) {
			fragment = cacheresults;
		}
	}
  // 如果没有缓存 则先创建文档片段 再调用jQuery.clean来完善刚刚创建的文档片段
	if ( !fragment ) {
		fragment = doc.createDocumentFragment();
		jQuery.clean( args, doc, fragment, scripts );
	}

  // 如果符合条件 更改缓存的状态
  // 第一次缓存前fragments[ first ] 肯定为 undefined， 即cacheresults = undefined
  // 第一次缓存完之后 fragments[ first ] 变为 1
  // 第二次cacheresults就是1了， 那么第二次缓存的时候 fragments[ first ]就是真正的文档片段了， 
  // 第三次再去读取时 以后就读到真正的文档片段了
	if ( cacheable ) {
		jQuery.fragments[ first ] = cacheresults ? fragment : 1;
	}

	return { fragment: fragment, cacheable: cacheable };
};
```

## jQuery.clean
在buildFragment方法中被调用，  用于将字符串转为真正的dom元素。 并且在$()的那个分支 最终的选择器就有clean创建而来 `selector = ( ret.cacheable ? jQuery.clone(ret.fragment) : ret.fragment ).childNodes;`， 这个childNodes正是clean方法最后传出来的。

下面的解释可能不是非常完全，保证可以看出clean方法的过程即可。
```
	clean: function( elems, context, fragment, scripts ) {
		var checkScriptType;

		context = context || document;

		// !context.createElement fails in IE with an error but returns typeof 'object'
		if ( typeof context.createElement === "undefined" ) {
			context = context.ownerDocument || context[0] && context[0].ownerDocument || document;
		}

		var ret = [], j;

		for ( var i = 0, elem; (elem = elems[i]) != null; i++ ) {
			if ( typeof elem === "number" ) {
				elem += "";
			}

			if ( !elem ) {
				continue;
			}

			// Convert html string into DOM nodes
			if ( typeof elem === "string" ) {
				// rhtml = /<|&#?\w+;/,
				// 如果不包含标签、字符代码和数字代码 直接调用createTextNode创建文本节点。
				// createTextNode的特点是不回转义字符
				// 但innerHTML可以
				if ( !rhtml.test( elem ) ) {
					elem = context.createTextNode( elem );
				} else {
					// Fix "XHTML"-style tags in all browsers
					// rxhtmlTag = /<(?!area|br|col|embed|hr|img|input|link|meta|param)(([\w:]+)[^>]*)\/>/ig,
					// 寻找不该自闭合的标签 并矫正
					// < \w /> 并且\w不是area|br|col|embed|hr|img|input|link|meta|param
					// 但也有问题哦  <div/></div>   --> <div></div></div>
					// 只是简单的把不该自闭合的变为自闭和
					elem = elem.replace(rxhtmlTag, "<$1></$2>");

					// Trim whitespace, otherwise indexOf won't work as expected
					// rtagName = /<([\w:]+)/,
					// wrapMap = {
					// 	option: [ 1, "<select multiple='multiple'>", "</select>" ],
					// 	legend: [ 1, "<fieldset>", "</fieldset>" ],
					// 	thead: [ 1, "<table>", "</table>" ],
					// 	tr: [ 2, "<table><tbody>", "</tbody></table>" ],
					// 	td: [ 3, "<table><tbody><tr>", "</tr></tbody></table>" ],
					// 	col: [ 2, "<table><tbody></tbody><colgroup>", "</colgroup></table>" ],
					// 	area: [ 1, "<map>", "</map>" ],
					// 	_default: [ 0, "", "" ]
					// },
					// 规范标签的包裹顺序
					var tag = ( rtagName.exec( elem ) || ["", ""] )[1].toLowerCase(),
						wrap = wrapMap[ tag ] || wrapMap._default,
						depth = wrap[0],
						div = context.createElement("div");
					// ie9以下有个神奇的bug，不支持的html5标签 需要先调教一次 他才会认识 
					// 所以有了safeFragment 针对document和createSafeFragment两个方法
					// Append wrapper element to unknown element safe doc fragment
					if ( context === document ) {
						// Use the fragment we've already created for this document
						safeFragment.appendChild( div );
					} else {
						// Use a fragment created with the owner document
						createSafeFragment( context ).appendChild( div );
					}

					// 正式将正确包裹顺序的html代码 使用 innerHTML插入div元素
					div.innerHTML = wrap[1] + elem + wrap[2];
					

					// 利用之前的depth来逐步将div再次指向真正元素的父层
					// 例如option元素就只想select
					while ( depth-- ) {
						div = div.lastChild;
						// console.log(div)
					}

					// Remove IE's autoinserted <tbody> from table fragments
					// 针对ie做自动插入tbody的处理 忽略
					if ( !jQuery.support.tbody ) {

						// String was a <table>, *may* have spurious <tbody>
						var hasBody = rtbody.test(elem),
							tbody = tag === "table" && !hasBody ?
								div.firstChild && div.firstChild.childNodes :

								// String was a bare <thead> or <tfoot>
								wrap[1] === "<table>" && !hasBody ?
									div.childNodes :
									[];

						for ( j = tbody.length - 1; j >= 0 ; --j ) {
							if ( jQuery.nodeName( tbody[ j ], "tbody" ) && !tbody[ j ].childNodes.length ) {
								tbody[ j ].parentNode.removeChild( tbody[ j ] );
							}
						}
					}

					// IE completely kills leading whitespace when innerHTML is used
					// ie678自动剔除前导空白符    忽略
					if ( !jQuery.support.leadingWhitespace && rleadingWhitespace.test( elem ) ) {
						div.insertBefore( context.createTextNode( rleadingWhitespace.exec(elem)[0] ), div.firstChild );
					}
					// 这里的elem一直是一个存在且使用的变量 在for循环定义时定义 只是进入循环时 为string， 现在要变为dom数组了
					elem = div.childNodes;
					// console.log(elem)
				}
			}

			// Resets defaultChecked for any radios and checkboxes
			// about to be appended to the DOM in IE 6/7 (#8060)
			//  修正插入时 checked的bug
			// 代码很简单 只是需要知道有这个bug即可
			var len;
			if ( !jQuery.support.appendChecked ) {
				if ( elem[0] && typeof (len = elem.length) === "number" ) {
					for ( j = 0; j < len; j++ ) {
						findInputs( elem[j] );
					}
				} else {
					findInputs( elem );
				}
			}

			if ( elem.nodeType ) {
				ret.push( elem );
			} else {
				ret = jQuery.merge( ret, elem );
			}
		}
		// 如果传入了fragment， 将ret的dom元素都插入到fragment中，并且将script标签全部提取出来
		if ( fragment ) {
			checkScriptType = function( elem ) {
				return !elem.type || rscriptType.test( elem.type );
			};
			for ( i = 0; ret[i]; i++ ) {
				if ( scripts && jQuery.nodeName( ret[i], "script" ) && (!ret[i].type || ret[i].type.toLowerCase() === "text/javascript") ) {
					scripts.push( ret[i].parentNode ? ret[i].parentNode.removeChild( ret[i] ) : ret[i] );

				} else {
					if ( ret[i].nodeType === 1 ) {
						var jsTags = jQuery.grep( ret[i].getElementsByTagName( "script" ), checkScriptType );

						ret.splice.apply( ret, [i + 1, 0].concat( jsTags ) );
					}
					fragment.appendChild( ret[i] );
				}
			}
		}

		return ret;
	},
```

接下来对于jQuery(变量1)还有一个大的内部方法 就是 extend 和一个extend方法的调用。

## jQuery.extend
`jQuery.extend = jQuery.fn.extend`一样决定了这个方法在原型链上。
extend其实就是扩展的意思，将属性从一个对象拷贝到target对象上面。 这样可以无限扩展jquery对象。 当然也可以用来扩展其它对象。 只是一个思路与想法。 同时jquery绝大多数的方法与外部插件都通过这种方法来实现
```
jQuery.extend = jQuery.fn.extend = function() {
		// 这里是默认配置
    var options, name, src, copy, copyIsArray, clone,
      target = arguments[0] || {},
      i = 1,
      length = arguments.length,
      deep = false;

    // 如果arguments[0]是一个布尔类型 则arguments[0] 表示deep变量， arguments[1]则顺位表示target
    if ( typeof target === "boolean" ) {
      deep = target;
      target = arguments[1] || {};
      // skip the boolean and the target
      i = 2;
    }

    // Handle case when target is a string or something (possible in deep copy)
		// 处理target不是{}的情况，默认一个空对象给target
    if ( typeof target !== "object" && !jQuery.isFunction(target) ) {
      target = {};
    }

    // extend jQuery itself if only one argument is passed
		// 如果发现arguments只有一个元素， 那么target就是this，
		// 一般而言 this无非jQuery 或者 jQuery.fn 
    if ( length === i ) {
      target = this;
      --i;
    }

		// 真正的主体是一个循环，用于处理target之后的arguments元素
    for ( ; i < length; i++ ) {
      // 不处理arguments[i]为null 或 undefined的情况
			// null == undefined --> true
			// null != undefined --> false
			// null == false --> false
			// null != false --> true
      if ( (options = arguments[ i ]) != null ) {
        // Extend the base object
				// 遍历options对象或数组的key
        for ( name in options ) {
					// src表示target上同名的值
					// copy表示options上的值
          src = target[ name ];
          copy = options[ name ];

					// 如果options上的值 就是target本身，那么再把这个值赋给target没意义 且会报错。
          // Prevent never-ending loop
          if ( target === copy ) {
            continue;
          }

          // Recurse if we're merging plain objects or arrays
          if ( deep && copy && ( jQuery.isPlainObject(copy) || (copyIsArray = jQuery.isArray(copy)) ) ) {
            if ( copyIsArray ) {
              copyIsArray = false;
              clone = src && jQuery.isArray(src) ? src : [];

            } else {
              clone = src && jQuery.isPlainObject(src) ? src : {};
            }

            // Never move original objects, clone them
						// 如果copy是数组或者对象 并且deep为true则递归的去拷贝
            target[ name ] = jQuery.extend( deep, clone, copy );

          // Don't bring in undefined values
          } else if ( copy !== undefined ) {
            // 只要不是深deep拷贝， 那么直接用copy覆盖或设置target同名属性即可。
						target[ name ] = copy;
          }
        }
      }
    }

    // Return the modified object
    return target;
  };
```
## jQuery.prototype 除init外的其余方法和属性

### selector、jquery、length、size、toArray、get
```
// 选择器 无实际意义，并不一定是可执行代码
selector: "",
// 当前jquery版本号
jquery: "1.7.1",
// 当前jQuery对象中元素个数
length: 0,
// 获取当前jQuery对象中元素个数的方法， 应该直接访问length，减少函数调用开销
size: function() {
	return this.length;
},
// 将jQuery对象的元素转为真正的数组
// jquery整个代码中有一点值得注意就是 this对象上面一只有一个length属性在动态的更改。且length值和[0],[1],[2]...息息相关
// 我们自己也可以实验，在对象上设置length值，然后调用Array.prototype.slice.call也可以获得length长度的数组	
toArray: function() {
	// slice = Array.prototype.slice,
	return slice.call( this, 0 );
},

// $('body').get(0)和 $('body')[0]其实是一样的
get: function( num ) {
	return num == null ? this.toArray() : ( num < 0 ? this[ this.length + num ] : this[ num ] );
},
```
看到这里需要看一下this到底是一个什么东东

找个地方随便打印一下this
会发现this有这么几个东西
```
// 以$('body')为例
this : {
	0 : body,
	context : document,
	length : 1,
	selector : "body",
}
```
this这鬼东西竟然是一个对象，但我们平时感觉他好像还有数组的特性。 this是一个对象这并不难理解，常识中this就应该是对象，但jquery中this的巧妙之处在于jquery一直在维护一个length属性，同时我们真正关心的对象元素的key都是0、1、2、3、4、5这样的连贯数字， 这样就让我们可以直接$('body')[0]来获取一个dom元素。 当然return this也实现了链式操作。 jquery设计的确实高明

### pushStack
```
pushStack: function( elems, name, selector ) {
	// 创建一个新的空的jquery对象
	var ret = this.constructor();
	// push = Array.prototype.push
	if ( jQuery.isArray( elems ) ) {
		// Array.prototype.push.apply( ret, elems );
		// 实现了将数组elems的对象添加到对象ret上面
		// ret是jquery的一个this值而已。
		// 与这句话相反的是从this到对象  toArray方法
		push.apply( ret, elems );
	} else {
		// 可以猜到这个merge方法肯定也是处理类似的事情，只是数组比较简单，如果是对象可能需要for遍历等
		jQuery.merge( ret, elems );
	}

	// 需要将当前的jquery this指针保存到新建德ret对象上， 同时把执行上下文赋给新的jquery对象
	ret.prevObject = this;
	ret.context = this.context;
	// 根据调用方式的不同，来更改一下ret的选择器selector记录
	if ( name === "find" ) {
		ret.selector = this.selector + ( this.selector ? " " : "" ) + selector;
	} else if ( name ) {
		ret.selector = this.selector + "." + name + "(" + selector + ")";
	}

	// 这一次返回了新的jquery对象，新的this。但是调用者感觉不到，因为jquery所有的方法都基本在prototype上， 无论新建多少个jquery实例都可以共享，只需要保证每次return的是jquery的this就可以保证链式操作
	return ret;
},
```
### end
在pushStack中会将原本存在的元素对象指针暂存到新的jQuery对象的prevObject属性。
end方法就是入栈的相反操作，推出栈顶部。
如果没有prevObject属性，就新建一个全空的全新的jQuery对象。
```
end: function() {
	return this.prevObject || this.constructor(null);
},
```
### each
each顾名思义也就是类似于遍历的概念，方法接受callback和args(回调函数和参数)
```
// 发现each调用jQuery.each？ 本身调用本身？仔细想一下其实不是的。
// 外层的这个each是jQuery.prototype.each(jQuery.fn.each)
// 内层的事jQuery.each
// 这里其实区别挺大的。外层的each其实是new jQuery.fn.init()实例来的this 当他身上没有挂载一个方法或者属性时他会寻找他的原型链。 也就是说$().each先寻找init方法上面each，失败，继续寻找 init.prototype.each就找到了 因为jQuery.fn.init.prototype = jQuery.fn;
// 内层的each其实是一个静态的方法。内层each的这个jQuery其实就是执行这句话的return new jQuery.fn.init( selector, context, rootjQuery )那个变量。但那个变量没有执行变为函数，也就是说没有this这个概念。 其实可以找到这个each是通过jQuery.extend({})挂载到jQuery身上的。 当然要访问也就是jQuery.each 相当于一个静态方法
// 例子
// let a = function(){ return this }
// a.prototype.each = '123'
// a.each = '456'
// a.each  --> 456
// (new a()).each    --> 123 靠原型链new时传递来实现的
each: function( callback, args ) {
	return jQuery.each( this, callback, args );
},
```
因此问题变成了 jQuery.each
jQuery.each就是为传入的数组或对象每个元素执行一遍回调函数而已， 思路想一下Array.prototype.map 和 Array.prototype.each即可
```
each: function( object, callback, args ) {
	var name, i = 0,
		length = object.length,
		// 判断object是不是有length属性 即是不是类数组对象或就是数组。
		// 类数组对象或者数组则为false， 数组等为true
		// 类数组对象 arguments是一个例子
		isObj = length === undefined || jQuery.isFunction( object );
		// 有没有传入参数 分为if else
		// 可以看出内部的循环便利每次出现false 则停止遍历break。 区别于continue 
	if ( args ) {
		// 是不是类数组或数组又分为 if else
		if ( isObj ) {
			// 非类数组和数组 用for in来遍历object
			for ( name in object ) {
				if ( callback.apply( object[ name ], args ) === false ) {
					break;
				}
			}
		} else {
			// 数组和类数组则直接用 下标来遍历
			for ( ; i < length; ) {
				if ( callback.apply( object[ i++ ], args ) === false ) {
					break;
				}
			}
		}

	// A special, fast, case for the most common use of each
	} else {
		if ( isObj ) {
			for ( name in object ) {
				if ( callback.call( object[ name ], name, object[ name ] ) === false ) {
					break;
				}
			}
		} else {
			for ( ; i < length; ) {
				// 这里传入的第一个参数和第三个参数是一样的，因为 i++ 是要等这个运算结束以后才会累加
				// 区别在于 i++ ++i， so easy
				if ( callback.call( object[ i ], i, object[ i++ ] ) === false ) {
					break;
				}
			}
		}
	}

	return object;
},
```
### ready 比较复杂 暂时不分析

### eq、first、last、slice
比较巧的是first、last调用的是eq，eq调用的是slice，因此先看slice，在看eq，最后看first和last
```
// slice这段调用的是pushStack， 而slice做的就是准备数据而已， 包括elems： slice.apply( this, arguments )， name： "slice"， selector：slice.call(arguments).join(",");
// 这段代码依旧利用了Array.prototype.slice 数组原型方法的方式来从类数组对象抽取数据
slice: function() {
	// slice = Array.prototype.slice,
	return this.pushStack( slice.apply( this, arguments ),
		"slice", slice.call(arguments).join(",") );
},
// eq就是从this中抽出指定的唯一下标的一个dom元素(其实不一定是dom元素，但是为了方便分析，毕竟jqeury获取的大部分都是dom元素)
eq: function( i ) {
	// 利用自加可以快速实现 字符串向阿拉伯数字转变
	i = +i;
	// 其实eq设定只能获取一个元素，而且调用的是slice方法，所以需要尽可能的为slice指定两个参数，否则会造成slice或得到一个数组
	return i === -1 ?
		this.slice( i ) :
		this.slice( i, i + 1 );
},

first: function() {
	return this.eq( 0 );
},

last: function() {
	return this.eq( -1 );
},
```
### map
又与each非常相似，内部调用了一个jQuery的静态同名方法map。 但这次还调用了pushStack这样一个入栈的方法。
可能与Array.prototype.map和Array.prototype.each类似，map不改变原值，需要返回新值，each直接在原值上更改，所以map需要新建一个jQuery对象，同时在prevObject保留了上一次的对象
```
map: function( callback ) {
	return this.pushStack( jQuery.map(this, function( elem, i ) {
		return callback.call( elem, i, elem );
	}));
},

// 这里有趣的一点是each和map里面的判断都变了，一个是判断isObj，一个是Array，而且对于缺省arg也不尽兴判断了。可能两段代码不是出自同一人之手。但基本效果是一致的
// 这里和each区别最大的就是 这里新建了一个ret空数组， 然后将遍历callback返回的数据都压入ret数组内，最后整个map返回的是扁平化之后的ret 组成的新数组。 这里新建的ret是空数组而不是jQuery空对象 可能是考虑到外层是pushStack调用，同时不应该直接修改，应该保存原有的到prevObject。
map: function( elems, callback, arg ) {
	var value, key, ret = [],
		i = 0,
		length = elems.length,
		// jquery objects are treated as arrays
		isArray = elems instanceof jQuery || length !== undefined && typeof length === "number" && ( ( length > 0 && elems[ 0 ] && elems[ length -1 ] ) || length === 0 || jQuery.isArray( elems ) ) ;

	// Go through the array, translating each of the items to their
	if ( isArray ) {
		for ( ; i < length; i++ ) {
			value = callback( elems[ i ], i, arg );

			if ( value != null ) {
				ret[ ret.length ] = value;
			}
		}

	// Go through every key on the object,
	} else {
		for ( key in elems ) {
			value = callback( elems[ key ], key, arg );

			if ( value != null ) {
				ret[ ret.length ] = value;
			}
		}
	}

	// Flatten any nested arrays
	return ret.concat.apply( [], ret );
},
```

## jQuery.extend调用 389-893行。
基于前面定义的jQuery.extend方法来扩展jQuery对象。

### noConflict 
```
noConflict: function( deep ) {
	// 这两歌变量是在整个jQuery执行前定义的，用于保存现有的window.jQuery和window.$同名对象。
	// _jQuery = window.jQuery,
	// _$ = window.$,
	// 当调用noConflict时 会把$ 和  jQuery的控制权重新移交给之前使用的库或者插件等。
	if ( window.$ === jQuery ) {
		window.$ = _$;
	}

	if ( deep && window.jQuery === jQuery ) {
		window.jQuery = _jQuery;
	}

	return jQuery;
},
```

### isReady、readyWait、holdReady、ready、bindReady 
这5个属性或者方法都是用于处理ready立即执行函数的。 后面再分析吧，先捡软柿子捏。

### isFunction、isArray、isWindow、isNumeric、type
可以看出isFunction、isArray都是依赖于type，isWindow、isNumeric来个方法比较纯粹。
```
// type内部先将undefined 和 null类型变量转为字符串形式 undefined 和 null，其余类型读取class2type， 如果读取不到则默认为object格式
type: function( obj ) {
	return obj == null ?
		String( obj ) :
		class2type[ toString.call(obj) ] || "object";
},
// class2type对象构建方式 利用jQuery.each 来构建
// 构建结果
<!--{
	[object Array] : "array"
	[object Boolean] : "boolean"
	[object Date] : "date"
	[object Function] : "function"
	[object Number] : "number"
	[object Object] : "object"
	[object RegExp] : "regexp"
	[object String] : "string
}-->
jQuery.each("Boolean Number String Function Array Date RegExp Object".split(" "), function(i, name) {
	class2type[ "[object " + name + "]" ] = name.toLowerCase();
});
isFunction: function( obj ) {
	return jQuery.type(obj) === "function";
},
isArray: Array.isArray || function( obj ) {
	return jQuery.type(obj) === "array";
},
// 如果是object， 则判断setInterval有没有在对象里， 如果在则就是window了， 利用一个小技巧曲线救国
// 当然如果一个对象也实现了setInterval方法，可能这个检测就失效了。 但谁会没事实现或命名这样一个属性或变量呢
isWindow: function( obj ) {
	return obj && typeof obj === "object" && "setInterval" in obj;
},
// 首先利用parseFloat来转换为浮点数，如果转换失败会返回NaN，证明这不是一个可转换的数字
// isNaN来判断 如果是NaN那么就证明这不是一个数字或可转换为数字
// isFinite来确定不是一个无限大或小的数
// NaN的判断比较奇怪， 如前比较稳定的判断方法是isNaN和Object.is()， 可以看另一片blog https://github.com/AllFE/blog/issues/6
isNumeric: function( obj ) {
	return !isNaN( parseFloat(obj) ) && isFinite( obj );
},
```

### isPlainObject、isEmptyObject
```
// 用于判断目标是不是一个真正的js纯对象。不可以是函数、dom对象等等
isPlainObject: function( obj ) {
	// 目标首先存在 并且jquery的类型检测是object。 在jquery类型检测中，一切不归于Boolean Number String Function Array Date RegExp 都算是object 因此最基本的还需要排除dom元素， 当认为nodeType就算dom元素，同时不可一世window对象
	但是这样 {nodeType: 'obj'}也会检测失败
	if ( !obj || jQuery.type(obj) !== "object" || obj.nodeType || jQuery.isWindow( obj ) ) {
		return false;
	}

	try {
		// 约定一个对象如果直接有constructor属性则必定不是纯Object， 但是constructor应该出现在其原型链上
		if ( obj.constructor &&
			!hasOwn.call(obj, "constructor") &&
			!hasOwn.call(obj.constructor.prototype, "isPrototypeOf") ) {
			return false;
		}
	} catch ( e ) {
		// IE8,9 Will throw exceptions on certain host objects #9897
		// 由于某些ie8、9原因需要用try--catch。 原因未知
		return false;
	}
	// 遍历对象的属性，如果是空对象则key最后为undefined，否则最后判断key应该是obj的自由属性，否则就是不是纯object
	var key;
	for ( key in obj ) {}
	return key === undefined || hasOwn.call( obj, key );
},
// 判断是不是空对象
isEmptyObject: function( obj ) {
	// 遍历 只要有name就报错
	for ( var name in obj ) {
		return false;
	}
	return true;
},
```

### error
```
// 快速的创建error实例
error: function( msg ) {
	throw new Error( msg );
},
```

### parseJSON、parseXML
```
parseJSON: function( data ) {
	// 传入数据不是字符串 则返回 null
	if ( typeof data !== "string" || !data ) {
		return null;
	}

	// 需要字符串先进行前后的空格去除
	data = jQuery.trim( data );

	// 尽可能的使用自带的JSON.parse功能去转译
	if ( window.JSON && window.JSON.parse ) {
		return window.JSON.parse( data );
	}

	// 验证json字符串的正确性 确保正确再进行转换 利用http://json.org/json2.js的思想逻辑
	// http://json.org/json-zh.html 详细介绍了json的规定。
	// 利用rvalidescape、rvalidtokens、rvalidbraces将一个合法的json字符串转换为仅剩\,:{}\s的字符串
	// rvalidchars = /^[\],:{}\s]*$/,
  // rvalidescape = /\\(?:["\\\/bfnrt]|u[0-9a-fA-F]{4})/g, "{\\r'qwer': '\\u2044\\b'}".replace(rvalidescape, '@') --> "{@'qwer': '@@'}"   主要解决转译码的问题
  // rvalidtokens = /"[^"\\\n\r]*"|true|false|null|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?/g,
	// rvalidtokens主要替换 ""包裹的字符串(内部不能有\n或者\r或者"存在) true false null 数字(数字比较麻烦，有小数还有科学技术)
  // rvalidbraces = /(?:^|:|,)(?:\s*\[)+/g,
	// rvalidbraces主要替换':   ['.replace(reg, '@') 开方括号 或者首位置的方括号
	// 经过三个替换，最后应该只剩下\,:{}\s的字符串 则用rvalidchars进行测试， 测试通过则就算是正确的json， 然后利用( new Function( "return " + data ) )()来实现json转换
	if ( rvalidchars.test( data.replace( rvalidescape, "@" )
		.replace( rvalidtokens, "]" )
		.replace( rvalidbraces, "")) ) {

		return ( new Function( "return " + data ) )();

	}
	// 走到这里就说明不是合法的json对象
	jQuery.error( "Invalid JSON: " + data );
},

parseXML: function( data ) {
	// 利用跨平台w3c标准和ie来进行xml解析
	var xml, tmp;
	try {
		if ( window.DOMParser ) { // Standard
			tmp = new DOMParser();
			xml = tmp.parseFromString( data , "text/xml" );
		} else { // IE
			xml = new ActiveXObject( "Microsoft.XMLDOM" );
			xml.async = "false";
			xml.loadXML( data );
		}
	} catch( e ) {
		xml = undefined;
	}
	// 报错
	if ( !xml || !xml.documentElement || xml.getElementsByTagName( "parsererror" ).length ) {
		jQuery.error( "Invalid XML: " + data );
	}
	return xml;
},
```

###  noop
空函数
```
noop: function() {},
```

### globalEval
```
globalEval: function( data ) {
	// rnotwhite = /\S/  保证非空
	if ( data && rnotwhite.test( data ) ) {
		// ie直接使用execScript或者使用匿名函数内部eval
		// 为什么使用匿名函数？
		// eval执行上下文只针对当前调用的环境， 但是匿名函数是直接挂在全局的，因此匿名函数内部的eval也就是全局的了
		// execScript直接就是全局
		( window.execScript || function( data ) {
			window[ "eval" ].call( window, data );
		} )( data );
	}
},
```

### camelCase
实现驼峰转换
```
// rmsPrefix = /^-ms-/,  如果以-ms-开头的则替换为ms-
// rdashAlpha = /-([a-z]|[0-9])/ig,  将以-开头的都替换为大写的并且去掉-
fcamelCase = function( all, letter ) {
	return ( letter + "" ).toUpperCase();
},
camelCase: function( string ) {
	return string.replace( rmsPrefix, "ms-" ).replace( rdashAlpha, fcamelCase );
},
```

### nodeName
判断nodeName
```
nodeName: function( elem, name ) {
	return elem.nodeName && elem.nodeName.toUpperCase() === name.toUpperCase();
},
```

### trim
```
// String.prototype.trim
// 下一行第二个trim是一个已经定义的变量
// 如果trim存在 直接调用浏览器的trim
// 否则用正则匹配掉左右的空白
// trimLeft = /^\s+/,
// trimRight = /\s+$/,

trim: trim ?
	function( text ) {
		return text == null ?
			"" :
			trim.call( text );
	} :
	function( text ) {
		return text == null ?
			"" :
			text.toString().replace( trimLeft, "" ).replace( trimRight, "" );
	},
```

### makeArray
将数据压入一个数组
```
makeArray: function( array, results ) {
	// 如果没有传入结果数组就新建
	var ret = results || [];
	// 当数据不为undefined、null时才进入，意味着false也进入
	if ( array != null ) {
		// The window, strings (and functions) also have 'length'
		// Tweaked logic slightly to handle Blackberry 4.7 RegExp issues #6930
		var type = jQuery.type( array );
		// 如果是简单类型直接push进数组
		// 否则调用merge进行数组拼接组合
		if ( array.length == null || type === "string" || type === "function" || type === "regexp" || jQuery.isWindow( array ) ) {
			push.call( ret, array );
		} else {
			jQuery.merge( ret, array );
		}
	}

	return ret;
},
```
### merge
主要用于拼接两个jQuery对象 或者将一个数组内容赋予给一个jQuery对象
jQuery对象有个特点就是有length属性，并且有从0开始的紧密下标
```
merge: function( first, second ) {
	var i = first.length,
		j = 0;
	// second是array的情况
	if ( typeof second.length === "number" ) {
		for ( var l = second.length; j < l; j++ ) {
			first[ i++ ] = second[ j ];
		}
	} else {
	// second是jQuery对象的情况
		while ( second[j] !== undefined ) {
			first[ i++ ] = second[ j++ ];
		}
	}
	first.length = i;
	return first;
},
```

### grep
过滤掉数组中满足条件的元素，并将其与不满足条件的返回，并且不修改原数组
```
grep: function( elems, callback, inv ) {
	// 定义一个ret新数组， 通过!!将inv一定能转为true|false
	var ret = [], retVal;
	inv = !!inv;

	// 同样的用!!来获得callback的执行结果true|false
	// 如果输出结果与inv相同则过滤掉 否则添加进ret
	for ( var i = 0, length = elems.length; i < length; i++ ) {
		retVal = !!callback( elems[ i ], i );
		if ( inv !== retVal ) {
			ret.push( elems[ i ] );
		}
	}

	return ret;
},
```

### proxy
返回一个拥有特定上下文的新函数，有点类似于bind
有两种形式 proxy(func, context)、 proxy(context, name)
第一种比较正常
第二种意味着指定context上name属性对应的函数上下文一定为context
```
proxy: function( fn, context ) {
	// 如果context为string则是第二种情况，需要修正fn和context
	if ( typeof context === "string" ) {
		var tmp = fn[ context ];
		context = fn;
		fn = tmp;
	}

	// 如果传入的函数有问题，不是可执行的返回undefined
	if ( !jQuery.isFunction( fn ) ) {
		return undefined;
	}

	// 模拟bind
	// 先将现在传入的参数提起出来保存到args中
	// 然后生成返回的更改了执行上下文的函数proxy
	// proxy内部又返回了一个新的函数，当新proxy执行的时候 内部新建的函数运动了apply
	// 1、直接使用apply是不可以的，因为apply是立即执行的，不符合需求
	// 2、新proxy内部既可以利用现在传入的参数，也可以利用真正执行时传入的新参数
	var args = slice.call( arguments, 2 ),
		proxy = function() {
			return fn.apply( context, args.concat( slice.call( arguments ) ) );
		};

	// 更改函数的guid，暂时不知道guid具体作用
	// 并且proxy.guid(|| jQuery.guid++;左边的)个人认为不会取到， 一直是undefined， 因为proxy才刚刚定义。。。。。不知道jquery怎么想的
	proxy.guid = fn.guid = fn.guid || proxy.guid || jQuery.guid++;

	return proxy;
},
```

### access
为集合elems中的元素设置一个或多个属性 或读第一个元素的属性值
```
access: function( elems, key, value, exec, fn, pass ) {
	var length = elems.length;

	// 如果key是一个对象 就是多个属性的时候 遍历设置
	if ( typeof key === "object" ) {
		for ( var k in key ) {
			jQuery.access( elems, k, key[k], exec, fn, value );
		}
		return elems;
	}

	// 当value有值时设置属性
	if ( value !== undefined ) {
		// value是函数 并且 exec为true， pass为false时
		exec = !pass && exec && jQuery.isFunction(value);
		// 遍历元素集合 为每一个元素设置key属性值，
		// 传入的第三个值由exec决定，可能是一个函数计算过的值 也可能是直接静态值
		// 如果是函数利用call来将执行上下文绑定到elems[i]上
		for ( var i = 0; i < length; i++ ) {
			fn( elems[i], key, exec ? value.call( elems[i], i, fn( elems[i], key ) ) : value, pass );
		}

		return elems;
	}

	// 执行这里说明是读取第一个元素的key属性值或者返回undefined
	return length ? fn( elems[0], key ) : undefined;
},
```

### now
返回当前时间的时间戳， 一个简写吧
```
now: function() {
	return ( new Date() ).getTime();
},
```

### uaMatch
获取当前浏览器的内核和版本号
```
// rwebkit = /(webkit)[ \/]([\w.]+)/,
// ropera = /(opera)(?:.*version)?[ \/]([\w.]+)/,
// rmsie = /(msie) ([\w.]+)/,
// rmozilla = /(mozilla)(?:.*? rv:([\w.]+))?/,

uaMatch: function( ua ) {
	ua = ua.toLowerCase();

	var match = rwebkit.exec( ua ) ||
		ropera.exec( ua ) ||
		rmsie.exec( ua ) ||
		ua.indexOf("compatible") < 0 && rmozilla.exec( ua ) ||
		[];

	return { browser: match[1] || "", version: match[2] || "0" };
},
```
配合一起看
// 假设navigator.userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
```
browserMatch = jQuery.uaMatch( userAgent );
if ( browserMatch.browser ) {
	jQuery.browser[ browserMatch.browser ] = true;
	jQuery.browser.version = browserMatch.version;
}

// safari也是webkit内核
if ( jQuery.browser.webkit ) {
	jQuery.browser.safari = true;
}
```
### sub
生成另一个jQuery备份， 但是这个不影响原有jQuery。
```
sub: function() {
	function jQuerySub( selector, context ) {
		return new jQuerySub.fn.init( selector, context );
	}
	jQuery.extend( true, jQuerySub, this );
	jQuerySub.superclass = this;
	jQuerySub.fn = jQuerySub.prototype = this();
	jQuerySub.fn.constructor = jQuerySub;
	jQuerySub.sub = this.sub;
	jQuerySub.fn.init = function init( selector, context ) {
		if ( context && context instanceof jQuery && !(context instanceof jQuerySub) ) {
			context = jQuerySub( context );
		}

		return jQuery.fn.init.call( this, selector, context, rootjQuerySub );
	};
	jQuerySub.fn.init.prototype = jQuerySub.fn;
	var rootjQuerySub = jQuerySub(document);
	return jQuerySub;
},
```
----
未完待续