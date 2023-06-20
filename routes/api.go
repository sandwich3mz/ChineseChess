package routes

import (
	"chesss/internal/logic"
	"github.com/gin-gonic/gin"
	"log"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	chess := router.Group("/chess")
	{
		chess.GET("/isValid/:before/:after", logic.IsValid)
		chess.GET("/move/:before", logic.Move)
		chess.POST("/memorize", logic.Memorize)
	}
	return router
}

// RunServer 启动服务器
func RunServer() {
	ginServer := setupRouter()
	err := ginServer.Run(":8083")
	if err != nil {
		log.Printf("Failed to run the ginServer: %v", err)
	}
}
