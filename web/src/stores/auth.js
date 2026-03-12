import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, register as registerApi } from '@/api/user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)
  const userName = computed(() => user.value?.username || '')
  const userRole = computed(() => user.value?.role || '')

  async function login(credentials) {
    const res = await loginApi(credentials)
    if (res.data?.token) {
      token.value = res.data.token
      user.value = {
        id: res.data.user_id,
        username: res.data.username,
        role: res.data.role
      }
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(user.value))
    }
    return res
  }

  async function register(credentials) {
    const res = await registerApi(credentials)
    return res
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    token,
    user,
    isLoggedIn,
    userName,
    userRole,
    login,
    register,
    logout
  }
})
