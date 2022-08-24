package repositories

import (
	"arvan-wallet-service/db/models"
	"arvan-wallet-service/types/constants"
	"arvan-wallet-service/types/structs"
	. "arvan-wallet-service/utils"
)

type TransactionRepository struct{}

func (r *TransactionRepository) UpdateTransactionAsLimited(id uint) structs.RepositoryResult[models.Transaction] {
	t := &models.Transaction{ID: id, Status: constants.TransactionStatusEnums.Requested}
	response := DB().Model(t).Select("*").Update("status", constants.TransactionStatusEnums.FailedForCreditCodeLimitations)
	return RepoResult[models.Transaction](t, response)
}

func (r *TransactionRepository) GetSuccessfulTransactions(mobile string, creditCode string) structs.RepositoryResult[models.Transaction] {
	t := []*models.Transaction{}
	response := DB().
		Where(&models.Transaction{
			CreditCodeText: creditCode,
			Mobile:         mobile,
			Status:         constants.TransactionStatusEnums.Successful,
		}).
		Find(&t)
	return RepoResult[models.Transaction](t, response)
}
