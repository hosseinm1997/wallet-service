package repositories

import (
	"arvan-wallet-service/db/models"
	"arvan-wallet-service/types/constants"
	"arvan-wallet-service/types/structs"
	. "arvan-wallet-service/utils"
)

type UserRepository struct{}

func (u *UserRepository) GetUser(mobile string) structs.RepositoryResult[models.User] {
	model := &models.User{Mobile: mobile}
	response := DB().Take(model)
	return RepoResult[models.User](model, response)
}

func (u *UserRepository) GetUserMayHaveTransaction(mobile string, creditCode string) structs.RepositoryResult[models.User] {
	model := &models.User{Mobile: mobile}
	response := DB().Preload("Transactions", models.Transaction{CreditCodeText: creditCode}).Take(model)
	return RepoResult[models.User](model, response)
}

func (u *UserRepository) CreateWithTransaction(mobile string, creditCode string) structs.RepositoryResult[models.User] {
	transaction := &models.Transaction{
		CreditCodeText: creditCode,
		Status:         constants.TransactionStatusEnums.Requested,
	}
	user := &models.User{Mobile: mobile, Transactions: []*models.Transaction{transaction}}
	response := DB().Create(user)
	return RepoResult[models.User](user, response)
}

func (u *UserRepository) AppendTransaction(user *models.User, creditCode string) structs.RepositoryResult[models.User] {
	transaction := &models.Transaction{CreditCodeText: creditCode, Status: 1}
	err := DB().Model(&user).Association("Transactions").Append(transaction)
	if err != nil {
		return structs.RepositoryResult[models.User]{
			Data:         nil,
			Error:        err,
			RowsAffected: 0,
			Model:        nil,
		}
	}

	return structs.RepositoryResult[models.User]{
		Data:         nil,
		Error:        nil,
		RowsAffected: 1,
		Model:        user,
	}
}

func (u *UserRepository) RunChargeBalanceDBFunc(user *models.User, amount uint, logId uint) error {
	trans := user.Transactions[0]
	response := DB().Exec("select update_user_balance_v1_1($1, $2, $3, $4)",
		trans.Mobile,
		amount,
		trans.ID,
		logId,
	)
	return response.Error
}
