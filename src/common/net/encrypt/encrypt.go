package encrypt

const MAX_SIZE = 206

type Encrypt struct {
	Code [MAX_SIZE]byte
	Step int
}

func NewEncrypt(key1, key2 byte) *Encrypt {
	inst := new(Encrypt)
	inst.ChangeCode(key1, key2)
	return inst
}

func (this *Encrypt) DoEncrypt(buf []byte) {
	for i := 0; i < len(buf); i++ {
		if this.Step >= MAX_SIZE {
			this.Step = 0
		}

		buf[i] ^= this.Code[this.Step]
		this.Step++
	}
}

func (this *Encrypt) DoDeencrypt(buf []byte) {
	this.DoEncrypt(buf)
}

func (this *Encrypt) ChangeCode(key1, key2 byte) {
	key := key1

	for i := 0; i < len(this.Code); i++ {
		this.Code[i] = key
		key = byte((int32(key)*int32(key1) + int32(key2)) % 256)

		if key == this.Code[i] {
			key = byte(int32(key+1) % 256)
		}
	}
}
