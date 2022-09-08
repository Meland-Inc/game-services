package application

type AppInterface interface {
	OnInit() error
	OnStart() error
	Run()
	OnExit()
}

type Application struct {
	appHandler AppInterface
}

var app Application

func Init(handler AppInterface) {
	app.appHandler = handler
	if err := app.appHandler.OnInit(); err != nil {
		panic(err)
	}
	if err := app.appHandler.OnStart(); err != nil {
		panic(err)
	}
}

func Run() {
	app.appHandler.Run()
}
