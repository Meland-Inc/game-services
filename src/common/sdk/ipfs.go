package sdk

import (
	"bytes"
	"io/ioutil"

	ipfs "github.com/ipfs/go-ipfs-api"
)

func newIpfsShell(ipfsUrl string) *ipfs.Shell {
	return ipfs.NewShell(ipfsUrl)
}

func IpfsAddFile(ipfsUrl string, body []byte) (fileHash string, err error) {
	sh := newIpfsShell(ipfsUrl)
	return sh.Add(bytes.NewBuffer(body))
}

func IpfsCatFile(ipfsUrl string, fileHash string) ([]byte, error) {
	sh := newIpfsShell(ipfsUrl)
	read, err := sh.Cat(fileHash)
	if err != nil {
		return []byte{}, err
	}
	return ioutil.ReadAll(read)
}
