package comm

import (
	"fmt"
	"net/url"

	"github.com/xm0onh/subspace_experiment/log"
)

type IComm interface {
	Send([]string)
	Recv() []string
	// Dial() error
	// Listen()
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
	fmt.Println("Message Broadcasted")

}

func (c *communication) Recv() (msg []string) {
	msg = <-c.sendAndRecv
	fmt.Println("Message Received", msg)
	return msg
}

func (c *communication) Close() {
	close(c.close)
}

// func (c *communication) Dial() error {
// 	fmt.Println("dialing ", c.uri.Host)
// 	conn, err := net.Dial("tcp", c.uri.Host)
// 	if err != nil {
// 		return err
// 	}
// 	go func(conn net.Conn) {
// 		encode := gob.NewEncoder(conn)
// 		defer conn.Close()
// 		for m := range c.sendAndRecv {
// 			fmt.Println(m)
// 			err := encode.Encode(&m)
// 			if err != nil {
// 				log.Fatal("error encoding message: ", err)
// 			}
// 		}
// 	}(conn)

// 	return nil
// }

// /// TCP Listener

// func (t *tcp) Listen() {
// 	log.Debug("listening on ", t.uri.Port())
// 	_, err := net.Listen("tcp", ":"+t.uri.Port())
// 	if err != nil {
// 		log.Fatal("TCP error listening: ", err)
// 	}
// }
