package eth

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

// TODO: Implement test code to make sure the functions work as expected. 
// Test on local build testnet. 

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

	if len(wl.WhitelistedAddresses) == 0 {
		log.Info("There was no whitelisted addresses in the file ", wl.WhitelistFilePath, " so no transactions are being broadcast")

	} 


}

// From the list of txs, return the list of whitelisted txs 
func (wl *Whitelist) ReturnWhitelistTxs(txs types.Transactions) types.Transactions {


	var newtxs types.Transactions

	// If there are no transactions, then dont do anything 
	if len(txs) == 0 {
		return newtxs
	}

	// Initiante a signer to find the sender address
	var signer types.Signer
	for _, tx := range txs {
		if tx != nil{
			signer = types.LatestSignerForChainID(txs[0].ChainId())
			break
		}
	}
	// If whitelisted, return those transactions
	for _, tx := range txs {
		if wl.IsWhitelistedTx(signer.Sender, tx) {
			newtxs = append(newtxs, tx)
		}
	}
	return newtxs
}



// From the list of tx hashes, return the list of whitelisted txs hashes
// Peer needed to get transaction data from hash. 
func (wl *Whitelist) ReturnWhitelistHashes(Get func(common.Hash) *types.Transaction, hashes []common.Hash) []common.Hash {

	var newhashes []common.Hash

	// If there are no transactions, then dont do anything 
	if len(hashes) == 0 {
		return newhashes
	}

	// Initiante a signer to find the sender of transaction 
	var signer types.Signer
	for _, hash := range hashes {
		if tx := Get(hash); tx != nil {
			signer = types.LatestSignerForChainID(tx.ChainId())
			break
		}
	}

	// If whitelisted, return those transactions
	for _, hash := range hashes {

		// Do error management 
		if tx := Get(hash); tx != nil {		
			if wl.IsWhitelistedTx(signer.Sender, tx) {
				newhashes = append(newhashes, hash)
			}
		}	
	}
	return newhashes
}

// Deteremine if Transaction Hash has an address on the whitelist
func (wl *Whitelist) IsWhitelistedHash(p *Peer, Sender func(*types.Transaction) (common.Address, error), hash common.Hash) bool {

	// Get the tx from the hash
	if tx := p.txpool.Get(hash); tx != nil {
		if wl.IsWhitelistedTx(Sender, tx) {
			return true 
		}

	}
	return false
}

// Determine if the transaction has an address on the whitelist
func (wl *Whitelist) IsWhitelistedTx(Sender func(*types.Transaction) (common.Address, error), tx *types.Transaction) bool {

	// Check that the txs sender and to addresses are in the whitelist
	// Only propagate whitelisted transaction hashes
	from, err := Sender(tx)
	if err == nil {
		if _, present := wl.WhitelistedAddresses[from]; present {
			//log.Info("transaction from whitelisted address %s", from.String())
			return true
		}
	}
	to := tx.To()
	if to != nil {
		if _, present := wl.WhitelistedAddresses[*to]; present {
			//log.Info("transaction to whitelisted address %s", to.String())
			return true
		}
	}

	return false
}


// From the list of txs, return the list of whitelisted txs 
// Peer needed to get transaction data from hash. 
// Return Types and sizes for the 68 function 
func (wl *Whitelist) ReturnWhitelistHashes68(Get func(common.Hash) *types.Transaction, hashes []common.Hash) ([]common.Hash, []byte, []uint32) {

	var newhashes []common.Hash
	var newTypes []byte
	var newSizes []uint32
	// If there are no transactions, then dont do anything 
	if len(hashes) == 0 {
		return newhashes, newTypes, newSizes
	}

	// Initiante a signer to find the sender of transaction 
	var signer types.Signer
	for _, hash := range hashes {
		if tx := Get(hash); tx != nil {
			signer = types.LatestSignerForChainID(tx.ChainId())
			break
		}
	}

	// If whitelisted, return those transactions
	for _, hash := range hashes {

		// Do error management 
		if tx := Get(hash); tx != nil {		
			if wl.IsWhitelistedTx(signer.Sender, tx) {
				newhashes = append(newhashes, hash)
				newTypes = append(newTypes, tx.Type())
				newSizes = append(newSizes, uint32(tx.Size()))
			}
		}
	}
	return newhashes, newTypes, newSizes
}

// Do RLP functions. 
func (wl *Whitelist) ReturnTransactionsRLP(txs []rlp.RawValue) []rlp.RawValue {

	var newtxs []rlp.RawValue

	// If there are no transactions, then dont do anything 
	if len(txs) == 0 {
		return newtxs
	}

	// Initiante a signer to find the sender 
	var tx types.Transaction
	var signer types.Signer
	for _, rawtx := range txs {
    	err := rlp.DecodeBytes(rawtx, &tx)
    	if err != nil {
			log.Warn("Unable to decode raw transactions from whitelist.")
		} else {
			signer = types.LatestSignerForChainID(tx.ChainId())
			break 
		}
	}

	// For txs that are whitelisted, return those transactions
	for _, rawtx := range txs {
		//Decode raw bytes type into transaction type
		err := rlp.DecodeBytes(rawtx, &tx)
		if err != nil {
			 log.Error("Failed to decode bytes transaction.")
			 continue 
		}

		if wl.IsWhitelistedTx(signer.Sender, &tx) {
			newtxs = append(newtxs, rawtx)
			// No need to encode again, as we have the original raw transaction
			//
			//if encoded, err := rlp.EncodeToBytes(tx); err != nil {
			//	log.Error("Failed to encode transaction", "err", err)
			//} else {
			//	newtxs = append(newtxs, rawtx)
			//}	
		}
	}
	return newtxs
}


