package repositories

import (
	"arvan-wallet-service/db/models"
	"arvan-wallet-service/types/structs"
)

type ITransactionRepository interface {
	UpdateTransactionAsLimited(id uint) structs.RepositoryResult[models.Transaction]
	GetSuccessfulTransactions(mobile string, creditCode string) structs.RepositoryResult[models.Transaction]
}
