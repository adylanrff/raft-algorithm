package raft

func GetRaftServerIDFromAddress(address string) string {
	// TODO: Implement some other server id / hash alg to generate server id
	// for now just use the address
	return address
}
