package services

import (
	"arvan-wallet-service/infrastructures"
	"arvan-wallet-service/types/exceptions"
	"arvan-wallet-service/types/interfaces/repositories"
	"arvan-wallet-service/types/interfaces/services"
	"arvan-wallet-service/types/structs"
	. "arvan-wallet-service/utils"
	"errors"
	"fmt"
)

type CircuitBreakerProxyChargeCreditByCode struct {
	userRepo    repositories.IUserRepository
	transRepo   repositories.ITransactionRepository
	service     services.IRequestToUtilizeCode
	stateHolder *creditCodeStateHolder
}

func (c *CircuitBreakerProxyChargeCreditByCode) SetRepositories(userRepo repositories.IUserRepository, transRepo repositories.ITransactionRepository) {
	c.userRepo = userRepo
	c.transRepo = transRepo
}

func (c *CircuitBreakerProxyChargeCreditByCode) SetService(service services.IRequestToUtilizeCode) {
	c.service = service
}

func (c *CircuitBreakerProxyChargeCreditByCode) Charge(mobile string, creditCode string) (uint, *structs.CustomError) {
	c.stateHolder = getStateHolder()

	c.stateHolder.mutex.Lock() // synchronize with other potential writers
	if c.checkAlreadyLimited(creditCode) {
		fmt.Println("prevent continue!!!!!!!!")
		c.stateHolder.mutex.Unlock()
		return 0, CustomError(
			structs.Categories.BusinessLogic,
			c.stateHolder.readValue(creditCode).lastErr.Error(),
		)
	}
	c.stateHolder.mutex.Unlock()

	amount, err, customErr := c.proxyToInternalService(mobile, creditCode)
	finalErr := c.handleErrors(err, customErr)

	if errors.Is(customErr, exceptions.CreditCodeLimitationError{}) {
		c.stateHolder.addValue(creditCode, 1, finalErr)
	}

	return amount, finalErr
}

func (c *CircuitBreakerProxyChargeCreditByCode) checkAlreadyLimited(creditCode string) bool {

	val := c.stateHolder.readValue(creditCode)
	if val == nil {
		inqErr := infrastructures.Resolve[services.IRequestToInquiryCode]().Send(creditCode)

		if inqErr != nil {
			c.stateHolder.storeValue(creditCode, 1, inqErr)
			return true
		} else {
			c.stateHolder.storeValue(creditCode, 0, nil)
			return false
		}

	} else if val.counter > 0 {
		fmt.Printf(" %d\n", val)
		return true
	}

	return false
}

func (c *CircuitBreakerProxyChargeCreditByCode) proxyToInternalService(mobile string, creditCode string) (uint, error, *structs.CustomError) {
	var (
		amount    uint
		customErr *structs.CustomError
	)

	_, err := c.stateHolder.cb.Execute(func() (interface{}, error) {

		s := &ChargeCreditByCode{}
		s.SetRepositories(c.userRepo, c.transRepo)
		s.SetService(c.service)
		amount, customErr = s.Charge(mobile, creditCode)

		if customErr != nil {
			return nil, customErr
		}

		return nil, nil
	})

	return amount, err, customErr
}

func (c *CircuitBreakerProxyChargeCreditByCode) handleErrors(err error, customErr *structs.CustomError) *structs.CustomError {
	if customErr != nil {
		return customErr
	}

	if err != nil {
		return CustomError(structs.Categories.Internal, err.Error())
	}

	return nil
}
