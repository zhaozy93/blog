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

### 排序
先搜索 再排序

子句： SQL语句由子句构成，有些子句是必需的，而有的是可选的。一个子句通常由一个关键字和所提供的数据组成。子句的例子有SELECT语句的FROM子句，我们在前一章看到过这个子句。 子句的排列顺序会影响最后的结果。
- ORDER BY: SELECT columnsName FROM tableName ORDER BY columnsName;  // 对搜索结果排序， 可以根据多个列进行排序
- DESC、ASC: SELECT columnsName FROM tableName ORDER BY columnName1 DESC, columnName2 ASC; // 指定排序方向， 升降序

### 过滤数据
先搜索 再过滤 最后排序

- WHERE: SELECT columnsName FROM tableName WHERE column = value
- 操作符(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql01.jpeg)
- 空值检测 WHERE column IS NULL
- AND： SELECT columnsName FROM tableName WHERE column1 = value1 AND column2 = value2
- OR: SELECT columnsName FROM tableName WHERE column1 = value1 OR column2 = value2
- IN: SELECT columnsName FROM tableName WHERE column1 IN ( value1, valu2, value3)  // IN操作符用来指定条件范围，范围中的每个条件都可以进行匹配。 类似于OR
- NOT: ex：NOT IN， NOT BETWEEN // 有且只有一个功能，那就是否定它之后所跟的任何条件。
- LIKE： SELECT columnsName FROM tableName WHERE column LIKE %value%
- REGEXP: SELECT columnsName FROM tableName WHERE column REGEXP 'regexpattern'

### 通配符
- %： 任何字符出现 任意次数  column LIKE '%value%'   %不能匹配NULL
- _: 下划线只匹配单个字符而不是多个字符。

### 正则匹配
- |： 或
- [1234]： 匹配其中一个
- [1-5]： 匹配范围
- 匹配字符类(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql02.jpeg)
- 多次匹配(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql03.jpeg)
- 定位符
(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql04.jpeg)

### 计算字段
- 拼接Concat:  SELECT Concat(string1, column1, string2, column2) FROM tableName
- 删除空格Trim、RTrim、LTrim: SELECT Trim(Concat(string1, column1, string2, column2)) FROM tableName
- 使用别名 AS: SELECT Trim(Concat(string1, column1, string2, column2)) AS new_name FROM tableName
- 算数计算 + - * /: SELECT Trim(Concat(string1, column1, string2, column2)) AS new_name, column1 * column2 AS new_name2 FROM tableName

### 数据处理函数
每个DBMS对数据处理函数的实现都不尽相同，因此不要过于依赖数据处理函数。
(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql05.jpeg)
(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql06.jpeg)
(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql07.jpeg)

### 汇总数据
- AVG() 返回某列的平均值  SELECT AVG(colum) as column_average FROM tableName
- COUNT() 返回某列的行数  COUNT(*) 会计算表内一共有多少行包含NULL， COUNT(column)则会忽略NULL 
- MAX() 返回某列的最大值
- MIN() 返回某列的最小值
- SUM() 返回某列值之和

### 分组数据
先搜索 再过滤 再分组 最后排序
GROUP BY 有时会报一个错误

`ERROR 1055 (42000): Expression #2 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'mysql.user.User' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by`

错误原因也比较简单就是因为GROUP BY后面的字段与前面的字段不是完全一致的`sql_mode=only_full_group_by`，可以通过设置sql_mode来关闭这个选项

- GROUP BY 对检索后数据进行分组 SELECT columnsName FROM tableName GROUP BY columnsName;
  - GROUP BY子句可以包含任意数目的列  
  - 如果分组列中具有NULL值，则NULL将作为一个分组返回。如果列中有多行NULL值，它们将分为一组。
- HAVING: 支持所有WHERE的搜索 SELECT columns FROM table GROUP BY column HAVING column > num
  - 唯一的差别是 WHERE过滤行，而HAVING过滤分组

图中搜索的含义是: WHERE子句过滤所有`prod_price`至少为10的 行。然后按`vend_id`分组数据，最后HAVING子句过滤计数为2或2以上的分组(过滤分组)
(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql08.jpeg)

### SELECT子句使用顺序
(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql09.jpeg)
(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql10.jpeg)

### 子查询
当我们需要先从A表中查出某个数据作为B表的筛选条件时。

`SELECT product_name FROM product WHERE product_id = ( SELECT product_id FROM order WHERE order_id = 1000 )`

先在订单表中查出订单1000对应的产品ID，再去产品表中查出对应ID的产品名称

### 联结查询
关系型数据库最核心的设计之一就是各表之间互相通过外键依赖来完成数据存储，那么查询的时候肯定不希望每次都像子查询那样麻烦。

- 使用WHERE来表达联结条件 SELECT table1.column1, table2.column2, table2.column3 FROM table1, table2 WHERE table1.column1 = table2.column1
- 明确指出内部联结 SELECT table1.column1, table2.column2, table2.column3 FROM table1 INNER JOIN table2 ON table1.column1 = table2.column1
- 使用别名 自联结 SELECT t1.column1, t1.column2 FROM table as t1, table as t2 WHERE t1.c1 = t2.c1 AND t2.c2 = 100
(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql11.jpeg)
- 外部联结 LEFT/RIGHT OUTER JOIN 外部联结还包括没 有关联行的行： SELECT table1.column1, table2.column2, table2.column3 FROM table1 LEFT OUTER JOIN table2 ON table1.column1 = table2.column1 
- 内部联结与外部联结区别： 内部联结如果table1中某一行没有在table2中匹配到，那么table1这一行是不会出现在结果中，但是外部联结会保留这一行。  外部联结的LEFT、RIGHT就是用来表明那一张表的每一行都应该被保留

### 组合查询
UNION 简单讲就是讲多个WHERE联合到一起
`SELECT ven_id, ven_name, produ_price FROM products WHERE produ_price <= 10 OR ven_id IN (1001, 1002)`  ==>

`SELECT ven_id, ven_name, produ_price FROM products WHERE produ_price <= 10 UNION SELECT ven_id, ven_name, produ_price FROM products ven_id IN (1001, 1002)`
- UNION必须由两条或两条以上的SELECT语句组成，语句之间用关键字UNION分隔。
- UNION中的每个查询必须包含相同的列、表达式或聚集函数(不过各个列不需要以相同的次序列出)。
- 列数据类型必须兼容:类型不必完全相同，但必须是DBMS可以隐含地转换的类型(例如，不同的数值类型或不同的日期类型)。
- UNION会自动取消重复的行， 如要保留多次 需使用UNION ALL
- 在用UNION组合查询时，只能使用一条ORDER BY子句，它必须出现在最后一条SELECT语句之后。

## 插入数据
INSERT INTO table VALUES(value1, valu2, value3);
- INSERT语句一般不会产生输出
- 存储到每个表列中的数据在VALUES子句中给出，对每个列必须提供一个值 
- 如果某个列没有值，应该使用NULL值(假定表允许对该列指定空值)
- 各个列必须以它们在表定义中出现的 次序填充

INSERT INTO table(column1, column2, column3) VALUES(value1, value2, value3);
- VALUES列表中的相应值填入列表中的对应项， 即使表的结构改变，此INSERT语句仍然能正确工作
- 可以省略列。这表示可以只给某些列提供值，给其他列不提供值。
- 一条INSERT插入多条记录 INSERT INTO table(columns) VALUES(values), VALUES(values), VALUES(values)
- INSERT INTO table(columns) SELECT columns FROM table2 请确保columns与columns位置一致， mysql不关心列名，在这里更关心列的位置

## 更新和删除数据
UPDATE table SET column2 = value2, column3 = value3 WHERE column1 = value1
- UPDATE 由三部分组成 更新的表名、列名和新值、更新行的过滤条件(除非更新整个表)
- 如果一次性更新多行，当有一行出现错误时，所有行会被恢复为原样 可以使用IGNORE 来强制更新继续并且忽略错误的那一行

DELETE FROM table WHERE column1 = value1
- 删除行而不是删除某一列 删除某一列请使用UPDATE 列 = NULL
- 当没有WHERE时代表删除整张表所有行
- 删除整张表效率较低， 可使用TRUNCATE table 来代替，他是删除表然后重建一张一样的表，但是DELETE FROM tabl逐行删除

## 表操作
### 创建表
```sql
    CREATE TABLE table (
        column1 type NOT NULL AUTO_INCREMENT,
        column2 type NOT NULL,
        column3 type NULL  DEFAULT 1,
        column4 type NULL,
        PRIMARY KEY(column1)
    ) ENGINE = InnoDB
```
注意事项
- 每个表只允许一个AUTO_INCREMENT列，而且它必须被索引(如，通过使它成为主键)。
- `SELECT last_insert_id()` 获取最后一次插入的自增的值

引擎类型
- InnoDB是一个可靠的事务处理引擎，它不支持全文本搜索;
- MEMORY在功能等同于MyISAM，但由于数据存储在内存(不是磁盘)中，速度很快(特别适合于临时表)
- MyISAM是一个性能极高的引擎，它支持全文本搜索，但不支持事务处理。

### 更改表结构
- ALERT TABLE table ADD column type;  增加字段
- ALERT TABLE table DROP COLUMN column; 删除字段
- ALTER TABLE table ADD INDEX indexName (column); 增加索引，为column增加索引，索引名称为indexName
- ALTER TABLE table ADD CONSTRAINT fkName FOREIGN KEY (column1) REFERENCES table2(column2);  增加外键 在table表中未column设置外键，约束来自于table2的column2字段，外键的名称叫做fkName
-  RENAME TABLE table1 TO table2; 修改表名称， 将表名有table1 修改至 table2

## 视图
视图是虚拟的表。与包含数据的表不一样，视图只包含使用时动态 检索数据的查询。 

视图特别像提前写好一段sql语句，将这段语句的返回结果作为一个临时的表(视图)。 之后可以拿这个临时表(视图)当做普通表一样进行sql查询等。 

不建议在视图上进行更新操作， 同时视图数量不宜过多

- 重用SQL语句。
- 简化复杂的SQL操作。在编写查询后，可以方便地重用它而不必知道它的基本查询细节。
- 使用表的组成部分而不是整个表。
- 保护数据。可以给用户授予表的特定部分的访问权限而不是整个表的访问权限。
- 更改数据格式和表示。视图可返回与底层表的表示和格式不同的数据。

- CREATE VIEW viewName AS SELECT columns FROM tables WHERE .....;   创建视图
- SHOW CREATE VIEW viewname; 显示视图简介
- DROP VIEW viewname; 删除视图
- CREATE OR REPLACE VIEW viewName AS .......; 更新视图

## 事务
事务处理(transaction processing)可以用来维护数据库的完整性，它保证成批的MySQL操作要么完全执行，要么完全不执行。
- 事务(transaction)指一组SQL语句;
- 回退(rollback)指撤销指定SQL语句的过程;
- 提交(commit)指将未存储的SQL语句结果写入数据库表;
- 保留点(savepoint)指事务处理中设置的临时占位符(place-holder)，你可以对它发布回退(与回退整个事务处理不同)

隐含提交(implicit commit)，即提交(写或保存)操作是自动进行的。 但事务过程中不会进行隐含提交，需要我们明确的显式提交。 

- START TRANSACTION;   开始一个事务
- ROLLBACK; 回滚
- SAVEPOINT name; 创建一个保留点，并命名。
- ROLLBACK TO pointname; 回滚到某一个保留点
- COMMIT; 提交前面所做的操作     

## 用户管理
- CREATE USER username IDENTIFIED BY 'password'; 创建用户
- RENAME USER old TO new;  更改用户名
- DROP USER name;   删除用户
- SHOW GRANTS FOR name; 显式某个用户的权限
- GRANT SELECT ON db.table FOR name; 为用户name分配db数据库下table表的select权限
- REVOKE SELECT ON db.table FOR name;  REVOKE是GRANT的反操作
- SET PASSWORD FOR username = Password('new password'); 更改其他用户的密码
- SET PASSWORD = Password('new password');  更改当前用户的密码

权限列表
(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql12.jpeg)
(https://raw.githubusercontent.com/zhaozy93/blog/master/img-bed/sql13.jpeg)



 