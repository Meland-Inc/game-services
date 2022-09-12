package component

type Model interface {
	OnInit() error
	OnStart() error
	OnTick() error
	OnStop() error
	OnExit() error
}
