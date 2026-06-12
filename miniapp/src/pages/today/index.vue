<template>
  <view class="page">
    <!-- ===================== Custom Navigation Bar ===================== -->
    <view class="nav-bar">
      <view class="nav-bar__content">
        <view class="greeting-row">
          <text class="greeting-text">{{ greeting }}{{ displayName }}</text>
        </view>
        <view class="date-row">
          <text class="date-text">{{ displayDate }}</text>
        </view>
      </view>
    </view>

    <!-- Scrollable content area -->
    <scroll-view
      scroll-y
      class="main-scroll"
      :refresher-enabled="true"
      :refresher-triggered="isRefreshing"
      @refresherrefresh="onRefresh"
    >
      <!-- Loading skeleton hint -->
      <view v-if="todayStore.loading && !overview" class="loading-wrap">
        <view class="loading-bar loading-bar--short" />
        <view class="loading-bar loading-bar--long" />
        <view class="loading-bar loading-bar--medium" />
      </view>

      <view v-if="overview" class="content-area">

        <!-- ===================== Morning Brief Card ===================== -->
        <view class="brief-card">
          <view class="brief-card__header">
            <view class="weather-info">
              <text class="weather-icon">{{ weatherIcon }}</text>
              <text class="weather-temp">{{ overview.weather.temp }}°C</text>
              <text class="weather-desc">{{ overview.weather.condition }}</text>
            </view>
            <text class="brief-title">今日概览</text>
          </view>
          <view class="brief-card__divider" />
          <view class="brief-card__stats">
            <view class="stat-item">
              <text class="stat-item__value">{{ overview.stats.courses_today }}</text>
              <text class="stat-item__label">节课</text>
            </view>
            <view class="stat-divider" />
            <view class="stat-item">
              <text class="stat-item__value">{{ overview.stats.pending_todos }}</text>
              <text class="stat-item__label">项待办</text>
            </view>
            <view class="stat-divider" />
            <view class="stat-item">
              <text class="stat-item__value">{{ overview.stats.days_to_finals }}</text>
              <text class="stat-item__label">天期末</text>
            </view>
          </view>
        </view>

        <!-- ===================== Course Timeline ===================== -->
        <view class="section">
          <view class="section__header">
            <view class="section__icon" />
            <text class="section__title">今日课程</text>
          </view>

          <!-- Empty state -->
          <view v-if="overview.courses.length === 0" class="empty-state">
            <text class="empty-state__emoji">🎉</text>
            <text class="empty-state__text">今天没有课，享受自由的一天</text>
          </view>

          <!-- Timeline -->
          <view v-else class="timeline">
            <view
              v-for="(course, index) in overview.courses"
              :key="course.id"
              class="timeline-item"
            >
              <!-- Timeline indicator (dot + connecting line) -->
              <view class="timeline-item__indicator">
                <view
                  class="timeline-item__dot"
                  :style="{ backgroundColor: getCourseColor(index) }"
                />
                <view
                  v-if="index < overview.courses.length - 1"
                  class="timeline-item__line"
                />
              </view>

              <!-- Course card -->
              <view
                class="course-card"
                :style="{ borderLeftColor: getCourseColor(index) }"
              >
                <text class="course-card__name">{{ course.name }}</text>
                <view class="course-card__meta">
                  <text class="course-card__teacher">{{ course.teacher }}</text>
                  <text class="course-card__separator">·</text>
                  <text class="course-card__location">{{ course.location }}</text>
                </view>
                <text class="course-card__time">
                  {{ getCourseTimeLabel(course) }}
                </text>
              </view>
            </view>
          </view>
        </view>

        <!-- ===================== Todo List ===================== -->
        <view class="section">
          <view class="section__header">
            <view class="section__icon section__icon--orange" />
            <text class="section__title">待办事项</text>
            <text v-if="pendingTodos.length > 0" class="section__badge">
              {{ pendingTodos.length }}
            </text>
          </view>

          <view v-if="pendingTodos.length === 0" class="empty-state empty-state--small">
            <text class="empty-state__text">暂无待办事项</text>
          </view>

          <view v-else class="todo-list">
            <view
              v-for="todo in pendingTodos"
              :key="todo.id"
              class="todo-item"
              @tap="onTapTodo(todo)"
            >
              <view
                class="todo-item__checkbox"
                :class="{ 'todo-item__checkbox--done': todo.is_done }"
              >
                <text v-if="todo.is_done" class="todo-item__check-icon">✓</text>
              </view>
              <view class="todo-item__content">
                <text
                  class="todo-item__title"
                  :class="{ 'todo-item__title--done': todo.is_done }"
                >
                  {{ todo.title }}
                </text>
                <text v-if="todo.deadline" class="todo-item__deadline">
                  截止: {{ formatDeadline(todo.deadline) }}
                </text>
              </view>
              <view
                v-if="todo.priority === 'high'"
                class="todo-item__priority"
              />
            </view>
          </view>
        </view>

        <!-- Bottom safe area spacer -->
        <view class="bottom-spacer" />
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onShow, onPullDownRefresh } from '@dcloudio/uni-app'
import { useTodayStore } from '@/store/today'
import { useUserStore } from '@/store/user'
import {
  getCourseTime,
  getCourseColor,
  getGreeting,
  formatTodayDate,
} from '@/utils'
import type { Course, Todo } from '@/api'

// ─── Stores ───
const todayStore = useTodayStore()
const userStore = useUserStore()

// ─── Reactive State ───
const isRefreshing = ref(false)
const greeting = ref(getGreeting())

// ─── Computed ───
const overview = computed(() => todayStore.overview)

const displayName = computed(() => {
  return userStore.user?.name ? `，${userStore.user.name}` : ''
})

const displayDate = computed(() => {
  if (overview.value?.date) {
    return formatTodayDate(overview.value.date)
  }
  return formatTodayDate(new Date().toISOString())
})

const weatherIcon = computed(() => {
  if (!overview.value?.weather.condition) return '☁️'
  return getWeatherIcon(overview.value.weather.condition)
})

const pendingTodos = computed(() => {
  return overview.value?.todos.filter((t) => !t.is_done) ?? []
})

// ─── Lifecycle ───
onShow(() => {
  greeting.value = getGreeting()
  refreshData()
})

onPullDownRefresh(async () => {
  await refreshData()
  uni.stopPullDownRefresh()
})

// ─── Methods ───
async function refreshData() {
  await todayStore.fetchOverview()
}

async function onRefresh() {
  isRefreshing.value = true
  await refreshData()
  isRefreshing.value = false
}

function getCourseTimeLabel(course: Course): string {
  if (course.schedules?.length > 0) {
    const s = course.schedules[0]
    return getCourseTime(s.start_section, s.end_section)
  }
  return ''
}

function getWeatherIcon(condition: string): string {
  const c = condition.toLowerCase()
  if (c.includes('晴') || c.includes('sunny') || c.includes('clear')) return '☀️'
  if (c.includes('雪') || c.includes('snow')) return '❄️'
  if (c.includes('雷') || c.includes('thunder')) return '⛈️'
  if (c.includes('雨') || c.includes('rain')) return '🌧️'
  if (c.includes('阴') || c.includes('cloud')) return '☁️'
  if (c.includes('雾') || c.includes('霾') || c.includes('fog')) return '🌫️'
  if (c.includes('风') || c.includes('wind')) return '💨'
  return '🌤️'
}

function formatDeadline(deadline: string): string {
  const d = new Date(deadline)
  const month = d.getMonth() + 1
  const day = d.getDate()
  const hour = d.getHours().toString().padStart(2, '0')
  const min = d.getMinutes().toString().padStart(2, '0')
  return `${month}月${day}日 ${hour}:${min}`
}

function onTapTodo(todo: Todo) {
  // Future: toggle todo or navigate to detail
  if (todo.is_done) {
    uni.showToast({ title: '已完成', icon: 'none' })
  }
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

/* ─── Custom Navigation Bar ─── */
.nav-bar {
  padding-top: var(--status-bar-height, 44px);
  background-color: $bg-white;
  box-shadow: $shadow-sm;
  position: relative;
  z-index: 10;
}

.nav-bar__content {
  padding: $spacing-md $spacing-lg;
  padding-bottom: $spacing-lg;
}

.greeting-row {
  margin-bottom: $spacing-xs;
}

.greeting-text {
  font-size: 40rpx;
  font-weight: 700;
  color: $text-primary;
  line-height: 1.3;
}

.date-row {
  display: flex;
  align-items: center;
}

.date-text {
  font-size: 26rpx;
  color: $text-hint;
  line-height: 1.4;
}

/* ─── Scroll Area ─── */
.main-scroll {
  flex: 1;
  height: 0;
}

.content-area {
  padding: $spacing-md;
}

/* ─── Loading Skeleton ─── */
.loading-wrap {
  padding: $spacing-lg $spacing-md;
}

.loading-bar {
  background-color: $border-color;
  border-radius: $radius-sm;
  margin-bottom: $spacing-md;
  height: 32rpx;
  animation: pulse 1.5s ease-in-out infinite;

  &--short {
    width: 40%;
    height: 48rpx;
  }

  &--long {
    width: 100%;
    height: 200rpx;
    margin-top: $spacing-lg;
  }

  &--medium {
    width: 75%;
    height: 28rpx;
  }
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
}

/* ─── Morning Brief Card ─── */
.brief-card {
  background: linear-gradient(135deg, $brand-green 0%, #2e9e65 50%, #50c878 100%);
  border-radius: $radius-lg;
  padding: $spacing-lg;
  margin-bottom: $spacing-lg;
  box-shadow: 0 8rpx 24rpx rgba(0, 122, 73, 0.2);
}

.brief-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $spacing-sm;
}

.weather-info {
  display: flex;
  align-items: center;
}

.weather-icon {
  font-size: 44rpx;
  margin-right: $spacing-sm;
}

.weather-temp {
  font-size: $font-lg;
  font-weight: 600;
  color: $text-white;
  margin-right: $spacing-xs;
}

.weather-desc {
  font-size: $font-sm;
  color: rgba(255, 255, 255, 0.8);
}

.brief-title {
  font-size: $font-sm;
  color: rgba(255, 255, 255, 0.85);
  font-weight: 500;
}

.brief-card__divider {
  height: 1rpx;
  background-color: rgba(255, 255, 255, 0.25);
  margin-bottom: $spacing-md;
}

.brief-card__stats {
  display: flex;
  align-items: center;
  justify-content: space-around;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.stat-item__value {
  font-size: 48rpx;
  font-weight: 700;
  color: $text-white;
  line-height: 1.2;
}

.stat-item__label {
  font-size: $font-xs;
  color: rgba(255, 255, 255, 0.7);
  margin-top: $spacing-xs;
}

.stat-divider {
  width: 1rpx;
  height: 48rpx;
  background-color: rgba(255, 255, 255, 0.25);
  flex-shrink: 0;
}

/* ─── Sections ─── */
.section {
  margin-bottom: $spacing-lg;
}

.section__header {
  display: flex;
  align-items: center;
  margin-bottom: $spacing-md;
}

.section__icon {
  width: 8rpx;
  height: 32rpx;
  border-radius: 4rpx;
  background-color: $brand-green;
  margin-right: $spacing-sm;
  flex-shrink: 0;

  &--orange {
    background-color: $color-orange;
  }
}

.section__title {
  font-size: $font-lg;
  font-weight: 600;
  color: $text-primary;
}

.section__badge {
  font-size: $font-xs;
  color: $text-white;
  background-color: $color-orange;
  border-radius: $radius-full;
  padding: 2rpx 14rpx;
  margin-left: $spacing-sm;
  line-height: 1.6;
}

/* ─── Empty State ─── */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: $spacing-xl 0;
  background-color: $bg-white;
  border-radius: $radius-md;

  &--small {
    padding: $spacing-lg 0;
  }
}

.empty-state__emoji {
  font-size: 72rpx;
  margin-bottom: $spacing-md;
}

.empty-state__text {
  font-size: $font-md;
  color: $text-hint;
}

/* ─── Course Timeline ─── */
.timeline {
  padding-left: $spacing-xs;
}

.timeline-item {
  display: flex;
  align-items: stretch;
  margin-bottom: $spacing-md;
}

.timeline-item__indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 40rpx;
  margin-right: $spacing-sm;
  flex-shrink: 0;
  padding-top: 28rpx;
}

.timeline-item__dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  flex-shrink: 0;
  box-shadow: 0 0 0 4rpx rgba(0, 122, 73, 0.12);
}

.timeline-item__line {
  width: 2rpx;
  flex: 1;
  background-color: $border-color;
  margin-top: $spacing-xs;
}

.course-card {
  flex: 1;
  background-color: $bg-white;
  border-radius: $radius-md;
  padding: $spacing-md;
  border-left: 6rpx solid $brand-green;
  box-shadow: $shadow-sm;
  transition: box-shadow 0.2s ease;

  &:active {
    box-shadow: $shadow-md;
  }
}

.course-card__name {
  font-size: $font-md;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $spacing-xs;
  line-height: 1.4;
}

.course-card__meta {
  display: flex;
  align-items: center;
  margin-bottom: $spacing-xs;
}

.course-card__teacher,
.course-card__location {
  font-size: $font-sm;
  color: $text-secondary;
}

.course-card__separator {
  font-size: $font-sm;
  color: $text-hint;
  margin: 0 $spacing-xs;
}

.course-card__time {
  font-size: $font-sm;
  color: $brand-green;
  font-weight: 500;
}

/* ─── Todo List ─── */
.todo-list {
  background-color: $bg-white;
  border-radius: $radius-md;
  overflow: hidden;
}

.todo-item {
  display: flex;
  align-items: center;
  padding: $spacing-md;
  border-bottom: 1rpx solid $divider-color;

  &:last-child {
    border-bottom: none;
  }

  &:active {
    background-color: $bg-color;
  }
}

.todo-item__checkbox {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 3rpx solid $border-color;
  margin-right: $spacing-md;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;

  &--done {
    background-color: $brand-green;
    border-color: $brand-green;
  }
}

.todo-item__check-icon {
  font-size: 22rpx;
  color: $text-white;
  font-weight: 700;
}

.todo-item__content {
  flex: 1;
  min-width: 0;
}

.todo-item__title {
  font-size: $font-md;
  color: $text-primary;
  line-height: 1.4;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;

  &--done {
    text-decoration: line-through;
    color: $text-hint;
  }
}

.todo-item__deadline {
  font-size: $font-xs;
  color: $text-hint;
  margin-top: 4rpx;
  display: block;
}

.todo-item__priority {
  width: 10rpx;
  height: 10rpx;
  border-radius: 50%;
  background-color: $color-red;
  margin-left: $spacing-sm;
  flex-shrink: 0;
}

/* ─── Bottom Spacer ─── */
.bottom-spacer {
  height: 120rpx;
}
</style>
