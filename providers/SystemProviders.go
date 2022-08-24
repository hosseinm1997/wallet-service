package providers

import (
	"arvan-wallet-service/http/middlewares"
	infra "arvan-wallet-service/infrastructures"
	infraInterfaces "arvan-wallet-service/infrastructures/interfaces"
	"net/http"
	//"arwan-wallet-service/interfaces"
)

func SystemProvides() {
	infra.Register[infraInterfaces.IChiRouter](
		func(params ...any) infraInterfaces.IChiRouter {
			return infra.ChiRouter()
		},
	)

	infra.Register[infraInterfaces.IResponseFormatter](
		func(params ...any) infraInterfaces.IResponseFormatter {
			return middlewares.GetResponseFormatter(params[0].(*http.Request))
		},
	)
}
