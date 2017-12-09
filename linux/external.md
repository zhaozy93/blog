- $PATH 环境变量
```
  echo $PATH 可以查看当前有多少环境变量
```

1. mkdir -pm
```
   mkdir -m 777 dirname  新建制定权限的文件夹
   mkdir -p dirname/dirname/dirname  递归建立所需要的文件夹
```

2. rmdir -p
```
   rmkdir -p dirname/dirname/dirname 递归的删除上层的空文件夹
```

3. ls -aAdfFhilnrRSt
```
   ls -a 列出所有
   ls -l 列出详细信息
   ls -lh  列出可读模式的容量大小
   ls --full-time 列出详细时间
   ls path/dirname 列出指定目录下的文件
```

4. cp(copy) -adfilprsu
```
   cp -d 若source为链接文件时，复制链接文件而非文件本身
   cp -f 若target已经存在，则移除后再复制
   cp -i 若target已存在，覆盖前询问
   cp -p 连同source的文件属性一起复制
   cp -r 递归的复制
   cp -u 若target比source旧才复制
```

5. rm(remove) -fir
```
   rm -i 询问互动式删除
   rm -r 递归删除
   rm -f 强制删除,忽略警告, 忽略不存在文件
```

6. mv(move) -fiu 改名的好方法
```
   mv -i 询问互动式移动
   mv -f 强制移动,忽略警告,即使target存在也直接覆盖
   mv -u 当source比较新时才覆盖
```

7. cat(concatenate) -AbEnTv
```
   cat -b 列出行号(空白行不计数)
   cat -E 列出结尾的断行字符$
   cat -n 列出行号
   cat -T 列出Tab键
   cat -v 列出一些看不出的特殊字符
```

8. tac 与cat一致 反向输出

9. nl 带行号输出

10. more 翻页形式读取文件
``` 
   空格 向下翻页
   回车 向下一行
   / 向下搜索
   :f 显示出文件名以及当前显示的行数
   q 退出
   b 向前翻页
```

11. less 一页一页翻动  man命令调用的less
```
   空格 向下翻页
   回车 向下一行
   / 向下搜索
   ? 向上搜索
   n 向后缩搜一个
   N 向前搜索一个
   q 退出
```

12. head 读取文件前几行 -n
```
   head -n 100 filename 读取文件前100行
```


13. tail 读取文件后几行 -nf
```
   tail -n 100 filename 读取文件后100行
   tail -f filename 持续侦测文件的改动，并显示在界面上
```

14. od -t [TYPE] 读取非纯文本文件 
```
   od -t a 利用默认字符输出
   od -t c 利用ASCII字符输出
   od -t d[size] 利用十进制(decimal)来输出数据，每个整数占用size bytes
   od -t f[size] 利用浮点数(floating)来输出数据，每个整数占用size bytes
   od -t o[size] 利用八进制(octal)来输出数据，每个整数占用size bytes
   od -t x[size] 利用十六进制(hexadecimal)来输出数据，每个整数占用size bytes
```

15. touch 修改文件时间或者新建文件
- modification time(mtime)： 当文件数据内容更改时，就会更新这个时间
- status time(ctime)： 当文件状态更改时，会更新这个时间。权限与属性更改时会更新
- access time(atime)：文件内容被读取时会更新这个时间

16. umask 预设权限
- 仅看后三位即可
- umask数值为拿掉的数值， 因此022表示拥有者全部，group和others拿掉写权限
```
   umask 看不懂
   umask -S 符号类型显示
   umask 000 直接跟数值表示更改umask
```

17. chattr、lsattr
- [chattr](http://www.ha97.com/5172.html)
- [csdn](http://blog.csdn.net/sailor201211/article/details/53215060)
- [man](https://linux.die.net/man/1/chattr)

18. 文件的特殊权限 SUID, SGID, SBIT
![ll](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/linux100.jpeg)
可以发现passwd文件除了rwx竟然还有一个s权限

- Set UID: 当s出现在owner的x权限上时，-rwsr-xr-x，称为Set UID。
  - SUID权限仅对二进制程序有效
  - 执行者对于该程序需要具有x的可执行权限
  - 本权限仅在执行该程序的过程中有效
  - 执行者将对该程序拥有owner权限
- Set GID: 当s出现在group的x权限上时，-rwx--s--x，称为Set GID。
  - SGID权限可对文件或者目录
  - 程序: 执行者对于该程序需要具有x的可执行权限
  - 程序: 执行者将对该程序拥有group权限
  - 目录: 用户对此目录具有r与x权限时，可以进入此目录
  - 目录: 用户在此目录下的群组将变为该目录的群组
- Sticky Bit: 当t出现在others的x权限上时，-rwxr-xr-t，称为Sticky Bit。
  - 仅对目录有效
  - 用户具有此目录的w、x权限时，亦具有写入权限时
  - 用户在此目录下新建的目录或文件，仅用户自己与root才有权删除

设定特殊权限
```
   // 以前是使用3位数字来表示rwx
   chmod 777 filename 
   // 特殊权限位于旧3位权限格式的前方，构成4位权限格式
   // 4 ==> SUID
   // 2 ==> SGID
   // 1 ==> SBIT
   chmod 7777 filename 
```

SUID
![ll](https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/linux101.jpeg)

19. file 查看文件类型
```
   file filename 查看文件属于ASCII或者data文件，或者是binary文件
```

20. which -a 寻找命令的位置 
``` 
   which -a node 例如所有指令位置
   which node 寻找node所在的位置
```

21. whereis -bmsu 文件名搜索 
```
   whereis -b 只找binary文件
   whereis -m 只找说明文档
   whereis -s 只找source文件
   whereis -u 搜索不属于上述三类的文件
```

22. locate -ir 定位文件
```
   locate keyword 寻找文件
   locate -i 忽略大小写
   locate -r 后面keyword为正则写法
```

23. find [PATH] [option] [action] 寻找文件
```
    find / -mtime n n天之前(就是那一天24小时而已)被修改的
    find / -mtime +n n天之前被修改的
    find / -mtime -n n天之内被修改的
    find / -user root 属于root用户的
    find / -group root 属于root用户组的
    find / -nouser 不归属于任何用户的文件
    find / -nogroup 不归属于任何用户组的文件
    find / -name filename 按名字搜索
    find / -size [+-]n 按大小搜索
    find / -perm [+-]modeNUM 按权限搜索
```
[find](http://www.cnblogs.com/peida/archive/2012/11/13/2767374.html)
