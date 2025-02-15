package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"leetroll/app/service/dto"
	"leetroll/auth"
	"leetroll/common/apis"
	"leetroll/config"
	"leetroll/db/handlers"
	"leetroll/db/models"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	apis.Api
}

// 查询用户信息
func (e User) GetUser(ctx *gin.Context) {
	req := dto.UserApiReq{}
	service := handlers.UserHandler{}
	err := e.MakeContext(ctx).
		MakeDB().
		Bind(&req, nil).
		MakeService(&service.Handler).
		Errors

	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	var user models.User
	err = service.FindById(req.GetId(), &user).Error
	if err != nil {
		e.Error(500, err, "fail")
		log.Fatal(err.Error())
		return
	}

	e.OK(user, "ok")

}

// 用户信息编辑
func (e User) UpdateUser(ctx *gin.Context) {

	req := dto.UserUpdateApiReq{}
	service := handlers.UserHandler{}
	err := e.MakeContext(ctx).
		MakeDB().
		Bind(&req, binding.JSON).
		MakeService(&service.Handler).
		Errors

	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	user := req.Generate()
	err = service.Update(&user).Error
	if err != nil {
		e.Error(500, err, "fail")
		return
	}

	e.OK(user, "ok")
}

// 修改头像
func (e User) UpdateUserAvatar(ctx *gin.Context) {
	req := dto.UserUpdateAvatarApiReq{}
	service := handlers.UserHandler{}
	err := e.MakeContext(ctx).
		MakeDB().
		Bind(&req, binding.JSON).
		MakeService(&service.Handler).
		Errors

	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	err = service.UpdateAvatar(req.ID, req.Avatar).Error
	if err != nil {
		e.Error(500, err, "fail")
		return
	}

	e.OK(nil, "ok")
}

func (e User) UpdateUserBG(ctx *gin.Context) {
	req := dto.UserUpdateBGApiReq{}
	service := handlers.UserHandler{}
	err := e.MakeContext(ctx).
		MakeDB().
		Bind(&req, binding.JSON).
		MakeService(&service.Handler).
		Errors

	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	err = service.UpdateBG(req.ID, req.BgImag).Error
	if err != nil {
		e.Error(500, err, "fail")
		return
	}

	e.OK(nil, "ok")
}

// 登录
func (e User) Login(ctx *gin.Context) {
	req := dto.UserLoginApiReq{}
	service := handlers.UserHandler{}
	err := e.MakeContext(ctx).
		MakeDB().
		Bind(&req, binding.JSON).
		MakeService(&service.Handler).
		Errors

	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	//check if exist
	var user models.User
	err = service.FindUserByPhone(&req, &user).Error
	if err != nil {
		e.Error(500, err, "user is not exist!")
		return
	}

	//create auth token
	token, err := auth.CreateToken(user.ID)

	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	saveErr := auth.CreateAuth(uint64(user.ID), token)
	if saveErr != nil {
		e.Error(500, saveErr, saveErr.Error())
		return
	}

	resp := map[string]interface{}{
		"User":         user,
		"AccessToken":  token.AccessToken,
		"RefreshToken": token.RefreshToken,
	}
	e.OK(resp, "ok")

}

// 登出
func (e User) Logout(c *gin.Context) {
	au, err := auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		e.Error(http.StatusUnauthorized, err, "unauthorized")
		return
	}
	deleted, delErr := auth.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		e.Error(http.StatusUnauthorized, err, "unauthorized")
		return
	}
	e.OK(http.StatusOK, "Successfully logged out")
}

// token刷新
func (e User) Refresh(ctx *gin.Context) {
	req := dto.UserTokenRefreshApiReq{}
	err := e.MakeContext(ctx).
		MakeDB().
		Bind(&req, binding.JSON).
		Errors

	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	token, err := auth.ParseToken(req.RefreshToken, config.JWTConfig.RefreshSecret)

	//if there is an error, the token must have expired
	if err != nil {
		e.Error(http.StatusUnauthorized, err, "Refresh token expired")
		return
	}

	//is token valid
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		e.Error(http.StatusUnauthorized, err, "Refresh token is not valid")
		return
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			e.Error(http.StatusUnprocessableEntity, nil, "StatusUnprocessableEntity")
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			e.Error(http.StatusUnprocessableEntity, err, "StatusUnprocessableEntity")
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := auth.DeleteAuth(refreshUuid)
		if delErr != nil || deleted != 0 {
			e.Error(http.StatusUnauthorized, delErr, "unauthorized")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := auth.CreateToken(int(userId))
		if createErr != nil {
			e.Error(http.StatusForbidden, createErr, "StatusForbidden")
			return
		}
		//save the tokens metadata to redis
		saveErr := auth.CreateAuth(userId, ts)
		if saveErr != nil {
			e.Error(http.StatusForbidden, saveErr, "StatusForbidden")
			return
		}
		tokens := map[string]string{
			"AccessToken":  ts.AccessToken,
			"RefreshToken": ts.RefreshToken,
		}
		e.OK(tokens, "refresh success")
	} else {
		e.Error(http.StatusUnauthorized, err, "Refresh token expired")
	}
}
