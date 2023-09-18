package election

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/xm0onh/subspace_experiment/identity"
)

type Rotation struct {
	peerNo int
	P      float64
	rng    *rand.Rand
	leader identity.NodeID
}

func NewRotation(peerNo int, P float64) *Rotation {
	src := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(src)

	return &Rotation{
		peerNo: peerNo,
		P:      P,
		rng:    rand,
	}
}

func (r *Rotation) GetLeader() identity.NodeID {
	fmt.Println(r.leader)
	return r.leader
}

func (r *Rotation) FindLeaderFor(view int) identity.NodeID {
	if view <= 3 {
		return identity.NewNodeID(r.peerNo)
	}
	h := sha1.New()
	h.Write([]byte(strconv.Itoa(int(view + 1))))
	bs := h.Sum(nil)
	data := binary.BigEndian.Uint64(bs)
	id := data%uint64(r.peerNo) + 1

	if rand.Float64() <= r.P {
		return identity.NewNodeID(int(id))
	}

	nextID := (id + 1) % uint64(r.peerNo)
	if nextID == 0 {
		nextID = uint64(r.peerNo)
	}
	return identity.NewNodeID(int(nextID))
}
