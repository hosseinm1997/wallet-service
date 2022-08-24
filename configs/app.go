package configs

import (
	"arvan-wallet-service/providers"
	"github.com/spf13/viper"
)

type Config map[string]any

func App() map[string]any {
	return Config{
		"environment": viper.GetString("APP_ENV"),
		"providers": []func(){
			providers.SystemProvides,
			providers.RepositoryProviders,
			providers.ServiceProvider,
			providers.EndpointProviders,
		},
	}
}
