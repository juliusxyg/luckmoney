1. 持久化到数据库
2. 记录日志方便分析
3. 幂等性
4. 打包数据
5. 定义结构化json接口
6. connection 的维护有问题


x M_envelop 需要加上 toJson 方法，方便存入mongo
x 不区分remain和opened，管理后台可能有用
? 刷入数据库，发红包，抢红包时, 用channel做？mongo index?
? clean内存表，已经抢光的话就可以用内存表中移除， 或者长时间没人访问，从mongo读取，冷热数据
需要一些统计变量：红包数，抢光的红包，遗留的红包，等

mongo 数据库结构

红包
{ID, money, size, created_at, opened [piece1, piece2, ..., pieceN]}
pieceN: {i, money, grabber, grabtime}
//可能要重构部分roll.go的代码， 这样可以实现定时刷新到mongo，piece太多怎么办？
//一大段json么，不过异步的

不需要linked list，slice就能满足需求， 链表用在需要在中间或头部插入元素的大列表情况下

//mongo
~: mongod --dbpath <path>

~: mongo
~: use <db name>
~: db.<collection name>.count()
~: db.<collection name>.find({"id": NumberLong("<envelop id>")})
~: db.<collection name>.find({"id": NumberLong("<envelop id>"), "pieces.grabber":"<grabber name>"}, {"pieces.$.grabtime": 1})
~: db.<collection name>.find({"id":NumberLong("<envelop id>")}, {"pieces":{ $elemMatch: {"grabber":"<grabber name>"}}})