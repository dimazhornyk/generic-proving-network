package logic

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
)

type Voting[K, V comparable] struct {
	VotingKey K
	Votes     map[peer.ID]V
}

type VotingMap[K, V comparable] map[K]Voting[K, V]

func NewVotingMap[K, V comparable]() VotingMap[K, V] {
	return make(VotingMap[K, V])
}

func (m VotingMap[K, V]) Add(key K, voter peer.ID, value V) {
	if _, ok := m[key]; !ok {
		m[key] = Voting[K, V]{
			VotingKey: key,
			Votes:     make(map[peer.ID]V),
		}
	}

	m[key].Votes[voter] = value
}

func (m VotingMap[K, V]) Get(key K) (Voting[K, V], error) {
	value, ok := m[key]
	if !ok {
		return Voting[K, V]{}, errors.New("unknown voting key")
	}

	return value, nil
}

func (m VotingMap[K, V]) Delete(key K) {
	delete(m, key)
}

func (m VotingMap[K, V]) GetWinner(key K) (*V, error) {
	voting, ok := m[key]
	if !ok {
		return nil, errors.New("unknown voting key")
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
		return nil, errors.New("voting has drawn")
	}

	return &winner, nil
}
