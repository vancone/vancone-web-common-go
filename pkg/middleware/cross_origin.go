package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CrossOriginMiddleware(ctx *gin.Context) {
	allowHeaders := strings.Join([]string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
		"Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"}, ", ")
	allowMethods := strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodDelete, http.MethodOptions}, ", ")

	// 回写 Origin 而不是 Referer：浏览器校验的是 Access-Control-Allow-Origin 与
	// Origin 头完全相等，带深路径的完整 Referer URL 不是合法 origin，会把响应
	// 整个判为 CORS 失败（浏览器侧拦截，服务端日志看不到任何错误）。
	origin := ctx.Request.Header.Get("Origin")
	if origin != "" && strings.Contains(origin, "vancone.com") {
		ctx.Header("Access-Control-Allow-Origin", origin)
		ctx.Header("Access-Control-Allow-Headers", allowHeaders)
		ctx.Header("Access-Control-Allow-Methods", allowMethods)
		ctx.Header("Access-Control-Allow-Credentials", "true")
	}

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
	}
}
