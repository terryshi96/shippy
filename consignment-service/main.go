// 货运微服务 主要的功能是记录当前所有需要托运的集装箱，以及对应的货运船。

package main

import (
	// 使用go-micro替代
	micro "github.com/micro/go-micro"
	//"log"
	//"net"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	// 引入生成的consignment.pb.go文件
	pb "./proto/consignment"
	"golang.org/x/net/context"
	"fmt"

)

// 理解为model
// Repository相当于一张表 - 模拟一个数据库，我们会在此后使用真正的数据库替代他
type Repository struct {
	consignments []*pb.Consignment
}

// 接口是可被实例化的类型，用来实例化一个struct并赋予类的实例方法，相当于在mode里定义增删改查操作
type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// 为结构体类型定义方法, 实现接口, 相当于一个ruby类的实例方法，实现增删改查操作
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}


// 理解为controller
// service要实现在proto中定义的所有方法。当你不确定时
// 可以去对应的*.pb.go文件里查看需要实现的方法及其定义
type service struct {
	repo IRepository
}
// CreateConsignment - 在proto中，我们只给这个微服务定一个了一个方法
// 就是这个CreateConsignment方法，它接受一个context以及proto中定义的
// Consignment消息，这个Consignment是由gRPC的服务器处理后提供给你的

//interface所包含的方法不变，但是各个方法所接受的参数发生了变化，返回的参数也不同了。
// 原始的 gRPC 代码有四种不同的方法申明，对应四种不同的 gRPC 数据传输手段。
// 而 go-micro 统一了四种接口，抽象出了 req 和 res。
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}
	// 注意，在这里我们使用go-micro的NewService方法来创建新的微服务服务器，
	// 而不是上一篇文章中所用的标准
	// micro.NewService() 抽象出了原本复杂的流程。
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		// 注意，Name方法的必须是你在proto文件中定义的package名字
		// 用于服务发现
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	// Init方法会解析命令行flags
	srv.Init()
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}