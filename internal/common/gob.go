package common

import (
	"bytes"
	"encoding/gob"
)

func InitGobModels() {
	gob.Register(ProverSelectionPayload{})
	gob.Register(ValidationPayload{})
	gob.Register(ProvingRequestMessage{})
	gob.Register(ZKProof{})
	gob.Register(RequestExtension{})
}

func GobEncodeMessage(msg any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func GobDecodeMessage(data []byte, dest any) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(dest)
}
