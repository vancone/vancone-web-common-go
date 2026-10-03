package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newCrossOriginRouter() *gin.Engine {
	r := gin.Default()
	r.Use(CrossOriginMiddleware)
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.OPTIONS("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func TestCrossOriginMiddleware_AllowedOrigin(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://www.vancone.com")
	w := httptest.NewRecorder()

	newCrossOriginRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "https://www.vancone.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, PATCH, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
}

// 回归：带深路径 Referer 的请求，ACAO 必须回写 Origin 本身，而不是把完整
// Referer URL 写进 Access-Control-Allow-Origin（那不是合法 origin，浏览器会
// 拦截响应，表现为接口“概率性失败”）。
func TestCrossOriginMiddleware_DeepPathRefererDoesNotLeakIntoAllowOrigin(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://console.vancone.com")
	req.Header.Set("Referer", "https://console.vancone.com/console/apps")
	w := httptest.NewRecorder()

	newCrossOriginRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "https://console.vancone.com", w.Header().Get("Access-Control-Allow-Origin"))
}

// 只有 Referer 没有 Origin 的请求（非 CORS 场景）不应回写跨域头。
func TestCrossOriginMiddleware_RefererOnlySetsNothing(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Referer", "https://www.vancone.com/path/")
	w := httptest.NewRecorder()

	newCrossOriginRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCrossOriginMiddleware_NotAllowedOrigin(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	newCrossOriginRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// 不应设置跨域头
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Methods"))
}

func TestCrossOriginMiddleware_OptionsMethod_Allowed(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://www.vancone.com")
	w := httptest.NewRecorder()

	newCrossOriginRouter().ServeHTTP(w, req)

	// 中间件应该直接返回204 No Content
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "https://www.vancone.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCrossOriginMiddleware_OptionsMethod_NotAllowed(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	newCrossOriginRouter().ServeHTTP(w, req)

	// 即使来源不允许，OPTIONS请求也应返回204
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCrossOriginMiddleware_NoOrigin(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	newCrossOriginRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}
