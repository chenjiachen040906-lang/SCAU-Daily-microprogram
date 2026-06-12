import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, AuthResponse } from '../api'
import { getMe } from '../api'
import { saveTokens, clearTokens } from '../api/request'

export const useUserStore = defineStore('user', () => {
  const user = ref<UserInfo | null>(null)
  const isLoggedIn = computed(() => !!user.value && !!user.value.student_id)

  function checkAuth() {
    const token = uni.getStorageSync('access_token')
    if (token) {
      // Token exists, try to fetch user info
      getMe()
        .then((u) => { user.value = u })
        .catch(() => {
          // Token invalid, clear
          clearTokens()
          user.value = null
        })
    }
  }

  function setAuth(resp: AuthResponse) {
    saveTokens(resp.access_token, resp.refresh_token)
    user.value = resp.user
  }

  function logout() {
    clearTokens()
    user.value = null
    uni.reLaunch({ url: '/pages/login/index' })
  }

  return { user, isLoggedIn, checkAuth, setAuth, logout }
})
