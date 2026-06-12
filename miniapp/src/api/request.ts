/**
 * Unified HTTP request wrapper for uni-app.
 * Handles token injection, automatic refresh, and error handling.
 */

const BASE_URL = 'https://api.scau-daily.com/api/v1' // TODO: replace with real domain
// const BASE_URL = 'http://localhost:8080/api/v1' // Local dev

interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  data?: Record<string, unknown>
  header?: Record<string, string>
  showLoading?: boolean
}

interface ApiResponse<T = unknown> {
  data: T
  error?: string
  message?: string
}

// Get stored tokens
function getToken(): string {
  return uni.getStorageSync('access_token') || ''
}

function getRefreshToken(): string {
  return uni.getStorageSync('refresh_token') || ''
}

// Store tokens after login/refresh
export function saveTokens(accessToken: string, refreshToken: string) {
  uni.setStorageSync('access_token', accessToken)
  uni.setStorageSync('refresh_token', refreshToken)
}

export function clearTokens() {
  uni.removeStorageSync('access_token')
  uni.removeStorageSync('refresh_token')
}

// Refresh access token
async function refreshToken(): Promise<boolean> {
  const refresh = getRefreshToken()
  if (!refresh) return false

  try {
    const res = await uni.request({
      url: `${BASE_URL}/auth/refresh`,
      method: 'POST',
      data: { refresh_token: refresh },
    })
    if (res.statusCode === 200 && res.data.access_token) {
      saveTokens(res.data.access_token, res.data.refresh_token)
      return true
    }
  } catch (e) {
    console.error('[API] refresh failed', e)
  }
  return false
}

// Main request function
export async function request<T = unknown>(options: RequestOptions): Promise<T> {
  const { url, method = 'GET', data, header = {}, showLoading = false } = options

  if (showLoading) {
    uni.showLoading({ title: '加载中...', mask: true })
  }

  // Inject auth token
  const token = getToken()
  if (token) {
    header['Authorization'] = `Bearer ${token}`
  }

  try {
    const res = await uni.request({
      url: `${BASE_URL}${url}`,
      method,
      data,
      header: {
        'Content-Type': 'application/json',
        ...header,
      },
    })

    // Success
    if (res.statusCode >= 200 && res.statusCode < 300) {
      return res.data as T
    }

    // Token expired — try refresh once
    if (res.statusCode === 401) {
      const refreshed = await refreshToken()
      if (refreshed) {
        // Retry original request with new token
        const retryToken = getToken()
        header['Authorization'] = `Bearer ${retryToken}`
        const retryRes = await uni.request({
          url: `${BASE_URL}${url}`,
          method,
          data,
          header: { 'Content-Type': 'application/json', ...header },
        })
        if (retryRes.statusCode >= 200 && retryRes.statusCode < 300) {
          return retryRes.data as T
        }
      }
      // Refresh failed — redirect to login
      clearTokens()
      uni.navigateTo({ url: '/pages/login/index' })
      throw new Error('登录已过期，请重新登录')
    }

    // Other errors
    const errMsg = (res.data as Record<string, string>)?.message || '请求失败'
    throw new Error(errMsg)
  } finally {
    if (showLoading) {
      uni.hideLoading()
    }
  }
}
