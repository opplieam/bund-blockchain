SHELL := /bin/bash

# Bund Address: 0x7D69D992d41542B81dc9663F1c79EDd5A0d62B54
# Opp Address: 0x00fb343D49B7E1cc90C03442EDD4468D4c4EAf49
# miner1 Address: 0xE45e25f67C6cf24CBBC39fA6c6d4a5ee5cEdBBB2
# miner2 Address: 0x2b5e8A61c178D7504f56C99e6dcf6275B871a95f
# miner3 Address: 0xfF75720644b5f40041C9dB0d4Cdc025A930EA939

# Generate Private key
# go run cmd/wallet/main.go generate
# View Address
# go run cmd/wallet/main.go account -a bund
# go run cmd/wallet/main.go account -a opp

run-scratch:
	go run cmd/scratch/main.go