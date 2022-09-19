package serviceMsgHandle

import (
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/spf13/cast"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func Web3DeductUserExpHandler(iMsg interface{}) {
	input, ok := iMsg.(*message.DeductUserExpInput)
	if !ok {
		serviceLog.Error("iMsg to UserEnterGameEvent failed")
		return
	}

	deductExp := cast.ToInt32(input.DeductExp)
	userId := cast.ToInt64(input.UserId)

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}

	dataModel.DeductExp(userId, deductExp)
}
