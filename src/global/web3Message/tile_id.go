package message

const tileIdTemplates int64 = 10000000

// land Id 和 xy 坐标系转换关系
func XyToTileId(x, y int32) int64 {
	// 例： 1000033300000004 RC 各占8位 总共16位
	// 10000333 第一位 1 代表负数  333 为R坐标 中间用0填充(1*10000000+333)
	// 00000004 第一位 0 代表正数  4   为C坐标 中间用0填充(0*10000000+333)
	var rTag, cTag int64
	if x < 0 {
		rTag = 1
		x = -x
	}
	if y < 0 {
		cTag = 1
		y = -y
	}

	rOffset := rTag*tileIdTemplates + int64(x)
	cOffset := cTag*tileIdTemplates + int64(y)
	return rOffset*tileIdTemplates*10 + cOffset
}

func TileIdToXy(landId int64) (x, y int32) {
	// 例： 1000033300000004 RC 各占8位 总共16位
	// 10000333 第一位 1 代表负数  333 为R坐标 中间用0填充(1*10000000+333)
	// 00000004 第一位 0 代表正数  4   为C坐标 中间用0填充(0*10000000+333)

	rOffset := landId / (tileIdTemplates * 10)
	cOffset := landId % (tileIdTemplates * 10)

	rTag := rOffset / tileIdTemplates
	x = int32(rOffset % tileIdTemplates)
	if rTag > 0 {
		x = -x
	}

	cTag := cOffset / tileIdTemplates
	y = int32(cOffset % tileIdTemplates)
	if cTag > 0 {
		y = -y
	}
	return
}
