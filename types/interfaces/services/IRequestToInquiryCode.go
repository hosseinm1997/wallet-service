package services

import "arvan-wallet-service/types/structs"

type IRequestToInquiryCode interface {
	Send(codeText string) *structs.CustomError
}
