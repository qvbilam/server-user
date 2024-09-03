package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"user/api"
	proto "user/api/qvbilam/user/v1"
	"user/global"
	"user/initialize"
	"user/utils"
)

func main() {
	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDatabase()
	initialize.InitRedis()
	initialize.InitElasticSearch()
	// 初始化 jaeger 链路追踪； 设置全局tracer (监听到 <-quit 后关闭
	tracer, jaegerCloser := initialize.InitJaeger()
	opentracing.SetGlobalTracer(tracer)

	// 注册服务(并且带入tracer
	//server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	server := grpc.NewServer(grpc.UnaryInterceptor(utils.ServerInterceptor(tracer)))
	proto.RegisterUserServer(server, &api.UserService{})
	proto.RegisterAccountServer(server, &api.AccountService{})

	Host := "0.0.0.0"
	Port, _ := utils.GetFreePort()
	Port = global.ServerConfig.Port
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", Host, Port))
	if err != nil {
		zap.S().Panicf("listen port error: %s", err)
	}

	zap.S().Infof("start %s server, host: %s:%d", global.ServerConfig.Name, Host, Port)
	go func() {
		if err := server.Serve(lis); err != nil {
			zap.S().Panicf("start server error: %s", err)
		}
	}()

	// 监听结束
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 关闭 jaeger 链路追踪
	_ = jaegerCloser.Close()
}
