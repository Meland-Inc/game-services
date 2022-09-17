package itemModel

import (
	"fmt"
	"game-message-core/proto"
)

func (p *ItemModel) playerItemsCacheKey(userId int64) string {
	return fmt.Sprintf("game_player_items_%d", userId)
}

func (p *ItemModel) GetPlayerItems(userId int64) (*PlayerItems, error) {
	cacheKey := p.playerItemsCacheKey(userId)
	iData, err := p.cache.GetOrStore(
		cacheKey,
		func() (interface{}, error) {
			items, err := p.LoadItems(userId)
			if err != nil {
				return nil, err
			}
			playerItems := &PlayerItems{
				UserId: userId,
				Items:  items,
			}
			return playerItems, err
		},
		p.cacheTTL,
	)
	if err != nil {
		return nil, err
	}

	p.cache.Touch(cacheKey, p.cacheTTL)
	return iData.(*PlayerItems), nil
}

func (p *ItemModel) ItemById(userId int64, nftId string) (*Item, error) {
	userItems, err := p.GetPlayerItems(userId)
	if err != nil {
		return nil, err
	}

	for _, it := range userItems.Items {
		if it.Id == nftId {
			return it, nil
		}
	}

	return nil, fmt.Errorf("Item not found")
}

func (p *ItemModel) UsingAvatars(userId int64) (avatars []*Item, err error) {
	userItems, err := p.GetPlayerItems(userId)
	if err != nil {
		return nil, err
	}
	for _, it := range userItems.Items {
		if it.Used &&
			it.AvatarPos >= int32(proto.AvatarPosition_AvatarPositionHead) &&
			it.AvatarPos <= int32(proto.AvatarPosition_AvatarPositionWeapon) {
			avatars = append(avatars, it)
		}
	}
	return avatars, err
}
