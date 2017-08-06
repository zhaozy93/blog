# chapter8 属性操作

对dom的操作才是jQuery 最擅长的事情，那么dom属性肯定是其中必不可少的一环。 jQuery中属性操作主要分为了4个部分

- HTML属性操作
- DOM属性操作
- 类操作

其中属性操作部分的接口也符合jQuery的一贯思路，同一个方法既是读又是取甚至可以删。

在这里jQuery做的工作并不多，主要是针对一些特殊的属性值做了某些兼容性封装处理，绝大多数还是依靠`setAttribute、getAttribute、elem.key`来完成工作。 不过这里面也有不少值得学习的地方，比如`~`运算符和`ttributeNode`的运用。同时也是对HTML属性和DOM属性区分的加深。

源码也分析了这么多了，可以看出一个规律，jQuery中被我们调用的接口或者说方法、API大都是添加到jQuery.fn上面的，通过`jQuery.fn.extend`方法，但真正做事情的要去jQuery上面找大致同名或相似名的方法。 其实这也就是利用继承的关系，因为jQuery.fn是继承的原型链的一环。


``` js
jQuery.fn.extend({
  attr: function( name, value ) {
    // access函数比较简单，做分发操作为多个接口提供支持
    // 当没有value时，返回第一个元素的属性值
    // 当有value时，循环为每个元素都设置属性值
    return jQuery.access( this, name, value, true, jQuery.attr );
  },

  removeAttr: function( name ) {
    return this.each(function() {
      jQuery.removeAttr( this, name );
    });
  },

  prop: function( name, value ) {
    return jQuery.access( this, name, value, true, jQuery.prop );
  },

  removeProp: function( name ) {
    name = jQuery.propFix[ name ] || name;
    return this.each(function() {
      // try/catch handles cases where IE balks (such as removing a property on window)
      // ie在do元素上删除属性会报错 所以先设undefined再删除
      try {
        this[ name ] = undefined;
        delete this[ name ];
      } catch( e ) {}
    });
  },

  // addClass逻辑简单，代码也比较简单
  // 如果原本没有类，则直接设置新类
  // 如果原来有 则匹配新的类，有就忽略，没有就添加到最后
  // 里面用到了 !~ 操作符合
  // !就是简单的取反
  // ~是先变换符号再减1
  addClass: function( value ) {
    var classNames, i, l, elem,
      setClass, c, cl;

    // 如果是传入的时方法
    // 为每个元素调用方法之后再执行addClass
    if ( jQuery.isFunction( value ) ) {
      return this.each(function( j ) {
        jQuery( this ).addClass( value.call(this, j, this.className) );
      });
    }

    if ( value && typeof value === "string" ) {
      classNames = value.split( rspace );

      for ( i = 0, l = this.length; i < l; i++ ) {
        elem = this[ i ];

        if ( elem.nodeType === 1 ) {
          if ( !elem.className && classNames.length === 1 ) {
            elem.className = value;

          } else {
            setClass = " " + elem.className + " ";

            for ( c = 0, cl = classNames.length; c < cl; c++ ) {
              if ( !~setClass.indexOf( " " + classNames[ c ] + " " ) ) {
                setClass += classNames[ c ] + " ";
              }
            }
            elem.className = jQuery.trim( setClass );
          }
        }
      }
    }

    return this;
  },

  // 内部用正则的方式去替换需要删掉的类名
  removeClass: function( value ) {
    var classNames, i, l, elem, className, c, cl;

    if ( jQuery.isFunction( value ) ) {
      return this.each(function( j ) {
        jQuery( this ).removeClass( value.call(this, j, this.className) );
      });
    }

    if ( (value && typeof value === "string") || value === undefined ) {
      classNames = ( value || "" ).split( rspace );

      for ( i = 0, l = this.length; i < l; i++ ) {
        elem = this[ i ];

        if ( elem.nodeType === 1 && elem.className ) {
          if ( value ) {
            className = (" " + elem.className + " ").replace( rclass, " " );
            for ( c = 0, cl = classNames.length; c < cl; c++ ) {
              className = className.replace(" " + classNames[ c ] + " ", " ");
            }
            elem.className = jQuery.trim( className );

          } else {
            elem.className = "";
          }
        }
      }
    }

    return this;
  },

  // 顾名思义， 切换类名
  // 最常用在动画中 比如展开某动画的时候需要添加类 .slider-show
  // 关闭的时候需要类名 .slider-hidden  两者应该只保留一个 也就是一直在切换
  toggleClass: function( value, stateVal ) {
    var type = typeof value,
      isBool = typeof stateVal === "boolean";

    if ( jQuery.isFunction( value ) ) {
      return this.each(function( i ) {
        jQuery( this ).toggleClass( value.call(this, i, this.className, stateVal), stateVal );
      });
    }

    return this.each(function() {
      if ( type === "string" ) {
        // toggle individual class names
        var className,
          i = 0,
          self = jQuery( this ),
          state = stateVal,
          classNames = value.split( rspace );
        
        // 原本有则删除、原本没有则添加
        while ( (className = classNames[ i++ ]) ) {
          // check each className given, space seperated list
          state = isBool ? state : !self.hasClass( className );
          self[ state ? "addClass" : "removeClass" ]( className );
        }

      } else if ( type === "undefined" || type === "boolean" ) {
        // 需要借助缓存系统来实现缓存， 在全部替换的时候再取出来
        if ( this.className ) {
          // store className if set
          jQuery._data( this, "__className__", this.className );
        }

        // toggle whole className
        this.className = this.className || value === false ? "" : jQuery._data( this, "__className__" ) || "";
      }
    });
  },
  
  // 判断一个元素是否含有某个类名
  hasClass: function( selector ) {
    var className = " " + selector + " ",
      i = 0,
      l = this.length;
    for ( ; i < l; i++ ) {
      if ( this[i].nodeType === 1 && (" " + this[i].className + " ").replace(rclass, " ").indexOf( className ) > -1 ) {
        return true;
      }
    }

    return false;
  },

  // 对elem.value的一个封装
  val: function( value ) {
    var hooks, ret, isFunction,
      elem = this[0];
    // 没有argumentts表示获取第一个元素的value值
    // 有hook时优先使用hook的get方法，没有再普通的elem.value
    if ( !arguments.length ) {
      if ( elem ) {
        hooks = jQuery.valHooks[ elem.nodeName.toLowerCase() ] || jQuery.valHooks[ elem.type ];

        if ( hooks && "get" in hooks && (ret = hooks.get( elem, "value" )) !== undefined ) {
          return ret;
        }

        ret = elem.value;

        return typeof ret === "string" ?
          // handle most common string cases
          ret.replace(rreturn, "") :
          // handle cases where value is null/undef or number
          ret == null ? "" : ret;
      }

      return;
    }

    isFunction = jQuery.isFunction( value );
    // 先修正参数， 保证是字符串
    // 再hook。set 最后 elem.value 
    return this.each(function( i ) {
      var self = jQuery(this), val;

      if ( this.nodeType !== 1 ) {
        return;
      }

      if ( isFunction ) {
        val = value.call( this, i, self.val() );
      } else {
        val = value;
      }

      // Treat null/undefined as ""; convert numbers to string
      if ( val == null ) {
        val = "";
      } else if ( typeof val === "number" ) {
        val += "";
      } else if ( jQuery.isArray( val ) ) {
        val = jQuery.map(val, function ( value ) {
          return value == null ? "" : value + "";
        });
      }

      hooks = jQuery.valHooks[ this.nodeName.toLowerCase() ] || jQuery.valHooks[ this.type ];

      // If set returns undefined, fall back to normal setting
      if ( !hooks || !("set" in hooks) || hooks.set( this, val, "value" ) === undefined ) {
        this.value = val;
      }
    });
  }
});

// 
jQuery.extend({
  valHooks: {
    // option元素的兼容性修复
    option: {
      get: function( elem ) {
        // attributes.value is undefined in Blackberry 4.7 but
        // uses .value. See #6932
        var val = elem.attributes.value;
        return !val || val.specified ? elem.value : elem.text;
      }
    },
    select: {
      // select有单选、多选
      // 内部option还有被disabled的情况，需要单独处理， 逻辑不复杂 代码也简单
      get: function( elem ) {
        var value, i, max, option,
          index = elem.selectedIndex,
          values = [],
          options = elem.options,
          one = elem.type === "select-one";

        // Nothing was selected
        if ( index < 0 ) {
          return null;
        }

        // Loop through all the selected options
        i = one ? index : 0;
        max = one ? index + 1 : options.length;
        for ( ; i < max; i++ ) {
          option = options[ i ];

          // Don't return options that are disabled or in a disabled optgroup
          if ( option.selected && (jQuery.support.optDisabled ? !option.disabled : option.getAttribute("disabled") === null) &&
              (!option.parentNode.disabled || !jQuery.nodeName( option.parentNode, "optgroup" )) ) {

            // Get the specific value for the option
            value = jQuery( option ).val();

            // We don't need an array for one selects
            if ( one ) {
              return value;
            }

            // Multi-Selects return an array
            values.push( value );
          }
        }

        // Fixes Bug #2551 -- select.val() broken in IE after form.reset()
        if ( one && !values.length && options.length ) {
          return jQuery( options[ index ] ).val();
        }

        return values;
      },
      // 把value组成一个数组
      // 循环所有option逐个修改selected值
      // 如果没有value则设置selectedIndex为-1
      set: function( elem, value ) {
        var values = jQuery.makeArray( value );

        jQuery(elem).find("option").each(function() {
          this.selected = jQuery.inArray( jQuery(this).val(), values ) >= 0;
        });

        if ( !values.length ) {
          elem.selectedIndex = -1;
        }
        return values;
      }
    }
  },

  attrFn: {
    val: true,
    css: true,
    html: true,
    text: true,
    data: true,
    width: true,
    height: true,
    offset: true
  },

  attr: function( elem, name, value, pass ) {
    var ret, hooks, notxml,
      nType = elem.nodeType;

    // 不是真正的dom元素 直接返回
    // 不知道为什么不直接判断nType === 1\9\11
    // 其余的nodetype 并不常用，不太熟悉
    if ( !elem || nType === 3 || nType === 8 || nType === 2 ) {
      return;
    }
    // 如果是几个特殊的属性值，那么直接调用对应的方法设置，有专用的方法
    // 这几个属性分别是val、css、html、text、da、ta、width、height、offset
    if ( pass && name in jQuery.attrFn ) {
      return jQuery( elem )[ name ]( value );
    }

    // 如果不支持getAttribute属性，则调用以一种思路
    // jQuery.prop
    if ( typeof elem.getAttribute === "undefined" ) {
      return jQuery.prop( elem, name, value );
    }

    // 如果元素不是最普通的dom元素 或者是当前是html环境
    notxml = nType !== 1 || !jQuery.isXMLDoc( elem );

    // 属性值必须替换为小写，因为html属性不支持大写
    // 又是一个修正部分，稍后解释
    if ( notxml ) {
      name = name.toLowerCase();
      hooks = jQuery.attrHooks[ name ] || ( rboolean.test( name ) ? boolHook : nodeHook );
    }

    // value不是undefined表明不是读取
    if ( value !== undefined ) {

      // value是null表示移除属性
      if ( value === null ) {
        jQuery.removeAttr( elem, name );
        return;
        // 如果有hooks并且hooks有set方法 并且设置后返回值不是undefined 则表示设置成功啦
      } else if ( hooks && "set" in hooks && notxml && (ret = hooks.set( elem, value, name )) !== undefined ) {
        return ret;

      } else {
        // 最普通的也是最后的选项就是setAttribute设置属性
        elem.setAttribute( name, "" + value );
        return value;
      }
      // 一样的表示读取属性
    } else if ( hooks && "get" in hooks && notxml && (ret = hooks.get( elem, name )) !== null ) {
      return ret;

    } else {
      // 最终的读取属性getAttribute
      ret = elem.getAttribute( name );

      // 读取属性我们不返回null，返回undefined
      return ret === null ?
        undefined :
        ret;
    }
  },

  removeAttr: function( elem, value ) {
    var propName, attrNames, name, l,
      i = 0;

    if ( value && elem.nodeType === 1 ) {
      // 	rspace = /\s+/,
      // 按照空格开划分， 支持一次移除多个属性
      attrNames = value.toLowerCase().split( rspace );
      l = attrNames.length;

      for ( ; i < l; i++ ) {
        name = attrNames[ i ];

        if ( name ) {
          propName = jQuery.propFix[ name ] || name;

          // 在移除属性之前，先将属性设置为空字符串 解决webkit无法移除style的问题
          // See #9699 for explanation of this approach (setting first, then removal)
          jQuery.attr( elem, name, "" );
          elem.removeAttribute( getSetAttribute ? name : propName );

          // Set corresponding property to false for boolean attributes
          if ( rboolean.test( name ) && propName in elem ) {
            elem[ propName ] = false;
          }
        }
      }
    }
  },

  // 这里的attrHooks并不完整
  // 还包括 tabIndex、width、height、contenteditable
  // 主要实现针对几个特殊html属性值的读取和设置，因为特殊
  // 并且每个属性都有set和get对应方法
  attrHooks: {
    type: {
      set: function( elem, value ) {
        // We can't allow the type property to be changed (since it causes problems in IE)
        // 在ie9以下如果检测到要修改的事input或者button并且还有父元素，那么不允许修改type
        // rtype = /^(?:button|input)$/i, 
        if ( rtype.test( elem.nodeName ) && elem.parentNode ) {
          jQuery.error( "type property can't be changed" );
        } else if ( !jQuery.support.radioValue && value === "radio" && jQuery.nodeName(elem, "input") ) {
          //  之前jQuery.support.radioValue如果为false表示设置type为radio会让value丢失，那么就应该先备份，再设置tyeop，最后再恢复value
          // Setting the type on a radio button after the value resets the value in IE6-9
          // Reset value to it's default in case type is set after value
          // This is for element creation
          var val = elem.value;
          elem.setAttribute( "type", value );
          if ( val ) {
            elem.value = val;
          }
          return value;
        }
      }
    },
    // Use the value property for back compat
    // Use the nodeHook for button elements in IE6/7 (#1954)
    // 在这里优先使用nodeHook来进行设置或读取
    value: {
      get: function( elem, name ) {
        if ( nodeHook && jQuery.nodeName( elem, "button" ) ) {
          return nodeHook.get( elem, name );
        }
        return name in elem ?
          elem.value :
          null;
      },
      set: function( elem, value, name ) {
        if ( nodeHook && jQuery.nodeName( elem, "button" ) ) {
          return nodeHook.set( elem, value, name );
        }
        // Does not return so that setAttribute is also used
        elem.value = value;
      }
    }
  },

  propFix: {
    tabindex: "tabIndex",
    readonly: "readOnly",
    "for": "htmlFor",
    "class": "className",
    maxlength: "maxLength",
    cellspacing: "cellSpacing",
    cellpadding: "cellPadding",
    rowspan: "rowSpan",
    colspan: "colSpan",
    usemap: "useMap",
    frameborder: "frameBorder",
    contenteditable: "contentEditable"
  },

  // dom属性处理与html属性处理非常类似
  // 某些类型节点不处理
  // 优先采用有hook的方法，再采用通用方法
  prop: function( elem, name, value ) {
    var ret, hooks, notxml,
      nType = elem.nodeType;

    // don't get/set properties on text, comment and attribute nodes
    // 与html属性处理一致，忽略注释、文本、属性节点
    if ( !elem || nType === 3 || nType === 8 || nType === 2 ) {
      return;
    }

    notxml = nType !== 1 || !jQuery.isXMLDoc( elem );

    if ( notxml ) {
      // Fix name and attach hooks
      name = jQuery.propFix[ name ] || name;
      hooks = jQuery.propHooks[ name ];
    }

    if ( value !== undefined ) {
      if ( hooks && "set" in hooks && (ret = hooks.set( elem, value, name )) !== undefined ) {
        return ret;

      } else {
        return ( elem[ name ] = value );
      }

    } else {
      if ( hooks && "get" in hooks && (ret = hooks.get( elem, name )) !== null ) {
        return ret;

      } else {
        return elem[ name ];
      }
    }
  },

  propHooks: {
    tabIndex: {
      get: function( elem ) {
        // elem.tabIndex doesn't always return the correct value when it hasn't been explicitly set
        // http://fluidproject.org/blog/2008/01/09/getting-setting-and-removing-tabindex-values-with-javascript/
        // 上面的链接说明了直接读取tabIndex并不能保证返回正确值，需要使用属性节点来保证正确性
        // 如果有返回值则以10进制形式返回，否则再检测元素是不是可点击或者可聚焦的元素，如果满足条件则默认其tabIndex为0， 其余的应该是undefined
        var attributeNode = elem.getAttributeNode("tabindex");
        // 	rfocusable = /^(?:button|input|object|select|textarea)$/i,
        // 	rclickable = /^a(?:rea)?$/i,
        return attributeNode && attributeNode.specified ?
          parseInt( attributeNode.value, 10 ) :
          rfocusable.test( elem.nodeName ) || rclickable.test( elem.nodeName ) && elem.href ?
            0 :
            undefined;
      }
    }
  }
});
```

几个辅助函数
``` js
// 用于检测属性值是不是属于boolean类型的
rboolean = /^(?:autofocus|autoplay|async|checked|controls|defer|disabled|hidden|loop|multiple|open|readonly|required|scoped|selected)$/i;
// Hook for boolean attributes
boolHook = {
  get: function( elem, name ) {
    // 先利用jQuery.prop去读取属性值， 入锅味true 直接返回属性名称
    // 否则再利用getAttributeNode方法去读取属性值
    // 辱国属性是true则返回属性名，否则返回undefined
    var attrNode,
      property = jQuery.prop( elem, name );
    return property === true || typeof property !== "boolean" && ( attrNode = elem.getAttributeNode(name) ) && attrNode.nodeValue !== false ?
      name.toLowerCase() :
      undefined;
  },
  set: function( elem, value, name ) {
    var propName;
    if ( value === false ) {
      // Remove boolean attributes when set to false
      jQuery.removeAttr( elem, name );
    } else {
      // 如果属性是某些特殊属性，还需要设置其相关属性一起为true
      // value is true since we know at this point it's type boolean and not false
      // Set boolean attributes to the same name and set the DOM property
      propName = jQuery.propFix[ name ] || name;
      if ( propName in elem ) {
        // Only set the IDL specifically if it already exists on the element
        elem[ propName ] = true;
      }

      elem.setAttribute( name, name.toLowerCase() );
    }
    return name;
  }
};


// jquery初始化时，如果setAttribute为false，则会保留nodeHook
// 目的就是来替代setAttribute和getAttribute、removeAttribute
// 主要是通过setAttributeNode与getAttributeNode的封装来达到setAttribute与getAttribute的效果
nodeHook = jQuery.valHooks.button = {
  get: function( elem, name ) {
    var ret;
    ret = elem.getAttributeNode( name );
    return ret && ( fixSpecified[ name ] ? ret.nodeValue !== "" : ret.specified ) ?
      ret.nodeValue :
      undefined;
  },
  set: function( elem, value, name ) {
    // Set the existing or create a new attribute node
    var ret = elem.getAttributeNode( name );
    if ( !ret ) {
      ret = document.createAttribute( name );
      elem.setAttributeNode( ret );
    }
    return ( ret.nodeValue = value + "" );
  }
};
```