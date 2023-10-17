package common

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
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
