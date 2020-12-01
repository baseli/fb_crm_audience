package server

import (
	"errors"
	"github.com/baseli/fb_crm_audience/api"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

type response struct {
	Code	int			`json:"code"`
	Data	interface{}	`json:"data"`
	Size	interface{}	`json:"size,omitempty"`
	Message string		`json:"message"`
}

func buildResponse(data interface{}, size int64, err error) *response {
	if err != nil {
		return &response{
			Code:    1,
			Data:    nil,
			Message: err.Error(),
		}
	}

	return &response{
		Code:    0,
		Data:    data,
		Size:    size,
		Message: "",
	}
}

func NewServer() {
	exPath := "."
	ex, err := os.Executable()
	if err == nil {
		exPath = filepath.Dir(ex)
	}

	r := gin.Default()
	r.Use(corsMiddleware())

	// 获取fb账号
	r.GET("/accounts", func(context *gin.Context) {
		accounts, err := api.GetAccounts()
		context.JSON(200, buildResponse(accounts, 0, err))
	})

	// 获取广告账号
	r.GET("/adAccounts", func(context *gin.Context) {
		context.JSON(200, buildResponse(api.GetAdAccounts(
			context.DefaultQuery("name", ""),
			context.DefaultQuery("id", ""),
			context.DefaultQuery("page", "1"),
			context.DefaultQuery("size", "100"),
		)))
	})

	// 获取任务列表
	r.GET("/task", func(context *gin.Context) {
		context.JSON(200, buildResponse(api.GetTask(
			context.DefaultQuery("status", ""),
			context.DefaultQuery("id", ""),
			context.DefaultQuery("page", ""),
			context.DefaultQuery("size", ""),
		)))
	})

	// 新建任务
	r.POST("/task", func(context *gin.Context) {
		var params struct{
			Files 		[]string 	`json:"files"`
			AdAccounts	[]string	`json:"adAccounts"`
		}

		err := context.BindJSON(&params)
		if err != nil {
			context.JSON(200, buildResponse(nil, 0, errors.New("参数错误")))
			return
		}

		context.JSON(200, buildResponse(nil, 0, api.CreateTask(params.Files, params.AdAccounts)))
	})

	// 同步广告账号
	r.POST("/bind", func(context *gin.Context) {
		go api.SyncAdAccount(context.Query("token"))

		context.JSON(200, response{
			Code:    0,
			Data:    nil,
			Message: "",
		})
	})

	// 清空任务
	r.DELETE("/task", func(context *gin.Context) {
		err := api.RemoveAll()

		context.JSON(200, buildResponse(nil, 0, err))
	})

	// 重试
	r.POST("/task/retry", func(context *gin.Context) {
		err := api.Retry()

		context.JSON(200, buildResponse(nil, 0, err))
	})

	// 获取代理信息
	r.GET("/proxy", func(context *gin.Context) {
		proxy, err := api.GetProxy()
		if err != nil {
			context.JSON(200, response{
				Code:    1,
				Data:    nil,
				Size:    nil,
				Message: err.Error(),
			})
			return
		}

		context.JSON(200, response{
			Code:    0,
			Data:    proxy,
			Size:    nil,
			Message: "",
		})
	})

	// 设置代理
	r.POST("/proxy", func(context *gin.Context) {
		err := api.StoreProxy(context.Query("proxy"))
		if err != nil {
			context.JSON(200, response{
				Code:    1,
				Data:    nil,
				Size:    nil,
				Message: err.Error(),
			})
			return
		}

		context.JSON(200, response{
			Code:    0,
			Data:    nil,
			Size:    nil,
			Message: "",
		})
	})

	// 静态资源
	r.Static("/web", exPath + "/web")

	_ = r.Run(":19597")
}
