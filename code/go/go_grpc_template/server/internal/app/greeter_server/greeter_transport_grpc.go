package greeter_server

import (
	"context"
	"errors"
	greeter "helloword/code/go/go_grpc_template/api/greeter/go"
	"reflect"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type (
	grpcServer struct {
		greeter grpctransport.Handler
	}
)

// NewGRPCServer 使一组端点可用作 gRPC Greeting 服务器。
func NewGRPCServer(endpoints Endpoints) greeter.GreeterServer {
	options := []grpctransport.ServerOption{

	}

	return &grpcServer{
		greeter: grpctransport.NewServer(
			endpoints.GreetingEndpoint,
			decodeGRPCGreetingRequest,
			encodeGRPCGreetingResponse,
			options...),
	}
}

// decodeGRPCGreetingRequest 是一个 transport/grpc.DecodeRequestFunc 将 gRPC Greeting 请求转换为用户域 Greeting 请求
func decodeGRPCGreetingRequest(context context.Context, grpcReq interface{}) (request interface{}, err error) {
	req, ok := grpcReq.(*greeter.GreetingRequest)
	if !ok {
		return nil, errors.New("grpcReq type should be (*greeter.GreetingRequest), but actual is " + reflect.TypeOf(grpcReq).String())
	}
	if grpcReq == nil {
		return nil, errors.New("grpcReq should not be nil")
	}
	if req.Name == "" {
		return nil, errors.New("grpcReq.Name should not be \"\" ")
	}
	return grpcReq, nil
}

// encodeGRPCGreetingResponse 是一个 transport/grpc.EncodeResponseFunc 将用户域
// 问Greeting响应转换为 gRPC Greeting 响应。
func encodeGRPCGreetingResponse(i context.Context, grpcResp interface{}) (response interface{}, err error) {
	_, ok := grpcResp.(*greeter.GreetingResponse)
	if !ok {
		return nil, errors.New("grpcResp type should be (*greeter.GreetingResponse), but actual is " + reflect.TypeOf(grpcResp).String())
	}
	return grpcResp, nil
}

// 实现 GreeterService.Greeting 接口
func (s *grpcServer) Greeting(ctx context.Context, req *greeter.GreetingRequest) (*greeter.GreetingResponse, error) {
	_, res, err := s.greeter.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*greeter.GreetingResponse), nil
}
