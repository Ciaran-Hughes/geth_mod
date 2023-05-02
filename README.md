# Whitelisted-Propagation Geth 

This repo contains a modified version of geth which only broadcasts transactions to peers if the "sender" and/or "to" addresses are on a user-defined whitelist. 

The official version of geth is [https://github.com/ethereum/go-ethereum](https://eips.ethereum.org/EIPS/eip-2464)

## Description

Ethereum transactions are propagated according to [EIP-2464](https://eips.ethereum.org/EIPS/eip-2464). This repo changes geth so that the node only sends transactions (or transaction hashes) to peers if the "sender" and/or "to" addresses of a transaction are in a user-defined whitelist. The whitelist is provided as a json array in a file called "whitelisted_addresses.json", in the run directory. 

The mempool is not altered in this version of geth, so that MEV strategies can be performed on pending mempool transactions, but any network overhead from propagating other users transactions is reduced. Note that not gossiping transactions is counter to the blockchain ethos. This whitelist-only node can also be used by solo users who want to submit transactions without participating in the blockchain. Further, to ensure forward compatibility with geth and reduce the risk of edge-cases, changes to the geth code base are kept to a minimum, and this new whitelist is included as an additional modular feature. 


## Usage 

To build a standalone version of geth with this whitelisted feature, clone this repo and
```
make geth 
```

Then run the geth executable as usual. Afterwards, whitelisted geth can be connected to MetaMask on a localnode. 

To build a docker-compose version of whitelist-propagation geth as the execution client, and the official lighthouse as the consensus client, follow the instructions at [https://github.com/Ciaran-Hughes/ETHNode_Docker](https://github.com/Ciaran-Hughes/ETHNode_Docker) with the docker files in the current repo. 

<!--If running in docker, put the whitelisted_address.json in the docker persistent directory where geth is run. The docker-compose is run as usual. -->


## License

This project is licensed under the MIT License as given in License.txt
