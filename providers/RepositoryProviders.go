package providers

import (
	"arvan-wallet-service/db/repositories"
	"arvan-wallet-service/infrastructures"
	repositoryInterfaces "arvan-wallet-service/types/interfaces/repositories"
)

func RepositoryProviders() {
	infrastructures.Register[repositoryInterfaces.IUserRepository](
		func(params ...any) repositoryInterfaces.IUserRepository {
			return &repositories.UserRepository{}
		},
	)

	infrastructures.Register[repositoryInterfaces.ITransactionRepository](
		func(params ...any) repositoryInterfaces.ITransactionRepository {
			return &repositories.TransactionRepository{}
		},
	)
}
