package common

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"math/big"
)

// maybe move to protobuf later
type Status int

const (
	StatusInit Status = iota
	StatusIdle
	StatusProving
	StatusShuttingDown
)

func (s Status) String() string {
	return [...]string{"StatusInit", "StatusIdle", "StatusProving", "StatusShuttingDown"}[s]
}

type StatusMessage struct {
	Status  Status `json:"status"`
	Payload any    `json:"payload"`
}

type ProvingRequestMessage struct {
	ID              RequestID `json:"request_id"`
	Reward          *big.Int  `json:"reward"`
	ConsumerName    string    `json:"consumer_name"`
	ConsumerAddress string    `json:"consumer_address"`
	Signature       []byte    `json:"signature"`
	Data            []byte    `json:"data"`
	Timestamp       int64     `json:"timestamp"`
}

type VotingMessageType int

const (
	VoteProverSelection = iota
	VoteValidation
)

type VotingMessage struct {
	Type    VotingMessageType `json:"type"`
	Payload any               `json:"payload"`
}

type Topic string

const (
	GlobalTopic   Topic = "global"
	RequestsTopic Topic = "requests"
	VotingTopic   Topic = "voting"
	ProofsTopic   Topic = "proofs"
)

func (t Topic) String() string {
	return string(t)
}

type ProverSelectionPayload struct {
	RequestID RequestID `json:"request_id"`
	PeerID    peer.ID   `json:"peer_id"`
}

type ProofSubmissionMessage struct {
	RequestID RequestID `json:"request_id"`
	ProofID   ProofID   `json:"proof_id"`
	Proof     []byte    `json:"proof"`
}

type ValidationPayload struct {
	RequestID           RequestID `json:"request_id"`
	ProverID            peer.ID   `json:"prover_id"`
	IsValid             bool      `json:"is_valid"`
	ValidationTimestamp int64     `json:"validation_timestamp,omitempty"`
	Signature           []byte    `json:"signature,omitempty"`
}

type DataToSign struct {
	RequestID     RequestID `json:"request_id"`
	ProverAddress string    `json:"prover_address"`
	IsValid       bool      `json:"is_valid"`
}
