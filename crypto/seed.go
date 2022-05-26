package crypto

import (
	"crypto/rand"
	"github.com/ChainSQL/go-chainsql-api/common"
	"strings"

	"github.com/buger/jsonparser"
)

type Seed struct {
	SeedHash Hash
	version  common.KeyType
}

func GenerateSeed(options string) (*Seed, error) {
	var err error
	sVersion, _ := jsonparser.GetString([]byte(options), "algorithm")
	var version common.KeyType
	seed := &Seed{}
	switch sVersion {
	case "ed25519":
		version = common.Ed25519
		break
	case "secp256k1":
		version = common.ECDSA
		break
	case "softGMAlg":
		version = common.SoftGMAlg
		break
	default:
		version = common.ECDSA
	}
	seed.version = version
	if version == common.SoftGMAlg {
		if strings.Contains(options, "secret") {
			secret, _ := (jsonparser.GetString([]byte(options), "secret"))
			seed.SeedHash, err = newHashFromString(secret)
		} else {
			seed.SeedHash = nil
		}
	} else {
		if strings.Contains(options, "secret") {
			secret, _ := (jsonparser.GetString([]byte(options), "secret"))
			seed.SeedHash, err = newHashFromString(secret)
		} else {
			rndBytes := make([]byte, 16)
			if _, err := rand.Read(rndBytes); err != nil {
				return nil, err
			}
			seed.SeedHash, err = GenerateFamilySeed(string(rndBytes))
		}
	}
	return seed, err
}

func (s *Seed) GenerateKey(keyType common.KeyType) (Key, error) {
	var (
		key Key
		err error
	)
	switch keyType {
	case common.Ed25519:
		key, err = NewEd25519Key(s.SeedHash.Payload())
	case common.ECDSA:
		key, err = NewECDSAKey(s.SeedHash.Payload())
	case common.SoftGMAlg:
		key, err = GenerateKeyPairBySeed(s.SeedHash.Payload())
	}
	if err != nil {
		return nil, err
	}
	return key, nil
}
