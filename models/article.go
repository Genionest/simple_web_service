package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	// 				不使用`json:"title"`那么json返回的会是Title而非title
	Title   string `binding:"required"`
	Content string `binding:"required"`
	Preview string `binding:"required"`
	// redis解决
	// Likes   int    `gorm:"default:0"`
}
