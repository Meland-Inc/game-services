package tcp

import (
	"net"

	"github.com/Meland-Inc/game-services/src/common/net/session"
)

type Client struct {
	addr               string
	session            *session.Session
	onReceivedCallback func(*session.Session, []byte)
	onCloseCallback    func(*session.Session)
}

func NewClient(
	addr string,
	onReceivedCallback func(*session.Session, []byte),
	onCloseCallback func(*session.Session),
) (*Client, error) {
	s := &Client{
		addr:               addr,
		onReceivedCallback: onReceivedCallback,
		onCloseCallback:    onCloseCallback,
	}

	err := s.Connect()
	return s, err
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	c.session = session.NewSession(conn)
	c.session.SetCallBack(c.onReceivedCallback, c.onCloseCallback)

	go func() {
		c.session.Run()
		c.Close()
	}()

	return nil
}

func (c *Client) Session() *session.Session {
	return c.session
}

func (c *Client) Write(data []byte) error {
	return c.session.Write(data)
}

func (c *Client) Close() {
	c.session.Stop()
}
