package conf

import (
	"github.com/achwanyusuf/carrent-accountsvc/src/domain"
	"github.com/achwanyusuf/carrent-accountsvc/src/handler/rest"
	"github.com/achwanyusuf/carrent-accountsvc/src/usecase"
	"github.com/achwanyusuf/carrent-lib/pkg/httpserver"
	"github.com/achwanyusuf/carrent-lib/pkg/psql"
	"github.com/achwanyusuf/carrent-lib/pkg/redis"
)

type Config struct {
	App     App            `mapstructure:"app"`
	Rest    rest.Config    `mapstructure:"rest"`
	Usecase usecase.Config `mapstructure:"usecase"`
	Domain  domain.Config  `mapstructure:"domain"`
}

type App struct {
	Env        string                `mapstructure:"env"`
	HTTPServer httpserver.HTTPServer `mapstructure:"http_server"`
	Swagger    httpserver.Swagger    `mapstructure:"swagger"`
	PSQL       psql.PSQL             `mapstructure:"psql"`
	Redis      redis.Redis           `mapstructure:"redis"`
}
