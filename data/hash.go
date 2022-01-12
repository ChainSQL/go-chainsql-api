package data

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/ChainSQL/go-chainsql-api/common"
	"github.com/ChainSQL/go-chainsql-api/crypto"
)

const (
	PUBKEY_LENGTH_GM     int = 65
	PUBKEYL_ENGTH_COMMON int = 33
)



type Hash128 [16]byte
type Hash160 [20]byte
type Hash256 [32]byte
type Vector256 []Hash256
type VariableLength []byte

//type PublicKey [33]byte
type PublicKey struct {
	KeyStore [65]byte
	KeyValue []byte
	KeyType  common.KeyType
}

type Account [20]byte
type RegularKey [20]byte

//type Seed [16]byte

var zero256 Hash256
var zeroAccount Account

//var zeroPublicKey PublicKey
//var zeroSeed Seed

func (h *Hash128) Bytes() []byte {
	if h == nil {
		return nil
	}
	return h[:]
}

func (h Hash128) String() string {
	return string(b2h(h[:]))
}

func (h *Hash160) Bytes() []byte {
	if h == nil {
		return nil
	}
	return h[:]
}

func (h Hash160) String() string {
	return string(b2h(h[:]))
}

func (h *Hash160) Account() *Account {
	if h == nil {
		return nil
	}
	var a Account
	copy(a[:], h[:])
	return &a
}

func (h *Hash160) Currency() *Currency {
	if h == nil {
		return nil
	}
	var c Currency
	copy(c[:], h[:])
	return &c
}

// Accepts either a hex string or a byte slice of length 32
func NewHash256(value interface{}) (*Hash256, error) {
	var h Hash256
	switch v := value.(type) {
	case []byte:
		if len(v) != 32 {
			return nil, fmt.Errorf("NewHash256: Wrong length %X", value)
		}
		copy(h[:], v)
	case string:
		n, err := hex.Decode(h[:], []byte(v))
		if err != nil {
			return nil, err
		}
		if n != 32 {
			return nil, fmt.Errorf("NewHash256: Wrong length %s", v)
		}
	default:
		return nil, fmt.Errorf("NewHash256: Wrong type %+v", v)
	}
	return &h, nil
}

func (h Hash256) IsZero() bool {
	return h == zero256
}

func (h Hash256) Xor(x Hash256) Hash256 {
	var xor Hash256
	for i := range h {
		xor[i] = h[i] ^ x[i]
	}
	return x
}

func (h Hash256) Compare(x Hash256) int {
	return bytes.Compare(h[:], x[:])
}

func (h *Hash256) Bytes() []byte {
	if h == nil {
		return nil
	}
	return h[:]
}

func (h Hash256) String() string {
	return string(b2h(h[:]))
}

func (h Hash256) TruncatedString(length int) string {
	return string(b2h(h[:length]))
}

func (v Vector256) String() string {
	var s []string
	for _, h := range v {
		s = append(s, h.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(s, ","))
}

func (v *VariableLength) String() string {
	if v != nil {
		b, _ := v.MarshalText()
		return string(b)
	}
	return ""
}

func (v *VariableLength) Bytes() []byte {
	if v != nil {
		return []byte(*v)
	}
	return []byte(nil)
}

func (p *PublicKey) SetKey(kType common.KeyType) {
	p.KeyType = kType
	if common.SoftGMAlg == kType {
		p.KeyValue = p.KeyStore[:PUBKEY_LENGTH_GM]
	} else {
		p.KeyValue = p.KeyStore[:PUBKEYL_ENGTH_COMMON]
	}

}
func (p PublicKey) NodePublicKey() string {
	hash, err := crypto.NewNodePublicKey(p.KeyValue[:])
	if err != nil {
		return "Bad node public key"
	}
	return hash.String()
}

func (p PublicKey) String() string {
	b, _ := p.MarshalText()
	return string(b)
}

func (p PublicKey) IsZero() bool {
	if p.KeyValue == nil {
		return true
	}
	return false
}

func (p *PublicKey) Bytes() []byte {
	if p != nil {
		return p.KeyValue[:]
	}
	return []byte(nil)
}

// Expects address in base58 form
func NewAccountFromAddress(s string) (*Account, error) {
	hash, err := crypto.NewRippleHashCheck(s, crypto.RIPPLE_ACCOUNT_ID)
	if err != nil {
		return nil, err
	}
	var account Account
	copy(account[:], hash.Payload())
	return &account, nil
}

func (a Account) Hash() (crypto.Hash, error) {
	return crypto.NewAccountId(a[:])
}

func (a Account) String() string {
	address, err := a.Hash()
	if err != nil {
		return fmt.Sprintf("Bad Address: %s", b2h(a[:]))
	}
	return address.String()
}

func (a Account) IsZero() bool {
	return a == zeroAccount
}

func (a *Account) Bytes() []byte {
	if a != nil {
		return a[:]
	}
	return []byte(nil)
}

func (a Account) Compare(b Account) int {
	return bytes.Compare(a[:], b[:])
}

func (a Account) Less(b Account) bool {
	return a.Compare(b) < 0
}

func (a Account) Equals(b Account) bool {
	return a == b
}

func (a Account) Hash256() Hash256 {
	var h Hash256
	copy(h[:], a[:])
	return h
}

// Expects address in base58 form
func NewRegularKeyFromAddress(s string) (*RegularKey, error) {
	hash, err := crypto.NewRippleHashCheck(s, crypto.RIPPLE_ACCOUNT_ID)
	if err != nil {
		return nil, err
	}
	var regKey RegularKey
	copy(regKey[:], hash.Payload())
	return &regKey, nil
}

func (r RegularKey) Hash() (crypto.Hash, error) {
	return crypto.NewAccountId(r[:])
}

func (r RegularKey) String() string {
	address, err := r.Hash()
	if err != nil {
		return fmt.Sprintf("Bad Address: %s", b2h(r[:]))
	}
	return address.String()
}

func (r *RegularKey) Bytes() []byte {
	if r != nil {
		return r[:]
	}
	return []byte(nil)
}

// Expects address in base58 form
// func NewSeedFromAddress(s string, version crypto.HashVersion) (*Seed, error) {
// 	keySeed := &Seed{}
// 	hash, err := crypto.NewRippleHashCheck(s, version)
// 	if err != nil {
// 		return nil, err
// 	}
// 	keySeed.seedHash = hash
// 	return keySeed, nil
// }

// func (s *Seed) Hash() (crypto.Hash, error) {
// 	return crypto.NewFamilySeed(s[:])
// }

// func (s Seed) String() string {
// 	address, err := s.Hash()
// 	if err != nil {
// 		return fmt.Sprintf("Bad Address: %s", b2h(s[:]))
// 	}
// 	return address.String()
// }

// func (s *Seed) Bytes() []byte {
// 	if s != nil {
// 		return s[:]
// 	}
// 	return []byte(nil)
// }

// func (s *Seed) Key(keyType KeyType) crypto.Key {
// 	var (
// 		key crypto.Key
// 		err error
// 	)
// 	switch keyType {
// 	case Ed25519:
// 		key, err = crypto.NewEd25519Key(s[:])
// 	case ECDSA:
// 		key, err = crypto.NewECDSAKey(s[:])
// 	case SoftGMAlg:
// 		key, err = crypto.GenerateKeyPairBySeed(s[:])
// 	}
// 	if err != nil {
// 		panic(fmt.Sprintf("bad seed: %v", err))
// 	}
// 	return key
// }

// func (s *Seed) AccountId(keyType KeyType, sequence *uint32) Account {
// 	var account Account
// 	copy(account[:], s.Key(keyType).Id(sequence))
// 	return account
// }

func KeyFromSecret(secret string) (crypto.Key, error) {
	var version crypto.HashVersion
	var err error
	var seed *crypto.Seed
	var regSoftGMSeed = "^[a-zA-Z1-9]{51,51}"
	r := regexp.MustCompile(regSoftGMSeed)
	if r.MatchString(secret) {
		version = crypto.RIPPLE_ACCOUNT_PRIVATE
		seed, err = crypto.NewRippleSeed(secret, version)
		if err != nil {
			return nil, err
		}
		return seed.GenerateKey(common.SoftGMAlg)
	} else {
		version = crypto.RIPPLE_FAMILY_SEED
		seed, err := crypto.NewRippleSeed(secret, version)
		if err != nil {
			return nil, err
		}
		return seed.GenerateKey(common.ECDSA)
	}
}
