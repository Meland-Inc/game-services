package playerModel

import (
	"fmt"

	"github.com/Meland-Inc/game-services/src/global/component"
)

func GetPlayerDataModel() (*PlayerDataModel, error) {
	iPlayerModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_PLAYER_DATA)
	if !exist {
		return nil, fmt.Errorf("player data model not found")
	}
	dataModel, _ := iPlayerModel.(*PlayerDataModel)
	return dataModel, nil
}
