# chapter5 队列Queue

队列是一种最基本的数据结构，队列就是一种特殊的线性表,只允许在队尾插入，队头删除，也就是我们常说的先进先出，后进后出。jQuery也实现了这样一种结构，通过`shift()`和`push()`对Array进行操作，来模拟队头删除和队尾插入。队列模块最重要的功能之一就是为后面的动画模块提供支持，毕竟很多动画是要按顺序执行，那这时候队列的作用就非常重要了。

## 整体结构
我们在这里这里先讨论一个基本问题，就是每次每个模块好像都有jQuery.extend 和 jQuery.fn.extend，为什么这样呢。
我们可以这样考虑这个问题，jQuery.extend是扩展到jQuery对象上面的，jQuery.fn.extend是扩展到jQuery原型链上面的。 当我们$()获得的对象实例是继承了原型链，但是$并没有继承原型链啊。也就是说jQuery.extend里面的方法，我们可以直接$.queue使用，但是jQuery.fn.extend必须是$().queue使用，虽然函数名是一样的，但完全不是一个东西啊。因为实例之后虽然原型链的东西继承了，但是原来那一套东西是不会跟着过去的。更详细的解释可以看另一篇文章[关于new的认识、附带lazyman和currying的理解](https://github.com/zhaozy93/blog/issues/6)。 举个例子
```
function alone(){};
alone.queue = 'abc';
alone.prototype.queue1 = 'qwer';
let a = new alone();
console.log(alone.queue)   // "abc"
console.log(alone.queue1)   // undefined
console.log(a.queue)   // undefined
console.log(a.queue1)   // "qwer"
```

这部分代码也没啥幺蛾子， 整体结构也比较容易理解，维护一个盛放队列的容器，并且对队列实现入列和出列的操作即可。
```
jQuery.extend({
  _mark: function( elem, type ) {},
  _unmark: function( force, elem, type ) {},
  // 函数入列，并返回队列
  queue: function( elem, type, data ) {},
  // 函数出列，并执行
  dequeue: function( elem, type ) {}
});
jQuery.fn.extend({
  // 取出函数队列，或者函数入列
  queue: function( type, data ) {},
  // 函数出列，并执行
  dequeue: function( type ) {},
  // 清空队列
  clearQueue: function( type ) {},
  // 观察函数队列和计数器是否完成，并返回异步队列的只读副本
  promise: function( type, object ) {}
});
```

## 源码 jQuery.extend
### queue
在这里我们就可以发现，队列模块是依靠数据缓存模块来维护队列的，将队列分散到各个dom元素或者js对象上，用一个单独的内部属性来作为盛放队列的容器。
```
queue: function( elem, type, data ) {
  var q;
  if ( elem ) {
    // 修正队列名称，并提供默认队列名称
    type = ( type || "fx" ) + "queue";
    q = jQuery._data( elem, type );
    
    // 如果没提供入列函数，就直接退出
    if ( data ) {
      // 如果之前没有队列缓存数据，或者传进来的是一个数组，那就将data变为数组然后进行整体存储或者替换， 否则就只是普通的数组push作为入列
      if ( !q || jQuery.isArray(data) ) {
        q = jQuery._data( elem, type, jQuery.makeArray(data) );
      } else {
        q.push( data );
      }
    }
    return q || [];
  }
},
```

### dequeue
用于从队列中取出队头的函数并执行
```
dequeue: function( elem, type ) {
  // 修正type名 并且查询当前type的队列，并取出列头第一个方法
  type = type || "fx";
  var queue = jQuery.queue( elem, type ),
    fn = queue.shift(),
    hooks = {};

  // 如果列头是一个inprogress，则再重新取出一个，这是动画正在执行的占位符号
  if ( fn === "inprogress" ) {
    fn = queue.shift();
  }

  if ( fn ) {
    // 如果当前type是fx即默认代表动画的队列则再头部插入一个inprogress表示占位符
    if ( type === "fx" ) {
      queue.unshift( "inprogress" );
    }

    // 新建了一个.run结尾的缓存数据，并把一个空对象存进去，引用类型哦
    jQuery._data( elem, type + ".run", hooks );
    // 执行这个刚刚取出来的函数， 为他增加两个变量，
    // 第一个封装过的next， 在调用next的时候会再执行jQuery.dequeue就可以保证衔接继续一直出列
    // hooks可以存储数据，然后之后的可以读取，相当于是队列前后执行函数进行数据沟通的桥梁
    fn.call( elem, function() {
      jQuery.dequeue( elem, type );
    }, hooks );
  }

  if ( !queue.length ) {
    // 如果队列内没有fn了，则应该清理一下缓存数据
    jQuery.removeData( elem, type + "queue " + type + ".run", true );
    handleQueueMarkDefer( elem, type, "queue" );
  }
}
```

### _mark、_unmark
在jQuery中有很多只在内部调用的方法，其名称都是以_开始的。
_mark 和 _unmark特别像计数器，实际上它们也只是为动画部分服务。
```
// _mark用于增加计数器个数
_mark: function( elem, type ) {
  if ( elem ) {
    // 修正type
    type = ( type || "fx" ) + "mark";
    // 新建或者为原本的缓存数据+1
    jQuery._data( elem, type, (jQuery._data( elem, type ) || 0) + 1 );
  }
},
// _unmark用于减少计数器个数
_unmark: function( force, elem, type ) {
  // 可以看出来如果接收的是三个参数，那么顺序是正确的
  // 如果接收的参数是两个，那么第一个是元素，第二个是type，force就是false
  if ( force !== true ) {
    type = elem;
    elem = force;
    force = false;
  }
  if ( elem ) {
    // 修正type
    type = type || "fx";
    var key = type + "mark",
    // 如果force是true则直接让count变为0， 否则先读区当前的count再减一
      count = force ? 0 : ( (jQuery._data( elem, key ) || 1) - 1 );
    // 有count值就设置count，没有的话就做删除数据的操作
    if ( count ) {
      jQuery._data( elem, key, count );
    } else {
      jQuery.removeData( elem, key, true );
      handleQueueMarkDefer( elem, type, "mark" );
    }
  }
},
```

## 源码 jQuery.fn.extend
这一部分可能是更常用的，直接获取完元素之后使用，可能更方便一点。

### queue
功能肯定与jQuery.extend中的同名方法类似。
```
queue: function( type, data ) {
  // 修正参数， 如果第一个参数不是string则证明只有一个参数，那么data就是要入列的函数， type则就是默认的fx
  if ( typeof type !== "string" ) {
    data = type;
    type = "fx";
  }
  // 如果修正之后都没有data，那就是一个参数也没有，那就是读取第一个元素的fx的数据
  if ( data === undefined ) {
    return jQuery.queue( this[0], type );
  }
  // 只要有data则就为每一个元素都 入列data方法
  // 同时 如果是默认的fx动画类型，并且第一个元素不是inprogress的占位符，则自动执行一次出列
  return this.each(function() {
    var queue = jQuery.queue( this, type, data );
    if ( type === "fx" && queue[0] !== "inprogress" ) {
      jQuery.dequeue( this, type );
    }
  });
},
```

### dequeue
为每个元素都执行一次出列
```
dequeue: function( type ) {
  return this.each(function() {
    jQuery.dequeue( this, type );
  });
},
```

### delay
```
delay: function( time, type ) {
  // time 有可能是动画的几个string类型
  // slow: 600, fast: 200, _default: 400
  // 先修正time和type
  time = jQuery.fx ? jQuery.fx.speeds[ time ] || time : time;
  type = type || "fx";

  // 入列一个延迟函数， 这个函数里面有next和hooks，这个可以去看之前那个dequeu的方法
  // 执行函数传入的第一个参数会是 再次调用一个dequeue以达到next继续调用的效果
  return this.queue( type, function( next, hooks ) {
    var timeout = setTimeout( next, time );
    hooks.stop = function() {
      clearTimeout( timeout );
    };
  });
},
```

### clearQueue
通过把队列设置为[]来起到置空队列的效果
```
clearQueue: function( type ) {
  return this.queue( type || "fx", [] );
},
```

### promise
这个方法的目的不是特别明确，只能看个代码的大概
```
promise: function( type, object ) {
  // 如果只有一个元素， 那么应该是object， 但是整个函数内部发现object没有被再次使用到。。。。真神奇
  if ( typeof type !== "string" ) {
    object = type;
    type = undefined;
  }
  // 修正type
  type = type || "fx";
  ／/ 新建一个异步队列
  var defer = jQuery.Deferred(),
    elements = this,
    i = elements.length,
    count = 1,
    deferDataKey = type + "defer",
    queueDataKey = type + "queue",
    markDataKey = type + "mark",
    tmp;
  function resolve() {
    if ( !( --count ) ) {
      defer.resolveWith( elements, [ elements ] );
    }
  }
  // 遍历当前this的所有元素， 为每个元素的type + "defer"缓存数据增加一个resolve方法，
  // 同时这里有个计步器count++
  // 这个resolve函数也比较有意思，一定要count为1的时候才会执行
  while( i-- ) {
    if (( tmp = jQuery.data( elements[ i ], deferDataKey, undefined, true ) ||
        ( jQuery.data( elements[ i ], queueDataKey, undefined, true ) ||
          jQuery.data( elements[ i ], markDataKey, undefined, true ) ) &&
        jQuery.data( elements[ i ], deferDataKey, jQuery.Callbacks( "once memory" ), true ) )) {
      count++;
      tmp.add( resolve );
    }
  }
  resolve();
  // 返回了这个队列的只读副本
  return defer.promise();
}
```

进过查证，这里有一个bug，最后一行其实是 `return defer.promise(object)`,这样这个异步只读副本就附加到object对象上

## handleQueueMarkDefer
负责检测元素的关联的队列和计数器是否完成，如果完成了就调用一次.promise，这样当前的count就会减一，以来继续判断是不是全部执行完毕。
```
function handleQueueMarkDefer( elem, type, src ) {
  // 新建局部变量，并且获取defer函数
  var deferDataKey = type + "defer",
    queueDataKey = type + "queue",
    markDataKey = type + "mark",
    defer = jQuery._data( elem, deferDataKey );
    // 根据传入的参数来决定 什么样的条件， 
    // 如果有关联的回调函数、计数器和队列才会进入
  if ( defer &&
    ( src === "queue" || !jQuery._data(elem, queueDataKey) ) &&
    ( src === "mark" || !jQuery._data(elem, markDataKey) ) ) {
    // 队列是空的了，并且计步器也是0了，那么久移除关联的回调函数，并且调用一次上面的promise
    setTimeout( function() {
      if ( !jQuery._data( elem, queueDataKey ) &&
        !jQuery._data( elem, markDataKey ) ) {
        jQuery.removeData( elem, deferDataKey, true );
        defer.fire();
      }
    }, 0 );
  }
}
```

## 总结
jQuery队列的内部实现也不复杂，就是规范一个固定出入方式的Array，这个Array还是由数据缓存部分提供基础支持。