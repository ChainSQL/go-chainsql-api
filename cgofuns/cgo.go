package cgofuns

/*
   #cgo CFLAGS: -I.
   #cgo windows LDFLAGS: -L ./clib/win/ -lsignature -lboost_regex -lcrypto -lssl -lgdi32 -lstdc++
   #cgo linux LDFLAGS: -L ./clib/linux/ -lsignature -lboost_regex -lcrypto -lssl -ldl -lstdc++
   #cgo darwin LDFLAGS: -L ./clib/darwin/ -lsignature -lboost_regex -lcrypto -lssl -lstdc++  
   #include "key_manager_api_c.h"
   #include "stdlib.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func cChar2byte(cPtr *C.uchar, cLen uint64, goSlice *[]byte) {
	p := uintptr(unsafe.Pointer(cPtr))
	for i := uint64(0); i < uint64(cLen); i++ {
		j := *(*byte)(unsafe.Pointer(p))
		*goSlice = append(*goSlice, j)
		p += unsafe.Sizeof(j)
	}
	// fmt.Println("go slice is : ", striIng(*goSlice))
}

//CGOFun is used to call the functions in c libry.
type CGOFun struct {
}

//SignTransaction : sign for common tranction, sTx should be json string
func (o *CGOFun) SignTransaction(sPrivateKey string, sTx string, pSignedData *[]byte, pHash *[]byte) bool {
	pKey := []byte(sPrivateKey)
	pTx := []byte(sTx)

	iHashLen := (C.ulong)(0)
	var pCData, pCHash *C.uchar
	iLenData := C.sign_for_common_transaction_ex((*C.uchar)(unsafe.Pointer(&pKey[0])), (C.ulong)(len(pKey)),
		(*C.uchar)(unsafe.Pointer(&pTx[0])), (C.ulong)(len(pTx)),
		&pCData,
		&pCHash, &iHashLen)

	if iLenData <= 0 {
		return false
	}
	defer C.free(unsafe.Pointer(pCData))
	defer C.free(unsafe.Pointer(pCHash))

	// fmt.Println("sign_for_common_transaction_ex ok, sigedData's length: ", iLenData)

	cChar2byte(pCData, uint64(iLenData), pSignedData)
	cChar2byte(pCHash, uint64(iHashLen), pHash)

	return true
}

//SignPlainData : sign for common data
func (o *CGOFun) SignPlainData(sPrivateKey string, sPlain string, pSignedData *[]byte) bool {
	pKey := []byte(sPrivateKey)
	pPlain := []byte(sPlain)

	var pCData *C.uchar
	iLenData := C.sign_common_data((*C.uchar)(unsafe.Pointer(&pKey[0])), (C.ulong)(len(pKey)),
		(*C.uchar)(unsafe.Pointer(&pPlain[0])), (C.ulong)(len(pPlain)),
		&pCData)
	if iLenData <= 0 {
		return false
	}
	defer C.free(unsafe.Pointer(pCData))

	// fmt.Println("sign_common_data ok, sigedData's length: ", iLenData)

	cChar2byte(pCData, uint64(iLenData), pSignedData)

	//fmt.Println(string(*pSignedData))

	return true
}

//GetValicBLCAddress : generate valid address which is used in blockchain
func (o *CGOFun) GetValicBLCAddress(pAccount *[]byte, pPublicKey *[]byte, pPublicKeyHex *[]byte, pPrivateKey *[]byte) bool {
	var iAccountLen, iPublicKey, iPublicKeyHex, iPrivateKey C.ulong
	var pCAccount, pCPubKey, pCPubKeyHex, pPriKey *C.uchar

	pSeed := (*C.uchar)(C.NULL)
	if (len(*pPrivateKey)) > 0 {
		pSeed = (*C.uchar)(unsafe.Pointer(&(*pPrivateKey)[0]))
	}
	bRet := C.get_valid_address(pSeed, (C.ulong)(len(*pPrivateKey)),
		&pCAccount, &iAccountLen,
		&pCPubKey, &iPublicKey,
		&pCPubKeyHex, &iPublicKeyHex,
		&pPriKey, &iPrivateKey)
	if !bRet {
		fmt.Println("fail to get address from c library")
		return false
	}
	defer C.free(unsafe.Pointer(pCAccount))
	defer C.free(unsafe.Pointer(pCPubKey))
	defer C.free(unsafe.Pointer(pCPubKeyHex))
	defer C.free(unsafe.Pointer(pPriKey))

	// fmt.Println("GetValicBLCAddress ok")

	cChar2byte(pCAccount, uint64(iAccountLen), pAccount)
	cChar2byte(pCPubKey, uint64(iPublicKey), pPublicKey)
	cChar2byte(pCPubKeyHex, uint64(iPublicKeyHex), pPublicKeyHex)
	if len(*pPrivateKey) <= 0 {
		// fmt.Println("test come here")
		cChar2byte(pPriKey, uint64(iPrivateKey), pPrivateKey)
	}

	return true
}
