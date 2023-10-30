package common

import "github.com/libp2p/go-libp2p/core/peer"

type CalculateProofRequest struct {
	ConsumerName    string
	ConsumerAddress string
	Signature       []byte
	Data            []byte
}

type NodeData struct {
	PeerID           peer.ID
	Status           Status
	CurrentRequestID *RequestID
	Commitments      []string
	AvailableSince   int64
}

type ZKProof struct {
	ProofID   ProofID
	Proof     []byte
	Timestamp int64
}

type RequestID = string
type ProofID = string

type Container struct {
	ID         string
	SourcePort string
}
