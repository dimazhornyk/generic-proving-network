package handlers

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/dimazhornyk/generic-proving-network/internal/common"
	"github.com/dimazhornyk/generic-proving-network/internal/connectors"
	"github.com/dimazhornyk/generic-proving-network/internal/logic"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"log/slog"
	"time"
)

type ProofsHandler struct {
	host     host.Host
	nodesMap logic.StatusMap
	storage  *logic.Storage
	service  *logic.Service
	pubsub   *connectors.PubSub
	key      *ecdsa.PrivateKey
}

func NewProofsHandler(key *ecdsa.PrivateKey, host host.Host, storage *logic.Storage, service *logic.Service, pubsub *connectors.PubSub, nodesMap logic.StatusMap) *ProofsHandler {
	return &ProofsHandler{
		host:     host,
		storage:  storage,
		service:  service,
		pubsub:   pubsub,
		key:      key,
		nodesMap: nodesMap,
	}
}

func (h *ProofsHandler) Handle(ctx context.Context, peerID peer.ID, msg common.ProofSubmissionMessage) {
	if peerID == h.host.ID() {
		return // no need to verify the proof that we just generated
	}

	reqData, err := h.storage.GetProvingRequestByID(msg.RequestID)
	if err != nil {
		slog.Error("error getting proving request data", slog.String("err", err.Error()))

		return
	}

	if err := h.storage.AddProof(msg.RequestID, peerID, msg.ProofID, msg.Proof); err != nil {
		slog.Error("error adding proof to storage", slog.String("err", err.Error()))

		return
	}

	// TODO: check if the node's status was Proving, so the state transition was correct

	valid, err := h.service.ValidateProof(msg.RequestID, reqData.ConsumerImage, reqData.Data, msg.Proof)
	if err != nil {
		slog.Error("error validating proof", slog.String("err", err.Error()))

		return
	}

	signature, err := h.getSignature(msg.RequestID, peerID, valid)
	if err != nil {
		slog.Error("error signing validation payload", slog.String("err", err.Error()))

		return
	}

	votingPayload := common.ValidationPayload{
		RequestID:           msg.RequestID,
		ProverID:            peerID,
		IsValid:             valid,
		Signature:           signature,
		ValidationTimestamp: time.Now().UnixNano(),
	}

	votingMsg := common.VotingMessage{
		Type:    common.VoteValidation,
		Payload: votingPayload,
	}

	if err := h.pubsub.Publish(ctx, common.VotingTopic, votingMsg); err != nil {
		slog.Error("error publishing validation message", slog.String("err", err.Error()))

		return
	}
}

func (h *ProofsHandler) getSignature(requestID common.RequestID, peerID peer.ID, isValid bool) ([]byte, error) {
	addr, err := common.PeerIDToEthAddress(peerID)
	if err != nil {
		return nil, errors.Wrap(err, "error converting peer ID to eth address")
	}

	dataToSign := common.DataToSign{
		RequestID:     requestID,
		ProverAddress: addr,
		IsValid:       isValid,
	}

	b, err := json.Marshal(dataToSign)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling a message")
	}

	hash := ethCrypto.Keccak256Hash(b)
	sig, err := ethCrypto.Sign(hash.Bytes(), h.key)
	if err != nil {
		panic(err)
	}

	return sig, nil
}
