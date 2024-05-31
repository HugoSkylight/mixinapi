package mixinapi

type MixinApp struct {
	config       Config        `description:"application configuration"`
	mixinRouters []MixinRouter `description:"routers"`
}

type Config struct {
	host        string  `description:"host"`
	port        string  `description:"port"`
	title       string  `description:"application title"`
	summary     string  `description:"application summary"`
	description string  `description:"application description"`
	version     string  `description:"application version"`
	openapiUrl  string  `description:"openapi url"`
	docsUrl     string  `description:"documentation url"`
	openapi     OpenAPI `description:"openapi"`
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

func NewMixinApp() *MixinApp {
	app := &MixinApp{}
	app.initialize()
	return app
}