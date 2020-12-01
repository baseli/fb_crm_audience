package main

import (
	"github.com/baseli/fb_crm_audience/api"
	"github.com/baseli/fb_crm_audience/internal/server"
	"github.com/imroc/req"
	"net/http"
	"net/url"
)

func main() {
	// 读取用户根目录配置文件，然后配置代理
	proxyConfig, err := api.GetProxy()
	if err == nil {
		req.SetProxy(func(request *http.Request) (*url.URL, error) {
			return url.Parse(proxyConfig)
		})
	}

	server.NewServer()
}
