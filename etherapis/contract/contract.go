// This file is an automatically generated Go binding. Do not modify as any
// change will likely be lost upon the next re-generation!

package contract

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EtherAPIsABI is the input ABI used to generate the binding from.
const EtherAPIsABI = `[{"constant":true,"inputs":[{"name":"from","type":"address"},{"name":"serviceID","type":"uint256"}],"name":"makeSubscriptionID","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":false,"inputs":[{"name":"serviceID","type":"uint256"}],"name":"subscribe","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"addr","type":"address"}],"name":"userServicesLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[],"name":"servicesLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionID","type":"bytes32"}],"name":"getSubscription","outputs":[{"name":"from","type":"address"},{"name":"serviceID","type":"uint256"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"cancelled","type":"bool"},{"name":"closedAt","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"name","type":"string"},{"name":"endpoint","type":"string"},{"name":"price","type":"uint256"},{"name":"cancellation","type":"uint256"}],"name":"addService","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionID","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"claim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"uint256"}],"name":"userServices","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionID","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verifyPayment","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"serviceID","type":"uint256"}],"name":"deleteService","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"serviceID","type":"uint256"}],"name":"enableService","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionID","type":"bytes32"}],"name":"getSubscriptionClosedAt","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionID","type":"bytes32"}],"name":"getSubscriptionNonce","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"from","type":"address"},{"name":"serviceID","type":"uint256"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"}],"name":"getHash","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":false,"inputs":[{"name":"serviceID","type":"uint256"}],"name":"disableService","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionID","type":"bytes32"}],"name":"getSubscriptionOwner","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionID","type":"bytes32"}],"name":"reclaim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionID","type":"bytes32"}],"name":"getSubscriptionValue","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionID","type":"bytes32"}],"name":"deposit","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionID","type":"bytes32"}],"name":"cancel","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"uint256"}],"name":"userSubscriptions","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":true,"inputs":[{"name":"addr","type":"address"}],"name":"userSubscriptionsLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionID","type":"bytes32"}],"name":"isValidSubscription","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionID","type":"bytes32"}],"name":"getSubscriptionServiceID","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"serviceID","type":"uint256"}],"name":"getService","outputs":[{"name":"name","type":"string"},{"name":"owner","type":"address"},{"name":"endpoint","type":"string"},{"name":"price","type":"uint256"},{"name":"cancellation","type":"uint256"},{"name":"enabled","type":"bool"},{"name":"deleted","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionID","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verifySignature","outputs":[{"name":"","type":"bool"}],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"name","type":"string"},{"indexed":true,"name":"owner","type":"address"},{"indexed":false,"name":"serviceID","type":"uint256"}],"name":"NewService","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"serviceID","type":"uint256"}],"name":"UpdateService","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"serviceID","type":"uint256"},{"indexed":false,"name":"subscriptionID","type":"bytes32"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"NewSubscription","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"subscriptionID","type":"bytes32"}],"name":"Deposit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionID","type":"bytes32"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"Redeem","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionID","type":"bytes32"},{"indexed":false,"name":"closedAt","type":"uint256"}],"name":"Cancel","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionID","type":"bytes32"}],"name":"Reclaim","type":"event"}]`

// EtherAPIsBin is the compiled bytecode used for deploying new contracts.
const EtherAPIsBin = `0x6060604052611702806100126000396000f3606060405236156101325760e060020a60003504630de607c381146101345780630f574ba7146101745780631d7c5cd11461018a5780631ebfdca0146101af5780631f32768e146101c55780633950fc56146102395780633e8b1dd7146102fc57806355404ace146103215780636012042e1461035457806374e29ee61461037b57806378fe2951146104285780638b91124d146104d15780638b91e9a2146104ee5780638ebac11b1461050b57806391499e2d1461055e57806393abc5301461060157806396afb365146106245780639840a6cd14610771578063b214faa51461078e578063c4d252f51461079f578063c95d6edc146107d0578063da2d7b7014610803578063dd8d11e214610827578063e3debbbe1461085b578063ef0e239b14610878578063f60744d514610a19575b005b6101b36004356024355b60408051600160a060020a03939093166c0100000000000000000000000002835260148301919091525190819003603401902090565b610132600435600060006000611200338561013e565b6101b3600435600160a060020a0381166000908152600260205260409020545b919050565b6001545b60408051918252519081900360200190f35b6004356000908152602081815260408051928190206008810154600b820154600983015483546001850154600a9590950154600160a060020a039190911688529587019390935285840191909152606085019190915260ff92909216608084015260a0830191909152519081900360c00190f35b6040805160206004803580820135601f8101849004840285018401909552848452610132949193602493909291840191908190840183828082843750506040805160208835808b0135601f81018390048302840183019094528383529799986044989297509190910194509092508291508401838280828437509496505093359350506064359150506001805480820180835560009291908290828015829011610c5f57600702816007028360005260206000209182019101610c5f9190610d2c565b61013260043560243560443560643560843560a43560006110d9878787878787610a2f565b6101b360043560243560026020526000828152604090208054829081101561000257506000908152602090200154905081565b6101b360043560243560443560643560843560a43560006000610ffd888888888888610a2f565b6101326004358033600160a060020a0316600160005082815481101561000257506000526000805160206116e2833981519152600783020154600160a060020a031614156104245760016001600050838154811015610002575050506000805160206116c28339815191526007820201805461ff00191661010017905560405181907fdfb66150893891bc499d2837280fff700363754123a8d780d6d4e543425a0e8590600090a25b5050565b6101326004358033600160a060020a0316600160005082815481101561000257506000526000805160206116e2833981519152600783020154600160a060020a0316141561042457600160016000508381548110156100025750506000805160206116c28339815191526007830201805460ff1916909117905560405181907fdfb66150893891bc499d2837280fff700363754123a8d780d6d4e543425a0e8590600090a25050565b6101b36004356000818152602081905260409020600b01546101aa565b6101b36004356000818152602081905260409020600801546101aa565b6101b36004356024356044356064355b60408051600160a060020a03959095166c010000000000000000000000000285526014850193909352603484019190915260548301525190819003607401902090565b6101326004358033600160a060020a0316600160005082815481101561000257506000526000805160206116e2833981519152600783020154600160a060020a0316141561042457600060016000508381548110156100025750506000805160206116c28339815191526007830201805460ff1916905560405182917fdfb66150893891bc499d2837280fff700363754123a8d780d6d4e543425a0e8591a25050565b610a6d600435600081815260208190526040902054600160a060020a03166101aa565b6101326004356000818152602081905260409020600b8101544290116104245760405160098201548254600160a060020a0316916000919082818181858883f19350505050506000600050600083815260200190815260200160002060006000820160006101000a815490600160a060020a03021916905560018201600060008201600050600090556001820160006101000a815490600160a060020a03021916905560028201600050805460018160011615610100020316600290046000825580601f106111c457505b5060038201600050805460018160011615610100020316600290046000825580601f106111e257505b5050600060048201819055600582018190556006909101805462ffffff191690556008830181905560098301819055600a8301805460ff19908116909155600b840191909155600c929092018054909216909155505050565b6101b36004356000818152602081905260409020600901546101aa565b61013260043560006111618261082e565b6101326004356000818152602081905260408120548190839033600160a060020a039081169116146111be57610002565b6101b360043560243560036020526000828152604090208054829081101561000257506000908152602090200154905081565b6101b3600435600160a060020a0381166000908152600360205260409020546101aa565b6101b36004355b6000818152602081905260408120600c81015460ff1680156108545750600b8101544290105b9392505050565b6101b36004356000818152602081905260409020600101546101aa565b610a8a600435602060405190810160405280600081526020015060006020604051908101604052806000815260200150600060006000600060006001600050898154811015610002575090819052600789026000805160206116e28339815191528101547fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cfa8201547fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cfb8301546000805160206116c28339815191528401547fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf88501805460408051602060029a841615610100908102600019019094169a909a04601f81018b90048b0282018b019092528181527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf68901999398600160a060020a0398909816977fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf90196959460ff8181169594909104169290918991830182828015610bae5780601f10610b8357610100808354040283529160200191610bae565b6101b360043560243560443560643560843560a4355b6000868152602081905260408120600c81015460ff168015610fef575080546001808301549091610f9991600160a060020a03909116908a8a61051b565b60408051600160a060020a03929092168252519081900360200190f35b604051808060200188600160a060020a031681526020018060200187815260200186815260200185815260200184815260200183810383528a8181518152602001915080519060200190808383829060006004602084601f0104600f02600301f150905090810190601f168015610b155780820380516001836020036101000a031916815260200191505b508381038252888181518152602001915080519060200190808383829060006004602084601f0104600f02600301f150905090810190601f168015610b6e5780820380516001836020036101000a031916815260200191505b50995050505050505050505060405180910390f35b820191906000526020600020905b815481529060010190602001808311610b9157829003601f168201915b5050604080518a54602060026001831615610100026000190190921691909104601f8101829004820283018201909352828252959c50948a94509092508401905082828015610c3e5780601f10610c1357610100808354040283529160200191610c3e565b820191906000526020600020905b815481529060010190602001808311610c2157829003601f168201915b50505050509450975097509750975097509750975050919395979092949650565b505050815481101561000257906000526020600020906007020160005060068101805462ff00001916620100001760ff191690556001805460001990810183558282018054600160a060020a0319163317905587516002848101805460008281526020908190209798509196958116156101000290940190931604601f90810183900484019391928a0190839010610ded57805160ff19168380011785555b50610e1d929150610db7565b5050600060048201819055600582015560068101805462ffffff191690556001015b80821115610dcb57600080825560018281018054600160a060020a031916905560028381018054848255909281161561010002600019011604601f819010610d9d57505b5060038201600050805460018160011615610100020316600290046000825580601f10610dcf5750610d0a565b601f016020900490600052602060002090810190610d7091905b80821115610dcb5760008155600101610db7565b5090565b601f016020900490600052602060002090810190610d0a9190610db7565b82800160010185558215610cfe579182015b82811115610cfe578251826000505591602001919060010190610dff565b505083816003016000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610e7857805160ff19168380011785555b50610ea8929150610db7565b82800160010185558215610e6c579182015b82811115610e6c578251826000505591602001919060010190610e8a565b5050600481018390556005810182905533600160a060020a031660009081526002602052604090208054600181018083558281838015829011610efe57818360005260206000209182019101610efe9190610db7565b50505060009283525060208083208454920191909155604051875133600160a060020a031693899383928582019290918291859183918691600490601f850104600f02600301f15090500191505060405180910390207f5906a2091185df1fc9aec1f6075d226ea7936b2dac0fbd8718beb5e65e2ca57a83600001600050546040518082815260200191505060405180910390a35050505050565b868686604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051808303816000866161da5a03f115610002575050604051518154600160a060020a039081169116145b91505b509695505050505050565b151561100c5760009150610ff2565b506000878152602081905260409020600b8101544290106110305760009150610ff2565b600881015487146110445760009150610ff2565b60019150610ff2565b6040516002820154600160a060020a031690600090879082818181858883f1505050600983018054889003905550505b6040805160088301548152905188917fc19bff313c99700dcf5a7a1351231739052237353454208b6f36ac3a97eeeeb2919081900360200190a26000878152602081905260409020600801805460010190555b50505050505050565b15156110e4576110d0565b50600086815260208190526040902060088101548614611103576110d0565b6002810154600160a060020a03908116339091161461112157610002565b600981015485111561104d5760405160028201546009830154600160a060020a0391909116916000919082818181858883f150505060098301555061107d565b151561116c57610002565b506000818152602081905260408082206009810180543401905590519091839133600160a060020a0316917f678afb2e81183654e6389bac063af1933c7935f97aceeae5aaa51bc54662cf8891a35050565b50505050565b601f0160209004906000526020600020908101906106ef9190610db7565b601f0160209004906000526020600020908101906107189190610db7565b6000818152602081905260409020600c810154919450925060ff1615156111be576001805485908110156100025790600052602060002090600702016000506040805160e0810182523381528151610100818101845284548252600185810154600160a060020a0316602084810191909152600287810180548851948116159095026000190190941604601f810182900482028301820187528083529697509395868501959394889486019391908301828280156112ff5780601f106112d4576101008083540402835291602001916112ff565b820191906000526020600020905b8154815290600101906020018083116112e257829003601f168201915b505050918352505060408051600384018054602060026001831615610100026000190190921691909104601f810182900482028401820190945283835293840193919290918301828280156113955780601f1061136a57610100808354040283529160200191611395565b820191906000526020600020905b81548152906001019060200180831161137857829003601f168201915b50505091835250506040805180820182526004840154815260058401546020828101919091528381019190915260069093015460ff8181168484015261010080830482166060868101919091526201000090930490911660809485015294865260008685018190523487840152908601819052918501829052600160a09590950185905288825281835280822086518154600160a060020a03199081169091178255878501518051838901908155818701516002858101805490951690911790935593810151805160038501805481895297899020959a939996989097938416159096026000190190921692909204601f908101859004840194919392909101908390106114b657805160ff19168380011785555b506114e6929150610db7565b828001600101855582156114aa579182015b828111156114aa5782518260005055916020019190600101906114c8565b50506060820151816003016000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061154557805160ff19168380011785555b50611575929150610db7565b82800160010185558215611539579182015b82811115611539578251826000505591602001919060010190611557565b505060808281015180516004840155602090810151600584015560a084810151600694909401805460c08781015160e0989098015160ff1992831690971761ff001916610100989098029790971762ff000019166201000090960295909517905560408781015160088801556060880151600988015592870151600a8701805486169091179055860151600b8601559490920151600c93909301805490911690921790915533600160a060020a0316600090815260039092529020805460018101808355828183801582901161165e5781836000526020600020918201910161165e9190610db7565b50505060009283525060209182902001849055600883015460408051868152928301919091528051869233600160a060020a0316927fc864b1ad6f1e3cc0c2b4a3a8a0c17e423ba2f01fd79c5591b01ff79edc09fc3992918290030190a35050505056b10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cfcb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf7`

// DeployEtherAPIs deploys a new Ethereum contract, binding an instance of EtherAPIs to it.
func DeployEtherAPIs(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EtherAPIs, error) {
	parsed, err := abi.JSON(strings.NewReader(EtherAPIsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EtherAPIsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EtherAPIs{EtherAPIsCaller: EtherAPIsCaller{contract: contract}, EtherAPIsTransactor: EtherAPIsTransactor{contract: contract}}, nil
}

// EtherAPIs is an auto generated Go binding around an Ethereum contract.
type EtherAPIs struct {
	EtherAPIsCaller     // Read-only binding to the contract
	EtherAPIsTransactor // Write-only binding to the contract
}

// EtherAPIsCaller is an auto generated read-only Go binding around an Ethereum contract.
type EtherAPIsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EtherAPIsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EtherAPIsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EtherAPIsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EtherAPIsSession struct {
	Contract     *EtherAPIs        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EtherAPIsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EtherAPIsCallerSession struct {
	Contract *EtherAPIsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// EtherAPIsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EtherAPIsTransactorSession struct {
	Contract     *EtherAPIsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// NewEtherAPIs creates a new instance of EtherAPIs, bound to a specific deployed contract.
func NewEtherAPIs(address common.Address, backend bind.ContractBackend) (*EtherAPIs, error) {
	contract, err := bindEtherAPIs(address, backend.(bind.ContractCaller), backend.(bind.ContractTransactor))
	if err != nil {
		return nil, err
	}
	return &EtherAPIs{EtherAPIsCaller: EtherAPIsCaller{contract: contract}, EtherAPIsTransactor: EtherAPIsTransactor{contract: contract}}, nil
}

// NewEtherAPIsCaller creates a new read-only instance of EtherAPIs, bound to a specific deployed contract.
func NewEtherAPIsCaller(address common.Address, caller bind.ContractCaller) (*EtherAPIsCaller, error) {
	contract, err := bindEtherAPIs(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &EtherAPIsCaller{contract: contract}, nil
}

// NewEtherAPIsTransactor creates a new write-only instance of EtherAPIs, bound to a specific deployed contract.
func NewEtherAPIsTransactor(address common.Address, transactor bind.ContractTransactor) (*EtherAPIsTransactor, error) {
	contract, err := bindEtherAPIs(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &EtherAPIsTransactor{contract: contract}, nil
}

// bindEtherAPIs binds a generic wrapper to an already deployed contract.
func bindEtherAPIs(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EtherAPIsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// GetHash is a free data retrieval call binding the contract method 0x8ebac11b.
//
// Solidity: function getHash(from address, serviceID uint256, nonce uint256, value uint256) constant returns(bytes32)
func (_EtherAPIs *EtherAPIsCaller) GetHash(opts *bind.CallOpts, from common.Address, serviceID *big.Int, nonce *big.Int, value *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "getHash", from, serviceID, nonce, value)
	return *ret0, err
}

// GetHash is a free data retrieval call binding the contract method 0x8ebac11b.
//
// Solidity: function getHash(from address, serviceID uint256, nonce uint256, value uint256) constant returns(bytes32)
func (_EtherAPIs *EtherAPIsSession) GetHash(from common.Address, serviceID *big.Int, nonce *big.Int, value *big.Int) ([32]byte, error) {
	return _EtherAPIs.Contract.GetHash(&_EtherAPIs.CallOpts, from, serviceID, nonce, value)
}

// GetHash is a free data retrieval call binding the contract method 0x8ebac11b.
//
// Solidity: function getHash(from address, serviceID uint256, nonce uint256, value uint256) constant returns(bytes32)
func (_EtherAPIs *EtherAPIsCallerSession) GetHash(from common.Address, serviceID *big.Int, nonce *big.Int, value *big.Int) ([32]byte, error) {
	return _EtherAPIs.Contract.GetHash(&_EtherAPIs.CallOpts, from, serviceID, nonce, value)
}

// EtherAPIsGetServiceResult is the result of the GetService invocation."
type EtherAPIsGetServiceResult struct {
	Name         string
	Owner        common.Address
	Endpoint     string
	Price        *big.Int
	Cancellation *big.Int
	Enabled      bool
	Deleted      bool
}

// GetService is a free data retrieval call binding the contract method 0xef0e239b.
//
// Solidity: function getService(serviceID uint256) constant returns(name string, owner address, endpoint string, price uint256, cancellation uint256, enabled bool, deleted bool)
func (_EtherAPIs *EtherAPIsCaller) GetService(opts *bind.CallOpts, serviceID *big.Int) (EtherAPIsGetServiceResult, error) {
	var (
		ret = new(EtherAPIsGetServiceResult)
	)
	out := ret
	err := _EtherAPIs.contract.Call(opts, out, "getService", serviceID)
	return *ret, err
}

// GetService is a free data retrieval call binding the contract method 0xef0e239b.
//
// Solidity: function getService(serviceID uint256) constant returns(name string, owner address, endpoint string, price uint256, cancellation uint256, enabled bool, deleted bool)
func (_EtherAPIs *EtherAPIsSession) GetService(serviceID *big.Int) (EtherAPIsGetServiceResult, error) {
	return _EtherAPIs.Contract.GetService(&_EtherAPIs.CallOpts, serviceID)
}

// GetService is a free data retrieval call binding the contract method 0xef0e239b.
//
// Solidity: function getService(serviceID uint256) constant returns(name string, owner address, endpoint string, price uint256, cancellation uint256, enabled bool, deleted bool)
func (_EtherAPIs *EtherAPIsCallerSession) GetService(serviceID *big.Int) (EtherAPIsGetServiceResult, error) {
	return _EtherAPIs.Contract.GetService(&_EtherAPIs.CallOpts, serviceID)
}

// EtherAPIsGetSubscriptionResult is the result of the GetSubscription invocation."
type EtherAPIsGetSubscriptionResult struct {
	From      common.Address
	ServiceID *big.Int
	Nonce     *big.Int
	Value     *big.Int
	Cancelled bool
	ClosedAt  *big.Int
}

// GetSubscription is a free data retrieval call binding the contract method 0x1f32768e.
//
// Solidity: function getSubscription(subscriptionID bytes32) constant returns(from address, serviceID uint256, nonce uint256, value uint256, cancelled bool, closedAt uint256)
func (_EtherAPIs *EtherAPIsCaller) GetSubscription(opts *bind.CallOpts, subscriptionID [32]byte) (EtherAPIsGetSubscriptionResult, error) {
	var (
		ret = new(EtherAPIsGetSubscriptionResult)
	)
	out := ret
	err := _EtherAPIs.contract.Call(opts, out, "getSubscription", subscriptionID)
	return *ret, err
}

// GetSubscription is a free data retrieval call binding the contract method 0x1f32768e.
//
// Solidity: function getSubscription(subscriptionID bytes32) constant returns(from address, serviceID uint256, nonce uint256, value uint256, cancelled bool, closedAt uint256)
func (_EtherAPIs *EtherAPIsSession) GetSubscription(subscriptionID [32]byte) (EtherAPIsGetSubscriptionResult, error) {
	return _EtherAPIs.Contract.GetSubscription(&_EtherAPIs.CallOpts, subscriptionID)
}

// GetSubscription is a free data retrieval call binding the contract method 0x1f32768e.
//
// Solidity: function getSubscription(subscriptionID bytes32) constant returns(from address, serviceID uint256, nonce uint256, value uint256, cancelled bool, closedAt uint256)
func (_EtherAPIs *EtherAPIsCallerSession) GetSubscription(subscriptionID [32]byte) (EtherAPIsGetSubscriptionResult, error) {
	return _EtherAPIs.Contract.GetSubscription(&_EtherAPIs.CallOpts, subscriptionID)
}

// GetSubscriptionClosedAt is a free data retrieval call binding the contract method 0x8b91124d.
//
// Solidity: function getSubscriptionClosedAt(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCaller) GetSubscriptionClosedAt(opts *bind.CallOpts, subscriptionID [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "getSubscriptionClosedAt", subscriptionID)
	return *ret0, err
}

// GetSubscriptionClosedAt is a free data retrieval call binding the contract method 0x8b91124d.
//
// Solidity: function getSubscriptionClosedAt(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsSession) GetSubscriptionClosedAt(subscriptionID [32]byte) (*big.Int, error) {
	return _EtherAPIs.Contract.GetSubscriptionClosedAt(&_EtherAPIs.CallOpts, subscriptionID)
}

// GetSubscriptionClosedAt is a free data retrieval call binding the contract method 0x8b91124d.
//
// Solidity: function getSubscriptionClosedAt(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCallerSession) GetSubscriptionClosedAt(subscriptionID [32]byte) (*big.Int, error) {
	return _EtherAPIs.Contract.GetSubscriptionClosedAt(&_EtherAPIs.CallOpts, subscriptionID)
}

// GetSubscriptionNonce is a free data retrieval call binding the contract method 0x8b91e9a2.
//
// Solidity: function getSubscriptionNonce(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCaller) GetSubscriptionNonce(opts *bind.CallOpts, subscriptionID [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "getSubscriptionNonce", subscriptionID)
	return *ret0, err
}

// GetSubscriptionNonce is a free data retrieval call binding the contract method 0x8b91e9a2.
//
// Solidity: function getSubscriptionNonce(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsSession) GetSubscriptionNonce(subscriptionID [32]byte) (*big.Int, error) {
	return _EtherAPIs.Contract.GetSubscriptionNonce(&_EtherAPIs.CallOpts, subscriptionID)
}

// GetSubscriptionNonce is a free data retrieval call binding the contract method 0x8b91e9a2.
//
// Solidity: function getSubscriptionNonce(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCallerSession) GetSubscriptionNonce(subscriptionID [32]byte) (*big.Int, error) {
	return _EtherAPIs.Contract.GetSubscriptionNonce(&_EtherAPIs.CallOpts, subscriptionID)
}

// GetSubscriptionOwner is a free data retrieval call binding the contract method 0x93abc530.
//
// Solidity: function getSubscriptionOwner(subscriptionID bytes32) constant returns(address)
func (_EtherAPIs *EtherAPIsCaller) GetSubscriptionOwner(opts *bind.CallOpts, subscriptionID [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "getSubscriptionOwner", subscriptionID)
	return *ret0, err
}

// GetSubscriptionOwner is a free data retrieval call binding the contract method 0x93abc530.
//
// Solidity: function getSubscriptionOwner(subscriptionID bytes32) constant returns(address)
func (_EtherAPIs *EtherAPIsSession) GetSubscriptionOwner(subscriptionID [32]byte) (common.Address, error) {
	return _EtherAPIs.Contract.GetSubscriptionOwner(&_EtherAPIs.CallOpts, subscriptionID)
}

// GetSubscriptionOwner is a free data retrieval call binding the contract method 0x93abc530.
//
// Solidity: function getSubscriptionOwner(subscriptionID bytes32) constant returns(address)
func (_EtherAPIs *EtherAPIsCallerSession) GetSubscriptionOwner(subscriptionID [32]byte) (common.Address, error) {
	return _EtherAPIs.Contract.GetSubscriptionOwner(&_EtherAPIs.CallOpts, subscriptionID)
}

// GetSubscriptionServiceID is a free data retrieval call binding the contract method 0xe3debbbe.
//
// Solidity: function getSubscriptionServiceID(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCaller) GetSubscriptionServiceID(opts *bind.CallOpts, subscriptionID [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "getSubscriptionServiceID", subscriptionID)
	return *ret0, err
}

// GetSubscriptionServiceID is a free data retrieval call binding the contract method 0xe3debbbe.
//
// Solidity: function getSubscriptionServiceID(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsSession) GetSubscriptionServiceID(subscriptionID [32]byte) (*big.Int, error) {
	return _EtherAPIs.Contract.GetSubscriptionServiceID(&_EtherAPIs.CallOpts, subscriptionID)
}

// GetSubscriptionServiceID is a free data retrieval call binding the contract method 0xe3debbbe.
//
// Solidity: function getSubscriptionServiceID(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCallerSession) GetSubscriptionServiceID(subscriptionID [32]byte) (*big.Int, error) {
	return _EtherAPIs.Contract.GetSubscriptionServiceID(&_EtherAPIs.CallOpts, subscriptionID)
}

// GetSubscriptionValue is a free data retrieval call binding the contract method 0x9840a6cd.
//
// Solidity: function getSubscriptionValue(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCaller) GetSubscriptionValue(opts *bind.CallOpts, subscriptionID [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "getSubscriptionValue", subscriptionID)
	return *ret0, err
}

// GetSubscriptionValue is a free data retrieval call binding the contract method 0x9840a6cd.
//
// Solidity: function getSubscriptionValue(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsSession) GetSubscriptionValue(subscriptionID [32]byte) (*big.Int, error) {
	return _EtherAPIs.Contract.GetSubscriptionValue(&_EtherAPIs.CallOpts, subscriptionID)
}

// GetSubscriptionValue is a free data retrieval call binding the contract method 0x9840a6cd.
//
// Solidity: function getSubscriptionValue(subscriptionID bytes32) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCallerSession) GetSubscriptionValue(subscriptionID [32]byte) (*big.Int, error) {
	return _EtherAPIs.Contract.GetSubscriptionValue(&_EtherAPIs.CallOpts, subscriptionID)
}

// IsValidSubscription is a free data retrieval call binding the contract method 0xdd8d11e2.
//
// Solidity: function isValidSubscription(subscriptionID bytes32) constant returns(bool)
func (_EtherAPIs *EtherAPIsCaller) IsValidSubscription(opts *bind.CallOpts, subscriptionID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "isValidSubscription", subscriptionID)
	return *ret0, err
}

// IsValidSubscription is a free data retrieval call binding the contract method 0xdd8d11e2.
//
// Solidity: function isValidSubscription(subscriptionID bytes32) constant returns(bool)
func (_EtherAPIs *EtherAPIsSession) IsValidSubscription(subscriptionID [32]byte) (bool, error) {
	return _EtherAPIs.Contract.IsValidSubscription(&_EtherAPIs.CallOpts, subscriptionID)
}

// IsValidSubscription is a free data retrieval call binding the contract method 0xdd8d11e2.
//
// Solidity: function isValidSubscription(subscriptionID bytes32) constant returns(bool)
func (_EtherAPIs *EtherAPIsCallerSession) IsValidSubscription(subscriptionID [32]byte) (bool, error) {
	return _EtherAPIs.Contract.IsValidSubscription(&_EtherAPIs.CallOpts, subscriptionID)
}

// MakeSubscriptionID is a free data retrieval call binding the contract method 0x0de607c3.
//
// Solidity: function makeSubscriptionID(from address, serviceID uint256) constant returns(bytes32)
func (_EtherAPIs *EtherAPIsCaller) MakeSubscriptionID(opts *bind.CallOpts, from common.Address, serviceID *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "makeSubscriptionID", from, serviceID)
	return *ret0, err
}

// MakeSubscriptionID is a free data retrieval call binding the contract method 0x0de607c3.
//
// Solidity: function makeSubscriptionID(from address, serviceID uint256) constant returns(bytes32)
func (_EtherAPIs *EtherAPIsSession) MakeSubscriptionID(from common.Address, serviceID *big.Int) ([32]byte, error) {
	return _EtherAPIs.Contract.MakeSubscriptionID(&_EtherAPIs.CallOpts, from, serviceID)
}

// MakeSubscriptionID is a free data retrieval call binding the contract method 0x0de607c3.
//
// Solidity: function makeSubscriptionID(from address, serviceID uint256) constant returns(bytes32)
func (_EtherAPIs *EtherAPIsCallerSession) MakeSubscriptionID(from common.Address, serviceID *big.Int) ([32]byte, error) {
	return _EtherAPIs.Contract.MakeSubscriptionID(&_EtherAPIs.CallOpts, from, serviceID)
}

// ServicesLength is a free data retrieval call binding the contract method 0x1ebfdca0.
//
// Solidity: function servicesLength() constant returns(uint256)
func (_EtherAPIs *EtherAPIsCaller) ServicesLength(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "servicesLength")
	return *ret0, err
}

// ServicesLength is a free data retrieval call binding the contract method 0x1ebfdca0.
//
// Solidity: function servicesLength() constant returns(uint256)
func (_EtherAPIs *EtherAPIsSession) ServicesLength() (*big.Int, error) {
	return _EtherAPIs.Contract.ServicesLength(&_EtherAPIs.CallOpts)
}

// ServicesLength is a free data retrieval call binding the contract method 0x1ebfdca0.
//
// Solidity: function servicesLength() constant returns(uint256)
func (_EtherAPIs *EtherAPIsCallerSession) ServicesLength() (*big.Int, error) {
	return _EtherAPIs.Contract.ServicesLength(&_EtherAPIs.CallOpts)
}

// UserServices is a free data retrieval call binding the contract method 0x55404ace.
//
// Solidity: function userServices( address,  uint256) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCaller) UserServices(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "userServices", arg0, arg1)
	return *ret0, err
}

// UserServices is a free data retrieval call binding the contract method 0x55404ace.
//
// Solidity: function userServices( address,  uint256) constant returns(uint256)
func (_EtherAPIs *EtherAPIsSession) UserServices(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _EtherAPIs.Contract.UserServices(&_EtherAPIs.CallOpts, arg0, arg1)
}

// UserServices is a free data retrieval call binding the contract method 0x55404ace.
//
// Solidity: function userServices( address,  uint256) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCallerSession) UserServices(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _EtherAPIs.Contract.UserServices(&_EtherAPIs.CallOpts, arg0, arg1)
}

// UserServicesLength is a free data retrieval call binding the contract method 0x1d7c5cd1.
//
// Solidity: function userServicesLength(addr address) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCaller) UserServicesLength(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "userServicesLength", addr)
	return *ret0, err
}

// UserServicesLength is a free data retrieval call binding the contract method 0x1d7c5cd1.
//
// Solidity: function userServicesLength(addr address) constant returns(uint256)
func (_EtherAPIs *EtherAPIsSession) UserServicesLength(addr common.Address) (*big.Int, error) {
	return _EtherAPIs.Contract.UserServicesLength(&_EtherAPIs.CallOpts, addr)
}

// UserServicesLength is a free data retrieval call binding the contract method 0x1d7c5cd1.
//
// Solidity: function userServicesLength(addr address) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCallerSession) UserServicesLength(addr common.Address) (*big.Int, error) {
	return _EtherAPIs.Contract.UserServicesLength(&_EtherAPIs.CallOpts, addr)
}

// UserSubscriptions is a free data retrieval call binding the contract method 0xc95d6edc.
//
// Solidity: function userSubscriptions( address,  uint256) constant returns(bytes32)
func (_EtherAPIs *EtherAPIsCaller) UserSubscriptions(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "userSubscriptions", arg0, arg1)
	return *ret0, err
}

// UserSubscriptions is a free data retrieval call binding the contract method 0xc95d6edc.
//
// Solidity: function userSubscriptions( address,  uint256) constant returns(bytes32)
func (_EtherAPIs *EtherAPIsSession) UserSubscriptions(arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	return _EtherAPIs.Contract.UserSubscriptions(&_EtherAPIs.CallOpts, arg0, arg1)
}

// UserSubscriptions is a free data retrieval call binding the contract method 0xc95d6edc.
//
// Solidity: function userSubscriptions( address,  uint256) constant returns(bytes32)
func (_EtherAPIs *EtherAPIsCallerSession) UserSubscriptions(arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	return _EtherAPIs.Contract.UserSubscriptions(&_EtherAPIs.CallOpts, arg0, arg1)
}

// UserSubscriptionsLength is a free data retrieval call binding the contract method 0xda2d7b70.
//
// Solidity: function userSubscriptionsLength(addr address) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCaller) UserSubscriptionsLength(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "userSubscriptionsLength", addr)
	return *ret0, err
}

// UserSubscriptionsLength is a free data retrieval call binding the contract method 0xda2d7b70.
//
// Solidity: function userSubscriptionsLength(addr address) constant returns(uint256)
func (_EtherAPIs *EtherAPIsSession) UserSubscriptionsLength(addr common.Address) (*big.Int, error) {
	return _EtherAPIs.Contract.UserSubscriptionsLength(&_EtherAPIs.CallOpts, addr)
}

// UserSubscriptionsLength is a free data retrieval call binding the contract method 0xda2d7b70.
//
// Solidity: function userSubscriptionsLength(addr address) constant returns(uint256)
func (_EtherAPIs *EtherAPIsCallerSession) UserSubscriptionsLength(addr common.Address) (*big.Int, error) {
	return _EtherAPIs.Contract.UserSubscriptionsLength(&_EtherAPIs.CallOpts, addr)
}

// VerifyPayment is a free data retrieval call binding the contract method 0x6012042e.
//
// Solidity: function verifyPayment(subscriptionID bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherAPIs *EtherAPIsCaller) VerifyPayment(opts *bind.CallOpts, subscriptionID [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "verifyPayment", subscriptionID, nonce, value, v, r, s)
	return *ret0, err
}

// VerifyPayment is a free data retrieval call binding the contract method 0x6012042e.
//
// Solidity: function verifyPayment(subscriptionID bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherAPIs *EtherAPIsSession) VerifyPayment(subscriptionID [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	return _EtherAPIs.Contract.VerifyPayment(&_EtherAPIs.CallOpts, subscriptionID, nonce, value, v, r, s)
}

// VerifyPayment is a free data retrieval call binding the contract method 0x6012042e.
//
// Solidity: function verifyPayment(subscriptionID bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherAPIs *EtherAPIsCallerSession) VerifyPayment(subscriptionID [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	return _EtherAPIs.Contract.VerifyPayment(&_EtherAPIs.CallOpts, subscriptionID, nonce, value, v, r, s)
}

// VerifySignature is a free data retrieval call binding the contract method 0xf60744d5.
//
// Solidity: function verifySignature(subscriptionID bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherAPIs *EtherAPIsCaller) VerifySignature(opts *bind.CallOpts, subscriptionID [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EtherAPIs.contract.Call(opts, out, "verifySignature", subscriptionID, nonce, value, v, r, s)
	return *ret0, err
}

// VerifySignature is a free data retrieval call binding the contract method 0xf60744d5.
//
// Solidity: function verifySignature(subscriptionID bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherAPIs *EtherAPIsSession) VerifySignature(subscriptionID [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	return _EtherAPIs.Contract.VerifySignature(&_EtherAPIs.CallOpts, subscriptionID, nonce, value, v, r, s)
}

// VerifySignature is a free data retrieval call binding the contract method 0xf60744d5.
//
// Solidity: function verifySignature(subscriptionID bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherAPIs *EtherAPIsCallerSession) VerifySignature(subscriptionID [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	return _EtherAPIs.Contract.VerifySignature(&_EtherAPIs.CallOpts, subscriptionID, nonce, value, v, r, s)
}

// AddService is a paid mutator transaction binding the contract method 0x3950fc56.
//
// Solidity: function addService(name string, endpoint string, price uint256, cancellation uint256) returns()
func (_EtherAPIs *EtherAPIsTransactor) AddService(opts *bind.TransactOpts, name string, endpoint string, price *big.Int, cancellation *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.contract.Transact(opts, "addService", name, endpoint, price, cancellation)
}

// AddService is a paid mutator transaction binding the contract method 0x3950fc56.
//
// Solidity: function addService(name string, endpoint string, price uint256, cancellation uint256) returns()
func (_EtherAPIs *EtherAPIsSession) AddService(name string, endpoint string, price *big.Int, cancellation *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.Contract.AddService(&_EtherAPIs.TransactOpts, name, endpoint, price, cancellation)
}

// AddService is a paid mutator transaction binding the contract method 0x3950fc56.
//
// Solidity: function addService(name string, endpoint string, price uint256, cancellation uint256) returns()
func (_EtherAPIs *EtherAPIsTransactorSession) AddService(name string, endpoint string, price *big.Int, cancellation *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.Contract.AddService(&_EtherAPIs.TransactOpts, name, endpoint, price, cancellation)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(subscriptionID bytes32) returns()
func (_EtherAPIs *EtherAPIsTransactor) Cancel(opts *bind.TransactOpts, subscriptionID [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.contract.Transact(opts, "cancel", subscriptionID)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(subscriptionID bytes32) returns()
func (_EtherAPIs *EtherAPIsSession) Cancel(subscriptionID [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.Contract.Cancel(&_EtherAPIs.TransactOpts, subscriptionID)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(subscriptionID bytes32) returns()
func (_EtherAPIs *EtherAPIsTransactorSession) Cancel(subscriptionID [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.Contract.Cancel(&_EtherAPIs.TransactOpts, subscriptionID)
}

// Claim is a paid mutator transaction binding the contract method 0x3e8b1dd7.
//
// Solidity: function claim(subscriptionID bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) returns()
func (_EtherAPIs *EtherAPIsTransactor) Claim(opts *bind.TransactOpts, subscriptionID [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.contract.Transact(opts, "claim", subscriptionID, nonce, value, v, r, s)
}

// Claim is a paid mutator transaction binding the contract method 0x3e8b1dd7.
//
// Solidity: function claim(subscriptionID bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) returns()
func (_EtherAPIs *EtherAPIsSession) Claim(subscriptionID [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.Contract.Claim(&_EtherAPIs.TransactOpts, subscriptionID, nonce, value, v, r, s)
}

// Claim is a paid mutator transaction binding the contract method 0x3e8b1dd7.
//
// Solidity: function claim(subscriptionID bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) returns()
func (_EtherAPIs *EtherAPIsTransactorSession) Claim(subscriptionID [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.Contract.Claim(&_EtherAPIs.TransactOpts, subscriptionID, nonce, value, v, r, s)
}

// DeleteService is a paid mutator transaction binding the contract method 0x74e29ee6.
//
// Solidity: function deleteService(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsTransactor) DeleteService(opts *bind.TransactOpts, serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.contract.Transact(opts, "deleteService", serviceID)
}

// DeleteService is a paid mutator transaction binding the contract method 0x74e29ee6.
//
// Solidity: function deleteService(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsSession) DeleteService(serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.Contract.DeleteService(&_EtherAPIs.TransactOpts, serviceID)
}

// DeleteService is a paid mutator transaction binding the contract method 0x74e29ee6.
//
// Solidity: function deleteService(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsTransactorSession) DeleteService(serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.Contract.DeleteService(&_EtherAPIs.TransactOpts, serviceID)
}

// Deposit is a paid mutator transaction binding the contract method 0xb214faa5.
//
// Solidity: function deposit(subscriptionID bytes32) returns()
func (_EtherAPIs *EtherAPIsTransactor) Deposit(opts *bind.TransactOpts, subscriptionID [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.contract.Transact(opts, "deposit", subscriptionID)
}

// Deposit is a paid mutator transaction binding the contract method 0xb214faa5.
//
// Solidity: function deposit(subscriptionID bytes32) returns()
func (_EtherAPIs *EtherAPIsSession) Deposit(subscriptionID [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.Contract.Deposit(&_EtherAPIs.TransactOpts, subscriptionID)
}

// Deposit is a paid mutator transaction binding the contract method 0xb214faa5.
//
// Solidity: function deposit(subscriptionID bytes32) returns()
func (_EtherAPIs *EtherAPIsTransactorSession) Deposit(subscriptionID [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.Contract.Deposit(&_EtherAPIs.TransactOpts, subscriptionID)
}

// DisableService is a paid mutator transaction binding the contract method 0x91499e2d.
//
// Solidity: function disableService(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsTransactor) DisableService(opts *bind.TransactOpts, serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.contract.Transact(opts, "disableService", serviceID)
}

// DisableService is a paid mutator transaction binding the contract method 0x91499e2d.
//
// Solidity: function disableService(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsSession) DisableService(serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.Contract.DisableService(&_EtherAPIs.TransactOpts, serviceID)
}

// DisableService is a paid mutator transaction binding the contract method 0x91499e2d.
//
// Solidity: function disableService(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsTransactorSession) DisableService(serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.Contract.DisableService(&_EtherAPIs.TransactOpts, serviceID)
}

// EnableService is a paid mutator transaction binding the contract method 0x78fe2951.
//
// Solidity: function enableService(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsTransactor) EnableService(opts *bind.TransactOpts, serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.contract.Transact(opts, "enableService", serviceID)
}

// EnableService is a paid mutator transaction binding the contract method 0x78fe2951.
//
// Solidity: function enableService(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsSession) EnableService(serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.Contract.EnableService(&_EtherAPIs.TransactOpts, serviceID)
}

// EnableService is a paid mutator transaction binding the contract method 0x78fe2951.
//
// Solidity: function enableService(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsTransactorSession) EnableService(serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.Contract.EnableService(&_EtherAPIs.TransactOpts, serviceID)
}

// Reclaim is a paid mutator transaction binding the contract method 0x96afb365.
//
// Solidity: function reclaim(subscriptionID bytes32) returns()
func (_EtherAPIs *EtherAPIsTransactor) Reclaim(opts *bind.TransactOpts, subscriptionID [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.contract.Transact(opts, "reclaim", subscriptionID)
}

// Reclaim is a paid mutator transaction binding the contract method 0x96afb365.
//
// Solidity: function reclaim(subscriptionID bytes32) returns()
func (_EtherAPIs *EtherAPIsSession) Reclaim(subscriptionID [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.Contract.Reclaim(&_EtherAPIs.TransactOpts, subscriptionID)
}

// Reclaim is a paid mutator transaction binding the contract method 0x96afb365.
//
// Solidity: function reclaim(subscriptionID bytes32) returns()
func (_EtherAPIs *EtherAPIsTransactorSession) Reclaim(subscriptionID [32]byte) (*types.Transaction, error) {
	return _EtherAPIs.Contract.Reclaim(&_EtherAPIs.TransactOpts, subscriptionID)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsTransactor) Subscribe(opts *bind.TransactOpts, serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.contract.Transact(opts, "subscribe", serviceID)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsSession) Subscribe(serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.Contract.Subscribe(&_EtherAPIs.TransactOpts, serviceID)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(serviceID uint256) returns()
func (_EtherAPIs *EtherAPIsTransactorSession) Subscribe(serviceID *big.Int) (*types.Transaction, error) {
	return _EtherAPIs.Contract.Subscribe(&_EtherAPIs.TransactOpts, serviceID)
}
