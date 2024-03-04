package role

import (
	"github.com/achwanyusuf/carrent-accountsvc/src/domain/role"
	"github.com/achwanyusuf/carrent-accountsvc/src/model"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/psqlmodel"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/svcerr"
	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
	"github.com/achwanyusuf/carrent-lib/pkg/hash"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
)

type RoleDep struct {
	log  logger.Logger
	conf Conf
	role role.RoleInterface
}

type Conf struct {
	SecretKey string `mapstructure:"secret_key"`
}

type RoleInterface interface {
	Create(ctx *gin.Context, v model.CreateRole) (model.Role, error)
	GetByParam(ctx *gin.Context, cacheControl string, v model.GetRolesByParam) ([]model.Role, model.Pagination, error)
	GetByID(ctx *gin.Context, cacheControl string, id int64) (model.Role, error)
	UpdateByID(ctx *gin.Context, id int64, v model.UpdateRole) (model.Role, error)
	DeleteByID(ctx *gin.Context, id int64, isHardDelete bool, vid int64) error
}

func New(conf Conf, logger *logger.Logger, role role.RoleInterface) RoleInterface {
	return &RoleDep{
		conf: conf,
		log:  *logger,
		role: role,
	}
}

func (r *RoleDep) Create(ctx *gin.Context, v model.CreateRole) (model.Role, error) {
	var result model.Role
	err := v.Validate()
	if err != nil {
		return result, err
	}

	secret, err := hash.EncAES(v.Sec, r.conf.SecretKey)
	if err != nil {
		return result, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error hash secret key")
	}

	role := &psqlmodel.Role{
		Scope:     v.Scope,
		Cid:       v.Cid,
		Sec:       secret,
		CreatedBy: int(v.CreatedBy),
		UpdatedBy: int(v.CreatedBy),
	}

	err = r.role.Insert(ctx, role)
	if err != nil {
		return result, err
	}

	return model.TransformPSQLSingleRole(role), nil
}

func (r *RoleDep) GetByParam(ctx *gin.Context, cacheControl string, v model.GetRolesByParam) ([]model.Role, model.Pagination, error) {
	roleSlice, pagination, err := r.role.GetByParam(ctx, cacheControl, &v)
	if err != nil {
		return []model.Role{}, model.Pagination{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get by param")
	}
	return model.TransformPSQLRole(&roleSlice), pagination, nil
}

func (r *RoleDep) GetByID(ctx *gin.Context, cacheControl string, id int64) (model.Role, error) {
	role, err := r.role.GetSingleByParam(ctx, cacheControl, &model.GetRoleByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.Role{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "data not found")
	}
	return model.TransformPSQLSingleRole(&role), nil
}

func (r *RoleDep) UpdateByID(ctx *gin.Context, id int64, v model.UpdateRole) (model.Role, error) {
	role, err := r.role.GetSingleByParam(ctx, model.MustRevalidate, &model.GetRoleByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.Role{}, err
	}

	if !v.Scope.Valid && !v.Cid.Valid && !v.Sec.Valid {
		return model.TransformPSQLSingleRole(&role), nil
	}

	if v.Scope.Valid {
		role.Scope = v.Scope.String
	}

	if v.Sec.Valid {
		v.Sec.String, err = hash.EncAES(v.Sec.String, r.conf.SecretKey)
		if err != nil {
			return model.Role{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error hash secret key")
		}
	}

	v.FillEntity(&role)
	role.UpdatedBy = int(v.UpdatedBy)

	err = r.role.Update(ctx, &role)
	if err != nil {
		return model.Role{}, err
	}

	return model.TransformPSQLSingleRole(&role), nil
}

func (r *RoleDep) DeleteByID(ctx *gin.Context, id int64, isHardDelete bool, vid int64) error {
	role, err := r.role.GetSingleByParam(ctx, model.MustRevalidate, &model.GetRoleByParam{
		ID: null.NewInt64(vid, true),
	})
	if err != nil {
		return err
	}
	return r.role.Delete(ctx, &role, id, isHardDelete)
}
