// Contains some pre-configured metadata for setting up an Ethereum client.

package geth

// EthereumNetwork is the variant of the Ethereum network to join.
type EthereumNetwork int

const (
	MainNet = 0 // Frontier live network
	TestNet = 2 // Morden test network
)

const (
	NodeName     = "Geth/EtherAPIs" // Client name to advertise on the Ethereum network
	NodeVersion  = "0.0.1"          // Client version to advertise on the Ethereum network
	NodePort     = 30303            // Listener port of the Ethereum P2P network
	NodeMaxPeers = 25               // Maximum number of peers connections to accept
)
