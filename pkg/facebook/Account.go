package facebook

import (
	"github.com/imroc/req"
)

type Account struct {
	AccessToken		string	`json:"access_token,omitempty"`
	Id				string	`json:"id"`
	Email			string	`json:"email"`
	Name			string	`json:"name"`
	ExpiresIn 		int32	`json:"expires_in,omitempty"`
}

type token struct {
	AccessToken 	string	`json:"access_token"`
	ExpiresIn 		int32	`json:"expires_in"`
}

func GetAccount(accessToken string) (*Account, error) {
	// 短期token换长期token
	res, err := req.Get(buildGraphUrl("/oauth/access_token", map[string]string{
		"grant_type": "fb_exchange_token",
		"client_id": appId,
		"client_secret": appSecret,
		"fb_exchange_token": accessToken,
	}))
	if err != nil {
		return nil, err
	}

	tokenResponse := new(token)
	err = res.ToJSON(&tokenResponse)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	res, err = req.Get(buildGraphUrl("/me", map[string]string{
		"fields": "id,name,email",
		"access_token": tokenResponse.AccessToken,
	}))
	if err != nil {
		return nil, err
	}

	account := new(Account)
	err = res.ToJSON(&account)
	if err != nil {
		return nil, err
	}

	account.AccessToken = tokenResponse.AccessToken
	account.ExpiresIn = tokenResponse.ExpiresIn
	return account, nil
}
