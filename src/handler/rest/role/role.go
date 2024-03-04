package role

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/achwanyusuf/carrent-accountsvc/src/model"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/svcerr"
	"github.com/achwanyusuf/carrent-accountsvc/src/usecase/role"
	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type RoleDep struct {
	log  logger.Logger
	role role.RoleInterface
	conf Conf
}

type Conf struct {
	TokenSecret  string `mapstructure:"token_secret"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type RoleInterface interface {
	Create(ctx *gin.Context)
	Read(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	UpdateByID(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
}

func New(conf Conf, log *logger.Logger, role role.RoleInterface) RoleInterface {
	return &RoleDep{
		conf: conf,
		log:  *log,
		role: role,
	}
}

// Create Role godoc
// @Summary Create Role
// @Description Create role data
// @Tags role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param data body model.CreateRole true "Role Data"
// @Success 200 {object} model.SingleRoleResponse
// @Success 400 {object} model.SingleRoleResponse
// @Success 500 {object} model.SingleRoleResponse
// @Router /role [post]
func (a *RoleDep) Create(ctx *gin.Context) {
	var (
		roleData model.CreateRole
		result   model.Role
		response model.SingleRoleResponse
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
	result, err = a.role.Create(ctx, roleData)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusCreated, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusCreated, nil)
	ctx.JSON(statusCode, response)
}

// Get Roles Data godoc
// @Summary Get roles data
// @Description Get roles data
// @Tags role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id query string false "search by id"
// @Param scope query string false "search by scope"
// @Param cid query string false "search by client_id"
// @Param sort_by query string false "sort result by attributes"
// @Param page query int false " "
// @Param limit query int false " "
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} model.RolesResponse
// @Success 400 {object} model.RolesResponse
// @Success 500 {object} model.RolesResponse
// @Router /role [get]
func (a *RoleDep) Read(ctx *gin.Context) {
	var (
		param    model.GetRolesByParam
		response model.RolesResponse
	)
	cacheControl := ctx.GetHeader("Cache-Control")
	var decoder = schema.NewDecoder()
	err := decoder.Decode(&param, ctx.Request.URL.Query())
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}
	roles, pagination, err := a.role.GetByParam(ctx, cacheControl, param)
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

// Get Roles Data godoc
// @Summary Get roles data
// @Description Get roles data
// @Tags role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "get by id"
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} model.SingleRoleResponse
// @Success 400 {object} model.SingleRoleResponse
// @Success 500 {object} model.SingleRoleResponse
// @Router /role/{id} [get]
func (a *RoleDep) GetByID(ctx *gin.Context) {
	var response model.SingleRoleResponse
	cacheControl := ctx.GetHeader("Cache-Control")
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(errormsg.Error400, err, "error get id"))
		ctx.JSON(statusCode, response)
		return
	}
	result, err := a.role.GetByID(ctx, cacheControl, id)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}

// Update Role Data godoc
// @Summary Update role data
// @Description Update role data
// @Tags role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "update by id"
// @Param data body model.UpdateRole true "Role Data"
// @Success 200 {object} model.SingleRoleResponse
// @Success 400 {object} model.SingleRoleResponse
// @Success 500 {object} model.SingleRoleResponse
// @Router /role/{id} [put]
func (a *RoleDep) UpdateByID(ctx *gin.Context) {
	var (
		updateData model.UpdateRole
		response   model.SingleRoleResponse
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
	updateData.UpdatedBy = ctx.Value("id").(int64)
	scope := ctx.Value("scope").(string)
	if scope != model.SuperAdminScope {
		id = ctx.Value("id").(int64)
	}
	result, err := a.role.UpdateByID(ctx, id, updateData)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	response.Data = result

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}

// Delete Role Data godoc
// @Summary Delete role data
// @Description Delete role data
// @Tags role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "delete by id"
// @Success 200 {object} model.EmptyResponse
// @Success 400 {object} model.EmptyResponse
// @Success 500 {object} model.EmptyResponse
// @Router /role/{id} [delete]
func (a *RoleDep) DeleteByID(ctx *gin.Context) {
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
	err = a.role.DeleteByID(ctx, ctx.Value("id").(int64), false, id)
	if err != nil {
		statusCode := response.Transform(ctx, a.log, http.StatusOK, err)
		ctx.JSON(statusCode, response)
		return
	}

	statusCode := response.Transform(ctx, a.log, http.StatusOK, nil)
	ctx.JSON(statusCode, response)
}
