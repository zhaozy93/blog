在现在的开发中，微服务大行其道，涉及到大量的rpc调用，通信双方之间的耗时或者通信调用的时间成本就需要考虑，在很多服务中也是至关重要的。那么网络通信环境那么复杂，需要dns寻址、tcp握手、网络传输，各个环节其实我们很难具体的去评判，有时候还会涉及到重定向，那么将更加复杂。 虽然我们在开发过程中很难去衡量或者测量这些过程具体的时间花费，例如在golang中，tcp三次握手都是封装在语言层面的，根本无法准确衡量。 那么当出现网络状况不良，调用服务延迟过高，通信双方互相扯皮的时候应该如何去诊断呢，除了我们自身做好服务耗时监控，对网络状况有清晰的了解也是必不可少的。

现在就使用最常见的curl命令来诊断双方的通信环境。

## 阅读curl manual
在Mac系统下可以使用`curl --manual` 来查看用户手册,如果阅读manual会发现里面介绍了如何使用curl以及最简单的示例。例如使用curl获取ftp、ftps、ssl、with-password等各种类型的网络信息。

## 学习curl -w
在manual中有这样一个参数，那就是-w，对它的介绍如下。
```
CUSTOM OUTPUT

  To better allow script programmers to get to know about the progress of
  curl, the -w/--write-out option was introduced. Using this, you can specify
  what information from the previous transfer you want to extract.

  To display the amount of bytes downloaded together with some text and an
  ending newline:

        curl -w 'We downloaded %{size_download} bytes\n' www.download.com
```
如果我们再多了解一些，其实`-w`不仅可以输出目标请求的size，还可以输出更多关于网络状况的信息。 
那么-w是如何使用的呢，其实就是在调用curl的时候为-w熟悉配置一串字符串，字符串中填入某些提前约定好的变量即可。例如manual中的示例`curl -w 'We downloaded %{size_download} bytes\n' www.download.com` 通过添加size_download变量，我们可以获取此次curl下载内容的大小。

## 诊断网络环境
那么针对网络环境，有哪些提前约定好的变量呢。 常用的有以下这些。
* time_namelookup：DNS 域名解析的时候，就是把https://zhihu.com转换成 ip 地址的过程
* time_connect：TCP 连接建立的时间，就是三次握手的时间
* time_appconnect：SSL/SSH 等上层协议建立连接的时间，比如 connect/handshake 的时间
* time_redirect：从开始到最后一个请求事务的时间
* time_pretransfer：从请求开始到响应开始传输的时间
* time_starttransfer：从请求开始到第一个字节将要传输的时间
* time_total：这次请求花费的全部时间

直接执行 `curl -w '%{time_namelookup}' 'www.baidu.com'`和`curl -w '%{time_namelookup}' 'www.github.com'`
可以发现baidu的域名解析速度远快于github。 百度在`0.012409`秒，而github在`0.514144`秒。 相差甚远。

根据刚才的manual，也可以将配置字符串用文件形式代替，编写一个配置查询文件
```
time_namelookup:  %{time_namelookup}
time_connect:  %{time_connect}
time_appconnect:  %{time_appconnect}
time_redirect:  %{time_redirect}
time_pretransfer:  %{time_pretransfer}
time_starttransfer:  %{time_starttransfer}
time_total:  %{time_total}
```
再执行`curl -w '@config.txt' 'www.baidu.com'`查看得到的结果：
```
time_namelookup:  0.012721
time_connect:  0.017005
time_appconnect:  0.000000
time_redirect:  0.000000
time_pretransfer:  0.017141
time_starttransfer:  0.024385
time_total:  0.026298
```
根据这个结果就可以有针对性的优化网络环境。