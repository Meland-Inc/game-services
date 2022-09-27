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
binary: web3-message account agent chat main manager task

.PHONY: binary-darwin
binary-darwin: web3-message account-darwin agent-darwin chat-darwin main-darwin manager-darwin task-darwin

.PHONY: binary-windows
binary-windows: web3-message account-windows agent-windows chat-windows main-windows manager-windows task-windows

.PHONY: binary-linux
binary-linux: web3-message account-linux agent-linux chat-linux main-linux manager-linux task-linux 
 

# ------------- download web3 dapr message -----------------
.PHONY: web3-message
web3-message: 
	@./scripts/down_web3_msg.sh



# ------------- build  account to target OS -----------------
.PHONY: account
account: 
	@./scripts/build_binary.sh 	account 	services/account
.PHONY: account-darwin
account-darwin: 
	@./scripts/build_binary.sh 	account 	services/account 	darwin
.PHONY: account-windows
account-windows: 
	@./scripts/build_binary.sh 	account 	services/account 	windows	
.PHONY: account-linux
account-linux: 
	@./scripts/build_binary.sh 	account 	services/account 	linux		

# ------------- build  agent to target OS -----------------
.PHONY: agent
agent: 
	@./scripts/build_binary.sh 	agent 	services/agent
.PHONY: agent-darwin
agent-darwin: 
	@./scripts/build_binary.sh 	agent 	services/agent 	darwin
.PHONY: agent-windows
agent-windows: 
	@./scripts/build_binary.sh 	agent 	services/agent 	windows	
.PHONY: agent-linux
agent-linux: 
	@./scripts/build_binary.sh 	agent 	services/agent 	linux	
 

# ------------- build  chat to target OS -----------------
.PHONY: chat
chat: 
	@./scripts/build_binary.sh 	chat 	services/chat
.PHONY: chat-darwin
chat-darwin: 
	@./scripts/build_binary.sh 	chat 	services/chat 	darwin
.PHONY: chat-windows
chat-windows: 
	@./scripts/build_binary.sh 	chat 	services/chat 	windows	
.PHONY: chat-linux
chat-linux: 
	@./scripts/build_binary.sh 	chat 	services/chat 	linux		 


# ------------- build  main to target OS -----------------
.PHONY: main
main: 
	@./scripts/build_binary.sh 	main 	services/main
.PHONY: main-darwin
main-darwin: 
	@./scripts/build_binary.sh 	main 	services/main 	darwin
.PHONY: main-windows
main-windows: 
	@./scripts/build_binary.sh 	main 	services/main 	windows	
.PHONY: main-linux
main-linux: 
	@./scripts/build_binary.sh 	main 	services/main 	linux			


# ------------- build  manager to target OS -----------------
.PHONY: manager
manager: 
	@./scripts/build_binary.sh 	manager 	services/manager
.PHONY: manager-darwin
manager-darwin: 
	@./scripts/build_binary.sh 	manager 	services/manager 	darwin
.PHONY: manager-windows
manager-windows: 
	@./scripts/build_binary.sh 	manager 	services/manager 	windows	
.PHONY: manager-linux
manager-linux: 
	@./scripts/build_binary.sh 	manager 	services/manager 	linux		


# ------------- build  task to target OS -----------------
.PHONY: task
task: 
	@./scripts/build_binary.sh 	task 	services/task
.PHONY: task-darwin
task-darwin: 
	@./scripts/build_binary.sh 	task 	services/task 	darwin
.PHONY: task-windows
task-windows: 
	@./scripts/build_binary.sh 	task 	services/task 	windows	
.PHONY: task-linux
task-linux: 
	@./scripts/build_binary.sh 	task 	services/task 	linux			
