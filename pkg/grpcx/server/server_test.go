package server

import (
	"testing"

	"github.com/IguoChan/go-project/internal/app/demo_app"
)

const (
	addr = ":8080"
)

func TestServer(t *testing.T) {
	//r, err := resolver.NewEtcdResolver("demo", &etcdx.Options{
	//	Addrs:       []string{"192.168.0.98:2379"},
	//	DialTimeout: 5 * time.Second,
	//	TTL:         30,
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}

	server, err := NewServer(&Config{
		Host: "192.168.0.101",
		Port: 9410,
	}, []PBServerRegister{demo_app.NewSimpleServer()})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	err = server.Serve()
	if err != nil {
		t.Fatal(err)
	}
}

func TestServer1(t *testing.T) {
	//r, err := resolver.NewEtcdResolver("demo", &etcdx.Options{
	//	Addrs:       []string{"192.168.0.98:2379"},
	//	DialTimeout: 5 * time.Second,
	//	TTL:         30,
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}

	server, err := NewServer(&Config{
		Host: "192.168.0.101",
		Port: 9411,
	}, []PBServerRegister{demo_app.NewSimpleServer()})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	err = server.Serve()
	if err != nil {
		t.Fatal(err)
	}
}
