#sudo apt install uglifyjs

#tinygo
#Copies the `wasm_exec.js` files from the locally installed Go repo to the `static/js` folder
#TINYGO_EXEC := /$$(readlink $$(which tinygo))../../targets/wasm_exec.js
TINYGO_EXEC := /usr/local/lib/tinygo/targets/wasm_exec.js
.PHONY: updatejs
updatejs:
	@echo UPDATEJS
	realpath $(TINYGO_EXEC)
#$(MAKE) -s __cpfile__ SRC=$(TINYGO_EXEC) DEST=$(STATIC_DIR)


#norm go
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
