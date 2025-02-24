package config

import (
	"fmt"
	"log"
	"time"

	"example.com/app/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() {
	conf := AppConfig
	str := "%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(str,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to initialize database, got error %v", err)
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(conf.Database.MaxIdelConns)
	sqlDB.SetMaxOpenConns(conf.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err != nil {
		log.Fatalf("Failed to configure database, got error %v", err)
	}

	global.DB = db

}
