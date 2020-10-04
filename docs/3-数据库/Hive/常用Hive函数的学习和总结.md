# 常用Hive函数的学习和总结

今天来小结一下工作中经常会使用到的一些Hive函数。关于Hive函数的总结，网上早有十分全面的版本。

参考：https://blog.csdn.net/doveyoung8/article/details/80014442。

本文主要从最常用和实用的角度出发，说明几个使用频率较高的函数，更注重使用函数组合来解决实际问题而不局限于单个函数的使用。所有数据都是虚构，代码均在本地的Hive环境上都通过测试。本文代码较多，需要各位看官耐心学习，可以收藏备查，欢迎补充和讨论。由于公众号对代码的支持不太友好，您可以在我的公众号(超哥的杂货铺) 后台回复“hive函数”获取本文的PDF版本，方便阅读。

### 1.json字符串处理：get_json_object，lateral view，explode，substr，json_tuple

先简要说明下几个函数的用法：

```
语法: get_json_object(string json_string, string path)
返回值: string
说明：解析json的字符串json_string，返回path指定的内容。如果输入的json字符串无效，那么返回NULL。
语法: explode(ARRAY)，经常和lateral view一起使用
返回值: 多行
说明: 将数组中的元素拆分成多行显示
语法: substr(string A, int start, int len),substring(string A, int start, int len)
返回值: string
说明：返回字符串 A 从 start 位置开始，长度为 len 的字符串
语法: json_tuple(string json_string, col1, col2, …) ，经常和lateral view一起使用
返回值: string
说明：同时解析多个json字符串中的多个字段
```

然后我们看实例：

```
--我们虚构的数据，jsondata.txt
1    {"store":{"fruit":[{"weight":8,"type":"apple"}, {"weight":9,"type":"pear"}], "bicycle":{"price":19.95,"color":"red"}}, "email":"amy@only_for_json_udf_test.net", "owner":"amy"}
2    {"store":{"fruit":[{"weight":8.1,"type":"apple"}, {"weight":9.2,"type":"pear"}], "bicycle":{"price":20.01,"color":"blue"}}, "email":"abc@example.com", "owner":"bob"}

hive> create table json_data(id int, data string) row format delimited fields terminated by '\t';
hive> load data local inpath 'jsondata.txt' into table json_data;
```

- 查询单层值

```
hive> select id, get_json_object(data, '$.owner') from json_data;
1    amy
2    bob
```

- 查询多层值1

```
#注意bicycle子串的格式同样是json格式
hive> select id, get_json_object(data, '$.store.bicycle.price') from json_data;
1    19.95
2    20.01
```

- 查询多层值2

```
#注意fruit子串的的格式是数组(带有方括号)，不是标准的json格式，下面语句取出fruit的值
hive> select id, get_json_object(data, '$.store.fruit') from json_data;
1    [{"weight":8,"type":"apple"},{"weight":9,"type":"pear"}]
2    [{"weight":8.1,"type":"apple"},{"weight":9.2,"type":"pear"}]


#可以使用索引访问数据里的值，如得到第一个fruit的数据，结果是一个json格式的
hive> select id, get_json_object(data, '$.store.fruit[0]') from json_data;
1    {"weight":8,"type":"apple"}
2    {"weight":8.1,"type":"apple"}

#在上面语句的基础上，可以获得weight和type值。
hive> select id, get_json_object(data, '$.store.fruit[0].weight') from json_data;
1    8
2    8.1
hive> select id, get_json_object(data, '$.store.fruit[1].weight') from json_data;
1    9
2    9.2
```

- 查询多层值3

如何同时获得weight和type值，有下面两种方式，效果一致。

方法1：和上一节一样，用数组方式
```
hive> select id, get_json_object(data, '$.store.fruit[0].weight'), get_json_object(data, '$.store.fruit[0].type')
    > from json_data;
  1    8   apple
  2    8.1 apple
```

方法2：用substr将数组的方括号截掉，转换为json
```
hive> select id,
    > get_json_object(substr(get_json_object(data, '$.store.fruit'), 2, length(get_json_object(data, '$.store.fruit')) - 2), '$.weight'),
    > get_json_object(substr(get_json_object(data, '$.store.fruit'), 2, length(get_json_object(data, '$.store.fruit')) - 2), '$.type')
    > from json_data;
1    8   apple
2    8.1 apple
```

- 查询多层值4

按照上面的两种方式，我们取到了fruit数组中第一个数据。但第二个数据，只能得到下面的效果（你可以试试看）：

```
1 8 apple 9 pear
2 8.1 apple 9.2 pear
```

能不能将相同的数排在一列，做出下面这样的效果，答案是肯定的。

```
1 8 apple
1 9 pear
2 8.1 apple
2 9.2 pear
```

思路是：取到fruit之后，substr截掉前后的方括号，使用split按照'},'对其分割，之后用explode行转列，再补全成完整的json，然后按照处理json的方式取数。步骤比较复杂，我们分三步看。

```
  #步骤1：截掉方括号，并用'},'分割，注意此时一三行不是完整的json，缺了右括号
  hive> select id, fruit
      > from json_data
      > lateral view explode(split(substr(get_json_object(data, '$.store.fruit'), 2, length(get_json_object(data, '$.store.fruit')) - 2), '},')) t as fruit;
  1    {"weight":8,"type":"apple"
  1    {"weight":9,"type":"pear"}
  2    {"weight":8.1,"type":"apple"
  2    {"weight":9.2,"type":"pear"}

  #步骤2：case when 补全json，可以看到一三行结果是json格式了
  hive> select a.id, case when substr(a.fruit, length(fruit), 1) = "}" then a.fruit else concat(a.fruit, '}') end as fruit_info
      > from
      > (
      >     select id, fruit
      >     from json_data
      >     lateral view explode(split(substr(get_json_object(data, '$.store.fruit'), 2, length(get_json_object(data, '$.store.fruit')) - 2), '},')) t as fruit
      > ) a;
  1    {"weight":8,"type":"apple"}
  1    {"weight":9,"type":"pear"}
  2    {"weight":8.1,"type":"apple"}
  2    {"weight":9.2,"type":"pear"}

  #步骤3：提取weight，type数据
  hive> select b.id, get_json_object(b.fruit_info, '$.weight'), get_json_object(b.fruit_info, '$.type')
      > from
      > (
      >     select a.id, case when substr(a.fruit, length(fruit), 1) = "}" then a.fruit else concat(a.fruit, '}') end as fruit_info
      >     from
      >     (
      >         select id, fruit
      >         from json_data
      >         lateral view explode(split(substr(get_json_object(data, '$.store.fruit'), 2, length(get_json_object(data, '$.store.fruit')) - 2), '},')) t as fruit
      >     ) a
      > ) b;
  1    8   apple
  1    9   pear
  2    8.1 apple
  2    9.2 pear
```

上面的步骤3，由于要获取两个字段，可以使用json_tuple函数代替，与later view连用，写法如下：

```
  hive> select b.id, c.weight, c.type
      > from
      > (
      >     select a.id, case when substr(a.fruit, length(fruit), 1) = "}" then a.fruit else concat(a.fruit, '}') end as fruit_info
      >     from
      >     (
      >         select id, fruit
      >         from json_data
      >         lateral view explode(split(substr(get_json_object(data, '$.store.fruit'), 2, length(get_json_object(data, '$.store.fruit')) - 2), '},')) t as fruit
      >     ) a
      > ) b
      > lateral view json_tuple(b.fruit_info,'weight', 'type') c as weight, type;
  1    8   apple
  1    9   pear
  2    8.1 apple
  2    9.2 pear
```

### 2.parse_url，regexp_replace，regexp_extract

```
语法: parse_url(string urlString, string partToExtract , string keyToExtract)
返回值: string
说明：返回 URL 中指定的部分。 partToExtract 的有效值为： HOST, PATH, QUERY, REF,
PROTOCOL, AUTHORITY, FILE, and USERINFO.

语法: regexp_replace(string A, string B, string C)
返回值: string
说明：将字符串 A 中的符合正则表达式 B 的部分替换为 C。

语法: regexp_extract(string subject, string pattern, int index)
返回值: string
说明：将字符串 subject 按照 pattern 正则表达式的规则拆分，返回 index 指定的字符。
```

下面看实例：

```
--我们虚构的数据，urldata.txt
1    https://ty.facebook.com/dwd/social?type=1&query=abc&id=1234
2    http://qq.tencent.com/dwd/category?type=2&query=def&id=5678#title1

hive> create table url_data(id int, data string) row format delimited fields terminated by '\t';
hive> load data local inpath 'urldata.txt' into table url_data;

#获取url协议
hive> select id, parse_url(data, 'PROTOCOL') from url_data;
1    https
2    http

#获取主机名
hive> select id, parse_url(data, 'HOST') from url_data;
1    ty.facebook.com
2    qq.tencent.com

#获取path
hive> select id, parse_url(data, 'PATH') from url_data;
1    /dwd/social
2    /dwd/category

#获取所有参数的序列
hive> select id, parse_url(data, 'QUERY') from url_data;
1    type=1&query=abc&id=1234
2    type=2&query=def&id=5678

#获取完整文件路径
hive> select id, parse_url(data, 'FILE') from url_data;
1    /dwd/social?type=1&query=abc&id=1234
2    /dwd/category?type=2&query=def&id=5678

#获取REF，没有的返回NULL值
hive> select id, parse_url(data, 'REF') from url_data;
1    NULL
2    title1
```

插曲：获取的参数序列是键值对的形式，能否将其拆分开呢？可以使用str_to_map函数.

```
语法: str_to_map(text, delimiter1, delimiter2)
返回值: map
说明：将字符串按照给定的分隔符转换成 map 结构。第一个分隔符在K-V之间分割，第二个分隔符分割K-V本身
hive> select id, parse_url(data, 'PROTOCOL'), parse_url(data, 'HOST'), parse_url(data, 'PATH'), str_to_map(parse_url(data, 'QUERY'), '&', '=')['type'],
    > str_to_map(parse_url(data, 'QUERY'), '&', '=')['query'], str_to_map(parse_url(data, 'QUERY'), '&', '=')['id']
    >  from url_data;
1    https   ty.facebook.com /dwd/social     1   abc 1234
2    http    qq.tencent.com  /dwd/category   2   def 5678
```

如果不使用parse_url，能否对相应的url子串进行截取，可以借助于regexp_extract，regexp_replace，不过可能调正则表达式需要用点功夫。

```
hive> select id, regexp_replace(data, 'dwd.+', "")
    > from url_data;
1    https://ty.facebook.com/
2    http://qq.tencent.com/

hive> select id, regexp_extract(data, 'query=(.*)', 1)
    > from url_data;
1    abc&id=1234
2    def&id=5678#title1

hive> select id, regexp_extract(data, 'query=(.*)&', 1)
    > from url_data;
1    abc
2    def
```

### 3.collect_set，collect_list，concat，concat_ws

```
语法: collect_set (col)
返回值: array
说明: 将 col 字段进行去重， 合并成一个数组。
语法: collect_list (col)
返回值: array
说明: 将 col 字段合并成一个数组,不去重
语法: concat(string A, string B…)
返回值: string
说明：返回输入字符串连接后的结果，支持任意个输入字符串
语法: concat_ws(string SEP, string A, string B…)
返回值: string
说明：返回输入字符串连接后的结果， SEP 表示各个字符串间的分隔符
```

```
--虚构的数据，fruitdata.txt
1001    apple
1001    pear
1001    banana
1001    pear
1002    blueberry
1002    bayberry

hive> create table fruit_data(id int, data string) row format delimited fields terminated by '\t';
hive> load data local inpath 'fruitdata.txt' into table fruit_data;

hive> select id, collect_set(data)
    > from fruit_data
    > group by id;
1001    ["apple","pear","banana"]
1002    ["blueberry","bayberry"]

hive> select id, collect_list(data)
    > from fruit_data
    > group by id;
1001    ["apple","pear","banana","pear"]
1002    ["blueberry","bayberry"]

--虚构的数据，userdata.txt。想想一个用户的粉丝在各个地域的分布情况
1001    area1   5%
1001    area2   20%
1001    area3   25%
1001    area4   50%
2001    area1   20%
2001    area2   50%
2001    area3   30%
hive> create table user_data(id int, area string, data string) row format delimited fields terminated by '\t';
hive> load data local inpath 'userdata.txt' into table user_data;

#按照每个用户一行进行排列
hive> select id, collect_set(concat_ws(':', area, data))
    > from user_data
    > group by id;
1001    ["area1:5%","area2:20%","area3:25%","area4:50%"]
2001    ["area1:20%","area2:50%","area3:30%"]

#下面使用concat能得到同样效果
hive> select id, collect_set(concat(area, ':', data))
    > from user_data
    > group by id;
1001    ["area1:5%","area2:20%","area3:25%","area4:50%"]
2001    ["area1:20%","area2:50%","area3:30%"]

#我们可以看到结果中，collect_set函数为我们加上了中括号和双引号，能不能去掉它们，我们来看下面的效果：
hive> select id, concat_ws(',', collect_set(concat( area, ':', data)))
    > from user_data
    > group by id;
1001    area1:5%,area2:20%,area3:25%,area4:50%
2001    area1:20%,area2:50%,area3:30%

#如果想变成map的格式，在此基础上可以再调用一下str_to_map即可
hive> select id, str_to_map(concat_ws(',', collect_set(concat( area, ':', data))), ",", ":")
    > from user_data
    > group by id;
1001    {"area1":"5%","area2":"20%","area3":"25%","area4":"50%"}
2001    {"area1":"20%","area2":"50%","area3":"30%"} 
```

### 4.datediff，from_unixtime，unix_timestamp，to_date

```
语法: datediff(string enddate, string startdate)
返回值: int
说明: 返回结束日期减去开始日期的天数。日期的格式需要是yyyy-MM-dd，或者yyyy-MM-dd HH:mm:ss

语法: from_unixtime(bigint unixtime[, string format])
返回值: string
说明: 转化 UNIX 时间戳（从 1970-01-01 00:00:00 UTC 到指定时间的秒数）到当前时区的时间格式，默认的format是yyyy-MM-dd HH:mm:ss，可以指定别的

语法: unix_timestamp(string date[, string format])
返回值: bigint
说明: 转换 pattern 格式的日期到 UNIX 时间戳。如果转化失败，则返回 0。默认的format是yyyy-MM-dd HH:mm:ss，可以指定别的。

语法: to_date(string timestamp)
返回值: string
说明: 返回日期时间字段中的日期部分。
```

下面看实例：

```
#虚构的数据datedata.txt，一共有10列，后9列是各种日期。
1    2019-02-03  20190203    2019-02-05  20190205    2019-03-03 10:38:24 2019-03-23 10:36:54 20190323 10:36:54   1551763940  1551763940267
2    2019-02-08  20190208    2019-02-18  20190218    2019-03-19 10:32:04 2019-03-31 10:39:15 20190331 10:39:15   1552632321  1551763940654

hive> create table date_data(id int, d1 string, d2 string, d3 string, d4 string, d5 string, d6 string, d7 string, d8 string, d9 string) row format delimited fields terminated by '\t';
hive> load data local inpath 'datedata.txt' into table date_data;
先看datediff的用法：

#yyyy-MM-dd的日期差
hive> select id, datediff(d3, d1)from date_data;
1    2
2    10
hive> select id, datediff(to_date(d3), to_date(d1)) from date_data;
1    2
2    10

#yyyyMMdd的日期差
hive> select id, datediff(d4, d2) from date_data;
1    NULL
2    NULL
#上面的写法不行，我们需要将日期转换为yyyy-MM-dd格式，使用截取拼接的套路进行
hive> select datediff(concat_ws('-', substr(d4, 1, 4), substr(d4, 5, 2), substr(d4, 7, 2)), concat_ws('-', substr(d2, 1, 4), substr(d2, 5, 2), substr(d2, 7, 2)))
    > from date_data;
2
10

#yyyy-MM-dd HH:mm:ss与yyyy-MM-dd的日期差
hive> select datediff(d5, d1) from date_data;
28
39
hive> select datediff(to_date(d5), d1) from date_data;
28
39

#yyyy-MM-dd HH:mm:ss与yyyy-MM-dd HH:mm:ss的日期差
hive> select datediff(d6, d5) from date_data;
20
12
hive> select datediff(to_date(d6), d5) from date_data;
20
12
```

再来看unix_timestamp的用法：

```
#yyyy-MM-dd HH:mm:ss转换为时间戳
hive> select unix_timestamp(d5) from date_data;
1551580704
1552962724
hive> select unix_timestamp(d5, 'yyyy-MM-dd HH:mm:ss') from date_data;
1551580704
1552962724

#yyyyMMdd HH:mm:ss转换为时间戳
hive> select unix_timestamp(d7, 'yyyyMMdd HH:mm:ss') from date_data;
1553308614
1553999955
最后看from_unixtime的用法：

#由于我们的表是string格式的，在转换之前需要转为bigint型
hive> select from_unixtime(cast(d8 as bigint)) from date_data;
2019-03-05 13:32:20
2019-03-15 14:45:21

hive> select from_unixtime(cast(d8 as bigint), 'yyyy-MM-dd HH:mm:ss') from date_data;
2019-03-05 13:32:20
2019-03-15 14:45:21

hive> select from_unixtime(cast(d8 as bigint), 'yyyyMMdd HH:mm:ss') from date_data;
20190305 13:32:20
20190315 14:45:21

hive> select from_unixtime(cast(d8 as bigint), 'yyyy-MM-dd') from date_data;
2019-03-05
2019-03-15

hive> select from_unixtime(cast(d8 as bigint), 'yyyyMMdd') from date_data;
20190305
20190315
```

我们经常会在业务中遇到13位的时间戳，10位的时间戳是精确到秒的，13位则是精确到毫秒的。这时只需除以1000并转化为整数即可。

```
hive> select from_unixtime(cast(d9/1000 as bigint)) from date_data;
2019-03-05 13:32:20
2019-03-05 13:32:20

hive> select from_unixtime(cast(d9/1000 as bigint), 'yyyy-MM-dd HH:mm:ss') from date_data;
2019-03-05 13:32:20
2019-03-05 13:32:20

hive> select from_unixtime(cast(d9/1000 as bigint), 'yyyyMMdd HH:mm:ss') from date_data;
20190305 13:32:20
20190305 13:32:20

hive> select from_unixtime(cast(d9/1000 as bigint), 'yyyy-MM-dd') from date_data;
2019-03-05
2019-03-05

hive> select from_unixtime(cast(d9/1000 as bigint), 'yyyyMMdd') from date_data;
20190305
20190305
```

### 5.coalesce

```
语法: COALESCE(T v1, T v2, …)
返回值: T
说明: 返回参数中的第一个非空值；如果所有值都为 NULL，那么返回 NULL
1    https://ty.facebook.com/dwd/social?type=1&query=abc&id=1234&task_id=1111
2    https://ty.facebook.com/dwd/social?type=1&query=abc&id=1234&taskid=2222

hive> create table exp_data(id int, data string) row format delimited fields terminated by '\t';
hive> load data local inpath 'expdata.txt' into table exp_data;
```

如果我们想提取出1111和2222这两个值，但一个是task_id，一个是taskid。如果直接用str_to_map，直接写的话，结果总会有一个空值：

```
hive> select str_to_map(data, '&', '=')['taskid'], str_to_map(data, '&', '=')['task_id']
    > from exp_data;
NULL    1111
2222    NULL
```

这个时候就可以用到coalesce

```
hive> select coalesce(str_to_map(data, '&', '=')['taskid'], str_to_map(data, '&', '=')['task_id'], "")
    > from exp_data;
1111
2222
```





