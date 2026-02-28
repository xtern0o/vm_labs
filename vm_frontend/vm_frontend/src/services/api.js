import axios from "axios";
import config from "./config.js";

const api = axios.create({
    baseURL: config.apiBaseUrl,
    timeout: config.timeout,
    headers: {
        'Content-Type': 'application/json',
    },
});

export default api;
export { config };