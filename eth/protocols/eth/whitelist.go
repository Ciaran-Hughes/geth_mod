package eth

import (
	"github.com/ethereum/go-ethereum/common"
)

// Read in the whitelist_addresses.json file
// Todo: Make this a command flag

const (
	whitelist_filename string = "whitelist_addresses.json"
)

type Whitelist struct {
	WhitelistFilePathFlag string
	WhiteListedAddresses  []common.Address
}
