### MASTER MAKEFILE ###

## EDIT THESE LINES TO CHANGE THE ENTRYPOINT & BINARY NAME ##
ENTRYPOINT_NAME := main.go
BIN_NAME := sms_server
BIN_FOLDER := .

####################################
#Variables
####################################
# Variables
#Adds a .exe if the OS is Windows
ifeq ($(OS), Windows_NT)
BIN_NAME := $(BIN_NAME).exe
endif

# Tools
FRESH := fresh
GO := go
GO_ENUM := go-enum
#MAKE := $(MAKE) -s

#Flags
GO_LD := -ldflags "-s -w"

#Define common commands
ifeq ($(OS), Windows_NT)
	CP := @copy /y
	MV := @move /y
	MKDIR := @mkdir
	RMDIR := @rmdir /S /Q
	RMFILE := @del /S /Q
	WHICH := where
	TOUCH := type nul >
else 
	CP := @cp
	MV := @mv
	MKDIR := @mkdir
	RMDIR := @rm -rf
	RMFILE := @rm -f
	WHICH := which
	TOUCH := touch
endif

####################################


####################################
#Functions
####################################
#Makes a new directory if it doesn't already exist. Call via `$(MAKE) -s __mkdir__ TARGET=<dirname>`
.PHONY: __mkdir__
__mkdir__:
ifeq ($(wildcard ./$(TARGET)/.),) #Test if target directory is nonexistant
	@echo MKDIR $(TARGET)
	$(MKDIR) $(TARGET)
endif

#Removes a directory recursively if it exists. Call via `$(MAKE) -s __rmdir__ TARGET=<dirname>`
.PHONY: __rmdir__
__rmdir__:
ifneq ($(wildcard ./$(TARGET)/.),) #Test if target directory is existant
	@echo RM -RF $(TARGET)
	$(RMDIR) $(TARGET)
endif

#Removes a file if it exists. Call via `$(MAKE) -s __rmfile__ TARGET=<filename>`
.PHONY: __rmfile__
__rmfile__:
ifneq ($(wildcard $(TARGET)),) # Test if target file is existant
	@echo RM -F $(TARGET)
	$(RMFILE) $(TARGET)
endif

#Copies a file to a destination if it exists. Call via `$(MAKE) -s __cpfile__ SRC=<source> DEST=<destination>`
.PHONY: __cpfile__
__cpfile__:
ifneq ($(wildcard $(SRC)),) # Test if target file is existant
	@echo CP $(SRC) $(DEST)
	$(CP) "$(SRC)" "$(DEST)"
endif

#Moves a file to a destination if it exists. Call via `$(MAKE) -s __mvfile__ SRC=<source> DEST=<destination>`
.PHONY: __mvfile__
__mvfile__:
ifneq ($(wildcard $(SRC)),) # Test if target file is existant
	@echo MV $(SRC) $(DEST)
	$(MV) "$(SRC)" "$(DEST)"
endif

####################################


####################################
# Phony Build Rules
####################################
#Builds the prohect
.PHONY: build
build:
	@echo BUILD $(ENTRYPOINT_NAME) "->" $(BIN_FOLDER)/$(BIN_NAME)
	$(GO) build $(GO_LD) -o $(BIN_FOLDER)/$(BIN_NAME) $(ENTRYPOINT_NAME)

#Cleans the directory of any binary files
.PHONY: clean
clean:
	@echo CLEAN $(BIN_FOLDER)/$(BIN_NAME)
	@$(MAKE) -s __rmfile__ TARGET=$(BIN_FOLDER)/$(BIN_NAME)

#Runs the binary with Fresh (watches for changes and auto-recompiles)
.PHONY: dev
dev:
	@echo FRESH .
	$(FRESH)

#Generates enum files based on `go-enum` syntax
.PHONY: enum
enum:
	@echo GO-ENUM .
	$(GO) generate ./...

#Runs the generated binary
.PHONY: run
run:
	@echo RUN $(BIN_FOLDER)/$(BIN_NAME)
	$(BIN_FOLDER)/$(BIN_NAME)