package configs

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
	"os"
	"testing"
)

func TestGetHttpConfigs(t *testing.T) {
	tests := []struct {
		name string
		want HttpConfigs
	}{{
		name: "Default Http Configs",
		want: HttpConfigs{
			Ip:   "",
			Port: 80,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getHttpConfigs())
		})
	}
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		want     logrus.Level
		envValue string
	}{{
		name: "LogLevel defaults to INFO",
		want: logrus.InfoLevel,
	}, {
		name:     "LogLevel defaults to INFO",
		want:     logrus.TraceLevel,
		envValue: "trace",
	}, {
		name:     "LogLevel defaults to INFO",
		want:     logrus.DebugLevel,
		envValue: "debug",
	}, {
		name:     "LogLevel defaults to INFO",
		want:     logrus.WarnLevel,
		envValue: "warn",
	}, {
		name:     "LogLevel defaults to INFO",
		want:     logrus.ErrorLevel,
		envValue: "error",
	}, {
		name:     "LogLevel defaults to INFO",
		want:     logrus.FatalLevel,
		envValue: "fatal",
	}, {
		name:     "LogLevel defaults to INFO",
		want:     logrus.PanicLevel,
		envValue: "panic",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				err := os.Setenv(EnvLogLevel, tt.envValue)
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.want, GetLogLevel())
			if tt.envValue != "" {
				err := os.Unsetenv(EnvLogLevel)
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetServerConfigs(t *testing.T) {
	tests := []struct {
		name string
		want ServerConfigs
	}{{
		name: "Default Server Configs",
		want: ServerConfigs{
			Http: HttpConfigs{
				Ip:   "",
				Port: 80,
			},
			Grpc: GrpcConfigs{
				Ip:   "",
				Port: 9090,
			},
			RateLimiter: RateLimiter{
				Burst: 5,
				Rate:  rate.Limit(2),
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetServerConfigs())
		})
	}
}

func TestGetRateLimiterConfigs(t *testing.T) {
	tests := []struct {
		name string
		want RateLimiter
	}{{
		name: "Default Rate Limiter Configs",
		want: RateLimiter{
			Burst: 5,
			Rate:  2,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetRateLimiterConfigs())
		})
	}
}

func TestGetGrpcConfigs(t *testing.T) {
	tests := []struct {
		name string
		want GrpcConfigs
	}{{
		name: "Default Grpc Configs",
		want: GrpcConfigs{
			Ip:   "",
			Port: 9090,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getGrpcConfigs())
		})
	}
}

func TestHttpConfigs_Format(t *testing.T) {
	type fields struct {
		Ip     string
		Port   int
		Config Config
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{{
		name: "Test Http Configs format",
		fields: fields{
			Ip:   "0.0.0.0",
			Port: 8080,
		},
		want: "0.0.0.0:8080",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpConfigs := &HttpConfigs{
				Ip:     tt.fields.Ip,
				Port:   tt.fields.Port,
				Config: tt.fields.Config,
			}
			assert.Equal(t, tt.want, httpConfigs.Format())
		})
	}
}

func TestServerConfigs_Format(t *testing.T) {
	type fields struct {
		Http        HttpConfigs
		Grpc        GrpcConfigs
		RateLimiter RateLimiter
		Config      Config
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{{
		name: "Test Http Configs format",
		fields: fields{
			Http: HttpConfigs{
				Ip:   "0.0.0.0",
				Port: 8080,
			},
			Grpc: GrpcConfigs{
				Ip:   "0.0.0.0",
				Port: 9090,
			},
			RateLimiter: RateLimiter{
				Rate:  10,
				Burst: 30,
			},
		},
		want: "Server running!!! http: 0.0.0.0:8080 |" +
			" gRPC: 0.0.0.0:9090 | rate limit: 10/30",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serverConfigs := &ServerConfigs{
				Http:        tt.fields.Http,
				Grpc:        tt.fields.Grpc,
				RateLimiter: tt.fields.RateLimiter,
				Config:      tt.fields.Config,
			}
			assert.Equal(t, tt.want, serverConfigs.Format())
		})
	}
}

func TestRateLimiter_Format(t *testing.T) {
	type fields struct {
		Burst  int
		Rate   rate.Limit
		Config Config
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{{
		name: "Rate Limit Configs Format",
		fields: fields{
			Burst: 10,
			Rate:  7,
		},
		want: "rate limit: 7/10",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rateLimiterConfigs := &RateLimiter{
				Burst:  tt.fields.Burst,
				Rate:   tt.fields.Rate,
				Config: tt.fields.Config,
			}
			assert.Equal(t, tt.want, rateLimiterConfigs.Format())
		})
	}
}

func TestGrpcConfigs_Format(t *testing.T) {
	type fields struct {
		Ip     string
		Port   int
		Config Config
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{{
		name: "Grpc Configs Format",
		fields: fields{
			Ip:   "1.1.1.1",
			Port: 53201,
		},
		want: "1.1.1.1:53201",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grpcConfigs := &GrpcConfigs{
				Ip:     tt.fields.Ip,
				Port:   tt.fields.Port,
				Config: tt.fields.Config,
			}
			assert.Equal(t, tt.want, grpcConfigs.Format())
		})
	}
}
