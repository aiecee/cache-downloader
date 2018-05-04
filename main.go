package main

import (
	"fmt"

	"github.com/mcollinge/cache-downloader/network"
)

func main() {
	client, err := network.NewClient("oldschool1.runescape.com:43594")
	if err != nil {
		panic(err)
	}

	handshakePacket := network.HandshakePacket{Version: 169}
	client.Send <- handshakePacket
	handshakeResponse := &network.HandshakeResponse{}
	client.Recieve <- handshakeResponse
	err = <-client.Error
	if err != nil {
		panic(err)
	}
	if handshakeResponse.Response != network.Ok {
		panic(fmt.Errorf("Invalid fesponse from server: %v", handshakeResponse.Response))
	}

	//client.Send <- network.EncryptionPacket{}

	indexRequestPacket := network.ArchiveRequestPacket{
		Priority: true,
		Index:    255,
		File:     7,
	}
	client.Send <- indexRequestPacket
	indexResponse := &network.ArchiveResponse{}
	client.Recieve <- indexResponse
	err = <-client.Error
	if err != nil {
		panic(err)
	}
}
