package environment

import (
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Server    Server
	Database  Database
	SecretKey string `envconfig:"SECRET_KEY"`
	// Wechat    Wechat
}

type Server struct {
	Port string `envconfig:"PORT" validate:"required"`
}

type Database struct {
	Sqlite   Sqlite
	Postgres Postgres
}

// type Wechat struct {
// 	CORPID     string `envconfig:"WECHAT_CORP_ID"`
// 	CORPSECRET string `envconfig:"WECHAT_CORP_SECRET"`
// 	AGENTID    string `envconfig:"WECHAT_AGENT_ID"`
// }

func Load() (*Configuration, error) {
	configuration := &Configuration{}
	_ = readEnv(configuration)
	return configuration, configuration.Validate()
}

func readEnv(cfg interface{}) error {
	return envconfig.Process("", cfg)
}

func (c Configuration) Validate() error {
	return validator.New().Struct(c.Server)
}
