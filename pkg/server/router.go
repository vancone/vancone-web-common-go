package server

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vancone/vancone-web-common-go/pkg/config"
	"github.com/vancone/vancone-web-common-go/pkg/logger"
	"go.uber.org/zap"
)

type Router interface {
	Load(engine *gin.Engine)
}

func StartServer(g *gin.Engine) {
	logger.Infof("Server started, listening on port %d", config.Server.Port)
	g.Use(GinZapLogger(logger.SugaredLogger.Desugar()), gin.RecoveryWithWriter(gin.DefaultErrorWriter))
	err := g.Run(fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		log.Panicln("Failed to start server:", err)
	}
}

// GinZapLogger 将 Gin 的访问日志通过 Zap 记录
func GinZapLogger(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 请求结束，记录日志
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		proto := c.Request.Proto
		userAgent := c.Request.UserAgent()

		// 构造日志字段
		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.Duration("latency", latency),
			zap.String("proto", proto),
		}

		// 根据状态码决定日志级别
		if statusCode >= 500 {
			zapLogger.Error("服务器错误", fields...)
		} else if statusCode >= 400 {
			zapLogger.Warn("客户端错误或未找到", fields...)
		} else {
			zapLogger.Info("HTTP 请求", fields...)
		}
	}
}
