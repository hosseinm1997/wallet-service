package interfaces

type IResponseFormatter interface {
	WithOkResult(data map[string]any)
	WithOkResultHavingMessage(data map[string]any, message string)
	WithError(err error)
	WithBusinessLogicExceptionResult(err error)
	WithInternalErrorResult(err error)
	WithUnAuthorizedError(err error)
}
