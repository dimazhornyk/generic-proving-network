package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.Handle("/prove", http.HandlerFunc(proveHandler))
	http.Handle("/validate", http.HandlerFunc(validateHandler))

	fmt.Println("starting a server on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}

type ProvingMessage struct {
	RequestID string `json:"request_id"`
	Data      []byte `json:"data"`
}

type ProvingResponse struct {
	Proof []byte `json:"proof"`
}

type ValidationProverMessage struct {
	RequestID string `json:"request_id"`
	Proof     []byte `json:"proof"`
	Data      []byte `json:"data"`
}

type ValidationResponse struct {
	Valid bool `json:"valid"`
}

func proveHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var msg ProvingMessage
	if err := json.Unmarshal(b, &msg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// hashing instead of proving for testing purposes
	hash := sha256.New()
	hash.Write([]byte(msg.RequestID))
	hash.Write(msg.Data)
	sha := hash.Sum(nil)

	resp := ProvingResponse{
		Proof: sha,
	}

	b, err = json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var msg ValidationProverMessage
	if err := json.Unmarshal(b, &msg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := sha256.New()
	hash.Write([]byte(msg.RequestID))
	hash.Write(msg.Data)
	sha := hash.Sum(nil)

	resp := ValidationResponse{
		Valid: bytes.Equal(sha, msg.Proof),
	}

	b, err = json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
