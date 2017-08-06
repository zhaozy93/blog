# chapter4 数据缓存Data

data部分提供为每个dom或js对象添加数据的功能。可能不是特别常用，但其为其他模块提供了支持，比如动画、队列、样式操作等模块。

这部分代码也并不复杂。全局有一个cache对象内部存储dom元素对应的对象，每个dom都会添加一个唯一标签，通过标签，cache内的对象可以与dom一一对应。为js对象添加则更加简单，为对象添加一个不会被占用的属性名即可。

先看整体结构，比较清晰，检测元素是不是允许设置数据，允许则先计算一些标签信息，再设置属性、最后再增加移除数据的接口。

```
jQuery.extend({
	cache: {}, // cache就是全局盛放dom数据的容器
	uuid: 0,   // dom唯一标签的计数器
	expando: "jQuery" + ( jQuery.fn.jquery + Math.random() ).replace( /\D/g, "" ),
	noData: { },
	hasData: function( elem ) // 检测元素是否有对应的数据
	data: function( elem, name, data, pvt /* Internal Use Only */ ) // 设置或者读取数据
	removeData: function( elem, name, pvt /* Internal Use Only */ ) // 移除数据
	_data: function( elem, name, data )  // 设置读取内部数据
	acceptData: function( elem ) // 检测是不是允许设置数据
});

jQuery.fn.extend({
	data: function( key, value )    // 对外接口，设置、读取数据

	removeData: function( key )  // 移除数据
});
```

下面还是一块一块看，先从jQuery.extend来看， 一般它都是jQuery.fn.extend的基础。

## 源码 jQuery.extend
### acceptData
主要是防止三类dom，embed、applet和flash，这三类不支持扩展属性，因此无法设置标签，没有标签就没办法让数据和dom进行统一，标签就像id一样是唯一的，是dom和cache沟通的桥梁。
```
// 判断一个dom元素是不是可以设置数据
acceptData: function( elem ) {
  // 如果不是dom则一定可以，直接返回true
  if ( elem.nodeName ) {
    // 看dom的name是不是在noData里面，如果再则取出来，不再则match=undefined
    var match = jQuery.noData[ elem.nodeName.toLowerCase() ];
    // 如果dom是embed或者applet则return false， 如果是object恰巧是falsh则也返回false
    if ( match ) {
      return !(match === true || elem.getAttribute("classid") !== match);
    }
  }
  return true;
}

noData: {
  "embed": true,
  // Ban all objects except for Flash (which handle expandos)
  "object": "clsid:D27CDB6E-AE6D-11cf-96B8-444553540000",
  "applet": true
},
```

### data 为js对象或者dom设置元素
```
data: function( elem, name, data, pvt /* Internal Use Only */ ) {
  // 不支持设置 直接return
  if ( !jQuery.acceptData( elem ) ) {
    return;
  }

  var privateCache, thisCache, ret,
    internalKey = jQuery.expando,
    getByName = typeof name === "string",

    // 判断是不是dom元素
    isNode = elem.nodeType,

    // 如果是dom元素则cache就是jQuery.cache这个对象，否则就是js对象本身
    // 分开设置是因为 js本身的数据可以直接被GC收集掉
    cache = isNode ? jQuery.cache : elem,

    // 如果是dom元素，则返回 elem[ internalKey ] 有可能是undefined
    // 如果是js对象，之前如果设置过数据则elem[ internalKey ]为真，返回internalKey， 否则false
    id = isNode ? elem[ internalKey ] : elem[ internalKey ] && internalKey,
    isEvents = name === "events";

  // 如果是要获取某个元素的数据，但没有数组则直接返回
  if ( (!id || !cache[id] || (!isEvents && !pvt && !cache[id].data)) && getByName && data === undefined ) {
    return;
  }

  if ( !id ) {
    // 如果之前没设置过数据
    // dom元素的话就增加一个固定属性名(internalKey的值)，值为自增的jQuery.uuid
    // 果然是js对象，则id就为internalKey的值
    if ( isNode ) {
      elem[ internalKey ] = id = ++jQuery.uuid;
    } else {
      id = internalKey;
    }
  }

  // 刚才也说过了，cache会是一个引用。jQuery.cache或者js元素本身
  if ( !cache[ id ] ) {
    cache[ id ] = {};

    // jQuery这部分数据不应该被字符串序列化时出现。
    if ( !isNode ) {
      cache[ id ].toJSON = jQuery.noop;
    }
  }

  // 如果传入的是对象，则把传入的内容扩展到现有的数据上
  // pvt控制是内部数据还是自定义数据
  if ( typeof name === "object" || typeof name === "function" ) {
    if ( pvt ) {
      cache[ id ] = jQuery.extend( cache[ id ], name );
    } else {
      cache[ id ].data = jQuery.extend( cache[ id ].data, name );
    }
  }

  privateCache = thisCache = cache[ id ];

  // 如果是不是pvt，则代表是自定义数据，则cache的真正目的地是cache.data
  // 过了这里之后也就是thisCache一定指向正确的cache对象。
  if ( !pvt ) {
    if ( !thisCache.data ) {
      thisCache.data = {};
    }

    thisCache = thisCache.data;
  }

  // camelCase: function( string ) {
  //    return string.replace( rmsPrefix, "ms-" ).replace( rdashAlpha, fcamelCase );
  //  },
  // rmsPrefix = /^-ms-/,
  // rdashAlpha = /-([a-z]|[0-9])/ig,
  // fcamelCase = function( all, letter ) {
  //    return ( letter + "" ).toUpperCase();
  //  },
  // 只要data为真就用驼峰的方式设置名字和数据。
  if ( data !== undefined ) {
    thisCache[ jQuery.camelCase( name ) ] = data;
  }

  // Users should not attempt to inspect the internal events object using jQuery.data,
  // it is undocumented and subject to change. But does anyone listen? No.
  // 一段非常骚的注释：
  // 用户不应该试图使用jQuery.data去查看内部events对象,
  // events是没有注解并且可能会更改的。但是谁会听呢？ 没人听。
  // 如果用户要查看events属性，并且不是自定义的那部分，则把真正的events返给用户.
  // 很无奈的注释
  if ( isEvents && !thisCache[ name ] ) {
    return privateCache.events;
  }

  // 通过查看为转换驼峰和转换驼峰的name来查看数据，有就返回， 
  // 如果都不是getByName则直接返回整个缓存的数据
  if ( getByName ) {
    // First Try to find as-is property data
    ret = thisCache[ name ];
    // Test for null|undefined property data
    if ( ret == null ) {
      // Try to find the camelCased property
      ret = thisCache[ jQuery.camelCase( name ) ];
    }
  } else {
    ret = thisCache;
  }

  return ret;
},
```

### _data惨不忍睹，专用于设置内部数据
```
_data: function( elem, name, data ) {
  return jQuery.data( elem, name, data, true );
},
```

### hasData
判断一个dom元素或者js对象有没有通过jQuery设置过数据，当然也不能是空对象
```
hasData: function( elem ) {
  elem = elem.nodeType ? jQuery.cache[ elem[jQuery.expando] ] : elem[ jQuery.expando ];
  return !!elem && !isEmptyDataObject( elem );
},
```

### removeData 移除数据
```
removeData: function( elem, name, pvt /* Internal Use Only */ ) {
  // 压根不会设置数据 直接return
  if ( !jQuery.acceptData( elem ) ) {
    return;
  }

  // 这段var就获取了dom的id或者internalKey，也就是cache的属性名
  var thisCache, i, l,
    internalKey = jQuery.expando,
    isNode = elem.nodeType,
    cache = isNode ? jQuery.cache : elem,
    id = isNode ? elem[ internalKey ] : internalKey;

  // 如果没有这个属性名对应的值，则直接返回
  if ( !cache[ id ] ) {
    return;
  }

  if ( name ) {
    // 根据pvt来修正 以让thisCache指向正确的cache对象
    thisCache = pvt ? cache[ id ] : cache[ id ].data;
    if ( thisCache ) {
      // 要将name变成数组的形式，不是数组则转换格式
      if ( !jQuery.isArray( name ) ) {
        // 如果恰巧是单个name字符串，恰巧也存在那么直接name = [name]
        if ( name in thisCache ) {
          name = [ name ];
        } else {
          // 再试试能不能先驼峰转换完了有， 有则[name]， 没有则安空格split成数组
          name = jQuery.camelCase( name );
          if ( name in thisCache ) {
            name = [ name ];
          } else {
            name = name.split( " " );
          }
        }
      }
      // 循环数组delete属性
      for ( i = 0, l = name.length; i < l; i++ ) {
        delete thisCache[ name[i] ];
      }

      // 如果删除之后不是空对象则return
      if ( !( pvt ? isEmptyDataObject : jQuery.isEmptyObject )( thisCache ) ) {
        return;
      }
    }
  }

  // 如果name有值，那么正常应该已经return了， 走到这里证明name的数据移除完，剩余的是个空数据容器了
  // 如果是删除自定义属性，则要删除自定义数据，同时删除自定义数据的容器.data
  if ( !pvt ) {
    delete cache[ id ].data;
    // 如果剩下的不是空对象则return
    if ( !isEmptyDataObject(cache[ id ]) ) {
      return;
    }
  }

  // 优雅降级的表现 能delete则delete，否则设置null
  if ( jQuery.support.deleteExpando || !cache.setInterval ) {
    delete cache[ id ];
  } else {
    cache[ id ] = null;
  }

  // 优雅降级的表现， 在dom上面， 如果支持delete扩展属性就delete，其次removeAttribute，最后是设置为null
  if ( isNode ) {
    if ( jQuery.support.deleteExpando ) {
      delete elem[ internalKey ];
    } else if ( elem.removeAttribute ) {
      elem.removeAttribute( internalKey );
    } else {
      elem[ internalKey ] = null;
    }
  }
},
```

## 源码 jQuery.fn.extend

### removeData 
先看简单的
很明显这是为$()获取到的数组应用each来为每个元素都执行removeData操作的，同时都是移除自定义数据
```
removeData: function( key ) {
  return this.each(function() {
    jQuery.removeData( this, key );
  });
}
```

### data
```
data: function( key, value ) {
  var parts, attr, name,
    data = null;
  // 如果key是undefined，返回第一个dom对应的缓存数据
  if ( typeof key === "undefined" ) {
    // 当前获取到了元素，长度才会大于0
    if ( this.length ) {
      // data设为第一个元素对应缓存数据
      data = jQuery.data( this[0] );
      // 是dom元素 且没有对应的data中没有parseAttrs属性
      if ( this[0].nodeType === 1 && !jQuery._data( this[0], "parsedAttrs" ) ) {
        // 获取dom元素所有属性的集合，
        attr = this[0].attributes;
        for ( var i = 0, l = attr.length; i < l; i++ ) {
          name = attr[i].name;
          if ( name.indexOf( "data-" ) === 0 ) {
            name = jQuery.camelCase( name.substring(5) );
            dataAttr( this[0], name, data[ name ] );
          }
        }
        // 在遍历完dom属性之后为dom设置属性名为parsedAttrs，值为true的内部属性
        jQuery._data( this[0], "parsedAttrs", true );
      }
    }
    return data;
  } else if ( typeof key === "object" ) {
    // 如果key是一个对象，则为每一个dom元素都执行jQuery.data来设置属性
    return this.each(function() {
      jQuery.data( this, key );
    });
  }

  // key既不是undefined也不是object， 那应该是string了
  parts = key.split(".");
  parts[1] = parts[1] ? "." + parts[1] : "";
  // 如果有key，且value是undefined则应该是获取数据
  if ( value === undefined ) {
    // 尝试用getData来获取数据， 暂时没看这个东西 太复杂
    data = this.triggerHandler("getData" + parts[1] + "!", [parts[0]]);

    // 上一步获取失败了，再从缓存数据里获取，再尝试dataAttr获取
    if ( data === undefined && this.length ) {
      data = jQuery.data( this[0], key );
      // 如果在缓存数据里找到了，下面这行是不会执行的，稍后会看到分析
      data = dataAttr( this[0], key, data );
    }
    // 如果这次找到data，直接返回，走则去掉parts[1]，再执行一个data尝试一次
    return data === undefined && parts[1] ?
      this.data( parts[0] ) :
      data;

  } else {
    // 如果有key和value 则为每一个元素都设置缓存数据，依旧是triggerHandler，暂时看不懂
    return this.each(function() {
      var self = jQuery( this ),
        args = [ parts[0], value ];

      self.triggerHandler( "setData" + parts[1] + "!", args );
      jQuery.data( this, key, value );
      self.triggerHandler( "changeData" + parts[1] + "!", args );
    });
  }
},
```
## 外部源码 不归属于jQuery.extend 和 fn.extend
### dataAttr
```
function dataAttr( elem, key, data ) {
  // 如果data不是undefined则直接返回undefined
	if ( data === undefined && elem.nodeType === 1 ) {
    // rmultiDash = /([A-Z])/g;
    // 修正变量名从驼峰到data-a-b-c
		var name = "data-" + key.replace( rmultiDash, "-$1" ).toLowerCase();
		data = elem.getAttribute( name );
		if ( typeof data === "string" ) {
      // 如果是字符串，尝试把字符串还原本身格式，true、false、null、数字、对象， 如果都不能还原则当作普通字符串返回
			try {
        // rbrace = /^(?:\{.*\}|\[.*\])$/
				data = data === "true" ? true :
				data === "false" ? false :
				data === "null" ? null :
				jQuery.isNumeric( data ) ? parseFloat( data ) :
					rbrace.test( data ) ? jQuery.parseJSON( data ) :
					data;
			} catch( e ) {}
			// 顺便再把刚刚解析出来的数据data-*的数据混存到jQuery.cache[id].data里面
			jQuery.data( elem, key, data );
		} else {
      // 只要data不是字符串(可能是null)那就是没找到，直接设为undefined返回
			data = undefined;
		}
	}
	return data;
}
```

### 
```
cleanData: function( elems ) {
  var data, id,
    cache = jQuery.cache,
    special = jQuery.event.special,
    deleteExpando = jQuery.support.deleteExpando;

  for ( var i = 0, elem; (elem = elems[i]) != null; i++ ) {
    if ( elem.nodeName && jQuery.noData[elem.nodeName.toLowerCase()] ) {
      continue;
    }

    id = elem[ jQuery.expando ];

    if ( id ) {
      data = cache[ id ];
      if ( data && data.events ) {
        // 遍历events内的数据
        for ( var type in data.events ) {
          //按两种方式移除事件监听或代理， 还不完全清楚event部分
          if ( special[ type ] ) {
            jQuery.event.remove( elem, type );
          } else {
            jQuery.removeEvent( elem, type, data.handle );
          }
        }
        // Null the DOM reference to avoid IE6/7/8 leak (#7054)
        // 设置为null以防ie678内存泄漏
        if ( data.handle ) {
          data.handle.elem = null;
        }
      }
      // 在dom上移除为缓存数据创建的标签
      if ( deleteExpando ) {
        delete elem[ jQuery.expando ];
      } else if ( elem.removeAttribute ) {
        elem.removeAttribute( jQuery.expando );
      }
      // 删除缓存数据占用的那个对象
      delete cache[ id ];
    }
  }
}
```

## 总结
对于数据缓存其实分两个部分 一个是dom，一个是js对象。 在整个页面初始化时，会得到一个唯一的jQuery开头的id字符串，这个id将作为之后数据缓存的标签id。

- js对象比较简单，直接在js对象的内部增加一个属性名为标签id的对象。
- dom也是为dom元素增加一个一个属性名为标签id的属性，但值为一个从0开始的全局自增的数字，同时会为jQuery.cache这个全局对象增加一个相同数字的属性，这样也为dom元素找到了一个盛放数据的地方。

无论dom还是js对象的缓存数据内部都分自定义数据和内部数据。 自定义数据放在data这个子属性内，内部数据直接挂在缓存对象上。

在代码里有一个地方处理的很巧妙，尽管分为dom和js两部分，但实现的时候先获取到承载数据的那个对象，只要获取的缓存数据的那个对象，后面的处理逻辑就都一样了，其实都是js object的操作了。