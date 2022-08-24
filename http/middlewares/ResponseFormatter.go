package middlewares

import (
	infrastructureInterfaces "arvan-wallet-service/infrastructures/interfaces"
	"arvan-wallet-service/types/structs"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseBag struct {
	HTTPCode int            `json:"-"`
	Result   bool           `json:"result"`
	Message  string         `json:"message"`
	Data     map[string]any `json:"data,omitempty"`
}

type ResFormat struct {
	req *http.Request
}

func (r *ResFormat) WithError(err error) {
	customErr, ok := err.(*structs.CustomError)

	if !ok {
		r.WithInternalErrorResult(err)
		return
	}

	switch customErr.Category {
	case structs.Categories.BusinessLogic:
		r.WithBusinessLogicExceptionResult(err)
	case structs.Categories.UnAuthorized:
		r.WithUnAuthorizedError(err)
	case structs.Categories.Internal:
		fallthrough
	default:
		r.WithInternalErrorResult(err)
	}

}

func (r *ResFormat) WithOkResult(data map[string]any) {
	bag := ResponseBag{
		HTTPCode: http.StatusOK,
		Result:   true,
		Message:  "ok",
		Data:     data,
	}

	r.setBag(bag)
}

func (r *ResFormat) WithBusinessLogicExceptionResult(err error) {
	bag := ResponseBag{
		HTTPCode: http.StatusUnprocessableEntity,
		Result:   false,
		Message:  err.Error(),
		Data:     nil,
	}

	r.setBag(bag)
}

func (r *ResFormat) WithInternalErrorResult(err error) {
	bag := ResponseBag{
		HTTPCode: http.StatusInternalServerError,
		Result:   false,
		Message:  "Internal server error",
		Data:     nil,
	}

	// todo: use logger service

	r.setBag(bag)

	fmt.Println(err)
}

func (r *ResFormat) WithOkResultHavingMessage(data map[string]any, message string) {
	bag := ResponseBag{
		HTTPCode: http.StatusOK,
		Result:   true,
		Message:  message,
		Data:     data,
	}

	r.setBag(bag)
}

func (r ResFormat) WithUnAuthorizedError(err error) {
	bag := ResponseBag{
		HTTPCode: http.StatusUnauthorized,
		Result:   false,
		Message:  err.Error(),
		Data:     nil,
	}

	r.setBag(bag)
}

func (r *ResFormat) setBag(bag ResponseBag) {
	*r.req = *r.req.WithContext(context.WithValue(r.req.Context(), "response", bag))
}

func ResponseFormatter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		finalWrite(w, r.Context().Value("response").(ResponseBag))
	}

	return http.HandlerFunc(fn)
}

func GetResponseFormatter(r *http.Request) infrastructureInterfaces.IResponseFormatter {
	rf := ResFormat{
		req: r,
	}
	return &rf
}

func finalWrite(w http.ResponseWriter, bag ResponseBag) {
	w.WriteHeader(bag.HTTPCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bag)
}
