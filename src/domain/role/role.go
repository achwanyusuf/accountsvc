package role

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/achwanyusuf/carrent-accountsvc/src/model"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/psqlmodel"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/svcerr"
	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"

	"github.com/gin-gonic/gin"
	goredislib "github.com/redis/go-redis/v9"
)

type RoleDep struct {
	Log   logger.Logger
	DB    *sql.DB
	Redis *goredislib.Client
	Conf  Conf
}

type Conf struct {
	DefaultPageLimit    int           `mapstructure:"page_limit"`
	RedisExpirationTime time.Duration `mapstructure:"expiration_time"`
}

type RoleInterface interface {
	Insert(ctx *gin.Context, data *psqlmodel.Role) error
	GetSingleByParam(ctx *gin.Context, cacheControl string, param *model.GetRoleByParam) (psqlmodel.Role, error)
	Update(ctx *gin.Context, v *psqlmodel.Role) error
	Delete(ctx *gin.Context, v *psqlmodel.Role, id int64, isHardDelete bool) error
	GetByParam(ctx *gin.Context, cacheControl string, param *model.GetRolesByParam) (psqlmodel.RoleSlice, model.Pagination, error)
}

func New(conf Conf, log *logger.Logger, db *sql.DB, rds *goredislib.Client) RoleInterface {
	return &RoleDep{
		Log:   *log,
		DB:    db,
		Redis: rds,
		Conf:  conf,
	}
}

func (r *RoleDep) Insert(ctx *gin.Context, data *psqlmodel.Role) error {
	return r.insertPSQL(ctx, data)
}

func (r *RoleDep) GetSingleByParam(ctx *gin.Context, cacheControl string, param *model.GetRoleByParam) (psqlmodel.Role, error) {
	str, err := json.Marshal(param)
	if err != nil {
		return psqlmodel.Role{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error marshal param")
	}

	key := fmt.Sprintf(model.GetSingleByParamRoleKey, str)
	if cacheControl != model.MustRevalidate {
		res, err := r.getSingleByParamRedis(ctx, key)
		if err != nil {
			if err == goredislib.Nil {
				res, err := r.getSingleByParamPSQL(ctx, param)
				if err == nil {
					dataStr, err := json.Marshal(&res)
					if err != nil {
						return res, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
					}
					err = r.setRedis(ctx, key, string(dataStr))
					if err != nil {
						return res, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error set redis")
					}
				}
				return res, err
			}
			return res, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error marshal param")
		}
		return res, nil
	}

	res, err := r.getSingleByParamPSQL(ctx, param)
	if err == nil {
		dataStr, err := json.Marshal(&res)
		if err != nil {
			return res, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
		}
		err = r.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error set redis")
		}
	}
	return res, nil
}

func (r *RoleDep) Update(ctx *gin.Context, v *psqlmodel.Role) error {
	return r.updatePSQL(ctx, v)
}

func (r *RoleDep) Delete(ctx *gin.Context, v *psqlmodel.Role, id int64, isHardDelete bool) error {
	return r.deletePSQL(ctx, v, id, isHardDelete)
}
func (r *RoleDep) GetByParam(ctx *gin.Context, cacheControl string, param *model.GetRolesByParam) (psqlmodel.RoleSlice, model.Pagination, error) {
	var pg model.Pagination
	var res psqlmodel.RoleSlice

	str, err := json.Marshal(param)
	if err != nil {
		return psqlmodel.RoleSlice{}, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error marshal param")
	}

	key := fmt.Sprintf(model.GetByParamRoleKey, str)
	keyPg := fmt.Sprintf(model.GetByParamRolePgKey, str)
	if cacheControl != model.MustRevalidate {
		res, err1 := r.getByParamRedis(ctx, key)
		pg, err2 := r.getByParamPaginationRedis(ctx, keyPg)
		if err1 != nil || err2 != nil {
			if err1 == goredislib.Nil || err2 == goredislib.Nil {
				res, pg, err := r.getByParamPSQL(ctx, param)
				if err == nil {
					dataStr, err := json.Marshal(&res)
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
					}
					err = r.setRedis(ctx, key, string(dataStr))
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error set redis")
					}
					dataStr, err = json.Marshal(&pg)
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
					}
					err = r.setRedis(ctx, key, string(dataStr))
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error set redis")
					}
				}
				return res, pg, err
			}
			return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error marshal param")
		}
		return res, pg, nil
	}

	res, pg, err = r.getByParamPSQL(ctx, param)
	if err == nil {
		dataStr, err := json.Marshal(&res)
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
		}
		err = r.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error set redis")
		}
		dataStr, err = json.Marshal(&pg)
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
		}
		err = r.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error set redis")
		}
	}
	return res, pg, err
}
