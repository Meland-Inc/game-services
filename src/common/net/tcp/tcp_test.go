package tcp

import (
	"fmt"
	"testing"
	"time"

	"github.com/Meland-Inc/game-services/src/common/net/session"
)

func ServerOnConnectCallback(s *session.Session) {
	fmt.Printf("session[%s][%s] ServerOnConnectCallback  ------\n", s.SessionId(), s.RemoteAddr())
	s.SetCallBack(
		ServerOnReceivedCallback,
		ServerOnCloseCallback,
	)
}
func ServerOnReceivedCallback(s *session.Session, data []byte) {
	fmt.Printf(
		"session[%s][%s]  ServerOnReceivedCallback data length[%d], dataStr:[%v]\n",
		s.SessionId(), s.RemoteAddr(), len(data), string(data),
	)
	s.Write([]byte("server response message is ServerOnReceivedCallback ########"))
}
func ServerOnCloseCallback(s *session.Session) {
	fmt.Printf("session[%s][%s]  ServerOnCloseCallback ======= \n", s.SessionId(), s.RemoteAddr())
}

func Test_TcpServer(t *testing.T) {
	tcpServer, err := NewTcpServer(
		":7659",
		100,
		180,
		ServerOnConnectCallback,
	)
	fmt.Println(err)
	fmt.Println(tcpServer)
	fmt.Printf("tcp server started  \n")
	time.Sleep(1 * time.Hour)
}

func ClientOnReceivedCallback(s *session.Session, data []byte) {
	fmt.Printf(
		"session[%s][%s]  ClientOnReceivedCallback data length[%d], dataStr:[%v]\n",
		s.SessionId(), s.RemoteAddr(), len(data), string(data),
	)
}
func ClientOnCloseCallback(s *session.Session) {
	fmt.Printf("session[%s][%s]  ClientOnCloseCallback ======= \n", s.SessionId(), s.RemoteAddr())
}
func Test_TcpClient(t *testing.T) {
	t.Log("TcpClient ------- begin --------")
	cli, err := NewClient(":7659", ClientOnReceivedCallback, ServerOnCloseCallback)
	if err != nil {
		panic(err)
	}

	t.Log("TcpClient ------- send --------")
	cli.Write([]byte("this is a test message from game services gcp client"))

	t.Log("TcpClient ------- send over --------")
	time.Sleep(30 * time.Second)
	cli.Close()
}
