package model

import (
	"regexp"
	"strings"
	"time"

	"github.com/achwanyusuf/carrent-accountsvc/src/model/psqlmodel"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/svcerr"
	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	GetSingleByParamAccountKey string = "gspAccount:%s"
	GetByParamAccountKey       string = "gpAccount:%s"
	GetByParamAccountPgKey     string = "gppgAccount:%s"
	MustRevalidate             string = "must-revalidate"
	RegExpEmail                string = `^[a-zA-Z0-9._+\-]+@[a-zA-Z0-9]+[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,10}$`
)

type Account struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	BaseInformation
}

type Login struct {
	Email        string `json:"username"`
	Password     string `json:"password"`
	ClientID     string `json:"-"`
	ClientSecret string `json:"-"`
}

type Auth struct {
	AccessToken string     `json:"access_token,omitempty"`
	TokenType   string     `json:"token_type,omitempty"`
	Exp         *time.Time `json:"exp,omitempty"`
	Scope       string     `json:"scope,omitempty"`
}

func (l *Login) Validate() error {
	if l.Email == "" {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidEmptyEmail, nil, "invalid empty name")
	}

	rg := regexp.MustCompile(RegExpEmail)
	if !rg.MatchString(l.Email) {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidEmailFormat, nil, "invalid email format")
	}

	if l.Password == "" {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidEmptyPassword, nil, "invalid empty password")
	}

	if len(l.Password) < 5 {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidMinimumPassword, nil, "invalid minimum password")
	}

	if len(l.Password) > 8 {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidMaximumPassword, nil, "invalid maximum password")
	}
	return nil
}

type Register struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	CreatedBy       int64  `json:"-"`
}

func (r *Register) Validate() error {
	if r.Name == "" {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidEmptyName, nil, "invalid empty name")
	}

	if r.Password == "" {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidEmptyPassword, nil, "invalid empty password")
	}

	if len(r.Password) < 5 {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidMinimumPassword, nil, "invalid minimum password")
	}

	if len(r.Password) > 8 {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidMaximumPassword, nil, "invalid maximum password")
	}
	if r.Password != r.ConfirmPassword {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidPasswordConfirmation, nil, "invalid password confirmation")
	}
	return nil
}

type GetAccountByParam struct {
	ID    null.Int64  `schema:"id" json:"id"`
	Email null.String `schema:"email" json:"email"`
	Name  null.String `schema:"name" json:"name"`
}

func (g *GetAccountByParam) GetQuery() []qm.QueryMod {
	var res []qm.QueryMod
	if g.ID.Valid {
		res = append(res, qm.Where("id=?", g.ID.Int64))
	}

	if g.Email.Valid {
		res = append(res, qm.Where("email=?", g.Email.String))
	}

	if g.Name.Valid {
		res = append(res, qm.Where("name=?", g.Name.String))
	}
	return res
}

type GetAccountsByParam struct {
	GetAccountByParam
	OrderBy null.String `schema:"order_by" json:"order_by"`
	Limit   int64       `schema:"limit" json:"limit"`
	Page    int64       `schema:"page" json:"page"`
}

func (g *GetAccountsByParam) GetQuery() []qm.QueryMod {
	var res []qm.QueryMod
	if g.ID.Valid {
		res = append(res, qm.Where("id=?", g.ID.Int64))
	}

	if g.Email.Valid {
		res = append(res, qm.Where("email like ?", g.Email.String+"%"))
	}

	if g.Name.Valid {
		res = append(res, qm.Where("name like ?", g.Name.String+"%"))
	}

	if g.OrderBy.Valid {
		order := strings.Split(g.OrderBy.String, ",")
		for _, o := range order {
			res = append(res, qm.OrderBy(o))
		}
	}

	return res
}

type GetAccounts struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdateAccountData struct {
	Name     string `json:"name"`
	UpdateBy int64  `json:"-"`
}

func (u *UpdateAccountData) Validate() error {
	if u.Name == "" {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidEmptyName, nil, "invalid empty name")
	}
	return nil
}

type UpdatePasswordData struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	UpdateBy        int64  `json:"-"`
}

func (u *UpdatePasswordData) IsValid() error {
	if u.Password == "" {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidEmptyPassword, nil, "invalid empty password")
	}

	if len(u.Password) < 5 {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidMinimumPassword, nil, "invalid minimum password")
	}

	if len(u.Password) > 8 {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidMaximumPassword, nil, "invalid maximum password")
	}
	if u.Password != u.ConfirmPassword {
		return errormsg.WrapErr(svcerr.AccountSVCInvalidPasswordConfirmation, nil, "invalid password confirmation")
	}
	return nil
}

func TransformPSQLSingleAccount(account *psqlmodel.Account) Account {
	creationInfo := BaseInformation{
		CreatedBy: int64(account.CreatedBy),
		CreatedAt: account.CreatedAt,
		UpdatedBy: int64(account.UpdatedBy),
		UpdatedAt: account.UpdatedAt,
		DeletedBy: int64(account.DeletedBy.Int),
		DeletedAt: account.DeletedAt.Time,
	}

	return Account{
		ID:              int64(account.ID),
		Name:            account.Name,
		Email:           account.Email,
		BaseInformation: creationInfo,
	}
}

func TransformPSQLAccount(account *psqlmodel.AccountSlice) []Account {
	var res []Account
	for _, v := range *account {
		creationInfo := BaseInformation{
			CreatedBy: int64(v.CreatedBy),
			CreatedAt: v.CreatedAt,
			UpdatedBy: int64(v.UpdatedBy),
			UpdatedAt: v.UpdatedAt,
			DeletedBy: int64(v.DeletedBy.Int),
			DeletedAt: v.DeletedAt.Time,
		}

		res = append(res, Account{
			ID:              int64(v.ID),
			Name:            v.Name,
			Email:           v.Email,
			BaseInformation: creationInfo,
		})
	}

	return res
}
