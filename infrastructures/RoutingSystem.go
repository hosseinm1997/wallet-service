package infrastructures

import (
	middlewares2 "arvan-wallet-service/http/middlewares"
	"arvan-wallet-service/infrastructures/interfaces"
	"github.com/go-chi/chi"
	"sync"
)

type Router struct{}

func (router *Router) InitRouter(routerFunc func(*chi.Mux)) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares2.HTTPReqPanicRecovery)
	r.Use(middlewares2.ResponseFormatter)
	routerFunc(r)
	return r
}

var (
	m          *Router
	routerOnce sync.Once
)

func ChiRouter() interfaces.IChiRouter {
	if m == nil {
		routerOnce.Do(func() {
			m = &Router{}
		})
	}
	return m
}
