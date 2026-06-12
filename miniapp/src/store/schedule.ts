import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Course } from '../api'
import { getAllCourses, getWeekCourses } from '../api'

export const useScheduleStore = defineStore('schedule', () => {
  const courses = ref<Course[]>([])
  const currentWeek = ref(1)
  const loading = ref(false)

  async function fetchAllCourses() {
    loading.value = true
    try {
      const res = await getAllCourses()
      courses.value = res.courses || []
    } catch (e) {
      console.error('[Schedule] fetch courses failed', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchWeekCourses(week: number) {
    loading.value = true
    try {
      const res = await getWeekCourses(week)
      courses.value = res.courses || []
      currentWeek.value = week
    } catch (e) {
      console.error('[Schedule] fetch week courses failed', e)
    } finally {
      loading.value = false
    }
  }

  return { courses, currentWeek, loading, fetchAllCourses, fetchWeekCourses }
})
