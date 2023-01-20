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
	Clients map[*websocket.Conn]bool
}

func NewCanvas() *Canvas {
	return &Canvas{
		Clients: make(map[*websocket.Conn]bool, 0),
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
	c.Clients[conn] = true
	c.mu.Unlock()
}

func (c *Canvas) remove(conn *websocket.Conn) {
	log.Println("remove client", conn.RemoteAddr())
	c.mu.Lock()
	delete(c.Clients, conn)
	c.mu.Unlock()
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
	for conn := range c.Clients {
		go func(client *websocket.Conn) {
			if _, err := client.Write(msg); err != nil {
				log.Println("error writing message", err.Error())
			}
		}(conn)
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
