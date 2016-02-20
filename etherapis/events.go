// Contains various events fired by Ether APIs itself.

package etherapis

import "github.com/ethereum/go-ethereum/common"

// NewAccountEvent is posted when a new account was created or imported.
type NewAccountEvent struct{ Address common.Address }

// DroppedAccountEvent is posted when an account is permanently deleted.
type DroppedAccountEvent struct{ Address common.Address }
