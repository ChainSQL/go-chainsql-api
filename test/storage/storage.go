// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package storage

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ChainSQL/go-chainsql-api/abigen/abi"
	"github.com/ChainSQL/go-chainsql-api/abigen/abi/bind"
	"github.com/ChainSQL/go-chainsql-api/common"
	"github.com/ChainSQL/go-chainsql-api/core"
	"github.com/ChainSQL/go-chainsql-api/data"
	"github.com/ChainSQL/go-chainsql-api/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = bind.Bind
	_ = common.Big1
)

// StorageMetaData contains all meta data concerning the Storage contract.
var StorageMetaData = &core.CtrMetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"numberChanges\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060bd8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e64cec11460375780636057361d14604c575b600080fd5b60005460405190815260200160405180910390f35b605b60573660046098565b605d565b005b60008190556040518181527f1161d67e3e40d64da0f22f41054120b745a28aa25e65d98d153fbaf4d31952519060200160405180910390a150565b60006020828403121560a957600080fd5b503591905056fea164736f6c6343000805000a",
}

// StorageABI is the input ABI used to generate the binding from.
// Deprecated: Use StorageMetaData.ABI instead.
var StorageABI = StorageMetaData.ABI

// StorageBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StorageMetaData.Bin instead.
var StorageBin = StorageMetaData.Bin

// DeployStorage deploys a new ChainSQL contract, binding an instance of Storage to it.
func DeployStorage(chainsql *core.Chainsql, auth *core.TransactOpts) (*core.DeployTxRet, *Storage, error) {
	parsed, err := StorageMetaData.GetAbi()
	if err != nil {
		return &core.DeployTxRet{}, nil, err
	}
	if parsed == nil {
		return &core.DeployTxRet{}, nil, errors.New("GetABI returned nil")
	}

	deployRet, contract, err := core.DeployContract(chainsql, auth, *parsed, common.FromHex(StorageBin))
	if err != nil {
		return &core.DeployTxRet{}, nil, err
	}
	return deployRet, &Storage{StorageCaller: StorageCaller{contract: contract}, StorageTransactor: StorageTransactor{contract: contract}, StorageFilterer: StorageFilterer{contract: contract}}, nil
}

// Storage is an auto generated Go binding around an ChainSQL contract.
type Storage struct {
	StorageCaller     // Read-only binding to the contract
	StorageTransactor // Write-only binding to the contract
	StorageFilterer   // Log filterer for contract events
}

// StorageCaller is an auto generated read-only Go binding around an ChainSQL contract.
type StorageCaller struct {
	contract *core.BoundContract // Generic contract wrapper for the low level calls
}

// StorageTransactor is an auto generated write-only Go binding around an ChainSQL contract.
type StorageTransactor struct {
	contract *core.BoundContract // Generic contract wrapper for the low level calls
}

// StorageFilterer is an auto generated log filtering Go binding around an ChainSQL contract events.
type StorageFilterer struct {
	contract *core.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSession is an auto generated Go binding around an ChainSQL contract,
// with pre-set call and transact options.
type StorageSession struct {
	Contract     *Storage          // Generic contract binding to set the session for
	CallOpts     core.CallOpts     // Call options to use throughout this session
	TransactOpts core.TransactOpts // Transaction auth options to use throughout this session
}

// StorageCallerSession is an auto generated read-only Go binding around an ChainSQL contract,
// with pre-set call options.
type StorageCallerSession struct {
	Contract *StorageCaller // Generic contract caller binding to set the session for
	CallOpts core.CallOpts  // Call options to use throughout this session
}

// StorageTransactorSession is an auto generated write-only Go binding around an ChainSQL contract,
// with pre-set transact options.
type StorageTransactorSession struct {
	Contract     *StorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts core.TransactOpts  // Transaction auth options to use throughout this session
}

// StorageRaw is an auto generated low-level Go binding around an ChainSQL contract.
type StorageRaw struct {
	Contract *Storage // Generic contract binding to access the raw methods on
}

// StorageCallerRaw is an auto generated low-level read-only Go binding around an ChainSQL contract.
type StorageCallerRaw struct {
	Contract *StorageCaller // Generic read-only contract binding to access the raw methods on
}

// StorageTransactorRaw is an auto generated low-level write-only Go binding around an ChainSQL contract.
type StorageTransactorRaw struct {
	Contract *StorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorage creates a new instance of Storage, bound to a specific deployed contract.
func NewStorage(chainsql *core.Chainsql, address string) (*Storage, error) {
	contract, err := bindStorage(chainsql, address)
	if err != nil {
		return nil, err
	}
	return &Storage{StorageCaller: StorageCaller{contract: contract}, StorageTransactor: StorageTransactor{contract: contract}, StorageFilterer: StorageFilterer{contract: contract}}, nil
}

// // NewStorageCaller creates a new read-only instance of Storage, bound to a specific deployed contract.
// func NewStorageCaller(address common.Address, caller bind.ContractCaller) (*StorageCaller, error) {
//   contract, err := bindStorage(address, caller, nil, nil)
//   if err != nil {
//     return nil, err
//   }
//   return &StorageCaller{contract: contract}, nil
// }

// // NewStorageTransactor creates a new write-only instance of Storage, bound to a specific deployed contract.
// func NewStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageTransactor, error) {
//   contract, err := bindStorage(address, nil, transactor, nil)
//   if err != nil {
//     return nil, err
//   }
//   return &StorageTransactor{contract: contract}, nil
// }

// // NewStorageFilterer creates a new log filterer instance of Storage, bound to a specific deployed contract.
// func NewStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageFilterer, error) {
//   contract, err := bindStorage(address, nil, nil, filterer)
//   if err != nil {
//     return nil, err
//   }
//   return &StorageFilterer{contract: contract}, nil
// }

// bindStorage binds a generic wrapper to an already deployed contract.
func bindStorage(chainsql *core.Chainsql, address string) (*core.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StorageABI))
	if err != nil {
		return nil, err
	}
	return core.NewBoundContract(chainsql, address, parsed), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
// func (_Storage *StorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
// 	return _Storage.Contract.StorageCaller.contract.Call(opts, result, method, params...)
// }

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
// func (_Storage *StorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
// 	return _Storage.Contract.StorageTransactor.contract.Transfer(opts)
// }

// Transact invokes the (paid) contract method with params as input values.
// func (_Storage *StorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
// 	return _Storage.Contract.StorageTransactor.contract.Transact(opts, method, params...)
// }

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
// func (_Storage *StorageCallerRaw) Call(opts *core.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
// 	return _Storage.Contract.contract.Call(opts, result, method, params...)
// }

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
// func (_Storage *StorageTransactorRaw) Transfer(opts *core.TransactOpts) (*types.Transaction, error) {
// 	return _Storage.Contract.contract.Transfer(opts)
// }

// Transact invokes the (paid) contract method with params as input values.
// func (_Storage *StorageTransactorRaw) Transact(opts *core.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
// 	return _Storage.Contract.contract.Transact(opts, method, params...)
// }

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_Storage *StorageCaller) Retrieve(opts *core.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Storage.contract.Call(opts, &out, "retrieve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_Storage *StorageSession) Retrieve() (*big.Int, error) {
	return _Storage.Contract.Retrieve(&_Storage.CallOpts)
}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_Storage *StorageCallerSession) Retrieve() (*big.Int, error) {
	return _Storage.Contract.Retrieve(&_Storage.CallOpts)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_Storage *StorageTransactor) Store(opts *core.TransactOpts, num *big.Int) (*common.TxResult, error) {
	return _Storage.contract.Transact(opts, "store", num)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_Storage *StorageSession) Store(num *big.Int) (*common.TxResult, error) {
	return _Storage.Contract.Store(&_Storage.TransactOpts, num)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_Storage *StorageTransactorSession) Store(num *big.Int) (*common.TxResult, error) {
	return _Storage.Contract.Store(&_Storage.TransactOpts, num)
}

// StorageNumberChangesIterator is returned from FilterNumberChanges and is used to iterate over the raw logs and unpacked data for NumberChanges events raised by the Storage contract.
type StorageNumberChangesIterator struct {
	Event *StorageNumberChanges // Event containing the contract specifics and raw log

	contract *core.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan data.Log      // Log channel receiving the found contract events
	sub  event.Subscription // Subscription for errors, completion and termination
	done bool               // Whether the subscription completed delivering logs
	fail error              // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StorageNumberChangesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StorageNumberChanges)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StorageNumberChanges)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StorageNumberChangesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StorageNumberChangesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StorageNumberChanges represents a NumberChanges event raised by the Storage contract.
type StorageNumberChanges struct {
	Num *big.Int
	Raw data.Log // Blockchain specific contextual infos
}

// FilterNumberChanges is a free log retrieval operation binding the contract event 0x1161d67e3e40d64da0f22f41054120b745a28aa25e65d98d153fbaf4d3195251.
//
// Solidity: event numberChanges(uint256 num)
// func (_Storage *StorageFilterer) FilterNumberChanges(opts *core.FilterOpts) (*StorageNumberChangesIterator, error) {
//
//

// 	logs, sub, err := _Storage.contract.FilterLogs(opts, "numberChanges")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &StorageNumberChangesIterator{contract: _Storage.contract, event: "numberChanges", logs: logs, sub: sub}, nil
// }

// WatchNumberChanges is a free log subscription operation binding the contract event 0x1161d67e3e40d64da0f22f41054120b745a28aa25e65d98d153fbaf4d3195251.
//
// Solidity: event numberChanges(uint256 num)
func (_Storage *StorageFilterer) WatchNumberChanges(opts *core.WatchOpts, sink chan<- *StorageNumberChanges) (event.Subscription, error) {

	sub, err := _Storage.contract.WatchLogs(opts, "numberChanges")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.UnSubscribe()
		for {
			select {
			case log := <-sub.EventMsgCh:
				// New log arrived, parse the event and forward to the user
				event := new(StorageNumberChanges)
				if err := _Storage.contract.UnpackLog(event, "numberChanges", *log); err != nil {
					return err
				}
				event.Raw = *log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNumberChanges is a log parse operation binding the contract event 0x1161d67e3e40d64da0f22f41054120b745a28aa25e65d98d153fbaf4d3195251.
//
// Solidity: event numberChanges(uint256 num)
func (_Storage *StorageFilterer) ParseNumberChanges(log data.Log) (*StorageNumberChanges, error) {
	event := new(StorageNumberChanges)
	if err := _Storage.contract.UnpackLog(event, "numberChanges", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
