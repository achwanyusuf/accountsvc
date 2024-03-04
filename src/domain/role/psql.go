package role

import (
	"database/sql"

	"github.com/achwanyusuf/carrent-accountsvc/src/model"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/psqlmodel"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/svcerr"
	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (r *RoleDep) insertPSQL(ctx *gin.Context, data *psqlmodel.Role) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return errormsg.WrapErr(svcerr.AccountSVCPSQLErrorTransaction, err, "error begin transaction")
	}

	err = data.Insert(ctx, tx, boil.Infer())
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			r.Log.Warn(ctx, errormsg.WrapErr(svcerr.AccountSVCPSQLErrorRollback, err, "error rollback"))
		}
		return errormsg.WrapErr(svcerr.AccountSVCPSQLErrorInsert, err, "error insert")
	}
	err = tx.Commit()
	if err != nil {
		return errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error commit")
	}
	return nil
}

func (r *RoleDep) getSingleByParamPSQL(ctx *gin.Context, param *model.GetRoleByParam) (psqlmodel.Role, error) {
	var res psqlmodel.Role
	qr := param.GetQuery()
	account, err := psqlmodel.Roles(qr...).One(ctx, r.DB)
	if err == sql.ErrNoRows {
		return res, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get accounts")
	}

	if err != nil {
		return res, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get accounts")
	}

	return *account, nil
}

func (r *RoleDep) updatePSQL(ctx *gin.Context, account *psqlmodel.Role) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return errormsg.WrapErr(svcerr.AccountSVCPSQLErrorTransaction, err, "error begin transaction")
	}

	_, err = account.Update(ctx, tx, boil.Infer())
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			r.Log.Warn(ctx, errormsg.WrapErr(svcerr.AccountSVCPSQLErrorRollback, err, "error rollback"))
		}
		return errormsg.WrapErr(svcerr.AccountSVCPSQLErrorUpdate, err, "error update")
	}
	err = tx.Commit()
	if err != nil {
		return errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error commit")
	}
	return nil
}

func (r *RoleDep) deletePSQL(ctx *gin.Context, account *psqlmodel.Role, id int64, isHardDelete bool) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return errormsg.WrapErr(svcerr.AccountSVCPSQLErrorTransaction, err, "error begin transaction")
	}

	_, err = account.Delete(ctx, tx, isHardDelete)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			r.Log.Warn(ctx, errormsg.WrapErr(svcerr.AccountSVCPSQLErrorRollback, err, "error rollback"))
		}
		return errormsg.WrapErr(svcerr.AccountSVCPSQLErrorUpdate, err, "error delete")
	}

	if !isHardDelete {
		account.DeletedBy = null.NewInt(int(id), true)
		_, err = account.Update(ctx, tx, boil.Infer())
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				r.Log.Warn(ctx, errormsg.WrapErr(svcerr.AccountSVCPSQLErrorRollback, err, "error rollback"))
			}
			return errormsg.WrapErr(svcerr.AccountSVCPSQLErrorUpdate, err, "error update")
		}
	}
	err = tx.Commit()
	if err != nil {
		return errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error commit")
	}
	return nil
}

func (r *RoleDep) getByParamPSQL(ctx *gin.Context, param *model.GetRolesByParam) (psqlmodel.RoleSlice, model.Pagination, error) {
	var totalPages int64 = 1
	if param.Limit == 0 {
		param.Limit = int64(r.Conf.DefaultPageLimit)
	}

	if param.Page == 0 {
		param.Page = 1
	}

	qr := param.GetQuery()
	count, err := psqlmodel.Roles(qr...).Count(ctx, r.DB)
	if err != nil {
		return psqlmodel.RoleSlice{}, model.Pagination{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error count data")
	}
	qr = append(qr, qm.Offset(int((param.Page-1)*param.Limit)))
	qr = append(qr, qm.Limit(int(param.Limit)))
	accounts, err := psqlmodel.Roles(qr...).All(ctx, r.DB)
	if err == sql.ErrNoRows {
		return accounts, model.Pagination{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get accounts")
	}
	if err != nil {
		return accounts, model.Pagination{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get accounts")
	}
	if count > 0 {
		totalPages = (count / param.Limit) + 1
	}
	return accounts, model.Pagination{
		CurrentPage:     param.Page,
		CurrentElements: int64(len(accounts)),
		TotalElements:   count,
		TotalPages:      totalPages,
		SortBy:          param.OrderBy.String,
	}, nil
}
