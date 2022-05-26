package crypto

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"github.com/ChainSQL/go-chainsql-api/common"
)

type ed25519key struct {
	priv ed25519.PrivateKey
	pub  ed25519.PublicKey
}

func checkSequenceIsNil(seq *uint32) {
	if seq != nil {
		fmt.Errorf("Ed25519 keys do not support account families")
	}
}

func (e *ed25519key) Id(seq *uint32) []byte {
	checkSequenceIsNil(seq)
	return Sha256RipeMD160(e.Public(seq))
}

func (e *ed25519key) Public(seq *uint32) []byte {
	checkSequenceIsNil(seq)
	return append([]byte{0xED}, e.priv[32:]...)
}

func (e *ed25519key) PUB(seq *uint32) (interface{}, error) {
	//checkSequenceIsNil(seq)
	//return append([]byte{0xED}, e.priv[32:]...), nil
	return e.pub, nil
}

func (e *ed25519key) Private(seq *uint32) []byte {
	checkSequenceIsNil(seq)
	return e.priv[:]
}

func (e *ed25519key) PK(seq *uint32) (interface{}, error) {
	//checkSequenceIsNil(seq)
	//return e.priv[:], nil
	return e.priv, nil
}

func (k *ed25519key) Type() common.KeyType {
	return common.Ed25519
}

func NewEd25519Key(seed []byte) (*ed25519key, error) {
	r := rand.Reader
	if seed != nil {
		r = bytes.NewReader(Sha512Half(seed))
	}
	pub, priv, err := ed25519.GenerateKey(r)
	if err != nil {
		return nil, err
	}
	key := &ed25519key{priv: priv, pub: pub}
	return key, nil
}
