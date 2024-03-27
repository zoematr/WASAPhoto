import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});
// axios.defaults.baseURL = __API_URL__;
export default instance;
