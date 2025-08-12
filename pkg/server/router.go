package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vancone/vancone-web-common-go/pkg/config"
)

type Router interface {
	Load(engine *gin.Engine)
}

func StartServer(g *gin.Engine) {
	log.Println("Server started, listening on port", config.Server.Port)
	err := g.Run(fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		log.Panicln("Failed to start server:", err)
	}
}
