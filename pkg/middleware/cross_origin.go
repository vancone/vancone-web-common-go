package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CrossOriginMiddleware(context *gin.Context) {
	method := context.Request.Method
	allowOrigin := context.Request.Referer()
	if strings.Contains(allowOrigin, "vancone.com") {
		if allowOrigin[len(allowOrigin)-1] == '/' {
			allowOrigin = allowOrigin[0 : len(allowOrigin)-1]
		}
		context.Header("Access-Control-Allow-Origin", allowOrigin)
		context.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		context.Header("Access-Control-Allow-Credentials", "true")
	}

	if method == "OPTIONS" {
		context.AbortWithStatus(http.StatusOK)
	}
}
