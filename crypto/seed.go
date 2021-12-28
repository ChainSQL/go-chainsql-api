package crypto

import (
	"crypto/rand"
	"strings"

	"github.com/buger/jsonparser"
)

type Seed struct {
	seedHash Hash
	version  string
}

type KeyType int

const (
	ECDSA     KeyType = 0
	Ed25519   KeyType = 1
	SoftGMAlg KeyType = 2
)

func GenerateSeed(options string) (*Seed, error) {
	var err error
	sVersion, _ := jsonparser.GetString([]byte(options), "algorithm")
	var version string
	seed := &Seed{}
	switch sVersion {
	case "ed25519":
		version = "ed25519"
		break
	case "secp256k1":
		version = "secp256k1"
		break
	case "softGMAlg":
		version = "softGMAlg"
		break
	default:
		version = "secp256k1"
	}
	seed.version = version
	if version == "softGMAlg" {
		if strings.Contains(options, "secret") {
			secret, _ := (jsonparser.GetString([]byte(options), "secret"))
			seed.seedHash, err = newHashFromString(secret)
		} else {
			seed.seedHash = nil
		}
	} else {
		if strings.Contains(options, "secret") {
			secret, _ := (jsonparser.GetString([]byte(options), "secret"))
			seed.seedHash, err = newHashFromString(secret)
		} else {
			rndBytes := make([]byte, 16)
			if _, err := rand.Read(rndBytes); err != nil {
				return nil, err
			}
			seed.seedHash, err = GenerateFamilySeed(string(rndBytes))
		}
	}
	return seed, err
}

func (s *Seed) GenerateKey(keyType KeyType) (Key, error) {
	var (
		key Key
		err error
	)
	switch keyType {
	case Ed25519:
		key, err = NewEd25519Key(s.seedHash.Payload())
	case ECDSA:
		key, err = NewECDSAKey(s.seedHash.Payload())
	case SoftGMAlg:
		key, err = GenerateKeyPairBySeed(s.seedHash.Payload())
	}
	if err != nil {
		return nil, err
	}
	return key, nil
}
