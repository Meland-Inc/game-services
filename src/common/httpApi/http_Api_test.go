package util

import (
	"encoding/json"
	"testing"
	"time"

)

func Test_HttpPublicApi(t *testing.T) {
	

	// {"Id":"149095152676290618","AccountId":715248,"RoleId":1001,"Name":"*只见花开","Gender":"male","RoleIcon":"icon_headPortrait_1",
	// 	"Feature":"{\"hair\":101002,\"clothes\":102002,\"glove\":105000,\"pants\":103001,\"face\":100003,\"shoes\":104003}",
	// 	"Lv":1,"Exp":0,"Guide":0,"CurHp":500,"HungryPoint":100,"ThirstyPoint":100,"FatiguePoint":3000,"CreatedAt":"2021-06-17T13:27:16Z",
	// 	"UpdatedAt":"2021-06-17T13:27:16Z","IsMgr":false,"BanTime":0,"LastloginTime":0,"Error":null}:
	//
	// 		first path segment in URL cannot contain colon,  url=https://account.wkcoding.com/1/updateplayer

	// urlstr := "https://account.wkcoding.com/1/updateplayer"
	// reqBody, err := json.Marshal(player)
	// t.Log(err)
	// t.Log(reqBody)
	// t.Log("-------------------------------------------------------------------------------------")

	// body, err1 := AuthHttpRequest("POST", string(reqBody[:]), urlstr)
	// t.Log(err1)
	// t.Log(body)
	// t.Log("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

	// ret := stData.UpdatePlayerRet{}
	// err = json.Unmarshal(body, &ret)
	// t.Log(ret)
	// t.Log(err)

}
