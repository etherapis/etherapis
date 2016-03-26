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

// EtherApisABI is the input ABI used to generate the binding from.
const EtherApisABI = `[{"constant":false,"inputs":[{"name":"serviceId","type":"uint256"}],"name":"subscribe","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"addr","type":"address"}],"name":"userServicesLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[],"name":"servicesLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscription","outputs":[{"name":"from","type":"address"},{"name":"serviceId","type":"uint256"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"cancelled","type":"bool"},{"name":"closedAt","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"name","type":"string"},{"name":"endpoint","type":"string"},{"name":"price","type":"uint256"},{"name":"cancellationTime","type":"uint256"}],"name":"addService","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"claim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"uint256"}],"name":"userServices","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verifyPayment","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"serviceId","type":"uint256"}],"name":"deleteService","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"serviceId","type":"uint256"}],"name":"enableService","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionClosedAt","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionNonce","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"from","type":"address"},{"name":"serviceId","type":"uint256"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"}],"name":"getHash","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":false,"inputs":[{"name":"serviceId","type":"uint256"}],"name":"disableService","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionOwner","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"reclaim","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionValue","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"from","type":"address"},{"name":"serviceId","type":"uint256"}],"name":"makeSubscriptionId","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"deposit","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"cancel","outputs":[],"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"uint256"}],"name":"userSubscriptions","outputs":[{"name":"","type":"bytes32"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"getSubscriptionServiceId","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"addr","type":"address"}],"name":"userSubscriptionsLength","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"}],"name":"isValidSubscription","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"serviceId","type":"uint256"}],"name":"getService","outputs":[{"name":"name","type":"string"},{"name":"owner","type":"address"},{"name":"endpoint","type":"string"},{"name":"price","type":"uint256"},{"name":"cancellationTime","type":"uint256"},{"name":"enabled","type":"bool"},{"name":"deleted","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"subscriptionId","type":"bytes32"},{"name":"nonce","type":"uint256"},{"name":"value","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}],"name":"verifySignature","outputs":[{"name":"","type":"bool"}],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"name","type":"string"},{"indexed":true,"name":"owner","type":"address"},{"indexed":false,"name":"serviceId","type":"uint256"}],"name":"NewService","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"serviceId","type":"uint256"}],"name":"UpdateService","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"serviceId","type":"uint256"},{"indexed":false,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"NewSubscription","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"subscriptionId","type":"bytes32"}],"name":"Deposit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"Redeem","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"},{"indexed":false,"name":"closedAt","type":"uint256"}],"name":"Cancel","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"subscriptionId","type":"bytes32"}],"name":"Reclaim","type":"event"}]`

// EtherApisBin is the compiled bytecode used for deploying new contracts.
const EtherApisBin = `6060604052612335806100126000396000f360606040523615610150576000357c0100000000000000000000000000000000000000000000000000000000900480630f574ba7146101525780631d7c5cd11461016a5780631ebfdca0146101965780631f32768e146101b95780633950fc561461021e5780633e8b1dd7146102cd57806355404ace146103125780636012042e1461034757806374e29ee6146103a057806378fe2951146103b85780638b91124d146103d05780638b91e9a2146103fc5780638ebac11b1461042857806391499e2d1461046f57806393abc5301461048757806396afb365146104c95780639840a6cd146104e15780639871e4f21461050d578063b214faa514610542578063c4d252f51461055a578063c95d6edc14610572578063d575af74146105a7578063da2d7b70146105d3578063dd8d11e2146105ff578063ef0e239b1461062b578063f60744d51461074757610150565b005b6101686004808035906020019091905050611205565b005b61018060048080359060200190919050506118aa565b6040518082815260200191505060405180910390f35b6101a360048050506118eb565b6040518082815260200191505060405180910390f35b6101cf6004808035906020019091905050610f60565b604051808773ffffffffffffffffffffffffffffffffffffffff168152602001868152602001858152602001848152602001838152602001828152602001965050505050505060405180910390f35b6102cb6004808035906020019082018035906020019191908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050909091908035906020019082018035906020019191908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050909091908035906020019091908035906020019091905050611ddf565b005b61031060048080359060200190919080359060200190919080359060200190919080359060200190919080359060200190919080359060200190919050506109ad565b005b610331600480803590602001909190803590602001909190505061115a565b6040518082815260200191505060405180910390f35b61038a600480803590602001909190803590602001909190803590602001909190803590602001909190803590602001909190803590602001909190505061092b565b6040518082815260200191505060405180910390f35b6103b66004808035906020019091905050611b21565b005b6103ce6004808035906020019091905050611c0b565b005b6103e660048080359060200190919050506110dc565b6040518082815260200191505060405180910390f35b610412600480803590602001909190505061102f565b6040518082815260200191505060405180910390f35b61045960048080359060200190919080359060200190919080359060200190919080359060200190919050506107a0565b6040518082815260200191505060405180910390f35b6104856004808035906020019091905050611cf5565b005b61049d600480803590602001909190505061105d565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6104df6004808035906020019091905050610cdd565b005b6104f76004808035906020019091905050611001565b6040518082815260200191505060405180910390f35b61052c60048080359060200190919080359060200190919050506122e3565b6040518082815260200191505060405180910390f35b6105586004808035906020019091905050610bca565b005b6105706004808035906020019091905050610c5b565b005b610591600480803590602001909190803590602001909190505061118f565b6040518082815260200191505060405180910390f35b6105bd60048080359060200190919050506110a8565b6040518082815260200191505060405180910390f35b6105e960048080359060200190919050506111c4565b6040518082815260200191505060405180910390f35b610615600480803590602001909190505061110a565b6040518082815260200191505060405180910390f35b6106416004808035906020019091905050611900565b60405180806020018873ffffffffffffffffffffffffffffffffffffffff1681526020018060200187815260200186815260200185815260200184815260200183810383528a8181518152602001915080519060200190808383829060006004602084601f0104600f02600301f150905090810190601f1680156106d95780820380516001836020036101000a031916815260200191505b508381038252888181518152602001915080519060200190808383829060006004602084601f0104600f02600301f150905090810190601f1680156107325780820380516001836020036101000a031916815260200191505b50995050505050505050505060405180910390f35b61078a6004808035906020019091908035906020019091908035906020019091908035906020019091908035906020019091908035906020019091905050610804565b6040518082815260200191505060405180910390f35b600084848484604051808573ffffffffffffffffffffffffffffffffffffffff166c01000000000000000000000000028152601401848152602001838152602001828152602001945050505050604051809103902090506107fc565b949350505050565b6000600060006000506000898152602001908152602001600020600050905080600c0160009054906101000a900460ff1680156109195750600161087b8260000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1683600101600050600001600050548a8a6107a0565b868686604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051808303816000866161da5a03f1156100025750506040518051906020015073ffffffffffffffffffffffffffffffffffffffff168160000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b9150610920565b509695505050505050565b6000600061093d888888888888610804565b151561094c57600091506109a2565b6000600050600089815260200190815260200160002060005090504281600b016000505410151561098057600091506109a2565b86816008016000505414151561099957600091506109a2565b600191506109a2565b509695505050505050565b60006109bd878787878787610804565b15156109c857610bc1565b6000600050600088815260200190815260200160002060005090508581600801600050541415156109f857610bc1565b3373ffffffffffffffffffffffffffffffffffffffff168160010160005060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141515610a5c57610002565b8060090160005054851115610ae1578060010160005060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660008260090160005054604051809050600060405180830381858888f193505050505060008160090160005081905550610b53565b8060010160005060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600086604051809050600060405180830381858888f193505050505084816009016000828282505403925050819055505b867fc19bff313c99700dcf5a7a1351231739052237353454208b6f36ac3a97eeeeb282600801600050546040518082815260200191505060405180910390a2600060005060008881526020019081526020016000206000506008016000818150548092919060010191905055505b50505050505050565b6000610bd58261110a565b1515610be057610002565b6000600050600083815260200190815260200160002060005090503481600901600082828250540192505081905550813373ffffffffffffffffffffffffffffffffffffffff167f678afb2e81183654e6389bac063af1933c7935f97aceeae5aaa51bc54662cf8860405180905060405180910390a35b5050565b60006000823373ffffffffffffffffffffffffffffffffffffffff166000600050600083815260200190815260200160002060005060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141515610cd657610002565b505b505050565b60006000600050600083815260200190815260200160002060005090504281600b0160005054111515610f5b578060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660008260090160005054604051809050600060405180830381858888f19350505050506000600050600083815260200190815260200160002060006000820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905560018201600060008201600050600090556001820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905560028201600050805460018160011615610100020316600290046000825580601f10610e0b5750610e48565b601f016020900490600052602060002090810190610e479190610e29565b80821115610e435760008181506000905550600101610e29565b5090565b5b5060038201600050805460018160011615610100020316600290046000825580601f10610e755750610eb2565b601f016020900490600052602060002090810190610eb19190610e93565b80821115610ead5760008181506000905550600101610e93565b5090565b5b506004820160006000820160005060009055600182016000506000905550506006820160006101000a81549060ff02191690556006820160016101000a81549060ff02191690556006820160026101000a81549060ff0219169055505060088201600050600090556009820160005060009055600a820160006101000a81549060ff0219169055600b820160005060009055600c820160006101000a81549060ff021916905550505b5b5050565b60006000600060006000600060006000600050600089815260200190815260200160002060005090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681600101600050600001600050548260080160005054836009016000505484600a0160009054906101000a900460ff1685600b0160005054965096509650965096509650610ff7565b5091939550919395565b60006000600050600083815260200190815260200160002060005060090160005054905061102a565b919050565b600060006000506000838152602001908152602001600020600050600801600050549050611058565b919050565b60006000600050600083815260200190815260200160002060005060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506110a3565b919050565b6000600060005060008381526020019081526020016000206000506001016000506000016000505490506110d7565b919050565b600060006000506000838152602001908152602001600020600050600b01600050549050611105565b919050565b6000600060006000506000848152602001908152602001600020600050905080600c0160009054906101000a900460ff16801561114d57504281600b0160005054105b9150611154565b50919050565b600260005060205281600052604060002060005081815481101561000257906000526020600020900160005b91509150505481565b600360005060205281600052604060002060005081815481101561000257906000526020600020900160005b91509150505481565b6000600360005060008373ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050805490509050611200565b919050565b60006000600061121533856122e3565b925060006000506000848152602001908152602001600020600050915081600c0160009054906101000a900460ff1615156118a357600160005084815481101561000257906000526020600020906007020160005b50905060e060405190810160405280338152602001826101006040519081016040529081600082016000505481526020016001820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600282016000508054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156113765780601f1061134b57610100808354040283529160200191611376565b820191906000526020600020905b81548152906001019060200180831161135957829003601f168201915b50505050508152602001600382016000508054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561141b5780601f106113f05761010080835404028352916020019161141b565b820191906000526020600020905b8154815290600101906020018083116113fe57829003601f168201915b5050505050815260200160048201600050604060405190810160405290816000820160005054815260200160018201600050548152602001505081526020016006820160009054906101000a900460ff1681526020016006820160019054906101000a900460ff1681526020016006820160029054906101000a900460ff1681526020015050815260200160008152602001348152602001600081526020016000815260200160018152602001506000600050600085815260200190815260200160002060005060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908302179055506020820151816001016000506000820151816000016000505560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908302179055506040820151816002016000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106115b257805160ff19168380011785556115e3565b828001600101855582156115e3579182015b828111156115e25782518260005055916020019190600101906115c4565b5b50905061160e91906115f0565b8082111561160a57600081815060009055506001016115f0565b5090565b50506060820151816003016000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061166557805160ff1916838001178555611696565b82800160010185558215611696579182015b82811115611695578251826000505591602001919060010190611677565b5b5090506116c191906116a3565b808211156116bd57600081815060009055506001016116a3565b5090565b50506080820151816004016000506000820151816000016000505560208201518160010160005055505060a08201518160060160006101000a81548160ff0219169083021790555060c08201518160060160016101000a81548160ff0219169083021790555060e08201518160060160026101000a81548160ff0219169083021790555050506040820151816008016000505560608201518160090160005055608082015181600a0160006101000a81548160ff0219169083021790555060a082015181600b016000505560c082015181600c0160006101000a81548160ff02191690830217905550905050600360005060003373ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050805480600101828181548183558181151161182957818360005260206000209182019101611828919061180a565b80821115611824576000818150600090555060010161180a565b5090565b5b5050509190906000526020600020900160005b8590919091505550833373ffffffffffffffffffffffffffffffffffffffff167fc864b1ad6f1e3cc0c2b4a3a8a0c17e423ba2f01fd79c5591b01ff79edc09fc39858560080160005054604051808381526020018281526020019250505060405180910390a35b5b50505050565b6000600260005060008373ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000508054905090506118e6565b919050565b600060016000508054905090506118fd565b90565b60206040519081016040528060008152602001506000602060405190810160405280600081526020015060006000600060006000600160005089815481101561000257906000526020600020906007020160005b509050806002016000508160010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1682600301600050836004016000506000016000505484600401600050600101600050548560060160009054906101000a900460ff168660060160019054906101000a900460ff16868054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015611a5f5780601f10611a3457610100808354040283529160200191611a5f565b820191906000526020600020905b815481529060010190602001808311611a4257829003601f168201915b50505050509650848054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015611afb5780601f10611ad057610100808354040283529160200191611afb565b820191906000526020600020905b815481529060010190602001808311611ade57829003601f168201915b505050505094509750975097509750975097509750611b15565b50919395979092949650565b803373ffffffffffffffffffffffffffffffffffffffff16600160005082815481101561000257906000526020600020906007020160005b5060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415611c06576001600160005083815481101561000257906000526020600020906007020160005b5060060160016101000a81548160ff02191690830217905550817fdfb66150893891bc499d2837280fff700363754123a8d780d6d4e543425a0e8560405180905060405180910390a25b505b50565b803373ffffffffffffffffffffffffffffffffffffffff16600160005082815481101561000257906000526020600020906007020160005b5060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415611cf0576001600160005083815481101561000257906000526020600020906007020160005b5060060160006101000a81548160ff02191690830217905550817fdfb66150893891bc499d2837280fff700363754123a8d780d6d4e543425a0e8560405180905060405180910390a25b505b50565b803373ffffffffffffffffffffffffffffffffffffffff16600160005082815481101561000257906000526020600020906007020160005b5060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415611dda576000600160005083815481101561000257906000526020600020906007020160005b5060060160006101000a81548160ff02191690830217905550817fdfb66150893891bc499d2837280fff700363754123a8d780d6d4e543425a0e8560405180905060405180910390a25b505b50565b60006001600050600160005080548091906001019090815481835581811511611f9957600702816007028360005260206000209182019101611f989190611e21565b80821115611f9457600060008201600050600090556001820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905560028201600050805460018160011615610100020316600290046000825580601f10611e895750611ec6565b601f016020900490600052602060002090810190611ec59190611ea7565b80821115611ec15760008181506000905550600101611ea7565b5090565b5b5060038201600050805460018160011615610100020316600290046000825580601f10611ef35750611f30565b601f016020900490600052602060002090810190611f2f9190611f11565b80821115611f2b5760008181506000905550600101611f11565b5090565b5b506004820160006000820160005060009055600182016000506000905550506006820160006101000a81549060ff02191690556006820160016101000a81549060ff02191690556006820160026101000a81549060ff021916905550600101611e21565b5090565b5b505050815481101561000257906000526020600020906007020160005b50905060018160060160026101000a81548160ff0219169083021790555060008160060160006101000a81548160ff021916908302179055506001600160005080549050038160000160005081905550338160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff0219169083021790555084816002016000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061208457805160ff19168380011785556120b5565b828001600101855582156120b5579182015b828111156120b4578251826000505591602001919060010190612096565b5b5090506120e091906120c2565b808211156120dc57600081815060009055506001016120c2565b5090565b505083816003016000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061213357805160ff1916838001178555612164565b82800160010185558215612164579182015b82811115612163578251826000505591602001919060010190612145565b5b50905061218f9190612171565b8082111561218b5760008181506000905550600101612171565b5090565b5050828160040160005060000160005081905550818160040160005060010160005081905550600260005060003373ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000508054806001018281815481835581811511612231578183600052602060002091820191016122309190612212565b8082111561222c5760008181506000905550600101612212565b5090565b5b5050509190906000526020600020900160005b8360000160005054909190915055503373ffffffffffffffffffffffffffffffffffffffff1685604051808280519060200190808383829060006004602084601f0104600f02600301f15090500191505060405180910390207f5906a2091185df1fc9aec1f6075d226ea7936b2dac0fbd8718beb5e65e2ca57a83600001600050546040518082815260200191505060405180910390a35b5050505050565b60008282604051808373ffffffffffffffffffffffffffffffffffffffff166c01000000000000000000000000028152601401828152602001925050506040518091039020905061232f565b9291505056`

// DeployEtherApis deploys a new Ethereum contract, binding an instance of EtherApis to it.
func DeployEtherApis(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EtherApis, error) {
	parsed, err := abi.JSON(strings.NewReader(EtherApisABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EtherApisBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EtherApis{EtherApisCaller: EtherApisCaller{contract: contract}, EtherApisTransactor: EtherApisTransactor{contract: contract}}, nil
}

// EtherApis is an auto generated Go binding around an Ethereum contract.
type EtherApis struct {
	EtherApisCaller     // Read-only binding to the contract
	EtherApisTransactor // Write-only binding to the contract
}

// EtherApisCaller is an auto generated read-only Go binding around an Ethereum contract.
type EtherApisCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EtherApisTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EtherApisTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EtherApisSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EtherApisSession struct {
	Contract     *EtherApis        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EtherApisCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EtherApisCallerSession struct {
	Contract *EtherApisCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// EtherApisTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EtherApisTransactorSession struct {
	Contract     *EtherApisTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// NewEtherApis creates a new instance of EtherApis, bound to a specific deployed contract.
func NewEtherApis(address common.Address, backend bind.ContractBackend) (*EtherApis, error) {
	contract, err := bindEtherApis(address, backend.(bind.ContractCaller), backend.(bind.ContractTransactor))
	if err != nil {
		return nil, err
	}
	return &EtherApis{EtherApisCaller: EtherApisCaller{contract: contract}, EtherApisTransactor: EtherApisTransactor{contract: contract}}, nil
}

// NewEtherApisCaller creates a new read-only instance of EtherApis, bound to a specific deployed contract.
func NewEtherApisCaller(address common.Address, caller bind.ContractCaller) (*EtherApisCaller, error) {
	contract, err := bindEtherApis(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &EtherApisCaller{contract: contract}, nil
}

// NewEtherApisTransactor creates a new write-only instance of EtherApis, bound to a specific deployed contract.
func NewEtherApisTransactor(address common.Address, transactor bind.ContractTransactor) (*EtherApisTransactor, error) {
	contract, err := bindEtherApis(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &EtherApisTransactor{contract: contract}, nil
}

// bindEtherApis binds a generic wrapper to an already deployed contract.
func bindEtherApis(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EtherApisABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// GetHash is a free data retrieval call binding the contract method 0x8ebac11b.
//
// Solidity: function getHash(from address, serviceId uint256, nonce uint256, value uint256) constant returns(bytes32)
func (_EtherApis *EtherApisCaller) GetHash(opts *bind.CallOpts, from common.Address, serviceId *big.Int, nonce *big.Int, value *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "getHash", from, serviceId, nonce, value)
	return *ret0, err
}

// GetHash is a free data retrieval call binding the contract method 0x8ebac11b.
//
// Solidity: function getHash(from address, serviceId uint256, nonce uint256, value uint256) constant returns(bytes32)
func (_EtherApis *EtherApisSession) GetHash(from common.Address, serviceId *big.Int, nonce *big.Int, value *big.Int) ([32]byte, error) {
	return _EtherApis.Contract.GetHash(&_EtherApis.CallOpts, from, serviceId, nonce, value)
}

// GetHash is a free data retrieval call binding the contract method 0x8ebac11b.
//
// Solidity: function getHash(from address, serviceId uint256, nonce uint256, value uint256) constant returns(bytes32)
func (_EtherApis *EtherApisCallerSession) GetHash(from common.Address, serviceId *big.Int, nonce *big.Int, value *big.Int) ([32]byte, error) {
	return _EtherApis.Contract.GetHash(&_EtherApis.CallOpts, from, serviceId, nonce, value)
}

// GetServiceResult is the result of the GetService invocation."
type GetServiceResult struct {
	Name             string
	Owner            common.Address
	Endpoint         string
	Price            *big.Int
	CancellationTime *big.Int
	Enabled          bool
	Deleted          bool
}

// GetService is a free data retrieval call binding the contract method 0xef0e239b.
//
// Solidity: function getService(serviceId uint256) constant returns(name string, owner address, endpoint string, price uint256, cancellationTime uint256, enabled bool, deleted bool)
func (_EtherApis *EtherApisCaller) GetService(opts *bind.CallOpts, serviceId *big.Int) (GetServiceResult, error) {
	var (
		ret = new(GetServiceResult)
	)
	out := ret
	err := _EtherApis.contract.Call(opts, out, "getService", serviceId)
	return *ret, err
}

// GetService is a free data retrieval call binding the contract method 0xef0e239b.
//
// Solidity: function getService(serviceId uint256) constant returns(name string, owner address, endpoint string, price uint256, cancellationTime uint256, enabled bool, deleted bool)
func (_EtherApis *EtherApisSession) GetService(serviceId *big.Int) (GetServiceResult, error) {
	return _EtherApis.Contract.GetService(&_EtherApis.CallOpts, serviceId)
}

// GetService is a free data retrieval call binding the contract method 0xef0e239b.
//
// Solidity: function getService(serviceId uint256) constant returns(name string, owner address, endpoint string, price uint256, cancellationTime uint256, enabled bool, deleted bool)
func (_EtherApis *EtherApisCallerSession) GetService(serviceId *big.Int) (GetServiceResult, error) {
	return _EtherApis.Contract.GetService(&_EtherApis.CallOpts, serviceId)
}

// GetSubscriptionResult is the result of the GetSubscription invocation."
type GetSubscriptionResult struct {
	From      common.Address
	ServiceId *big.Int
	Nonce     *big.Int
	Value     *big.Int
	Cancelled bool
	ClosedAt  *big.Int
}

// GetSubscription is a free data retrieval call binding the contract method 0x1f32768e.
//
// Solidity: function getSubscription(subscriptionId bytes32) constant returns(from address, serviceId uint256, nonce uint256, value uint256, cancelled bool, closedAt uint256)
func (_EtherApis *EtherApisCaller) GetSubscription(opts *bind.CallOpts, subscriptionId [32]byte) (GetSubscriptionResult, error) {
	var (
		ret = new(GetSubscriptionResult)
	)
	out := ret
	err := _EtherApis.contract.Call(opts, out, "getSubscription", subscriptionId)
	return *ret, err
}

// GetSubscription is a free data retrieval call binding the contract method 0x1f32768e.
//
// Solidity: function getSubscription(subscriptionId bytes32) constant returns(from address, serviceId uint256, nonce uint256, value uint256, cancelled bool, closedAt uint256)
func (_EtherApis *EtherApisSession) GetSubscription(subscriptionId [32]byte) (GetSubscriptionResult, error) {
	return _EtherApis.Contract.GetSubscription(&_EtherApis.CallOpts, subscriptionId)
}

// GetSubscription is a free data retrieval call binding the contract method 0x1f32768e.
//
// Solidity: function getSubscription(subscriptionId bytes32) constant returns(from address, serviceId uint256, nonce uint256, value uint256, cancelled bool, closedAt uint256)
func (_EtherApis *EtherApisCallerSession) GetSubscription(subscriptionId [32]byte) (GetSubscriptionResult, error) {
	return _EtherApis.Contract.GetSubscription(&_EtherApis.CallOpts, subscriptionId)
}

// GetSubscriptionClosedAt is a free data retrieval call binding the contract method 0x8b91124d.
//
// Solidity: function getSubscriptionClosedAt(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisCaller) GetSubscriptionClosedAt(opts *bind.CallOpts, subscriptionId [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "getSubscriptionClosedAt", subscriptionId)
	return *ret0, err
}

// GetSubscriptionClosedAt is a free data retrieval call binding the contract method 0x8b91124d.
//
// Solidity: function getSubscriptionClosedAt(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisSession) GetSubscriptionClosedAt(subscriptionId [32]byte) (*big.Int, error) {
	return _EtherApis.Contract.GetSubscriptionClosedAt(&_EtherApis.CallOpts, subscriptionId)
}

// GetSubscriptionClosedAt is a free data retrieval call binding the contract method 0x8b91124d.
//
// Solidity: function getSubscriptionClosedAt(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisCallerSession) GetSubscriptionClosedAt(subscriptionId [32]byte) (*big.Int, error) {
	return _EtherApis.Contract.GetSubscriptionClosedAt(&_EtherApis.CallOpts, subscriptionId)
}

// GetSubscriptionNonce is a free data retrieval call binding the contract method 0x8b91e9a2.
//
// Solidity: function getSubscriptionNonce(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisCaller) GetSubscriptionNonce(opts *bind.CallOpts, subscriptionId [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "getSubscriptionNonce", subscriptionId)
	return *ret0, err
}

// GetSubscriptionNonce is a free data retrieval call binding the contract method 0x8b91e9a2.
//
// Solidity: function getSubscriptionNonce(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisSession) GetSubscriptionNonce(subscriptionId [32]byte) (*big.Int, error) {
	return _EtherApis.Contract.GetSubscriptionNonce(&_EtherApis.CallOpts, subscriptionId)
}

// GetSubscriptionNonce is a free data retrieval call binding the contract method 0x8b91e9a2.
//
// Solidity: function getSubscriptionNonce(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisCallerSession) GetSubscriptionNonce(subscriptionId [32]byte) (*big.Int, error) {
	return _EtherApis.Contract.GetSubscriptionNonce(&_EtherApis.CallOpts, subscriptionId)
}

// GetSubscriptionOwner is a free data retrieval call binding the contract method 0x93abc530.
//
// Solidity: function getSubscriptionOwner(subscriptionId bytes32) constant returns(address)
func (_EtherApis *EtherApisCaller) GetSubscriptionOwner(opts *bind.CallOpts, subscriptionId [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "getSubscriptionOwner", subscriptionId)
	return *ret0, err
}

// GetSubscriptionOwner is a free data retrieval call binding the contract method 0x93abc530.
//
// Solidity: function getSubscriptionOwner(subscriptionId bytes32) constant returns(address)
func (_EtherApis *EtherApisSession) GetSubscriptionOwner(subscriptionId [32]byte) (common.Address, error) {
	return _EtherApis.Contract.GetSubscriptionOwner(&_EtherApis.CallOpts, subscriptionId)
}

// GetSubscriptionOwner is a free data retrieval call binding the contract method 0x93abc530.
//
// Solidity: function getSubscriptionOwner(subscriptionId bytes32) constant returns(address)
func (_EtherApis *EtherApisCallerSession) GetSubscriptionOwner(subscriptionId [32]byte) (common.Address, error) {
	return _EtherApis.Contract.GetSubscriptionOwner(&_EtherApis.CallOpts, subscriptionId)
}

// GetSubscriptionServiceId is a free data retrieval call binding the contract method 0xd575af74.
//
// Solidity: function getSubscriptionServiceId(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisCaller) GetSubscriptionServiceId(opts *bind.CallOpts, subscriptionId [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "getSubscriptionServiceId", subscriptionId)
	return *ret0, err
}

// GetSubscriptionServiceId is a free data retrieval call binding the contract method 0xd575af74.
//
// Solidity: function getSubscriptionServiceId(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisSession) GetSubscriptionServiceId(subscriptionId [32]byte) (*big.Int, error) {
	return _EtherApis.Contract.GetSubscriptionServiceId(&_EtherApis.CallOpts, subscriptionId)
}

// GetSubscriptionServiceId is a free data retrieval call binding the contract method 0xd575af74.
//
// Solidity: function getSubscriptionServiceId(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisCallerSession) GetSubscriptionServiceId(subscriptionId [32]byte) (*big.Int, error) {
	return _EtherApis.Contract.GetSubscriptionServiceId(&_EtherApis.CallOpts, subscriptionId)
}

// GetSubscriptionValue is a free data retrieval call binding the contract method 0x9840a6cd.
//
// Solidity: function getSubscriptionValue(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisCaller) GetSubscriptionValue(opts *bind.CallOpts, subscriptionId [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "getSubscriptionValue", subscriptionId)
	return *ret0, err
}

// GetSubscriptionValue is a free data retrieval call binding the contract method 0x9840a6cd.
//
// Solidity: function getSubscriptionValue(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisSession) GetSubscriptionValue(subscriptionId [32]byte) (*big.Int, error) {
	return _EtherApis.Contract.GetSubscriptionValue(&_EtherApis.CallOpts, subscriptionId)
}

// GetSubscriptionValue is a free data retrieval call binding the contract method 0x9840a6cd.
//
// Solidity: function getSubscriptionValue(subscriptionId bytes32) constant returns(uint256)
func (_EtherApis *EtherApisCallerSession) GetSubscriptionValue(subscriptionId [32]byte) (*big.Int, error) {
	return _EtherApis.Contract.GetSubscriptionValue(&_EtherApis.CallOpts, subscriptionId)
}

// IsValidSubscription is a free data retrieval call binding the contract method 0xdd8d11e2.
//
// Solidity: function isValidSubscription(subscriptionId bytes32) constant returns(bool)
func (_EtherApis *EtherApisCaller) IsValidSubscription(opts *bind.CallOpts, subscriptionId [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "isValidSubscription", subscriptionId)
	return *ret0, err
}

// IsValidSubscription is a free data retrieval call binding the contract method 0xdd8d11e2.
//
// Solidity: function isValidSubscription(subscriptionId bytes32) constant returns(bool)
func (_EtherApis *EtherApisSession) IsValidSubscription(subscriptionId [32]byte) (bool, error) {
	return _EtherApis.Contract.IsValidSubscription(&_EtherApis.CallOpts, subscriptionId)
}

// IsValidSubscription is a free data retrieval call binding the contract method 0xdd8d11e2.
//
// Solidity: function isValidSubscription(subscriptionId bytes32) constant returns(bool)
func (_EtherApis *EtherApisCallerSession) IsValidSubscription(subscriptionId [32]byte) (bool, error) {
	return _EtherApis.Contract.IsValidSubscription(&_EtherApis.CallOpts, subscriptionId)
}

// MakeSubscriptionId is a free data retrieval call binding the contract method 0x9871e4f2.
//
// Solidity: function makeSubscriptionId(from address, serviceId uint256) constant returns(bytes32)
func (_EtherApis *EtherApisCaller) MakeSubscriptionId(opts *bind.CallOpts, from common.Address, serviceId *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "makeSubscriptionId", from, serviceId)
	return *ret0, err
}

// MakeSubscriptionId is a free data retrieval call binding the contract method 0x9871e4f2.
//
// Solidity: function makeSubscriptionId(from address, serviceId uint256) constant returns(bytes32)
func (_EtherApis *EtherApisSession) MakeSubscriptionId(from common.Address, serviceId *big.Int) ([32]byte, error) {
	return _EtherApis.Contract.MakeSubscriptionId(&_EtherApis.CallOpts, from, serviceId)
}

// MakeSubscriptionId is a free data retrieval call binding the contract method 0x9871e4f2.
//
// Solidity: function makeSubscriptionId(from address, serviceId uint256) constant returns(bytes32)
func (_EtherApis *EtherApisCallerSession) MakeSubscriptionId(from common.Address, serviceId *big.Int) ([32]byte, error) {
	return _EtherApis.Contract.MakeSubscriptionId(&_EtherApis.CallOpts, from, serviceId)
}

// ServicesLength is a free data retrieval call binding the contract method 0x1ebfdca0.
//
// Solidity: function servicesLength() constant returns(uint256)
func (_EtherApis *EtherApisCaller) ServicesLength(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "servicesLength")
	return *ret0, err
}

// ServicesLength is a free data retrieval call binding the contract method 0x1ebfdca0.
//
// Solidity: function servicesLength() constant returns(uint256)
func (_EtherApis *EtherApisSession) ServicesLength() (*big.Int, error) {
	return _EtherApis.Contract.ServicesLength(&_EtherApis.CallOpts)
}

// ServicesLength is a free data retrieval call binding the contract method 0x1ebfdca0.
//
// Solidity: function servicesLength() constant returns(uint256)
func (_EtherApis *EtherApisCallerSession) ServicesLength() (*big.Int, error) {
	return _EtherApis.Contract.ServicesLength(&_EtherApis.CallOpts)
}

// UserServices is a free data retrieval call binding the contract method 0x55404ace.
//
// Solidity: function userServices( address,  uint256) constant returns(uint256)
func (_EtherApis *EtherApisCaller) UserServices(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "userServices", arg0, arg1)
	return *ret0, err
}

// UserServices is a free data retrieval call binding the contract method 0x55404ace.
//
// Solidity: function userServices( address,  uint256) constant returns(uint256)
func (_EtherApis *EtherApisSession) UserServices(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _EtherApis.Contract.UserServices(&_EtherApis.CallOpts, arg0, arg1)
}

// UserServices is a free data retrieval call binding the contract method 0x55404ace.
//
// Solidity: function userServices( address,  uint256) constant returns(uint256)
func (_EtherApis *EtherApisCallerSession) UserServices(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _EtherApis.Contract.UserServices(&_EtherApis.CallOpts, arg0, arg1)
}

// UserServicesLength is a free data retrieval call binding the contract method 0x1d7c5cd1.
//
// Solidity: function userServicesLength(addr address) constant returns(uint256)
func (_EtherApis *EtherApisCaller) UserServicesLength(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "userServicesLength", addr)
	return *ret0, err
}

// UserServicesLength is a free data retrieval call binding the contract method 0x1d7c5cd1.
//
// Solidity: function userServicesLength(addr address) constant returns(uint256)
func (_EtherApis *EtherApisSession) UserServicesLength(addr common.Address) (*big.Int, error) {
	return _EtherApis.Contract.UserServicesLength(&_EtherApis.CallOpts, addr)
}

// UserServicesLength is a free data retrieval call binding the contract method 0x1d7c5cd1.
//
// Solidity: function userServicesLength(addr address) constant returns(uint256)
func (_EtherApis *EtherApisCallerSession) UserServicesLength(addr common.Address) (*big.Int, error) {
	return _EtherApis.Contract.UserServicesLength(&_EtherApis.CallOpts, addr)
}

// UserSubscriptions is a free data retrieval call binding the contract method 0xc95d6edc.
//
// Solidity: function userSubscriptions( address,  uint256) constant returns(bytes32)
func (_EtherApis *EtherApisCaller) UserSubscriptions(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "userSubscriptions", arg0, arg1)
	return *ret0, err
}

// UserSubscriptions is a free data retrieval call binding the contract method 0xc95d6edc.
//
// Solidity: function userSubscriptions( address,  uint256) constant returns(bytes32)
func (_EtherApis *EtherApisSession) UserSubscriptions(arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	return _EtherApis.Contract.UserSubscriptions(&_EtherApis.CallOpts, arg0, arg1)
}

// UserSubscriptions is a free data retrieval call binding the contract method 0xc95d6edc.
//
// Solidity: function userSubscriptions( address,  uint256) constant returns(bytes32)
func (_EtherApis *EtherApisCallerSession) UserSubscriptions(arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	return _EtherApis.Contract.UserSubscriptions(&_EtherApis.CallOpts, arg0, arg1)
}

// UserSubscriptionsLength is a free data retrieval call binding the contract method 0xda2d7b70.
//
// Solidity: function userSubscriptionsLength(addr address) constant returns(uint256)
func (_EtherApis *EtherApisCaller) UserSubscriptionsLength(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "userSubscriptionsLength", addr)
	return *ret0, err
}

// UserSubscriptionsLength is a free data retrieval call binding the contract method 0xda2d7b70.
//
// Solidity: function userSubscriptionsLength(addr address) constant returns(uint256)
func (_EtherApis *EtherApisSession) UserSubscriptionsLength(addr common.Address) (*big.Int, error) {
	return _EtherApis.Contract.UserSubscriptionsLength(&_EtherApis.CallOpts, addr)
}

// UserSubscriptionsLength is a free data retrieval call binding the contract method 0xda2d7b70.
//
// Solidity: function userSubscriptionsLength(addr address) constant returns(uint256)
func (_EtherApis *EtherApisCallerSession) UserSubscriptionsLength(addr common.Address) (*big.Int, error) {
	return _EtherApis.Contract.UserSubscriptionsLength(&_EtherApis.CallOpts, addr)
}

// VerifyPayment is a free data retrieval call binding the contract method 0x6012042e.
//
// Solidity: function verifyPayment(subscriptionId bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherApis *EtherApisCaller) VerifyPayment(opts *bind.CallOpts, subscriptionId [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "verifyPayment", subscriptionId, nonce, value, v, r, s)
	return *ret0, err
}

// VerifyPayment is a free data retrieval call binding the contract method 0x6012042e.
//
// Solidity: function verifyPayment(subscriptionId bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherApis *EtherApisSession) VerifyPayment(subscriptionId [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	return _EtherApis.Contract.VerifyPayment(&_EtherApis.CallOpts, subscriptionId, nonce, value, v, r, s)
}

// VerifyPayment is a free data retrieval call binding the contract method 0x6012042e.
//
// Solidity: function verifyPayment(subscriptionId bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherApis *EtherApisCallerSession) VerifyPayment(subscriptionId [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	return _EtherApis.Contract.VerifyPayment(&_EtherApis.CallOpts, subscriptionId, nonce, value, v, r, s)
}

// VerifySignature is a free data retrieval call binding the contract method 0xf60744d5.
//
// Solidity: function verifySignature(subscriptionId bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherApis *EtherApisCaller) VerifySignature(opts *bind.CallOpts, subscriptionId [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EtherApis.contract.Call(opts, out, "verifySignature", subscriptionId, nonce, value, v, r, s)
	return *ret0, err
}

// VerifySignature is a free data retrieval call binding the contract method 0xf60744d5.
//
// Solidity: function verifySignature(subscriptionId bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherApis *EtherApisSession) VerifySignature(subscriptionId [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	return _EtherApis.Contract.VerifySignature(&_EtherApis.CallOpts, subscriptionId, nonce, value, v, r, s)
}

// VerifySignature is a free data retrieval call binding the contract method 0xf60744d5.
//
// Solidity: function verifySignature(subscriptionId bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) constant returns(bool)
func (_EtherApis *EtherApisCallerSession) VerifySignature(subscriptionId [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (bool, error) {
	return _EtherApis.Contract.VerifySignature(&_EtherApis.CallOpts, subscriptionId, nonce, value, v, r, s)
}

// AddService is a paid mutator transaction binding the contract method 0x3950fc56.
//
// Solidity: function addService(name string, endpoint string, price uint256, cancellationTime uint256) returns()
func (_EtherApis *EtherApisTransactor) AddService(opts *bind.TransactOpts, name string, endpoint string, price *big.Int, cancellationTime *big.Int) (*types.Transaction, error) {
	return _EtherApis.contract.Transact(opts, "addService", name, endpoint, price, cancellationTime)
}

// AddService is a paid mutator transaction binding the contract method 0x3950fc56.
//
// Solidity: function addService(name string, endpoint string, price uint256, cancellationTime uint256) returns()
func (_EtherApis *EtherApisSession) AddService(name string, endpoint string, price *big.Int, cancellationTime *big.Int) (*types.Transaction, error) {
	return _EtherApis.Contract.AddService(&_EtherApis.TransactOpts, name, endpoint, price, cancellationTime)
}

// AddService is a paid mutator transaction binding the contract method 0x3950fc56.
//
// Solidity: function addService(name string, endpoint string, price uint256, cancellationTime uint256) returns()
func (_EtherApis *EtherApisTransactorSession) AddService(name string, endpoint string, price *big.Int, cancellationTime *big.Int) (*types.Transaction, error) {
	return _EtherApis.Contract.AddService(&_EtherApis.TransactOpts, name, endpoint, price, cancellationTime)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(subscriptionId bytes32) returns()
func (_EtherApis *EtherApisTransactor) Cancel(opts *bind.TransactOpts, subscriptionId [32]byte) (*types.Transaction, error) {
	return _EtherApis.contract.Transact(opts, "cancel", subscriptionId)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(subscriptionId bytes32) returns()
func (_EtherApis *EtherApisSession) Cancel(subscriptionId [32]byte) (*types.Transaction, error) {
	return _EtherApis.Contract.Cancel(&_EtherApis.TransactOpts, subscriptionId)
}

// Cancel is a paid mutator transaction binding the contract method 0xc4d252f5.
//
// Solidity: function cancel(subscriptionId bytes32) returns()
func (_EtherApis *EtherApisTransactorSession) Cancel(subscriptionId [32]byte) (*types.Transaction, error) {
	return _EtherApis.Contract.Cancel(&_EtherApis.TransactOpts, subscriptionId)
}

// Claim is a paid mutator transaction binding the contract method 0x3e8b1dd7.
//
// Solidity: function claim(subscriptionId bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) returns()
func (_EtherApis *EtherApisTransactor) Claim(opts *bind.TransactOpts, subscriptionId [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _EtherApis.contract.Transact(opts, "claim", subscriptionId, nonce, value, v, r, s)
}

// Claim is a paid mutator transaction binding the contract method 0x3e8b1dd7.
//
// Solidity: function claim(subscriptionId bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) returns()
func (_EtherApis *EtherApisSession) Claim(subscriptionId [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _EtherApis.Contract.Claim(&_EtherApis.TransactOpts, subscriptionId, nonce, value, v, r, s)
}

// Claim is a paid mutator transaction binding the contract method 0x3e8b1dd7.
//
// Solidity: function claim(subscriptionId bytes32, nonce uint256, value uint256, v uint8, r bytes32, s bytes32) returns()
func (_EtherApis *EtherApisTransactorSession) Claim(subscriptionId [32]byte, nonce *big.Int, value *big.Int, v *big.Int, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _EtherApis.Contract.Claim(&_EtherApis.TransactOpts, subscriptionId, nonce, value, v, r, s)
}

// DeleteService is a paid mutator transaction binding the contract method 0x74e29ee6.
//
// Solidity: function deleteService(serviceId uint256) returns()
func (_EtherApis *EtherApisTransactor) DeleteService(opts *bind.TransactOpts, serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.contract.Transact(opts, "deleteService", serviceId)
}

// DeleteService is a paid mutator transaction binding the contract method 0x74e29ee6.
//
// Solidity: function deleteService(serviceId uint256) returns()
func (_EtherApis *EtherApisSession) DeleteService(serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.Contract.DeleteService(&_EtherApis.TransactOpts, serviceId)
}

// DeleteService is a paid mutator transaction binding the contract method 0x74e29ee6.
//
// Solidity: function deleteService(serviceId uint256) returns()
func (_EtherApis *EtherApisTransactorSession) DeleteService(serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.Contract.DeleteService(&_EtherApis.TransactOpts, serviceId)
}

// Deposit is a paid mutator transaction binding the contract method 0xb214faa5.
//
// Solidity: function deposit(subscriptionId bytes32) returns()
func (_EtherApis *EtherApisTransactor) Deposit(opts *bind.TransactOpts, subscriptionId [32]byte) (*types.Transaction, error) {
	return _EtherApis.contract.Transact(opts, "deposit", subscriptionId)
}

// Deposit is a paid mutator transaction binding the contract method 0xb214faa5.
//
// Solidity: function deposit(subscriptionId bytes32) returns()
func (_EtherApis *EtherApisSession) Deposit(subscriptionId [32]byte) (*types.Transaction, error) {
	return _EtherApis.Contract.Deposit(&_EtherApis.TransactOpts, subscriptionId)
}

// Deposit is a paid mutator transaction binding the contract method 0xb214faa5.
//
// Solidity: function deposit(subscriptionId bytes32) returns()
func (_EtherApis *EtherApisTransactorSession) Deposit(subscriptionId [32]byte) (*types.Transaction, error) {
	return _EtherApis.Contract.Deposit(&_EtherApis.TransactOpts, subscriptionId)
}

// DisableService is a paid mutator transaction binding the contract method 0x91499e2d.
//
// Solidity: function disableService(serviceId uint256) returns()
func (_EtherApis *EtherApisTransactor) DisableService(opts *bind.TransactOpts, serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.contract.Transact(opts, "disableService", serviceId)
}

// DisableService is a paid mutator transaction binding the contract method 0x91499e2d.
//
// Solidity: function disableService(serviceId uint256) returns()
func (_EtherApis *EtherApisSession) DisableService(serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.Contract.DisableService(&_EtherApis.TransactOpts, serviceId)
}

// DisableService is a paid mutator transaction binding the contract method 0x91499e2d.
//
// Solidity: function disableService(serviceId uint256) returns()
func (_EtherApis *EtherApisTransactorSession) DisableService(serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.Contract.DisableService(&_EtherApis.TransactOpts, serviceId)
}

// EnableService is a paid mutator transaction binding the contract method 0x78fe2951.
//
// Solidity: function enableService(serviceId uint256) returns()
func (_EtherApis *EtherApisTransactor) EnableService(opts *bind.TransactOpts, serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.contract.Transact(opts, "enableService", serviceId)
}

// EnableService is a paid mutator transaction binding the contract method 0x78fe2951.
//
// Solidity: function enableService(serviceId uint256) returns()
func (_EtherApis *EtherApisSession) EnableService(serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.Contract.EnableService(&_EtherApis.TransactOpts, serviceId)
}

// EnableService is a paid mutator transaction binding the contract method 0x78fe2951.
//
// Solidity: function enableService(serviceId uint256) returns()
func (_EtherApis *EtherApisTransactorSession) EnableService(serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.Contract.EnableService(&_EtherApis.TransactOpts, serviceId)
}

// Reclaim is a paid mutator transaction binding the contract method 0x96afb365.
//
// Solidity: function reclaim(subscriptionId bytes32) returns()
func (_EtherApis *EtherApisTransactor) Reclaim(opts *bind.TransactOpts, subscriptionId [32]byte) (*types.Transaction, error) {
	return _EtherApis.contract.Transact(opts, "reclaim", subscriptionId)
}

// Reclaim is a paid mutator transaction binding the contract method 0x96afb365.
//
// Solidity: function reclaim(subscriptionId bytes32) returns()
func (_EtherApis *EtherApisSession) Reclaim(subscriptionId [32]byte) (*types.Transaction, error) {
	return _EtherApis.Contract.Reclaim(&_EtherApis.TransactOpts, subscriptionId)
}

// Reclaim is a paid mutator transaction binding the contract method 0x96afb365.
//
// Solidity: function reclaim(subscriptionId bytes32) returns()
func (_EtherApis *EtherApisTransactorSession) Reclaim(subscriptionId [32]byte) (*types.Transaction, error) {
	return _EtherApis.Contract.Reclaim(&_EtherApis.TransactOpts, subscriptionId)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(serviceId uint256) returns()
func (_EtherApis *EtherApisTransactor) Subscribe(opts *bind.TransactOpts, serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.contract.Transact(opts, "subscribe", serviceId)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(serviceId uint256) returns()
func (_EtherApis *EtherApisSession) Subscribe(serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.Contract.Subscribe(&_EtherApis.TransactOpts, serviceId)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(serviceId uint256) returns()
func (_EtherApis *EtherApisTransactorSession) Subscribe(serviceId *big.Int) (*types.Transaction, error) {
	return _EtherApis.Contract.Subscribe(&_EtherApis.TransactOpts, serviceId)
}
