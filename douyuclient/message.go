package douyuclient

import (
	"encoding/binary"
)

// Message contains the timestamp and  content of the received message.
type Message struct {
	NanoTimestamp int64
	RawContent    []byte
}

func (message *Message) Head() []byte {
	return message.RawContent[:12]
}

func (message *Message) Body() []byte {
	return message.RawContent[12:]
}

func (message *Message) Text() string {
	body := message.Body()
	if body[len(body)-1] == 0 {
		return string(body[:len(body)-1])
	}
	return string(body)
}

func CreateRawSendMsg(data string) []byte {
	return CreateRawMsg([]byte(data), 689)
}

func CreateRawMsg(data []byte, msgType uint16) []byte {
	if data[len(data)-1] != 0 {
		data = append(data, 0)
	}
	dataLen := len(data)

	res := make([]byte, 12+dataLen)
	binary.LittleEndian.PutUint32(res, uint32(8+dataLen))
	binary.LittleEndian.PutUint32(res[4:], uint32(8+dataLen))
	binary.LittleEndian.PutUint16(res[8:], msgType)
	copy(res[12:], data)
	return res
}
