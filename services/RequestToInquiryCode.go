package services

import (
	"arvan-wallet-service/types/structs"
	. "arvan-wallet-service/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

type RequestToInquiryCode struct{}

func (r RequestToInquiryCode) Send(codeText string) *structs.CustomError {

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, r.getUrl(codeText), bytes.NewBuffer(make([]byte, 0)))
	req.Close = true
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Client-token", viper.GetString("MS_CREDIT_SECRET"))
	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		return CustomError(structs.Categories.Internal, err.Error())
	}

	var response map[string]any
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return CustomError(structs.Categories.Internal, err.Error())
	}

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnprocessableEntity:
		return CustomError(structs.Categories.BusinessLogic, response["message"].(string))
	case http.StatusInternalServerError:
		fallthrough
	default:
		return CustomError(structs.Categories.Internal, "unexpected response")
	}
}

func (r RequestToInquiryCode) getUrl(code string) string {
	return fmt.Sprintf("%s/credit/code/%s/inquiry",
		viper.GetString("MS_CREDIT_URL"),
		code,
	)
}
