package constants

const (
	requested = iota + 1
	successful
	failedForCreditCodeLimitations
)

var TransactionStatusEnums = struct {
	Requested                      uint
	Successful                     uint
	FailedForCreditCodeLimitations uint
}{
	requested,
	successful,
	failedForCreditCodeLimitations,
}
