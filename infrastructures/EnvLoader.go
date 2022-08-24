package infrastructures

import (
	"github.com/spf13/viper"
	"sync"
)

type IEnvLoader interface {
	LoadFromFile(filename string) error
}

type envLoader struct{}

func (e *envLoader) LoadFromFile(filename string) error {
	viper.SetConfigFile(filename)

	return viper.ReadInConfig()
}

var (
	el     *envLoader
	elOnce sync.Once
)

func InitEnvLoader() *envLoader {
	if el == nil {
		elOnce.Do(func() {
			el = &envLoader{}
		})
	}

	return el
}
