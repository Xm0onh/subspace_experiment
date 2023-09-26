package election

import (
	"crypto/sha1"
	"encoding/binary"
	"strconv"

	"github.com/xm0onh/subspace_experiment/identity"
)

type Rotation struct {
	peerNo int
}

func NewRotation(peerNo int) *Rotation {
	return &Rotation{
		peerNo: peerNo,
	}
}

func (r *Rotation) IsLeader(id identity.NodeID, view int) bool {
	if view <= 3 {
		if id.Node() < r.peerNo {
			return false
		}
		return true
	}
	h := sha1.New()
	h.Write([]byte(strconv.Itoa(int(view) + 1)))
	bs := h.Sum(nil)
	data := binary.BigEndian.Uint64(bs)
	return data%uint64(r.peerNo) == uint64(id.Node()-1)
}

func (r *Rotation) FindLeaderFor(view int) identity.NodeID {
	// if view <= 3 {
	// 	return identity.NewNodeID(r.peerNo)
	// }
	// h := sha1.New()
	// h.Write([]byte(strconv.Itoa(int(view + 1))))
	// bs := h.Sum(nil)
	// data := binary.BigEndian.Uint64(bs)
	// id := data%uint64(r.peerNo) + 1

	// if rand.Float64() <= r.P {
	// 	return identity.NewNodeID(int(id))
	// }

	// nextID := (id + 1) % uint64(r.peerNo)
	// if nextID == 0 {
	// 	nextID = uint64(r.peerNo)
	// }
	// return identity.NewNodeID(int(nextID))
	if view <= 3 {
		return identity.NewNodeID(r.peerNo)
	}
	h := sha1.New()
	h.Write([]byte(strconv.Itoa(int(view + 1))))
	bs := h.Sum(nil)
	data := binary.BigEndian.Uint64(bs)
	id := data%uint64(r.peerNo) + 1
	// id := rand.Intn(r.peerNo)

	return identity.NewNodeID(int(id))
}
