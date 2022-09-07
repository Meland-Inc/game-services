package encrypt

import (
	"fmt"
	"testing"
)

func Test_Encrypt(t *testing.T) {
	encServer := NewEncrypt(113, 126)
	encClient := NewEncrypt(113, 126)

	msg := []byte{'a', 'b', 'c', 1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Println("msg: ", msg)

	encServer.DoEncrypt(msg)
	fmt.Println("msg: ", msg)

	encClient.DoEncrypt(msg)
	fmt.Println("msg: ", msg)
}
