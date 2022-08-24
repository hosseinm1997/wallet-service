# Wallet service part of "Arvan Cloud" challenge

This service is responsible for handling user requests for charging balance by credit codes and as a proxy to `credit service`.

## Quick start
To quickly jump into the main logic go to following links:

### Visualized documents:
- [**Process Flow**](https://online.visual-paradigm.com/community/share/arvan-challenge-flow)
- [**ERD**](https://online.visual-paradigm.com/community/share/arvan-challenge-erd)



### Main Endpoints:
- [Charge user balance by amount API logic](https://github.com/hosseinm1997/wallet-service/blob/30e2401267eb551462005b84bfadcab71a1e876e/http/endpoints/user/CreditEndpoint.go#L15)
- [Show user balance API logic](https://github.com/hosseinm1997/wallet-service/blob/30e2401267eb551462005b84bfadcab71a1e876e/http/endpoints/user/CreditEndpoint.go#L37)
- [List charged users API logic](https://github.com/hosseinm1997/wallet-service/blob/30e2401267eb551462005b84bfadcab71a1e876e/http/endpoints/user/CreditEndpoint.go#L37)

## Overview

### Approach

One of the most important challenge of wallet service was overcomming `credit service` outage. My approach to solving this problem is using circuit breaker pattern.

### Architecture
I Use the dual write strategy for communication between these two microservices. It's better to use distributed transaction management patterns, especially the Saga pattern via the Orchestrator model. 


I visualize all inter-microservice communications via **BPMN** notation language.
[**Process Flow**](https://online.visual-paradigm.com/community/share/arvan-challenge-flow)

You can see the **ERD** of this service via this like:
[**ERD**](https://online.visual-paradigm.com/community/share/arvan-challenge-erd)

Open and click Edit This Design

### Framework
This service was made based on a simple framework made by myself (in a limited time). I'm not interested in `reinvent the wheel` myself!! My idea behind this is to dig into the Go language deeper. It has following features:

- IoC implemented using service container, created by new `generic` feature of go 1.18. [see ServiceContainer.go](https://github.com/hosseinm1997/credit-service/blob/main/infrastructures/ServiceContainer.go)
- Routing system using middlewares and contextes. [see RoutingSystem.go](https://github.com/hosseinm1997/credit-service/blob/main/infrastructures/RoutingSystem.go)
- Easy exception handling with `Respond()` helper function. [see ResponseFormatter.go](https://github.com/hosseinm1997/credit-service/blob/main/http/middlewares/ResponseFormatter.go), [see an example](https://github.com/hosseinm1997/credit-service/blob/ab1eda279aa9e2a4d02b4d752e09de0e0f3da42f/http/endpoints/SpendCodeEndpoint.go#L71)
- Handling env variables
- Service and repository pattern considered

### Packages
Direct packages used:

- Viper for managing env variables
- Go-chi for routing
- Gorm for using as an ORM and query builder
- Gobreaker for implementing ciruite breaker

### Framework TODOs:
- [ ] Use gRPC for communication
- [ ] Use Saga pattern instead of dual write strategy
- [ ] Segregate reads and writes via CQRS pattern 
- [ ] Use gorm migration
- [ ] Use unit testing with high code coverage
- [ ] Pass context into internal services
- [ ] Use swagger for API documentation

### Wallet Service TODOs:
- [ ] Make the circuit breaker service more precise

<br/>
<br/>
<br/>

## Wallet Service Docs

This service is standing between user and `credit service` as a proxy. It charge the balance, show the balance and has admin api for showing list of charged user by a specific credit code

### - Charge API:
- **signature**: `/user/{mobile}/credit/code/{code}` 
- **inputs**:
    `{mobile}` (string): user mobile received from user.

    `{code}` (string): credit code text received from user.

- **description**: This endpoint request for user and transactions to the database. If no transaction found, it stores user and
his transaction with status: `requested`, then make request to the `credit service` and after getting an log id make the status: `successful` and
charge user balance. If error occures between microservices, circuit pattern state will be changed to open. 
With a specific interval it will be change to half-open state. **`Circuite breaker`** implemented via **`Proxy pattern`** over main service.


### - Balance API:
- **signature**: `/user/{mobile}/credit}` 
- **inputs**:

    `{mobile}` (string): user mobile received from user.

- **description**: This service is used to get the user balance from `users` table if exists.


### - List API:
- **signature**: `/admin/credit/{code}/list` 
- **inputs**:

    `{code}` (string): credit code text received from admin.

- **description**: This endpoint will load all transactions of following credit code with status: `successful`.

