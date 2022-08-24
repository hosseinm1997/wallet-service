package infrastructures

type IKernel interface {
	loadEnvVars()
	initializeServiceContainer(config map[string]any)
	initializeDatabase()
	initializeRouter()
}

type kernel struct{}

func (k *kernel) loadEnvVars() {
	if err := InitEnvLoader().LoadFromFile(".env"); err != nil {
		println(err)
		// todo: implement error handling
	}

}

func (k *kernel) initializeDatabase() {
	DatabaseSetup()
}

func (k *kernel) initializeServiceContainer(config map[string]any) {
	InitializeServiceContainer(config)
}

func (k *kernel) initializeRouter() {
	ChiRouter()
}
