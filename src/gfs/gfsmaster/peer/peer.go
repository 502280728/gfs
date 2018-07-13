// peer
package peer

type PeerType uint8

const (
	Leader PeerType = 1 << iota
	Follower
)

type Peer struct {
	Type PeerType
}
