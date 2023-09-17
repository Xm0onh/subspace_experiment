package election

import (
	"crypto/sha1"
	"encoding/binary"
	"math/rand"
	"strconv"
	"time"

	"github.com/xm0onh/subspace_experiment/identity"
)

type Rotation struct {
	peerNo int
	P      float64
}

func NewRotation(peerNo int, P float64) *Rotation {
	src := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(src)
	_ = rand.Float64()
	return &Rotation{
		peerNo: peerNo,
		P:      P,
	}
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
