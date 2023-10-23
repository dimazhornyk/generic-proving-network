package common

import "github.com/libp2p/go-libp2p/core/peer"

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

type ProverSelectionMessage struct {
	RequestID RequestID `json:"request_id"`
	PeerID    peer.ID   `json:"peer_id"`
}

type ProofSubmissionMessage struct {
	RequestID RequestID `json:"request_id"`
	ProofID   ProofID   `json:"proof_id"`
	Proof     []byte    `json:"proof"`
}
