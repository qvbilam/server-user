package initialize

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"user/global"
)

func InitJaeger() (opentracing.Tracer, io.Closer) {
	// 配置
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,                          // 始终做出相同决定的采样器类型。
			Param: float64(global.ServerConfig.JaegerConfig.Output), // 0 不采样, 1全部采样
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           global.ServerConfig.JaegerConfig.IsLog, // 是否打印日志
			LocalAgentHostPort: fmt.Sprintf("%s:%s", global.ServerConfig.JaegerConfig.Host, global.ServerConfig.JaegerConfig.Port),
		},
		ServiceName: global.ServerConfig.JaegerConfig.Server,
	}
	// 生成链路Tracer
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(any(err))
	}
	// 关闭链路(在main 收到 <-quit 后调用
	return tracer, closer
}
