package services

import (
	"arvan-wallet-service/types/interfaces/repositories"
	"arvan-wallet-service/types/structs"
)

type IChargeCreditByCode interface {
	SetRepositories(userRepo repositories.IUserRepository, transRepo repositories.ITransactionRepository)
	SetService(service IRequestToUtilizeCode)
	Charge(mobile string, creditCode string) (uint, *structs.CustomError)
}
