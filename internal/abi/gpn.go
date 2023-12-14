// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gpn

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// NetworkConsumerView is an auto generated low-level Go binding around an user-defined struct.
type NetworkConsumerView struct {
	Addr          common.Address
	Balance       *big.Int
	ContainerName string
}

// NetworkProverView is an auto generated low-level Go binding around an user-defined struct.
type NetworkProverView struct {
	Addr    common.Address
	Balance *big.Int
}

// ProvingNetworkMetaData contains all meta data concerning the ProvingNetwork contract.
var ProvingNetworkMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"depositEth\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_containerName\",\"type\":\"string\"}],\"name\":\"registerConsumer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registerProver\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"StringsInsufficientHexLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint8[]\",\"name\":\"vs\",\"type\":\"uint8[]\"}],\"name\":\"submitSignedProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawConsumer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"consumerAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"consumers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"containerName\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConsumers\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"containerName\",\"type\":\"string\"}],\"internalType\":\"structNetwork.ConsumerView[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProvers\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structNetwork.ProverView[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_ETH_AMOUNT_CONSUMER\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_ETH_AMOUNT_PROVER\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"payoutRequestIds\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"payouts\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"claimableAfterTimestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"proverAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"provers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ProvingNetworkABI is the input ABI used to generate the binding from.
// Deprecated: Use ProvingNetworkMetaData.ABI instead.
var ProvingNetworkABI = ProvingNetworkMetaData.ABI

// ProvingNetwork is an auto generated Go binding around an Ethereum contract.
type ProvingNetwork struct {
	ProvingNetworkCaller     // Read-only binding to the contract
	ProvingNetworkTransactor // Write-only binding to the contract
	ProvingNetworkFilterer   // Log filterer for contract events
}

// ProvingNetworkCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProvingNetworkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProvingNetworkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProvingNetworkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProvingNetworkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProvingNetworkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProvingNetworkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProvingNetworkSession struct {
	Contract     *ProvingNetwork   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProvingNetworkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProvingNetworkCallerSession struct {
	Contract *ProvingNetworkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ProvingNetworkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProvingNetworkTransactorSession struct {
	Contract     *ProvingNetworkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ProvingNetworkRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProvingNetworkRaw struct {
	Contract *ProvingNetwork // Generic contract binding to access the raw methods on
}

// ProvingNetworkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProvingNetworkCallerRaw struct {
	Contract *ProvingNetworkCaller // Generic read-only contract binding to access the raw methods on
}

// ProvingNetworkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProvingNetworkTransactorRaw struct {
	Contract *ProvingNetworkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProvingNetwork creates a new instance of ProvingNetwork, bound to a specific deployed contract.
func NewProvingNetwork(address common.Address, backend bind.ContractBackend) (*ProvingNetwork, error) {
	contract, err := bindProvingNetwork(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ProvingNetwork{ProvingNetworkCaller: ProvingNetworkCaller{contract: contract}, ProvingNetworkTransactor: ProvingNetworkTransactor{contract: contract}, ProvingNetworkFilterer: ProvingNetworkFilterer{contract: contract}}, nil
}

// NewProvingNetworkCaller creates a new read-only instance of ProvingNetwork, bound to a specific deployed contract.
func NewProvingNetworkCaller(address common.Address, caller bind.ContractCaller) (*ProvingNetworkCaller, error) {
	contract, err := bindProvingNetwork(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProvingNetworkCaller{contract: contract}, nil
}

// NewProvingNetworkTransactor creates a new write-only instance of ProvingNetwork, bound to a specific deployed contract.
func NewProvingNetworkTransactor(address common.Address, transactor bind.ContractTransactor) (*ProvingNetworkTransactor, error) {
	contract, err := bindProvingNetwork(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProvingNetworkTransactor{contract: contract}, nil
}

// NewProvingNetworkFilterer creates a new log filterer instance of ProvingNetwork, bound to a specific deployed contract.
func NewProvingNetworkFilterer(address common.Address, filterer bind.ContractFilterer) (*ProvingNetworkFilterer, error) {
	contract, err := bindProvingNetwork(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProvingNetworkFilterer{contract: contract}, nil
}

// bindProvingNetwork binds a generic wrapper to an already deployed contract.
func bindProvingNetwork(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProvingNetworkMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProvingNetwork *ProvingNetworkRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProvingNetwork.Contract.ProvingNetworkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProvingNetwork *ProvingNetworkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProvingNetwork.Contract.ProvingNetworkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProvingNetwork *ProvingNetworkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProvingNetwork.Contract.ProvingNetworkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProvingNetwork *ProvingNetworkCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProvingNetwork.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProvingNetwork *ProvingNetworkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProvingNetwork.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProvingNetwork *ProvingNetworkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProvingNetwork.Contract.contract.Transact(opts, method, params...)
}

// MINETHAMOUNTCONSUMER is a free data retrieval call binding the contract method 0xb99fd95a.
//
// Solidity: function MIN_ETH_AMOUNT_CONSUMER() view returns(uint256)
func (_ProvingNetwork *ProvingNetworkCaller) MINETHAMOUNTCONSUMER(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ProvingNetwork.contract.Call(opts, &out, "MIN_ETH_AMOUNT_CONSUMER")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINETHAMOUNTCONSUMER is a free data retrieval call binding the contract method 0xb99fd95a.
//
// Solidity: function MIN_ETH_AMOUNT_CONSUMER() view returns(uint256)
func (_ProvingNetwork *ProvingNetworkSession) MINETHAMOUNTCONSUMER() (*big.Int, error) {
	return _ProvingNetwork.Contract.MINETHAMOUNTCONSUMER(&_ProvingNetwork.CallOpts)
}

// MINETHAMOUNTCONSUMER is a free data retrieval call binding the contract method 0xb99fd95a.
//
// Solidity: function MIN_ETH_AMOUNT_CONSUMER() view returns(uint256)
func (_ProvingNetwork *ProvingNetworkCallerSession) MINETHAMOUNTCONSUMER() (*big.Int, error) {
	return _ProvingNetwork.Contract.MINETHAMOUNTCONSUMER(&_ProvingNetwork.CallOpts)
}

// MINETHAMOUNTPROVER is a free data retrieval call binding the contract method 0x4e448c2f.
//
// Solidity: function MIN_ETH_AMOUNT_PROVER() view returns(uint256)
func (_ProvingNetwork *ProvingNetworkCaller) MINETHAMOUNTPROVER(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ProvingNetwork.contract.Call(opts, &out, "MIN_ETH_AMOUNT_PROVER")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINETHAMOUNTPROVER is a free data retrieval call binding the contract method 0x4e448c2f.
//
// Solidity: function MIN_ETH_AMOUNT_PROVER() view returns(uint256)
func (_ProvingNetwork *ProvingNetworkSession) MINETHAMOUNTPROVER() (*big.Int, error) {
	return _ProvingNetwork.Contract.MINETHAMOUNTPROVER(&_ProvingNetwork.CallOpts)
}

// MINETHAMOUNTPROVER is a free data retrieval call binding the contract method 0x4e448c2f.
//
// Solidity: function MIN_ETH_AMOUNT_PROVER() view returns(uint256)
func (_ProvingNetwork *ProvingNetworkCallerSession) MINETHAMOUNTPROVER() (*big.Int, error) {
	return _ProvingNetwork.Contract.MINETHAMOUNTPROVER(&_ProvingNetwork.CallOpts)
}

// ConsumerAddresses is a free data retrieval call binding the contract method 0x47ee6e3c.
//
// Solidity: function consumerAddresses(uint256 ) view returns(address)
func (_ProvingNetwork *ProvingNetworkCaller) ConsumerAddresses(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ProvingNetwork.contract.Call(opts, &out, "consumerAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ConsumerAddresses is a free data retrieval call binding the contract method 0x47ee6e3c.
//
// Solidity: function consumerAddresses(uint256 ) view returns(address)
func (_ProvingNetwork *ProvingNetworkSession) ConsumerAddresses(arg0 *big.Int) (common.Address, error) {
	return _ProvingNetwork.Contract.ConsumerAddresses(&_ProvingNetwork.CallOpts, arg0)
}

// ConsumerAddresses is a free data retrieval call binding the contract method 0x47ee6e3c.
//
// Solidity: function consumerAddresses(uint256 ) view returns(address)
func (_ProvingNetwork *ProvingNetworkCallerSession) ConsumerAddresses(arg0 *big.Int) (common.Address, error) {
	return _ProvingNetwork.Contract.ConsumerAddresses(&_ProvingNetwork.CallOpts, arg0)
}

// Consumers is a free data retrieval call binding the contract method 0x0bf53668.
//
// Solidity: function consumers(address ) view returns(uint256 balance, string containerName)
func (_ProvingNetwork *ProvingNetworkCaller) Consumers(opts *bind.CallOpts, arg0 common.Address) (struct {
	Balance       *big.Int
	ContainerName string
}, error) {
	var out []interface{}
	err := _ProvingNetwork.contract.Call(opts, &out, "consumers", arg0)

	outstruct := new(struct {
		Balance       *big.Int
		ContainerName string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Balance = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ContainerName = *abi.ConvertType(out[1], new(string)).(*string)

	return *outstruct, err

}

// Consumers is a free data retrieval call binding the contract method 0x0bf53668.
//
// Solidity: function consumers(address ) view returns(uint256 balance, string containerName)
func (_ProvingNetwork *ProvingNetworkSession) Consumers(arg0 common.Address) (struct {
	Balance       *big.Int
	ContainerName string
}, error) {
	return _ProvingNetwork.Contract.Consumers(&_ProvingNetwork.CallOpts, arg0)
}

// Consumers is a free data retrieval call binding the contract method 0x0bf53668.
//
// Solidity: function consumers(address ) view returns(uint256 balance, string containerName)
func (_ProvingNetwork *ProvingNetworkCallerSession) Consumers(arg0 common.Address) (struct {
	Balance       *big.Int
	ContainerName string
}, error) {
	return _ProvingNetwork.Contract.Consumers(&_ProvingNetwork.CallOpts, arg0)
}

// GetConsumers is a free data retrieval call binding the contract method 0x3b729b86.
//
// Solidity: function getConsumers() view returns((address,uint256,string)[])
func (_ProvingNetwork *ProvingNetworkCaller) GetConsumers(opts *bind.CallOpts) ([]NetworkConsumerView, error) {
	var out []interface{}
	err := _ProvingNetwork.contract.Call(opts, &out, "getConsumers")

	if err != nil {
		return *new([]NetworkConsumerView), err
	}

	out0 := *abi.ConvertType(out[0], new([]NetworkConsumerView)).(*[]NetworkConsumerView)

	return out0, err

}

// GetConsumers is a free data retrieval call binding the contract method 0x3b729b86.
//
// Solidity: function getConsumers() view returns((address,uint256,string)[])
func (_ProvingNetwork *ProvingNetworkSession) GetConsumers() ([]NetworkConsumerView, error) {
	return _ProvingNetwork.Contract.GetConsumers(&_ProvingNetwork.CallOpts)
}

// GetConsumers is a free data retrieval call binding the contract method 0x3b729b86.
//
// Solidity: function getConsumers() view returns((address,uint256,string)[])
func (_ProvingNetwork *ProvingNetworkCallerSession) GetConsumers() ([]NetworkConsumerView, error) {
	return _ProvingNetwork.Contract.GetConsumers(&_ProvingNetwork.CallOpts)
}

// GetProvers is a free data retrieval call binding the contract method 0xc0bfd036.
//
// Solidity: function getProvers() view returns((address,uint256)[])
func (_ProvingNetwork *ProvingNetworkCaller) GetProvers(opts *bind.CallOpts) ([]NetworkProverView, error) {
	var out []interface{}
	err := _ProvingNetwork.contract.Call(opts, &out, "getProvers")

	if err != nil {
		return *new([]NetworkProverView), err
	}

	out0 := *abi.ConvertType(out[0], new([]NetworkProverView)).(*[]NetworkProverView)

	return out0, err

}

// GetProvers is a free data retrieval call binding the contract method 0xc0bfd036.
//
// Solidity: function getProvers() view returns((address,uint256)[])
func (_ProvingNetwork *ProvingNetworkSession) GetProvers() ([]NetworkProverView, error) {
	return _ProvingNetwork.Contract.GetProvers(&_ProvingNetwork.CallOpts)
}

// GetProvers is a free data retrieval call binding the contract method 0xc0bfd036.
//
// Solidity: function getProvers() view returns((address,uint256)[])
func (_ProvingNetwork *ProvingNetworkCallerSession) GetProvers() ([]NetworkProverView, error) {
	return _ProvingNetwork.Contract.GetProvers(&_ProvingNetwork.CallOpts)
}

// PayoutRequestIds is a free data retrieval call binding the contract method 0x15a0a89a.
//
// Solidity: function payoutRequestIds(uint256 ) view returns(string)
func (_ProvingNetwork *ProvingNetworkCaller) PayoutRequestIds(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _ProvingNetwork.contract.Call(opts, &out, "payoutRequestIds", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// PayoutRequestIds is a free data retrieval call binding the contract method 0x15a0a89a.
//
// Solidity: function payoutRequestIds(uint256 ) view returns(string)
func (_ProvingNetwork *ProvingNetworkSession) PayoutRequestIds(arg0 *big.Int) (string, error) {
	return _ProvingNetwork.Contract.PayoutRequestIds(&_ProvingNetwork.CallOpts, arg0)
}

// PayoutRequestIds is a free data retrieval call binding the contract method 0x15a0a89a.
//
// Solidity: function payoutRequestIds(uint256 ) view returns(string)
func (_ProvingNetwork *ProvingNetworkCallerSession) PayoutRequestIds(arg0 *big.Int) (string, error) {
	return _ProvingNetwork.Contract.PayoutRequestIds(&_ProvingNetwork.CallOpts, arg0)
}

// Payouts is a free data retrieval call binding the contract method 0xd7be0a06.
//
// Solidity: function payouts(string ) view returns(address consumer, uint256 claimableAfterTimestamp)
func (_ProvingNetwork *ProvingNetworkCaller) Payouts(opts *bind.CallOpts, arg0 string) (struct {
	Consumer                common.Address
	ClaimableAfterTimestamp *big.Int
}, error) {
	var out []interface{}
	err := _ProvingNetwork.contract.Call(opts, &out, "payouts", arg0)

	outstruct := new(struct {
		Consumer                common.Address
		ClaimableAfterTimestamp *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Consumer = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ClaimableAfterTimestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Payouts is a free data retrieval call binding the contract method 0xd7be0a06.
//
// Solidity: function payouts(string ) view returns(address consumer, uint256 claimableAfterTimestamp)
func (_ProvingNetwork *ProvingNetworkSession) Payouts(arg0 string) (struct {
	Consumer                common.Address
	ClaimableAfterTimestamp *big.Int
}, error) {
	return _ProvingNetwork.Contract.Payouts(&_ProvingNetwork.CallOpts, arg0)
}

// Payouts is a free data retrieval call binding the contract method 0xd7be0a06.
//
// Solidity: function payouts(string ) view returns(address consumer, uint256 claimableAfterTimestamp)
func (_ProvingNetwork *ProvingNetworkCallerSession) Payouts(arg0 string) (struct {
	Consumer                common.Address
	ClaimableAfterTimestamp *big.Int
}, error) {
	return _ProvingNetwork.Contract.Payouts(&_ProvingNetwork.CallOpts, arg0)
}

// ProverAddresses is a free data retrieval call binding the contract method 0xd2c7f2ac.
//
// Solidity: function proverAddresses(uint256 ) view returns(address)
func (_ProvingNetwork *ProvingNetworkCaller) ProverAddresses(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ProvingNetwork.contract.Call(opts, &out, "proverAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProverAddresses is a free data retrieval call binding the contract method 0xd2c7f2ac.
//
// Solidity: function proverAddresses(uint256 ) view returns(address)
func (_ProvingNetwork *ProvingNetworkSession) ProverAddresses(arg0 *big.Int) (common.Address, error) {
	return _ProvingNetwork.Contract.ProverAddresses(&_ProvingNetwork.CallOpts, arg0)
}

// ProverAddresses is a free data retrieval call binding the contract method 0xd2c7f2ac.
//
// Solidity: function proverAddresses(uint256 ) view returns(address)
func (_ProvingNetwork *ProvingNetworkCallerSession) ProverAddresses(arg0 *big.Int) (common.Address, error) {
	return _ProvingNetwork.Contract.ProverAddresses(&_ProvingNetwork.CallOpts, arg0)
}

// Provers is a free data retrieval call binding the contract method 0x1dec844b.
//
// Solidity: function provers(address ) view returns(uint256 balance)
func (_ProvingNetwork *ProvingNetworkCaller) Provers(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ProvingNetwork.contract.Call(opts, &out, "provers", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Provers is a free data retrieval call binding the contract method 0x1dec844b.
//
// Solidity: function provers(address ) view returns(uint256 balance)
func (_ProvingNetwork *ProvingNetworkSession) Provers(arg0 common.Address) (*big.Int, error) {
	return _ProvingNetwork.Contract.Provers(&_ProvingNetwork.CallOpts, arg0)
}

// Provers is a free data retrieval call binding the contract method 0x1dec844b.
//
// Solidity: function provers(address ) view returns(uint256 balance)
func (_ProvingNetwork *ProvingNetworkCallerSession) Provers(arg0 common.Address) (*big.Int, error) {
	return _ProvingNetwork.Contract.Provers(&_ProvingNetwork.CallOpts, arg0)
}

// DepositEth is a paid mutator transaction binding the contract method 0x439370b1.
//
// Solidity: function depositEth() payable returns()
func (_ProvingNetwork *ProvingNetworkTransactor) DepositEth(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProvingNetwork.contract.Transact(opts, "depositEth")
}

// DepositEth is a paid mutator transaction binding the contract method 0x439370b1.
//
// Solidity: function depositEth() payable returns()
func (_ProvingNetwork *ProvingNetworkSession) DepositEth() (*types.Transaction, error) {
	return _ProvingNetwork.Contract.DepositEth(&_ProvingNetwork.TransactOpts)
}

// DepositEth is a paid mutator transaction binding the contract method 0x439370b1.
//
// Solidity: function depositEth() payable returns()
func (_ProvingNetwork *ProvingNetworkTransactorSession) DepositEth() (*types.Transaction, error) {
	return _ProvingNetwork.Contract.DepositEth(&_ProvingNetwork.TransactOpts)
}

// RegisterConsumer is a paid mutator transaction binding the contract method 0x55bd3610.
//
// Solidity: function registerConsumer(string _containerName) payable returns()
func (_ProvingNetwork *ProvingNetworkTransactor) RegisterConsumer(opts *bind.TransactOpts, _containerName string) (*types.Transaction, error) {
	return _ProvingNetwork.contract.Transact(opts, "registerConsumer", _containerName)
}

// RegisterConsumer is a paid mutator transaction binding the contract method 0x55bd3610.
//
// Solidity: function registerConsumer(string _containerName) payable returns()
func (_ProvingNetwork *ProvingNetworkSession) RegisterConsumer(_containerName string) (*types.Transaction, error) {
	return _ProvingNetwork.Contract.RegisterConsumer(&_ProvingNetwork.TransactOpts, _containerName)
}

// RegisterConsumer is a paid mutator transaction binding the contract method 0x55bd3610.
//
// Solidity: function registerConsumer(string _containerName) payable returns()
func (_ProvingNetwork *ProvingNetworkTransactorSession) RegisterConsumer(_containerName string) (*types.Transaction, error) {
	return _ProvingNetwork.Contract.RegisterConsumer(&_ProvingNetwork.TransactOpts, _containerName)
}

// RegisterProver is a paid mutator transaction binding the contract method 0x4fab5637.
//
// Solidity: function registerProver() payable returns()
func (_ProvingNetwork *ProvingNetworkTransactor) RegisterProver(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProvingNetwork.contract.Transact(opts, "registerProver")
}

// RegisterProver is a paid mutator transaction binding the contract method 0x4fab5637.
//
// Solidity: function registerProver() payable returns()
func (_ProvingNetwork *ProvingNetworkSession) RegisterProver() (*types.Transaction, error) {
	return _ProvingNetwork.Contract.RegisterProver(&_ProvingNetwork.TransactOpts)
}

// RegisterProver is a paid mutator transaction binding the contract method 0x4fab5637.
//
// Solidity: function registerProver() payable returns()
func (_ProvingNetwork *ProvingNetworkTransactorSession) RegisterProver() (*types.Transaction, error) {
	return _ProvingNetwork.Contract.RegisterProver(&_ProvingNetwork.TransactOpts)
}

// SubmitSignedProof is a paid mutator transaction binding the contract method 0x52668861.
//
// Solidity: function submitSignedProof(string requestId, uint256 reward, bytes32[] rs, bytes32[] ss, uint8[] vs) returns()
func (_ProvingNetwork *ProvingNetworkTransactor) SubmitSignedProof(opts *bind.TransactOpts, requestId string, reward *big.Int, rs [][32]byte, ss [][32]byte, vs []uint8) (*types.Transaction, error) {
	return _ProvingNetwork.contract.Transact(opts, "submitSignedProof", requestId, reward, rs, ss, vs)
}

// SubmitSignedProof is a paid mutator transaction binding the contract method 0x52668861.
//
// Solidity: function submitSignedProof(string requestId, uint256 reward, bytes32[] rs, bytes32[] ss, uint8[] vs) returns()
func (_ProvingNetwork *ProvingNetworkSession) SubmitSignedProof(requestId string, reward *big.Int, rs [][32]byte, ss [][32]byte, vs []uint8) (*types.Transaction, error) {
	return _ProvingNetwork.Contract.SubmitSignedProof(&_ProvingNetwork.TransactOpts, requestId, reward, rs, ss, vs)
}

// SubmitSignedProof is a paid mutator transaction binding the contract method 0x52668861.
//
// Solidity: function submitSignedProof(string requestId, uint256 reward, bytes32[] rs, bytes32[] ss, uint8[] vs) returns()
func (_ProvingNetwork *ProvingNetworkTransactorSession) SubmitSignedProof(requestId string, reward *big.Int, rs [][32]byte, ss [][32]byte, vs []uint8) (*types.Transaction, error) {
	return _ProvingNetwork.Contract.SubmitSignedProof(&_ProvingNetwork.TransactOpts, requestId, reward, rs, ss, vs)
}

// WithdrawConsumer is a paid mutator transaction binding the contract method 0x13799aa9.
//
// Solidity: function withdrawConsumer() returns()
func (_ProvingNetwork *ProvingNetworkTransactor) WithdrawConsumer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProvingNetwork.contract.Transact(opts, "withdrawConsumer")
}

// WithdrawConsumer is a paid mutator transaction binding the contract method 0x13799aa9.
//
// Solidity: function withdrawConsumer() returns()
func (_ProvingNetwork *ProvingNetworkSession) WithdrawConsumer() (*types.Transaction, error) {
	return _ProvingNetwork.Contract.WithdrawConsumer(&_ProvingNetwork.TransactOpts)
}

// WithdrawConsumer is a paid mutator transaction binding the contract method 0x13799aa9.
//
// Solidity: function withdrawConsumer() returns()
func (_ProvingNetwork *ProvingNetworkTransactorSession) WithdrawConsumer() (*types.Transaction, error) {
	return _ProvingNetwork.Contract.WithdrawConsumer(&_ProvingNetwork.TransactOpts)
}
