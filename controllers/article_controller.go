package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"example.com/app/global"
	"example.com/app/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var cacheKey = "articles"

func CreateArticle(ctx *gin.Context) {
	var article models.Article
	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := global.DB.AutoMigrate(&article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := global.DB.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 删除缓存
	if err := global.RedisDB.Del(cacheKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, article)
}

func GetArticles(ctx *gin.Context) {
	var articles []models.Article

	// 使用缓存
	cachedData, err := global.RedisDB.Get(cacheKey).Result()
	if err == redis.Nil {
		// fmt.Println("redis.Nil")
		if err := global.DB.Find(&articles).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		articleJSON, err := json.Marshal(articles)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := global.RedisDB.Set(cacheKey, articleJSON, 10*time.Minute).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		err := json.Unmarshal([]byte(cachedData), &articles)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}

	ctx.JSON(http.StatusOK, articles)

}

func GetArticleByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var article models.Article
	if err := global.DB.Where("id = ?", id).First(&article).Error; err != nil {
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

	ctx.JSON(http.StatusOK, article)
}
