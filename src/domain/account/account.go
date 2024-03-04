package account

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

type AccountDep struct {
	Log   logger.Logger
	DB    *sql.DB
	Redis *goredislib.Client
	Conf  Conf
}

type Conf struct {
	DefaultPageLimit    int           `mapstructure:"page_limit"`
	RedisExpirationTime time.Duration `mapstructure:"expiration_time"`
}

type AccountInterface interface {
	Insert(ctx *gin.Context, data *psqlmodel.Account) error
	GetSingleByParam(ctx *gin.Context, cacheControl string, param *model.GetAccountByParam) (psqlmodel.Account, error)
	Update(ctx *gin.Context, account *psqlmodel.Account) error
	Delete(ctx *gin.Context, account *psqlmodel.Account, id int64, isHardDelete bool) error
	GetByParam(ctx *gin.Context, cacheControl string, param *model.GetAccountsByParam) (psqlmodel.AccountSlice, model.Pagination, error)
}

func New(conf Conf, log *logger.Logger, db *sql.DB, rds *goredislib.Client) AccountInterface {
	return &AccountDep{
		Log:   *log,
		DB:    db,
		Redis: rds,
		Conf:  conf,
	}
}

func (a *AccountDep) Insert(ctx *gin.Context, data *psqlmodel.Account) error {
	return a.insertPSQL(ctx, data)
}

func (a *AccountDep) GetSingleByParam(ctx *gin.Context, cacheControl string, param *model.GetAccountByParam) (psqlmodel.Account, error) {
	str, err := json.Marshal(param)
	if err != nil {
		return psqlmodel.Account{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error marshal param")
	}

	key := fmt.Sprintf(model.GetSingleByParamAccountKey, str)
	if cacheControl != model.MustRevalidate {
		res, err := a.getSingleByParamRedis(ctx, key)
		if err != nil {
			if err == goredislib.Nil {
				res, err := a.getSingleByParamPSQL(ctx, param)
				if err == nil {
					dataStr, err := json.Marshal(&res)
					if err != nil {
						return res, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
					}
					err = a.setRedis(ctx, key, string(dataStr))
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

	res, err := a.getSingleByParamPSQL(ctx, param)
	if err == nil {
		dataStr, err := json.Marshal(&res)
		if err != nil {
			return res, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
		}
		err = a.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error set redis")
		}
	}
	return res, nil
}

func (a *AccountDep) Update(ctx *gin.Context, account *psqlmodel.Account) error {
	return a.updatePSQL(ctx, account)
}

func (a *AccountDep) Delete(ctx *gin.Context, account *psqlmodel.Account, id int64, isHardDelete bool) error {
	return a.deletePSQL(ctx, account, id, isHardDelete)
}
func (a *AccountDep) GetByParam(ctx *gin.Context, cacheControl string, param *model.GetAccountsByParam) (psqlmodel.AccountSlice, model.Pagination, error) {
	var pg model.Pagination
	var res psqlmodel.AccountSlice

	str, err := json.Marshal(param)
	if err != nil {
		return psqlmodel.AccountSlice{}, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error marshal param")
	}

	key := fmt.Sprintf(model.GetByParamAccountKey, str)
	keyPg := fmt.Sprintf(model.GetByParamAccountPgKey, str)
	if cacheControl != model.MustRevalidate {
		res, err1 := a.getByParamRedis(ctx, key)
		pg, err2 := a.getByParamPaginationRedis(ctx, keyPg)
		if err1 != nil || err2 != nil {
			if err1 == goredislib.Nil || err2 == goredislib.Nil {
				res, pg, err := a.getByParamPSQL(ctx, param)
				if err == nil {
					dataStr, err := json.Marshal(&res)
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
					}
					err = a.setRedis(ctx, key, string(dataStr))
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error set redis")
					}
					dataStr, err = json.Marshal(&pg)
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
					}
					err = a.setRedis(ctx, key, string(dataStr))
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

	res, pg, err = a.getByParamPSQL(ctx, param)
	if err == nil {
		dataStr, err := json.Marshal(&res)
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
		}
		err = a.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error set redis")
		}
		dataStr, err = json.Marshal(&pg)
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get psql")
		}
		err = a.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error set redis")
		}
	}
	return res, pg, err
}
