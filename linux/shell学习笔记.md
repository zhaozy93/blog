# Shell 学习笔记

## 1. Sheel 是什么
Shell是一种解释性脚本语言。那么必然需要解释器
- Unix/Linux上常见的脚本解释器有bash、sh\csh等
- bash是Linux默认的标准shell解释器

Shell 有两种执行方式
- 交互式Interactive: 解释执行用户的命令，用户输入一条执行一条
- 批处理Batch: 用户事先写一个Shell脚本，Shell一次性把命令执行完，而不必一条条的敲入

## 2. Shell编写-基础知识
### 2.1 第一行内容 和 注释
``` shell
#! /bin/bash
# 第一行内容是约定内容，表示寻找指定的shell解释器
# 可以使用 /bin/sh test.sh 或 ./test.sh 方式去执行test.sh脚本
# test.sh 不可以执行，因为他回去PATH里面寻找目标文件，会找不到
```

### 2.2 变量
#### 2.2.1 定义变量
- 定义变量时变量名不建议以$开头
- 字母开头
- 不能出现空格、不能使用关键字
- 变量和=之间、=和value之间不能出现空格
#### 2.2.2 读取变量
- 读取变量需要在变量之前加$符号
- 或者${variableName} 使用{}包裹变量名
- 使用{}是一个好习惯
``` shell
#! /bin/bash
variable1="my name"
variable2=100
echo ${variable1}
echo ${variable2}
```
#### 2.2.3 重定义变量值
- 普通情况下不需要使用$
- 使用了$会有意想不到的结果
```shell
#! /bin/bash
variable1="variable2"
variable2=100
echo ${variable1}
echo ${variable2}

variable2=200
echo ${variable1}
echo ${variable2}

$variable1 = 300
echo ${variable1}
echo ${variable2}
```

#### 2.2.4 定义只读变量
```shell
#! /bin/bash
variable1=100
echo ${variable1}
readonly variable1
variable1=200
```
#### 2.2.5 删除变量
- unset 删除普通变量
- unset 不能删除只读变量
```shell
#! /bin/bash
variable1=200
echo ${variable1}
unset variable1
echo ${variable1}

variable1=200
echo ${variable1}
readonly variable1
unset variable1
echo ${variable1}
```
#### 2.2.6 特殊变量
|特殊变量|含义|
|:----:|:----:|
|$0|当前脚本的文件名|
|$n|传递给脚本或函数的参数,n表示1....的参数数位置|
|$#|传递给脚本或函数的参数个数|
|$*|传递给脚本或函数的全部参数|
|$@|传递给脚本或函数的所有参数.被双引号(" ")包含时,与 $* 稍有不同|
|$?|上个命令的退出状态,或函数的返回值|
|$$|当前Shell进程ID|

```shell
#! /bin/bash
echo $0
echo $#
echo $1 
echo $2
echo $*
echo $@
echo $$
echo $?
# ./test.sh golang shell
```
#### 2.2.7 $@与$*区别
```shell
#! /bin/bash
# 第一部分循环
echo "print each param from \$*"
for var in $*
do
echo "$var"
done
echo "print each param from \$@"
for var in $@
do
echo "$var"
done
# 第二部分循环
echo "print each param from \"\$*\""
for var in "$*"
do
echo "$var"
done
echo "print each param from \"\$@\""
for var in "$@"
do
echo "$var"
done 

# /bin/bash test.sh "golang shell" 两者没有任何区别
# /bin/bash test.sh golang shell 第一部分循环没区别，第二部分有区别  "$*" 会把内容看做一个整体
```

#### 2.2.8 变量替换
- 在双引号中使用变量是一种替换
- \n之类的转义符也属于
- -e是允许转义  -E则是禁止转义

```shell
#! /bin/bash
name="golang"
echo "My name is $name\n"
echo -e "My name is $name\n"
echo -E "My name is $name\n"
```
#### 2.2.9 变量替换命令
- 使用``把命令赋给某个变量
- 使用$variable来执行某个命令
```shell
#! /bin/bash
variable=`date`
echo $variable 
```
#### 变量检测
- ${var} 变量本身
- ${var:-word}  如果var不存在，返回word值
- ${var:=word}  如果var不存在，返回word值同时将var设为word值
- ${var:?message} 如果var不存在，那么输出message，同时结束脚本执行
- ${var:+word} 如果var存在，返回word值

```shell
#! /bin/bash
var1="golang"
var2="shell"
echo ${var1}
echo ${var3:-$var1}
echo ${var3}
echo ${var3:=$var1}
echo ${var3}
echo ${var4:?"sss"}
echo ${var1:+$var2}
```

### 2.3 运算符
#### 2.3.1 基础运算
|运算符|表示|
|:----:|:----:|
|+| \`expr $a + $b\` |
|-| \`expr $a - $b\` |
|*| \`expr $a \* $b\`|
|/| \`expr $a / $b\` |
|%| \`expr $a % $b\` |


#### 2.3.2 关系运算符
|运算符|表示|
|:----:|:----:|
|==| if [ $a == $b ]  |
|!=| if [ $a != $b ]  |
|备注| 下面的只支持数字了 |
|==| if [ $a -eq $b ] |
|!=| if [ $a -ne $b ] |
|< | if [ $a -lt $b ] |
|> | if [ $a -gt $b ] |
|<=| if [ $a -le $b ] |
|>=| if [ $a -ge $b ] |

#### 2.3.4 与或运算符
|运算符|表示|
|:----:|:----:|
| && | if [ part1 -a part2 ] | 
| \|\| | if [ part1 -o part2 ] | 

```shell 
#! /bin/bash
a=100
b=50

if [ $a == $b ]
then
echo "a is equal to b"
else
echo "a is not equal to b"
fi

if [ $a -gt 100 -a $a -gt 90 ]
then 
echo "a is large than 100 and 90"
else
echo "a is not large than 100 and 90"
fi 

if [ $b -gt 100 -o $b -gt 10 ]
then 
echo "b is large equal than 100 and 10"
else
echo "b is not large equal than 100 and 10"
fi 
```

#### 2.3.5 字符串检测运算符
|运算符|表示|
|:----:|:----:|
| 长度为0 | if [ -z $s ] | 
| 长度非0 | if [ -n $s ] |
| 长度非0 | if [ $s ]    |

#### 2.3.6 文件检测运算符
|运算符|表示|
|:----:|:----:|
| 是否块文件    | if [ -b $s ] | 
| 是否设备文件  | if [ -c $s ] |
| 是否目录     | if [ -d $s ]    |
| 是否普通文件  | if [ -f $s ]    |
| 是否可读文件  | if [ -r $s ]    |
| 是否可写文件  | if [ -w $s ]    |
| 是否可执行文件 | if [ -x $s ]    |
| 是否空文件    | if [ -s $s ]    |
| 是否存在      | if [ -e $s ]    |


### 2.4 字符串
#### 2.4.1单引号
- 单引号的字符串会原样输出
- 单引号的字符串内部不能出现单引号,使用转义符\也不可以
- 单引号内的变量是无效的

#### 2.4.2双引号
- 内部可以使用转义符
- 内部可以出现变量
- 内部可以嵌套双引号""

```shell
#! /bin/bash
name="golang"
echo "hello $name !" 
echo "hello ${name} !"
echo "hello "${name}" !"
echo "hello "$name" !" 
```

#### 2.4.3 字符串操作
```shell
#! /bin/bash
name="golang"
# 获取字符串长度
echo ${#name}
# 截取字符串  variable:start_position_length
echo ${name:0:2}

# 替换
echo ${name/"go"/"sql"} # 只替换一次
echo ${name//"go"/"sql"} # 替换所有
```

### 2.5 数组
#### 2.5.1 数组定义
```shell
#! /bin/bash
array_variable=(
  a
  b
  c
  d
)
echo ${array_variable}
array_variable[0]=1
array_variable[100]=1
echo ${array_variable}
echo ${array_variable[100]}
echo ${array_variable[99]}
echo ${array_variable[*]}
echo ${array_variable[@]}
```

#### 2.5.2 数组长度
```shell
#! /bin/bash
array_variable=(
  abc
  b
  c
  d
)
echo ${#array_variable[@]}
echo ${#array_variable[*]}
echo ${#array_variable[0]} #输出单个元素的长度
```

### 2.6 逻辑判断
#### 2.6.1 if...else...
```shell 
#! /bin/bash

if []
then
else 
fi

if []
then 
elif []
then
else
fi
```

#### 2.6.2 case esac 
```shell
#!/bin/bash
option="${1}"
case ${option} in
-f) FILE="${2}"
echo "File name is $FILE"
;;
-d) DIR="${2}"
echo "Dir name is $DIR"
;;
*)
echo "`basename ${0}`:usage: [-f file] | [-d directory]"
exit 1 # Command to come out of the program with status 1
;;
esac
```

#### 2.6.3 for
```shell
#! /bin/bash
for loop in 1 2 3 4 5
do
  echo $loop
done

for str in 'This is a string'
do
echo $str
done

for FILE in $HOME/.bash*
do
echo $FILE
done
```
#### 2.6.4 while
```shell
#! /bin/bash
cnt=0
while [ cnt -lt 5 ]
do 
cnt=`expr $cnt + 1 `
echo $cnt
done
```
#### 2.6.5 until
```shell
#! /bin/bash
cnt=0
until [ $cnt -gt 5 ]
do 
cnt=`expr $cnt + 1 `
echo $cnt
done
```
#### 2.6.6 continue break



### 2.7 函数function
```shell
# 如果没有return会把最后一条命令作为结果返回
# 函数返回值只能是整数 一般用0表示函数正常执行
function function_name(){
  return value 
}

# 函数也可以像遍历一样被删除
unset .f function_name
```

实战一个函数
```shell
#! /bin/bash
function f1(){
  echo "input one number \n"
  read num1
  echo "input another number \n"
  read num2
  return $(($num1 + $num2))
}

f1
echo "result is sum: $?"
```
带参数的函数调用
```shell
#!/bin/bash
funWithParam(){
echo "The value of the first parameter is $1 !"
echo "The value of the second parameter is $2 !"
echo "The value of the tenth parameter is ${10} !"
echo "The value of the eleventh parameter is ${11} !"
echo "The amount of the parameters is $# !" # 参数个数
echo "The string of the parameters is $* !" # 传递给函数的所有参数
}
funWithParam 1 2 3 4 5 6 7 8 9 34 73
echo $?
```

