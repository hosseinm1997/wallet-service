package user

import "net/http"

type ICreditEndpoint interface {
	Charge(res http.ResponseWriter, req *http.Request)
	Balance(res http.ResponseWriter, req *http.Request)
}
