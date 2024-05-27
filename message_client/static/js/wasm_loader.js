/**
 * Defines a Go WASM instance. Objects of this class are returned by the
 * `loadWasm()` function and allow interaction with the WASM module and
 * instance plus the Golang RT post-initialization. Refer to the following
 * docs for usage:
 * https://github.com/golang/go/blob/b2fcfc1a50fbd46556f7075f7f1fbf600b5c9e5d/misc/wasm/wasm_exec.html
 */
class WasmInst {
	constructor(go, mod, inst){
		this.go = go;
		this.mod = mod;
		this.inst = inst;
	}
}

/**
 * Compiles a given WASM file and instantiates the environment. After this
 * function is called, the exported WASM functions from Golang will become
 * available to any JS code block needing to use it. Simply call the target
 * function as normal and it'll run the Go code in the fetched WASM binary.
 * This function MUST be called AFTER `wasm_exec.js` is loaded. Otherwise,
 * it will fail, since the Golang bindings will be unavailable.
 * See: https://stackoverflow.com/a/76082718
 * See: https://github.com/golang/go/blob/b2fcfc1a50fbd46556f7075f7f1fbf600b5c9e5d/misc/wasm/wasm_exec.html
 * 
 * @param wasmURL The URL from which to load the WASM binary
 * @param runNow Whether the Go RT should be run immediately or later when using the instance
 * @return A `WasmInst` object wrapped in a `Promise`; ie: `Promise<WasmInst>`
 */
function loadWasm(wasmURL, runNow){
	//Setup a polyfill if `WebAssembly.instantiateStreaming()` is unavailable for whatever reason
	if(!WebAssembly.instantiateStreaming){
		WebAssembly.instantiateStreaming = async (resp, importObject) => {
			const source = await (await resp).arrayBuffer();
			return await WebAssembly.instantiate(source, importObject);
		};
	}
	
	//Create a new Golang WASM handler
	const go = new Go(); //Defined in `wasm_exec.js`
	let mod, inst;

	//Initialize the Promise that will be returned; Thanks CGPT
	return new Promise((resolve, reject) => {
		//Fetch the WASM file and instantiate it
		WebAssembly.instantiateStreaming(fetch(wasmURL), go.importObject).then(async (result) => {
			//Initialize the mod and inst variables
			mod = await result.module;
			inst = await result.instance;

			//Initialize a new `WasmInst` object
			const obj = await new WasmInst(go, mod, inst);

			//Run the Go code if requested
			if(runNow){
				obj.go.run(obj.inst);
				console.log("Ran Go WASM RT; instance: ", obj);
			}

			//Resolve the promise with the object
			resolve(obj);
		}).catch((err) => {
			//Report the error and reject the promise
			console.error(err);
			reject(error);
		});
	});
}