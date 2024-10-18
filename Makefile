SHELL := /bin/bash

# Bund Address: 0x7D69D992d41542B81dc9663F1c79EDd5A0d62B54
# Opp Address: 0x00fb343D49B7E1cc90C03442EDD4468D4c4EAf49
# Kate Address: 0xcb3a4F9DdA0c18a5b5e183b3bb01259D29FE7B37
# Atom Address: 0x4d418d2FE56D2368c657EeB0596f96e31C9AE98d
# Book Address: 0x51508b2AF867f85e2f6D5035DEf65b913d1ED440
# miner1 Address: 0xE45e25f67C6cf24CBBC39fA6c6d4a5ee5cEdBBB2
# miner2 Address: 0x2b5e8A61c178D7504f56C99e6dcf6275B871a95f
# miner3 Address: 0xfF75720644b5f40041C9dB0d4Cdc025A930EA939

# Generate Private key
# go run cmd/wallet/main.go generate -a bund
# View Address
# go run cmd/wallet/main.go account -a bund
# go run cmd/wallet/main.go account -a opp

up:
	go run cmd/node/*.go miner1

up2:
	go run cmd/node/*.go miner2

up3:
	go run cmd/node/*.go miner3

# ==============================================================================
# Transactions

load:
	go run cmd/wallet/main.go send -a bund -n 1 -f 0x7D69D992d41542B81dc9663F1c79EDd5A0d62B54 -t 0xcb3a4F9DdA0c18a5b5e183b3bb01259D29FE7B37 -v 100
	go run cmd/wallet/main.go send -a opp -n 1 -f 0x00fb343D49B7E1cc90C03442EDD4468D4c4EAf49 -t 0xcb3a4F9DdA0c18a5b5e183b3bb01259D29FE7B37 -v 75
	go run cmd/wallet/main.go send -a bund -n 2 -f 0x7D69D992d41542B81dc9663F1c79EDd5A0d62B54 -t 0x4d418d2FE56D2368c657EeB0596f96e31C9AE98d -v 150
	go run cmd/wallet/main.go send -a opp -n 2 -f 0x00fb343D49B7E1cc90C03442EDD4468D4c4EAf49 -t 0x51508b2AF867f85e2f6D5035DEf65b913d1ED440 -v 125
	go run cmd/wallet/main.go send -a bund -n 3 -f 0x7D69D992d41542B81dc9663F1c79EDd5A0d62B54 -t 0x51508b2AF867f85e2f6D5035DEf65b913d1ED440 -v 200
	go run cmd/wallet/main.go send -a opp -n 3 -f 0x00fb343D49B7E1cc90C03442EDD4468D4c4EAf49 -t 0x4d418d2FE56D2368c657EeB0596f96e31C9AE98d -v 250

load2:
	go run cmd/wallet/main.go send -a bund -n 4 -f 0x7D69D992d41542B81dc9663F1c79EDd5A0d62B54 -t 0xcb3a4F9DdA0c18a5b5e183b3bb01259D29FE7B37 -v 100
	go run cmd/wallet/main.go send -a opp -n 4 -f 0x00fb343D49B7E1cc90C03442EDD4468D4c4EAf49 -t 0xcb3a4F9DdA0c18a5b5e183b3bb01259D29FE7B37 -v 75

load3:
	go run cmd/wallet/main.go send -a bund -n 5 -f 0x7D69D992d41542B81dc9663F1c79EDd5A0d62B54 -t 0x4d418d2FE56D2368c657EeB0596f96e31C9AE98d -v 150
	go run cmd/wallet/main.go send -a opp -n 5 -f 0x00fb343D49B7E1cc90C03442EDD4468D4c4EAf49 -t 0x51508b2AF867f85e2f6D5035DEf65b913d1ED440 -v 125
	go run cmd/wallet/main.go send -a bund -n 6 -f 0x7D69D992d41542B81dc9663F1c79EDd5A0d62B54 -t 0x51508b2AF867f85e2f6D5035DEf65b913d1ED440 -v 200
	go run cmd/wallet/main.go send -a opp -n 6 -f 0x00fb343D49B7E1cc90C03442EDD4468D4c4EAf49 -t 0x4d418d2FE56D2368c657EeB0596f96e31C9AE98d -v 250
