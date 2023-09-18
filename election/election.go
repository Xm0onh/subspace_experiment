package election

import "github.com/xm0onh/subspace_experiment/identity"

type Election interface {
	GetLeader() identity.NodeID
	FindLeaderFor(view int) identity.NodeID
}
