package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

type compressionType byte

const (
	None compressionType = 0
	BZ2  compressionType = 1
	GZ   compressionType = 2

	ChunkSize uint32 = 512
)

type ArchiveRequestPacket struct {
	Priority bool
	Index    byte
	File     int16
}

func (ar ArchiveRequestPacket) Encode() []byte {
	var buffer bytes.Buffer
	if ar.Priority {
		buffer.WriteByte(1)
	} else {
		buffer.WriteByte(0)
	}
	buffer.WriteByte(ar.Index)
	file := make([]byte, 2)
	binary.BigEndian.PutUint16(file, uint16(ar.File))
	buffer.Write(file)
	return buffer.Bytes()
}

type ArchiveResponse struct {
	Index byte
	File  uint16
	Data  []byte
}

func (a *ArchiveResponse) Decode(reader *bufio.Reader) error {
	metaData, err := reader.Peek(8)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	buffer.Write(metaData)

	index, file, compression, compressedSize, err := a.getMetaData(&buffer)
	if err != nil {
		return err
	}

	fmt.Printf("Index: %v File: %v Compression: %v Compressed size: %v", index, file, compression, compressedSize)

	compressionLength := uint32(0)
	if compression != None {
		compressionLength = 4
	}
	fileSize := compressedSize + 5 + compressionLength

	breaks := a.calculateBreaks(fileSize)

	readableBytes := reader.Buffered() // 8 for bytes read (index 1, file 2, compression 1, compressedSize 4)
	if fileSize+3+breaks > uint32(readableBytes) {
		return fmt.Errorf("Index %v Archive %v: Not enough data %v > %v", index, file, (fileSize + 3 + breaks), reader.Buffered())
	}

	discarded, err := reader.Discard(3)
	if err != nil || discarded != 3 {
		return fmt.Errorf("Could not bypass index and file")
	}
	var totalRead uint32 = 3

	var buf bytes.Buffer

	for i := uint32(0); i < breaks+1; {
		byteInBlock := ChunkSize - (totalRead % ChunkSize)
		bytesToRead := math.Min(float64(byteInBlock), float64(fileSize-uint32(buf.Cap())))

		data := make([]byte, int(bytesToRead))
		_, err = reader.Read(data)
		buf.Write(data)

		totalRead += uint32(bytesToRead)
		i++

		if i < breaks {
			endByte, err := reader.ReadByte()
			if err != nil {
				return fmt.Errorf("Could not read chunk end byte: %v", err)
			}
			if endByte != 0xFF {
				return fmt.Errorf("Something went wrong, chunk end does not equal 0xFF")
			}
			totalRead++
		}
	}
	a.Index = index
	a.File = file
	a.Data = buf.Bytes()

	return nil
}

func (a ArchiveResponse) getMetaData(reader *bytes.Buffer) (byte, uint16, compressionType, uint32, error) {
	var index byte
	var file uint16
	var compression compressionType
	var fileSize uint32

	// Read index
	index, err := reader.ReadByte()
	if err != nil {
		return index, file, compression, fileSize, fmt.Errorf("Unable to read file index: %v", err)
	}

	// Read file
	buffer := make([]byte, 2)
	_, err = reader.Read(buffer)
	if err != nil {
		return index, file, compression, fileSize, fmt.Errorf("Unable to read file: %v", err)
	}
	file = binary.BigEndian.Uint16(buffer)

	// Compression type
	compressionByte, err := reader.ReadByte()
	if err != nil {
		return index, file, compression, fileSize, fmt.Errorf("Unable to read compression type: %v", err)
	}
	compression = compressionType(compressionByte)

	// File size
	buffer = make([]byte, 4)
	_, err = reader.Read(buffer)
	if err != nil {
		return index, file, compression, fileSize, fmt.Errorf("Unable to read compressed size: %v", err)
	}
	fileSize = binary.BigEndian.Uint32(buffer)

	return index, file, compression, fileSize, nil
}

func (a ArchiveResponse) calculateBreaks(size uint32) uint32 {
	initialSize := ChunkSize - 3
	if size < initialSize {
		return 0
	}
	left := size - initialSize
	if left%(ChunkSize-1) == 0 {
		return left / (ChunkSize - 1)
	}
	return left/(ChunkSize-1) + 1
}
