package api

import (
	"errors"
	"github.com/baseli/fb_crm_audience/api/schema"
	"github.com/baseli/fb_crm_audience/api/util"
	"github.com/baseli/fb_crm_audience/pkg/facebook"
	"strconv"
	"sync"
	"time"
)

type resultType struct {
	schema.Task
	AdAccountName 	string	`gorm:"-"`
}

type resultChan struct {
	Id		uint
	Err		error
}

func CreateTask(files, adAccounts []string) error {
	db, err := util.NewDatabase()
	if err != nil {
		return err
	}

	var list []schema.Task
	for _, file := range files {
		for _, adAccount := range adAccounts {
			list = append(list, schema.Task{
				AdAccountId: adAccount,
				File:        file,
				Msg:         "",
				Status:      0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
		}
	}

	if len(list) > 0 {
		err = db.Create(&list).Error
		if err != nil {
			return err
		}

		go uploadAudience()
		return nil
	}

	return nil
}

func GetTask(status, id, page, size string) ([]resultType, int64, error) {
	db, err := util.NewDatabase()
	if err != nil {
		return nil, 0, err
	}

	where := db.Model(&schema.Task{}).Order("id desc")
	if status != "" {
		where.Where("status = ?", status)
	}

	if id != "" {
		where.Where("ad_account_id = ?", id)
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

	var result []resultType
	err = where.Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	if len(result) > 0 {
		// 关联数据处理
		var accounts []schema.AdAccount
		var adAccountIds []string
		for _, item := range result {
			adAccountIds = append(adAccountIds, item.AdAccountId)
		}

		err = db.Where("ad_account_id in ?", adAccountIds).Find(&accounts).Error
		if err != nil {
			return nil, 0, err
		}

		nameMap := make(map[string]string)
		for _, account := range accounts {
			nameMap[account.AdAccountId] = account.Name
		}

		for i := range result {
			result[i].AdAccountName = nameMap[result[i].AdAccountId]
		}
	}

	return result, allSize, nil
}

func RemoveAll() error {
	db, err := util.NewDatabase()
	if err != nil {
		return err
	}

	return db.Model(&schema.Task{}).Where("1 = 1").Delete(schema.Task{}).Error
}

func Retry() error {
	db, err := util.NewDatabase()
	if err != nil {
		return err
	}

	err = db.Model(&schema.Task{}).Where("status = ?", 2).Update("status", 0).Error
	if err != nil {
		return err
	}

	go uploadAudience()
	return nil
}

func uploadAudience() {
	db, err := util.NewDatabase()
	if err != nil {
		return
	}

	var list []schema.Task
	err = db.Where("status = ?", 0).Order("file asc, ad_account_id asc").Find(&list).Error
	if err != nil {
		return
	}

	var taskId []uint
	wg := sync.WaitGroup{}
	ret := make(chan resultChan, len(list))

	for index, item := range list {
		taskId = append(taskId, item.Id)
		if index + 1 < len(list) {
			if item.File != list[index + 1].File {
				wg.Add(1)

				go reallyUpload(taskId, item.File, ret, &wg)
				taskId = nil
			}
		} else {
			wg.Add(1)

			go reallyUpload(taskId, item.File, ret, &wg)
			taskId = nil
		}
	}

	wg.Wait()
	close(ret)
	for result := range ret {
		if result.Err != nil {
			db.Exec("update tasks set status = 2, msg = ? where id = ?", result.Err.Error(), result.Id)
		} else {
			db.Exec("update tasks set status = 1 where id = ?", result.Id)
		}
	}
}

func reallyUpload(taskId []uint, file string, chanResult chan resultChan, wg *sync.WaitGroup) {
	defer wg.Done()
	db, err := util.NewDatabase()
	if err != nil {
		return
	}

	type token struct {
		Token		string
		AdAccountId	string
	}

	for _, task := range taskId {
		var tokenValue token
		err = db.Raw("select accounts.token, tasks.ad_account_id from tasks " +
			"left join ad_accounts on tasks.ad_account_id = ad_accounts.ad_account_id " +
			"left join accounts on ad_accounts.account_id = accounts.account_id " +
			"where tasks.id = ? limit 1", task).First(&tokenValue).Error

		if err != nil {
			chanResult <- resultChan{
				Id:  task,
				Err: errors.New("token failed"),
			}
			continue
		}

		err = facebook.CreateCustomAudienceByFile(file, tokenValue.AdAccountId, tokenValue.Token)
		if err != nil {
			chanResult <- resultChan{
				Id:  task,
				Err: err,
			}
			continue
		}

		chanResult <- resultChan{
			Id:  task,
			Err: nil,
		}
	}
}
