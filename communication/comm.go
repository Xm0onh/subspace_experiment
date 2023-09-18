package comm

import (
	"encoding/gob"
	"fmt"
	"net"
	"net/url"

	"github.com/xm0onh/subspace_experiment/log"
)

type IComm interface {
	Send([]string)
	Recv() []string
	Dial() error
	Listen()
	Close()
}

type communication struct {
	uri         *url.URL
	sendAndRecv chan []string
	recv        chan []string
	close       chan struct{}
}

type tcp struct {
	*communication
}

func NewComm(addr string) IComm {
	uri, err := url.Parse(addr)
	if err != nil {
		log.Fatalf("error parsing address $s : $s\n", addr, err)
	}

	communication := &communication{
		uri:         uri,
		sendAndRecv: make(chan []string, 10240),
		recv:        make(chan []string, 10240),
		close:       make(chan struct{}),
	}

	c := new(tcp)
	c.communication = communication
	return c
}

func (c *communication) Send(msg []string) {
	c.sendAndRecv <- msg
	c.Dial()
}

func (c *communication) Recv() (msg []string) {
	msg, ok := <-c.recv
	fmt.Println("Message Received", ok)
	if !ok {
		return nil
	}
	return msg
}

func (c *communication) Close() {
	close(c.close)
}

func (c *communication) Dial() error {
	fmt.Println("dialing ", c.uri.Host)
	conn, err := net.Dial("tcp", c.uri.Host)
	if err != nil {
		return err
	}
	go func(conn net.Conn) {
		encode := gob.NewEncoder(conn)
		defer conn.Close()
		for m := range c.sendAndRecv {
			fmt.Println(m)
			err := encode.Encode(&m)
			if err != nil {
				log.Fatal("error encoding message: ", err)
			}
		}
	}(conn)

	return nil
}

/// TCP Listener

func (t *tcp) Listen() {
	log.Debug("listening on ", t.uri.Port())

	listener, err := net.Listen("tcp", ":"+t.uri.Port())
	if err != nil {
		log.Fatal("TCP error listening: ", err)
	}

	go func(listener net.Listener) {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Error("TCP error accepting: ", err) // Changed from Fatal to Error so the server doesn't crash
				continue
			}
			go func(conn net.Conn) {
				defer conn.Close()
				var msg []string
				decoder := gob.NewDecoder(conn)
				err := decoder.Decode(&msg)
				if err != nil {
					log.Error("Error decoding message: ", err)
					return
				}
				t.recv <- msg
			}(conn)
		}
	}(listener)
}
