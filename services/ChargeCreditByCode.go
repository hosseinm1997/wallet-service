package services

import (
	"arvan-wallet-service/types/constants"
	"arvan-wallet-service/types/exceptions"
	"arvan-wallet-service/types/interfaces/repositories"
	"arvan-wallet-service/types/interfaces/services"
	"arvan-wallet-service/types/structs"
	. "arvan-wallet-service/utils"
	"encoding/json"
	"errors"
	"github.com/jackc/pgconn"
)

type ChargeCreditByCode struct {
	userRepo  repositories.IUserRepository
	transRepo repositories.ITransactionRepository
	service   services.IRequestToUtilizeCode
}

func (c *ChargeCreditByCode) SetRepositories(userRepo repositories.IUserRepository, transRepo repositories.ITransactionRepository) {
	c.userRepo = userRepo
	c.transRepo = transRepo
}

func (c *ChargeCreditByCode) SetService(service services.IRequestToUtilizeCode) {
	c.service = service
}

func (c *ChargeCreditByCode) Charge(mobile string, creditCode string) (uint, *structs.CustomError) {
	if c.userRepo == nil {
		return 0, CustomError(structs.Categories.Internal, "no repository set")
	}

	rr := c.userRepo.GetUserMayHaveTransaction(mobile, creditCode)

	if rr.RowsAffected == 0 {
		rr = c.userRepo.CreateWithTransaction(mobile, creditCode)
		if rr.Error != nil {
			return 0, CustomError(structs.Categories.Internal, "unable to create user with error: %s", rr.Error.Error())
		}
	} else if len(rr.Model.Transactions) == 0 {
		rr = c.userRepo.AppendTransaction(rr.Model, creditCode)
		if rr.Error != nil {
			return 0, CustomError(structs.Categories.Internal, "unable to create transaction with error: %s", rr.Error.Error())
		}
	} else {
		switch rr.Model.Transactions[0].Status {
		case constants.TransactionStatusEnums.Requested:
			return 0, CustomError(structs.Categories.BusinessLogic, "your request is in progress")
		case constants.TransactionStatusEnums.Successful:
			return 0, CustomError(structs.Categories.BusinessLogic, "already used this code")
		case constants.TransactionStatusEnums.FailedForCreditCodeLimitations:
			return 0, CustomError(structs.Categories.BusinessLogic, "credit code limitation reached")
		default:
			return 0, CustomError(structs.Categories.BusinessLogic, "undefined transaction status")
		}
	}

	response := c.service.Send(creditCode, rr.Model.Transactions[0].ID)

	if response.Err != nil {
		if err := c.processMSResponseErr(response.Err, rr.Model.Transactions[0].ID); err != nil {
			return 0, err
		}

		return 0, response.Err
	}

	if err := c.userRepo.RunChargeBalanceDBFunc(rr.Model, response.Amount, response.LogId); err != nil {
		return 0, c.processDBFuncErr(err)
	}

	return response.Amount, response.Err
}

func (c *ChargeCreditByCode) processMSResponseErr(err *structs.CustomError, transId uint) *structs.CustomError {

	if errors.Is(err, exceptions.CreditCodeLimitationError{}) {
		return nil
	}

	res := c.transRepo.UpdateTransactionAsLimited(transId)
	if res.RowsAffected == 0 {
		return CustomError(structs.Categories.Internal, "can not set transaction [%s] as limited", transId)
	}

	return nil

}

func (c ChargeCreditByCode) processDBFuncErr(err error) *structs.CustomError {

	decodeError := CustomError(structs.Categories.Internal, "can not decode database function [charge_by_credit_code] error details. error: %s", err.Error())
	pgErr, ok := err.(*pgconn.PgError)

	if !ok {
		return decodeError
	}

	var errorDetail map[string]interface{}
	if err = json.Unmarshal([]byte(pgErr.Detail), &errorDetail); err == nil {
		switch errorDetail["code"].(float64) {
		case 1:
			return CustomError(structs.Categories.Internal, "direct call of db func [charge_by_credit_code]")
		case 2, 3:
			fallthrough
		default:
			return CustomError(structs.Categories.Internal, "database function [charge_by_credit_code] error: %s ", err.Error())

		}
	}

	return decodeError
}
