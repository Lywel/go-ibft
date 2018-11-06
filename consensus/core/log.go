package core

import (
	"fmt"

	"bitbucket.org/ventureslash/go-ibft/consensus"
)

// Logger helps printing clears logs
type Logger struct {
	address consensus.Address
}

// Log prints the message with the address at the beginning
func (l *Logger) Log(args ...interface{}) {
	fmt.Printf("%s: ", l.address.String())
	fmt.Println(args...)
}
