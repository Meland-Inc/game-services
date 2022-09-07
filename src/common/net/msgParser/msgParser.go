package msgParser

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/Meland-Inc/game-services/src/common/net/encrypt"
)

var (
	HEAD_SIZE = 4
	ID_SIZE   = uint32(2)
	DOENCRYPT = false // 配置是否进行加密解密
	MSG_LIMIT = 8192
)

// --------------
// | len | data |
// --------------
type MsgParser struct {
	minMsgLen uint32
	maxMsgLen uint32

	encClient *encrypt.Encrypt // 客户端解密
	encServer *encrypt.Encrypt // 服务端加密
}

func NewMsgParser() *MsgParser {
	return &MsgParser{
		minMsgLen: ID_SIZE,
		maxMsgLen: 256 * 256,
		encClient: encrypt.NewEncrypt(119, 127),
		encServer: encrypt.NewEncrypt(119, 127),
	}
}

func (this *MsgParser) Transform(value []byte) uint32 {
	return uint32(binary.LittleEndian.Uint16(value))
}

// 解码
func (this *MsgParser) Decode(head []byte, reader *bufio.Reader) ([]byte, error) {
	_, err := io.ReadFull(reader, head)
	if err != nil {
		return nil, err
	}

	if DOENCRYPT {
		this.encClient.DoEncrypt(head)
	}

	dataLen := this.Transform(head)

	if dataLen > this.maxMsgLen {
		return nil, errors.New(fmt.Sprintf("message too long:%d", dataLen))
	}

	if dataLen < this.minMsgLen {
		return nil, errors.New(fmt.Sprintf("message too short:%d", dataLen))
	}

	data := make([]byte, dataLen)

	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, err
	}

	if DOENCRYPT {
		this.encClient.DoEncrypt(data)
	}

	return data, nil
}

// 编码
func (this *MsgParser) Encode(data []byte) ([]byte, error) {
	msgLen := uint32(len(data))

	if msgLen > this.maxMsgLen {
		return nil, errors.New(fmt.Sprintf("message too long.max:%v cur:%v", this.maxMsgLen, msgLen))
	} else if msgLen < this.minMsgLen {
		return nil, errors.New(fmt.Sprintf("message too short.min:%v cur:%v", this.minMsgLen, msgLen))
	}

	head := make([]byte, HEAD_SIZE)
	binary.LittleEndian.PutUint16(head, uint16(msgLen))

	msg := make([]byte, len(head)+int(msgLen))
	l := 0
	copy(msg[l:], head)
	l += len(head)
	copy(msg[l:], data)

	if DOENCRYPT {
		this.encServer.DoEncrypt(msg)
	}

	return msg, nil
}
