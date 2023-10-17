package logic

import (
	"github.com/libp2p/go-libp2p/core"
	"log/slog"
	"time"
)

type ConnectionMap map[core.PeerID]core.Conn

func NewConnectionMap() ConnectionMap {
	m := make(ConnectionMap)
	go m.oldConnectionsChecker()

	return m
}

func (m ConnectionMap) Add(conn core.Conn) {
	if oldConn, ok := m[conn.RemotePeer()]; ok {
		if err := oldConn.Close(); err != nil {
			slog.Error("error on closing old connection", err)
		}
	}

	slog.Info("new connection", slog.String("peerID", conn.RemotePeer().String()))
	m[conn.RemotePeer()] = conn
}

func (m ConnectionMap) oldConnectionsChecker() {
	ticker := time.NewTicker(time.Millisecond * 500)

	for range ticker.C {
		for peer, conn := range m {
			if conn.IsClosed() {
				delete(m, peer)
			}
		}
	}
}
