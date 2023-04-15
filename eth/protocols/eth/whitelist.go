package eth

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// Hardcode the whitelist_addresses.json file path
// Todo: Make this a command line flag parameter

// The whitelist file should be a json array of addresses
const (
	whitelist_filename string = "whitelist_addresses.json"
)

// Whitelist contains relevant information about the Whitelist parameters
type Whitelist struct {
	WhitelistFilePath    string                  // Path of the whitelist file.
	WhitelistedAddresses map[common.Address]bool // Use a map for optimized complexity of operations
}

// Make a New Whitelist and initialize it
func NewWhitelist() *Whitelist {

	wl := &Whitelist{
		WhitelistFilePath: whitelist_filename,
	}

	wl.ReadWhitelistFile()
	log.Info("Read in whitelist file ", wl.WhitelistFilePath,
		". Only allowing these addresses in transaction/hash propagation")

	return wl
}

// Set the file path for the whitelist object
func (wl *Whitelist) SetFilePath() {
	wl.WhitelistFilePath = whitelist_filename
}

// Set the file path for the whitelist object
func (wl *Whitelist) ReadWhitelistFile() {

	// Read in the whitelist file
	bytes, err := os.ReadFile(wl.WhitelistFilePath)
	if err != nil {
		log.Warn("There was an Error ", err, " Opening the Whitelist File Path : ",
			wl.WhitelistFilePath, ". Defaulting to not using a whitelist")
	}

	var wl_array []common.Address

	// Unmarshal the bytes
	if err := json.Unmarshal(bytes, &wl_array); err != nil {
		log.Warn("There was an Error ", err, " Unmarshaling the Whitelist File from Path : ",
			wl.WhitelistFilePath, ". Note the format of the file is expected to be ",
			" a json array of addresses as strings. Defaulting to not using a whitelist")
	}

	for _, address := range wl_array {
		wl.WhitelistedAddresses[address] = true
	}

}
