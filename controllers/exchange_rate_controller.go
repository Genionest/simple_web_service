package controllers

import (
	"net/http"
	"time"

	"example.com/app/global"
	"example.com/app/models"
	"github.com/gin-gonic/gin"
)

func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate

	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	exchangeRate.Date = time.Now()

	if err := global.DB.AutoMigrate(&exchangeRate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := global.DB.Create(&exchangeRate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, exchangeRate)
}

func GetExchangeRates(ctx *gin.Context) {
	var exchangeRates []models.ExchangeRate

	if err := global.DB.Find(&exchangeRates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, exchangeRates)
}
