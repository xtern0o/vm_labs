const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000'

const REQUEST_TIMEOUT = 30000

export const config = {
  apiBaseUrl: API_BASE_URL,
  timeout: REQUEST_TIMEOUT,
}

export default config
