package sum

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	grpc_template "helloword/code/go/go_grpc_template/api/sum/go"
	"helloword/code/go/go_grpc_template/internal/pkg/grpc_logger"
	"net"
	"runtime/debug"
	"time"
)

type RpcServer struct {
	listen net.Listener
	server *grpc.Server

	RetargetHandler grpc_transport.Handler
}

//NewRpcServer
//创建基于 go-kit 的 rpc 服务，MakeRetargetDspEndpoint 是主要业务逻辑实现。
//在 go-kit 推荐使用中，decodeGRPCRequest 和 encodeGRPCResponse 两部分应承担请求到 rtContext 和 rtContext 到响应的两个 pipeline。这里未实现。
func NewRpcServer() (*RpcServer, error) {
	var err error

	//build Transport
	service := new(RpcServer)
	service.RetargetHandler = grpc_transport.NewServer(
		//增加速率限制，每50毫秒补充一次，设置容量12
		NewTokenBucketLimiterWithBuildIn(rate.NewLimiter(rate.Every(time.Millisecond*time.Duration(50)), 12),
		)(MakeRetargetDspEndpoint(&SumHandler{})),
		decodeGRPCRequest,
		encodeGRPCResponse,
	)

	service.server = grpc.NewServer()
	grpc_template.RegisterRpcTemplateServiceServer(service.server, service)

	service.listen, err = net.Listen("tcp", ":"+"13000")

	return service, err
}

//Start
//启动 rpc 服务
func (rs *RpcServer) Start(context.Context) error {
	if err := rs.server.Serve(rs.listen); err != nil {
		return fmt.Errorf("failed to serve: %s", err.Error())
	}
	return nil
}

//使用 go-lit 推荐将其参数透传给 ServeGRPC
func (rs *RpcServer) Sum(ctx context.Context, request *grpc_template.RpcTemplateRequest) (*grpc_template.RpcTemplateResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RpcServer Sum error:", err)
			debug.PrintStack()
		}
	}()

	//go-kit 的 rpc 服务函数，执行 decodeGRPCRequest、MakeRetargetDspEndpoint、encodeGRPCResponse 流程
	_, response, err := rs.RetargetHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*grpc_template.RpcTemplateResponse), err
}

func decodeGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	if request == nil {
		return nil, errors.New("RpcServer request is nil")
	}

	return &SumHandler{
		Items: request.(*grpc_template.RpcTemplateRequest).Items,
	}, nil
}

func encodeGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	if response == nil {
		return nil, errors.New("RpcServer response is nil")
	}

	return &grpc_template.RpcTemplateResponse{
		Result: response.(*SumHandler).Result,
	}, nil
}

// [2]中间件：限流服务
var ErrLimitExceed = errors.New("limit exceed error")

// NewTokenBucketLimiterWithBuildIn 使用x/time/rate创建限流中间件
func NewTokenBucketLimiterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

//MakeRetargetDspEndpoint
//主要业务逻辑实现代码，触发 Rt Serve Pipeline，记录请求日志
func MakeRetargetDspEndpoint(srv *SumHandler) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (response interface{}, err error) {
		handler := req.(*SumHandler)

		//pipeline
		handler.Sum()

		//write log
		handler.WriteLog()

		return handler, err
	}
}

type SumHandler struct {
	Items []int32

	Result int32

	elapsed string
}

func (sh *SumHandler) Sum() {

	defer func(begin time.Time) {
		sh.elapsed = time.Since(begin).String()
	}(time.Now())

	for i, _ := range sh.Items {
		sh.Result += sh.Items[i]
	}
}

func (sh *SumHandler) WriteLog() {
	grpc_logger.Grpc_logger.Req(map[string]interface{}{
		"func":    "Sum",
		"elapsed": sh.elapsed,
		"items":   sh.Items,
		"reslut":  sh.Result,
	})
}
