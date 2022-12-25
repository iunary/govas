package main

import (
	"io"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	bufferSize = 1024
)

type Canvas struct {
	mu      sync.Mutex
	Clients []*websocket.Conn
}

func NewCanvas() *Canvas {
	return &Canvas{
		Clients: make([]*websocket.Conn, 0),
	}
}

func (c *Canvas) handleWS(conn *websocket.Conn) {
	log.Println("new client joined", conn.RemoteAddr())
	defer conn.Close()
	c.join(conn)
	c.reader(conn)
}

func (c *Canvas) join(conn *websocket.Conn) {
	c.mu.Lock()
	c.Clients = append(c.Clients, conn)
	c.mu.Unlock()
}

func (c *Canvas) remove(conn *websocket.Conn) {
	log.Println("remove client", conn.RemoteAddr())
	for index, client := range c.Clients {
		if client == conn {
			c.Clients = append(c.Clients[:index], c.Clients[index+1:]...)
		}
	}
}

func (c *Canvas) reader(conn *websocket.Conn) {
	buf := make([]byte, bufferSize)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("error", err.Error())
			if err == io.EOF {
				c.remove(conn)
				break
			}
			c.remove(conn)
			continue
		}
		msg := buf[:n]
		log.Println("recv", string(msg))
		c.broadcast(msg)
	}
}

func (c *Canvas) broadcast(msg []byte) {
	for _, client := range c.Clients {
		go func(client *websocket.Conn) {
			if _, err := client.Write(msg); err != nil {
				log.Println("error writing message", err.Error())
			}
		}(client)
	}
}

func main() {
	canvas := NewCanvas()
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.Handle("/ws", websocket.Handler(canvas.handleWS))
	log.Println("running on http://0.0.0.0:4444")
	log.Panic(http.ListenAndServe(":4444", nil))
}
