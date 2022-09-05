# 是否为 CI 环境？
INCI ?= false

.PHONY: all
all: help

.PHONY: help
help:
	@echo "Usage: make <command>"
	@echo "The commands are:"
	@echo "    web3-message		down web3 Message"
	@echo "    binary			Build binaries"
	@echo "    binary-linux		Build linux"	
	@echo "    binary-windows		Build windows"	
	@echo "    binary-darwin		Build mac"


.PHONY: binary
binary:  web3-message account

.PHONY: binary-darwin
binary-darwin: web3-message account-darwin 

.PHONY: binary-windows
binary-windows: web3-message account-windows 

.PHONY: binary-linux
binary-linux: web3-message account-linux 


#  download web3 dapr message 
.PHONY: web3-message
web3-message: 
	@./scripts/down_web3_msg.sh



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
