package sync

import (
	"bufio"
	"context"
	"fmt"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"log/slog"
	"time"
)

const delim = '\n'
const minConnectionsForSync = 3

type InitialSyncer struct {
	host        host.Host
	protocolID  core.ProtocolID
	storage     *logic.Storage
	connections *logic.ConnectionHolder
}

type chanResp struct {
	peerID string
	hash   string
	err    error
}

func NewInitialSyncer(cfg *common.Config, connections *logic.ConnectionHolder, storage *logic.Storage, host host.Host) *InitialSyncer {
	return &InitialSyncer{
		host:        host,
		storage:     storage,
		connections: connections,
		protocolID:  cfg.ProtocolID,
	}
}

func (is *InitialSyncer) Sync(ctx context.Context) error {
	ticker := time.NewTicker(time.Millisecond * 500)
	timeout := time.NewTimer(time.Second * 10)

outer:
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timeout.C:
			break outer
		case <-ticker.C:
			if is.connections.Len() >= minConnectionsForSync {
				break outer
			}
		}
	}

	if is.connections.Len() == 0 { // first node in the network, no need to sync
		slog.Info("first node, no need to sync")

		return nil
	}

	peers := is.connections.GetPeerIDs()
	for _, p := range peers {
		if err := is.syncStorage(ctx, p); err != nil {
			slog.Error("error on syncing storage", slog.String("peerID", p.String()), slog.String("error", err.Error()))

			continue
		} else {
			if err := is.checkStateChecksum(ctx, peers); err != nil {
				return errors.Wrap(err, "error checking state checksum")
			}
			slog.Info("state checksum is correct")

			return nil
		}
	}

	return errors.New("unable to sync storage")
}

func (is *InitialSyncer) syncStorage(ctx context.Context, p peer.ID) error {
	stream, err := is.host.NewStream(ctx, p, is.protocolID)
	if err != nil {
		return errors.Wrap(err, "error on creating a new stream")
	}

	if err := is.sendInitMessage(bufio.NewWriter(stream)); err != nil {
		return errors.Wrap(err, "error sending init message")
	}

	reader := bufio.NewReader(stream)
	if err := is.readRequestsData(reader); err != nil {
		return errors.Wrap(err, "error reading requests data")
	}

	if err := is.readLatestProofsData(reader); err != nil {
		return errors.Wrap(err, "error reading latest proofs data")
	}

	if err := stream.Close(); err != nil {
		return errors.Wrap(err, "error closing a stream")
	}

	slog.Info("synced with a peer", slog.String("peerID", p.String()))

	return nil
}

func (is *InitialSyncer) checkStateChecksum(ctx context.Context, peers []peer.ID) error {
	streams := make([]network.Stream, 0, len(peers))
	for _, p := range peers {
		stream, err := is.host.NewStream(ctx, p, is.protocolID)
		if err != nil {
			return errors.Wrapf(err, "error on creating a new stream, peerID: %s", p.String())
		}

		streams = append(streams, stream)
	}

	defer func() {
		for _, s := range streams {
			if err := s.Close(); err != nil {
				slog.Error("error closing a stream", slog.String("error", err.Error()))
			}
		}
	}()

	storageHash, err := is.storage.GetStorageHash()
	if err != nil {
		return errors.Wrap(err, "error getting storage hash")
	}

	if err := is.requestHashes(streams); err != nil {
		return errors.Wrap(err, "error requesting hashes")
	}

	ch := make(chan chanResp, len(streams))
	for _, s := range streams {
		s := s // avoid capturing loop variable
		go is.listenForStorageHash(s, ch)
	}

	for i := 0; i < len(streams); i++ {
		resp := <-ch
		if resp.err != nil {
			return resp.err
		}

		if resp.hash != storageHash {
			return fmt.Errorf("state checksum is not matching, peerID: %s", resp.peerID)
		}
	}

	return nil
}

func (is *InitialSyncer) requestHashes(streams []network.Stream) error {
	msg := Message{Type: RequestStorageHash}
	b, err := common.GobEncodeMessage(msg)
	if err != nil {
		return errors.Wrap(err, "error encoding a message")
	}

	b = append(b, delim)
	for _, stream := range streams {
		wr := bufio.NewWriter(stream)
		if _, err := wr.Write(b); err != nil {
			return errors.Wrap(err, "error writing to a stream")
		}

		if err := wr.Flush(); err != nil {
			return errors.Wrap(err, "error flushing a stream")
		}
	}

	return nil
}

func (is *InitialSyncer) listenForStorageHash(stream network.Stream, ch chan chanResp) {
	reader := bufio.NewReader(stream)
	jobResponse := chanResp{
		peerID: stream.Conn().RemotePeer().String(),
	}

	b, err := reader.ReadBytes(delim)
	if err != nil {
		jobResponse.err = errors.Wrap(err, "error reading hash from a stream")
		ch <- jobResponse

		return
	}

	var resp Message
	if err := common.GobDecodeMessage(b, &resp); err != nil {
		jobResponse.err = errors.Wrap(err, "error decoding a message")
		ch <- jobResponse

		return
	}

	hash := resp.Payload.(string)
	jobResponse.hash = hash
	ch <- jobResponse
}

func (is *InitialSyncer) sendInitMessage(wr *bufio.Writer) error {
	msg := Message{Type: InitSync}
	b, err := common.GobEncodeMessage(msg)
	if err != nil {
		return errors.Wrap(err, "error encoding a message")
	}

	b = append(b, delim)
	if _, err := wr.Write(b); err != nil {
		return errors.Wrap(err, "error writing to a stream")
	}

	if err := wr.Flush(); err != nil {
		return errors.Wrap(err, "error flushing a stream")
	}

	return nil
}

func (is *InitialSyncer) readRequestsData(r *bufio.Reader) error {
	b, err := r.ReadBytes(delim)
	if err != nil {
		return errors.Wrap(err, "error reading requests from a stream")
	}

	var resp Message
	if err := common.GobDecodeMessage(b, &resp); err != nil {
		return errors.Wrap(err, "error decoding a message")
	}

	requests, ok := resp.Payload.(map[common.RequestID]common.RequestExtension)
	if !ok {
		return errors.New("error decoding requests data")
	}

	is.storage.SetRequests(requests)

	return nil
}

func (is *InitialSyncer) readLatestProofsData(r *bufio.Reader) error {
	b, err := r.ReadBytes(delim)
	if err != nil {
		return errors.Wrap(err, "error reading requests from a stream")
	}

	var resp Message
	if err := common.GobDecodeMessage(b, &resp); err != nil {
		return errors.Wrap(err, "error decoding a message")
	}

	latestProofs, ok := resp.Payload.(map[string]common.ZKProof)
	if !ok {
		return errors.New("error decoding latest proofs data")
	}

	is.storage.SetLatestProofs(latestProofs)

	return nil
}

func (is *InitialSyncer) ProvideData() {
	is.host.SetStreamHandler(is.protocolID, func(stream network.Stream) {
		slog.Info("new stream", slog.String("peerID", stream.Conn().RemotePeer().String()))
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

		for {
			data, err := rw.ReadBytes(delim)
			if err != nil {
				slog.Error("error reading from a stream", slog.String("error", err.Error()))

				return
			}

			var msg Message
			if err := common.GobDecodeMessage(data, &msg); err != nil {
				slog.Error("error decoding a message", slog.String("error", err.Error()))

				return
			}

			switch msg.Type {
			case InitSync:
				err = is.shareStorage(rw)
			case RequestStorageHash:
				err = is.sendStorageHash(rw)
			default:
				slog.Error("unknown message type", slog.Int("type", int(msg.Type)))

				return
			}

			if err != nil {
				slog.Error("error sending a message", slog.String("error", err.Error()))

				return
			}
		}
	})
}

func (is *InitialSyncer) shareStorage(rw *bufio.ReadWriter) error {
	requests := is.storage.GetRequests()
	proofs := is.storage.GetLatestProofs()

	msg := Message{
		Type:    SendData,
		Payload: requests,
	}

	b, err := common.GobEncodeMessage(msg)
	if err != nil {
		return errors.Wrap(err, "error encoding a message")
	}

	b = append(b, delim)
	if _, err := rw.Write(b); err != nil {
		return errors.Wrap(err, "error writing to a stream")
	}

	if err := rw.Flush(); err != nil {
		return errors.Wrap(err, "error flushing a stream")
	}

	msg = Message{
		Type:    SendData,
		Payload: proofs,
	}

	b, err = common.GobEncodeMessage(msg)
	if err != nil {
		return errors.Wrap(err, "error encoding a message")
	}

	b = append(b, delim)
	if _, err := rw.Write(b); err != nil {
		return errors.Wrap(err, "error writing to a stream")
	}

	if err := rw.Flush(); err != nil {
		return errors.Wrap(err, "error flushing a stream")
	}

	return nil
}

func (is *InitialSyncer) sendStorageHash(rw *bufio.ReadWriter) error {
	hash, err := is.storage.GetStorageHash()
	if err != nil {
		return errors.Wrap(err, "error getting storage hash")
	}

	msg := Message{
		Type:    SendStorageHash,
		Payload: hash,
	}

	b, err := common.GobEncodeMessage(msg)
	if err != nil {
		return errors.Wrap(err, "error encoding a message")
	}

	b = append(b, delim)
	if _, err := rw.Write(b); err != nil {
		return errors.Wrap(err, "error writing to a stream")
	}

	if err := rw.Flush(); err != nil {
		return errors.Wrap(err, "error flushing a stream")
	}

	return nil
}
