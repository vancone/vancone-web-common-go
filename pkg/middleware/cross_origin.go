package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CrossOriginMiddleware(ctx *gin.Context) {
	allowOrigin := ctx.Request.Referer()
	allowHeaders := strings.Join([]string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
		"Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"}, ", ")
	allowMethods := strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodDelete, http.MethodOptions}, ", ")

	if strings.Contains(allowOrigin, "vancone.com") {
		if allowOrigin[len(allowOrigin)-1] == '/' {
			allowOrigin = allowOrigin[0 : len(allowOrigin)-1]
		}
		ctx.Header("Access-Control-Allow-Origin", allowOrigin)
		ctx.Header("Access-Control-Allow-Headers", allowHeaders)
		ctx.Header("Access-Control-Allow-Methods", allowMethods)
		ctx.Header("Access-Control-Allow-Credentials", "true")
	}

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
	}
}
