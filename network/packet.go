package network

import (
	"bytes"
	"encoding/binary"
)

type Packet interface {
	Encode() []byte
}

type EncryptionPacket struct {
	Key byte
}

func (e EncryptionPacket) Encode() []byte {
	var buffer bytes.Buffer
	buffer.WriteByte(4)
	buffer.WriteByte(e.Key)
	short := make([]byte, 2)
	binary.BigEndian.PutUint16(short, 0)
	buffer.Write(short)
	return buffer.Bytes()
}
