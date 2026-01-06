import { boot } from 'quasar/wrappers'
import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080'
})

export default boot(() => {
  const token = localStorage.getItem('token')
  if (token) {
    api.defaults.headers.common.Authorization = `Bearer ${token}`
  }
})

export { api }
