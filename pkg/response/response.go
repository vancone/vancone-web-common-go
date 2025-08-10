package response

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type Page[T any] struct {
	List       []T   `json:"list"`
	PageNo     int64 `json:"pageNo"`
	PageSize   int64 `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
	TotalPage  int64 `json:"totalPage"`
}

func Paginate[T any](db *gorm.DB, pageNo int, pageSize int, result *Page[T]) error {
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Query totalCount count
	var totalCount int64
	if err := db.Count(&totalCount).Error; err != nil {
		return err
	}

	// Calculate totalCount pages
	totalPage := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	// Query data with pagination
	offset := (pageNo - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&result.List).Error; err != nil {
		return err
	}

	// Set response data
	result.TotalCount = totalCount
	result.TotalPage = int64(totalPage)
	result.PageSize = int64(pageSize)
	result.PageNo = int64(pageNo)

	return nil
}

func Success[T any](ctx *gin.Context, data T) {
	resp := Response[T]{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	ctx.JSON(http.StatusOK, resp)
}

func Fail[T any](ctx *gin.Context, code int, err error) {
	resp := Response[T]{
		Code:    code,
		Message: err.Error(),
	}
	ctx.JSON(http.StatusInternalServerError, resp)
}

func New[T any](ctx *gin.Context, data T, err error) {
	if err != nil {
		Fail[T](ctx, -1, err)
	} else {
		Success(ctx, data)
	}
}
