// 货船服务  根据集装箱的重量，数量去寻找合适的货运船
package main

import (
	"errors"
	"fmt"
	//与本地路径一致 dockerfile中workdir也要与此一致
	pb "shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"context"
)

// 定义Model
type VesselRepository struct {
	vessels []*pb.Vessel
}

// 定义model类方法
type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

// 实现model类方法
// FindAvailable - 根据Specification，从若干货船中挑选出合适的货船来运送货物
// 如果货物的数量和重量都没有超过一个货船的数量和重量上限，
// 那么我们就返回这个货船
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	// 遍历所有的货船
	for _, vessel := range repo.vessels {
		// 返回数量和重量上限量大于需求的货船
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("No vessel found by that spec")
}

// 定义controller，在controller方法中调用model类方法
type service struct {
	repo Repository
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}
	res.Vessel = vessel
	return nil
}
func main() {

	// 创建一套货船记录 后续使用数据库代替
	vessels := []*pb.Vessel{
		&pb.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselRepository{vessels}
	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)
	srv.Init()
	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
