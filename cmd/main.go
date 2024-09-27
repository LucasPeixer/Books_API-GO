package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	conn, err := db.connectDB()
	if(err != nil){
		panic(err)
	}

	server := gin.Default()

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
				"message": "pong",
		})
	})

	server.Run(":8000")
}