package dbData

type UsingNft struct {
	NftId     string `gorm:"primaryKey;autoIncrement" json:"nftId,string"`
	UserId    int64  `json:"userId"`
	Cid       int32  `json:"cid"`
	AvatarPos int32  `json:"avatarPos"`
}
