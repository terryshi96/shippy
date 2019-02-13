// 控制器
package main
import (
"log"
"context"
pb "github.com/terryshi96/shippy/consignment-service/proto"
vesselProto "github.com/terryshi96/shippy/vessel-service/proto"
	"gopkg.in/mgo.v2"
)
type service struct {
	session *mgo.Session
	vesselClient vesselProto.VesselServiceClient
}

//从效果上看，我们在创建了 master 会话之后，其实就没有再真正用过它了，
// 因为在之后的每次数据库请求中，我们都首先使用 Clone 生成了一个新的会话。
// 虽然我在代码中有过一段与之相关的注释，但我觉得有必要在这仔细讨论下原因。
// 当你每次只使用 master 会话来发起请求时时，在底层，你是在用同一个socket的同一个连接。
// 这意味着你的部分请求会被某个正在进行的请求阻塞，这是对 Golang 强大并发能力的浪费。

//为了不阻塞请求，mgo 支持使用 Copy() 或者 Clone() 来复制一个会话，这样你就能并发的处理请求了。
// Copy 和 Clone 功能尽管差不多，但有其细微且重要的区别。Clone 后的会话将使用和 master 会话相同的 socket，
// 但会使用一个新的连接，这既达到了并发的效果，还减少了新创一个 socket 的开销。这点非常适用于那些快速的写入操作。
// 但某些需要长时间处理的操作，比如复杂的询问，大数据操作等，可能会阻塞其他试图使用此 socket 的 goroutine。
// 而 Copy 则是会生成一个新的socket，相对 Clone, 它的开销就稍微大一点了。
func (s *service) GetRepo() Repository {
	// 模型实例化
	return &ConsignmentRepository{s.session.Clone()}
}

// 实现proto中定义的函数
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	repo := s.GetRepo()
	defer repo.Close()
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity: int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}
	// We set the VesselId as the vessel we got back from our
	// vessel service
	req.VesselId = vesselResponse.Vessel.Id
	// Save our consignment
	err = repo.Create(req)
	if err != nil {
		return err
	}
	// Return matching the `Response` message we created in our
	// protobuf definition.
	res.Created = true
	res.Consignment = req
	return nil
}
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	repo := s.GetRepo()
	defer repo.Close()
	consignments, err := repo.GetAll()
	if err != nil {
		return err
	}
	res.Consignments = consignments
	return nil
}
