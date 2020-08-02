package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	greeter "helloword/code/go/go_grpc_template/api/greeter/go"
	"helloword/code/go/go_grpc_template/client/internal/app/greeter_client"
	"time"
)

func main() {

	var (
		discover = greeter_client.ConsulDiscover()
	)

	node := discover.DiscoverNode()
	if node == nil || len(node.Address) == 0 {
		fmt.Println(node)
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
}
