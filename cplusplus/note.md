# C++ Primer 5
## Chapter 1
### 初识输入输出
```c++
#include <iostream>
int main(){
	std::cout << "Enter two number" << std::endl;
	int v1 = 0, v2 = 0;
	std::cin >> v1 >> v2;
	std::cout << "Numbers are " << v1 << " and " << v2 << std::endl;
	return 0
}
```
`<<`输出运算符，左侧为一个ostream对象，右侧为要打印的值。
`>>`输入运算符，左侧为一个istream对象，右侧为接收对象。
`<< 与 >>`的返回结果都是左侧iostream对象，当 >> 遇到 EOF或者无效输入时会变成无效。
写入cout的数据会进入缓冲区，调用std::endl会刷入设备。 cerr的数据不进缓冲区，直接进设备。调用cin会刷新cout。

## Chapter 2
### 基本内置类型
c++基本内置类型包含 算数类型 和 空类型。 算数类型又可分为整型和浮点型。
算数类型包含: bool、char、wchar_t、char16_t、char32_t、short、int、long、long long、float、double、long double。
#内置类型的及其实现#
- 机器内部总是以比特序列存储数据，即0、1。 一个字节包含8比特(位)。
- 存储的基本单元称为 位。 常见一位等于4个或8个字节，即32或64比特(位)。
- 内存中每一个字节与一个数字(地址)关联
- 每个内存地址必须赋予一个数据类型才有意义。 
#### 类型转换
- 有符号值和无符号值在一起时，有符号值会转为无符号值再计算
- 非布尔值赋给布尔类型时，0为false
- 布尔值赋给非布尔类型时，false为0，其余为1
- 浮点数赋值给整数，近似处理，仅保留整数部分
- 整数赋值给浮点数，小数部分为0
#### 字符串字面值与字符
字符串字面值的类型实际上由常量字符构成的数组，编译器会在每个字符串的结尾处添加一个空字符 ‘\0’，因此，字符串字面值实际长度比他的内容多1.
#### 指定字面值的类型
通过添加前缀与后缀可以指定字面量的类型
前/后缀	含义				类型
u 前 	Unicode16字符		char16_t
U 前		Unicode32字符		char32_t
L 前		宽字符				wchar_t
u8 前	UTF-8				char
u或U 后         unsigned  最小匹配类型			
l 或L 后		long	  最小匹配类型	
Ll 或 LL后	long long 最小匹配类型	
f 或 F		float		   最小匹配类型	
L 或 L 		long float  最小匹配类型	
### 变量
变量提供一个具名的、可供程序操作的存储空间。
在一条定义语句中，可以用先定义的变量值去初始化后定义的其他变量
#### 列表初始化
在c++11中才可以使用 新特性
当用列表初始化且存在精度丢失风险时，编译器会报错
```c++
int sold = 0;
int sold = {0};
int sold{0};
int sold(0);
double a = 3.1415926;
int sold{a};  // 报错，因为存在精度丢失风险
int sold = a  // 不报错
int sold(a) 	// 不报错
```
#### 默认初始化
定义变量时没有指定初始值，则变量被默认初始化，默认值由变量类型决定，同时定义位置也会有影响。
- 内置类型的默认初始化由位置决定，任何函数体外的变量初始化为0。 定义在函数体内部的内置类型变量不被初始化，此时变量的值是未定义的，访问或拷贝此值将引发错误
- 类各自决定初始化对象的方式，是否允许不经初始化就定义对象也由类决定
#### 变量声明和定义的关系
分离式编译机制: 允许将程序分割为若干个文件，每个文件可被独立编译
- 声明: 使得名字为程序所知，一个文件如果想使用别处定义的名字则必须包含对那个名字的声明
- 定义: 负责创建与名字关联的实体，还申请存储空间、也可能为变量赋一个初始值
``` c++
// extern 表示声明，不能有初始值(const除外)
// 函数体内部初始化一个由extern标记的变量将报错
extern int i;
int j;
extern int i = 10; // 定义了 
```
### 复合类型
#### 引用(左值引用)
为对象起另外一个名字，引用类型引用另外一种类型。
定义引用时，程序把引用和他的初始值绑定bind在一起而不是拷贝，因此必须有初始值！
- 引用的类型必须与其所引用对象的类型一致（有两个例外）
- 引用本身不是一个对象，因此不存在引用的引用，也不存在指向引用的指针
- 引用比如绑定在对象上，不能是字面量(除特殊情况外)
```c++
int ival = 1024;
int &refVal = ival;
int &refVal2 ; //报错 未初始化
```
#### 指针
指针是指向另外一种类型的复合类型。 指针也实现了对其他对象的间接访问，但指针本身也是一个对象，允许对指针赋值和拷贝。 指针无需在定义时赋值。
指针类型必须与其所指向对象的类型一致（有两个例外）。
```c++
int *p1, *p2;   // 两个指向int类型的指针
double dp, *dp2;  // dp是double类型对象，dp2是指向double类型的指针
int val = 42;
int *p = &val;   // 获取对象地址
```
#### 指针值
- 指向一个对象
- 指向紧邻对象所占空间的下一个位置
- 空指针，意味着指针没有指向任何对象
- 无效指针，上述情况外的其他值
#### 空指针
- int *p1 = nullptr;
- int *p2 = 0;
- int *p3 = NULL;   (需引入cstdlib，等价于 *p3 = 0 )
#### void*指针
void*指针是一种特殊的指针类型，可以存放任意对象的地址。 但我们并不知道里面存放的是什么类型的对象，因此不能对其直接进行操作，仅仅是一个内存空间而已。
```c++
double obj = 3.14, *pd = &obj;
void *pv = &obj;  // 正确，可以存放任意类型的对象
pv = pd; // 可以存放任意类型的指针
```
### const限定符
- 默认状态下，const对象仅在文件内有效
- extern const int bufSize = 512  // 可以多个文件共享同一个const变量
#### 常量的引用
- 常量的引用不能被用于修改它所绑定的对象
- 初始化常量引用时允许用任意表达式作为初始值，只要该表达式的结果能转换成引用的类型即可
```c++
const int ci = 1024;
const int &r = ci;
```
```c++
// const引用 不同类型变量
// 会发现r1的值没有更改.
// 原因是因为aa和r1不是同类，会产生一个临时变量 int temp = aa
	double aa = 42.00;
	const int &r1 = aa;
	std::cout << "R1 is " << r1;
	aa = 20;
	std::cout << "R1 is " << r1;
```
```c++
// const引用同类变量， 尽管不能通过引用去更改变量值，但是可以直接去修改 
	int aa = 42.00;
	const int &r1 = aa;
	std::cout << "R1 is " << r1;
	aa = 20;
	std::cout << "R1 is " << r1;
```
#### 指向常量的指针
- 指向常量的指针不能用于改变其所指对象的值
- 允许一个指向常量的指针指向一个非常量对象
```c++
const double pi = 3.14;
const double *ptr = &pi;
```
#### const指针
指针是对象，因此存在const指针。 
- const指针是不允许更改指针所存放的那个地址，但是地址对应的值可以更改
- 常量的指针不允许更改对应地址对应的常量的值，但是可以更换地址。
```c++
const double pi = 3.14;
double *const ptr = &pi;
```
#关于const指针的解释#
[c++ - What is the difference between const int*, const int * const, and int const *? - Stack Overflow](https://stackoverflow.com/questions/1143262/what-is-the-difference-between-const-int-const-int-const-and-int-const)
#### constexpr和常量表达式 c++11
常量表达式是指值不会改变并且在编译过程就能得到计算结果的表达式。
用const定义的不一定是常量表达式如`const int sz = get_size()`，虽然不报错，但是sz需要依赖执行后才能知道结果，因此不是常量表达式。 
因此c++11增加了 constexpr类型`constexpr int sz = get_size()`，只有get_size也是constexpr类型的函数时，才能编译通过。
### 处理类型
#### 类型别名
类型别名是一个名字，它是某种类型的同义词。有typedef和using两种方式
尽量不要对符合类型与指针使用类型别名，很难理解！！！
```c++
typedef doublle wages;  // wages是double的同义词
typedef wages base, *p;// base是double的同义词，p是double*的同义词
using SI = Sales_item; // SI是Sales_item的同义词
```
### auto类型
编译器会根据计算结果自动推算变量的类型，但在定义时使用auto要求同一行定义中类型一致。 auto会自动丢掉顶层const
```
auto item = val1 + val2; 
auto i = 0, *p = &i;  // 正确 int类型
auto i = 0, a = 2.14;   // 错误，类型不一致
```
#### decltype类型
希望用表达式的类型推断出要定义的变量的类型，但是不想用该表达式的值初始化变量。 会保留顶层const
- 表达式的内容是解引用操作，那么decltype将得到引用类型
- decltype((variable)) 得到的永远是variable的引用类型
- decltype( i=x ) 得到的是x类型的引用
```c++
decltype(f()) sum = x; // 不会真的调用f,只是拿到f的返回值类型
const int ci = 0, &cj = ci;
decltype(ci) x = 0;  // const int
decltype(cj) y = x;  // const int引用
decltype(cj) z;   // 错误，引用必须初始化
```
### 自定义数据结构
```c++
// example
struct Sale_data {
	std::string bookNo;
	unsigned units_sold = 0;
	double revenue = 0.0;
};
```
#### 预处理器
预处理器是在编译之前执行的一段程序，可以部分的更改我们所写的程序。
- include 是一种预处理标记
- 头文件保护符也是一种预处理标记
##### 头文件保护符
使用头文件保护符来确保唯一性
- define 把一个名字设定为预处理变量
- ifdef 当变量已定义时为真
- ifndef 当变量未定义时为真
- endif 结束指令
```c++
#ifndef SALES_DATA_H
#define SALES_DATA_H
#include<string>
struct Sales_data {
	std::string bookNo;
};
#endif
```
## Chapter3
### 命名空间的using声明
之前每次使用cin、cout都必须显式指出属于命名空间std。通过using声明可以简单的直接引入。 但using不应被使用在头文件中，如果使用在头文件中，那每个使用了该头文件的文件就都有这个声明，容易出错。
```c++
using std::cin;
using std::cout;
using std::endl;
```
### 标准库类型string
```c++
#include <string>
using std::string;
string s1;    // 得到空字符串，默认初始化
string s2(s1);    // s2是s1的副本
string s2 = s1; // 等价于s2(s1)
string s3("value");    //字面值value的副本，除了字面值最后那个空字符
string s3 = "value";   // 等价于s3("value")
string s4(n, 'c');    // 把s4初始化为连续的n个字符c组成的串
```
#### 直接初始化与拷贝初始化
- 如果使用 = 初始化一个变量，实际执行的是拷贝初始化 copy initialization，编译器把右侧的初始值拷贝到新创建的对象中去。 
- 如果不使用等号，执行的是直接初始化 direct initialization
#### string上的操作
##### 读写string对象
- 与内置类型的输入输出操作一致，string对象的此类操作也是返回运算符左侧的运算对象作为结果，因此多个输入输出是可以写在一起的。
- 普通的cin >> 操作遇到空白符会打断，如果想获取一整行数据则使用getline(cin, s)
```c++
	string s;
	cin >> s;
	cout << s << endl;
```
```c++
// 对比普通的cin与getline
 	while(cin >> s){
	 	cout << s << endl;
 	}

	while(getline(cin, s)){
		cout << s << endl;
	}
```
##### empty和size
调用size返回的是string::size_type类型， 它是一种无符号类型的整数。因此不能将它与负数进行比较，也没意义。 
auto len = s.size() 来获取一个string::size_type类型的变量
```c++
// 只输出非空行
	while(getline(cin, s)){
		if (s.empty()){
			continue;
		}
		cout << s << endl;
	}
// 只输出大于80个字符的行
	while(getline(cin, s)){
		if (s.size()<80){
			continue;
		}
		cout << s << endl;
	}
```
##### 字符串比较
== 与 !==检验两个string对象是否相等，要求长度相同且包含的字符也全部相同，大小写敏感。
<、<=、>、>=按照字典序进行比较
- 如果两个string长度不同，且较短的string对象每个字符都与较长的string对象对应位置字符相同，则较短的string小于较长的string
- 如果两个string对象在某个对应位置不一致，则strig对象比较结果其实是第一个不相同字符的比较结果
##### 循环string对象 c++11
```c++
// 这里其实循环出的每个变量类型是char
string str("a string");
for (auto c : str){
	cout << c << endl;
}
```
```c++
// 通过引用可以修改string对象的内容
string str("a string");
for (auto &c : str){
	&c = toupper(c)
}
```
##### 下标形式访问
进行下表访问时不要越界， 不出现负数，尽可能使用 string::size_type类型
```c++
string s("a string");
s[0]
s[s.size()-1]

for(decltype(s.size()) index = 0; index < s.size()-1; ++index){
	s[index] = toupper(s[index])
}
```
### 标准库类型vector
vector表示对象的集合，其中所有对象的类型都相同。集合中的每个对象有一个索引。 vector也称为容器
```c++
#include <vector>
using std::vector;
vector <T> v1;
vector <T> v2(v1);
vector <T> v2 = v1;
vector <T> v3(n, val);
vector <T> v4(n);
vector <T> v5{a,b,c,d,e.....};
vector <T> v5={a,b,c,d,e.....}
```
#### 定义和初始化vector对象
```c++
vector <vector<int>> v1; // c++11没有问题
vector <vector<int> > v2; // 旧版本需要有一个空格
```
#### 向vector添加元素
```c++
vector <int> v1;
v2.push_back(100);
```
#### 其他vector操作
```c++
v.empty();
v.size();
v.push_back(i);
<  <=  >   >=   ==  !=
v[2];//下标访问不可以越界同时不可以使用下标添加元素，只能使用下标修改元素
```
### 迭代器介绍
为循环遍历容器二产生的机制，迭代器。
#### 使用迭代器
有迭代器的类型同时拥有返回迭代器的成员。 
```c++
// b表示v的第一个元素，e表示尾元素的下一个位置
// 空集合则b、e为同一个位置
auto b = v.begin(), e = v.end();
```
#### 迭代器操作
```c++
*iter  //返回迭代器所指元素的引用
iter->mem  // 解引iter并获取mem成员 (*iter).mem
++iter   // 令iter指向容器下一个元素
--iter   // 令iter指向容器上一个元素
iter1 == iter2  // 判断两个迭代器是否相同
iter1 != iter2  // 
iter + n // 向后移动N个位置，或者指向 尾元素的下一个位置
iter -n  
iter1 - iter2 // 两个迭代器相减就是他们之间的距离差
```
### 数组
类似于vector但是在定义时必须指定长度，且在编译阶段就必须指定。
数组不允许直接进行拷贝赋值
函数内部定义了某种内置类型的数组，那么默认初始化数组含有未定义的值
```c++
int a[] = {0, 1, 2};
int a1[3] = {0, 1, 2};
int a2 = a;  // 错误
a2 = a;   // 错误
```
#### 数组与指针
很多情况使用数组变量其实是指向数组第一个元素的指针
```c++
int a[] = {0, 1, 2, 3};
auto p1(a);   // 得到的是  int *p1 = &a[0]
int *p2 = a;    // 得到的是  int *p1 = &a[0]
```
##### 数组的指针也是迭代器
对应的首指针就是0位元素，尾指针就是最后一位元素的下一个索引位置。
```c++
int a[] = {0, 1, 2, 3};
int *p = a;  // 指向第一个元素 a[0]
++p;        // 指向第一个元素 a[1]
int *e = &a[4]; // 4不存在，尾指针，不要去赋值 不要去读取
```
向上面这样手动去取首尾指针极易出错。 此外还有迭代器类型的begin()、end()方式。与string和vector一致。
```c++
int *s = a.begin();
int *e = a.end();
```
