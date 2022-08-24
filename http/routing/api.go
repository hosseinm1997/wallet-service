package routing

import (
	"arvan-wallet-service/infrastructures"
	adminEndpointInterfaces "arvan-wallet-service/types/interfaces/endpoints/admin"
	userEndpointInterfaces "arvan-wallet-service/types/interfaces/endpoints/user"
	"github.com/go-chi/chi"
)

func Routes(r *chi.Mux) {
	userCreditEndpoint := infrastructures.Resolve[userEndpointInterfaces.ICreditEndpoint]()
	adminCreditEndpoint := infrastructures.Resolve[adminEndpointInterfaces.ICreditEndpoint]()

	r.Post("/user/{mobile}/credit/code/{code}", userCreditEndpoint.Charge)
	r.Get("/user/{mobile}/credit", userCreditEndpoint.Balance)
	r.Get("/admin/credit/{codeText}/list", adminCreditEndpoint.List)
}
