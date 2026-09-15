package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCrossOriginMiddleware_AllowedOrigin(t *testing.T) {
	// 初始化Gin测试引擎
	r := gin.Default()
	r.Use(CrossOriginMiddleware)
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// 创建测试请求 - 允许的来源
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Referer", "https://www.vancone.com/path/")
	w := httptest.NewRecorder()

	// 发送请求
	r.ServeHTTP(w, req)

	// 检查响应头
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "https://www.vancone.com/path", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, PATCH, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestCrossOriginMiddleware_AllowedOriginWithoutTrailingSlash(t *testing.T) {
	r := gin.Default()
	r.Use(CrossOriginMiddleware)
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// 来源没有尾部斜杠
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Referer", "https://sub.vancone.com")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "https://sub.vancone.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCrossOriginMiddleware_NotAllowedOrigin(t *testing.T) {
	r := gin.Default()
	r.Use(CrossOriginMiddleware)
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// 不允许的来源
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Referer", "https://example.com")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// 不应设置跨域头
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Methods"))
}

func TestCrossOriginMiddleware_OptionsMethod_Allowed(t *testing.T) {
	r := gin.Default()
	r.Use(CrossOriginMiddleware)
	r.OPTIONS("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// OPTIONS请求，允许的来源
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Referer", "https://www.vancone.com")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// 中间件应该直接返回204 No Content
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "https://www.vancone.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCrossOriginMiddleware_OptionsMethod_NotAllowed(t *testing.T) {
	r := gin.Default()
	r.Use(CrossOriginMiddleware)
	r.OPTIONS("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// OPTIONS请求，不允许的来源
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Referer", "https://example.com")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// 即使来源不允许，OPTIONS请求也应返回204
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCrossOriginMiddleware_NoReferer(t *testing.T) {
	r := gin.Default()
	r.Use(CrossOriginMiddleware)
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// 没有Referer头的请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}
