package playerModel

import "fmt"

func (p *PlayerModel) getPlayerSceneDataKey(userId int64) string {
	return fmt.Sprintf("player_scene_data_key_%d", userId)
}
