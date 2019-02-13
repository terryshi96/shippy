// 模型
package main

import (
	pb "github.com/terryshi96/shippy/consignment-service/proto"
	"gopkg.in/mgo.v2"
)

// todo 使用配置文件
const (
	dbName = "shippy"
	consignmentCollection = "consignments"
)

// 相当于ruby定义模型的实例方法
type Repository interface {
	Create(*pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
	Close()
}

// 数据库表结构
type ConsignmentRepository struct {
	session *mgo.Session
}

// 实例方法实现
// Create a new consignment
func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) error {

	// 数据库插入数据
	return repo.collection().Insert(consignment)
}
// GetAll consignments
func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error) {
	var consignments []*pb.Consignment
	// Find()通常接受一个询问条件(query)，但我们想要所有的货运任务，所以在这里用nil
	// 然后把找到的所有货运任务通过All()赋值给consignment
	// 另外在mgo中，One可以处理单个结果

	// 数据库查询
	err := repo.collection().Find(nil).All(&consignments)
	return consignments, err
}
// Close closes the database session after each query has ran.
// Mgo creates a 'master' session on start-up, it's then good practice
// to copy a new session for each request that's made. This means that
// each request has its own database session. This is safer and more efficient,
// as under the hood each session has its own database socket and error handling.
// Using one main database socket means requests having to wait for that session.
// I.e this approach avoids locking and allows for requests to be processed concurrently. Nice!
// But... it does mean we need to ensure each session is closed on completion. Otherwise
// you'll likely build up loads of dud connections and hit a connection limit. Not nice!
// （我认为作者这里的描述不太准确，且有点混乱，故放上英文原文。）
// Close()会在所有的询问都结束后关闭数据库会话(session)
// Mgo会在程序启动时创建一个'master'会话
// 一个好习惯就是为每一个数据库请求复制一个新的会话
// 这即更安全也更有效率。
// 因为，在底层，每一个数据库会话都有他自己的数据库socket和错误的处理机制(handling)。
// 让每个请求都只使用同一个数据库socket，这意味着某些请求需要等待socket的使用权。
// 这即排除了锁死的可能，也能更好的并发处理数据库请求。
// 随之而来的是，我们要确保每个会话结束后关闭会话，不然等待我们的就是一大堆无用的连接了！
func (repo *ConsignmentRepository) Close() {
	repo.session.Close()
}

// 定义具体操作哪个数据库的哪张表
func (repo *ConsignmentRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(consignmentCollection)
}