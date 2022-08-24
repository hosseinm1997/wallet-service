package admin

import "net/http"

type ICreditEndpoint interface {
	List(res http.ResponseWriter, req *http.Request)
}
