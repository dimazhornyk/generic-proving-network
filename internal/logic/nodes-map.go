package logic

import (
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"time"
)

type NodesMap map[peer.ID]common.NodeData

func NewNodesMap() NodesMap {
	return make(NodesMap)
}

func (m NodesMap) Add(peerID peer.ID, status common.Status, commitments []string) {
	m[peerID] = common.NodeData{
		PeerID:         peerID,
		Status:         status,
		Commitments:    commitments,
		AvailableSince: time.Now().UnixNano(),
	}
}

func (m NodesMap) UpdateStatus(peerID peer.ID, status common.Status) error {
	node, ok := m[peerID]
	if !ok {
		return errors.New("unknown peerID")
	}

	node.Status = status
	m[peerID] = node

	return nil
}

func (m NodesMap) CountNodesForConsumer(consumerName string) int {
	count := 0

	for _, node := range m {
		for _, commitment := range node.Commitments {
			if commitment == consumerName {
				count++

				break
			}
		}
	}

	return count
}
