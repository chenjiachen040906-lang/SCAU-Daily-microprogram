<template>
  <view class="page">
    <!-- ===================== Week Switcher ===================== -->
    <view class="week-switcher">
      <view class="week-switcher__pill">
        <view class="week-arrow" @tap="prevWeek">
          <text class="week-arrow__icon">◀</text>
        </view>
        <text class="week-switcher__text">第 {{ currentWeek }} 周</text>
        <view class="week-arrow" @tap="nextWeek">
          <text class="week-arrow__icon">▶</text>
        </view>
      </view>
    </view>

    <!-- Loading state -->
    <view v-if="scheduleStore.loading" class="loading-state">
      <view class="loading-spinner" />
      <text class="loading-text">加载中...</text>
    </view>

    <!-- Empty state -->
    <view
      v-else-if="!hasCourses"
      class="empty-state"
    >
      <text class="empty-state__emoji">📭</text>
      <text class="empty-state__text">本周暂无课程</text>
    </view>

    <!-- ===================== Schedule Grid ===================== -->
    <scroll-view v-else scroll-x scroll-y class="grid-scroll">
      <view class="schedule-grid">
        <!-- Header row: weekday names -->
        <view class="grid-header">
          <view class="grid-header__time-col">
            <text class="grid-header__label">节次</text>
          </view>
          <view
            v-for="day in 5"
            :key="day"
            class="grid-header__day"
            :class="{ 'grid-header__day--today': day === todayWeekday }"
          >
            <text class="grid-header__day-name">{{ WEEKDAY_NAMES[day] }}</text>
          </view>
        </view>

        <!-- Grid body: sections × days -->
        <view class="grid-body">
          <view
            v-for="section in displaySections"
            :key="section"
            class="grid-row"
          >
            <!-- Time column -->
            <view class="grid-row__time">
              <text class="grid-row__section-num">{{ section }}</text>
              <text class="grid-row__time-range">{{ getSectionTimeRange(section) }}</text>
            </view>

            <!-- Day cells -->
            <view
              v-for="day in 5"
              :key="day"
              class="grid-cell"
              :class="{ 'grid-cell--today': day === todayWeekday }"
            >
              <view class="grid-cell__inner">
                <!--
                  Only render course block at its start_section.
                  The block height spans all occupied sections.
                -->
                <template v-if="!isCellHidden(section, day)">
                  <view
                    v-if="getCourseBlock(section, day)"
                    class="course-block"
                    :style="courseBlockStyle(getCourseBlock(section, day)!)"
                    @tap="showCourseDetail(getCourseBlock(section, day)!.course)"
                  >
                    <text class="course-block__name">
                      {{ abbreviateName(getCourseBlock(section, day)!.course.name) }}
                    </text>
                    <text class="course-block__location">
                      {{ getCourseBlock(section, day)!.course.location }}
                    </text>
                  </view>
                  <view v-else class="grid-cell__empty">
                    <text class="grid-cell__dash">-</text>
                  </view>
                </template>
              </view>
            </view>
          </view>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useScheduleStore } from '@/store/schedule'
import {
  SECTION_TIMES,
  WEEKDAY_NAMES,
  getCourseColor,
} from '@/utils'
import type { Course } from '@/api'

// ─── Constants ───
const MAX_WEEK = 25
const MIN_WEEK = 1
const SECTION_HEIGHT = 100 // rpx per section
const CELL_GAP = 4 // rpx gap inside course block

// ─── Stores ───
const scheduleStore = useScheduleStore()

// ─── Reactive State ───
const currentWeek = ref(1)

// ─── Grid Data Structure ───
interface GridCellData {
  course: Course
  colorIndex: number
  startSection: number
  endSection: number
  spanCount: number
}

/**
 * Build a 2D lookup: section × day → course block info.
 * A cell is "hidden" when a course that started at an earlier section
 * still occupies it (multi-section span).
 */
const gridData = computed(() => {
  const data: Record<string, GridCellData> = {}
  const hidden: Record<string, boolean> = {}

  scheduleStore.courses.forEach((course, courseIndex) => {
    course.schedules?.forEach((schedule) => {
      // Only include courses for the current week
      if (!isCourseInWeek(schedule.weeks, currentWeek.value)) return
      // Only Mon–Fri
      if (schedule.day_of_week < 1 || schedule.day_of_week > 5) return

      const startSec = schedule.start_section
      const endSec = schedule.end_section
      const spanCount = endSec - startSec + 1

      // Place the course block at the start section
      data[`${startSec}-${schedule.day_of_week}`] = {
        course,
        colorIndex: courseIndex,
        startSection: startSec,
        endSection: endSec,
        spanCount,
      }

      // Mark subsequent sections as hidden
      for (let s = startSec + 1; s <= endSec; s++) {
        hidden[`${s}-${schedule.day_of_week}`] = true
      }
    })
  })

  return { data, hidden }
})

/** Whether any courses exist for the current week. */
const hasCourses = computed(() => {
  return Object.keys(gridData.value.data).length > 0
})

/** Current weekday (1 = Mon ... 5 = Fri). */
const todayWeekday = computed(() => {
  const d = new Date().getDay()
  return d === 0 ? 7 : d
})

/**
 * Determine which section rows to display.
 * Shows only the range that has courses, with sensible defaults.
 */
const displaySections = computed(() => {
  const sections: number[] = []
  let minSec = 11
  let maxSec = 1

  // Find the range of sections that have courses
  Object.values(gridData.value.data).forEach((block) => {
    minSec = Math.min(minSec, block.startSection)
    maxSec = Math.max(maxSec, block.endSection)
  })

  // Fallback to 1–8 if no courses
  if (minSec > maxSec) {
    minSec = 1
    maxSec = 8
  }

  for (let i = minSec; i <= maxSec; i++) {
    sections.push(i)
  }
  return sections
})

// ─── Lifecycle ───
onMounted(() => {
  loadWeek(currentWeek.value)
})

onShow(() => {
  // Refresh data when returning to this page
  if (scheduleStore.courses.length === 0) {
    loadWeek(currentWeek.value)
  }
})

// ─── Helpers ───

/**
 * Determine whether a course schedule's `weeks` string includes a given week.
 * Supports formats like "1-16", "1-16(单)", "9-16(双)", "1,3,5,7".
 */
function isCourseInWeek(weeksStr: string, week: number): boolean {
  if (!weeksStr) return false

  // Handle comma-separated weeks: "1,3,5,7"
  if (weeksStr.includes(',') && !weeksStr.includes('-')) {
    return weeksStr.split(',').some((w) => parseInt(w.trim(), 10) === week)
  }

  // Match range patterns like "1-16", "1-16(单)", "1-16(双)"
  const match = weeksStr.match(/(\d+)\s*-\s*(\d+)(?:\s*[\(（](单|双)[\)）])?/)
  if (match) {
    const start = parseInt(match[1], 10)
    const end = parseInt(match[2], 10)
    const parity = match[3] // '单' (odd) | '双' (even) | undefined

    if (week < start || week > end) return false

    if (parity === '单') return week % 2 === 1
    if (parity === '双') return week % 2 === 0
    return true
  }

  // Fallback: try direct number match
  return weeksStr.includes(String(week))
}

// ─── Methods ───
function loadWeek(week: number) {
  currentWeek.value = week
  scheduleStore.fetchWeekCourses(week)
}

function prevWeek() {
  if (currentWeek.value > MIN_WEEK) {
    loadWeek(currentWeek.value - 1)
  }
}

function nextWeek() {
  if (currentWeek.value < MAX_WEEK) {
    loadWeek(currentWeek.value + 1)
  }
}

function getCourseBlock(section: number, day: number): GridCellData | undefined {
  return gridData.value.data[`${section}-${day}`]
}

function isCellHidden(section: number, day: number): boolean {
  return !!gridData.value.hidden[`${section}-${day}`]
}

function courseBlockStyle(block: GridCellData): Record<string, string> {
  const height = block.spanCount * SECTION_HEIGHT - CELL_GAP * 2
  const color = getCourseColor(block.colorIndex)
  return {
    height: `${height}rpx`,
    backgroundColor: color,
  }
}

function getSectionTimeRange(section: number): string {
  const time = SECTION_TIMES[section]
  if (!time) return ''
  return `${time.start}\n${time.end}`
}

function abbreviateName(name: string): string {
  // Truncate long course names to fit the narrow cell
  if (name.length > 6) {
    return name.substring(0, 6) + '…'
  }
  return name
}

function showCourseDetail(course: Course) {
  const schedule = course.schedules?.[0]
  const timeRange = schedule
    ? `${WEEKDAY_NAMES[schedule.day_of_week]} 第${schedule.start_section}-${schedule.end_section}节`
    : '未知时间'

  // Build detail lines
  const lines: string[] = []
  if (course.teacher) lines.push(`教师：${course.teacher}`)
  if (course.location) lines.push(`地点：${course.location}`)
  if (course.credit) lines.push(`学分：${course.credit}`)
  if (course.course_type) lines.push(`类型：${course.course_type}`)
  lines.push(`时间：${timeRange}`)

  if (course.schedules?.length) {
    const time = SECTION_TIMES[course.schedules[0].start_section]
    if (time) {
      const endTime = SECTION_TIMES[course.schedules[0].end_section]
      lines.push(`时段：${time.start} - ${endTime?.end || time.end}`)
    }
  }

  uni.showModal({
    title: course.name,
    content: lines.join('\n'),
    showCancel: false,
    confirmText: '知道了',
    confirmColor: '#007A49',
  })
}
</script>

<style lang="scss" scoped>
@use '../../styles/variables' as *;

/* ─── Page Container ─── */
.page {
  min-height: 100vh;
  background-color: $bg-color;
  display: flex;
  flex-direction: column;
}

/* ─── Week Switcher ─── */
.week-switcher {
  display: flex;
  justify-content: center;
  padding: $spacing-md 0;
  background-color: $bg-white;
  border-bottom: 1rpx solid $divider-color;
}

.week-switcher__pill {
  display: flex;
  align-items: center;
  background-color: $brand-green-light;
  border-radius: $radius-full;
  padding: $spacing-xs $spacing-sm;
}

.week-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56rpx;
  height: 56rpx;
  border-radius: 50%;
  transition: background-color 0.15s;

  &:active {
    background-color: rgba(0, 122, 73, 0.15);
  }
}

.week-arrow__icon {
  font-size: $font-xs;
  color: $brand-green;
}

.week-switcher__text {
  font-size: $font-lg;
  font-weight: 600;
  color: $brand-green;
  margin: 0 $spacing-md;
  min-width: 140rpx;
  text-align: center;
}

/* ─── Loading State ─── */
.loading-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-spinner {
  width: 60rpx;
  height: 60rpx;
  border: 4rpx solid $border-color;
  border-top-color: $brand-green;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: $spacing-md;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  font-size: $font-sm;
  color: $text-hint;
}

/* ─── Empty State ─── */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 160rpx 0;
}

.empty-state__emoji {
  font-size: 80rpx;
  margin-bottom: $spacing-lg;
}

.empty-state__text {
  font-size: $font-lg;
  color: $text-hint;
}

/* ─── Scroll Container ─── */
.grid-scroll {
  flex: 1;
  height: 0;
}

/* ─── Schedule Grid ─── */
.schedule-grid {
  min-width: 100%;
}

/* Grid Header */
.grid-header {
  display: flex;
  background-color: $bg-white;
  border-bottom: 2rpx solid $border-color;
  position: sticky;
  top: 0;
  z-index: 5;
}

.grid-header__time-col {
  width: 96rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: $spacing-sm 0;
  border-right: 1rpx solid $divider-color;
}

.grid-header__label {
  font-size: $font-xs;
  color: $text-hint;
}

.grid-header__day {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: $spacing-sm 0;
  border-right: 1rpx solid $divider-color;

  &:last-child {
    border-right: none;
  }

  &--today {
    background-color: rgba(0, 122, 73, 0.06);
  }
}

.grid-header__day-name {
  font-size: $font-sm;
  font-weight: 600;
  color: $text-secondary;

  .grid-header__day--today & {
    color: $brand-green;
    font-weight: 700;
  }
}

/* Grid Body */
.grid-body {
  background-color: $bg-white;
}

.grid-row {
  display: flex;
  border-bottom: 1rpx solid $divider-color;
}

.grid-row__time {
  width: 96rpx;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $spacing-xs 0;
  border-right: 1rpx solid $divider-color;
  background-color: #fafafa;
}

.grid-row__section-num {
  font-size: $font-md;
  font-weight: 600;
  color: $text-primary;
  line-height: 1.2;
}

.grid-row__time-range {
  font-size: 18rpx;
  color: $text-hint;
  text-align: center;
  line-height: 1.3;
  margin-top: 4rpx;
  white-space: pre-line;
}

/* Grid Cells */
.grid-cell {
  flex: 1;
  min-width: 0;
  height: 100rpx;
  border-right: 1rpx solid $divider-color;
  position: relative;

  &:last-child {
    border-right: none;
  }

  &--today {
    background-color: rgba(0, 122, 73, 0.03);
  }
}

.grid-cell__inner {
  position: relative;
  width: 100%;
  height: 100%;
}

.grid-cell__empty {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.grid-cell__dash {
  font-size: $font-sm;
  color: #ddd;
}

/* Course Blocks */
.course-block {
  position: absolute;
  top: 4rpx;
  left: 4rpx;
  right: 4rpx;
  border-radius: $radius-sm;
  padding: $spacing-xs $spacing-xs;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  z-index: 2;
  box-shadow: 0 2rpx 6rpx rgba(0, 0, 0, 0.1);
  transition: opacity 0.15s;

  &:active {
    opacity: 0.85;
  }
}

.course-block__name {
  font-size: $font-xs;
  font-weight: 600;
  color: $text-white;
  text-align: center;
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  word-break: break-all;
}

.course-block__location {
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.85);
  text-align: center;
  margin-top: 4rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}
</style>
