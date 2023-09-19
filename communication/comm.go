package comm

import (
	"encoding/gob"
	"fmt"
	"net"
	"net/url"

	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/log"
)

type IComm interface {
	Send(identity.NodeID, interface{})
	Recv() interface{}
	Dial() error
	Listen()
	Close()
}

type communication struct {
	uri         *url.URL
	sendAndRecv chan interface{}
	recv        chan interface{}
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
		sendAndRecv: make(chan interface{}, 10240),
		recv:        make(chan interface{}, 10240),
		close:       make(chan struct{}),
	}

	c := new(tcp)
	c.communication = communication
	return c
}

func (c *communication) Send(from identity.NodeID, m interface{}) {
	c.sendAndRecv <- m
	// c.Dial()
}

func (c *communication) Recv() (m interface{}) {
	msg, ok := <-c.recv
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
		// w := bufio.NewWriter(conn)
		// codec := NewCodec(config.Codec, conn)
		encoder := gob.NewEncoder(conn)
		defer conn.Close()
		for m := range c.sendAndRecv {
			err := encoder.Encode(&m)
			if err != nil {
				log.Error(err)
			}
		}
	}(conn)

	return nil
}

/// TCP Listener

func (t *tcp) Listen() {
	log.Debug("start listening ", t.uri.Port())

	listener, err := net.Listen("tcp", ":"+t.uri.Port())
	if err != nil {
		log.Fatal("TCP Listener error: ", err)
	}

	go func(listener net.Listener) {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Error("TCP Accept error: ", err)
				continue
			}

			go func(conn net.Conn) {
				// codec := NewCodec(config.Codec, conn)
				decoder := gob.NewDecoder(conn)
				defer conn.Close()
				//r := bufio.NewReader(conn)
				for {
					select {
					case <-t.close:
						return
					default:
						var m interface{}
						err := decoder.Decode(&m)
						if err != nil {
							log.Error(err)
							continue
						}
						t.recv <- m
					}
				}
			}(conn)

		}
	}(listener)
}
