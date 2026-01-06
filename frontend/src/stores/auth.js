import { defineStore } from 'pinia'
import { api } from 'boot/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token')
  }),

  getters: {
    isLogged: (s) => !!s.token
  },

  actions: {
    init () {
      if (this.token) {
        api.defaults.headers.common.Authorization = `Bearer ${this.token}`
      }
    },

    async login (email, password) {
      const { data } = await api.post('/auth/login', { email, password })
      this.token = data.token
      localStorage.setItem('token', data.token)
      api.defaults.headers.common.Authorization = `Bearer ${data.token}`
    },

    logout () {
      this.token = null
      localStorage.removeItem('token')
      delete api.defaults.headers.common.Authorization
    }
  }
})
