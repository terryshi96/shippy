package main

// 客户端代码。在这里我们要一个命令行工具, 它会读取JSON文件并和我们的服务器交互
import (
	"log"
	"os"
	"golang.org/x/net/context"
	"encoding/json"
	"io/ioutil"
	pb "../consignment-service/proto/consignment"
	"github.com/micro/go-micro/cmd"
	microclient "github.com/micro/go-micro/client"
	//"grpc_demo/app/client"
	//"google.golang.org/grpc"
)
const (
	//address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

// 解析json文件，从文件中读取数据并返回，构造用于创建货运记录的数据，相当于向服务端post表单数据
func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}
func main() {
	// Set up a connection to the server.
	// 以不验证身份的方式连接grpc server
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("Did not connect: %v", err)
	//}
	//defer conn.Close()
	//client := pb.NewShippingServiceClient(conn)

	cmd.Init()
	// Create new greeter client
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)

	// Contact the server and print out its response.
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	// 调用grcp server的函数
	r, err := client.CreateConsignment(context.TODO(), consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
