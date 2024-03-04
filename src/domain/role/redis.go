package role

import (
	"encoding/json"

	"github.com/achwanyusuf/carrent-accountsvc/src/model"
	"github.com/achwanyusuf/carrent-accountsvc/src/model/psqlmodel"
	"github.com/gin-gonic/gin"
)

func (r *RoleDep) getSingleByParamRedis(ctx *gin.Context, key string) (psqlmodel.Role, error) {
	var res psqlmodel.Role
	data, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *RoleDep) setRedis(ctx *gin.Context, key string, data string) error {
	expTime := r.Conf.RedisExpirationTime
	if r.Conf.RedisExpirationTime == 0 {
		expTime = model.DefaultRedisExpiration
	}
	_, err := r.Redis.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	_, err = r.Redis.Set(ctx, key, data, expTime).Result()
	return err
}

func (r *RoleDep) getByParamRedis(ctx *gin.Context, key string) (psqlmodel.RoleSlice, error) {
	var res psqlmodel.RoleSlice
	data, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *RoleDep) getByParamPaginationRedis(ctx *gin.Context, key string) (model.Pagination, error) {
	var res model.Pagination
	data, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
