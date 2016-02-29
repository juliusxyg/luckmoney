1. 持久化到数据库
2. 记录日志方便分析
3. 幂等性
4. 打包数据
5. 定义结构化json接口
6. connection 的维护有问题


x M_envelop 需要加上 toJson 方法，方便存入mongo
x 不区分remain和opened，管理后台可能有用
x 刷入数据库，发红包，抢红包时, 用channel做？
? mongo index?
x clean内存表，已经抢光的话就可以用内存表中移除， 
? 或者长时间没人访问，从mongo读取，冷热数据
x 需要一些统计变量：红包数，抢光的红包，遗留的红包，等
？加个网关，网关处理多个服务器后端
x 优化uniqueid (each golang machine can only work with one process?!)
x 一个红包堆一次只能由一个人抓，保证原子性
x 判断抓到过的人不能再抓
x benchmark test (create red packet)
? test script to simulate multi-user playing
? refactor module

#########################################################################
? Share memory by communicating, don't communicate by sharing memory.
e.g.
package main

import "fmt"

type UpdateOp struct {
	key   int
	value string
}

func applyUpdate(data map[int]string, op UpdateOp) {
	data[op.key] = op.value
}

func main() {
	m := make(map[int]string)
	m[2] = "asdf"
	
	ch := make(chan UpdateOp)
	
	go func(ch chan UpdateOp) {
		ch <- UpdateOp{2, "New Value"}
	}(ch)
	
	applyUpdate(m, <- ch)
	fmt.Printf("%s\n", m[2])
}
#########################################################################

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