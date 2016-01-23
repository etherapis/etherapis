// Contains the custom HTTP headers defined by the payment proxy.

package proxy

const (
	SubscriptionHeader  = "etherapi-subscripton"  // Used by a client to authorize payment of a particular subscription
	AuthorizationHeader = "etherapi-authorize"    // Cumulative amount of payments to authorize (previous + current)
	SignatureHeader     = "etherapi-signature"    // Client signature to verify the payment authorization
	UnauthorizedHeader  = "etherapi-unauthorized" // Header sent back by the proxy with the authorization failure
)
