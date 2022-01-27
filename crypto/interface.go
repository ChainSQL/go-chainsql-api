package crypto

import (
	"github.com/ChainSQL/go-chainsql-api/common"
	"math/big"
)

type Key interface {
	Private(*uint32) []byte
	Id(*uint32) []byte
	Public(*uint32) []byte
	Type() common.KeyType
	PK(*uint32) (interface{}, error)
	PUB(*uint32) (interface{}, error)
	//Hasher() Hash
}

type Hash interface {
	Version() HashVersion
	Payload() []byte
	PayloadTrimmed() []byte
	Value() *big.Int
	String() string
	Clone() Hash
	MarshalText() ([]byte, error)
}
