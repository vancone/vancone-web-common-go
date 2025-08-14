package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetPagination_NormalParams(t *testing.T) {
	// 初始化Gin引擎
	gin.SetMode(gin.TestMode)

	// 创建测试请求
	req := httptest.NewRequest(http.MethodGet, "/?pageNo=2&pageSize=20", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// 调用函数
	pageNo, pageSize := GetPagination(ctx)

	// 验证结果
	assert.Equal(t, 2, pageNo)
	assert.Equal(t, 20, pageSize)
}

func TestGetPagination_DefaultWhenPageNoMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodGet, "/?pageSize=30", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	pageNo, pageSize := GetPagination(ctx)

	assert.Equal(t, 1, pageNo)    // 默认值
	assert.Equal(t, 30, pageSize) // 传入的值
}

func TestGetPagination_DefaultWhenPageSizeMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodGet, "/?pageNo=5", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	pageNo, pageSize := GetPagination(ctx)

	assert.Equal(t, 5, pageNo)    // 传入的值
	assert.Equal(t, 10, pageSize) // 默认值
}

func TestGetPagination_DefaultWhenBothMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	pageNo, pageSize := GetPagination(ctx)

	assert.Equal(t, 1, pageNo)    // 默认值
	assert.Equal(t, 10, pageSize) // 默认值
}

func TestGetPagination_InvalidPageNo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodGet, "/?pageNo=abc&pageSize=20", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	pageNo, pageSize := GetPagination(ctx)

	assert.Equal(t, 1, pageNo)    // 无效时使用默认值
	assert.Equal(t, 20, pageSize) // 有效参数
}

func TestGetPagination_InvalidPageSize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodGet, "/?pageNo=3&pageSize=xyz", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	pageNo, pageSize := GetPagination(ctx)

	assert.Equal(t, 3, pageNo)    // 有效参数
	assert.Equal(t, 10, pageSize) // 无效时使用默认值
}

func TestGetPagination_NegativeValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodGet, "/?pageNo=-1&pageSize=-5", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	pageNo, pageSize := GetPagination(ctx)

	// 注意：原函数对负数不做处理，会直接返回转换后的值
	assert.Equal(t, -1, pageNo)
	assert.Equal(t, -5, pageSize)
}
