package uuid

import (
	"github.com/google/uuid"
)

func NextUUid() string {
	// V4 基于随机数
	u4 := uuid.New()
	return u4.String()
}
