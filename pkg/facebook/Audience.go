package facebook

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/imroc/req"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type createCustomAudienceForm struct {
	Name				string	`json:"name"`
	Subtype				string	`json:"subtype"`
	Description 		string	`json:"description"`
	CustomerFileSource	string	`json:"customer_file_source"`
}

type addContentForm struct {
	Session		session		`json:"session"`
	Payload		payload		`json:"payload"`
}

type session struct {
	SessionId	int			`json:"session_id"`
	Batch		int			`json:"batch_seq"`
	IsLast		bool		`json:"last_batch_flag"`
}

type payload struct {
	Schema		string		`json:"schema"`
	Data		[]string	`json:"data"`
}

type createAudienceResponse struct {
	Id		string			`json:"id,omitempty"`
	Error	*errorResponse	`json:"error,omitempty"`
}

type addDataResponse struct {
	Error	*errorResponse	`json:"error,omitempty"`
}

type audienceData struct {
	Id		string		`json:"id"`
}

type errorResponse struct {
	Message			string		`json:"message"`
	Type			string		`json:"type,omitempty"`
	Code			int			`json:"code,omitempty"`
	ErrorSubCode	int			`json:"error_subcode,omitempty"`
}

func AddContentToAudience(audience, accessToken, accountId string, data []string, isLast bool, batch int) error {
	id, err := strconv.Atoi(accountId)
	if err != nil {
		return err
	}

	form := &addContentForm{
		Session: session{
			SessionId: id,
			Batch:     batch,
			IsLast:    isLast,
		},
		Payload: payload{
			Schema: "EMAIL_SHA256",
			Data: data,
		},
	}

	res, err := req.Post(buildGraphUrl("/" + audience + "/users", map[string]string{
		"access_token": accessToken,
	}), req.BodyJSON(&form))
	if err != nil {
		return err
	}

	var response addDataResponse
	err = res.ToJSON(&response)
	if err != nil {
		return err
	}

	if response.Error != nil {
		return errors.New(response.Error.Message)
	}

	return err
}

func AddContentToAudienceFromFile(adAccountId, audienceId, accessToken, filePath string) error {
	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 按行读取，每次10000条处理
	var data []string
	buff := bufio.NewReader(file)
	batch := 1
	for {
		line, err := buff.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if len(bytes.TrimSpace(line)) != 0 {
					data = append(data, fmt.Sprintf("%x", sha256.Sum256(bytes.ToLower(bytes.TrimSpace(line)))))
				}
				break
			}

			return err
		}

		data = append(data, fmt.Sprintf("%x", sha256.Sum256(bytes.ToLower(bytes.TrimSpace(line)))))
		if len(data) == 10000 {
			// 保存，并清空
			err = AddContentToAudience(audienceId, accessToken, adAccountId, data, false, batch)
			if err != nil {
				return err
			}

			data = nil
			batch++
		}
	}

	if len(data) > 0 {
		err = AddContentToAudience(audienceId, accessToken, adAccountId, data, true, batch)
		if err != nil {
			return err
		}
	}

	return nil
}

// create custom audience from file
func CreateCustomAudienceByFile(file string, adAccountId string, accessToken string) error {
	fileName := filepath.Base(file)

	param := &createCustomAudienceForm{
		Name:               strings.TrimSuffix(fileName, filepath.Ext(fileName)),
		Subtype:            "CUSTOM",
		Description:        "",
		CustomerFileSource: "USER_PROVIDED_ONLY",
	}

	res, err := req.Post(buildGraphUrl("/act_" + adAccountId + "/customaudiences", map[string]string{
		"access_token": accessToken,
	}), req.BodyJSON(param))
	if err != nil {
		return err
	}

	var response createAudienceResponse
	err = res.ToJSON(&response)
	if err != nil {
		return err
	}

	if response.Error != nil {
		return errors.New(response.Error.Message)
	}
	return AddContentToAudienceFromFile(adAccountId, response.Id, accessToken, file)
}
