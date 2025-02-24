package controllers

import (
	"net/http"

	"example.com/app/global"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func LikeArticle(ctx *gin.Context) {
	articleID := ctx.Param("id")

	likeKey := "article:" + articleID + ":like"

	if err := global.RedisDB.Incr(likeKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successful like Article",
	})
}

func GetArticleLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")

	likeKey := "article:" + articleID + ":like"

	likes, err := global.RedisDB.Get(likeKey).Result()

	if err == redis.Nil {
		likes = "0"
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"likes": likes,
	})
}
