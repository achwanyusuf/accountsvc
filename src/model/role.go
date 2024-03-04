package model

import (
	"strings"

	"github.com/achwanyusuf/carrent-accountsvc/src/model/psqlmodel"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/svcerr"
	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	GetSingleByParamRoleKey string = "gspRole:%s"
	GetByParamRoleKey       string = "gpRole:%s"
	GetByParamRolePgKey     string = "gppgRole:%s"
	SuperAdminScope         string = "sup"
	StoreScope              string = "sto"
	CustomerScope           string = "cus"
)

type GetRoleByParam struct {
	ID    null.Int64  `schema:"id" json:"id"`
	Scope null.String `schema:"scope" json:"scope"`
	Cid   null.String `schema:"cid" json:"cid"`
}

func (g *GetRoleByParam) GetQuery() []qm.QueryMod {
	var res []qm.QueryMod
	if g.ID.Valid {
		res = append(res, qm.Where("id=?", g.ID.Int64))
	}

	if g.Scope.Valid {
		res = append(res, qm.Where("scope=?", g.Scope.String))
	}

	if g.Cid.Valid {
		res = append(res, qm.Where("cid=?", g.Cid.String))
	}
	return res
}

type GetRolesByParam struct {
	GetRoleByParam
	OrderBy null.String `schema:"order_by" json:"order_by"`
	Limit   int64       `schema:"limit" json:"limit"`
	Page    int64       `schema:"page" json:"page"`
}

func (g *GetRolesByParam) GetQuery() []qm.QueryMod {
	var res []qm.QueryMod
	if g.ID.Valid {
		res = append(res, qm.Where("id=?", g.ID.Int64))
	}

	if g.Scope.Valid {
		res = append(res, qm.Where("scope like ?", g.Scope.String+"%"))
	}

	if g.Cid.Valid {
		res = append(res, qm.Where("cid like ?", g.Cid.String+"%"))
	}

	if g.OrderBy.Valid {
		order := strings.Split(g.OrderBy.String, ",")
		for _, o := range order {
			res = append(res, qm.OrderBy(o))
		}
	}

	return res
}

type CreateRole struct {
	Scope     string `json:"scope"`
	Cid       string `json:"client_id"`
	Sec       string `json:"client_secret"`
	CreatedBy int64  `json:"-"`
}

func (v *CreateRole) Validate() error {
	if v.Scope == "" {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidScope, nil, "invalid scope")
	}

	if v.Cid == "" {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidClientIDClientSecret, nil, "invalid scope")
	}

	if v.Cid == "" {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidClientIDClientSecret, nil, "invalid scope")
	}
	return nil
}

type UpdateRole struct {
	Scope     null.String `json:"scope"`
	Cid       null.String `json:"client_id"`
	Sec       null.String `json:"client_secret"`
	UpdatedBy int64       `json:"-"`
}

func (v *UpdateRole) FillEntity(role *psqlmodel.Role) {
	if v.Scope.Valid {
		role.Scope = v.Scope.String
	}

	if v.Cid.Valid {
		role.Cid = v.Cid.String
	}

	if v.Sec.Valid {
		role.Sec = v.Sec.String
	}
}

type Role struct {
	ID    int64  `json:"id"`
	Scope string `json:"scope"`
	Cid   string `json:"client_id"`
	BaseInformation
}

func TransformPSQLSingleRole(role *psqlmodel.Role) Role {
	creationInfo := BaseInformation{
		CreatedBy: int64(role.CreatedBy),
		CreatedAt: role.CreatedAt,
		UpdatedBy: int64(role.UpdatedBy),
		UpdatedAt: role.UpdatedAt,
		DeletedBy: int64(role.DeletedBy.Int),
		DeletedAt: role.DeletedAt.Time,
	}

	return Role{
		ID:              int64(role.ID),
		Scope:           role.Scope,
		Cid:             role.Cid,
		BaseInformation: creationInfo,
	}
}

func TransformPSQLRole(role *psqlmodel.RoleSlice) []Role {
	var res []Role
	for _, v := range *role {
		creationInfo := BaseInformation{
			CreatedBy: int64(v.CreatedBy),
			CreatedAt: v.CreatedAt,
			UpdatedBy: int64(v.UpdatedBy),
			UpdatedAt: v.UpdatedAt,
			DeletedBy: int64(v.DeletedBy.Int),
			DeletedAt: v.DeletedAt.Time,
		}

		res = append(res, Role{
			ID:              int64(v.ID),
			Scope:           v.Scope,
			Cid:             v.Cid,
			BaseInformation: creationInfo,
		})
	}

	return res
}
