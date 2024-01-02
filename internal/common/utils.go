package common

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"math/rand"
	"net"
)

// BytesToRandom converts a byte slice to a random number generator
func BytesToRandom(zkProof []byte) (*rand.Rand, error) {
	hash := sha256.New()
	hash.Write(zkProof)
	sha := hash.Sum(nil)

	n := binary.BigEndian.Uint64(sha)
	src := rand.NewSource(int64(n))

	return rand.New(src), nil
}

func Map[T, E any](slice []T, f func(T) E) []E {
	res := make([]E, len(slice))

	for i, v := range slice {
		res[i] = f(v)
	}

	return res
}

func Filter[T any](slice []T, f func(T) bool) []T {
	res := make([]T, 0)

	for _, v := range slice {
		if f(v) {
			res = append(res, v)
		}
	}

	return res
}

func AvailablePort() (string, error) {
	server, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", errors.Wrap(err, "error listening to a port")
	}
	defer server.Close()

	hostString := server.Addr().String()

	_, port, err := net.SplitHostPort(hostString)
	if err != nil {
		return "", errors.Wrap(err, "error splitting host and port")
	}

	return port, nil
}

func PeerIDToEthAddress(peerID peer.ID) (string, error) {
	pubkey, err := peerID.ExtractPublicKey()
	if err != nil {
		return "", errors.Wrap(err, "error marshalling ID")
	}

	b, err := pubkey.Raw()
	if err != nil {
		return "", errors.Wrap(err, "error getting raw public key")
	}

	x, y := secp256k1.DecompressPubkey(b)
	keccak := ethCrypto.Keccak256(append(x.Bytes(), y.Bytes()...))

	return "0x" + hex.EncodeToString(keccak[len(keccak)-20:]), nil
}

func GetRSV(signature []byte) ([32]byte, [32]byte, uint8, error) {
	if len(signature) != 65 {
		return [32]byte{}, [32]byte{}, 0, errors.New("wrong signature length")
	}

	return [32]byte(signature[:32]), [32]byte(signature[32:64]), 27 + signature[64], nil
}
