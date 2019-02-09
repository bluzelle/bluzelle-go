package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Conn represents the persistent layer for Bluzelle.
type Conn struct {
	// Endpoint represents the entry point for the bluzelle network.
	Endpoint string

	// IncomingMsg
	IncomingMsg chan []byte

	conn *websocket.Conn
}

// NewConn creates a new conn
func NewConn(endpoint string) *Conn {
	return &Conn{
		Endpoint:    endpoint,
		IncomingMsg: make(chan []byte),
		conn:        nil,
	}
}

func (conn *Conn) Dial() error {
	u := url.URL{Scheme: "ws", Host: conn.Endpoint}
	log.Println("Connecting to: ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	conn.conn = c
	go func() {
		for {
			_, r, err := c.ReadMessage()
			if err != nil {
				log.Print(err)
			}
			conn.IncomingMsg <- r
		}
	}()

	return nil
}

func (conn *Conn) ReadMsg() <-chan []byte {
	return conn.IncomingMsg
}

func (conn *Conn) SendMsg(data []byte) error {
	return conn.conn.WriteMessage(websocket.TextMessage, data)
}
