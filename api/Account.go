package api

import (
	"github.com/baseli/fb_crm_audience/api/schema"
	"github.com/baseli/fb_crm_audience/api/util"
)

func GetAccounts() ([]schema.Account, error) {
	db, err := util.NewDatabase()
	if err != nil {
		return nil, err
	}

	var result []schema.Account
	err = db.Raw("select * from accounts").Find(&result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}
