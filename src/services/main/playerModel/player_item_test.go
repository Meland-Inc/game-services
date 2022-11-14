package playerModel

import (
	"fmt"
	"game-message-core/protoTool"
	"testing"

	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func NFT() message.NFT {
	color := "6ecb67"
	image := "https://meland-inc.github.io/icon2s3/data/MV01_eq_dagger_01_sheepHorn_ic.png"
	return message.NFT{
		Address:        "0x0000000000000000000000000000000000000000",
		Amount:         99,
		AmountOfChange: 0,
		Id:             "local#0x0000000000000000000000000000000000000000#156949617517294774568639930436394721876",
		IsMelandAI:     true,
		ItemId:         "71010001",
		Network:        "local",
		ProductId:      "0x45abc8537f1ecc9082bd3fc433b815a9f84fad0e33a0264fe6eb3858af209036",
		TokenId:        "156949617517294774568639930436394721876",
		TokenURL:       "",
		Metadata: &message.NFTMetadata_1{
			BackgroundColor: &color,
			Description:     "Sheep Horn Dagger",
			Image:           &image,
			Name:            "Sheep Horn Dagger",
			MelandAttributes: []message.MelandAttribute{
				message.MelandAttribute{TraitType: "Quality", Value: "Enhanced"},
				message.MelandAttribute{TraitType: "Type", Value: "Dagger"},
				message.MelandAttribute{TraitType: "Attack", Value: "11"},
				message.MelandAttribute{TraitType: "Attack Speed", Value: "10"},
				message.MelandAttribute{TraitType: "CoreSkillId", Value: "810002"},
				// message.MelandAttribute{TraitType: "Crit Damage", Value: "100"},
				// message.MelandAttribute{TraitType: "Crit Points", Value: "80"},
				// message.MelandAttribute{TraitType: "Defence", Value: "10"},
				// message.MelandAttribute{TraitType: "Dodge Points", Value: "45"},
				// message.MelandAttribute{TraitType: "Gender", Value: "2"},
				// message.MelandAttribute{TraitType: "Get Buff", Value: "200001"},
				// message.MelandAttribute{TraitType: "HP Recovery", Value: "10"},
				// message.MelandAttribute{TraitType: "Hit Points", Value: "90"},
				// message.MelandAttribute{TraitType: "Learn Recipe", Value: "88"},
				// message.MelandAttribute{TraitType: "MaxHP", Value: "90"},
				// message.MelandAttribute{TraitType: "Move Speed", Value: "5"},
				// message.MelandAttribute{TraitType: "Rarity", Value: "3"},
				// message.MelandAttribute{TraitType: "Requires level", Value: "5"},
				// message.MelandAttribute{TraitType: "Restore HP", Value: "10"},
				// message.MelandAttribute{TraitType: "SkillLevel", Value: "3"},
			},
		},
	}
}

func getItem() *Item {
	return NFTToItem(798, NFT())
}

func Test_ItemData(t *testing.T) {
	item := getItem()

	// 全量nft json数据 测试
	pbData1 := item.ToNetItem()
	pbData1.NftData = nil
	msgBs1, err := protoTool.MarshalProto(pbData1)
	t.Log(fmt.Sprintf("\n 全量nft json数据长度[%d], err:%+v", len(msgBs1), err))

	// 使用 proto nft data数据 测试
	pbData2 := item.ToNetItem()
	// pbData2.NftJsonData = ""
	// pbData2.Attribute = nil
	msgBs2, err := protoTool.MarshalProto(pbData2)
	t.Log(fmt.Sprintf("\n proto nft data 数据长度[%d], err:%+v", len(msgBs2), err))

}
