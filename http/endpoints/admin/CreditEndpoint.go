package admin

import (
	"arvan-wallet-service/infrastructures"
	"arvan-wallet-service/types/interfaces/repositories"
	. "arvan-wallet-service/utils"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

type CreditEndpoint struct {
}

func (r CreditEndpoint) List(res http.ResponseWriter, req *http.Request) {
	mobile := chi.URLParam(req, "mobile")
	creditCode := chi.URLParam(req, "code")

	repo := infrastructures.Resolve[repositories.ITransactionRepository]()
	result := repo.GetSuccessfulTransactions(mobile, creditCode)

	if result.RowsAffected == 0 {
		Respond(req).WithBusinessLogicExceptionResult(fmt.Errorf("no successful transaction"))
		return
	}

	Respond(req).WithOkResult(map[string]any{
		"transactions": result.Data,
	})
}
