package accountrole

import (
	"github.com/achwanyusuf/carrent-accountsvc/src/domain/accountrole"
	"github.com/achwanyusuf/carrent-accountsvc/src/model"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/psqlmodel"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/svcerr"
	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
)

type AccountRoleDep struct {
	log         logger.Logger
	conf        Conf
	accountRole accountrole.AccountRoleInterface
}

type Conf struct{}

type AccountRoleInterface interface {
	Create(ctx *gin.Context, v model.CreateAccountRole) (model.AccountRole, error)
	GetByParam(ctx *gin.Context, cacheControl string, v model.GetAccountRolesByParam) ([]model.AccountRole, model.Pagination, error)
	GetByID(ctx *gin.Context, cacheControl string, id int64) (model.AccountRole, error)
	DeleteByID(ctx *gin.Context, id int64, isHardDelete bool, vid int64) error
}

func New(conf Conf, logger *logger.Logger, accountRole accountrole.AccountRoleInterface) AccountRoleInterface {
	return &AccountRoleDep{
		conf:        conf,
		log:         *logger,
		accountRole: accountRole,
	}
}

func (a *AccountRoleDep) Create(ctx *gin.Context, v model.CreateAccountRole) (model.AccountRole, error) {
	var result model.AccountRole
	err := v.Validate()
	if err != nil {
		return result, err
	}

	role := &psqlmodel.AccountRole{
		AccountID: int(v.AccountID),
		RoleID:    int(v.RoleID),
		CreatedBy: int(v.CreatedBy),
		UpdatedBy: int(v.CreatedBy),
	}

	err = a.accountRole.Insert(ctx, role)
	if err != nil {
		return result, err
	}

	return model.TransformPSQLSingleAccountRole(role), nil
}

func (a *AccountRoleDep) GetByParam(ctx *gin.Context, cacheControl string, v model.GetAccountRolesByParam) ([]model.AccountRole, model.Pagination, error) {
	accountRoleSlice, pagination, err := a.accountRole.GetByParam(ctx, cacheControl, &v)
	if err != nil {
		return []model.AccountRole{}, model.Pagination{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get by param")
	}
	return model.TransformPSQLAccountRole(&accountRoleSlice), pagination, nil
}

func (a *AccountRoleDep) GetByID(ctx *gin.Context, cacheControl string, id int64) (model.AccountRole, error) {
	accountRole, err := a.accountRole.GetSingleByParam(ctx, cacheControl, &model.GetAccountRoleByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.AccountRole{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "data not found")
	}
	return model.TransformPSQLSingleAccountRole(&accountRole), nil
}

func (a *AccountRoleDep) DeleteByID(ctx *gin.Context, id int64, isHardDelete bool, vid int64) error {
	accountRole, err := a.accountRole.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountRoleByParam{
		ID: null.NewInt64(vid, true),
	})
	if err != nil {
		return err
	}
	return a.accountRole.Delete(ctx, &accountRole, id, isHardDelete)
}
