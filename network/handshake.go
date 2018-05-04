package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

type ResponseType byte

const (
	Unknown                       ResponseType = 101
	Ok                            ResponseType = 0
	LoggedIn                      ResponseType = 2
	InvalidCredentials            ResponseType = 3
	AccountDisabled               ResponseType = 4
	AccountOnline                 ResponseType = 5
	Outdated                      ResponseType = 6
	WorldFull                     ResponseType = 7
	ServerOffline                 ResponseType = 8
	LimitExceeded                 ResponseType = 9
	BadSessionID                  ResponseType = 10
	AccountHijack                 ResponseType = 11
	MembersWorld                  ResponseType = 12
	CouldNotCompleteLogin         ResponseType = 13
	ServerUpdating                ResponseType = 14
	TooManyAttempts               ResponseType = 16
	MembersOnlyArea               ResponseType = 17
	AccountLocked                 ResponseType = 18
	ClosedBeta                    ResponseType = 19
	InvalidLoginServer            ResponseType = 20
	ProfileTransfer               ResponseType = 21
	MalformedPacket               ResponseType = 22
	NoReplyFromLoginServer        ResponseType = 23
	ErrorLoadingProfile           ResponseType = 24
	UnexpectedLoginServerResponse ResponseType = 25
	IPBanned                      ResponseType = 26
	ServiceUnavailable            ResponseType = 27
	NoDisplayName                 ResponseType = 32
	BillingError                  ResponseType = 32
	AccountInaccessible           ResponseType = 37
	VoteToPlay                    ResponseType = 38
	NonEligible                   ResponseType = 55
	NeedAuthenticator             ResponseType = 56
	AuthenticatorCodeWrong        ResponseType = 57
)

type HandshakePacket struct {
	Version int
}

func (h HandshakePacket) Encode() []byte {
	var buffer bytes.Buffer
	buffer.WriteByte(15)
	version := make([]byte, 4)
	binary.BigEndian.PutUint32(version, uint32(h.Version))
	buffer.Write(version)
	return buffer.Bytes()
}

type HandshakeResponse struct {
	Response ResponseType
}

func (h *HandshakeResponse) Decode(reader *bufio.Reader) error {
	h.Response = Unknown
	responseByte, err := reader.ReadByte()
	if err != nil {
		return err
	}
	h.Response = ResponseType(responseByte)
	return nil
}
