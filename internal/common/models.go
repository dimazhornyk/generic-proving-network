package common

import "github.com/libp2p/go-libp2p/core/peer"

type CalculateProofRequest struct {
	ID              string
	ConsumerName    string
	ConsumerAddress string
	Signature       []byte // signature has to be done of the requestID
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

type Consumer struct {
	Name    string
	Address string
	Image   string
}

type ValidationSignature struct {
	PeerID    peer.ID
	Signature []byte
}
