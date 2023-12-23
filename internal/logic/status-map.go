package logic

import (
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"time"
)

type StatusMap map[peer.ID]common.NodeData

func NewStatusMap() StatusMap {
	return make(StatusMap)
}

func (m StatusMap) Add(peerID peer.ID, status common.Status, commitments []string) {
	m[peerID] = common.NodeData{
		PeerID:         peerID,
		Status:         status,
		Commitments:    commitments,
		AvailableSince: time.Now().UnixNano(),
	}
}

func (m StatusMap) UpdateStatus(peerID peer.ID, status common.Status) error {
	node, ok := m[peerID]
	if !ok {
		return errors.New("unknown peerID")
	}

	node.Status = status
	m[peerID] = node

	return nil
}
