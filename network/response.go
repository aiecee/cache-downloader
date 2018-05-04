package network

import (
	"bufio"
)

type Response interface {
	Decode(reader *bufio.Reader) error
}
