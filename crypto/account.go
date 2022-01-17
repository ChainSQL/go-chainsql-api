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
	Address      string `json:"address"`
	PublicKey    string `json:"publicKey"`
	PublicKeyHex string `json:"publicKeyHex"`
	PrivateKey   string `json:"privateKey"`
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
	generated := Account{
		Address:      account.String(),
		PublicKey:    publicKey.String(),
		PublicKeyHex: fmt.Sprintf("%X", key.Public(&sequenceZero)),
		PrivateKey:   seed.String(),
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
		key, err = NewEd25519Key(seed.seedHash.Payload())
		break
	case common.SoftGMAlg:
		key, err = GenerateKeyPair(seed)
		break
	case common.ECDSA:
		key, err = NewECDSAKey(seed.seedHash.Payload())
		break
	default:
		key, err = NewECDSAKey(seed.seedHash.Payload())
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
	var privateKey Hash
	if sVersion == common.SoftGMAlg {
		privateKey, err = AccountPrivateKey(key, &sequenceZero)
		if err != nil {
			return "", err
		}
	} else {
		privateKey = seed.seedHash
	}
	generated := Account{
		Address:      account.String(),
		PublicKey:    publicKey.String(),
		PublicKeyHex: fmt.Sprintf("%X", key.Public(&sequenceZero)),
		PrivateKey:   privateKey.String(),
	}
	jsonStr, err := json.Marshal(generated)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
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
