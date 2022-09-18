package playerModel

import "fmt"

func (p *PlayerDataModel) getPlayerSceneDataKey(userId int64) string {
	return fmt.Sprintf("player_scene_data_key_%d", userId)
}

func (p *PlayerDataModel) playerItemsCacheKey(userId int64) string {
	return fmt.Sprintf("game_player_items_%d", userId)
}

func (p *PlayerDataModel) playerItemsSlotCacheKey(userId int64) string {
	return fmt.Sprintf("game_player_items_slot_%d", userId)
}
