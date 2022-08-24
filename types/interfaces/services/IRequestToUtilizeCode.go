package services

import "arvan-wallet-service/types/structs"

type IRequestToUtilizeCode interface {
	Send(codeText string, referenceId uint) structs.UtilizeCodeResponse
}
