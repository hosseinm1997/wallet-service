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

type RequestToUtilizeCode struct{}

func (r *RequestToUtilizeCode) Send(codeText string, referenceId uint) structs.UtilizeCodeResponse {
	data := map[string]any{
		"code_text":    codeText,
		"reference_id": referenceId,
	}

	var jsonData []byte
	jsonData, _ = json.Marshal(data)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, r.getUrl(codeText, referenceId), bytes.NewBuffer(jsonData))
	req.Close = true
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Client-token", viper.GetString("MS_CREDIT_SECRET"))
	res, err := client.Do(req)
	if err != nil {
		return structs.UtilizeCodeResponse{
			LogId:  0,
			Amount: 0,
			Err:    CustomError(structs.Categories.Internal, err.Error()),
		}
	}

	defer res.Body.Close()

	var response map[string]any
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return structs.UtilizeCodeResponse{
			LogId:  0,
			Amount: 0,
			Err:    CustomError(structs.Categories.Internal, err.Error()),
		}
	}

	switch res.StatusCode {
	case http.StatusOK:
		if !IsInMap[map[string]any](response, "data") ||
			!IsInMap[float64](response["data"].(map[string]any), "amount") ||
			!IsInMap[float64](response["data"].(map[string]any), "log_id") {

			return structs.UtilizeCodeResponse{
				LogId:  0,
				Amount: 0,
				Err:    CustomError(structs.Categories.Internal, "unexpected response"),
			}
		}

		return structs.UtilizeCodeResponse{
			LogId:  uint(response["data"].(map[string]any)["log_id"].(float64)),
			Amount: uint(response["data"].(map[string]any)["amount"].(float64)),
			Err:    nil,
		}

	case http.StatusUnprocessableEntity:

		if !IsInMap[string](response, "message") {
			return structs.UtilizeCodeResponse{
				LogId:  0,
				Amount: 0,
				Err:    CustomError(structs.Categories.Internal, "unexpected response"),
			}
		}

		return structs.UtilizeCodeResponse{
			LogId:  0,
			Amount: 0,
			Err:    CustomError(structs.Categories.BusinessLogic, response["message"].(string)),
		}

	case http.StatusInternalServerError:
		fallthrough
	default:
		return structs.UtilizeCodeResponse{
			LogId:  0,
			Amount: 0,
			Err:    CustomError(structs.Categories.Internal, "unexpected response"),
		}
	}
}

func (r *RequestToUtilizeCode) getUrl(code string, referenceId uint) string {
	return fmt.Sprintf("%s/credit/code/%s/%d",
		viper.GetString("MS_CREDIT_URL"),
		code,
		referenceId,
	)
}
