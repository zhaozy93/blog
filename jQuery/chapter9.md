# chapter9 DOM遍历

在jQuery如日中天的年代，DOM操作才是页面更新的几乎唯一途径，如今在各大框架React, VUE, Angular的蚕食下，大家都在降低对DOM的直接操作，转为对数据的操作，由框架去决定如何更新DOM。但DOM的操作仍是一名前端最基础的要求。DOM操作的前提是获取DOM元素，除了常规的`getElementById, getElementsByClassName, $()`，另一种途径就是根据当前元素，通过目标元素与当前元素的关系如兄弟、孩子、祖先等来寻找目标元素，也就是现在要讲的DOM遍历。

下面10个方法是jQuery暴漏出来的10个对兄弟、祖先、孩子元素获取的方法。可以看的出，10个方法都是通过一个模版函数挂载到jQuery.fn上的，模版函数内公用部分完成过滤、排序、驱虫操作，最后将结果再以jQuery对象的形式返回。


``` js
// 10个方法大部分依赖了dir、nth、sibling基础方法去寻找目标元素
jQuery.each({
  // 返回一个元素的父亲
  parent: function( elem ) {
    var parent = elem.parentNode;
    return parent && parent.nodeType !== 11 ? parent : null;
  },
  // 返回一个元素所有父亲
  parents: function( elem ) {
    return jQuery.dir( elem, "parentNode" );
  },
  // 返回一个元素所有父亲，知道碰到until元素为止
  parentsUntil: function( elem, i, until ) {
    return jQuery.dir( elem, "parentNode", until );
  },
  // 返回一个元素的下一个兄弟
  next: function( elem ) {
    return jQuery.nth( elem, 2, "nextSibling" );
  },
  // 返回一个元素的上一个兄弟
  prev: function( elem ) {
    return jQuery.nth( elem, 2, "previousSibling" );
  },
  // 返回一个元素的所有后面兄弟
  nextAll: function( elem ) {
    return jQuery.dir( elem, "nextSibling" );
  },
  // 返回一个元素的所有前面兄弟
  prevAll: function( elem ) {
    return jQuery.dir( elem, "previousSibling" );
  },
  // 返回一个元素的所有后面兄弟或者碰到某个元素为止
  nextUntil: function( elem, i, until ) {
    return jQuery.dir( elem, "nextSibling", until );
  },
  // 返回一个元素的所有前面兄弟或者
  prevUntil: function( elem, i, until ) {
    return jQuery.dir( elem, "previousSibling", until );
  },
  // 元素所有兄弟元素
  siblings: function( elem ) {
    return jQuery.sibling( elem.parentNode.firstChild, elem );
  },
  // 返回一个元素所有子元素
  children: function( elem ) {
    return jQuery.sibling( elem.firstChild );
  },
  // 返回元素的所有子元素包含各种注释、文本节点，如果是iframe元素则返回window对象
  contents: function( elem ) {
    return jQuery.nodeName( elem, "iframe" ) ?
      elem.contentDocument || elem.contentWindow.document :
      jQuery.makeArray( elem.childNodes );
  }
}, function( name, fn ) {
  jQuery.fn[ name ] = function( until, selector ) {
    var ret = jQuery.map( this, fn, until );
    // var runtil = /Until$/,
    // 如果函数名是不是以UNtil结尾的，表示不需要until参数
    // 修正参数
    if ( !runtil.test( name ) ) {
      selector = until;
    }
    // 如果有参数selector，并且参数selector是字符串，则执行过滤操作
    if ( selector && typeof selector === "string" ) {
      ret = jQuery.filter( selector, ret );
    }

    // 当过滤后仍有多于1个dom元素，还需要进行去重操作
    //   guaranteedUnique = {
    //    children: true,
    //    contents: true,
    //    next: true,
    //    prev: true
    //  }; 
    ret = this.length > 1 && !guaranteedUnique[ name ] ? jQuery.unique( ret ) : ret;

    //  rmultiselector = /,/
    //   rparentsprev = /^(?:parents|prevUntil|prevAll)/
    // 需要对特殊的操作进行倒序操作
    if ( (this.length > 1 || rmultiselector.test( selector )) && rparentsprev.test( name ) ) {
      ret = ret.reverse();
    }
    // 用找到的dom数组构建新的jQuery对象并返回
    return this.pushStack( ret, name, slice.call( arguments ).join(",") );
  };
});

jQuery.extend({
  // 过滤操作， 调用的事Sizzle接口
  filter: function( expr, elems, not ) {
    if ( not ) {
      expr = ":not(" + expr + ")";
    }
    return elems.length === 1 ?
      jQuery.find.matchesSelector(elems[0], expr) ? [ elems[0] ] : [] :
      jQuery.find.matches(expr, elems);
  },
  // 上面可以看出dir有三个值 parentNode、previousSibling、nextSibling
  dir: function( elem, dir, until ) {
    var matched = [],
      cur = elem[ dir ];
    // 循环调用elem[dir]， 然后一直查找，知道找到window 或者 不存在 或者 与 until一致则停止
    while ( cur && cur.nodeType !== 9 && (until === undefined || cur.nodeType !== 1 || !jQuery( cur ).is( until )) ) {
      if ( cur.nodeType === 1 ) {
        matched.push( cur );
      }
      cur = cur[dir];
    }
    return matched;
  },
  // result 有 2
  // dir 有nextSibling、parentNode、previousSibling、nextSibling
  nth: function( cur, result, dir, elem ) {
    result = result || 1;
    var num = 0;
    // 按照一个方向，向下走几次
    // result表示走几次
    // dir表示方向
    // 当result为0时，cur就是当前的元素本身， 那么会执行++num ==> 1，cur=cur[dir]
    // 在此进入循环 这一次 ++num ==> 2，也就会break，就达到了寻找紧邻的兄弟的目的
    for ( ; cur; cur = cur[dir] ) {
      if ( cur.nodeType === 1 && ++num === result ) {
        break;
      }
    }

    return cur;
  },

  // 获取元素n的所有后面兄弟元素，但是不能包含elem元素
  // 当n为一个父元素的第一个子元素，那么就是获得了当前父元素任意一个子元素的所有兄弟元素
  sibling: function( n, elem ) {
    var r = [];

    for ( ; n; n = n.nextSibling ) {
      if ( n.nodeType === 1 && n !== elem ) {
        r.push( n );
      }
    }

    return r;
  }
});
```