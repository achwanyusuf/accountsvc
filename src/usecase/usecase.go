package usecase

import (
	"github.com/achwanyusuf/carrent-accountsvc/src/domain"
	"github.com/achwanyusuf/carrent-accountsvc/src/usecase/account"
	"github.com/achwanyusuf/carrent-accountsvc/src/usecase/accountrole"
	"github.com/achwanyusuf/carrent-accountsvc/src/usecase/role"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"
)

type UsecaseDep struct {
	Conf   Config
	Log    *logger.Logger
	Domain *domain.DomainInterface
}

type Config struct {
	Account     account.Conf     `mapstructure:"account"`
	Role        role.Conf        `mapstructure:"role"`
	AccountRole accountrole.Conf `mapstructure:"account_role"`
}

type UsecaseInterface struct {
	Account     account.AccountInterface
	Role        role.RoleInterface
	AccountRole accountrole.AccountRoleInterface
}

func New(u *UsecaseDep) *UsecaseInterface {
	return &UsecaseInterface{
		account.New(u.Conf.Account, u.Log, u.Domain.Account, u.Domain.Role, u.Domain.AccountRole),
		role.New(u.Conf.Role, u.Log, u.Domain.Role),
		accountrole.New(u.Conf.AccountRole, u.Log, u.Domain.AccountRole),
	}
}
