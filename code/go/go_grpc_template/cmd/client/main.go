package main

import (
	"context"
	"fmt"
	"github.com/bj-wangjia/go-kit/balancer"
	"google.golang.org/grpc"
	greeter "helloword/code/go/go_grpc_template/api/greeter/go"
	"time"
)

func main() {
	consulResolver, err := balancer.NewConsulResolver(
		":8500",
		"go-kit-srv-greeter",
		"go-kit-client-greeter",
		time.Duration(1)*time.Second,
		3,
		60,
		"",
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	node := consulResolver.DiscoverNode()
	if node == nil || len(node.Address) == 0 {
		return
	}
	fmt.Println(node)

	ctxConn, celConn := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(100))
	defer celConn()
	conn, err := grpc.DialContext(ctxConn, node.Address, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dial context error:", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	c := greeter.NewGreeterClient(conn)

	ctxRead, celRead := context.WithTimeout(context.Background(), time.Duration(300)*time.Millisecond)
	defer celRead()

	response, err := c.Greeting(ctxRead, &greeter.GreetingRequest{Name: "go-kit"})
	if err != nil {
		fmt.Println("Read context error:", err)
		return
	}
	fmt.Println(response)

	time.Sleep(1 * time.Minute)
}
