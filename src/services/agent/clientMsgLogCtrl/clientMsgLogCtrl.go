package clientMsgLogCtrl

import "game-message-core/proto"

var pbMsgLogCtrl = make(map[proto.EnvelopeType]bool) // map{[EnvelopeType]=needShow}

func init() {
	registerMsgLogCtrl(proto.EnvelopeType_Ping, false)
	registerMsgLogCtrl(proto.EnvelopeType_BroadCastItemAdd, false)
	registerMsgLogCtrl(proto.EnvelopeType_BroadCastEntityMove, false)
	registerMsgLogCtrl(proto.EnvelopeType_UpdateSelfLocation, false)
	registerMsgLogCtrl(proto.EnvelopeType_BroadCastMapEntityUpdate, false)
	// registerMsgLogCtrl(proto.EnvelopeType_UseSkill, false)
	// registerMsgLogCtrl(proto.EnvelopeType_BroadCastEntityCombat, false)
}

func registerMsgLogCtrl(msgType proto.EnvelopeType, show bool) {
	pbMsgLogCtrl[msgType] = show
}

func PrintCliMsgLog(msgType proto.EnvelopeType) bool {
	show, exist := pbMsgLogCtrl[msgType]
	if !exist {
		return true
	}
	return show
}
