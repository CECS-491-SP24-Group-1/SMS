//https://stackoverflow.com/a/53256982

//It is assumed that Axios already exists in the global namespace
//import axios from "axios";

//Sets the maximum delta between the current time and the expiry time before auto-refreshes occur (in seconds).
const MAX_DELTA_TO_EXPIRY = 43200; //12 hours

//Create a new Axios instance; this is needed to prevent infinite loops from the inner request
const axiosInst = axios.create({
	withCredentials: true,
	validateStatus: () => true //https://stackoverflow.com/a/76240990
});

//Export for use elsewhere
const taxios = axiosInst;
//export default taxios;