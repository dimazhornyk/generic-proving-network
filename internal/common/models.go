package common

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"math/big"
)

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
	Address ethcommon.Address
	Balance *big.Int
	Image   string
}

type ValidationSignature struct {
	PeerID    peer.ID
	Signature []byte
}

type RequestExtension struct {
	ProvingRequestMessage
	ProvingPeers         []peer.ID
	Proofs               map[peer.ID]ZKProof
	ValidationSignatures [][]byte
}
