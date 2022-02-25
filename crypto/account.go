package crypto

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/ChainSQL/go-chainsql-api/common"
	"log"
	"strings"
)

//Account define the account format
type Account struct {
	Address         string      `json:"address"`
	PublicKeyBase58 string      `json:"publicKeyBase58"`
	PublicKeyHex    string      `json:"publicKeyHex"`
	PrivateSeed     string      `json:"privateSeed"`
	PrivateKey      interface{} `json:"privateKey"`
	PublicKey       interface{} `json:"publicKey"`
}

type SeedKey struct {
	Seed      string `json:"seed"`
	PublicKey string `json:"publicKey"`
}

//生成特殊地址仍然使用此方法
func GenerateAccount(args ...string) (string, error) {
	var seed Hash
	var err error
	var key *ecdsaKey
	if len(args) == 1 {
		seed, err = NewRippleHash(args[0])
		if err != nil {
			return "", err
		}
	} else {
		rndBytes := make([]byte, 16)
		if _, err := rand.Read(rndBytes); err != nil {
			return "", err
		}
		seed, err = GenerateFamilySeed(string(rndBytes))
		if err != nil {
			return "", err
		}
	}
	key, err = NewECDSAKey(seed.Payload())
	if err != nil {
		log.Println(err)
		return "", err
	}

	sequenceZero := uint32(0)
	account, _ := AccountId(key, &sequenceZero)
	publicKey, _ := AccountPublicKey(key, &sequenceZero)
	pk, err := key.PK(&sequenceZero)
	if err != nil {
		return "", err
	}

	pub, err := key.PUB(&sequenceZero)
	if err != nil {
		return "", err
	}
	generated := Account{
		Address:         account.String(),
		PublicKeyBase58: publicKey.String(),
		PublicKeyHex:    fmt.Sprintf("%X", key.Public(&sequenceZero)),
		PrivateSeed:     seed.String(),
		PrivateKey:      pk,
		PublicKey:       pub,
	}
	jsonStr, err := json.Marshal(generated)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}

func GenerateAddress(options string) (string, error) {
	var seed *Seed
	var err error
	var key Key
	if strings.Contains(options, "secret") && !strings.Contains(options, "algorithm") {
		return "", fmt.Errorf("Invalid parameter")
	}
	seed, err = GenerateSeed(options)
	if err != nil {
		return "", err
	}
	sVersion := seed.version
	switch sVersion {
	case common.Ed25519:
		key, err = NewEd25519Key(seed.SeedHash.Payload())
		break
	case common.SoftGMAlg:
		key, err = GenerateKeyPair(seed)
		break
	case common.ECDSA:
		key, err = NewECDSAKey(seed.SeedHash.Payload())
		break
	default:
		key, err = NewECDSAKey(seed.SeedHash.Payload())
	}
	if err != nil {
		return "", err
	}

	sequenceZero := uint32(0)
	account, err := AccountId(key, &sequenceZero)
	if err != nil {
		return "", err
	}
	publicKey, err := AccountPublicKey(key, &sequenceZero)
	if err != nil {
		return "", err
	}
	var privSeed Hash
	if sVersion == common.SoftGMAlg {
		privSeed, err = AccountPrivateKey(key, &sequenceZero)
		if err != nil {
			return "", err
		}
	} else {
		privSeed = seed.SeedHash
	}
	pk, err := key.PK(&sequenceZero)
	if err != nil {
		return "", err
	}

	pub, err := key.PUB(&sequenceZero)
	if err != nil {
		return "", err
	}

	generated := Account{
		Address:         account.String(),
		PublicKeyBase58: publicKey.String(),
		PublicKeyHex:    fmt.Sprintf("%X", key.Public(&sequenceZero)),
		PrivateSeed:     privSeed.String(),
		PrivateKey:      pk,
		PublicKey:       pub,
	}

	jsonStr, err := json.Marshal(generated)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}

func GenerateAddressObj(options string) (*Account, error) {
	var seed *Seed
	var err error
	var key Key
	if strings.Contains(options, "secret") && !strings.Contains(options, "algorithm") {
		return nil, fmt.Errorf("Invalid parameter")
	}
	seed, err = GenerateSeed(options)
	if err != nil {
		return nil, err
	}
	sVersion := seed.version
	sequenceZero := uint32(0)
	switch sVersion {
	case common.Ed25519:
		key, err = NewEd25519Key(seed.SeedHash.Payload())
		break
	case common.SoftGMAlg:
		key, err = GenerateKeyPair(seed)
		break
	case common.ECDSA:
		key1 := &ecdsaKey{}
		key1, err = NewECDSAKey(seed.SeedHash.Payload())
		key = key1.GenerateEcdsaKey(sequenceZero)
		break
	default:
		key, err = NewECDSAKey(seed.SeedHash.Payload())
	}

	if err != nil {
		return nil, err
	}

	account, err := AccountId(key, nil)
	if err != nil {
		return nil, err
	}
	publicKey, err := AccountPublicKey(key, nil)
	if err != nil {
		return nil, err
	}
	var privSeed Hash
	if sVersion == common.SoftGMAlg {
		privSeed, err = AccountPrivateKey(key, nil)
		if err != nil {
			return nil, err
		}
	} else {
		privSeed = seed.SeedHash
	}
	pk, err := key.PK(nil)
	if err != nil {
		return nil, err
	}

	pub, err := key.PUB(nil)
	if err != nil {
		return nil, err
	}

	return &Account{
		Address:         account.String(),
		PublicKeyBase58: publicKey.String(),
		PublicKeyHex:    fmt.Sprintf("%X", key.Public(nil)),
		PrivateSeed:     privSeed.String(),
		PrivateKey:      pk,
		PublicKey:       pub,
	}, nil
}

func ValidationCreate() (string, error) {
	var seed Hash
	var err error
	var key *ecdsaKey
	rndBytes := make([]byte, 16)
	if _, err := rand.Read(rndBytes); err != nil {
		return "", err
	}
	seed, err = GenerateFamilySeed(string(rndBytes))
	if err != nil {
		return "", err
	}
	key, err = NewECDSAKey(seed.Payload())
	if err != nil {
		log.Println(err)
		return "", err
	}
	publicKey, _ := NodePublicKey(key)
	generated := SeedKey{
		Seed:      seed.String(),
		PublicKey: publicKey.String(),
	}
	jsonStr, err := json.Marshal(generated)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}
