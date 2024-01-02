package logic

import (
	"github.com/libp2p/go-libp2p/core"
	"log/slog"
	"sync"
	"time"
)

type ConnectionHolder struct {
	m  map[core.PeerID]core.Conn
	mu sync.RWMutex
}

func NewConnectionHolder() *ConnectionHolder {
	ch := new(ConnectionHolder)
	go ch.oldConnectionsChecker()

	return ch
}

func (c *ConnectionHolder) Add(conn core.Conn) {
	c.mu.Lock()
	oldConn, ok := c.m[conn.RemotePeer()]
	c.m[conn.RemotePeer()] = conn
	c.mu.Unlock()

	slog.Info("new connection", slog.String("peerID", conn.RemotePeer().String()))

	if ok {
		if err := oldConn.Close(); err != nil {
			slog.Error("error on closing old connection", err)
		}
	}
}

func (c *ConnectionHolder) oldConnectionsChecker() {
	ticker := time.NewTicker(time.Millisecond * 500)

	for range ticker.C {
		c.mu.Lock()
		for peer, conn := range c.m {
			if conn.IsClosed() {
				delete(c.m, peer)
			}
		}
		c.mu.Unlock()
	}
}

func (c *ConnectionHolder) GetPeerIDs() []core.PeerID {
	c.mu.RLock()
	defer c.mu.RUnlock()

	peerIDs := make([]core.PeerID, 0, len(c.m))
	for peerID := range c.m {
		peerIDs = append(peerIDs, peerID)
	}

	return peerIDs
}

func (c *ConnectionHolder) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.m)
}
