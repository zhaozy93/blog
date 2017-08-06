# chapter3 Deferred Object 异步队列

jQuery(1.5以后)的ajax这类异步任务都被重写，将ajax与回调函数解耦，实现高内聚低耦合。

jQuery的异步队列部分也比较简洁，代码主要集中在1040-1400之中， 主要实现jQuery.Callbacks 和 jQuery.Deferred与jQuery.when三个方法。

jQuery.Callbacks又是在1.7引入，返回一个链式工具对象，用于管理所有回调函数。 后两个是通过jQuery.extends添加上去，因此我们从jQuery.Callbacks开始，看jQuery的Deferred Object异步队列都做了什么，又是怎么实现的

## jQuery.Callbacks
先来看整个结构
jQuery.Callbacks定义了一堆变量 + 两个方法 add、fire与一个对象self，self内部属性也是一堆方法，最终整个函数以self作为返回值。
```
jQuery.Callbacks = function( flags ) {
  flags = flags ? ( flagsCache[ flags ] || createFlags( flags ) ) : {};
  var list = [],
    stack = [],
    memory,
    firing,
    firingStart,
    firingLength,
    firingIndex,
    add = function( args ) {},
    fire = function( context, args ){},
    self = {
      add: function() { },
      remove: function() { },
      has: function( fn ) { },
      empty: function() { },
      disable: function() { },
      disabled: function() { },
      lock: function() { },
      locked: function() { },
      fireWith: function( context, args ) { },
      fire: function() { },
      fired: function() { }
    };
  return self;
};
```

### createFlags
```
var flagsCache = {};
// 将字符串类型的flag转换为对象类型的
// 传入参数应该是 'once memory unique stopOnFalse'这样的格式，空格分割的单词
// 在一开始就设置了 flagsCache[ flags ] = object， 这样在之后操作object的时候，也起到了相同flag字符串的cache效果
// 使用正则匹配空格分割字符串 然后为object添加属性，整个过程没问题，perfect。
function createFlags( flags ) {
  var object = flagsCache[ flags ] = {},
    i, length;
  flags = flags.split( /\s+/ );
  for ( i = 0, length = flags.length; i < length; i++ ) {
    object[ flags[i] ] = true;
  }
  return object;
}
```

### add
这个方法比较简单，可以看出来就是将args数组内部元素一个个添加进list， list也是jQuery.callbacks内的一个数组元素。 如果args内部还有数组则递归解构， 同时，也会根据flag的标签来检测是不是重复性。
```
add = function( args ) {
  var i,
    length,
    elem,
    type,
    actual;
  for ( i = 0, length = args.length; i < length; i++ ) {
    elem = args[ i ];
    type = jQuery.type( elem );
    if ( type === "array" ) {
      // Inspect recursively
      add( elem );
    } else if ( type === "function" ) {
      // Add if not in unique mode and callback is not in
      if ( !flags.unique || !self.has( elem ) ) {
        list.push( elem );
      }
    }
  }
},
```

### fire
使用指定的上下文和参数来执行list中的回调函数
这里有一个变量特别乱 memory， 它的初始值是undefined， flags.memory表示在建callbacks时是否传入了memory这个flag。
```
var memory, firing, firingStart, firingLength, firingIndex;
fire = function( context, args ) {
  args = args || [];
  // 如果传入了memory flag则变量memory为 [context, args]
  memory = !flags.memory || [ context, args ];
  firing = true;
  firingIndex = firingStart || 0;
  firingStart = 0;
  firingLength = list.length;
  for ( ; list && firingIndex < firingLength; firingIndex++ ) {
    // 这里利用if内部来实现list的中的每个函数都被执行
    if (list[ firingIndex ].apply( context, args ) === false && flags.stopOnFalse ) {
      // 这里利用break来实现一旦有返回值是false，并且有stopOnFalse flag， 则后面的不再执行
      memory = true;
      break;
    }
  }
  // fireing就是表示是不是在执行回调函数
  firing = false;
  if ( list ) {
    if ( !flags.once ) {
      // 没有once flag
      if ( stack && stack.length ) {
        // 这一段有点意外，需要结合最后fireWith一起看，如果调用fireWith的时候，恰巧正在执行即firing为true， 则将这次之行指令压入stack, 待当前的任务执行完毕后，再调用stack内的上下文和参数 继续执行。
        memory = stack.shift();
        self.fireWith( memory[ 0 ], memory[ 1 ] );
      }
    } else if ( memory === true ) {
      // 有once 这个flag， 并且memory严格等于true就是刚才上面有stopOnFalse flag并且某一个回调函数返回值为false
      self.disable();
    } else {
      // 有once 这个flag，但刚刚没有回调函数返回false或者没有stopOnFalse flag
      list = [];
    }
  }
},
```
即是看完这段代码也对memory有点似懂非懂，再加上最后once这个flag。 不过可以看出来，是通过这几个flag来实现代码逻辑的分流操作。
基本可以看出来，只要list回调函数执行过之后，memory就被变成了[ , ]或true

### self.add()
这个add大概可以帮助我们理解刚才上面那个fire函数内部复杂的逻辑。

首先调用add()将要添加的方法添加进list,

1）紧接着检测当前回调函数式是不是在执行过程，如果在执行过程，就把firingLength重新计算一下，这样就会把刚刚添加进去的方法给执行了。

2）如果没有在执行，并且memory既不是undefined也不是true，则表示曾经执行过，那么更改firingStart为之前的list.length,这样这一次就会只执行刚添加进去的方法，同时立即调用fire方法，再把刚刚添加进去的方法立即执行， 通过参数可以看出来，刚才memory那个数组就是保存context和参数。
    // With memory, if we're not firing then
    // we should call right away, unless previous
    // firing was halted (stopOnFalse)
```
add: function() {
  if ( list ) {
    var length = list.length;
    add( arguments );
    if ( firing ) {
      firingLength = list.length;
    } else if ( memory && memory !== true ) {
      firingStart = length;
      fire( memory[ 0 ], memory[ 1 ] );
    }
  }
  return this;
},
```

### self.remove()
我们发现add和remove都有个前提就是list存在，而且不存在的话不会使用默认的[]，这表明，如果list不存在了，就代表执行过了，例如有once flag的情况， 表示只执行一次回调函数

remove对应于add表示要移除某些回调函数。
```
remove: function() {
  if ( list ) {
    var args = arguments,
      argIndex = 0,
      argLength = args.length;
    // 两个for循环来对比寻找要移除的对象在list中的位置。
    for ( ; argIndex < argLength ; argIndex++ ) {
      for ( var i = 0; i < list.length; i++ ) {
        if ( args[ argIndex ] === list[ i ] ) {
          
          if ( firing ) {
            // 如果当前正在执行回调函数list， 如果当前要剔除的元素位置在执行的范围内，将执行总长度减一
            // 如果已经执行过了，需要把当前执行的下标减一，因为list马上就要变化了。
            if ( i <= firingLength ) {
              firingLength--;
              if ( i <= firingIndex ) {
                firingIndex--;
              }
            }
          }
          // 移除完了元素之后，记得要将i减少一位，否则会跳过紧邻的下一个元素， 也正是这个原因，for循环那里是 i< list.length 而不是提前将list.length提取为变量
          list.splice( i--, 1 );
          // 如果有unique这个flag， 那么寻找到一次就不需要内部for循环继续了，可以直接跳过内部这个小循环，直接进入下一个要移除的函数的寻找匹配工作
          if ( flags.unique ) {
            break;
          }
        }
      }
    }
  }
  return this;
},
```

### self.has、empty、disable、disabled、lock、locked
```
// 检测一个函数是否在list内部， 指针类型检测，easy
has: function( fn ) {
  if ( list ) {
    var i = 0,
      length = list.length;
    for ( ; i < length; i++ ) {
      if ( fn === list[ i ] ) {
        return true;
      }
    }
  }
  return false;
},
// 将list置空 easy
empty: function() {
  list = [];
  return this;
},
// 将当前callbacks设为残疾状态， 通过几个关键指标设为undefined
disable: function() {
  list = stack = memory = undefined;
  return this;
},
// 检测当前callbacks是否为残疾状态
disabled: function() {
  return !list;
},
// 锁定当前callbacks， 通过设置stack为undefined， 同时如果memory为undefined或者严格等于true，顺便执行致残行为
// 如果是memory模式，则对list无影响，同时可以继续添加移除和执行回调函数，只是不能更改上下文和参数
lock: function() {
  stack = undefined;
  if ( !memory || memory === true ) {
    self.disable();
  }
  return this;
},
// 检测当前callbacks是否为锁定状态
locked: function() {
  return !stack;
},
```

### self.fireWith、fire、fired
```
// 使用给定的上下文来执行list中回调函数
fireWith: function( context, args ) {
  // 调用时会先检测stack，stack默认的是[], 只有self.locak()和self.disable()会置为undefined
  // 也就是说一个callbacks在锁定和残疾状态下是没有办法执行fireWith的
  if ( stack ) {
    if ( firing ) {
      if ( !flags.once ) {
        // 如果正在执行 且不是once模式，则设置stack
        stack.push( [ context, args ] );
      }
    } else if ( !( flags.once && memory ) ) {
      // 不在执行状态，且不是once flag或者没执行过， 则执行list回调函数
      fire( context, args );
    }
  }
  return this;
},
// 使用当前的上下文执行回调函数列表list
fire: function() {
  self.fireWith( this, arguments );
  return this;
},
// 查看当前list是否执行过
fired: function() {
  return !!memory;
}
```

## jQuery.Deferred
借助Callbacks 我们来看一下jQuery.Deferred是如何实现的， Deferred是jquery中很重要的一个延迟解决方案，它和promise思想很类似，也是jQuery.ajax的重要基础。因此，deferred就是jQuery的回调函数解决方案

我们都知道promise有一个很重要的特征就是有一个状态，初始值是pending,成功后变为resolved, 失败后变为rejected， 并且这个状态是单向不可逆的，同时只能更改一次，也就是说也么从pending-->resolved, 要么从pending-->rejected。 然后对应的会执行成功，或者失败的回调函数。

带着这个基础再来看deferred的实现。
```
function( func ) {
  // 首先内部定义了三个回调函数队列，利用callbacks，看得出来三个队列分别是成功(done),失败(fail),运行中(progress),具体名字可能和刚才的resolved、rejected。不一致，但意思可以理解。只是这里多了一个progress，运行中的回调函数队列。
  // 主意这里的flag， once表示只能执行一次 memory表示会记录之前执行的参数
  var doneList = jQuery.Callbacks( "once memory" ),
    failList = jQuery.Callbacks( "once memory" ),
    progressList = jQuery.Callbacks( "memory" ),
    // 定义初始值未pending
    state = "pending",
    // 定义了一个lists对象， 好像不知道具体用意，后面会有答案
    lists = {
      resolve: doneList,
      reject: failList,
      notify: progressList
    },
    // 定义了一个promise， interesting
    promise = {
      // 实现了done、fail、progress三个方法， 分别是向三个回调函数列表添加函数的方法。
      done: doneList.add,
      fail: failList.add,
      progress: progressList.add,
      // 返回当前state的方法
      state: function() {
        return state;
      },
      // isResolved 判断doneList这个Callbacks是不是执行过，执行过也大概意味着整个deferred成功了。应该也是判断deferred状态的方法。
      isResolved: doneList.fired,
      isRejected: failList.fired,

      // then 接受三个参数，但其实内部调用的反而是done、fail、progress，也是向三个回调函数列表添加函数的方法，只是一个简写而已。
      // 回去查看callbacks的代码能看出来 add方法return的是this， 也就是说这里可以支持链式写法
      then: function( doneCallbacks, failCallbacks, progressCallbacks ) {
        deferred.done( doneCallbacks ).fail( failCallbacks ).progress( progressCallbacks );
        return this;
      },
      // always是同时方法添加给done和fail，意味着不管失败还是成功都要调用这个方法
      always: function() {
        deferred.done.apply( deferred, arguments ).fail.apply( deferred, arguments );
        return this;
      },
      // 这个pipe略复杂
      // 首先 return是一个Deferred对象，并且调用了promise方法，大概意思就是复制一下这个对象， 后面就可以看到。
      // 内部用each来实现对done、fail和progress根据传入的三个参数向新产生的deferred对象添加回调函数，不过这个方法具体真没看懂具体的用意，查了一下1.8就被废掉了，只存在了0.1个小版本，就不费时间了。
      pipe: function( fnDone, fnFail, fnProgress ) {
        return jQuery.Deferred(function( newDefer ) {
          jQuery.each( {
            done: [ fnDone, "resolve" ],
            fail: [ fnFail, "reject" ],
            progress: [ fnProgress, "notify" ]
          }, function( handler, data ) {
            var fn = data[ 0 ],
              action = data[ 1 ],
              returned;
            if ( jQuery.isFunction( fn ) ) {
              deferred[ handler ](function() {
                returned = fn.apply( this, arguments );
                if ( returned && jQuery.isFunction( returned.promise ) ) {
                  returned.promise().then( newDefer.resolve, newDefer.reject, newDefer.notify );
                } else {
                  newDefer[ action + "With" ]( this === deferred ? newDefer : this, [ returned ] );
                }
              });
            } else {
              deferred[ handler ]( newDefer[ action ] );
            }
          });
        }).promise();
      },
      // 这里看着有点奇怪， 内部也有一个promise，其实内部这个promise指外层这个大的promise， 也就是说要把这个大的promise赋给传入的对象， 如果没有传入就新建一个， 有就把done、fail...属性赋给指定的obj，让他也有异步这个能力
      promise: function( obj ) {
        if ( obj == null ) {
          obj = promise;
        } else {
          for ( var key in promise ) {
            obj[ key ] = promise[ key ];
          }
        }
        return obj;
      }
    },
    // 这里用promise.promise得到一个半成品的promise，因为可以看出来现在为止，得到的这个deferred只能添加done\fail\progress函数，获取状态，但是不能更改state的能力
    deferred = promise.promise({}),
    key;
  // 这里就给刚刚声明的deferred增加了触发done、fail、progress函数队列的能力，还是利用Callbacks的基础方法
  for ( key in lists ) {
    deferred[ key ] = lists[ key ].fire;
    deferred[ key + "With" ] = lists[ key ].fireWith;
  }

  // 再为deferred添加成功失败必有得回调函数，更改状态的函数
  deferred.done( function() {
    state = "resolved";
  }, failList.disable, progressList.lock ).fail( function() {
    state = "rejected";
  }, doneList.disable, progressList.lock );

  // 如果整个deferred在声明的时候传入了一个函数，那就执行它
  if ( func ) {
    func.call( deferred, deferred );
  }

  // 一切都做好了，就返回这个deferred。
  return deferred;
},
```
这段代码看起来也不是很复杂，维护了三个回调函数队列和一个state，根据不同的方法去执行相对应的回调函数队列，并维护好state。一切的基础都在Callbacks。 可能有点不太明白的就是为什么要先声明一个promise再赋给deferred， 感觉直接声明deferred好像也没什么不可以。

## when
when这段代码更容易理解。 when就是接受多个deferred对象，当所有都resolved之后，认为整个when是成功的，执行成功回调函数，否则，任何一个deferred参数fail了，则整个when是fail的。

```
var sliceDeferred = [].slice;
when: function( firstParam ) {
  var args = sliceDeferred.call( arguments, 0 ),
    i = 0,
    length = args.length,
    pValues = new Array( length ),
    count = length,
    pCount = length,
    // 声明了deferred作为when的主deferred对象，如果只有一个参数则不新建
    deferred = length <= 1 && firstParam && jQuery.isFunction( firstParam.promise ) ?
      firstParam :
      jQuery.Deferred(),
    promise = deferred.promise();
  function resolveFunc( i ) {
    return function( value ) {
      args[ i ] = arguments.length > 1 ? sliceDeferred.call( arguments, 0 ) : value;
      if ( !( --count ) ) {
        deferred.resolveWith( deferred, args );
      }
    };
  }
  function progressFunc( i ) {
    return function( value ) {
      pValues[ i ] = arguments.length > 1 ? sliceDeferred.call( arguments, 0 ) : value;
      deferred.notifyWith( promise, pValues );
    };
  }
  if ( length > 1 ) {
    // 遍历所有的deferred参数，为每个参数添加一个成功的回调函数，和失败的回调函数，还有progress回调函数， 同时检测正确的deferred有结果，用于修正count
    // 就是让每一个deferred添加成功回调函数，然后这个函数在一直修正count，直到0则执行主deferred的回调函数。
    for ( ; i < length; i++ ) {
      if ( args[ i ] && args[ i ].promise && jQuery.isFunction( args[ i ].promise ) ) {
        args[ i ].promise().then( resolveFunc(i), deferred.reject, progressFunc(i) );
      } else {
        --count;
      }
    }
    // 如果一直修正count，最终是!0了，那就证明所有的参数都不是正确的deferred，那么立即执resolved了
    if ( !count ) {
      deferred.resolveWith( deferred, args );
    }
    // 如果length不大于1，同时主deferred也不是传入参数， 也证明没有正确的deferred，那么直接执行resolved好了
  } else if ( deferred !== firstParam ) {
    deferred.resolveWith( deferred, length ? [ firstParam ] : [] );
  }
  return promise;
}
```


## Summary
deferred部分在jQuery中起到了函数异步管理的作用，但内部并没有用到任何异步的内容，它的基础是jQuery.Callbacks对象。

jQuery.Callbacks对象：它在内部维护一个管理函数的list，同时有一个很重要的firing状态用来表示当前list内的函数是不是在执行中，对list的任何操作都需要考虑firing状态，用于保证list内的函数执行过程中即不略过也不被重复执行。

deferred也是一个很有意思的内容，deferred内部维护三个jQuery.Callbacks对象的list，对应done、fail、progress三个状态。并且deferred的三个状态是模拟promise的状态，逻辑处理也是按照promise的逻辑来处理，状态只能更改一次。

when是一个多deferred管理的方法，接受多个deferred作为参数，但实现一个主deferred对象同时when.count变量作为整体状态的记录，并且为每个deferred添加一个done和fail方法，done用于一直维护when.count， 然后fail用于终止并调用主deferred的fail。

这三个对象的设计可谓精妙。让使用者在主观上完成了异步、回调函数的轻松管理，再加上其链式写法让整个回调的管理更加简单。 但研究它的源码，也不过区区一两百行，但实现的东西却是极为丰富的。