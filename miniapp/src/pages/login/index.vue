<template>
  <view class="login-page">
    <!-- Background decorative elements -->
    <view class="bg-decoration">
      <view class="bg-circle bg-circle--1" />
      <view class="bg-circle bg-circle--2" />
      <view class="bg-circle bg-circle--3" />
    </view>

    <!-- Brand area -->
    <view class="brand">
      <view class="brand__logo">
        <text class="brand__logo-icon">S</text>
      </view>
      <text class="brand__title">SCAU Daily</text>
      <text class="brand__subtitle">华南农业大学校园助手</text>
    </view>

    <!-- Login form -->
    <view class="form">
      <view class="form__card">
        <!-- Student ID -->
        <view class="form__field">
          <view class="form__field-icon">
            <text class="icon-text">ID</text>
          </view>
          <input
            v-model="studentId"
            class="form__input"
            type="number"
            placeholder="请输入学号"
            placeholder-class="form__placeholder"
            :maxlength="12"
          />
        </view>

        <!-- Password -->
        <view class="form__field">
          <view class="form__field-icon">
            <text class="icon-text">PW</text>
          </view>
          <input
            v-model="password"
            class="form__input"
            :password="!showPassword"
            placeholder="请输入密码"
            placeholder-class="form__placeholder"
            :maxlength="32"
          />
          <view class="form__field-toggle" @tap="showPassword = !showPassword">
            <text class="toggle-text">{{ showPassword ? '隐藏' : '显示' }}</text>
          </view>
        </view>

        <!-- Login button -->
        <button
          class="btn btn--primary"
          :loading="loading"
          :disabled="loading || !studentId || !password"
          @tap="handleLogin"
        >
          登录
        </button>

        <!-- Divider -->
        <view class="divider">
          <view class="divider__line" />
          <text class="divider__text">或</text>
          <view class="divider__line" />
        </view>

        <!-- WeChat login -->
        <button class="btn btn--wechat" @tap="handleWxLogin">
          <text class="btn__wechat-icon">W</text>
          <text>微信一键登录</text>
        </button>
      </view>
    </view>

    <!-- Bind dialog overlay -->
    <view v-if="showBind" class="bind-overlay" @tap.self="showBind = false">
      <view class="bind-dialog">
        <text class="bind-dialog__title">绑定学号</text>
        <text class="bind-dialog__desc">首次微信登录，请绑定您的学号</text>
        <view class="form__field" style="margin-top: 32rpx;">
          <view class="form__field-icon">
            <text class="icon-text">ID</text>
          </view>
          <input
            v-model="bindStudentId"
            class="form__input"
            type="number"
            placeholder="请输入学号"
            placeholder-class="form__placeholder"
            :maxlength="12"
          />
        </view>
        <view class="form__field">
          <view class="form__field-icon">
            <text class="icon-text">PW</text>
          </view>
          <input
            v-model="bindPassword"
            class="form__input"
            password
            placeholder="请输入教务系统密码"
            placeholder-class="form__placeholder"
            :maxlength="32"
          />
        </view>
        <button class="btn btn--primary" :loading="bindLoading" @tap="handleBind">
          确认绑定
        </button>
      </view>
    </view>

    <!-- Footer -->
    <view class="footer">
      <text class="footer__text">登录即表示同意</text>
      <text class="footer__link">用户协议</text>
      <text class="footer__text">与</text>
      <text class="footer__link">隐私政策</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/store/user'
import { login, wxLogin, bindStudent } from '@/api'

const userStore = useUserStore()

// Form fields
const studentId = ref('')
const password = ref('')
const showPassword = ref(false)
const loading = ref(false)

// Bind dialog
const showBind = ref(false)
const bindStudentId = ref('')
const bindPassword = ref('')
const bindLoading = ref(false)
const wxCodeCache = ref('')

/**
 * Student ID + password login
 */
async function handleLogin() {
  if (!studentId.value || !password.value) {
    uni.showToast({ title: '请输入学号和密码', icon: 'none' })
    return
  }

  loading.value = true
  uni.showLoading({ title: '登录中...' })

  try {
    const resp = await login(studentId.value, password.value)
    userStore.setAuth(resp)
    uni.showToast({ title: '登录成功', icon: 'success' })
    setTimeout(() => {
      uni.switchTab({ url: '/pages/today/index' })
    }, 500)
  } catch (err: any) {
    uni.showToast({
      title: err?.message || '登录失败，请检查学号和密码',
      icon: 'none'
    })
  } finally {
    loading.value = false
    uni.hideLoading()
  }
}

/**
 * WeChat one-click login
 */
async function handleWxLogin() {
  loading.value = true
  uni.showLoading({ title: '正在获取微信授权...' })

  try {
    // Step 1: Get wx login code
    const loginRes = await new Promise<UniApp.LoginRes>((resolve, reject) => {
      uni.login({
        provider: 'weixin',
        success: (res) => resolve(res),
        fail: (err) => reject(new Error('微信授权失败，请重试')),
      })
    })

    const code = loginRes.code
    wxCodeCache.value = code

    uni.showLoading({ title: '正在登录...' })

    // Step 2: Call backend wx login
    const resp = await wxLogin(code)

    if (resp.need_bind) {
      // Need to bind student ID
      showBind.value = true
    } else if (resp.data) {
      userStore.setAuth(resp.data)
      uni.showToast({ title: '登录成功', icon: 'success' })
      setTimeout(() => {
        uni.switchTab({ url: '/pages/today/index' })
      }, 500)
    }
  } catch (err: any) {
    uni.showToast({
      title: err?.message || '微信登录失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
    uni.hideLoading()
  }
}

/**
 * Bind student account after WeChat login
 */
async function handleBind() {
  if (!bindStudentId.value || !bindPassword.value) {
    uni.showToast({ title: '请输入学号和密码', icon: 'none' })
    return
  }

  bindLoading.value = true
  uni.showLoading({ title: '绑定中...' })

  try {
    const resp = await bindStudent(
      bindStudentId.value,
      bindPassword.value
    )
    userStore.setAuth(resp)
    showBind.value = false
    uni.showToast({ title: '绑定成功', icon: 'success' })
    setTimeout(() => {
      uni.switchTab({ url: '/pages/today/index' })
    }, 500)
  } catch (err: any) {
    uni.showToast({
      title: err?.message || '绑定失败，请检查信息',
      icon: 'none'
    })
  } finally {
    bindLoading.value = false
    uni.hideLoading()
  }
}
</script>

<style lang="scss" scoped>
@use '../../styles/variables' as *;

.login-page {
  position: relative;
  min-height: 100vh;
  background: linear-gradient(160deg, #f0f9f4 0%, #ffffff 40%, #f5f6fa 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 60rpx;
  overflow: hidden;
}

/* Background decorative circles */
.bg-decoration {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  overflow: hidden;
}

.bg-circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.12;

  &--1 {
    width: 500rpx;
    height: 500rpx;
    background: #007a49;
    top: -120rpx;
    right: -100rpx;
  }

  &--2 {
    width: 300rpx;
    height: 300rpx;
    background: #00a86b;
    top: 200rpx;
    left: -80rpx;
  }

  &--3 {
    width: 200rpx;
    height: 200rpx;
    background: #007a49;
    bottom: 300rpx;
    right: -40rpx;
  }
}

/* Brand area */
.brand {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 180rpx;
  margin-bottom: 64rpx;
  z-index: 1;

  &__logo {
    width: 120rpx;
    height: 120rpx;
    border-radius: 32rpx;
    background: linear-gradient(135deg, #007a49, #00a86b);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 28rpx;
    box-shadow: 0 16rpx 40rpx rgba(0, 122, 73, 0.25);
  }

  &__logo-icon {
    font-size: 56rpx;
    font-weight: 800;
    color: #ffffff;
    font-family: 'Georgia', serif;
  }

  &__title {
    font-size: 56rpx;
    font-weight: 800;
    color: #1a1a2e;
    letter-spacing: 2rpx;
  }

  &__subtitle {
    font-size: 28rpx;
    color: #666680;
    margin-top: 12rpx;
    letter-spacing: 4rpx;
  }
}

/* Form area */
.form {
  width: 100%;
  z-index: 1;

  &__card {
    background: #ffffff;
    border-radius: 32rpx;
    padding: 48rpx 40rpx;
    box-shadow: 0 8rpx 40rpx rgba(0, 0, 0, 0.06);
  }

  &__field {
    display: flex;
    align-items: center;
    height: 104rpx;
    background: #f5f6fa;
    border-radius: 20rpx;
    padding: 0 28rpx;
    margin-bottom: 24rpx;
    border: 2rpx solid transparent;
    transition: border-color 0.2s;

    &:focus-within {
      border-color: #007a49;
      background: #f8faf9;
    }
  }

  &__field-icon {
    width: 56rpx;
    height: 56rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 20rpx;
    flex-shrink: 0;

    .icon-text {
      font-size: 24rpx;
      font-weight: 700;
      color: #007a49;
      background: rgba(0, 122, 73, 0.1);
      padding: 8rpx 12rpx;
      border-radius: 8rpx;
    }
  }

  &__input {
    flex: 1;
    height: 104rpx;
    font-size: 30rpx;
    color: #1a1a2e;
  }

  &__placeholder {
    color: #b0b0c0;
    font-size: 28rpx;
  }

  &__field-toggle {
    padding: 12rpx 16rpx;
    flex-shrink: 0;

    .toggle-text {
      font-size: 24rpx;
      color: #007a49;
      font-weight: 500;
    }
  }
}

/* Buttons */
.btn {
  width: 100%;
  height: 96rpx;
  border-radius: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  font-weight: 600;
  border: none;
  margin: 0;
  padding: 0;
  line-height: 96rpx;

  &::after {
    border: none;
  }

  &--primary {
    background: linear-gradient(135deg, #007a49, #00a86b);
    color: #ffffff;
    margin-top: 16rpx;
    box-shadow: 0 12rpx 32rpx rgba(0, 122, 73, 0.3);

    &[disabled] {
      opacity: 0.5;
      box-shadow: none;
    }
  }

  &--wechat {
    background: #f5f6fa;
    color: #333345;
    font-weight: 500;
    gap: 12rpx;

    &:active {
      background: #ecedf2;
    }
  }

  &__wechat-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 44rpx;
    height: 44rpx;
    background: #07c160;
    color: #ffffff;
    font-size: 26rpx;
    font-weight: 800;
    border-radius: 12rpx;
    font-family: 'Arial', sans-serif;
  }
}

/* Divider */
.divider {
  display: flex;
  align-items: center;
  margin: 40rpx 0;

  &__line {
    flex: 1;
    height: 1rpx;
    background: #e8e8f0;
  }

  &__text {
    padding: 0 28rpx;
    font-size: 26rpx;
    color: #999aad;
  }
}

/* Bind dialog overlay */
.bind-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
  padding: 60rpx;
}

.bind-dialog {
  background: #ffffff;
  border-radius: 32rpx;
  padding: 56rpx 44rpx;
  width: 100%;

  &__title {
    font-size: 36rpx;
    font-weight: 700;
    color: #1a1a2e;
    display: block;
    text-align: center;
  }

  &__desc {
    font-size: 26rpx;
    color: #666680;
    display: block;
    text-align: center;
    margin-top: 12rpx;
  }

  .btn--primary {
    margin-top: 40rpx;
  }
}

/* Footer */
.footer {
  position: fixed;
  bottom: 60rpx;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;

  &__text {
    font-size: 22rpx;
    color: #999aad;
  }

  &__link {
    font-size: 22rpx;
    color: #007a49;
    margin: 0 4rpx;
  }
}
</style>
