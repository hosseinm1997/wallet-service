package providers

import (
	adminEndpoint "arvan-wallet-service/http/endpoints/admin"
	userEndpoint "arvan-wallet-service/http/endpoints/user"
	"arvan-wallet-service/infrastructures"
	adminEndpointInterface "arvan-wallet-service/types/interfaces/endpoints/admin"
	userEndpointInterface "arvan-wallet-service/types/interfaces/endpoints/user"
)

func EndpointProviders() {
	infrastructures.Register[userEndpointInterface.ICreditEndpoint](
		func(params ...any) userEndpointInterface.ICreditEndpoint {
			return userEndpoint.CreditEndpoint{}
		},
	)

	infrastructures.Register[adminEndpointInterface.ICreditEndpoint](
		func(params ...any) adminEndpointInterface.ICreditEndpoint {
			return adminEndpoint.CreditEndpoint{}
		},
	)
}
