package server

import "github.com/gin-gonic/gin"

type Router interface {
	Load(engine *gin.Engine)
}
