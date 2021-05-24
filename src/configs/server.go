package configs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

const (
	EnvAppIpKey            = "APP_IP"
	EnvGrpcPortKey         = "APP_GRPC_PORT"
	EnvHttpPortKey         = "APP_HTTP_PORT"
	EnvLogLevel            = "ENV_LOG_LEVEL"
	EnvRateLimiterBurstKey = "RATE_LIMITER_BURST"
	EnvRateLimiterRateKey  = "RATE_LIMITER_RATE"
)

type ServerConfigs struct {
	Http        HttpConfigs
	Grpc        GrpcConfigs
	RateLimiter RateLimiter
	Config
}

type HttpConfigs struct {
	Ip   string
	Port int
	Config
}

type GrpcConfigs struct {
	Ip   string
	Port int
	Config
}

type RateLimiter struct {
	Burst int
	Rate  rate.Limit
	Config
}

func (serverConfigs *ServerConfigs) Format() string {
	return fmt.Sprintf(
		"Server running!!! http: %s | gRPC: %s | %s",
		serverConfigs.Http.Format(),
		serverConfigs.Grpc.Format(),
		serverConfigs.RateLimiter.Format(),
	)
}
func GetServerConfigs() ServerConfigs {
	return ServerConfigs{
		Http:        getHttpConfigs(),
		Grpc:        getGrpcConfigs(),
		RateLimiter: GetRateLimiterConfigs(),
	}
}

func (httpConfigs *HttpConfigs) Format() string {
	return fmt.Sprintf(
		"%s:%d",
		httpConfigs.Ip,
		httpConfigs.Port,
	)
}
func getHttpConfigs() HttpConfigs {
	return HttpConfigs{
		Ip:   GetEnvStr(EnvAppIpKey, ""),
		Port: GetEnvInt(EnvHttpPortKey, 80),
	}
}

func (grpcConfigs *GrpcConfigs) Format() string {
	return fmt.Sprintf(
		"%s:%d",
		grpcConfigs.Ip,
		grpcConfigs.Port,
	)
}
func getGrpcConfigs() GrpcConfigs {
	return GrpcConfigs{
		Ip:   GetEnvStr(EnvAppIpKey, ""),
		Port: GetEnvInt(EnvGrpcPortKey, 9090),
	}
}

func (rateLimiterConfigs *RateLimiter) Format() string {
	return fmt.Sprintf(
		"rate limit: %.0f/%d",
		rateLimiterConfigs.Rate,
		rateLimiterConfigs.Burst,
	)
}
func GetRateLimiterConfigs() RateLimiter {
	return RateLimiter{
		Burst: GetEnvInt(EnvRateLimiterBurstKey, 5),
		Rate:  rate.Limit(GetEnvInt(EnvRateLimiterRateKey, 2)),
	}
}

func GetLogLevel() logrus.Level {
	defaultLevel, _ := logrus.Level.MarshalText(logrus.InfoLevel)
	level, _ := logrus.ParseLevel(GetEnvStr(EnvLogLevel, string(defaultLevel)))
	return level
}
