package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vancone/vancone-web-common-go/pkg/logger"
)

func GetPagination(ctx *gin.Context) (pageNo int, pageSize int) {
	var err error
	pageNo, err = strconv.Atoi(ctx.Query("pageNo"))
	if err != nil {
		logger.Errorf("Failed to read pageNo: %v", err)
		pageNo = 1
	}
	pageSize, err = strconv.Atoi(ctx.Query("pageSize"))
	if err != nil {
		logger.Errorf("Failed to read pageSize: %v", err)
		pageSize = 10
	}
	return
}
