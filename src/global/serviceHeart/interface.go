package serviceHeart

import "github.com/Meland-Inc/game-services/src/global/component"

type ServiceHeartInterface interface {
	component.ModelInterface

	SendHeart(curMs int64) error
}
