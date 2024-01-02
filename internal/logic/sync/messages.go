package sync

type MessageType int

const (
	InitSync MessageType = iota
	SendData
	RequestStorageHash
	SendStorageHash
)

type Message struct {
	Type    MessageType
	Payload any
}
