package common

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP    ConfigHTTP         `yaml:"http"`
	GRPC    ConfigGRPC         `yaml:"grpc"`
	Logging ConfigLogging      `yaml:"logging"`
	Common  ConfigCommon       `yaml:"common"`
	Tasks   TaskDefinitionList `yaml:"tasks"`
}

type ConfigCommon struct {
	ConfigFile  *string       `arg:"--config-file,env:CONFIG_FILE"`
	GracePeriod time.Duration `arg:"env:GRACE_PERIOD" yaml:"grace-period" default:"15s"`
}

type ConfigLogging struct {
	Level       string `arg:"--log-level,env:LOG_LEVEL" yaml:"level" default:"debug"`
	Development bool   `arg:"--log-dev,env:LOG_DEV" yaml:"development" default:"false"`
}

type ConfigHTTP struct {
	ListenHost  string  `arg:"--http-listen-host,env:HTTP_LISTEN_HOST" yaml:"host" default:""`
	ListenPort  uint    `arg:"--http-listen-port,env:HTTP_LISTEN_PORT" yaml:"port" default:"3000"`
	Tls         bool    `arg:"--http-tls,env:HTTP_TLS" yaml:"tls" default:"false"`
	TlsKey      *string `arg:"--http-tls-key,env:HTTP_TLS_KEY" yaml:"tls-key"`
	TlsCert     *string `arg:"--http-tls-cert,env:HTTP_TLS_CERT" yaml:"tls-cert"`
	StaticPath  *string `arg:"--http-static,env:HTTP_STATIC" yaml:"static-path"`
	PathHealthz *string `arg:"--http-healthz-path,env:HTTP_HEALTHZ_PATH" yaml:"healthz-path"`
	PathReadyz  *string `arg:"--http-readyz-path,env:HTTP_READYZ_PATH" yaml:"readyz-path"`
}

type ConfigGRPC struct {
	ListenHost string  `arg:"--grpc-listen-host,env:GRPC_LISTEN_HOST" yaml:"host" default:""`
	ListenPort uint    `arg:"--grpc-listen-port,env:GRPC_LISTEN_PORT" yaml:"port" default:"3001"`
	Tls        bool    `arg:"--grpc-tls,env:GRPC_TLS" yaml:"tls" default:"false"`
	TlsKey     *string `arg:"--grpc-tls-key,env:GRPC_TLS_KEY" yaml:"tls-key"`
	TlsCert    *string `arg:"--grpc-tls-cert,env:GRPC_TLS_CERT" yaml:"tls-cert"`
	Reflection bool    `arg:"--grpc-reflection,env:GRPC_REFLECTION" yaml:"reflection"`
}

func (c *ConfigHTTP) Address() string {
	return net.JoinHostPort(c.ListenHost, fmt.Sprint(c.ListenPort))
}

func (c *ConfigGRPC) Address() string {
	return net.JoinHostPort(c.ListenHost, fmt.Sprint(c.ListenPort))
}

func (l TaskDefinitionList) Find(id string) *TaskDefinition {
	for _, t := range l {
		if t.Id == id {
			return &t
		}
	}

	return nil
}

func (c *ConfigLogging) AtomicLevel() zap.AtomicLevel {
	l, err := zap.ParseAtomicLevel(c.Level)
	if err != nil {
		panic(fmt.Errorf("log level not parsable - %w", err))
	}

	return l
}

func (c *Config) MustParse() {
	arg.MustParse(&c.Logging)
	arg.MustParse(&c.Common)
	arg.MustParse(&c.GRPC)
	arg.MustParse(&c.HTTP)

	if c.Common.ConfigFile != nil {
		cf := *c.Common.ConfigFile
		if _, err := os.Stat(cf); err == nil {
			if cd, err := os.ReadFile(cf); err == nil {
				if err := yaml.Unmarshal(cd, c); err != nil {
					panic(err)
				}
			}
		}
	}
}
