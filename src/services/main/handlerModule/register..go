package handlerModule

import (
	"game-message-core/grpc"
	"game-message-core/proto"

	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/Meland-Inc/game-services/src/services/main/handlerModule/clientHandler"
	"github.com/Meland-Inc/game-services/src/services/main/handlerModule/serviceHandler"
	"github.com/Meland-Inc/game-services/src/services/main/handlerModule/web3Handler"
)

func (p *HandlerModule) RegisterClientEvent() {
	// sing in message
	p.AddClientEvent(proto.EnvelopeType_SigninPlayer, clientHandler.SingInHandler)

	// item message
	p.AddClientEvent(proto.EnvelopeType_ItemGet, clientHandler.ItemGetHandler)
	p.AddClientEvent(proto.EnvelopeType_ItemUse, clientHandler.ItemUseHandler)
	p.AddClientEvent(proto.EnvelopeType_UpdateAvatar, clientHandler.LoadAvatarHandler)
	p.AddClientEvent(proto.EnvelopeType_UnloadAvatar, clientHandler.UnloadAvatarHandler)

	//item slot message
	p.AddClientEvent(proto.EnvelopeType_GetItemSlot, clientHandler.ItemSlotGetHandler)
	p.AddClientEvent(proto.EnvelopeType_UpgradeItemSlot, clientHandler.ItemSlotUpgradeHandler)

	// player level
	p.AddClientEvent(proto.EnvelopeType_UpgradePlayerLevel, clientHandler.PlayerLevelUpgradeHandler)

	// granary message
	p.AddClientEvent(proto.EnvelopeType_QueryGranary, clientHandler.QueryGranaryHandler)
	p.AddClientEvent(proto.EnvelopeType_GranaryCollect, clientHandler.GranaryCollectHandler)

	// land message
	p.AddClientEvent(proto.EnvelopeType_QueryLands, clientHandler.QueryLandsHandler)
	p.AddClientEvent(proto.EnvelopeType_Build, clientHandler.LandBuildHandler)
	p.AddClientEvent(proto.EnvelopeType_Recycling, clientHandler.RecyclingHandler)
	p.AddClientEvent(proto.EnvelopeType_MintBattery, clientHandler.MintBatteryHandler)
	p.AddClientEvent(proto.EnvelopeType_Charged, clientHandler.ChargedHandler)
	p.AddClientEvent(proto.EnvelopeType_Harvest, clientHandler.HarvestHandler)
	p.AddClientEvent(proto.EnvelopeType_Collection, clientHandler.CollectionHandler)
	p.AddClientEvent(proto.EnvelopeType_SelfNftBuilds, clientHandler.SelfNftBuildsHandler)

}

func (p *HandlerModule) RegisterGameServiceDaprCall() {
	p.AddGameServiceDaprCall(
		string(grpc.MainServiceActionGetHomeData),
		serviceHandler.GRPCGetHomeDataHandler,
	)
	p.AddGameServiceDaprCall(
		string(grpc.MainServiceActionGetAllBuild),
		serviceHandler.GRPCGetAllBuildHandler,
	)

	p.AddGameServiceDaprCall(
		string(grpc.UserActionGetUserData),
		serviceHandler.GRPCGetUserDataHandler,
	)
	p.AddGameServiceDaprCall(
		string(grpc.MainServiceActionMintNFT),
		serviceHandler.GRPCMintUserNftHandler,
	)
	p.AddGameServiceDaprCall(
		string(grpc.MainServiceActionTakeNFT),
		serviceHandler.GRPCTakeUserNftHandler,
	)

}

func (p *HandlerModule) RegisterGameServiceDaprEvent() {
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventSaveHomeData),
		serviceHandler.GRPCSaveHomeDataEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventGranaryStockpile),
		serviceHandler.GRPCGranaryStockpileEvent,
	)

	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventUserEnterGame),
		serviceHandler.GRPCUserEnterGameEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventUserLeaveGame),
		serviceHandler.GRPCUserLeaveGameEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventSavePlayerData),
		serviceHandler.GRPCSavePlayerDataEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventKillMonster),
		serviceHandler.GRPCKillMonsterEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventPlayerDeath),
		serviceHandler.GRPCPlayerDeathEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventUserTaskReward),
		serviceHandler.GRPCUserTaskRewardEvent,
	)

	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventUserChangeService),
		serviceHandler.GRPCUserChangeServiceEvent,
	)

}

func (p *HandlerModule) RegisterWeb3DaprCall() {
	p.AddWeb3DaprCall(
		string(message.GameDataServiceActionDeductUserExp),
		web3Handler.Web3DeductUserExpHandler,
	)
	p.AddWeb3DaprCall(
		string(message.GameDataServiceActionGetPlayerInfoByUserId),
		web3Handler.Web3GetPlayerDataHandler,
	)

}

func (p *HandlerModule) RegisterWeb3DaprEvent() {
	p.AddWeb3DaprEvent(
		string(message.SubscriptionEventMultiLandDataUpdateEvent),
		web3Handler.Web3MultiLandDataUpdateEvent,
	)
	p.AddWeb3DaprEvent(
		string(message.SubscriptionEventMultiRecyclingEvent),
		web3Handler.Web3MultiRecyclingEvent,
	)
	p.AddWeb3DaprEvent(
		string(message.SubscriptionEventMultiBuildUpdateEvent),
		web3Handler.Web3MultiBuildUpdateEvent,
	)

	p.AddWeb3DaprEvent(
		string(message.SubscriptionEventUpdateUserNFT),
		web3Handler.Web3UpdateUserNftEvent,
	)
	p.AddWeb3DaprEvent(
		string(message.SubscriptionEventMultiUpdateUserNFT),
		web3Handler.Web3MultiUpdateUserNftEvent,
	)

}
