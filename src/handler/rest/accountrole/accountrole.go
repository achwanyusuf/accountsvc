package accountrole

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/achwanyusuf/carrent-accountsvc/src/model"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/svcerr"
	"github.com/achwanyusuf/carrent-accountsvc/src/usecase/accountrole"
	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type AccountRoleDep struct {
	log         logger.Logger
	accountrole accountrole.AccountRoleInterface
	conf        Conf
}

type Conf struct {
	TokenSecret  string `mapstructure:"token_secret"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type AccountRoleInterface interface {
	Create(ctx *gin.Context)
	Read(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
}

func New(conf Conf, log *logger.Logger, accountrole accountrole.AccountRoleInterface) AccountRoleInterface {
	return &AccountRoleDep{
		conf:        conf,
		log:         *log,
		accountrole: accountrole,
	}
}

// Create AccountRole godoc
// @Summary Create AccountRole
// @Description Create account role data
// @Tags account-role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param data body model.CreateAccountRole true "AccountRole Data"
// @Success 200 {object} model.SingleAccountRoleResponse
// @Success 400 {object} model.SingleAccountRoleResponse
// @Success 500 {object} model.SingleAccountRoleResponse
// @Router /account-role [post]
func (a *AccountRoleDep) Create(ctx *gin.Context) {
	var (
		roleData model.CreateAccountRole
		result   model.AccountRole
		response model.SingleAccountRoleResponse
	)

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusCreated, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error read body"))
		ctx.JSON(statusCode, response)
		return
	}

	if err = json.Unmarshal(body, &roleData); err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusCreated, errormsg.WrapErr(svcerr.AccountSVCBadRequest, err, "error unmarshal body"))
		ctx.JSON(statusCode, response)
		return
	}

	roleData.CreatedBy = ctx.Value("id").(int64)
	result, err = a.accountrole.Create(ctx, roleData)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusCreated, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusCreated, nil)
	ctx.JSON(statusCode, response)
}

// Get AccountRoles Data godoc
// @Summary Get account roles data
// @Description Get account roles data
// @Tags account-role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id query string false "search by id"
// @Param account_id query int false "search by account id"
// @Param role_id query int false "search by role id"
// @Param sort_by query string false "sort result by attributes"
// @Param page query int false " "
// @Param limit query int false " "
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} model.AccountRolesResponse
// @Success 400 {object} model.AccountRolesResponse
// @Success 500 {object} model.AccountRolesResponse
// @Router /account-role [get]
func (a *AccountRoleDep) Read(ctx *gin.Context) {
	var (
		param    model.GetAccountRolesByParam
		response model.AccountRolesResponse
	)
	cacheControl := ctx.GetHeader("Cache-Control")
	var decoder = schema.NewDecoder()
	err := decoder.Decode(&param, ctx.Request.URL.Query())
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}
	roles, pagination, err := a.accountrole.GetByParam(ctx, cacheControl, param)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = roles
	response.Pagination = pagination

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}

// Get AccountRoles Data godoc
// @Summary Get account role by id data
// @Description Get account role by id data
// @Tags account-role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "get by id"
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} model.SingleAccountRoleResponse
// @Success 400 {object} model.SingleAccountRoleResponse
// @Success 500 {object} model.SingleAccountRoleResponse
// @Router /account-role/{id} [get]
func (a *AccountRoleDep) GetByID(ctx *gin.Context) {
	var response model.SingleAccountRoleResponse
	cacheControl := ctx.GetHeader("Cache-Control")
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(errormsg.Error400, err, "error get id"))
		ctx.JSON(statusCode, response)
		return
	}
	result, err := a.accountrole.GetByID(ctx, cacheControl, id)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}

// Delete AccountRole Data godoc
// @Summary Delete account role data
// @Description Delete account role data
// @Tags account-role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "delete by id"
// @Success 200 {object} model.EmptyResponse
// @Success 400 {object} model.EmptyResponse
// @Success 500 {object} model.EmptyResponse
// @Router /account-role/{id} [delete]
func (a *AccountRoleDep) DeleteByID(ctx *gin.Context) {
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
	err = a.accountrole.DeleteByID(ctx, ctx.Value("id").(int64), false, id)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}
