// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ChainSQL/go-chainsql-api/abigen/abi"
	. "github.com/ChainSQL/go-chainsql-api/abigen/abi/bind"
	"github.com/ChainSQL/go-chainsql-api/common"
	"github.com/ChainSQL/go-chainsql-api/crypto"
	"github.com/ChainSQL/go-chainsql-api/data"
	. "github.com/ChainSQL/go-chainsql-api/data"
	"github.com/ChainSQL/go-chainsql-api/export"
	"github.com/ChainSQL/go-chainsql-api/net"
	"github.com/buger/jsonparser"
)

// SignerFn is a signer function callback when a contract requires a method to
// sign the transaction before submission.
// type SignerFn func(common.Address, *types.Transaction) (*types.Transaction, error)

// CallOpts is the collection of options to fine tune a contract call request.
type CallOpts struct {
	LedgerIndex int64 // Optional the block number on which the call should be performed
}

type CallReq struct {
	common.RequestBase
	Account         string `json:"account"`
	ContractAddress string `json:"contract_address"`
	ContractData    string `json:"contract_data"`
	LedgerIndex     uint32 `json:"ledger_index"`
}

// TransactOpts is the collection of authorization data required to create a
// valid ChainSQL transaction.
type TransactOpts struct {
	ContractValue int64 // Funds to transfer along the transaction (nil = 0 = no funds)
	Gas           uint32
	Expectation   string
}

type DeployTxRet struct {
	common.TxResult
	ContractAddress string `json:"contractAddress"`
}

// FilterOpts is the collection of options to fine tune filtering for events
// within a bound contract.
type FilterOpts struct {
	Start uint64  // Start of the queried range
	End   *uint64 // End of the range (nil = latest)

	Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

// WatchOpts is the collection of options to fine tune subscribing for events
// within a bound contract.
type WatchOpts struct {
	Start   *uint64         // Start of the queried range (nil = latest)
	Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

// CtrMetaData collects all metadata for a bound contract.
type CtrMetaData struct {
	mu   sync.Mutex
	Sigs map[string]string
	Bin  string
	ABI  string
	ab   *abi.ABI
}

func (m *CtrMetaData) GetAbi() (*abi.ABI, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ab != nil {
		return m.ab, nil
	}
	if parsed, err := abi.JSON(strings.NewReader(m.ABI)); err != nil {
		return nil, err
	} else {
		m.ab = &parsed
	}
	return m.ab, nil
}

// BoundContract is the base wrapper object that reflects a contract on the
// Ethereum network. It contains a collection of methods that are used by the
// higher level contract bindings to operate.
type BoundContract struct {
	SubmitBase
	address    string             // Deployment address of the contract on the ChainSQL blockchain
	abi        abi.ABI            // Reflect based ABI to access the correct ChainSQL methods
	caller     ContractCaller     // Read interface to interact with the blockchain
	transactor ContractTransactor // Write interface to interact with the blockchain
	filterer   ContractFilterer   // Event filtering to interact with the blockchain
	TransactOpts
	ContractOpType   uint16
	ContractData     []byte
	isFirstSubscribe bool
	ctrEventCache    map[string]export.Callback
}

// NewBoundContract creates a low level contract interface through which calls
// and transactions may be made through.
// func NewBoundContract(chainsql *Chainsql, address string, abi abi.ABI, caller ContractCaller, transactor ContractTransactor, filterer ContractFilterer) *BoundContract {
func NewBoundContract(chainsql *Chainsql, address string, abi abi.ABI) *BoundContract {
	// bCtr := &BoundContract{
	// 	address:    address,
	// 	abi:        abi,
	// 	caller:     caller,
	// 	transactor: transactor,
	// 	filterer:   filterer,
	// }
	bCtr := &BoundContract{
		address:          address,
		abi:              abi,
		ContractOpType:   2,
		isFirstSubscribe: true,
	}
	bCtr.client = chainsql.client
	bCtr.IPrepare = bCtr
	return bCtr
}

// DeployContract deploys a contract onto the ChainSQL blockchain and binds the
// deployment address with a Go wrapper.
// func DeployContract(chainsql *Chainsql, opts *TransactOpts, abi abi.ABI, bytecode []byte, backend ContractBackend, params ...interface{}) (*DeployTxRet, *BoundContract, error) {
func DeployContract(chainsql *Chainsql, opts *TransactOpts, abi abi.ABI, bytecode []byte, params ...interface{}) (*DeployTxRet, *BoundContract, error) {
	// Otherwise try to deploy the contract
	// c := NewBoundContract(chainsql, "", abi, backend, backend, backend)
	c := NewBoundContract(chainsql, "", abi)
	c.ContractOpType = 1

	input, err := c.abi.Pack("", params...)
	if err != nil {
		return nil, nil, err
	}
	txRet, err := c.transact(opts, append(bytecode, input...))
	if err != nil {
		return nil, nil, err
	}

	deployTxRet := &DeployTxRet{}
	deployTxRet.Status = txRet.Status
	deployTxRet.TxHash = txRet.TxHash

	ret, err := c.client.GetTransaction(deployTxRet.TxHash)
	if err == nil {
		txAddr, err := jsonparser.GetString([]byte(ret), "result", "Account")
		if err != nil {
			return deployTxRet, c, err
		}
		txAddrSeq, err := jsonparser.GetInt([]byte(ret), "result", "Sequence")
		if err != nil {
			return deployTxRet, c, err
		}
		txRetMeta, _, _, err := jsonparser.Get([]byte(ret), "result", "meta", "AffectedNodes")
		if err != nil {
			return deployTxRet, c, err
		}

		var txDetailCtrAddr string
		_, _ = jsonparser.ArrayEach(txRetMeta, func(value []byte, dataType jsonparser.ValueType, offset int, errin error) {
			ctrAddr, err := jsonparser.GetString(value, "CreatedNode", "NewFields", "Account")
			if err != nil {
				return
			}
			txDetailCtrAddr = ctrAddr
			// return
		})
		c.address, err = crypto.CreateContractAddr(txAddr, uint32(txAddrSeq))
		if err != nil {
			deployTxRet.ErrorMessage = err.Error()
			return deployTxRet, c, err
		}
		if txDetailCtrAddr == "" || c.address != txDetailCtrAddr {
			err := errors.New("mismatch, can't find correct contract address")
			deployTxRet.ErrorMessage = err.Error()
			return deployTxRet, c, err
		}
		deployTxRet.ContractAddress = c.address
	} else {
		err := errors.New("Get Tx detail failed, can't find correct contract address")
		deployTxRet.ErrorMessage = err.Error()
		return deployTxRet, c, err
	}

	return deployTxRet, c, nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (c *BoundContract) Call(opts *CallOpts, results *[]interface{}, method string, params ...interface{}) error {
	// Don't crash on a lazy user
	if opts == nil {
		opts = new(CallOpts)
	}
	if results == nil {
		results = new([]interface{})
	}
	// Pack the input, call and unpack the results
	input, err := c.abi.Pack(method, params...)
	if err != nil {
		return err
	}
	inputHexStr := fmt.Sprintf("%x", input)
	callReq := &CallReq{
		Account:         c.client.Auth.Address,
		ContractAddress: c.address,
		ContractData:    inputHexStr,
	}
	if opts.LedgerIndex == 0 {
		_, lastLedgerSeq, _ := net.PrepareLastLedgerSeqAndFee(c.client)

		callReq.LedgerIndex = lastLedgerSeq - 20
	}
	// c.client.cmdIDs++
	// callReq.ID = 2
	callReq.Command = "contract_call"
	request := c.client.SyncRequest(callReq)

	err = c.client.ParseResponseError(request)
	if err != nil {
		return err
	}

	contractCallRet, err := jsonparser.GetString([]byte(request.Response.Value), "result", "contract_call_result")
	if err != nil {
		return err
	}
	ctrCallRetHex, err := hex.DecodeString(contractCallRet[2:])
	if err != nil {
		return err
	}
	if len(*results) == 0 {
		res, err := c.abi.Unpack(method, ctrCallRetHex)
		*results = res
		return err
	}
	res := *results
	return c.abi.UnpackIntoInterface(res[0], method, ctrCallRetHex)
}

// Transact invokes the (paid) contract method with params as input values.
// func (c *BoundContract) Transact(opts *TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
func (c *BoundContract) Transact(opts *TransactOpts, method string, params ...interface{}) (*common.TxResult, error) {
	// Otherwise pack up the parameters and invoke the contract
	input, err := c.abi.Pack(method, params...)
	if err != nil {
		return nil, err
	}
	// todo(rjl493456442) check the method is payable or not,
	// reject invalid transaction at the first place
	// return c.transact(opts, &c.address, input)
	return c.transact(opts, input)
}

// // RawTransact initiates a transaction with the given raw calldata as the input.
// // It's usually used to initiate transactions for invoking **Fallback** function.
// func (c *BoundContract) RawTransact(opts *TransactOpts, calldata []byte) (*types.Transaction, error) {
// 	// todo(rjl493456442) check the method is payable or not,
// 	// reject invalid transaction at the first place
// 	return c.transact(opts, &c.address, calldata)
// }

// // Transfer initiates a plain transaction to move funds to the contract, calling
// // its default method if one is available.
// func (c *BoundContract) Transfer(opts *TransactOpts) (*types.Transaction, error) {
// 	// todo(rjl493456442) check the payable fallback or receive is defined
// 	// or not, reject invalid transaction at the first place
// 	return c.transact(opts, &c.address, nil)
// }

// func (c *BoundContract) estimateGasLimit(opts *TransactOpts, contract *common.Address, input []byte, gasPrice, gasTipCap, gasFeeCap, value *big.Int) (uint64, error) {
// 	if contract != nil {
// 		// Gas estimation cannot succeed without code for method invocations.
// 		if code, err := c.transactor.PendingCodeAt(ensureContext(opts.Context), c.address); err != nil {
// 			return 0, err
// 		} else if len(code) == 0 {
// 			return 0, ErrNoCode
// 		}
// 	}
// 	msg := ethereum.CallMsg{
// 		From:      opts.From,
// 		To:        contract,
// 		GasPrice:  gasPrice,
// 		GasTipCap: gasTipCap,
// 		GasFeeCap: gasFeeCap,
// 		Value:     value,
// 		Data:      input,
// 	}
// 	return c.transactor.EstimateGas(ensureContext(opts.Context), msg)
// }

// func (c *BoundContract) getNonce(opts *TransactOpts) (uint64, error) {
// 	if opts.Nonce == nil {
// 		return c.transactor.PendingNonceAt(ensureContext(opts.Context), opts.From)
// 	} else {
// 		return opts.Nonce.Uint64(), nil
// 	}
// }

func (c *BoundContract) PrepareTx() (Signer, error) {
	contractTxObj := &ContractTx{}

	contractTxObj.TransactionType = CONTRACT
	account, err := NewAccountFromAddress(c.client.Auth.Address)
	if err != nil {
		return nil, err
	}
	contractTxObj.Account = *account
	seq, err := net.PrepareRipple(c.client)
	if err != nil {
		return nil, err
	}
	contractTxObj.Sequence = seq
	contractTxObj.ContractOpType = c.ContractOpType

	// if contractTxObj.ContractOpType == 1 {
	// 	var zeroAccount Account
	// 	contractTxObj.ContractAddress = zeroAccount
	// } else {
	if contractTxObj.ContractOpType == 2 {
		contracAddress, err := NewAccountFromAddress(c.address)
		if err != nil {
			return nil, err
		}
		contractTxObj.ContractAddress = *contracAddress
	}

	contractTxObj.ContractData = c.ContractData

	contractValue, _ := NewNativeValue(c.ContractValue)
	currencyZxc, _ := NewCurrency("ZXC")
	contractAmount := Amount{
		Value:    contractValue,
		Currency: currencyZxc,
	}
	contractTxObj.ContractValue = contractAmount
	contractTxObj.Gas = c.Gas

	fee, lastLedgerSeq, err := net.PrepareLastLedgerSeqAndFee(c.client)
	if err != nil {
		return nil, err
	}

	contractTxObj.LastLedgerSequence = &lastLedgerSeq
	finalFee, err := NewNativeValue(fee)
	if err != nil {
		return nil, err
	}
	contractTxObj.Fee = *finalFee

	return contractTxObj, nil
}

// transact executes an actual transaction invocation, first deriving any missing
// authorization fields, and then scheduling the transaction for execution.
// func (c *BoundContract) transact(opts *TransactOpts, contract *common.Address, input []byte) (*common.TxResult, error) {
func (c *BoundContract) transact(opts *TransactOpts, input []byte) (*common.TxResult, error) {
	// Create the transaction
	c.ContractValue = opts.ContractValue
	c.Gas = opts.Gas
	c.ContractData = input

	return c.Submit(opts.Expectation), nil
}

// FilterLogs filters contract logs for past blocks, returning the necessary
// channels to construct a strongly typed bound iterator on top of them.
// func (c *BoundContract) FilterLogs(opts *FilterOpts, name string, query ...[]interface{}) (chan types.Log, event.Subscription, error) {
// 	// Don't crash on a lazy user
// 	if opts == nil {
// 		opts = new(FilterOpts)
// 	}
// 	// Append the event selector to the query parameters and construct the topic set
// 	query = append([][]interface{}{{c.abi.Events[name].ID}}, query...)

// 	topics, err := abi.MakeTopics(query...)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	// Start the background filtering
// 	logs := make(chan types.Log, 128)

// 	config := ethereum.FilterQuery{
// 		Addresses: []common.Address{c.address},
// 		Topics:    topics,
// 		FromBlock: new(big.Int).SetUint64(opts.Start),
// 	}
// 	if opts.End != nil {
// 		config.ToBlock = new(big.Int).SetUint64(*opts.End)
// 	}
// 	/* TODO(karalabe): Replace the rest of the method below with this when supported
// 	sub, err := c.filterer.SubscribeFilterLogs(ensureContext(opts.Context), config, logs)
// 	*/
// 	buff, err := c.filterer.FilterLogs(ensureContext(opts.Context), config)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	sub, err := event.NewSubscription(func(quit <-chan struct{}) error {
// 		for _, log := range buff {
// 			select {
// 			case logs <- log:
// 			case <-quit:
// 				return nil
// 			}
// 		}
// 		return nil
// 	}), nil

// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return logs, sub, nil
// }
type EventSub struct {
	eventSign  string
	EventMsgCh chan *data.Log
	err        chan error
	ctr        *BoundContract
}

func (e *EventSub) UnSubscribe() {
	e.ctr.client.UnRegisterCtrEvent(e.eventSign, e.EventMsgCh)
}

func (e *EventSub) Err() chan error {
	return e.err
}

// WatchLogs filters subscribes to contract logs for future blocks, returning a
// subscription object that can be used to tear down the watcher.
// func (c *BoundContract) WatchLogs(opts *WatchOpts, name string, query ...[]interface{}) (chan data.Log, event.Subscription, error) {
func (c *BoundContract) WatchLogs(opts *WatchOpts, name string, query ...[]interface{}) (*EventSub, error) {
	// Don't crash on a lazy user
	if opts == nil {
		opts = new(WatchOpts)
	}
	// Append the event selector to the query parameters and construct the topic set
	query = append([][]interface{}{{c.abi.Events[name].ID}}, query...)

	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return nil, err
	}
	log.Println(topics[0][0])

	// Start the background filtering
	// logs := make(chan data.Log, 128)

	if c.isFirstSubscribe {
		c.client.SubscribeCtrAddr(c.address, true)
		c.isFirstSubscribe = false
	}

	eventMsgCh := make(chan *data.Log)
	errCh := make(chan error)
	eventSig := topics[0][0].String()
	c.client.RegisterCtrEvent(eventSig, eventMsgCh)

	sub := &EventSub{
		eventSign:  eventSig,
		EventMsgCh: eventMsgCh,
		ctr:        c,
		err:        errCh,
	}

	return sub, nil
}

// UnpackLog unpacks a retrieved log into the provided output structure.
func (c *BoundContract) UnpackLog(out interface{}, event string, log data.Log) error {
	if log.Topics[0] != c.abi.Events[event].ID {
		return fmt.Errorf("event signature mismatch")
	}
	if len(log.Data) > 0 {
		if err := c.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return abi.ParseTopics(out, indexed, log.Topics[1:])
}

// // UnpackLogIntoMap unpacks a retrieved log into the provided map.
// func (c *BoundContract) UnpackLogIntoMap(out map[string]interface{}, event string, log types.Log) error {
// 	if log.Topics[0] != c.abi.Events[event].ID {
// 		return fmt.Errorf("event signature mismatch")
// 	}
// 	if len(log.Data) > 0 {
// 		if err := c.abi.UnpackIntoMap(out, event, log.Data); err != nil {
// 			return err
// 		}
// 	}
// 	var indexed abi.Arguments
// 	for _, arg := range c.abi.Events[event].Inputs {
// 		if arg.Indexed {
// 			indexed = append(indexed, arg)
// 		}
// 	}
// 	return abi.ParseTopicsIntoMap(out, indexed, log.Topics[1:])
// }

// // ensureContext is a helper method to ensure a context is not nil, even if the
// // user specified it as such.
// func ensureContext(ctx context.Context) context.Context {
// 	if ctx == nil {
// 		return context.Background()
// 	}
// 	return ctx
// }
