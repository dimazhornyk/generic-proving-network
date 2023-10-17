package logic

import (
	"context"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/pkg/errors"
	"log/slog"
	"multi-proving-client/internal/common"
)

const connectivityFactor = 3

type Discovery struct {
	host        host.Host
	dht         *dht.IpfsDHT
	namespace   string
	protocolID  core.ProtocolID
	connections ConnectionMap
}

func NewDiscovery(host host.Host, dht *dht.IpfsDHT, cfg common.Config) *Discovery {
	return &Discovery{
		host:        host,
		dht:         dht,
		namespace:   cfg.Namespace,
		protocolID:  cfg.ProtocolID,
		connections: NewConnectionMap(),
	}
}

func (d *Discovery) Start(ctx context.Context) error {
	discovery := routing.NewRoutingDiscovery(d.dht)
	util.Advertise(ctx, discovery, d.namespace)

	peersCh, err := discovery.FindPeers(ctx, d.namespace)
	if err != nil {
		return errors.Wrap(err, "error from discovery find peers")
	}

	go d.listen(ctx, peersCh)

	return nil
}

func (d *Discovery) listen(ctx context.Context, ch <-chan peer.AddrInfo) {
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-ch:
			if p.ID == d.host.ID() {
				continue
			}

			slog.Info("Found peer", slog.String("peerID", p.ID.String()))
			if d.host.Network().Connectedness(p.ID) != network.Connected && len(d.connections) < connectivityFactor {
				conn, err := d.host.Network().DialPeer(ctx, p.ID)
				if err != nil {
					slog.Error("error on dialing peer", err, slog.String("peerID", p.ID.String()))
					continue
				}
				slog.Info("Connected to:", p.ID.String())

				d.connections.Add(conn)
			}
		}
	}
}
