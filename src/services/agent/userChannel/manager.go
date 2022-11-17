package userChannel

import "sync"

var mgrInstance *UserChannelManager

func GetInstance() *UserChannelManager {
	if mgrInstance == nil {
		NewUserChannelManager()
	}
	return mgrInstance
}

func NewUserChannelManager() *UserChannelManager {
	mgrInstance = &UserChannelManager{}
	return mgrInstance
}

type UserChannelManager struct {
	userChannelsByOwner sync.Map
	userChannelsById    sync.Map
	count               int32
}

func (mgr *UserChannelManager) OnlineCount() int32 { return mgr.count }

func (mgr *UserChannelManager) UserChannelById(id string) *UserChannel {
	iChannel, exist := mgr.userChannelsById.Load(id)
	if !exist {
		return nil
	}
	chanel, _ := iChannel.(*UserChannel)
	return chanel
}

func (mgr *UserChannelManager) UserChannelByOwner(owner int64) *UserChannel {
	iChannel, exist := mgr.userChannelsByOwner.Load(owner)
	if !exist {
		return nil
	}
	chanel, _ := iChannel.(*UserChannel)
	return chanel
}

func (mgr *UserChannelManager) AddUserChannelById(channel *UserChannel) {
	if channel == nil {
		return
	}

	channelId := channel.GetId()
	_, exist := mgr.userChannelsById.Load(channelId)
	if !exist {
		mgr.count++
	}
	mgr.userChannelsById.Store(channelId, channel)
}

func (mgr *UserChannelManager) AddUserChannelByOwner(channel *UserChannel) {
	if channel == nil || channel.GetOwner() == 0 {
		return
	}

	_, exist := mgr.userChannelsById.Load(channel.GetId())
	if !exist {
		mgr.AddUserChannelById(channel)
	}
	mgr.userChannelsByOwner.Store(channel.GetOwner(), channel)
}

func (mgr *UserChannelManager) RemoveUserChannel(channel *UserChannel) {
	if channel == nil {
		return
	}

	channelId := channel.GetId()
	if _, exist := mgr.userChannelsById.Load(channelId); exist {
		mgr.count--
	}
	mgr.userChannelsById.Delete(channelId)
	mgr.userChannelsByOwner.Delete(channel.GetOwner())
}

func (mgr *UserChannelManager) Range(f func(channel *UserChannel) bool) {
	mgr.userChannelsById.Range(func(key, value interface{}) bool {
		cha := value.(*UserChannel)
		return f(cha)
	})
}
