package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/xm0onh/subspace_experiment/identity"
)

type Config struct {
	Addrs     map[identity.NodeID]string `json:"addrs"`
	HTTPAddrs map[identity.NodeID]string `json:"http_addrs"`

	n int // total number of nodes
}

var Configuration Config

func (c *Config) Load() {

	c.Addrs = make(map[identity.NodeID]string)
	c.HTTPAddrs = make(map[identity.NodeID]string)
	// load ips
	ip_file, err := os.Open("ips.txt")
	if err != nil {
		fmt.Println(err)
	}

	defer ip_file.Close()

	scanner := bufio.NewScanner(ip_file)
	i := 1
	for scanner.Scan() {
		id := identity.NewNodeID(i)
		port := strconv.Itoa(3734 + i)
		addr := "tcp://" + scanner.Text() + ":" + port
		portHttp := strconv.Itoa(8069 + i)
		addrHttp := "http://" + scanner.Text() + ":" + portHttp
		c.Addrs[id] = addr
		c.HTTPAddrs[id] = addrHttp
		i++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Addrs", c.Addrs)
	c.n = len(c.Addrs)
}
func GetConfig() Config {
	return Configuration
}
