package core

import (
	"fmt"
	"strconv"

	"github.com/buger/jsonparser"

	. "github.com/ChainSQL/go-chainsql-api/data"
	"github.com/ChainSQL/go-chainsql-api/net"
)

// Base is a struct
type Base struct {
}

//OpInfo is the opearting details
// type TransactionRequest struct {
// 	TransactionType string
// 	Amount          Amount
// 	Destination     string
// 	Query           []interface{}
// }

//TxInfo is the opearting details
type TxInfo struct {
	Raw    string
	TxType TransactionType
	Query  []interface{}
}

// type Amount struct {
// 	Value    string `json:"value"`
// 	Currency string `json:"currency"`
// 	Account  string `json:"account"`
// }

func (b *Base) Say() {
	fmt.Println("base")
}

func (b *Base) Say2() {
	fmt.Println("base2")
}

// Ripple is aaa
type Ripple struct {
	*Base
	client *net.Client
	op     *TxInfo
	SubmitBase
}

func (r *Ripple) Say() {
	fmt.Println("Ripple")
}

func NewRipple(client *net.Client) *Ripple {
	ripple := &Ripple{
		Base:   &Base{},
		client: client,
		op: &TxInfo{
			Query: make([]interface{}, 0),
		},
	}
	ripple.SubmitBase.client = ripple.client
	ripple.SubmitBase.IPrepare = ripple
	return ripple
}

func (r *Ripple) Pay(accountId string, value int64) *Ripple {
	r.op.TxType = PAYMENT

	r.op.Raw = "{\"AccountId\": \"" + accountId + "\", \"Value\": \"" + strconv.FormatInt(value, 10) + "\"}"
	return r
}

func (r *Ripple) pay(raw string) (Signer, error) {
	accountId, _ := jsonparser.GetString([]byte(raw), "AccountId")
	strValue, _ := jsonparser.GetString([]byte(raw), "Value")
	value, _ := strconv.ParseInt(strValue, 10, 64)

	valueTemp, _ := NewNativeValue(value)
	currency_zxc, _ := NewCurrency("ZXC")
	amount := Amount{
		Value:    valueTemp,
		Currency: currency_zxc,
	}
	return r.PayToNode(accountId, amount)
}

func (r *Ripple) PayToNode(accountId string, amount Amount) (Signer, error) {

	// if !amount.Currency.IsNative() {
	// 	accountData, err := r.client.GetAccountInfo(string(amount.Issuer))
	// 	if err != nil {
	// 		log.Println("get issuer %s", err)
	// 	}

	// 	if accountData != "" {
	// 		//var feeMin, feeMax = "", ""
	// 		//var lFeeRate = Value(0)
	// 		var mapObj map[string]interface{}
	// 		va := amount.Value
	// 		json.Unmarshal([]byte(accountData), &mapObj)
	// 		feeMin := mapObj["TransferFeeMin"].(Value)
	// 		feeMax := mapObj["TransferFeeMax"].(Value)
	// 		lFeeRate := mapObj["TransferRate"].(Value)
	// 		fee := Value()
	// 		if feeMin.IsZero() || feeMax.IsZero() || lFeeRate.IsZero() {
	// 			if feeMin == feeMax {
	// 				fee = feeMin.Float()
	// 			} else if !lFeeRate.IsZero() {
	// 				//	fee = FloatOperation.accMul(parseFloat(value), data.rate - 1);
	// 				fee = va.Multiply(lFeeRate)
	// 				if !feeMin.IsZero() {
	// 					if
	// 					fee =
	// 				}
	// 				if feeMax.IsZero() {
	// 					fee = Math.min(fee, parseFloat(data.max))
	// 				}
	// 				//
	// 				value = value.add(fee)
	// 			}
	// 		}
	// 	}
	// }
	destination, err := NewAccountFromAddress(accountId)
	if err != nil {
		return nil, err
	}
	payment := &Payment{
		//SendMax:     nil,
		Destination: *destination,
		Amount:      amount,
	}
	payment.TransactionType = PAYMENT
	account, err := NewAccountFromAddress(r.client.Auth.Address)
	if err != nil {
		return nil, err
	}
	payment.Account = *account
	seq, err := net.PrepareRipple(r.client)
	if err != nil {
		return nil, err
	}
	payment.Sequence = seq
	valueTemp, _ := NewNativeValue(10)
	payment.Fee = *valueTemp
	var sign Signer = payment
	return sign, nil
}

//PrepareTx prepare tx json for submit
func (r *Ripple) PrepareTx() (Signer, error) {
	var tx Signer
	var err error
	switch r.op.TxType {
	case PAYMENT:
		tx, err = r.pay(r.op.Raw)
		break
	default:
	}
	if err != nil {
		return nil, err
	}
	return tx, nil

}
