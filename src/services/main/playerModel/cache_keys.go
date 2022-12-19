package playerModel

import "fmt"

func (p *PlayerDataModel) playerItemsCacheKey(userId int64) string {
	return fmt.Sprintf("game_player_items_%d", userId)
}

func (p *PlayerDataModel) playerItemsSlotCacheKey(userId int64) string {
	return fmt.Sprintf("game_player_items_slot_%d", userId)
}
