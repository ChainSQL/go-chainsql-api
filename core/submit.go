package core

import (
	"fmt"
	"log"
	"strings"
	"sync"

	. "github.com/ChainSQL/go-chainsql-api/common"
	"github.com/ChainSQL/go-chainsql-api/crypto"
	. "github.com/ChainSQL/go-chainsql-api/data"
	"github.com/ChainSQL/go-chainsql-api/export"
	"github.com/ChainSQL/go-chainsql-api/net"
	"github.com/ChainSQL/go-chainsql-api/util"
	"github.com/buger/jsonparser"
)

// TxJSON is a type indicate a transaction json
type TxJSON interface {
}

// TxSigned is signed transaction
type TxSigned struct {
	blob string
	hash string
}

// IPrepare is an interface that a struct call submit() method must implment
type IPrepare interface {
	PrepareTx() (Signer, error)
}

// SubmitBase base struct
type SubmitBase struct {
	expect   string
	callback export.Callback
	client   *net.Client
	IPrepare
}

// // Submit submit a tx with a cocurrent expect
// func (s *SubmitBase) Submit(cond string) string {
// 	s.expect = cond
// 	ret := s.doSubmit()
// 	jsonRet, _ := json.Marshal(ret)
// 	return string(jsonRet)
// }
// Submit submit a tx with a cocurrent expect
func (s *SubmitBase) Submit(cond string) (txRet *TxResult) {
	if cond == "" {
		s.expect = util.SendSuccess
	} else {
		s.expect = cond
	}
	return s.doSubmit()
}

//SubmitAsync submit a transaction and got response asynchronously
func (s *SubmitBase) SubmitAsync(callback export.Callback) {
	s.callback = callback
	s.doSubmit()
}

func (s *SubmitBase) doSubmit() *TxResult {
	tx, err := s.PrepareTx()
	if err != nil {
		log.Printf("doSubmit error:%s\n", err)
		return &TxResult{
			ErrorCode:    "errPrepareTx",
			ErrorMessage: err.Error(),
		}
	}
	// str, err := json.Marshal(tx)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(string(str))

	key, err := KeyFromSecret(s.client.Auth.Secret)
	if err != nil {
		log.Printf("doSubmit error:%s\n", err)
		return &TxResult{
			ErrorCode:    "errGenerateKey",
			ErrorMessage: err.Error(),
		}
	}
	/*var hasher hash.Hash
	if key.Type() == common.SoftGMAlg {
		hasher = sm3.New()
	}else {
		hasher = sha512.New()
	}
	*/
	sequenceZero := uint32(0)
	err = Sign(tx, key, &sequenceZero, key.Type())
	if err != nil {
		log.Printf("doSubmit error:%s\n", err)
		return &TxResult{
			ErrorCode:    "errSign",
			ErrorMessage: err.Error(),
		}
	}

	txSigned := &TxSigned{
		blob: fmt.Sprintf("%X", *tx.GetBlob()),
		hash: string(crypto.B2H32(*tx.GetHash())),
	}

	// log.Printf("hash:%s\n", txSigned.hash)
	// log.Printf("blob:%s\n", txSigned.blob)

	return s.handleSignedTx(txSigned)
}

func (s *SubmitBase) checkWaitGroupDone(wait *sync.WaitGroup, countDone *int, maxDone int) bool {
	if *countDone < maxDone {
		defer wait.Done()
		(*countDone)++
		return true
	} else {
		return false
	}
}

// handleSignedTx handles signed transaction submit
// Chainsql will re-use this function when commit called
func (s *SubmitBase) handleSignedTx(tx *TxSigned) *TxResult {
	ret := &TxResult{}
	wait := new(sync.WaitGroup)
	countDone := 0
	maxDone := 0
	if s.expect != util.SendSuccess {
		// subscribe for result
		s.client.SubscribeTx(tx.hash, func(msg string) {
			// log.Println(msg)
			if !s.checkWaitGroupDone(wait, &countDone, maxDone) {
				// log.Printf("Already %d times of wait.Done,msg=%s \n", countDone, msg)
				return
			}
			status, err := jsonparser.GetString([]byte(msg), "status")
			if err != nil {
				log.Printf("handleSignedTx error:%s\n", err)
				return
			}

			ret.Status = status
			ret.TxHash = tx.hash
			doneTwice := false
			//db_success come before validate_success
			// validate_success will not come any more
			if status == s.expect && status == util.DbSuccess && countDone == 1 {
				doneTwice = true
			}
			if status != s.expect {
				if status != util.ValidateSuccess && status != util.DbSuccess {
					if status == util.ValidateError {
						errCode, _ := jsonparser.GetString([]byte(msg), "error")
						errorMessage, _ := jsonparser.GetString([]byte(msg), "error_message")
						ret.ErrorCode = errCode
						ret.ErrorMessage = errorMessage
					}
					if s.expect == util.DbSuccess &&
						(strings.Contains(status, "validate_") ||
							(strings.Contains(status, "db_") && countDone == 1)) {
						doneTwice = true
					}
				}
			}
			if doneTwice {
				if !s.checkWaitGroupDone(wait, &countDone, maxDone) {
					// log.Printf("Already %d times of wait.Done,msg=%s \n", countDone, msg)
					return
				}
			}
		})
	}

	//submit transaction
	response := s.client.Submit(tx.blob)
	status, err := jsonparser.GetString([]byte(response), "status")
	if err != nil {
		log.Printf("handleSignedTx error:%s\n", err)
	}
	if status == "error" {
		log.Printf("Send tx error:%s\n", response)
		errorCode, _ := jsonparser.GetString([]byte(response), "error")
		errorMessage, _ := jsonparser.GetString([]byte(response), "error_message")
		return &TxResult{
			Status:       util.SendError,
			TxHash:       tx.hash,
			ErrorCode:    errorCode,
			ErrorMessage: errorMessage,
		}
	}
	result, err := jsonparser.GetString([]byte(response), "result", "engine_result")
	if err != nil {
		log.Printf("handleSignedTx error:%s\n", err)
	}
	if result == "tesSUCCESS" {
		if s.expect == util.SendSuccess {
			return &TxResult{
				Status: util.SendSuccess,
				TxHash: tx.hash,
			}
		} else {
			//waiting for subscribe result
			if s.expect == util.ValidateSuccess {
				wait.Add(1)
				maxDone = 1
			} else {
				wait.Add(2)
				maxDone = 2
			}
			wait.Wait()
			s.client.UnSubscribeTx(tx.hash)
		}
	} else {
		if s.expect != util.SendSuccess {
			s.client.UnSubscribeTx(tx.hash)
		}

		err, _ := jsonparser.GetString([]byte(response), "result", "engine_result")
		errorMessage, _ := jsonparser.GetString([]byte(response), "result", "engine_result_message")
		return &TxResult{
			Status:       util.SendError,
			TxHash:       tx.hash,
			ErrorCode:    err,
			ErrorMessage: errorMessage,
		}
	}

	return ret
}
