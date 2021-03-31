package core

import (
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/ChainSQL/go-chainsql-api/cgofuns"
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

// TxResult is tx submit response
type TxResult struct {
	Status       string `json:"status"`
	TxHash       string `json:"hash"`
	ErrorCode    string `json:"error,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// IPrepare is an interface that a struct call submit() method must implment
type IPrepare interface {
	PrepareTx() (TxJSON, error)
}

// SubmitBase base struct
type SubmitBase struct {
	expect   string
	callback export.Callback
	client   *net.Client
	IPrepare
}

// Submit submit a tx with a cocurrent expect
func (s *SubmitBase) Submit(cond string) string {
	s.expect = cond
	ret := s.doSubmit()
	jsonRet, _ := json.Marshal(ret)
	return string(jsonRet)
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

	jsonStr, err := json.Marshal(tx)
	// log.Printf("Tx json before sign:%s\n",string(jsonStr))
	if err != nil {
		log.Printf("doSubmit error:%s\n", err)
		return &TxResult{
			ErrorCode:    "errJsonMarshal",
			ErrorMessage: "Error when json.Marshal",
		}
	}

	o := new(cgofuns.CGOFun)
	var signedData []byte
	var txHash []byte
	signed := o.SignTransaction(s.client.Auth.Secret, string(jsonStr), &signedData, &txHash)
	if !signed {
		return &TxResult{
			ErrorCode:    "errSign",
			ErrorMessage: "Error when sign transaction",
		}
	}

	// log.Printf("Sign Result: hash=%s\n", string(txHash))

	txSigned := &TxSigned{
		blob: string(signedData),
		hash: string(txHash),
	}

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
