package client

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/IguoChan/go-project/pkg/etcdx"
	"github.com/IguoChan/go-project/pkg/grpcx/resolver"

	"github.com/IguoChan/go-project/api/genproto/demo_app/simplepb"
)

func TestClientWithEtcdResolver(t *testing.T) {
	r, err := resolver.NewEtcdResolver("demo", &etcdx.Options{
		Addrs:       []string{"192.168.0.98:2379"},
		DialTimeout: 5 * time.Second,
		TTL:         30,
	})
	if err != nil {
		t.Fatal(err)
	}

	//resolver.SetDefaultScheme("dns")
	client, err := NewClient(&Config{
		//target: "lalala:443",
		//target: ":9410",
	}, WithEtcdDiscovery(r))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	grpcClient := simplepb.NewSimpleClient(client.Conn())
	req := simplepb.SimpleRequest{
		Data: "grpc stupid",
	}
	for i := 0; i < 10000; i++ {
		res, err := grpcClient.Route(context.Background(), &req)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
	ch := make(chan struct{})
	<-ch
}

func TestClientWithEmbedResolver(t *testing.T) {

	//resolver.SetDefaultScheme("dns")
	client, err := NewClient(&Config{
		//target: "lalala:443",
		target: "",
	}, SetEndpointsDiscovery([]string{":9412", ":9413", ":9414", ":9415"}))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	grpcClient := simplepb.NewSimpleClient(client.Conn())
	req := simplepb.SimpleRequest{
		Data: "grpc stupid1",
	}
	for i := 0; i < 100; i++ {
		res, err := grpcClient.Route(context.Background(), &req)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
	ch := make(chan struct{})
	<-ch
}

func TestTree(t *testing.T) {
	fmt.Println("Hello World!")

	leaf1 := &TreeNode{
		Val: 100,
	}
	leaf2 := &TreeNode{
		Val: 101,
	}

	tree := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val:   2,
			Left:  leaf1,
			Right: leaf2,
		},
		Right: &TreeNode{},
	}

	node := f(tree, leaf1, leaf2)
	fmt.Println(node.Val)
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func f(root, leaf1, leaf2 *TreeNode) *TreeNode {
	var isP1 func(root *TreeNode) bool
	isP1 = func(root *TreeNode) bool {
		if root == nil {
			return false
		}
		if root == leaf1 {
			return true
		}
		return isP1(root.Left) || isP1(root.Right)
	}

	var isP2 func(root *TreeNode) bool
	isP2 = func(root *TreeNode) bool {
		if root == nil {
			return false
		}
		if root == leaf2 {
			return true
		}
		return isP2(root.Left) || isP2(root.Right)
	}

	ans := make([]*TreeNode, 0)
	var isP func(root *TreeNode)
	isP = func(root *TreeNode) {
		if root == nil {
			return
		}

		if isP1(root) && isP2(root) {
			ans = append(ans, root)
		}

		isP(root.Left)
		isP(root.Right)
	}

	isP(root)
	return ans[len(ans)-1]
}

func TestWithLogger1(t *testing.T) {
	aa := []int{5, 1, 4, 2}

	sema := semaphore.NewWeighted(3)
	wg := sync.WaitGroup{}
	for i := 0; i < len(aa); i++ {
		wg.Add(1)
		err := sema.Acquire(context.Background(), 1)
		if err != nil {
			t.Error(err)
			continue
		}
		idx := i
		go func() {
			defer wg.Done()
			defer sema.Release(1)
			f1(aa[idx])
		}()
	}
	wg.Wait()
}

func f1(a int) int {
	time.Sleep(time.Duration(a) * time.Second)
	fmt.Println(a)
	return a
}

func TestPaste(t *testing.T) {
	a1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	a2 := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i'}

	c1 := make(chan struct{})
	c2 := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(2)

	// 打印 a1
	go func() {
		defer wg.Done()
		for _, v := range a1 {
			fmt.Println(v)
			c1 <- struct{}{}
			<-c2
		}
	}()

	// 打印a2
	go func() {
		defer wg.Done()
		for _, v := range a2 {
			<-c1
			fmt.Printf("%c\n", v)
			c2 <- struct{}{}
		}
	}()

	wg.Wait()
}

func TestDIDI(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6}
	b := f3(a)
	t.Log(b)
}

func f3(a []int) []int {
	p1, p2 := 0, 0
	for p2 < len(a) {
		if a[p2]%2 == 1 {
			a[p1], a[p2] = a[p2], a[p1]
			p1++
		}
		p2++
	}

	return a
}
