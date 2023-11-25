package server

import (
	"strconv"
	"testing"
	"time"

	"github.com/IguoChan/go-project/pkg/etcdx"
	"github.com/IguoChan/go-project/pkg/grpcx/resolver"

	"github.com/IguoChan/go-project/internal/app/demo_app"
)

const (
	addr = ":8080"
)

func TestServerWithEtcdResolver(t *testing.T) {
	r, err := resolver.NewEtcdResolver("demo", &etcdx.Options{
		Addrs:       []string{"192.168.0.98:2379"},
		DialTimeout: 5 * time.Second,
		TTL:         30,
	})
	if err != nil {
		t.Fatal(err)
	}

	server, err := NewServer(&Config{
		Host: "",
		Port: 9410,
	}, []PBServerRegister{demo_app.NewSimpleServer()}, WithEtcdRegister(r))
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	err = server.Serve()
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerWithEtcdResolver1(t *testing.T) {
	r, err := resolver.NewEtcdResolver("demo", &etcdx.Options{
		Addrs:       []string{"192.168.0.98:2379"},
		DialTimeout: 5 * time.Second,
		TTL:         30,
	})
	if err != nil {
		t.Fatal(err)
	}

	server, err := NewServer(&Config{
		Host: "",
		Port: 9412,
	}, []PBServerRegister{demo_app.NewSimpleServer()}, WithEtcdRegister(r))
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	err = server.Serve()
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerWithEmbedResolver(t *testing.T) {
	server, err := NewServer(&Config{
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

func TestServerWithEndpointsResolver(t *testing.T) {
	serve := func(port int) {
		server, err := NewServer(&Config{
			Port: port,
		}, []PBServerRegister{demo_app.SimpleServer{Ext: strconv.Itoa(port)}})
		if err != nil {
			t.Fatal(err)
		}
		defer server.Stop()

		err = server.Serve()
		if err != nil {
			t.Fatal(err)
		}
	}

	go serve(9412)
	go serve(9413)
	go serve(9414)
	go serve(9415)

	ch := make(chan struct{})
	<-ch
}
