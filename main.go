package main

import (
	"arvan-wallet-service/configs"
	"arvan-wallet-service/http/routing"
	"arvan-wallet-service/infrastructures"
	"arvan-wallet-service/infrastructures/interfaces"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	infrastructures.KernelBuilder().Build(configs.App())
	routingSystem := infrastructures.Resolve[interfaces.IChiRouter]()
	http.ListenAndServe(
		getAddress(),
		routingSystem.InitRouter(routing.Routes),
	)
}

func getAddress() string {
	return fmt.Sprintf("%s:%s", viper.GetString("SERVER_HOST"), viper.GetString("SERVER_PORT"))
}
