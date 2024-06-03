### MASTER MAKEFILE ###

## EDIT THESE LINES TO CHANGE THE ENTRYPOINT & BINARY NAME ##
ENTRYPOINT_NAME := js_exports.go
BIN_NAME := ed25519_keygen
BIN_EXT := wasm

####################################
#Variables
####################################
# Variables
#Sets the full binary names
BIN_NORM := $(BIN_NAME).$(BIN_EXT)
BIN_MIN := $(BIN_NAME).min.$(BIN_EXT)

#Denotes where to drop the resultant WASM files
STATIC_DIR := ../static
OUT_DIR := $(STATIC_DIR)/wasm
JS_DIR := $(STATIC_DIR)/js

# Tools
GO_BUILD := go build
TINYGO_BUILD := tinygo build
#MAKE := $(MAKE) -s

#Flags
GO_OS_ARCH := GOARCH=wasm GOOS=js
GO_LD := -ldflags "-s -w"
TINYGO_ARCH := -target wasm
TINYGO_LD := -gc leaking -no-debug

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
#Builds a Golang WASM file using the standard Go compiler
.PHONY: build
build:
	@echo BUILD $(ENTRYPOINT_NAME) "->" $(BIN_NORM)
	$(GO_OS_ARCH) $(GO_BUILD) $(GO_LD) -o $(BIN_NORM) $(ENTRYPOINT_NAME)

#Builds a Golang WASM file using the TinyGo compiler
.PHONY: buildtiny
buildtiny:
	@echo BUILD $(ENTRYPOINT_NAME) "->" $(BIN_MIN)
	$(TINYGO_BUILD) $(TINYGO_ARCH) $(TINYGO_LD) -o $(BIN_MIN) $(ENTRYPOINT_NAME)
	

#Cleans the directory of any WASM files
.PHONY: clean
clean:
	@echo CLEAN $(BIN_NORM) $(BIN_MIN)
	@$(MAKE) -s __rmfile__ TARGET=$(BIN_NORM)
	@$(MAKE) -s __rmfile__ TARGET=$(BIN_MIN)

#Copies the generated WASM files to the static site folder for testing
.PHONY: install
install:
	@echo INSTALL $(BIN_NORM) $(BIN_MIN)
	@$(MAKE) -s __cpfile__ SRC=$(BIN_NORM) DEST=$(OUT_DIR)
	@$(MAKE) -s __cpfile__ SRC=$(BIN_MIN) DEST=$(OUT_DIR)
