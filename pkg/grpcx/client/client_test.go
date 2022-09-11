package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/IguoChan/go-project/api/genproto/demo_app/simplepb"
)

func TestClient(t *testing.T) {
	//r, err := resolver.NewEtcdResolver("demo", &etcdx.Options{
	//	Addrs:       []string{"192.168.0.98:2379"},
	//	DialTimeout: 5 * time.Second,
	//	TTL:         30,
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}

	client, err := NewClient(&Config{
		target: "lalala:9410",
	}) //, SetEndpointsDiscovery([]string{":9411", ":9411", ":9412", ":9413", ":9414", ":9415", ":9416", ":9410", ":9418", ":9419", ":9420", ":9421", ":9422"}))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	grpcClient := simplepb.NewSimpleClient(client.Conn())
	req := simplepb.SimpleRequest{
		Data: "grpc stupid",
	}
	for i := 0; i < 100; i++ {
		res, err := grpcClient.Route(context.Background(), &req)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
}
