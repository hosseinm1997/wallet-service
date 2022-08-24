package repositories

import (
	"arvan-wallet-service/db/models"
	"arvan-wallet-service/types/structs"
)

type IUserRepository interface {
	GetUser(mobile string) structs.RepositoryResult[models.User]
	GetUserMayHaveTransaction(mobile string, creditCode string) structs.RepositoryResult[models.User]
	CreateWithTransaction(mobile string, creditCode string) structs.RepositoryResult[models.User]
	AppendTransaction(user *models.User, creditCode string) structs.RepositoryResult[models.User]
	RunChargeBalanceDBFunc(user *models.User, amount uint, logId uint) error
}
