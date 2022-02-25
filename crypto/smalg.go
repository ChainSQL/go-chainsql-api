package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ChainSQL/go-chainsql-api/common"
	"strings"

	"github.com/peersafe/gm-crypto/sm2"
	"github.com/peersafe/gm-crypto/sm3"
)

type smKey struct {
	PrivateKey string
	PublicKey  string
}

func generateKeyPair() (*smKey, error) {
	priv, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	priKey := leftPad(priv.D.Text(16), 64)
	Px := leftPad(priv.PublicKey.X.Text(16), 64)
	Py := leftPad(priv.PublicKey.Y.Text(16), 64)
	pubKey := "47" + Px + Py
	private := &smKey{}
	private.PrivateKey = strings.ToUpper(priKey)
	private.PublicKey = strings.ToUpper(pubKey)
	return private, nil
}

func GenerateKeyPairBySeed(seed []byte) (*smKey, error) {
	var err error
	seedStr := hex.EncodeToString(seed)
	priv, err := sm2.GenerateKeyBySeed(rand.Reader, seedStr)
	if err != nil {
		return nil, err
	}
	priKey := leftPad(priv.D.Text(16), 64)
	Px := leftPad(priv.PublicKey.X.Text(16), 64)
	Py := leftPad(priv.PublicKey.Y.Text(16), 64)
	pubKey := "47" + Px + Py
	private := &smKey{}
	private.PrivateKey = strings.ToUpper(priKey)
	private.PublicKey = strings.ToUpper(pubKey)
	return private, nil
}

func PrivKeyFromBytes(privKey []byte) (*sm2.PrivateKey, error) {
	seedStr := hex.EncodeToString(privKey)
	return sm2.GenerateKeyBySeed(rand.Reader, seedStr)
}

/*func sm2KeyPairToChainsqlKeyPair(priv *PrivateKey)(*smKey, error){
	priKey := leftPad(priv.D.Text(16), 64)
	Px := leftPad(priv.PublicKey.X.Text(16), 64)
	Py := leftPad(priv.PublicKey.Y.Text(16), 64)
	pubKey := "47"+ Px + Py
	private := &smKey{}
	private.PrivateKey = strings.ToUpper(priKey)
	private.PublicKey = strings.ToUpper(pubKey)
	return private,nil
}*/

func GenerateKeyPair(seed *Seed) (*smKey, error) {
	if seed.SeedHash == nil {
		smKeyPair, err := generateKeyPair()
		if err != nil {
			return nil, err
		}
		return smKeyPair, nil
	} else {
		// 在国密算法内部添加新的方法
		smKeyPair, err := GenerateKeyPairBySeed(seed.SeedHash.Payload())
		if err != nil {
			return nil, err
		}
		return smKeyPair, nil
	}
}

/**
 * 补全16进制字符串
 */
func leftPad(input string, num int) string {
	if len([]byte(input)) >= num {
		return input
	}
	length := num - len([]byte(input))
	for i := 0; i < length; i++ {
		input = "0" + input
	}
	return input
}

// GM算法没有seed生成公私钥，因此不存在sequence
func (k *smKey) Id(sequence *uint32) []byte {
	return Sha256RipeMD160(k.Public(sequence))
}

func (k *smKey) Private(sequence *uint32) []byte {
	privateByte, err := hex.DecodeString(k.PrivateKey) // 转码
	if err != nil {
		fmt.Printf("PrivateKey transcoding exception ")
	}
	return privateByte
}

func (k *smKey) PK(sequence *uint32) (interface{}, error) {
	privateByte, err := hex.DecodeString(k.PrivateKey) // 转码
	if err != nil {
		fmt.Printf("PrivateKey transcoding exception ")
		return nil, fmt.Errorf("PrivateKey transcoding exception ")
	}
	pk, err := PrivKeyFromBytes(privateByte)
	if err != nil {
		fmt.Printf("sm2 PrivKeyFromBytes err : %v\n", err)
		return nil, err
	}
	return pk, nil
}

func (k *smKey) Public(sequence *uint32) []byte {
	pubkeyByte, err := hex.DecodeString(k.PublicKey)
	if err != nil {
		fmt.Printf("PublicKey transcoding exception ")
	}
	return []byte(pubkeyByte)
}

func (k *smKey) PUB(sequence *uint32) (interface{}, error) {
	privateByte, err := hex.DecodeString(k.PrivateKey) // 转码
	if err != nil {
		fmt.Printf("PrivateKey transcoding exception ")
		return nil, fmt.Errorf("PrivateKey transcoding exception ")
	}
	pk, err := PrivKeyFromBytes(privateByte)
	if err != nil {
		fmt.Printf("sm2 PrivKeyFromBytes err : %v\n", err)
		return nil, err
	}
	return pk.PublicKey, nil
}

func (k *smKey) Type() common.KeyType {
	return common.SoftGMAlg
}

func SmSign(k *sm2.PrivateKey, msg []byte) (string, error) {
	hashed := sm3.SumSM3(msg)

	r, s, err := sm2.SignWithDigest(rand.Reader, k, hashed)
	if err != nil {
		return "", err
	}
	/*if !sm2.VerifyWithDigest(&k.PublicKey, hashed, r, s) {
		log.Println("err")
	}*/
	sigValueHex := leftPad(r.Text(16), 64) + leftPad(s.Text(16), 64)
	return strings.ToUpper(sigValueHex), nil
}
