package election

import "github.com/xm0onh/subspace_experiment/identity"

type Election interface {
	IsLeader(id identity.NodeID, view int) bool
	FindLeaderFor(view int) identity.NodeID
}
