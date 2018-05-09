## String
### string 向 byte转换

str := "我是中国人"

bytes := []byte(str)

转换后byte是一个数组 表示UTF8的编码

[230 136 145 230 152 175 228 184 173 229 155 189 228 186 186]

会发现转换后的长度变得很长， 从5个汉字变成了15个长度的数组，而且看不懂什么意思

转换步骤

* 汉字 --> unicode
    我 --> \u6211
* unicode --> 二进制
6211 --> 110001000010001
* 二进制--> ASCII码
110 001000 010001 --> 11100110 10001000 10010001
* ASCII码 --> 十进制
11100110 10001000 10010001 --> 230 136 145

这样转换完成后 我 就变成了 三个十进制长度的数组


注： 二进制转ASCII码
http://www.ruanyifeng.com/blog/2007/10/ascii_unicode_and_utf-8.html

### string的slice
str := "我是中国人"

len(str) // 15 直接调用的byte的长度

utf8.RuneCountInString(v_str) // 5

str_1 = []rune(str)

len(str_1) // 5 这个是可以直接使用切片的

str_1[1:2] // 26159 这样得到的是对应字的10进制，转换成16进制就是unicode

26159 --> 662f --> \u662f --> 是

### string 向 rune转换