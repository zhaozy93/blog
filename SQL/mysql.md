# 基础

## 概念
数据库(database)： 保存有组织的数据的容器(通常是一个文 件或一组文件)。

DBMS(数据库管理系统)： 数据库是通过DBMS创建和操纵的容器。

表(table)：某种特定类型数据的结构化清单。

模式(schema)： 关于数据库和表的布局及特性的信息。

列(column): 表中的一个字段。所有表都是由一个或多个列组成的。

数据类型(datatype): 所容许的数据的类型。每个表列都有相应的数据类型，它限制(或容许)该列中存储的数据。

行(row): 表中的一个记录。

主键(primary key): 一列(或一组列)，其值能够唯一区分表中每个行。
- 任意两行都不具有相同的主键值;
- 每个行都必须具有一个主键值(主键列不允许NULL值)。
好习惯
- 不更新主键列中的值;
- 不重用主键列的值;
- 不在主键列中使用可能会更改的值。(例如，如果使用一个名字作为主键以标识某个供应商，当该供应商合并和更改其 名字时，必须更改这个主键。)


操作符(operator)：用来联结或改变WHERE子句中的子句的关键 字。也称为逻辑操作符(logical operator)。

通配符(wildcard)： 用来匹配值的一部分的特殊字符。 

搜索模式(search pattern)： 由字面值、通配符或两者组合构成的搜索条件。


SQL： 是结构化查询语言(Structured Query Language)的缩写。SQL是一种专门用来与数据库通信的语言。

DBMS可分为两类:一类为基于共享文件系统的DBMS，另一类为基于客户机—服务器的DBMS。

客户机—服务器应用分为两个不同的部分。
- 服务器部分是负责所有数据访问和处理的一个软件。这个软件运行在称为数据库服务器的计算机上。与数据文件打交道的只有服务器软件。
- 客户机是与用户打交道的软件。向服务器部分发送数据添加、删除、更新等请求。

## 基本命令
sql命令不区分大小写，不过将指令大写，表名、列名小写是一种好的习惯
- SHOW DATABSES;   // 显示当前数数据库列表
- USE tableName;   // 使用(选择)某个数据库
- SHOW TABLES;     // 显示数据库中表列表
- SHOW COLUMNS FROM tableName;   或 DESCRIBE tableName;   // 显示表中列的详细信息
- SHOW STATUS，用于显示广泛的服务器状态信息;
- SHOW ERRORS和SHOW WARNINGS，用来显示服务器错误或警告消息。

## 检索数据
- SELECT columnName FROM tableName;  // 从某张表中搜索1个列
- SELECT columnName1, columnName2, columnName3, columnName4 FROM tableName; // 从某张表中搜索多个列
- *: SELECT * FROM tableName;  // 从某张表中搜索所有列
- DISTINCT: SELECT DISTINCT columnName FROM tableName;   // 只检索出不同的列，  即 如果有两条记录 值相同，则只返回一条
- Limit: SELECT columnName FROM tableName Limit n;  // 检索结果只返回限制的N条记录
- Limit: SELECT columnName FROM tableName Limit n1, n2;  // 检索结果只返回限制的从n1条的n2条记录， 记录是从0开始计数的
- Limit OFFSET: SELECT columnName FROM tableName Limit n1 OFFSET n2;  // 检索结果只返回限制的从n2条的n1条记录

## 排序
先搜索 再排序

子句： SQL语句由子句构成，有些子句是必需的，而有的是可选的。一个子句通常由一个关键字和所提供的数据组成。子句的例子有SELECT语句的FROM子句，我们在前一章看到过这个子句。 子句的排列顺序会影响最后的结果。
- ORDER BY: SELECT columnsName FROM tableName ORDER BY columnsName;  // 对搜索结果排序， 可以根据多个列进行排序
- DESC、ASC: SELECT columnsName FROM tableName ORDER BY columnName1 DESC, columnName2 ASC; // 指定排序方向， 升降序

## 过滤数据
先搜索 再过滤 最后排序

- WHERE: SELECT columnsName FROM tableName WHERE column = value
- 操作符(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql01.jpeg)
- 空值检测 WHERE column IS NULL
- AND： SELECT columnsName FROM tableName WHERE column1 = value1 AND column2 = value2
- OR: SELECT columnsName FROM tableName WHERE column1 = value1 OR column2 = value2
- IN: SELECT columnsName FROM tableName WHERE column1 IN ( value1, valu2, value3)  // IN操作符用来指定条件范围，范围中的每个条件都可以进行匹配。 类似于OR
- NOT: ex：NOT IN， NOT BETWEEN // 有且只有一个功能，那就是否定它之后所跟的任何条件。
- LIKE： SELECT columnsName FROM tableName WHERE column LIKE %value%

## 通配符
- %： 任何字符出现 任意次数  column LIKE '%value%'   %不能匹配NULL
- _: 下划线只匹配单个字符而不是多个字符。