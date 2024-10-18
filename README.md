# Bund-Chain

## Overview
Learn by doing: build a blockchain from scratch using Go, inspired by Bitcoin and Ethereum. 
Note that this project is not production-ready.

### Feature List

- **Consensus Mechanism:** Proof of Work or Proof of Authority
- **Pool Selector:** Best Tip
- **File Storage Format:** JSON
- **Node Communication:** HTTP
- **Peer Discovery Method:** Known Peers (similar to Ethereum)
- **Transaction Validation:** Merkle Tree
- **Digital Signature:** Custom Stamp before Encryption (similar to Bitcoin)
- **Name Service:** Transform address to readable name base on private key file (For develop only)
- **Public API:** HTTP


## Project Structure
```
├── bin                 # binary directory
├── cmd             
│   ├── node            # Bund-chain node cmd
│   └── wallet          # Helper tools for Bund-chain wallet
│       ├── cli
├── conf                # genesis block and miner config
│   ├── accounts        # private keys directory
├── data
│   ├── miner1          # storage in json for node 1
│   ├── miner2          # storage in json for node 2
│   └── miner3          # storage in json for node 3
└── internal
    ├── blockchain
    │   ├── database    # Operation of Bund-Chain from account to TX
    │   ├── genesis     # Load config for Genesis block
    │   ├── mempool     # Pool of transaction
    │   │   └── selector
    │   ├── merkle      # Merkel tree 
    │   ├── peer        # Keep track peer info
    │   ├── state       # Centralize API
    │   ├── storage     # Writing to storage operation
    │   │   └── disk
    │   └── worker      # Concurreny handler
    ├── handler         # Public and P2P API
    ├── nameservice     # Traslate address to name (Readability)
    └── utils
        ├── getenv      # Get OS environment
        └── signature   # Digital signature utils
```

## How to run

I have included everything that needs to be run, including the wallet and private key. 
If you want to run manually or start from scratch, please go to the next section.

We run three miner nodes with a default Proof of Work consensus. 
If you want to run Proof of Authority, please edit `conf/miner1.env`, `conf/miner2.env`, and `conf/miner3.env`.

You can also edit the Genesis config block in `conf/genesis.json`
Increase the difficulty level if it's progressing too quickly.

Terminal 1

`make up`

Terminal 2

`make up2`

Terminal 3

`make up3`

We can now send six transactions using the command-line interface (CLI).

Terminal 4

`make load`

If you want to send more `make load2` `make load3`

The mined block will be stored in `data/miner1`, `data/miner2`, and `data/miner3` due to synchronization (replication).

You can list the Accounts by accessing API `http://localhost:3000/accounts/list`
You can also list the pool `http://localhost:3000/tx/uncommitted/list`

For more routes, Please check `cmd/node/routes.go`

## How to run from scratch

1. Edit the Genesis block configuration in the file `conf/genesis.conf`.
2. Create a wallet for the user by running the following command: `go run cmd/wallet/main.go generate -a YOUR_NAME`. 
3. This will generate a private key stored in `conf/accounts` with a `.ecdsa` extension.
3. To view the wallet address, use this command: `go run cmd/wallet/main.go account -a YOUR_NAME`.
4. To create a miner wallet, repeat steps 2 and 3.
5. Create an ENV variable for the miner you created in step 4 under `conf/YOUR_MINER.env`:
```
WEB_ADDR="0.0.0.0:3000"
WEB_PRIVATE_ADDR="0.0.0.0:3030"
BENEFICIARY="YOUR_MINER"
DB_PATH="data/YOUR_MINER/"
CONSENSUS="POW"
```
6. Run a node using `go run cmd/node/*.go YOUR_MINER`. 
The arguments must match the exact name of the ENV file located in `conf/YOUR_MINER.env`.
7. Send a transaction using the CLI.
      `go run cmd/wallet/main.go send -a bund -n 1 -f 0x7D69D992d41542B81dc9663F1c79EDd5A0d62B54 -t 0xcb3a4F9DdA0c18a5b5e183b3bb01259D29FE7B37 -v 100`
```
-a = account
-f = from account address
-t = to account address
-v = amount
```