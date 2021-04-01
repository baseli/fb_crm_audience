package facebook

import (
	"github.com/imroc/req"
)

type AdAccount struct {
	Id				string	`json:"id"`
	Name			string	`json:"name"`
	AccountId		string	`json:"account_id"`
	AccountStatus	int		`json:"account_status"`
	Owner			string	`json:"owner,omitempty"`
}

type pager struct {
	Next		string		`json:"next,omitempty"`
	Previous	string		`json:"previous,omitempty"`
}

type adAccountPagination struct {
	Data	[]AdAccount	`json:"data"`
	Pager	pager		`json:"paging"`
}

// 获取fb广告账号
func GetAdAccounts(accountId string, accessToken string) (*adAccountPagination, error) {
	result := &adAccountPagination{
		Data:  []AdAccount{},
		Pager: pager{},
	}

	res, err := req.Get(buildGraphUrl("/" + accountId + "/adaccounts", map[string]string{
		"fields": "id,name,account_id,account_status",
		"access_token": accessToken,
		"limit": "200",
	}))
	if err != nil {
		return nil, err
	}

	err = res.ToJSON(&result)
	return result, err
}

// next page
func (pagination *adAccountPagination) GetNext() (*adAccountPagination, error) {
	if pagination.Pager.Next != "" {
		res, err := req.Get(pagination.Pager.Next)
		if err != nil {
			return nil, err
		}

		result := &adAccountPagination{
			Data:  []AdAccount{},
			Pager: pager{},
		}
		err = res.ToJSON(&result)
		return result, err
	}

	pagination.Data = []AdAccount{}

	return pagination, nil
}
