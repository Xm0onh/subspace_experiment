package identity

import (
	"strconv"

	"github.com/xm0onh/subspace_experiment/log"
)

type NodeID string
type IDs []NodeID

func NewNodeID(node int) NodeID {
	if node < 0 {
		node = -node
	}
	return NodeID(strconv.Itoa(node))
}

func (i NodeID) Node() int {
	var s string = string(i)
	node, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Errorf("Failed to convert Node %s to int\n", s)
		return 0
	}
	return int(node)
}

func (a IDs) Len() int      { return len(a) }
func (a IDs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
