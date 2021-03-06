package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iissy/goweb/src/cli"
	"github.com/iissy/goweb/src/model"
	"github.com/iissy/goweb/src/redis"
	"github.com/iissy/goweb/src/utils"
	"github.com/micro/go-micro/v2/config"
	"time"
)

func Login(ctx *gin.Context) {
	user := new(model.Account)
	err := ctx.BindJSON(&user)
	if ok := utils.WriteErrorLog(ctx.FullPath(), err); ok {
		fmt.Print(0)
	}

	user.Password = utils.GetMd5String(user.Password)
	user.LastLoginDate = time.Now().Format("2006-01-02 15:04:05")
	result := new(model.Account)
	err = cli.Call("Login", user, result)
	if ok := utils.WriteErrorLog(ctx.FullPath(), err); ok {
		fmt.Print(0)
	}

	if result.ID > 0 {
		token := utils.Random62String(64)
		ctx.Header(utils.ASYUSERID, result.UserId)
		ctx.Header(utils.ASYTOKEN, token)
		ctx.SetCookie(utils.ASYUSERID, result.UserId, 3600, "/", config.Get("domain").String("localhost"), false, true)
		ctx.SetCookie(utils.ASYTOKEN, token, 3600, "/", config.Get("domain").String("localhost"), false, true)
		err = redis.Set(result.UserId, token)
		utils.WriteErrorLog(ctx.FullPath(), err)
	}

	ctx.JSON(200, result)
}
