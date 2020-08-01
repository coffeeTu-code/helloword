package greeter_server

import (
	"context"
	"encoding/json"
	"errors"
	greeter "helloword/code/go/go_grpc_template/api/greeter/go"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// 当缺少预期的路径变量时，将返回 ErrBadRouting。
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

// NewHTTPHandler返回一个使一组端点在预定义路径上可用的HTTP处理程序。
func NewHTTPHandler(endpoints Endpoints) http.Handler {
	m := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	// GET /health         获取服务健康信息
	// GET /greeting?name  获取 greeting

	// NewServer 方法需要端点，解码器，编码器作为参数
	m.Methods("GET").Path("/health").Handler(
		httptransport.NewServer(
			endpoints.HealthEndpoint,
			DecodeHTTPHealthRequest,
			EncodeHTTPHealthResponse,
			options...),
	)
	m.Methods("GET").Path("/greeting").Handler(
		httptransport.NewServer(
			endpoints.GreetingEndpoint,
			DecodeHTTPGreetingRequest,
			EncodeHTTPGreetingResponse,
			options...),
	)
	return m
}

// 解码 Health HTTP 请求的方法
func DecodeHTTPHealthRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return HealthRequest{}, nil
}

// 解码 Health HTTP 请求的方法
func EncodeHTTPHealthResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(writer).Encode(response)
}

// 解码 Greeting HTTP 请求的方法
func DecodeHTTPGreetingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := r.URL.Query()
	names, exists := vars["name"]
	if !exists || len(names) != 1 {
		return nil, ErrBadRouting
	}
	req := greeter.GreetingRequest{Name: names[0]}
	return req, nil
}

func EncodeHTTPGreetingResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	if f, ok := response.(Failer); ok && f.Failed() != nil {
		encodeError(ctx, f.Failed(), writer)
		return nil
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(writer).Encode(response)
}

// errorWrapper 将 error 封装为一个 json 结构体方便转换为 json
type errorWrapper struct {
	Error string `json:"error"`
}

// 编码错误的方法
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

// err2code 函数将 error 转换为对应的 http 状态码
func err2code(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}
