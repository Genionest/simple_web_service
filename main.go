package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/app/config"
	"example.com/app/router"
)

func main() {
	config.InitConfig()
	// fmt.Printf("%v\n", config.AppConfig.App.Port)
	// r := gin.Default()
	// r.GET("/ping", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"message": "pong\n",
	// 	})
	// })
	host := config.AppConfig.App.Host
	port := config.AppConfig.App.Port

	r := router.SetupRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", host, port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	r.Run(fmt.Sprintf("%v:%v", host, port))
}
