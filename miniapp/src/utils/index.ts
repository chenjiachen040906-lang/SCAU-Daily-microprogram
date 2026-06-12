// Section time mapping — SCAU standard schedule
export const SECTION_TIMES: Record<number, { start: string; end: string }> = {
  1: { start: '08:00', end: '08:45' },
  2: { start: '08:55', end: '09:40' },
  3: { start: '10:00', end: '10:45' },
  4: { start: '10:55', end: '11:40' },
  5: { start: '14:30', end: '15:15' },
  6: { start: '15:25', end: '16:10' },
  7: { start: '16:30', end: '17:15' },
  8: { start: '17:25', end: '18:10' },
  9: { start: '19:30', end: '20:15' },
  10: { start: '20:25', end: '21:10' },
  11: { start: '21:20', end: '22:05' },
}

// Weekday names
export const WEEKDAY_NAMES = ['', '周一', '周二', '周三', '周四', '周五', '周六', '周日']

// Get course time range string
export function getCourseTime(startSection: number, endSection: number): string {
  const start = SECTION_TIMES[startSection]?.start || ''
  const end = SECTION_TIMES[endSection]?.end || ''
  return `${start} - ${end}`
}

// Assign a consistent color to a course based on its index
const COURSE_COLORS = [
  '#007A49', '#1890FF', '#722ED1', '#FF9C00', '#FF4D4F',
  '#13C2C2', '#EB2F96', '#2F54EB', '#FA8C16', '#52C41A',
]

export function getCourseColor(index: number): string {
  return COURSE_COLORS[index % COURSE_COLORS.length]
}

// Format greeting based on hour
export function getGreeting(): string {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了'
  if (hour < 11) return '早上好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  return '晚上好'
}

// Format today's date as readable string
export function formatTodayDate(dateStr: string): string {
  const d = new Date(dateStr)
  const month = d.getMonth() + 1
  const day = d.getDate()
  const weekday = WEEKDAY_NAMES[(d.getDay() || 7)] // 0=Sunday -> 7
  return `${month}月${day}日 ${weekday}`
}
