package message

const landIdTemplates int64 = 10000000

// land Id 和 xy 坐标系转换关系

func XyToLandId(r, c int16) int64 {
	// 例： 1000033300000004 RC 各占8位 总共16位
	// 10000333 第一位 1 代表负数  333 为R坐标 中间用0填充(1*10000000+333)
	// 00000004 第一位 0 代表正数  4   为C坐标 中间用0填充(0*10000000+333)
	var rTag, cTag int64
	if r < 0 {
		rTag = 1
		r = -r
	}
	if c < 0 {
		cTag = 1
		c = -c
	}

	rOffset := rTag*landIdTemplates + int64(r)
	cOffset := cTag*landIdTemplates + int64(c)
	return rOffset*landIdTemplates*10 + cOffset
}

func LandIdToXy(landId int64) (r, c int16) {
	// 例： 1000033300000004 RC 各占8位 总共16位
	// 10000333 第一位 1 代表负数  333 为R坐标 中间用0填充(1*10000000+333)
	// 00000004 第一位 0 代表正数  4   为C坐标 中间用0填充(0*10000000+333)

	rOffset := landId / (landIdTemplates * 10)
	cOffset := landId % (landIdTemplates * 10)

	rTag := rOffset / landIdTemplates
	r = int16(rOffset % landIdTemplates)
	if rTag > 0 {
		r = -r
	}

	cTag := cOffset / landIdTemplates
	c = int16(cOffset % landIdTemplates)
	if cTag > 0 {
		c = -c
	}
	return
}
