package account

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/achwanyusuf/carrent-accountsvc/src/model"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/svcerr"
	"github.com/achwanyusuf/carrent-accountsvc/src/usecase/account"
	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type AccountDep struct {
	log     logger.Logger
	account account.AccountInterface
	conf    Conf
}

type Conf struct {
	TokenSecret  string `mapstructure:"token_secret"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type AccountInterface interface {
	Oauth2(ctx *gin.Context)
	CurrentAccount(ctx *gin.Context)
	UpdateCurrentAccount(ctx *gin.Context)
	UpdatePasswordAccount(ctx *gin.Context)
	Register(ctx *gin.Context)
	Create(ctx *gin.Context)
	Read(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	UpdateByID(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
}

func New(conf Conf, log *logger.Logger, acc account.AccountInterface) AccountInterface {
	return &AccountDep{
		conf:    conf,
		log:     *log,
		account: acc,
	}
}

// Oauth2 godoc
// @Summary OAUTH2 Authorization
// @Description OAUTH2 Authorization Code flow will show generated token to access apps
// @Tags account
// @Accept x-www-form-urlencoded
// @Produce json
// @Param client_id header string true "Client ID"
// @Param client_secret header string true "Client Secret"
// @Param username formData string true "Account Email"
// @Param password formData string true "Account Password"
// @Success 200 {object} model.LoginResponse
// @Success 400 {object} model.LoginResponse
// @Success 401 {object} model.LoginResponse
// @Success 500 {object} model.LoginResponse
// @Router /oauth2 [post]
func (a *AccountDep) Oauth2(ctx *gin.Context) {
	var (
		response               model.LoginResponse
		authorization          = ctx.GetHeader("Authorization")
		clientID, clientSecret string
	)

	if authorization != "" {
		clientID, clientSecret = a.decodeClient(ctx, authorization)
	} else {
		clientID = ctx.GetHeader("client_id")
		clientSecret = ctx.GetHeader("client_secret")
	}

	loginData := model.Login{
		Email:        ctx.Request.FormValue("username"),
		Password:     ctx.Request.FormValue("password"),
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	if loginData.Email == "" {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error read body"))
			ctx.JSON(statusCode, response)
			return
		}
		if err = json.Unmarshal(body, &loginData); err != nil {
			statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error unmarshal body"))
			ctx.JSON(statusCode, response)
			return
		}
	}

	auth, err := a.account.Oauth2(ctx, loginData)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Auth = auth

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}

func (a *AccountDep) decodeClient(ctx *gin.Context, auth string) (string, string) {
	client := strings.Split(auth, " ")
	if len(client) < 2 {
		return "", ""
	}
	rawDecodedText, err := base64.StdEncoding.DecodeString(client[1])
	if err != nil {
		a.log.Warn(ctx, err)
	}
	secret := strings.Split(string(rawDecodedText), ":")
	if len(secret) < 2 {
		return "", ""
	}
	return secret[0], secret[1]
}

// Current Account godoc
// @Summary Get current account data
// @Description Get current account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} model.SingleAccountResponse
// @Success 400 {object} model.SingleAccountResponse
// @Success 401 {object} model.SingleAccountResponse
// @Success 500 {object} model.SingleAccountResponse
// @Router /me [get]
func (a *AccountDep) CurrentAccount(ctx *gin.Context) {
	var response model.SingleAccountResponse
	cacheControl := ctx.GetHeader("Cache-Control")
	result, err := a.account.GetByID(ctx, cacheControl, ctx.Value("id").(int64))
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}

// Update Current Account godoc
// @Summary Update current account data
// @Description Update current account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} model.SingleAccountResponse
// @Success 400 {object} model.SingleAccountResponse
// @Success 401 {object} model.SingleAccountResponse
// @Success 500 {object} model.SingleAccountResponse
// @Router /me [put]
func (a *AccountDep) UpdateCurrentAccount(ctx *gin.Context) {
	var (
		updateData model.UpdateAccountData
		response   model.SingleAccountResponse
	)

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error read body"))
		ctx.JSON(statusCode, response)
		return
	}

	if err = json.Unmarshal(body, &updateData); err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error unmarshal body"))
		ctx.JSON(statusCode, response)
		return
	}

	updateData.UpdateBy = ctx.Value("id").(int64)
	result, err := a.account.UpdateByID(ctx, ctx.Value("id").(int64), updateData)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}

// Update Password Account godoc
// @Summary Update password account data
// @Description Update password account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} model.SingleAccountResponse
// @Success 400 {object} model.SingleAccountResponse
// @Success 401 {object} model.SingleAccountResponse
// @Success 500 {object} model.SingleAccountResponse
// @Router /me/password [put]
func (a *AccountDep) UpdatePasswordAccount(ctx *gin.Context) {
	var (
		updateData model.UpdatePasswordData
		response   model.SingleAccountResponse
	)

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error read body"))
		ctx.JSON(statusCode, response)
		return
	}

	if err = json.Unmarshal(body, &updateData); err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error unmarshal body"))
		ctx.JSON(statusCode, response)
		return
	}

	updateData.UpdateBy = ctx.Value("id").(int64)
	result, err := a.account.UpdatePasswordByID(ctx, ctx.Value("id").(int64), updateData)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}

// Register godoc
// @Summary Register account
// @Description Register to create access from guest
// @Tags account
// @Accept json
// @Produce json
// @Param data body model.Register true "Account Data"
// @Success 200 {object} model.RegisterResponse
// @Success 400 {object} model.RegisterResponse
// @Success 500 {object} model.RegisterResponse
// @Router /account [post]
func (a *AccountDep) Register(ctx *gin.Context) {
	var (
		registerData model.Register
		result       model.Account
		response     model.RegisterResponse
	)
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusCreated, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error read body"))
		ctx.JSON(statusCode, response)
		return
	}

	if err = json.Unmarshal(body, &registerData); err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusCreated, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error unmarshal body"))
		ctx.JSON(statusCode, response)
		return
	}

	result, err = a.account.Create(ctx, registerData)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusCreated, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusCreated, nil)
	ctx.JSON(statusCode, response)
}

// Create Account godoc
// @Summary Create account
// @Description Create account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param data body model.Register true "Account Data"
// @Success 200 {object} model.RegisterResponse
// @Success 400 {object} model.RegisterResponse
// @Success 500 {object} model.RegisterResponse
// @Router /account [post]
func (a *AccountDep) Create(ctx *gin.Context) {
	var (
		registerData model.Register
		result       model.Account
		response     model.RegisterResponse
	)

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusCreated, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error read body"))
		ctx.JSON(statusCode, response)
		return
	}

	if err = json.Unmarshal(body, &registerData); err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusCreated, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error unmarshal body"))
		ctx.JSON(statusCode, response)
		return
	}

	registerData.CreatedBy = ctx.Value("id").(int64)
	result, err = a.account.Create(ctx, registerData)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusCreated, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusCreated, nil)
	ctx.JSON(statusCode, response)
}

// Get Accounts Data godoc
// @Summary Get accounts data
// @Description Get accounts data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id query string false "search by id"
// @Param name query string false "search by name"
// @Param email query string false "search by email"
// @Param sort_by query string false "sort result by attributes"
// @Param page query int false " "
// @Param limit query int false " "
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} model.AccountsResponse
// @Success 400 {object} model.AccountsResponse
// @Success 500 {object} model.AccountsResponse
// @Router /account [get]
func (a *AccountDep) Read(ctx *gin.Context) {
	var (
		param    model.GetAccountsByParam
		response model.AccountsResponse
	)
	cacheControl := ctx.GetHeader("Cache-Control")
	var decoder = schema.NewDecoder()
	err := decoder.Decode(&param, ctx.Request.URL.Query())
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}
	accounts, pagination, err := a.account.GetByParam(ctx, cacheControl, param)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = accounts
	response.Pagination = pagination

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}

// Get Accounts Data godoc
// @Summary Get accounts data
// @Description Get accounts data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "get by id"
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} model.SingleAccountResponse
// @Success 400 {object} model.SingleAccountResponse
// @Success 500 {object} model.SingleAccountResponse
// @Router /account/{id} [get]
func (a *AccountDep) GetByID(ctx *gin.Context) {
	var response model.SingleAccountResponse
	cacheControl := ctx.GetHeader("Cache-Control")
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(errormsg.Error400, err, "error get id"))
		ctx.JSON(statusCode, response)
		return
	}
	result, err := a.account.GetByID(ctx, cacheControl, id)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}

// Update Account Data godoc
// @Summary Update account data
// @Description Update account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "update by id"
// @Param data body model.UpdateAccountData true "Account Data"
// @Success 200 {object} model.SingleAccountResponse
// @Success 400 {object} model.SingleAccountResponse
// @Success 500 {object} model.SingleAccountResponse
// @Router /account/{id} [put]
func (a *AccountDep) UpdateByID(ctx *gin.Context) {
	var (
		updateData model.UpdateAccountData
		response   model.SingleAccountResponse
	)

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(errormsg.Error400, err, "error get id"))
		ctx.JSON(statusCode, response)
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error read body"))
		ctx.JSON(statusCode, response)
		return
	}

	if err = json.Unmarshal(body, &updateData); err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error unmarshal body"))
		ctx.JSON(statusCode, response)
		return
	}
	updateData.UpdateBy = ctx.Value("id").(int64)
	scope := ctx.Value("scope").(string)
	if scope != model.SuperAdminScope {
		id = ctx.Value("id").(int64)
	}
	result, err := a.account.UpdateByID(ctx, id, updateData)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}

// Delete Account Data godoc
// @Summary Delete account data
// @Description Delete account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "delete by id"
// @Success 200 {object} model.EmptyResponse
// @Success 400 {object} model.EmptyResponse
// @Success 500 {object} model.EmptyResponse
// @Router /account/{id} [delete]
func (a *AccountDep) DeleteByID(ctx *gin.Context) {
	var (
		response model.EmptyResponse
	)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(errormsg.Error400, err, "error get id"))
		ctx.JSON(statusCode, response)
		return
	}
	scope := ctx.Value("scope").(string)
	if scope != model.SuperAdminScope {
		id = ctx.Value("id").(int64)
	}
	err = a.account.DeleteByID(ctx, ctx.Value("id").(int64), false, id)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}
