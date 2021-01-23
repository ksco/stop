package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Client struct {
	conn    *websocket.Conn
	receive chan string
	send    chan string
	done    chan bool
}

func (c Client) read() {
	defer func() {
		close(c.done)
		close(c.receive)
		_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		_ = c.conn.Close()
	}()

	messageType, message, err := c.conn.ReadMessage()
	if err != nil || messageType != websocket.TextMessage {
		return
	}
	c.receive <- string(message)

	for {
		messageType, _, err := c.conn.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			return
		}
		continue
	}
}

func (c Client) write() {
	for {
		select {
		case message := <-c.send:
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				return
			}
		case <-c.done:
			return
		}
	}
}

func (c Client) ReceiveChan() <-chan string {
	return c.receive
}

func (c Client) SendChan() chan<- string {
	return c.send
}

func (c Client) Close() {
	_ = c.conn.Close()
}

func (c Client) DoneChan() chan bool {
	return c.done
}

func NewClient(w http.ResponseWriter, r *http.Request) (*Client, error) {
	conn, err := ws.Conn(w, r)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:    conn,
		receive: make(chan string),
		send:    make(chan string),
		done:    make(chan bool),
	}

	go client.read()
	go client.write()
	return client, nil
}
