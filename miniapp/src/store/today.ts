import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { TodayOverview } from '../api'
import { getTodayOverview } from '../api'

export const useTodayStore = defineStore('today', () => {
  const overview = ref<TodayOverview | null>(null)
  const loading = ref(false)
  const error = ref('')

  async function fetchOverview() {
    loading.value = true
    error.value = ''
    try {
      overview.value = await getTodayOverview()
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : '加载失败'
      error.value = msg
    } finally {
      loading.value = false
    }
  }

  return { overview, loading, error, fetchOverview }
})
