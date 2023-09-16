package config

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/xm0onh/subspace_experiment/identity"
)

var configFile = flag.String("config", "config.json", "Configuration file for Subspace Experiment config.json.")

type Config struct {
	Addrs     map[identity.NodeID]string `json:"addrs"`
	HTTPAddrs map[identity.NodeID]string `json:"http_addrs"`

	n int // total number of nodes
}

func (c *Config) Load() {
	file, err := os.Open(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		log.Fatal(err)
	}

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

	c.n = len(c.Addrs)
}
