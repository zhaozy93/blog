# chapter11 Style操作

一般我们不太建议直接对元素的style进行操作，建议通过更改类的方式来实现更改css的目的，在[chapter8](https://github.com/zhaozy93/blog/issues/17)中也提到过如何对类进行增、删、交替更换等操作。但有时还是免不了需要直接和属性打交道，在很多时候都要知道用户的鼠标路径、某个盒子的大小尺寸等信息。

jQuery中大概将样式部分分为了 css样式设置计算、类样式、坐标、尺寸这几类，类样式之前已经介绍过。不可否认的是在样式部分，各浏览器的兼容尤其是低版本的IE简直惨不忍睹，阅读一下早期版本的jQuery还是有不少帮助的。

## css样式设置计算

``` js
// 定义了一堆变量
var ralpha = /alpha\([^)]*\)/i,
  ropacity = /opacity=([^)]*)/,
  // fixed for IE9, see #8346
  rupper = /([A-Z]|^ms)/g,
  rnumpx = /^-?\d+(?:px)?$/i,
  rnum = /^-?\d/,
  rrelNum = /^([\-+])=([\-+.\de]+)/,
  cssShow = { position: "absolute", visibility: "hidden", display: "block" },
  cssWidth = [ "Left", "Right" ],
  cssHeight = [ "Top", "Bottom" ],
  curCSS,
  getComputedStyle,
  currentStyle;

jQuery.fn.css = function( name, value ) {
  // Setting 'undefined' is a no-op
  // 设置undefine给某一个属性是没意义的， 不执行任何操作
  if ( arguments.length === 2 && value === undefined ) {
    return this;
  }

  // 遍历所有当前元素，执行回调函数
  // 如果传入了value就是设置  jQuery.style( elem, name, value )
  // 没有value就是读取 jQuery.css( elem, name )
  return jQuery.access( this, name, value, true, function( elem, name, value ) {
    return value !== undefined ?
      jQuery.style( elem, name, value ) :
      jQuery.css( elem, name );
  });
};

jQuery.extend({
  // Add in style property hooks for overriding the default
  // behavior of getting and setting a style property
  cssHooks: {
    // 这便是默认的支持opacity浏览器的对opacity的hook
    // 确保返回的一定是一个数字
    opacity: {
      get: function( elem, computed ) {
        if ( computed ) {
          // We should always get a number back from opacity
          var ret = curCSS( elem, "opacity", "opacity" );
          return ret === "" ? "1" : ret;

        } else {
          return elem.style.opacity;
        }
      }
    }
  },

  // Exclude the following css properties to add px
  cssNumber: {
    "fillOpacity": true,
    "fontWeight": true,
    "lineHeight": true,
    "opacity": true,
    "orphans": true,
    "widows": true,
    "zIndex": true,
    "zoom": true
  },

  // Add in properties whose names you wish to fix before
  // setting or getting the value
  cssProps: {
    // normalize float css property
    "float": jQuery.support.cssFloat ? "cssFloat" : "styleFloat"
  },

  // Get and set the style property on a DOM Node
  // 在jQuery.fn.css中提到用style方法设置属性
  style: function( elem, name, value, extra ) {
    // Don't set styles on text and comment nodes
    // 不对文本节点和注释节点操作， 也不对不支持样式的节点进行操作
    if ( !elem || elem.nodeType === 3 || elem.nodeType === 8 || !elem.style ) {
      return;
    }

    // Make sure that we're working with the right name
    // 第一步就是将名字转为驼峰写法
    var ret, type, origName = jQuery.camelCase( name ),
      style = elem.style, hooks = jQuery.cssHooks[ origName ];

    // 继续修正属性名
    // float名字比较特殊，ie下面是styleFloat，火狐chrome则是cssFloat
    name = jQuery.cssProps[ origName ] || origName;

    // Check if we're setting a value
    // 如果是设置属性
    if ( value !== undefined ) {
      type = typeof value;

      // convert relative number strings (+= or -=) to relative numbers. #7345
      // rrelNum = /^([\-+])=([\-+.\de]+)/
      // 如果value是字符串的情况下，我们检测他是不是 += 或者-=, 这种情况我们需要先进行计算
      if ( type === "string" && (ret = rrelNum.exec( value )) ) {
        value = ( +( ret[1] + 1) * +ret[2] ) + parseFloat( jQuery.css( elem, name ) );
        // Fixes bug #9237
        type = "number";
      }

      // Make sure that NaN and null values aren't set. See: #7116
      if ( value == null || type === "number" && isNaN( value ) ) {
        return;
      }

      // If a number was passed in, add 'px' to the (except for certain CSS properties)
      // 修正value值，是否需要增加 px后缀
      if ( type === "number" && !jQuery.cssNumber[ origName ] ) {
        value += "px";
      }

      // If a hook was provided, use that value, otherwise just set the specified value
      // 如果有hooks传入，则优先使用hooks.set方法进行属性设置
      // 最后再调用elem.style进行属性设置。
      if ( !hooks || !("set" in hooks) || (value = hooks.set( elem, value )) !== undefined ) {
        // Wrapped to prevent IE from throwing errors when 'invalid' values are provided
        // Fixes bug #5509
        try {
          style[ name ] = value;
        } catch(e) {}
      }

    } else {
      // 如果是读取属性
      // 现场时用hooks.get读取\再使用elem.style读取，不过这样读出来的仅仅是内连样式哦哦哦哦！！！！！！
      // If a hook was provided get the non-computed value from there
      if ( hooks && "get" in hooks && (ret = hooks.get( elem, false, extra )) !== undefined ) {
        return ret;
      }

      // Otherwise just get the value from the style object
      return style[ name ];
    }
  },

  // 在jQuery.fn.css中提到用css方法读取属性
  // 与style方法中读取属性基本一致，只是最后是使用curCss来读取计算后的最终属性，而非内连属性
  // curCSS = getComputedStyle || currentStyle;
  // 而非jQuery.curCSS
  css: function( elem, name, extra ) {
    var ret, hooks;

    // Make sure that we're working with the right name
    name = jQuery.camelCase( name );
    hooks = jQuery.cssHooks[ name ];
    name = jQuery.cssProps[ name ] || name;

    // cssFloat needs a special treatment
    if ( name === "cssFloat" ) {
      name = "float";
    }

    // If a hook was provided get the computed value from there
    if ( hooks && "get" in hooks && (ret = hooks.get( elem, true, extra )) !== undefined ) {
      return ret;

    // Otherwise, if a way to get the computed value exists, use that
    } else if ( curCSS ) {
      return curCSS( elem, name );
    }
  },
   
  // 为某些情况提前准备环境，
  // 先将options的属性复制给elem， 然后执行callback
  // 最后再将原本覆盖的属性恢复回去
  // A method for quickly swapping in/out CSS properties to get correct calculations
  swap: function( elem, options, callback ) {
    var old = {};

    // Remember the old values, and insert the new ones
    // 先记录原本的属性值
    // 然后将options中的属性值附上去
    for ( var name in options ) {
      old[ name ] = elem.style[ name ];
      elem.style[ name ] = options[ name ];
    }
    // 执行callback
    callback.call( elem );

    // 再还原属性值
    // Revert the old values
    for ( name in options ) {
      elem.style[ name ] = old[ name ];
    }
  }
});

// DEPRECATED, Use jQuery.css() instead
jQuery.curCSS = jQuery.css;

jQuery.each(["height", "width"], function( i, name ) {
  jQuery.cssHooks[ name ] = {
    get: function( elem, computed, extra ) {
      var val;
      // 只支持读取计算样式，那么css调用该方法会得到undefined
      if ( computed ) {
        // 如果该元素可见，那么调用getWH， 否则调用swap方法
        if ( elem.offsetWidth !== 0 ) {
          return getWH( elem, name, extra );
        } else {
          // cssShow = { position: "absolute", visibility: "hidden", display: "block" },
          jQuery.swap( elem, cssShow, function() {
            val = getWH( elem, name, extra );
          });
        }

        return val;
      }
    },

    set: function( elem, value ) {
      // rnumpx来检测属性是不是不是以px结尾
      // rnumpx = /^-?\d+(?:px)?$/i,
      // 只有纯数字或者以px结尾可以通过这个测试
      // parseFloat 则可以忽略px后缀，只返回带符号的数字部分
      if ( rnumpx.test( value ) ) {
        // ignore negative width and height values #1599
        value = parseFloat( value );

        if ( value >= 0 ) {
          return value + "px";
        }

      } else {
        return value;
      }
    }
  };
});

// 这里是对浏览器不支持opacity的hook
// 浏览器不支持opacity的时候需要使用filter来代替
if ( !jQuery.support.opacity ) {
  jQuery.cssHooks.opacity = {
    get: function( elem, computed ) {
      // IE uses filters for opacity
      // computed为真时使用计算样式、否则使用内连样式， 使用ropacity检测是否设置透明度
      // 设置了则返回百分比之后的值、否则计算的话显示1， 内敛的返回''空字符串
      // ropacity = /opacity=([^)]*)/,
      return ropacity.test( (computed && elem.currentStyle ? elem.currentStyle.filter : elem.style.filter) || "" ) ?
        ( parseFloat( RegExp.$1 ) / 100 ) + "" :
        computed ? "1" : "";
    },

    set: function( elem, value ) {
      var style = elem.style,
        currentStyle = elem.currentStyle,
        opacity = jQuery.isNumeric( value ) ? "alpha(opacity=" + value * 100 + ")" : "",
        filter = currentStyle && currentStyle.filter || style.filter || "";

      // IE has trouble with opacity if it does not have layout
      // Force it by setting the zoom level
      style.zoom = 1;

      // 如果是设置filter为1， 那么直接移除元素属性filter就好
      // if setting opacity to 1, and no other filters exist - attempt to remove filter attribute #6652
      if ( value >= 1 && jQuery.trim( filter.replace( ralpha, "" ) ) === "" ) {

        // Setting style.filter to null, "" & " " still leave "filter:" in the cssText
        // if "filter:" is present at all, clearType is disabled, we want to avoid this
        // style.removeAttribute is IE Only, but so apparently is this code path...
        style.removeAttribute( "filter" );
        // 如果成功移除了 return即可
        // if there there is no filter style applied in a css rule, we are done
        if ( currentStyle && !currentStyle.filter ) {
          return;
        }
      }

      // 乖乖的设置filter属性
      // otherwise, set new filter values
      // ralpha = /alpha\([^)]*\)/i, 
      style.filter = ralpha.test( filter ) ?
        filter.replace( ralpha, opacity ) :
        filter + " " + opacity;
    }
  };
}

jQuery(function() {
  // This hook cannot be added until DOM ready because the support test
  // for it is not run until after DOM ready
  if ( !jQuery.support.reliableMarginRight ) {
    jQuery.cssHooks.marginRight = {
      // 对于margin-right紊乱的情况，先讲属性设置为inline-block即可正确获取它的marinRight
      get: function( elem, computed ) {
        // WebKit Bug 13343 - getComputedStyle returns wrong value for margin-right
        // Work around by temporarily setting element display to inline-block
        var ret;
        jQuery.swap( elem, { "display": "inline-block" }, function() {
          if ( computed ) {
            ret = curCSS( elem, "margin-right", "marginRight" );
          } else {
            ret = elem.style.marginRight;
          }
        });
        return ret;
      }
    };
  }
});

// 封装getComputedStyle方法
if ( document.defaultView && document.defaultView.getComputedStyle ) {
  getComputedStyle = function( elem, name ) {
    var ret, defaultView, computedStyle;

    // rupper = /([A-Z]|^ms)/g,
    name = name.replace( rupper, "-$1" ).toLowerCase();

    if ( (defaultView = elem.ownerDocument.defaultView) &&
        (computedStyle = defaultView.getComputedStyle( elem, null )) ) {
      ret = computedStyle.getPropertyValue( name );
      if ( ret === "" && !jQuery.contains( elem.ownerDocument.documentElement, elem ) ) {
        // computedStyle没有计算到元素则尝试读取内连属性
        ret = jQuery.style( elem, name );
      }
    }

    return ret;
  };
}

if ( document.documentElement.currentStyle ) {
  currentStyle = function( elem, name ) {
    var left, rsLeft, uncomputed,
      ret = elem.currentStyle && elem.currentStyle[ name ],
      style = elem.style;

    // Avoid setting ret to empty string here
    // so we don't default to auto
    // 如果currentStyle没有拿到属性，则尝试使用elem.style读取内连属性
    if ( ret === null && style && (uncomputed = style[ name ]) ) {
      ret = uncomputed;
    }

    // From the awesome hack by Dean Edwards
    // http://erik.eae.net/archives/2007/07/27/18.54.15/#comment-102291

    // If we're not dealing with a regular pixel number
    // but a number that has a weird ending, we need to convert it to pixels
    // rnumpx = /^-?\d+(?:px)?$/i,
    // rnum = /^-?\d/,
    // rnumpx来检测属性是不是不是以px结尾，  rnum测试属性是不是数字
    // 如果是这样的话例如 1rem、1em、65%， 我们需要将它们转为px单位
    // 先记录属性原本的left值。 利用runtimeStyle、pixelLeft来计算百分比， 最后再恢复原本的left值
    if ( !rnumpx.test( ret ) && rnum.test( ret ) ) {

      // Remember the original values
      left = style.left;
      rsLeft = elem.runtimeStyle && elem.runtimeStyle.left;

      // Put in the new values to get a computed value out
      if ( rsLeft ) {
        elem.runtimeStyle.left = elem.currentStyle.left;
      }
      style.left = name === "fontSize" ? "1em" : ( ret || 0 );
      ret = style.pixelLeft + "px";

      // Revert the changed values
      style.left = left;
      if ( rsLeft ) {
        elem.runtimeStyle.left = rsLeft;
      }
    }

    return ret === "" ? "auto" : ret;
  };
}

curCSS = getComputedStyle || currentStyle;

```

## 坐标系统计算

``` js
var rtable = /^t(?:able|d|h)$/i,
  rroot = /^(?:body|html)$/i;

// jQuery.fn.offset在支持getBoundingClientRect情况下调用getBoundingClientRect方法，
// 同时也为不支持getBoundingClientRect方法的做了兼容处理
// 其实IE678910全部支持getBoundingClientRect方法，并不知道要兼容哪些浏览器
if ( "getBoundingClientRect" in document.documentElement ) {
  jQuery.fn.offset = function( options ) {
    var elem = this[0], box;

    if ( options ) {
      // 如果有options表示为每个元素设置位置
      return this.each(function( i ) {
        jQuery.offset.setOffset( this, options, i );
      });
    }

    if ( !elem || !elem.ownerDocument ) {
      return null;
    }
    // body的位置需要单独计算
    if ( elem === elem.ownerDocument.body ) {
      return jQuery.offset.bodyOffset( elem );
    }

    try {
      box = elem.getBoundingClientRect();
    } catch(e) {}

    var doc = elem.ownerDocument,
      docElem = doc.documentElement;

    // Make sure we're not dealing with a disconnected DOM node
    // 如果当前元素不在文档中，返回{0,0}
    if ( !box || !jQuery.contains( docElem, elem ) ) {
      return box ? { top: box.top, left: box.left } : { top: 0, left: 0 };
    }

    // 距离文档上标 = 距窗口上标 + 垂直滚动偏移 -文档左上边框厚度
    // 距离文档左标 = 距窗口左标 + 水平滚动便宜 -文档左边框厚度
    var body = doc.body,
      win = getWindow(doc),
      clientTop  = docElem.clientTop  || body.clientTop  || 0,
      clientLeft = docElem.clientLeft || body.clientLeft || 0,
      scrollTop  = win.pageYOffset || jQuery.support.boxModel && docElem.scrollTop  || body.scrollTop,
      scrollLeft = win.pageXOffset || jQuery.support.boxModel && docElem.scrollLeft || body.scrollLeft,
      top  = box.top  + scrollTop  - clientTop,
      left = box.left + scrollLeft - clientLeft;

    return { top: top, left: left };
  };

} else {
  // 在不支持getBoundingClientRect的浏览器中，利用offsetParent来一直循环计算top、left数值
  jQuery.fn.offset = function( options ) {
    var elem = this[0];

    if ( options ) {
      return this.each(function( i ) {
        jQuery.offset.setOffset( this, options, i );
      });
    }

    if ( !elem || !elem.ownerDocument ) {
      return null;
    }

    if ( elem === elem.ownerDocument.body ) {
      return jQuery.offset.bodyOffset( elem );
    }

    var computedStyle,
      offsetParent = elem.offsetParent,
      prevOffsetParent = elem,
      doc = elem.ownerDocument,
      docElem = doc.documentElement,
      body = doc.body,
      defaultView = doc.defaultView,
      prevComputedStyle = defaultView ? defaultView.getComputedStyle( elem, null ) : elem.currentStyle,
      top = elem.offsetTop,
      left = elem.offsetLeft;

    while ( (elem = elem.parentNode) && elem !== body && elem !== docElem ) {
      if ( jQuery.support.fixedPosition && prevComputedStyle.position === "fixed" ) {
        break;
      }

      computedStyle = defaultView ? defaultView.getComputedStyle(elem, null) : elem.currentStyle;
      top  -= elem.scrollTop;
      left -= elem.scrollLeft;

      if ( elem === offsetParent ) {
        top  += elem.offsetTop;
        left += elem.offsetLeft;

        if ( jQuery.support.doesNotAddBorder && !(jQuery.support.doesAddBorderForTableAndCells && rtable.test(elem.nodeName)) ) {
          top  += parseFloat( computedStyle.borderTopWidth  ) || 0;
          left += parseFloat( computedStyle.borderLeftWidth ) || 0;
        }

        prevOffsetParent = offsetParent;
        offsetParent = elem.offsetParent;
      }

      if ( jQuery.support.subtractsBorderForOverflowNotVisible && computedStyle.overflow !== "visible" ) {
        top  += parseFloat( computedStyle.borderTopWidth  ) || 0;
        left += parseFloat( computedStyle.borderLeftWidth ) || 0;
      }

      prevComputedStyle = computedStyle;
    }

    if ( prevComputedStyle.position === "relative" || prevComputedStyle.position === "static" ) {
      top  += body.offsetTop;
      left += body.offsetLeft;
    }

    if ( jQuery.support.fixedPosition && prevComputedStyle.position === "fixed" ) {
      top  += Math.max( docElem.scrollTop, body.scrollTop );
      left += Math.max( docElem.scrollLeft, body.scrollLeft );
    }

    return { top: top, left: left };
  };
}

jQuery.offset = {
  // 计算body的文档坐标
  bodyOffset: function( body ) {
    var top = body.offsetTop,
      left = body.offsetLeft;

    if ( jQuery.support.doesNotIncludeMarginInBodyOffset ) {
      top  += parseFloat( jQuery.css(body, "marginTop") ) || 0;
      left += parseFloat( jQuery.css(body, "marginLeft") ) || 0;
    }

    return { top: top, left: left };
  },

  setOffset: function( elem, options, i ) {
    var position = jQuery.css( elem, "position" );

    // set position first, in-case top/left are set even on static elem
    // 如果当前position为static，则设置为relative，否则不起效
    if ( position === "static" ) {
      elem.style.position = "relative";
    }

    var curElem = jQuery( elem ),
      curOffset = curElem.offset(), // 当前元素的文档坐标
      curCSSTop = jQuery.css( elem, "top" ),  // 当前元素的top
      curCSSLeft = jQuery.css( elem, "left" ),  // 当前元素的left值
      calculatePosition = ( position === "absolute" || position === "fixed" ) && jQuery.inArray("auto", [curCSSTop, curCSSLeft]) > -1,
      props = {}, curPosition = {}, curTop, curLeft;

    // need to be able to calculate position if either top or left is auto and position is either absolute or fixed
    // 继续修正参数
    if ( calculatePosition ) {
      curPosition = curElem.position();
      curTop = curPosition.top;
      curLeft = curPosition.left;
    } else {
      curTop = parseFloat( curCSSTop ) || 0;
      curLeft = parseFloat( curCSSLeft ) || 0;
    }

    if ( jQuery.isFunction( options ) ) {
      options = options.call( elem, i, curOffset );
    }

    // 内连top = 目标文档坐标top - 当前文档坐标top + 计算样式top
    // 内连left = 目标文档坐标left - 当前文档坐标left + 计算样式left
    if ( options.top != null ) {
      props.top = ( options.top - curOffset.top ) + curTop;
    }
    if ( options.left != null ) {
      props.left = ( options.left - curOffset.left ) + curLeft;
    }

    if ( "using" in options ) {
      options.using.call( elem, props );
    } else {
      curElem.css( props );
    }
  }
};


jQuery.fn.extend({
  // 寻找第一个元素相对于定位祖先元素的位置
  // 利用元素本身的offset和定位元素的offset
  // 但是元素本身需要减去对应的margin、 父元素需要增加border的宽度
  position: function() {
    if ( !this[0] ) {
      return null;
    }

    var elem = this[0],

    // Get *real* offsetParent
    offsetParent = this.offsetParent(),

    // Get correct offsets
    offset       = this.offset(),
    parentOffset = rroot.test(offsetParent[0].nodeName) ? { top: 0, left: 0 } : offsetParent.offset();

    // Subtract element margins
    // note: when an element has margin: auto the offsetLeft and marginLeft
    // are the same in Safari causing offset.left to incorrectly be 0
    offset.top  -= parseFloat( jQuery.css(elem, "marginTop") ) || 0;
    offset.left -= parseFloat( jQuery.css(elem, "marginLeft") ) || 0;

    // Add offsetParent borders
    parentOffset.top  += parseFloat( jQuery.css(offsetParent[0], "borderTopWidth") ) || 0;
    parentOffset.left += parseFloat( jQuery.css(offsetParent[0], "borderLeftWidth") ) || 0;

    // Subtract the two offsets
    return {
      top:  offset.top  - parentOffset.top,
      left: offset.left - parentOffset.left
    };
  },

  offsetParent: function() {
    // 循环寻找定位祖先元素
    // 两个退出条件、定位祖先不能是static、或者祖先是body
    return this.map(function() {
      var offsetParent = this.offsetParent || document.body;
      while ( offsetParent && (!rroot.test(offsetParent.nodeName) && jQuery.css(offsetParent, "position") === "static") ) {
        offsetParent = offsetParent.offsetParent;
      }
      return offsetParent;
    });
  }
});


// Create scrollLeft and scrollTop methods
jQuery.each( ["Left", "Top"], function( i, name ) {
  var method = "scroll" + name;

  jQuery.fn[ method ] = function( val ) {
    var elem, win;
    // 读取第一个元素的Scroll值
    if ( val === undefined ) {
      elem = this[ 0 ];

      if ( !elem ) {
        return null;
      }

      win = getWindow( elem );

      // Return the scroll offset
      return win ? ("pageXOffset" in win) ? win[ i ? "pageYOffset" : "pageXOffset" ] :
        jQuery.support.boxModel && win.document.documentElement[ method ] ||
          win.document.body[ method ] :
        elem[ method ];
    }

    // Set the scroll offset
    // 为每个元素设置scroll值
    return this.each(function() {
      win = getWindow( this );

      if ( win ) {
        win.scrollTo(
          !i ? val : jQuery( win ).scrollLeft(),
           i ? val : jQuery( win ).scrollTop()
        );

      } else {
        this[ method ] = val;
      }
    });
  };
});
```

## 尺寸计算
``` js
// getWH是整个尺寸计算的基础函数
function getWH( elem, name, extra ) {

  // Start with offset property
  // cssWidth = [ "Left", "Right" ],
  // cssHeight = [ "Top", "Bottom" ],
  // offset包含content、padding、border、但是不包含margin
  var val = name === "width" ? elem.offsetWidth : elem.offsetHeight,
    which = name === "width" ? cssWidth : cssHeight,
    i = 0,
    len = which.length;
  // 表示元素可见
  if ( val > 0 ) {
    if ( extra !== "border" ) {
      for ( ; i < len; i++ ) {
        // 没有extra表示只计算content， 需要先减去padding，再减去border
        if ( !extra ) {
          val -= parseFloat( jQuery.css( elem, "padding" + which[ i ] ) ) || 0;
        }
        // 如果extra是border，则加上border的宽度
        if ( extra === "margin" ) {
          val += parseFloat( jQuery.css( elem, extra + which[ i ] ) ) || 0;
        } else {
          val -= parseFloat( jQuery.css( elem, "border" + which[ i ] + "Width" ) ) || 0;
        }
      }
    }
    // 如果extra是border则刚好不需要任何修复
    return val + "px";
  }

  // Fall back to computed then uncomputed css if necessary
  // 先计算样式再内连样式，再根绝extra修正
  val = curCSS( elem, name, name );
  if ( val < 0 || val == null ) {
    val = elem.style[ name ] || 0;
  }
  // Normalize "", auto, and prepare for extra
  val = parseFloat( val ) || 0;

  // Add padding, border, margin
  if ( extra ) {
    for ( ; i < len; i++ ) {
      val += parseFloat( jQuery.css( elem, "padding" + which[ i ] ) ) || 0;
      if ( extra !== "padding" ) {
        val += parseFloat( jQuery.css( elem, "border" + which[ i ] + "Width" ) ) || 0;
      }
      if ( extra === "margin" ) {
        val += parseFloat( jQuery.css( elem, extra + which[ i ] ) ) || 0;
      }
    }
  }

  return val + "px";
}

// Create width, height, innerHeight, innerWidth, outerHeight and outerWidth methods
jQuery.each([ "Height", "Width" ], function( i, name ) {

  var type = name.toLowerCase();

  // innerHeight and innerWidth
  // 获取第一个元素的inner属性
  // 没有元素返回null
  // 有元素有style属性，返回parseFloat( jQuery.css( elem, type, "padding" ) )
  // 有元素没style返回elem.width()
  // jQuery.css( elem, type, "padding" ) 其实还是最终调用hooks--> getWH
  //
  jQuery.fn[ "inner" + name ] = function() {
    var elem = this[0];
    return elem ?
      elem.style ?
      parseFloat( jQuery.css( elem, type, "padding" ) ) :
      this[ type ]() :
      null;
  };

  // outerHeight and outerWidth
  jQuery.fn[ "outer" + name ] = function( margin ) {
    var elem = this[0];
    return elem ?
      elem.style ?
      parseFloat( jQuery.css( elem, type, margin ? "margin" : "border" ) ) :
      this[ type ]() :
      null;
  };

  jQuery.fn[ type ] = function( size ) {
    // Get window width or height
    var elem = this[0];
    if ( !elem ) {
      return size == null ? null : this;
    }

    if ( jQuery.isFunction( size ) ) {
      return this.each(function( i ) {
        var self = jQuery( this );
        self[ type ]( size.call( this, i, self[ type ]() ) );
      });
    }
    // 计算window元素的width和height
    if ( jQuery.isWindow( elem ) ) {
      // Everyone else use document.documentElement or document.body depending on Quirks vs Standards mode
      // 3rd condition allows Nokia support, as it supports the docElem prop but not CSS1Compat
      var docElemProp = elem.document.documentElement[ "client" + name ],
        body = elem.document.body;
      return elem.document.compatMode === "CSS1Compat" && docElemProp ||
        body && body[ "client" + name ] || docElemProp;

    // document对象的width和height
    } else if ( elem.nodeType === 9 ) {
      // Either scroll[Width/Height] or offset[Width/Height], whichever is greater
      return Math.max(
        elem.documentElement["client" + name],
        elem.body["scroll" + name], elem.documentElement["scroll" + name],
        elem.body["offset" + name], elem.documentElement["offset" + name]
      );

    // Get or set width or height on the element
    // 直接css读取witdh和heigt对于普通元素
    } else if ( size === undefined ) {
      var orig = jQuery.css( elem, type ),
        ret = parseFloat( orig );

      return jQuery.isNumeric( ret ) ? ret : orig;

    // Set the width or height on the element (default to pixels if value is unitless)
    } else {
      // 为普通元素添加width和 height属性
      return this.css( type, typeof size === "string" ? size : size + "px" );
    }
  };

});
```