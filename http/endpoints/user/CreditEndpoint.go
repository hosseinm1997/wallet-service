package user

import (
	"arvan-wallet-service/infrastructures"
	"arvan-wallet-service/types/interfaces/repositories"
	"arvan-wallet-service/types/interfaces/services"
	. "arvan-wallet-service/utils"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

type CreditEndpoint struct{}

func (r CreditEndpoint) Charge(res http.ResponseWriter, req *http.Request) {

	mobile := chi.URLParam(req, "mobile")
	creditCode := chi.URLParam(req, "code")

	// todo: mobile and code validation

	service := infrastructures.Resolve[services.IChargeCreditByCode]()

	amount, err := service.Charge(mobile, creditCode)

	if err != nil {
		Respond(req).WithError(err)
		return
	}

	Respond(req).WithOkResultHavingMessage(
		map[string]any{"charged_amount": amount},
		fmt.Sprintf("Your balance has been charged by %d R", amount),
	)
}

func (r CreditEndpoint) Balance(res http.ResponseWriter, req *http.Request) {
	mobile := chi.URLParam(req, "mobile")
	repo := infrastructures.Resolve[repositories.IUserRepository]()
	result := repo.GetUser(mobile)

	if result.RowsAffected == 0 {
		Respond(req).WithBusinessLogicExceptionResult(fmt.Errorf("user not found"))
		return
	}

	Respond(req).WithOkResult(map[string]any{
		"balance": result.Model.Balance,
	})
}
