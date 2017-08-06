# chapter10 DOM操作

对DOM的操作在jQuery中有插入、删除、复制、替换、包裹元素五类，基于原生的insertBefore、appendChild、removeChild、cloneNode4个方法，其中replaceChild没有用到，替换元素则是基于已经实现的删除与插入来完成。

DOM操作部分代码并不是特别复杂，只有插入那块用到了之前分析过的一两段复杂的构建html的代码，复制部分需要对兼容性修正，其余部分都是利用原生方法封装为更方便易用的接口。

``` js
// 这段代码的方法是jQuery里面最基本的方法了。很熟悉的感觉
jQuery.fn.extend({
  text: function( text ) {
    // text为函数， 则each调用函数之后的值
    if ( jQuery.isFunction(text) ) {
      return this.each(function(i) {
        var self = jQuery( this );

        self.text( text.call(this, i, self.text()) );
      });
    }

    // 为元素插入文本节点
    if ( typeof text !== "object" && text !== undefined ) {
      return this.empty().append( (this[0] && this[0].ownerDocument || document).createTextNode( text ) );
    }
    // 返回第一个元素的文本信息
    return jQuery.text( this );
  },

  wrapAll: function( html ) {
    // 传入方法则为每个元素用方法执行后的结果执行warpAll
    if ( jQuery.isFunction( html ) ) {
      return this.each(function(i) {
        jQuery(this).wrapAll( html.call(this, i) );
      });
    }

    if ( this[0] ) {
      // The elements to wrap the target around
      // 先获取包裹元素的副本
      var wrap = jQuery( html, this[0].ownerDocument ).eq(0).clone(true);

      // 如果有第一个匹配元素有父元素，则将包裹元素插入到第一个元素的前面
      if ( this[0].parentNode ) {
        wrap.insertBefore( this[0] );
      }

      // 找到wrap最内层的元素，然后把当前元素插进去
      wrap.map(function() {
        var elem = this;

        while ( elem.firstChild && elem.firstChild.nodeType === 1 ) {
          elem = elem.firstChild;
        }

        return elem;
      }).append( this );
    }

    return this;
  },

  wrapInner: function( html ) {
    // 累
    if ( jQuery.isFunction( html ) ) {
      return this.each(function(i) {
        jQuery(this).wrapInner( html.call(this, i) );
      });
    }
    // 遍历当前元素
    // 如果当前元素有子元素则子元素调用wrapAll
    // 否则直接把wrap插入当前元素
    return this.each(function() {
      var self = jQuery( this ),
        contents = self.contents();

      if ( contents.length ) {
        contents.wrapAll( html );

      } else {
        self.append( html );
      }
    });
  },

  wrap: function( html ) {
    var isFunction = jQuery.isFunction( html );
    // 在匹配元素的每个元素执行wrapAll
    return this.each(function(i) {
      jQuery( this ).wrapAll( isFunction ? html.call(this, i) : html );
    });
  },

  unwrap: function() {
    // 为每个匹配元素的父元素执行replaceWidth， 用子元素直接替代它，就相当于解除包裹了
    return this.parent().each(function() {
      if ( !jQuery.nodeName( this, "body" ) ) {
        jQuery( this ).replaceWith( this.childNodes );
      }
    }).end();
  },

  append: function() {
    // 调用domManip完成dom的转换，然后再执行appendChild完成插入操作
    return this.domManip(arguments, true, function( elem ) {
      if ( this.nodeType === 1 ) {
        this.appendChild( elem );
      }
    });
  },

  prepend: function() {
    // 完成在头部插入
    return this.domManip(arguments, true, function( elem ) {
      if ( this.nodeType === 1 ) {
        this.insertBefore( elem, this.firstChild );
      }
    });
  },

  before: function() {
    // 在指定元素前面插入元素
    if ( this[0] && this[0].parentNode ) {
      return this.domManip(arguments, false, function( elem ) {
        this.parentNode.insertBefore( elem, this );
      });
    } else if ( arguments.length ) {
      var set = jQuery.clean( arguments );
      set.push.apply( set, this.toArray() );
      return this.pushStack( set, "before", arguments );
    }
  },

  after: function() {
    // 在指定元素的后方插入元素
    if ( this[0] && this[0].parentNode ) {
      return this.domManip(arguments, false, function( elem ) {
        this.parentNode.insertBefore( elem, this.nextSibling );
      });
    } else if ( arguments.length ) {
      var set = this.pushStack( this, "after", arguments );
      set.push.apply( set, jQuery.clean(arguments) );
      return set;
    }
  },

  // keepData is for internal use only--do not document
  // 先移除后代元素和关联的数据和事件，以防止内存泄漏
  remove: function( selector, keepData ) {
    for ( var i = 0, elem; (elem = this[i]) != null; i++ ) {
      if ( !selector || jQuery.filter( selector, [ elem ] ).length ) {
        if ( !keepData && elem.nodeType === 1 ) {
          jQuery.cleanData( elem.getElementsByTagName("*") );
          jQuery.cleanData( [ elem ] );
        }

        if ( elem.parentNode ) {
          elem.parentNode.removeChild( elem );
        }
      }
    }

    return this;
  },

  // 移除所有后代
  // 也要先移除数据和事件
  // 不过最后有个监测措施
  empty: function() {
    for ( var i = 0, elem; (elem = this[i]) != null; i++ ) {
      // Remove element nodes and prevent memory leaks
      if ( elem.nodeType === 1 ) {
        jQuery.cleanData( elem.getElementsByTagName("*") );
      }

      // Remove any remaining nodes
      while ( elem.firstChild ) {
        elem.removeChild( elem.firstChild );
      }
    }

    return this;
  },
  // 负责决定两个参数，是否复制元素本身的数据和事件
  // 是否复制后代的元素和事件
  clone: function( dataAndEvents, deepDataAndEvents ) {
    dataAndEvents = dataAndEvents == null ? false : dataAndEvents;
    deepDataAndEvents = deepDataAndEvents == null ? dataAndEvents : deepDataAndEvents;

    return this.map( function () {
      return jQuery.clone( this, dataAndEvents, deepDataAndEvents );
    });
  },

  // rinlinejQuery = / jQuery\d+="(?:\d+|null)"/g,
  html: function( value ) {
    // 如果没有传入值，那么表示读取第一个元素innerHTML
    if ( value === undefined ) {
      return this[0] && this[0].nodeType === 1 ?
        this[0].innerHTML.replace(rinlinejQuery, "") :
        null;

    // See if we can take a shortcut and just use innerHTML
    // 在不需要任何修正的情况下
    } else if ( typeof value === "string" && !rnoInnerhtml.test( value ) &&
      (jQuery.support.leadingWhitespace || !rleadingWhitespace.test( value )) &&
      !wrapMap[ (rtagName.exec( value ) || ["", ""])[1].toLowerCase() ] ) {

      value = value.replace(rxhtmlTag, "<$1></$2>");

      try {
        for ( var i = 0, l = this.length; i < l; i++ ) {
          // Remove element nodes and prevent memory leaks
          // 遍历当前所有匹配元素，尝试把内部的所有元素的数据缓存都删掉，然后重设他们的innerHTML
          if ( this[i].nodeType === 1 ) {
            jQuery.cleanData( this[i].getElementsByTagName("*") );
            this[i].innerHTML = value;
          }
        }

      // If using innerHTML throws an exception, use the fallback method
      } catch(e) {
        // 前面出了问题，直接调用empty方法，然后再插入新值
        this.empty().append( value );
      }

    } else if ( jQuery.isFunction( value ) ) {
      // 如果发现value是函数，则each为每个匹配元素调用html方法
      this.each(function(i){
        var self = jQuery( this );

        self.html( value.call(this, i, self.html()) );
      });

    } else {
      this.empty().append( value );
    }

    return this;
  },
  // 替换元素
  replaceWith: function( value ) {
    // 当前有元素且元素有父元素
    if ( this[0] && this[0].parentNode ) {
      // Make sure that the elements are removed from the DOM before they are inserted
      // this can help fix replacing a parent with child elements
      // 如果value是函数 则便利当前匹配元素，继续用value执行后的结果调用replaceWith
      if ( jQuery.isFunction( value ) ) {
        return this.each(function(i) {
          var self = jQuery(this), old = self.html();
          self.replaceWith( value.call( this, i, old ) );
        });
      }
      // 如果value不是方法、不是字符串可能是dom活着jQuery对象，则先把value从文档中移除掉
      if ( typeof value !== "string" ) {
        value = jQuery( value ).detach();
      }
      
      // 先把元素移除掉，再用value插入、可能是通过后面兄弟前面插入、也可能是通过父亲直接最后插入
      return this.each(function() {
        var next = this.nextSibling,
          parent = this.parentNode;

        jQuery( this ).remove();

        if ( next ) {
          jQuery(next).before( value );
        } else {
          jQuery(parent).append( value );
        }
      });
    } else {
      // 有length表示有匹配元素，但第一个元素没有父元素，执行this.pushStack( jQuery(jQuery.isFunction(value) ? value() : value), "replaceWith", value )
      // 否则直接意味着当前没有匹配元素 直接返回this本身
      // 第一种情况，直接拿value构建一个新的jQuery对象返回
      return this.length ?
        this.pushStack( jQuery(jQuery.isFunction(value) ? value() : value), "replaceWith", value ) :
        this;
    }
  },

  detach: function( selector ) {
    return this.remove( selector, true );
  },

  // 被多个插入方法调用的基本方法
  // 就是转换dom元素，并且调用真正的回调函数，把转换后的dom元素插入进去
  // checked="checked" or checked
  // rchecked = /checked\s*(?:[^=]|=\s*.checked.)/i, 
  domManip: function( args, table, callback ) {
    var results, first, fragment, parent,
      value = args[0],
      scripts = [];

    // 浏览器能否正常复制含有checked参数的dom元素
    // 迭代执行来解决这个问题
    if ( !jQuery.support.checkClone && arguments.length === 3 && typeof value === "string" && rchecked.test( value ) ) {
      return this.each(function() {
        jQuery(this).domManip( args, table, callback, true );
      });
    }

    // 如果发现args数组内是函数，则遍历当前匹配的元素集合，将执行后的数据继续执行domManip
    if ( jQuery.isFunction(value) ) {
      return this.each(function(i) {
        var self = jQuery(this);
        args[0] = value.call(this, i, table ? self.html() : undefined);
        self.domManip( args, table, callback );
      });
    }

    if ( this[0] ) {
      parent = value && value.parentNode;

      // If we're in a fragment, just use that instead of building a new one
      // 这里有个坑， jQuery.support没有测试parentNode， 所以这个if永远不会执行
      if ( jQuery.support.parentNode && parent && parent.nodeType === 11 && parent.childNodes.length === this.length ) {
        results = { fragment: parent };

      } else {
        // 调用buildFragment将字符串转为dom元素
        // 顺便把script提取出来了
        results = jQuery.buildFragment( args, this, scripts );
      }

      fragment = results.fragment;

      if ( fragment.childNodes.length === 1 ) {
        first = fragment = fragment.firstChild;
      } else {
        first = fragment.firstChild;
      }

      if ( first ) {
        // 查看当前元素是不是tr
        table = table && jQuery.nodeName( first, "tr" );

        for ( var i = 0, l = this.length, lastIndex = l - 1; i < l; i++ ) {
          callback.call(
            table ?
              root(this[i], first) :
              this[i],
            // Make sure that we do not leak memory by inadvertently discarding
            // the original fragment (which might have attached data) instead of
            // using it; in addition, use the original fragment object for the last
            // item instead of first because it can end up being emptied incorrectly
            // in certain situations (Bug #8070).
            // Fragments from the fragment cache must always be cloned and never used
            // in place.
            // 如果该dom是可缓存的，则总是插入它的副本
            // 如果当前含有多个匹配的dom元素，则前面插入副本，最后一个插入本身
            results.cacheable || ( l > 1 && i < lastIndex ) ?
              jQuery.clone( fragment, true, true ) :
              fragment
          );
        }
      }
      // 如果有提取到script标签，则执行该标签
      if ( scripts.length ) {
        jQuery.each( scripts, evalScript );
      }
    }

    return this;
  }
});

// 修正table的tbody，
// 返回tbody元素
function root( elem, cur ) {
  return jQuery.nodeName(elem, "table") ?
    (elem.getElementsByTagName("tbody")[0] ||
    elem.appendChild(elem.ownerDocument.createElement("tbody"))) :
    elem;
}
```

``` js
// appendTo\prependTo\insertBefore\insertAfter目的恰恰和前面不带To的相反，一个主动、一个被动的区别
jQuery.each({
  appendTo: "append",
  prependTo: "prepend",
  insertBefore: "before",
  insertAfter: "after",
  replaceAll: "replaceWith"
}, function( name, original ) {
  jQuery.fn[ name ] = function( selector ) {
    var ret = [],
      insert = jQuery( selector ),
      parent = this.length === 1 && this[0].parentNode;

    if ( parent && parent.nodeType === 11 && parent.childNodes.length === 1 && insert.length === 1 ) {
      insert[ original ]( this[0] );
      return this;

    } else {
      for ( var i = 0, l = insert.length; i < l; i++ ) {
        var elems = ( i > 0 ? this.clone(true) : this ).get();
        jQuery( insert[i] )[ original ]( elems );
        ret = ret.concat( elems );
      }

      return this.pushStack( ret, name, insert.selector );
    }
  };
});
```


``` js
jQuery.extend({
  clone: function( elem, dataAndEvents, deepDataAndEvents ) {
    var srcElements,
      destElements,
      i,
      //   rnoshimcache = new RegExp("<(?:" + nodeNames + ")", "i"),
      // 浏览器支持html5元素 或者不包含html5元素 调用原生的cloneNode、否则调用shimCloneNode
      // IE<=8 does not properly clone detached, unknown element nodes
      clone = jQuery.support.html5Clone || !rnoshimcache.test( "<" + elem.nodeName ) ?
        elem.cloneNode( true ) :
        shimCloneNode( elem );

    // 如果浏览器不支持事件复制、或者不能正确复制checked状态
    // 则要借助cloneFixAttributes方法
    if ( (!jQuery.support.noCloneEvent || !jQuery.support.noCloneChecked) &&
        (elem.nodeType === 1 || elem.nodeType === 11) && !jQuery.isXMLDoc(elem) ) {
      // IE copies events bound via attachEvent when using cloneNode.
      // Calling detachEvent on the clone will also remove the events
      // from the original. In order to get around this, we use some
      // proprietary methods to clear the events. Thanks to MooTools
      // guys for this hotness.

      cloneFixAttributes( elem, clone );

      // Using Sizzle here is crazy slow, so we use getElementsByTagName instead
      // getAll = getElementsByTagName || querySelectorAll
      srcElements = getAll( elem );
      destElements = getAll( clone );

      // Weird iteration because IE will replace the length property
      // with an element if you are cloning the body and one of the
      // elements on the page has a name or id of "length"
      // 为每一个字元素都执行cloneFixAttributes修正方法
      for ( i = 0; srcElements[i]; ++i ) {
        // Ensure that the destination node is not null; Fixes #9587
        if ( destElements[i] ) {
          cloneFixAttributes( srcElements[i], destElements[i] );
        }
      }
    }

    // Copy the events from the original to the clone
    // 复制事件 按需求 是否深层复制事件
    // cloneCopyEvent
    if ( dataAndEvents ) {
      cloneCopyEvent( elem, clone );

      if ( deepDataAndEvents ) {
        srcElements = getAll( elem );
        destElements = getAll( clone );

        for ( i = 0; srcElements[i]; ++i ) {
          cloneCopyEvent( srcElements[i], destElements[i] );
        }
      }
    }

    srcElements = destElements = null;

    // Return the cloned set
    return clone;
  },
});
```

cloneFixAttributes 方法
``` js
function cloneFixAttributes( src, dest ) {
  var nodeName;

  // We do not need to do anything for non-Elements
  if ( dest.nodeType !== 1 ) {
    return;
  }

  // clearAttributes\mergeAttributes仅在IE678下支持， 会清除掉元素的属性和attachEvent的事件，但mergeAttributes只会把属性复制，不会复制事件
  // clearAttributes removes the attributes, which we don't want,
  // but also removes the attachEvent events, which we *do* want
  if ( dest.clearAttributes ) {
    dest.clearAttributes();
  }

  // mergeAttributes, in contrast, only merges back on the
  // original attributes, not the events
  if ( dest.mergeAttributes ) {
    dest.mergeAttributes( src );
  }

  nodeName = dest.nodeName.toLowerCase();

  // IE6-8 fail to clone children inside object elements that use
  // the proprietary classid attribute value (rather than the type
  // attribute) to identify the type of content to display
  // 如果一个元素师object，那么直接原封的文本复制即可
  if ( nodeName === "object" ) {
    dest.outerHTML = src.outerHTML;

  } else if ( nodeName === "input" && (src.type === "checkbox" || src.type === "radio") ) {
    // 如果是input、cehcekbox、radio则复制其checked、value
    // IE6-8 fails to persist the checked state of a cloned checkbox
    // or radio button. Worse, IE6-7 fail to give the cloned element
    // a checked appearance if the defaultChecked value isn't also set
    if ( src.checked ) {
      dest.defaultChecked = dest.checked = src.checked;
    }

    // IE6-7 get confused and end up setting the value of a cloned
    // checkbox/radio button to an empty string instead of "on"
    if ( dest.value !== src.value ) {
      dest.value = src.value;
    }

  // IE6-8 fails to return the selected option to the default selected
  // state when cloning options
  // 修正option的selected状态
  } else if ( nodeName === "option" ) {
    dest.selected = src.defaultSelected;

  // IE6-8 fails to set the defaultValue to the correct value when
  // cloning other types of input fields
  } else if ( nodeName === "input" || nodeName === "textarea" ) {
    dest.defaultValue = src.defaultValue;
  }

  // Event data gets referenced instead of copied if the expando
  // gets copied too
  dest.removeAttribute( jQuery.expando );
}
```
cloneCopyEvent 方法
``` js
function cloneCopyEvent( src, dest ) {
  // 检查src元素是否有事件
  if ( dest.nodeType !== 1 || !jQuery.hasData( src ) ) {
    return;
  }

  // 先将原始数据缓存复制到目标元素上
  var type, i, l,
    oldData = jQuery._data( src ),
    curData = jQuery._data( dest, oldData ),
    events = oldData.events;
  // 如果原始数据中有事件数据
  // 先将目标元素上面的事件数据清空掉
  if ( events ) {
    delete curData.handle;
    curData.events = {};
    // 然后将原始数据中的事件一个一个遍历添加到新的目标元素上
    for ( type in events ) {
      for ( i = 0, l = events[ type ].length; i < l; i++ ) {
        jQuery.event.add( dest, type + ( events[ type ][ i ].namespace ? "." : "" ) + events[ type ][ i ].namespace, events[ type ][ i ], events[ type ][ i ].data );
      }
    }
  }

  // make the cloned public data object a copy from the original
  if ( curData.data ) {
    curData.data = jQuery.extend( {}, curData.data );
  }
}
```

``` js
// 获取元素的所有子节点元素
contents: function( elem ) {
  return jQuery.nodeName( elem, "iframe" ) ?
    elem.contentDocument || elem.contentWindow.document :
    jQuery.makeArray( elem.childNodes );
}
```