package logic

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
)

var errUnknownVotingKey = errors.New("unknown request")
var errVotingHasDrawn = errors.New("voting has drawn")

type Voting[K, V comparable] struct {
	VotingKey K
	Votes     map[peer.ID]V
}

type VotingMap[K, V comparable] map[K]Voting[K, V]

func (m VotingMap[K, V]) Add(key K, voter peer.ID, value V) bool {
	_, ok := m[key]
	if !ok {
		m[key] = Voting[K, V]{
			VotingKey: key,
			Votes:     make(map[peer.ID]V),
		}
	}

	m[key].Votes[voter] = value

	return ok
}

func (m VotingMap[K, V]) Get(key K) (Voting[K, V], error) {
	value, ok := m[key]
	if !ok {
		return Voting[K, V]{}, errUnknownVotingKey
	}

	return value, nil
}

func (m VotingMap[K, V]) Delete(key K) {
	delete(m, key)
}

func (m VotingMap[K, V]) GetWinner(key K) (*V, error) {
	voting, ok := m[key]
	if !ok {
		return nil, errUnknownVotingKey
	}

	opts := make(map[V]int)
	for _, v := range voting.Votes {
		opts[v]++
	}

	var winner V
	var maxVotes int
	var hasEqual bool

	for v, count := range opts {
		if count == maxVotes {
			hasEqual = true
		}

		if count > maxVotes {
			maxVotes = count
			winner = v
			hasEqual = false
		}
	}

	if hasEqual {
		return nil, errVotingHasDrawn
	}

	return &winner, nil
}
