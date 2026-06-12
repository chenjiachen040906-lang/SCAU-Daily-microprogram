<script setup lang="ts">
import { onLaunch, onShow } from '@dcloudio/uni-app'
import { useUserStore } from './store/user'

onLaunch(() => {
  console.log('[App] launched')
  const userStore = useUserStore()
  userStore.checkAuth()
})

onShow(() => {
  console.log('[App] visible')
})
</script>

<style lang="scss">
@import './styles/variables.scss';

page {
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Helvetica Neue',
    'Microsoft YaHei', sans-serif;
  background-color: $bg-color;
  color: $text-primary;
  font-size: 28rpx;
  line-height: 1.6;
}

/* Global utility classes */
.container {
  padding: 24rpx;
}

.card {
  background: #fff;
  border-radius: 20rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.text-primary {
  color: $brand-green;
}

.text-secondary {
  color: $text-secondary;
}

.text-hint {
  color: $text-hint;
  font-size: 24rpx;
}
</style>
