//https://gist.github.com/waynegraham/5766565

/*
//Set a destination for system logs
if(typeof console != "undefined"){
	console.olog = typeof console.log != "undefined" ? console.log : function(){};
}

console.log = function(message){
	console.olog(message);
	document.getElementById("result").innerHTML += "<span class='line' style='display: block; word-wrap: break-word'>" + (message + "").replaceAll("\n", "<br>") + "</span>";
};
console.error = console.debug = console.info = console.log;
*/

//Define a function to capture data sent to a console stream
const sinker = function(message){
	document.getElementById("result").innerHTML += "<span class='line' style='display: block; word-wrap: break-word'>" + (message + "").replaceAll("\n", "<br>") + "</span>";
}

//Define a sink initializer
const sinkInit = function(sinkName){
	//Check if console exists and has the sink method
	if (typeof console !== "undefined" && typeof console[sinkName] === "function"){
		//Get the original sink function
		const osink = console[sinkName];

		//Remap the sinked function to a wrapper
		console[sinkName] = function(message){
			osink(message);
			sinker(message);
		};
	}
	else console.warn(`Console or console.${sinkName} is not available.`);
}

// Sink the standard streams
const streams = ["log", "debug", "error", "info", "warn"];
streams.forEach((e) => sinkInit(e));
