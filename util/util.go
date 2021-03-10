package util

import(
	"encoding/json"
	"errors"
	"github.com/go-chainsql-api/cgofuns"
	"github.com/go-chainsql-api/common"
)

//GetExtraFee get fee by data for chainsql tx
func GetExtraFee(raw string, dropsPerByte int) int {
	zxcDrops := 1000
	zxcDrops += len(raw) * dropsPerByte
	return zxcDrops
}

// IsChainsqlType justify if a transaction is chainsql type
func IsChainsqlType(t string) bool {
	if t == TableListSet ||
		t == SQLStatement ||
		t == SQLTransaction {
		return true
	}
	return false
}

//SignPlainData sign a plain text and return the signature
func SignPlainData(privateKey string,data string) (string,error){
	o := new(cgofuns.CGOFun)
	var signedData []byte
	ret := o.SignPlainData(privateKey, data, &signedData)
	if !ret{
		return "",errors.New("generate account failed")
	}
	return string(signedData),nil
}

// GenerateAccount generate an account with the format:
// {
//		"address":"zxY4HEbEDSivZwouzwzqHQBA9QbJYdqDTg",
//		"publicKey":"cBPjenRgb2qzoYTnXmPV934kq5wpj2czHoz6rscHtzL34NqZN3KA",
//		"publicKeyHex":"02EA30B2A25844D4AFBAF6020DA9C9FED573AA0058791BFC8642E69888693CF8EA",
//		"privateKey":"xniMQKhxZTMbfWb8scjRPXa5Zv6HB",
// }
func GenerateAccount(args ...string) (string,error) {
	o := new(cgofuns.CGOFun)
	var account, publicKey, publicKeyHex, privateKey []byte
	if len(args) == 1{
		privateKey = []byte(args[0])
	}
	ret := o.GetValicBLCAddress(&account, &publicKey, &publicKeyHex, &privateKey)
	if !ret{
		return "",errors.New("generate account failed")
	}
	generated := common.Account{
		Address : 		string(account),
		PublicKey:		string(publicKey),
		PublicKeyHex:	string(publicKeyHex),
		PrivateKey:		string(privateKey),
	}
	jsonStr, err := json.Marshal(generated)
	if err != nil{
		return "",err
	}
	return string(jsonStr),nil
}