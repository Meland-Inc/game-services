package sdk

import (
	"bytes"
	"io/ioutil"
	"testing"

	"game-message-core/proto"

	shell "github.com/ipfs/go-ipfs-api"
)

func Test_IpfsAdd(t *testing.T) {
	testdata := &proto.BigWorldTile{
		R: 999,
		C: 999,
	}
	t.Log(testdata)
	bs, err := testdata.Marshal()
	t.Log(err)
	t.Log(bs)

	xxx := &proto.BigWorldTile{}
	err = xxx.Unmarshal(bs)
	t.Log(err)
	t.Log(xxx)

	// Where your local node is running on localhost:5001
	sh := shell.NewShell("localhost:5001")
	hash, err := sh.Add(bytes.NewBuffer(bs))
	t.Log(err)
	t.Log(hash)

	read, err := sh.Cat(hash)
	t.Log(err)
	t.Log(read)
	body, err := ioutil.ReadAll(read)

	t.Log(err)
	t.Log(body)

	vv := &proto.BigWorldTile{}
	err = vv.Unmarshal(body)
	t.Log(err)
	t.Log(vv)
}
