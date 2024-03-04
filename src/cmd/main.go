package main

import (
	"context"
	"flag"
	"os"

	"github.com/achwanyusuf/carrent-accountsvc/conf"
	"github.com/achwanyusuf/carrent-accountsvc/docs"
	"github.com/achwanyusuf/carrent-accountsvc/src/domain"
	"github.com/achwanyusuf/carrent-accountsvc/src/handler/rest"
	"github.com/achwanyusuf/carrent-accountsvc/src/usecase"
	"github.com/achwanyusuf/carrent-lib/pkg/httpserver"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"
	"github.com/achwanyusuf/carrent-lib/pkg/migration"
	"github.com/achwanyusuf/carrent-lib/pkg/psql"
	"github.com/achwanyusuf/carrent-lib/pkg/redis"
)

// @contact.name   CarRent Support
// @contact.url		https://www.carrent.com/support
// @contact.email 	support@carrent.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl %s
var (
	staticConfPath, Namespace, BuildTime, Version string
	migrateup, migratedown                        bool
	OAuth2PasswordTokenUrl                        string
)

func main() {
	flag.StringVar(&staticConfPath, "staticConfPath", "./conf/conf.yaml", "config path")
	flag.BoolVar(&migrateup, "migrateup", false, "run migration up")
	flag.BoolVar(&migratedown, "migratedown", false, "run migration up")
	flag.Parse()
	cfg, err := conf.New(staticConfPath)
	if err != nil {
		panic(err)
	}

	if Version == "" {
		Version = "v1.0.0"
	}
	// setup logger
	log := logger.New(&logger.Config{
		IsFile: false,
		Level:  logger.LevelDebug,
		CustomFields: map[string]interface{}{
			"namespace":  Namespace,
			"version":    Version,
			"build_time": BuildTime,
			"pid":        os.Getpid(),
		},
	})

	// setup db connection
	psql := psql.PsqlConnect(cfg.App.PSQL)

	if migrateup {
		err = migration.Migrate(psql, cfg.App.PSQL.MigrationPath, true)
		if err != nil {
			log.Error(context.Background(), err)
		}
		return
	}

	if migratedown {
		err = migration.Migrate(psql, cfg.App.PSQL.MigrationPath, false)
		if err != nil {
			log.Error(context.Background(), err)
		}
		return
	}

	// setup redis connection
	redis := redis.RedisConnect(cfg.App.Redis)

	// init domain
	dom := domain.New(&domain.DomainDep{
		Conf:  cfg.Domain,
		Log:   &log,
		DB:    psql,
		Redis: redis,
	})

	// init usecase
	uc := usecase.New(&usecase.UsecaseDep{
		Conf:   cfg.Usecase,
		Log:    &log,
		Domain: dom,
	})
	cfg.App.Swagger.Title = Namespace
	cfg.App.Swagger.Version = Version
	// setup http server
	http := httpserver.HTTPSetting{
		Env:     cfg.App.Env,
		Conf:    cfg.App.HTTPServer,
		Swagger: cfg.App.Swagger,
		Log:     log,
	}

	gin := http.NewHTTPServer()
	http.SetSwaggo(docs.SwaggerInfo)

	// init http router
	restCfg := rest.RestDep{
		Conf:    cfg.Rest,
		Log:     &log,
		Usecase: uc,
		Gin:     gin,
	}
	handler := rest.New(&restCfg)

	restCfg.Serve(handler)

	http.Run()
}
