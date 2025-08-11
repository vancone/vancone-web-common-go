package request

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPagination(ctx *gin.Context) (pageNo int, pageSize int) {
	var err error
	pageNo, err = strconv.Atoi(ctx.Query("pageNo"))
	if err != nil {
		log.Println("Failed to read pageNo", err)
		pageNo = 1
	}
	pageSize, err = strconv.Atoi(ctx.Query("pageSize"))
	if err != nil {
		log.Println("Failed to read pageSize", err)
		pageSize = 10
	}
	return
}
