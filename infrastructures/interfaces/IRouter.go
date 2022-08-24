package interfaces

import "github.com/go-chi/chi"

type IChiRouter interface {
	InitRouter(func(*chi.Mux)) *chi.Mux
}
