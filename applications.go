package mixinapi

type MixinApp struct {
	config Config `description:"application configuration"`
	openapi OpenAPI `description:"openapi"`
	mixinRouters []
}

type Config struct {
	debug bool `description:"enable debug mode"`
	host string `description:"host"`
	port string `description:"port"`
	title string `description:"application title"`
	summary string `description:"application summary"`
	description string `description:"application description"`
	version string `description:"application version"`
	openapiUrl string `description:"openapi url"`
	docsUrl string `description:"documentation url"`
}	

func (app *MixinApp) GetConfig() Config {
	return app.config
}

func (app *MixinApp) Start(host, port string) {
	app.config.host = host
	app.config.port = port

}


func (app *MixinApp) initialize() *MixinApp {
	return app
}