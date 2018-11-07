package core

import (
	"errors"
)

var (
	errNotAValidatorAddress   = errors.New("This address is not a validator")
	errInvalidMessage         = errors.New("Invalid message")
	errOldMessage             = errors.New("Old message")
	errFutureMessage          = errors.New("Future message")
	errNotFromProposer        = errors.New("Preprepare not from proposer")
	errFailedDecodePreprepare = errors.New("decode preprepare failed")
	errFailedDecodePrepare    = errors.New("decode prepare failed")
	errFailedDecodeCommit     = errors.New("decode commit failed")
	errSubjectsDoNotMatch     = errors.New("subjects do not match")
	errUnauthorized           = errors.New("address not authorized")
)
