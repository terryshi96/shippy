// 货运微服务 主要的功能是记录当前所有需要托运的集装箱，以及对应的货运船。

package main

import (
	// 使用go-micro替代
	//micro "github.com/micro/go-micro"
	//"net"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	// 引入生成的consignment.pb.go文件
	//pb "./proto/consignment"
	//"context"
	"log"
	"fmt"
	"os"
	pb "github.com/terryshi96/shippy/consignment-service/proto"
	vesselProto "github.com/terryshi96/shippy/vessel-service/proto"
	"github.com/terryshi96/shippy/common"
	"github.com/micro/go-micro"
)

// v1
//// 理解为model
//// Repository相当于一张表 - 模拟一个数据库，我们会在此后使用真正的数据库替代他
//type Repository struct {
//	consignments []*pb.Consignment
//}
//
//// 接口是可被实例化的类型，用来实例化一个struct并赋予类的实例方法，相当于在mode里定义增删改查操作
//type IRepository interface {
//	Create(*pb.Consignment) (*pb.Consignment, error)
//	GetAll() []*pb.Consignment
//}
//
//// 为结构体类型定义方法, 实现接口, 相当于一个ruby类的实例方法，实现增删改查操作
//func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
//	updated := append(repo.consignments, consignment)
//	repo.consignments = updated
//	return consignment, nil
//}
//
//func (repo *Repository) GetAll() []*pb.Consignment {
//	return repo.consignments
//}
//
//
//// 理解为controller
//// service要实现在proto中定义的所有方法。当你不确定时
//// 可以去对应的*.pb.go文件里查看需要实现的方法及其定义
//type service struct {
//	repo IRepository
//	// 请注意，我们在这里记录了一个货船服务的客户端对象，这里出现了一个微服务的方法调用了另一个微服务的方法
//	vesselClient vesselProto.VesselServiceClient
//}
//// CreateConsignment - 在proto中，我们只给这个微服务定一个了一个方法
//// 就是这个CreateConsignment方法，它接受一个context以及proto中定义的
//// Consignment消息，这个Consignment是由gRPC的服务器处理后提供给你的
//
////interface所包含的方法不变，但是各个方法所接受的参数发生了变化，返回的参数也不同了。
//// 原始的 gRPC 代码有四种不同的方法申明，对应四种不同的 gRPC 数据传输手段。
//// 而 go-micro 统一了四种接口，抽象出了 req 和 res。
//func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
//	// 这里，我们通过货船服务的客户端对象，向货船服务发出了一个请求
//	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
//		MaxWeight: req.Weight,
//		Capacity: int32(len(req.Containers)),
//	})
//	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
//	if err != nil {
//		return err
//	}
//
//	// 维护关联关系
//	req.VesselId = vesselResponse.Vessel.Id
//
//	// 存在可用的货船才能创建货运订单
//	consignment, err := s.repo.Create(req)
//	if err != nil {
//		return err
//	}
//	res.Created = true
//	res.Consignment = consignment
//	return nil
//}
//
//func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
//	consignments := s.repo.GetAll()
//	res.Consignments = consignments
//	return nil
//}
//
//func main() {
//	repo := &Repository{}
//	// 注意，在这里我们使用go-micro的NewService方法来创建新的微服务服务器，
//	// 而不是上一篇文章中所用的标准
//	// micro.NewService() 抽象出了原本复杂的流程。
//	srv := micro.NewService(
//		// This name must match the package name given in your protobuf definition
//		// 注意，Name方法的必须是你在proto文件中定义的package名字
//		// 用于服务发现
//		micro.Name("go.micro.srv.consignment"),
//		micro.Version("latest"),
//	)
//
//	// 我们在这里使用预置的方法生成了一个货船服务的客户端对象
//	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
//	// Init方法会解析命令行flags
//	srv.Init()
//	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})
//	if err := srv.Run(); err != nil {
//		fmt.Println(err)
//	}
//}

// v2
// todo 数据库连接配置考虑使用配置文件维护

const (
	defaultHost = "localhost:27017"
)
func main() {
	// Database host from the environment variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = defaultHost
	}
	session, err := common.CreateSession(host)
	// 确保在main退出前关闭会话
	defer session.Close()
	if err != nil {
		log.Panicf("Could not connect to datastore with host %s - %v", host, err)
	}
	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
	// Init will parse the command line flags.
	srv.Init()
	// Register handler
	// go 中所有的指针参数，都是对原值的直接操作，而不是值传递是新建一个变量保存值
	pb.RegisterShippingServiceHandler(srv.Server(), &service{session, vesselClient})
	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}