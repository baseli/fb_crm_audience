package api

import (
	"github.com/baseli/fb_crm_audience/api/schema"
	"github.com/baseli/fb_crm_audience/api/util"
	"github.com/baseli/fb_crm_audience/pkg/facebook"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

type adAccountList struct {
	schema.AdAccount
	AccountName string	`gorm:"-"`
}

func GetAdAccounts(adAccountName, adAccountId, page, size string) ([]adAccountList, int64, error) {
	db, err := util.NewDatabase()
	if err != nil {
		return nil, 0, err
	}

	where := db.Model(&schema.AdAccount{})

	if adAccountName != "" {
		where.Where("name like ?", "%" + adAccountName + "%")
	}

	if adAccountId != "" {
		where.Where("ad_account_id = ?", adAccountId)
	}

	limit, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		limit = 100
	}

	offset, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		offset = 0
	} else {
		offset = (offset - 1) * limit
	}

	var allSize int64
	err = where.Count(&allSize).Error
	if err != nil {
		return nil, 0, err
	}

	where.Limit(int(limit)).Offset(int(offset))

	var result []adAccountList
	err = where.Find(&result).Error

	if err != nil {
		return nil, 0, err
	}

	// 关联数据处理
	var accounts []schema.Account
	err = db.Model(&schema.Account{}).Find(&accounts).Error
	if err != nil {
		return nil, 0, err
	}

	nameMap := make(map[string]string)
	for _, account := range accounts {
		nameMap[account.AccountId] = account.Name
	}

	for index, adAccount := range result {
		result[index].AccountName = nameMap[adAccount.AccountId]
	}

	return result, allSize, nil
}

func SyncAdAccount(accessToken string) {
	log, err := util.NewLog()
	if err != nil {
		return
	}
	defer log.FileHandler.Close()

	db, err := util.NewDatabase()
	if err != nil {
		log.Error().Msg("Error: get database error, " + err.Error())
		return
	}

	log.Debug().Msg("Start sync ad account.")
	facebookAccount, err := facebook.GetAccount(accessToken)
	if err != nil {
		log.Error().Msg("Sync error: " + err.Error())
		return
	}

	var facebookAccountExists schema.Account
	err = db.Where("account_id = ?", facebookAccount.Id).First(&facebookAccountExists).Error
	if err != nil {
		facebookAccountExists = schema.Account{
			AccountId:    facebookAccount.Id,
			Email:        facebookAccount.Email,
			Token:        facebookAccount.AccessToken,
			Name:         facebookAccount.Name,
			AuthTime:     time.Now(),
			NextAuthTime: time.Now().Add(time.Duration(rand.Int31n(facebookAccount.ExpiresIn)) * time.Second),
		}
		db.Create(&facebookAccountExists)
	} else {
		facebookAccountExists.AccountId = facebookAccount.Id
		facebookAccountExists.Email = facebookAccount.Email
		facebookAccountExists.Token = facebookAccount.AccessToken
		facebookAccountExists.Name = facebookAccount.Name
		facebookAccountExists.AuthTime = time.Now()
		facebookAccountExists.NextAuthTime = time.Now().Add(time.Duration(rand.Int31n(facebookAccount.ExpiresIn)) * time.Second)

		db.Where("account_id", facebookAccount.Id).Updates(&facebookAccountExists)
	}

	adAccount, err := facebook.GetAdAccounts(facebookAccount.Id, accessToken)
	if err != nil {
		log.Error().Msg("Error: get ad accounts error, " + err.Error())
		return
	}
	saveData(adAccount.Data, db, facebookAccount.Id)

	for adAccount.Pager.Next != "" {
		log.Debug().Msg("Sync next page ad account: " + adAccount.Pager.Next)
		adAccount, err = adAccount.GetNext()
		if err != nil {
			log.Error().Msg("Error: get ad accounts error, " + err.Error())
			break
		}

		if len(adAccount.Data) == 0 {
			break
		}

		saveData(adAccount.Data, db, facebookAccount.Id)
	}

	log.Debug().Msg("Finish sync ad account.")
}

func saveData(adAccounts []facebook.AdAccount, db *gorm.DB, accountId string) {
	var adAccountIds []string
	var data []schema.AdAccount

	for _, account := range adAccounts {
		adAccountIds = append(adAccountIds, account.AccountId)
		data = append(data, schema.AdAccount{
			AdAccountId: account.AccountId,
			AccountId:   accountId,
			Name:        account.Name,
		})
	}

	db.Transaction(func(tx *gorm.DB) error {
		if len(adAccountIds) > 0 && len(data) > 0 {
			tx.Where("ad_account_id in (?)", adAccountIds).Where("account_id", accountId).Delete(schema.AdAccount{})
			tx.Create(&data)
		}

		return nil
	})
}
