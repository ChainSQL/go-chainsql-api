package crypto

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
)

//Account define the account format
type Account struct {
	Address      string `json:"address"`
	PublicKey    string `json:"publicKey"`
	PublicKeyHex string `json:"publicKeyHex"`
	PrivateKey   string `json:"privateKey"`
}

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

func ValidationCreate() (string, error) {
	generated := Account{}
	jsonStr, err := json.Marshal(generated)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}

func GetAccountInfo(address string) (string, error) {
	generated := Account{}
	jsonStr, err := json.Marshal(generated)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}
