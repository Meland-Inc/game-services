# 是否为 CI 环境？
INCI ?= false

.PHONY: all
all: help

.PHONY: help
help:
	@echo "Usage: make <command>"
	@echo "The commands are:"
	@echo "    binary		Build binaries"
	@echo "    binary-linux	Build linux"	
	@echo "    binary-windows	Build windows"	
	@echo "    binary-darwin	Build mac"



.PHONY: binary
binary: account 

.PHONY: binary-darwin
binary: account-darwin 

.PHONY: binary-windows
binary: account-windows 

.PHONY: binary-linux
binary: account-linux 



#  build account  to target OS
.PHONY: account
account: 
	@./scripts/build_binary.sh 	account 	src/cmd/account

.PHONY: account-darwin
account-darwin: 
	@./scripts/build_binary.sh 	account 	src/cmd/account 	darwin

.PHONY: account-windows
account-windows: 
	@./scripts/build_binary.sh 	account 	src/cmd/account	 	windows
	
.PHONY: account-linux
account-linux: 
	@./scripts/build_binary.sh 	account 	src/cmd/account  	linux
