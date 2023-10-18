package common

import (
	"crypto/sha256"
	"encoding/binary"
	"github.com/pkg/errors"
	"math/rand"
	"net"
)

func ZKPToRandom(zkProof []byte) (*rand.Rand, error) {
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
