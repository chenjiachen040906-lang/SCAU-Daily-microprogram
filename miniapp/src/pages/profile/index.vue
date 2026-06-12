<template>
  <view class="profile-page">
    <!-- User header card -->
    <view class="user-card">
      <view class="user-card__bg" />
      <view class="user-card__content">
        <!-- Avatar -->
        <view class="avatar" @tap="goLogin">
          <image
            v-if="userStore.user?.avatar_url"
            class="avatar__img"
            :src="userStore.user.avatar_url"
            mode="aspectFill"
          />
          <view v-else class="avatar__placeholder">
            <text class="avatar__initial">{{ avatarInitial }}</text>
          </view>
        </view>

        <!-- User info -->
        <view class="user-info" v-if="isLoggedIn">
          <text class="user-info__name">{{ userStore.user?.name || '同学' }}</text>
          <text class="user-info__id">{{ userStore.user?.student_id }}</text>
          <text class="user-info__dept" v-if="userStore.user?.department">
            {{ userStore.user.department }}
          </text>
        </view>
        <view class="user-info" v-else @tap="goLogin">
          <text class="user-info__name user-info__name--login">点击登录</text>
          <text class="user-info__id">登录后查看更多功能</text>
        </view>

        <!-- Arrow indicator for not-logged-in -->
        <view class="user-card__arrow" v-if="!isLoggedIn">
          <text class="arrow-text">›</text>
        </view>
      </view>
    </view>

    <!-- Menu section -->
    <view class="section">
      <text class="section__title">功能设置</text>
      <view class="card">
        <!-- Sync schedule -->
        <view class="cell" @tap="handleSyncSchedule">
          <view class="cell__icon cell__icon--sync">
            <text class="cell__icon-text">↻</text>
          </view>
          <view class="cell__body">
            <text class="cell__label">同步课表</text>
            <text class="cell__desc">从教务系统拉取最新课表</text>
          </view>
          <text class="cell__arrow">›</text>
        </view>

        <!-- Semester settings -->
        <view class="cell" @tap="handleSemesterSetting">
          <view class="cell__icon cell__icon--semester">
            <text class="cell__icon-text">▦</text>
          </view>
          <view class="cell__body">
            <text class="cell__label">学期设置</text>
            <text class="cell__desc">选择当前学期</text>
          </view>
          <text class="cell__arrow">›</text>
        </view>

        <!-- Reminder settings -->
        <view class="cell" @tap="handleReminderSetting">
          <view class="cell__icon cell__icon--reminder">
            <text class="cell__icon-text">◈</text>
          </view>
          <view class="cell__body">
            <text class="cell__label">提醒设置</text>
            <text class="cell__desc">App 中可开启系统推送</text>
          </view>
          <text class="cell__arrow">›</text>
        </view>

        <!-- About -->
        <view class="cell cell--last" @tap="handleAbout">
          <view class="cell__icon cell__icon--about">
            <text class="cell__icon-text">ⓘ</text>
          </view>
          <view class="cell__body">
            <text class="cell__label">关于我们</text>
            <text class="cell__desc">版本信息与反馈</text>
          </view>
          <text class="cell__arrow">›</text>
        </view>
      </view>
    </view>

    <!-- Download App card -->
    <view class="section">
      <text class="section__title">获取完整体验</text>
      <view class="download-card">
        <view class="download-card__header">
          <text class="download-card__title">下载 App</text>
          <view class="download-card__badge">
            <text class="download-card__badge-text">推荐</text>
          </view>
        </view>
        <text class="download-card__subtitle">解锁小程序无法提供的功能</text>

        <view class="download-card__features">
          <view class="feature-item">
            <view class="feature-item__dot" />
            <text class="feature-item__text">课前推送提醒</text>
          </view>
          <view class="feature-item">
            <view class="feature-item__dot" />
            <text class="feature-item__text">离线查看课表</text>
          </view>
          <view class="feature-item">
            <view class="feature-item__dot" />
            <text class="feature-item__text">桌面小组件</text>
          </view>
          <view class="feature-item">
            <view class="feature-item__dot" />
            <text class="feature-item__text">无限制 AI 对话</text>
          </view>
        </view>

        <button class="download-card__btn" @tap="handleDownload">
          前往下载
        </button>
      </view>
    </view>

    <!-- Logout button -->
    <view class="logout-section" v-if="isLoggedIn">
      <button class="logout-btn" @tap="handleLogout">
        退出登录
      </button>
    </view>

    <!-- Bottom safe area -->
    <view class="safe-bottom" />
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useUserStore } from '@/store/user'
import { syncSchedule } from '@/api'

const userStore = useUserStore()

/**
 * Whether user is logged in
 */
const isLoggedIn = computed(() => userStore.isLoggedIn)

/**
 * Avatar initial character
 */
const avatarInitial = computed(() => {
  const name = userStore.user?.name
  if (name) return name.charAt(0)
  return '?'
})

/**
 * Navigate to login page
 */
function goLogin() {
  if (!isLoggedIn.value) {
    uni.navigateTo({ url: '/pages/login/index' })
  }
}

/**
 * Sync schedule from educational system
 */
async function handleSyncSchedule() {
  if (!isLoggedIn.value) {
    goLogin()
    return
  }

  // Prompt for password via modal
  uni.showModal({
    title: '同步课表',
    content: '',
    editable: true,
    placeholderText: '请输入教务系统密码',
    confirmText: '开始同步',
    cancelText: '取消',
    success: async (res) => {
      if (res.confirm && res.content) {
        const pwd = res.content.trim()
        if (!pwd) {
          uni.showToast({ title: '请输入密码', icon: 'none' })
          return
        }

        uni.showLoading({ title: '同步中...' })
        try {
          await syncSchedule(pwd)
          uni.showToast({ title: '同步成功', icon: 'success' })
        } catch (err: any) {
          uni.showToast({
            title: err?.message || '同步失败，请检查密码',
            icon: 'none'
          })
        } finally {
          uni.hideLoading()
        }
      }
    }
  })
}

/**
 * Semester settings placeholder
 */
function handleSemesterSetting() {
  if (!isLoggedIn.value) {
    goLogin()
    return
  }
  uni.showToast({ title: '功能开发中', icon: 'none' })
}

/**
 * Reminder settings placeholder
 */
function handleReminderSetting() {
  if (!isLoggedIn.value) {
    goLogin()
    return
  }
  uni.showToast({ title: '功能开发中', icon: 'none' })
}

/**
 * About page placeholder
 */
function handleAbout() {
  uni.showToast({ title: 'SCAU Daily v1.0.0', icon: 'none' })
}

/**
 * Download app
 */
function handleDownload() {
  uni.setClipboardData({
    data: 'https://scau-daily.app/download',
    success: () => {
      uni.showToast({ title: '下载链接已复制', icon: 'none' })
    }
  })
}

/**
 * Logout
 */
function handleLogout() {
  uni.showModal({
    title: '退出登录',
    content: '确定要退出当前账号吗？',
    confirmText: '退出',
    confirmColor: '#e53935',
    success: (res) => {
      if (res.confirm) {
        userStore.logout()
        uni.showToast({ title: '已退出', icon: 'success' })
      }
    }
  })
}
</script>

<style lang="scss" scoped>
@use '../../styles/variables' as *;

.profile-page {
  min-height: 100vh;
  background: #f5f6fa;
  padding-bottom: env(safe-area-inset-bottom);
}

/* User header card */
.user-card {
  position: relative;
  margin: 24rpx 32rpx;
  border-radius: 28rpx;
  overflow: hidden;
  background: #ffffff;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.06);

  &__bg {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 160rpx;
    background: linear-gradient(135deg, #007a49, #00a86b);
  }

  &__content {
    position: relative;
    display: flex;
    align-items: flex-start;
    padding: 40rpx 36rpx 36rpx;
    z-index: 1;
  }

  &__arrow {
    align-self: center;
    margin-left: auto;

    .arrow-text {
      font-size: 48rpx;
      color: #c0c0d0;
      font-weight: 300;
    }
  }
}

/* Avatar */
.avatar {
  flex-shrink: 0;
  margin-right: 28rpx;

  &__img {
    width: 120rpx;
    height: 120rpx;
    border-radius: 50%;
    border: 6rpx solid #ffffff;
    box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.1);
  }

  &__placeholder {
    width: 120rpx;
    height: 120rpx;
    border-radius: 50%;
    background: linear-gradient(135deg, #007a49, #00a86b);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 6rpx solid #ffffff;
    box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.1);
  }

  &__initial {
    font-size: 48rpx;
    font-weight: 700;
    color: #ffffff;
  }
}

/* User info */
.user-info {
  display: flex;
  flex-direction: column;
  padding-top: 8rpx;

  &__name {
    font-size: 36rpx;
    font-weight: 700;
    color: #1a1a2e;
    line-height: 1.3;

    &--login {
      color: #007a49;
    }
  }

  &__id {
    font-size: 26rpx;
    color: #888898;
    margin-top: 6rpx;
    letter-spacing: 1rpx;
  }

  &__dept {
    font-size: 24rpx;
    color: #007a49;
    margin-top: 8rpx;
    background: rgba(0, 122, 73, 0.08);
    padding: 6rpx 16rpx;
    border-radius: 8rpx;
    align-self: flex-start;
  }
}

/* Section */
.section {
  margin: 32rpx 32rpx 0;

  &__title {
    font-size: 26rpx;
    font-weight: 600;
    color: #999aad;
    text-transform: uppercase;
    letter-spacing: 2rpx;
    margin-bottom: 16rpx;
    margin-left: 8rpx;
    display: block;
  }
}

/* Card */
.card {
  background: #ffffff;
  border-radius: 24rpx;
  overflow: hidden;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.04);
}

/* Cell */
.cell {
  display: flex;
  align-items: center;
  padding: 32rpx 32rpx;
  position: relative;

  &:not(.cell--last)::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 108rpx;
    right: 32rpx;
    height: 1rpx;
    background: #f0f0f5;
  }

  &:active {
    background: #fafafd;
  }

  &__icon {
    width: 72rpx;
    height: 72rpx;
    border-radius: 20rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 24rpx;
    flex-shrink: 0;

    &--sync {
      background: rgba(0, 122, 73, 0.1);
    }

    &--semester {
      background: rgba(33, 150, 243, 0.1);
    }

    &--reminder {
      background: rgba(255, 152, 0, 0.1);
    }

    &--about {
      background: rgba(156, 39, 176, 0.1);
    }
  }

  &__icon-text {
    font-size: 32rpx;
  }

  &__body {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  &__label {
    font-size: 30rpx;
    font-weight: 600;
    color: #1a1a2e;
    line-height: 1.4;
  }

  &__desc {
    font-size: 24rpx;
    color: #999aad;
    margin-top: 4rpx;
    line-height: 1.4;
  }

  &__arrow {
    font-size: 36rpx;
    color: #c8c8d8;
    margin-left: 16rpx;
    font-weight: 300;
    flex-shrink: 0;
  }
}

/* Download card */
.download-card {
  background: #ffffff;
  border-radius: 24rpx;
  padding: 40rpx 36rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.04);

  &__header {
    display: flex;
    align-items: center;
    gap: 16rpx;
  }

  &__title {
    font-size: 34rpx;
    font-weight: 700;
    color: #1a1a2e;
  }

  &__badge {
    background: linear-gradient(135deg, #007a49, #00a86b);
    padding: 4rpx 16rpx;
    border-radius: 8rpx;
  }

  &__badge-text {
    font-size: 20rpx;
    color: #ffffff;
    font-weight: 600;
  }

  &__subtitle {
    font-size: 24rpx;
    color: #999aad;
    margin-top: 8rpx;
    display: block;
  }

  &__features {
    margin-top: 28rpx;
    display: flex;
    flex-direction: column;
    gap: 20rpx;
  }

  &__btn {
    margin-top: 32rpx;
    width: 100%;
    height: 88rpx;
    border-radius: 44rpx;
    background: linear-gradient(135deg, #007a49, #00a86b);
    color: #ffffff;
    font-size: 30rpx;
    font-weight: 600;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    line-height: 88rpx;
    box-shadow: 0 8rpx 24rpx rgba(0, 122, 73, 0.25);

    &::after {
      border: none;
    }

    &:active {
      opacity: 0.9;
    }
  }
}

/* Feature items */
.feature-item {
  display: flex;
  align-items: center;

  &__dot {
    width: 12rpx;
    height: 12rpx;
    border-radius: 50%;
    background: #007a49;
    margin-right: 16rpx;
    flex-shrink: 0;
  }

  &__text {
    font-size: 28rpx;
    color: #444460;
    font-weight: 500;
  }
}

/* Logout */
.logout-section {
  margin: 48rpx 32rpx 0;
}

.logout-btn {
  width: 100%;
  height: 92rpx;
  border-radius: 46rpx;
  background: #ffffff;
  color: #e53935;
  font-size: 30rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid #fde8e8;
  line-height: 92rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);

  &::after {
    border: none;
  }

  &:active {
    background: #fff5f5;
  }
}

/* Bottom safe area spacing */
.safe-bottom {
  height: 60rpx;
}
</style>
