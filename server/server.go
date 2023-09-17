package main

import (
	"fmt"
	"sync"

	"github.com/xm0onh/subspace_experiment"
	"github.com/xm0onh/subspace_experiment/config"
	"github.com/xm0onh/subspace_experiment/identity"
	"github.com/xm0onh/subspace_experiment/operator"
)

func initOperator(id identity.NodeID) {
	fmt.Println("init operator", id)
	o := operator.NewOperator(id)
	o.Start()
}
func main() {
	subspace_experiment.Init()
	var wg sync.WaitGroup
	wg.Add(1)
	for id := range config.GetConfig().Addrs {
		fmt.Println(id)
		go initOperator(id)
	}
	wg.Wait()
	// log.Infof("node %v starting...", 1)
}
