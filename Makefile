SHELL := /bin/bash

# Bund Address: 0x7D69D992d41542B81dc9663F1c79EDd5A0d62B54
# Opp Address: 0x00fb343D49B7E1cc90C03442EDD4468D4c4EAf49

# Generate Private key
# go run cmd/wallet/main.go generate
# View Address
# go run cmd/wallet/main.go account -a bund
# go run cmd/wallet/main.go account -a opp

run-scratch:
	go run cmd/scratch/main.go