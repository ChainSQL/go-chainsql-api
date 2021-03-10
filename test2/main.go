package main

import (
	"fmt"

	"github.com/go-chainsql-api/cgofuns"
)

func main() {
	fmt.Println("Hello, World!")
	o := new(cgofuns.CGOFun)
	var signedData []byte
	o.SignPlainData("xh2PnvQnCC9LAJvTH7f8gu82EK7ez", "This is a withdraw for xrp", &signedData)
	fmt.Println("signed data is : ", string(signedData))

	var account, publicKey, publicKeyHex, privateKey []byte
	o.GetValicBLCAddress(&account, &publicKey, &publicKeyHex, &privateKey)
	fmt.Println("\nget valid address \naccount: ", string(account), "\npublicKey: ", string(publicKey), "\npublicKeyHex: ", string(publicKeyHex), "\nprivateKey: ", string(privateKey))
}
