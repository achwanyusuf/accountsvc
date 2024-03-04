package account

import (
	"time"

	"github.com/achwanyusuf/carrent-accountsvc/src/domain/account"
	"github.com/achwanyusuf/carrent-accountsvc/src/domain/accountrole"
	"github.com/achwanyusuf/carrent-accountsvc/src/domain/role"
	"github.com/achwanyusuf/carrent-accountsvc/src/model"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/psqlmodel"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/svcerr"
	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
	"github.com/achwanyusuf/carrent-lib/pkg/hash"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/volatiletech/null/v8"
)

type AccountDep struct {
	log         logger.Logger
	conf        Conf
	account     account.AccountInterface
	role        role.RoleInterface
	accountRole accountrole.AccountRoleInterface
}

type Conf struct {
	TokenTimeout time.Duration `mapstructure:"token_timeout"`
	TokenSecret  string        `mapstructure:"token_secret"`
	AESSecret    string        `mapstructure:"aes_secret"`
}

type AccountInterface interface {
	Oauth2(ctx *gin.Context, v model.Login) (model.Auth, error)
	Create(ctx *gin.Context, v model.Register) (model.Account, error)
	GetByParam(ctx *gin.Context, cacheControl string, v model.GetAccountsByParam) ([]model.Account, model.Pagination, error)
	GetByID(ctx *gin.Context, cacheControl string, id int64) (model.Account, error)
	UpdateByID(ctx *gin.Context, id int64, v model.UpdateAccountData) (model.Account, error)
	UpdatePasswordByID(ctx *gin.Context, id int64, v model.UpdatePasswordData) (model.Account, error)
	DeleteByID(ctx *gin.Context, id int64, isHardDelete bool, vid int64) error
}

func New(conf Conf, logger *logger.Logger, account account.AccountInterface, role role.RoleInterface, accountRole accountrole.AccountRoleInterface) AccountInterface {
	return &AccountDep{
		conf:        conf,
		log:         *logger,
		account:     account,
		role:        role,
		accountRole: accountRole,
	}
}

func (a *AccountDep) Oauth2(ctx *gin.Context, v model.Login) (model.Auth, error) {
	var auth model.Auth
	err := v.Validate()
	if err != nil {
		return auth, err
	}
	role, err := a.role.GetSingleByParam(ctx, "", &model.GetRoleByParam{
		Cid: null.NewString(v.ClientID, true),
	})
	if err != nil {
		return auth, errormsg.WrapErr(svcerr.AccountSVCNotAuthorized, err, "role not found")
	}

	if match := hash.CompareAES(role.Sec, a.conf.AESSecret, v.ClientSecret); !match {
		return auth, errormsg.WrapErr(svcerr.AccountSVCNotAuthorized, err, "invalid client id/client secret")
	}

	account, err := a.account.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountByParam{
		Email: null.NewString(v.Email, true),
	})
	if err != nil {
		return auth, errormsg.WrapErr(svcerr.AccountSVCNotAuthorized, err, "account not found")
	}

	_, err = a.accountRole.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountRoleByParam{
		AccountID: null.NewInt64(int64(account.ID), true),
		RoleID:    null.NewInt64(int64(role.ID), true),
	})
	if err != nil {
		return auth, errormsg.WrapErr(svcerr.AccountSVCNotAuthorized, err, "invalid client id/client secret")
	}

	err = hash.Compare(account.Password, v.Password)
	if err != nil {
		return auth, errormsg.WrapErr(svcerr.AccountSVCInvalidPasswordNotMatch, err, "password not match")
	}

	token := jwt.New(jwt.SigningMethodHS512)
	expired := time.Now().Add(a.conf.TokenTimeout)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = account.ID
	claims["username"] = account.Email
	claims["exp"] = expired.Unix()
	claims["scope"] = role.Scope
	t, err := token.SignedString([]byte(a.conf.TokenSecret))
	if err != nil {
		return auth, errormsg.WrapErr(svcerr.AccountSVCInvalidPasswordNotMatch, err, "invalid token")
	}

	auth = model.Auth{
		AccessToken: t,
		Exp:         &expired,
		TokenType:   model.TokenTypeBearer,
		Scope:       role.Scope,
	}

	return auth, nil
}

func (a *AccountDep) Create(ctx *gin.Context, v model.Register) (model.Account, error) {
	var result model.Account
	err := v.Validate()
	if err != nil {
		return result, err
	}

	pwd, err := hash.Hash(v.Password)
	if err != nil {
		return result, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error hash password")
	}

	account := &psqlmodel.Account{
		Name:      v.Name,
		Email:     v.Email,
		Password:  pwd,
		CreatedBy: int(v.CreatedBy),
		UpdatedBy: int(v.CreatedBy),
	}

	err = a.account.Insert(ctx, account)
	if err != nil {
		return result, err
	}

	return model.TransformPSQLSingleAccount(account), nil
}

func (a *AccountDep) GetByParam(ctx *gin.Context, cacheControl string, v model.GetAccountsByParam) ([]model.Account, model.Pagination, error) {
	accountSlice, pagination, err := a.account.GetByParam(ctx, cacheControl, &v)
	if err != nil {
		return []model.Account{}, model.Pagination{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error get by param")
	}
	return model.TransformPSQLAccount(&accountSlice), pagination, nil
}

func (a *AccountDep) GetByID(ctx *gin.Context, cacheControl string, id int64) (model.Account, error) {
	account, err := a.account.GetSingleByParam(ctx, cacheControl, &model.GetAccountByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.Account{}, errormsg.WrapErr(svcerr.AccountSVCNotFound, err, "data not found")
	}
	return model.TransformPSQLSingleAccount(&account), nil
}

func (a *AccountDep) UpdateByID(ctx *gin.Context, id int64, v model.UpdateAccountData) (model.Account, error) {
	if err := v.Validate(); err != nil {
		return model.Account{}, err
	}

	account, err := a.account.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.Account{}, err
	}

	if v.Name == account.Name {
		return model.TransformPSQLSingleAccount(&account), nil
	}

	account.Name = v.Name
	account.UpdatedBy = int(v.UpdateBy)

	err = a.account.Update(ctx, &account)
	if err != nil {
		return model.Account{}, err
	}

	return model.TransformPSQLSingleAccount(&account), nil
}

func (a *AccountDep) UpdatePasswordByID(ctx *gin.Context, id int64, v model.UpdatePasswordData) (model.Account, error) {
	if err := v.IsValid(); err != nil {
		return model.Account{}, err
	}

	account, err := a.account.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.Account{}, err
	}
	pwd, err := hash.Hash(v.Password)
	if err != nil {
		return model.Account{}, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error hash password")
	}

	account.Password = pwd
	account.UpdatedBy = int(v.UpdateBy)

	err = a.account.Update(ctx, &account)
	if err != nil {
		return model.Account{}, err
	}
	return model.TransformPSQLSingleAccount(&account), nil
}

func (a *AccountDep) DeleteByID(ctx *gin.Context, id int64, isHardDelete bool, vid int64) error {
	account, err := a.account.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountByParam{
		ID: null.NewInt64(vid, true),
	})
	if err != nil {
		return err
	}
	return a.account.Delete(ctx, &account, id, isHardDelete)
}
