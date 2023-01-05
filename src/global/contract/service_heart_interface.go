package contract

type IServiceHeartInterface interface {
	IModuleInterface
	SendHeart(curMs int64) error
}
