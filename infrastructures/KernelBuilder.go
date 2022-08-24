package infrastructures

import (
	"sync"
)

type IKernelBuilder interface {
	Build(config map[string]any) IKernel
}

type kernelBuilder struct{}

var (
	k          *kernel
	kernelOnce sync.Once
)

func (kb *kernelBuilder) Build(config map[string]any) IKernel {
	if k == nil {
		kernelOnce.Do(func() {
			k = &kernel{}
			k.loadEnvVars()
			k.initializeServiceContainer(config)
			k.initializeDatabase()
			k.initializeRouter()
		})
	}

	return k
}

func KernelBuilder() IKernelBuilder {
	return &kernelBuilder{}
}
