package raft

import "time"

type RaftConfig struct {
	ElectionTimeout time.Duration
	IdleTimeout     time.Duration
	// ClusterMemberAddresses list all of the cluster member address
	// take note to also include our own address in this
	ClusterMemberAddreses []string
}
