package providers

import (
	"arvan-wallet-service/infrastructures"
	"arvan-wallet-service/services"
	"arvan-wallet-service/types/interfaces/repositories"
	serviceInterfaces "arvan-wallet-service/types/interfaces/services"
)

func ServiceProvider() {

	infrastructures.Register[serviceInterfaces.IRequestToInquiryCode](
		func(params ...any) serviceInterfaces.IRequestToInquiryCode {
			return &services.RequestToInquiryCode{}
		},
	)

	infrastructures.Register[serviceInterfaces.IRequestToUtilizeCode](
		func(params ...any) serviceInterfaces.IRequestToUtilizeCode {
			return &services.RequestToUtilizeCode{}
		},
	)

	infrastructures.Register[serviceInterfaces.IChargeCreditByCode](
		func(params ...any) serviceInterfaces.IChargeCreditByCode {
			s := &services.CircuitBreakerProxyChargeCreditByCode{}
			s.SetRepositories(
				infrastructures.Resolve[repositories.IUserRepository](),
				infrastructures.Resolve[repositories.ITransactionRepository](),
			)
			s.SetService(infrastructures.Resolve[serviceInterfaces.IRequestToUtilizeCode]())
			return s
		},
	)

}
