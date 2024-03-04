package rest

import (
	"github.com/achwanyusuf/carrent-accountsvc/src/handler/rest/account"
	"github.com/achwanyusuf/carrent-accountsvc/src/handler/rest/accountrole"
	"github.com/achwanyusuf/carrent-accountsvc/src/handler/rest/role"
	"github.com/achwanyusuf/carrent-accountsvc/src/model"
	"github.com/achwanyusuf/carrent-accountsvc/src/usecase"
	"github.com/achwanyusuf/carrent-lib/pkg/httpserver"
	"github.com/achwanyusuf/carrent-lib/pkg/jwt"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"
	"github.com/gin-gonic/gin"
)

type RestDep struct {
	Conf    Config
	Log     *logger.Logger
	Usecase *usecase.UsecaseInterface
	Gin     *gin.Engine
}

type Config struct {
	Account     account.Conf     `mapstructure:"account"`
	Role        role.Conf        `mapstructure:"role"`
	AccountRole accountrole.Conf `mapstructure:"account_role"`
}

type RestInterface struct {
	Account     account.AccountInterface
	Role        role.RoleInterface
	AccountRole accountrole.AccountRoleInterface
}

func New(r *RestDep) *RestInterface {
	return &RestInterface{
		account.New(r.Conf.Account, r.Log, r.Usecase.Account),
		role.New(r.Conf.Role, r.Log, r.Usecase.Role),
		accountrole.New(r.Conf.AccountRole, r.Log, r.Usecase.AccountRole),
	}
}

func (r *RestDep) Serve(handler *RestInterface) {
	api := r.Gin.Group("/api")
	api.POST("/oauth2", handler.Account.Oauth2)
	api.POST("/register", handler.Account.Register)

	api.Use(jwt.JWT(*r.Log, []byte(r.Conf.Account.TokenSecret)))
	{
		api.GET("/me", handler.Account.CurrentAccount)
		api.PUT("/me", handler.Account.UpdateCurrentAccount)
		api.PUT("/me/password", handler.Account.UpdatePasswordAccount)
		api.POST("/account", handler.Account.Create)
		api.GET("/account", httpserver.ValidateScope(*r.Log, []string{model.SuperAdminScope}), handler.Account.Read)
		api.GET("/account/:id", handler.Account.GetByID)
		api.PUT("/account/:id", handler.Account.UpdateByID)
		api.DELETE("/account/:id", handler.Account.DeleteByID)

		api.POST("/role", httpserver.ValidateScope(*r.Log, []string{model.SuperAdminScope}), handler.Role.Create)
		api.GET("/role", httpserver.ValidateScope(*r.Log, []string{model.SuperAdminScope}), handler.Role.Read)
		api.GET("/role/:id", httpserver.ValidateScope(*r.Log, []string{model.SuperAdminScope}), handler.Role.GetByID)
		api.PUT("/role/:id", httpserver.ValidateScope(*r.Log, []string{model.SuperAdminScope}), handler.Role.UpdateByID)
		api.DELETE("/role/:id", httpserver.ValidateScope(*r.Log, []string{model.SuperAdminScope}), handler.Role.DeleteByID)

		api.POST("/account-role", handler.AccountRole.Create)
		api.GET("/account-role", httpserver.ValidateScope(*r.Log, []string{model.SuperAdminScope}), handler.AccountRole.Read)
		api.GET("/account-role/:id", httpserver.ValidateScope(*r.Log, []string{model.SuperAdminScope}), handler.AccountRole.GetByID)
		api.DELETE("/account-role/:id", httpserver.ValidateScope(*r.Log, []string{model.SuperAdminScope}), handler.AccountRole.DeleteByID)
	}
}
