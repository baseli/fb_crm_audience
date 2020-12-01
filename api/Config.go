package api

import (
	"encoding/json"
	"github.com/baseli/fb_crm_audience/api/util"
	"io/ioutil"
	"path"
)

type config struct {
	Proxy	string	`json:"proxy"`
}

func StoreProxy(proxy string) error {
	filepath, err := getConfigPath()
	if err != nil {
		return err
	}

	proxyConfig := config{Proxy: proxy}
	content, err := json.Marshal(proxyConfig)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, content, 0700)
}

func GetProxy() (string, error) {
	filepath, err := getConfigPath()
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	var proxyConfig config
	err = json.Unmarshal(content, &proxyConfig)
	if err != nil {
		return "", err
	}

	return proxyConfig.Proxy, nil
}

func getConfigPath() (string, error) {
	home, err := util.GetHomePath()
	if err != nil {
		return "", err
	}

	return path.Join(home, "config.json"), nil
}
