/**
 * API endpoint functions — calls the backend through the unified request wrapper.
 */
import { request, saveTokens } from './request'

// ===== Types =====

export interface UserInfo {
  id: string
  student_id: string
  name: string
  department: string
  major: string
  grade: string
  avatar_url: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  user: UserInfo
}

export interface Course {
  id: string
  name: string
  teacher: string
  location: string
  course_type: string
  credit: number
  exam_type: string
  schedules: CourseSchedule[]
}

export interface CourseSchedule {
  day_of_week: number
  start_section: number
  end_section: number
  weeks: string
}

export interface Todo {
  id: string
  title: string
  description: string
  deadline: string | null
  is_done: boolean
  priority: string
  source: string
  created_at: string
}

export interface Notification {
  id: string
  title: string
  content: string
  source: string
  category: string
  published_at: string
  ai_summary: string
  ai_deadline: string | null
  raw_url: string
}

export interface TodayOverview {
  date: string
  weekday: number
  current_week: number
  semester: string
  weather: {
    temp: number
    condition: string
    icon: string
  }
  courses: Course[]
  todos: Todo[]
  notifications: Notification[]
  stats: {
    courses_today: number
    pending_todos: number
    days_to_finals: number
  }
}

export interface WxLoginResponse {
  need_bind: boolean
  data: AuthResponse | null
}

// ===== Auth =====

export function login(studentId: string, password: string) {
  return request<AuthResponse>({
    url: '/auth/login',
    method: 'POST',
    data: { student_id: studentId, password },
  })
}

export function wxLogin(code: string) {
  return request<WxLoginResponse>({
    url: '/auth/wx-login',
    method: 'POST',
    data: { code },
  })
}

export function bindStudent(studentId: string, password: string) {
  return request<AuthResponse>({
    url: '/auth/bind-student',
    method: 'POST',
    data: { student_id: studentId, password },
  })
}

export function getMe() {
  return request<UserInfo>({ url: '/auth/me' })
}

// ===== Schedule =====

export function syncSchedule(password: string) {
  return request<{ message: string; courses: Course[]; count: number }>({
    url: '/schedule/sync',
    method: 'POST',
    data: { password },
    showLoading: true,
  })
}

export function getTodayCourses(date?: string) {
  return request<{ date: string; courses: Course[] }>({
    url: '/schedule/today',
    method: 'GET',
    data: date ? { date } : undefined,
  })
}

export function getWeekCourses(week: number) {
  return request<{ week: number; courses: Course[] }>({
    url: '/schedule/week',
    method: 'GET',
    data: { week: String(week) },
  })
}

export function getAllCourses() {
  return request<{ courses: Course[] }>({
    url: '/schedule/courses',
    method: 'GET',
  })
}

// ===== Todos =====

export function getTodos(status?: string) {
  return request<{ todos: Todo[] }>({
    url: '/todos',
    method: 'GET',
    data: status ? { status } : undefined,
  })
}

export function createTodo(data: { title: string; description?: string; deadline?: string; priority?: string }) {
  return request<Todo>({
    url: '/todos',
    method: 'POST',
    data,
  })
}

export function updateTodo(id: string, data: Partial<Todo>) {
  return request<Todo>({
    url: `/todos/${id}`,
    method: 'PATCH',
    data,
  })
}

export function deleteTodo(id: string) {
  return request({
    url: `/todos/${id}`,
    method: 'DELETE',
  })
}

// ===== Today =====

export function getTodayOverview() {
  return request<TodayOverview>({
    url: '/today/overview',
    method: 'GET',
  })
}
