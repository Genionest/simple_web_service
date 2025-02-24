package controllers

import (
	"errors"
	"net/http"

	"example.com/app/global"
	"example.com/app/models"
	"example.com/app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {
	var user models.User
	// 绑定发送的数据
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"info":  "ShouldBindJSON",
		})
		return
	}
	// 加密
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"info":  "HashPassword",
		})
		return
	}
	user.Password = hashedPwd
	// JWT: Json Web Token
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"info":  "GenerateJWT",
		})
		return
	}

	// 自动在database里生成table
	if err := global.DB.AutoMigrate(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"info":  "AutoMigrate",
		})
		return
	}

	// 创建数据
	if err := global.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"info":  "DB.Create",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 验证用户名
	var user models.User
	if err := global.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		// ctx.JSON(http.StatusNotFound, gin.H{
		// 	"error": "wrong credentials",
		// })
		// return
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "wrong credentials2",
		})
		return
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"info":  "GenerateJWT",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
