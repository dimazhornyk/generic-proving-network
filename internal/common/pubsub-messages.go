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
	Status  Status
	Payload any
}

type ProvingRequestMessage struct {
	ID              RequestID
	ConsumerName    string
	ConsumerAddress string
	Signature       []byte
	Data            []byte
	Timestamp       int64
}

type VotingMessageType int

const (
	InitProverSelectionVoting = iota
	VoteProverSelection
	InitValidationVoting
	VoteValidation
)

type VotingMessage struct {
	Type    VotingMessageType
	Payload any
}

type Topic string

const (
	GlobalTopic   Topic = "global"
	RequestsTopic Topic = "requests"
	VotingTopic   Topic = "voting"
)

func (t Topic) String() string {
	return string(t)
}

type ProverSelectionMessage struct {
	RequestID RequestID
	PeerID    peer.ID
}
