package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	greeter "helloword/code/go/go_grpc_template/api/greeter/go"
	"helloword/code/go/go_grpc_template/server/internal/app/greeter_server"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"
)

func main() {
	fs := flag.NewFlagSet("greetersvc", flag.ExitOnError)
	var (
		httpPort = fs.String("http.port", "9110", "HTTP Listen Port")
		grpcPort = fs.String("grpc-addr", "9120", "gRPC listen address")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	// 日志相关
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var service greeter_server.Service

	var (
		endpoints   = greeter_server.MakeServerEndpoints(service)
		httpHandler = greeter_server.NewHTTPHandler(endpoints)
		registar    = greeter_server.ConsulRegister(*grpcPort)
		grpcServer  = greeter_server.NewGRPCServer(endpoints)
	)

	var g group.Group
	{
		// The service discovery registration.
		httpServer := &http.Server{
			Addr:           ":" + *httpPort,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Minute,
			MaxHeaderBytes: 1 << 20,
			Handler:        httpHandler,
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "port", *httpPort)
			return httpServer.ListenAndServe()
		}, func(err error) {
			_ = httpServer.Close()
		})
	}
	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", ":"+*grpcPort)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			registar.Register()
			logger.Log("transport", "gRPC", "port", *grpcPort)
			baseServer := grpc.NewServer()
			greeter.RegisterGreeterServer(baseServer, grpcServer)
			grpc_health_v1.RegisterHealthServer(baseServer, &greeter_server.Health{})
			return baseServer.Serve(grpcListener)
		}, func(err error) {
			registar.Deregister()
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(err error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
