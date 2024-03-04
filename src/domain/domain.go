package domain

import (
	"database/sql"

	"github.com/achwanyusuf/carrent-accountsvc/src/domain/account"
	"github.com/achwanyusuf/carrent-accountsvc/src/domain/accountrole"
	"github.com/achwanyusuf/carrent-accountsvc/src/domain/role"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"
	goredislib "github.com/redis/go-redis/v9"
)

type DomainDep struct {
	Conf  Config
	Log   *logger.Logger
	DB    *sql.DB
	Redis *goredislib.Client
}

type Config struct {
	Account     account.Conf     `mapstructure:"account"`
	Role        role.Conf        `mapstructure:"role"`
	AccountRole accountrole.Conf `mapstructure:"account_role"`
}

type DomainInterface struct {
	Account     account.AccountInterface
	Role        role.RoleInterface
	AccountRole accountrole.AccountRoleInterface
}

func New(d *DomainDep) *DomainInterface {
	return &DomainInterface{
		account.New(d.Conf.Account, d.Log, d.DB, d.Redis),
		role.New(d.Conf.Role, d.Log, d.DB, d.Redis),
		accountrole.New(d.Conf.AccountRole, d.Log, d.DB, d.Redis),
	}
}
