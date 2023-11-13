package sync

type MessageType int

const (
	InitSync MessageType = iota
	SendData
	RequestChecksum
	SendChecksum
)

type Message struct {
	Type    MessageType
	Payload any
}
