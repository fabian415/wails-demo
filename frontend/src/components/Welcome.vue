<template>
  <div class="welcome-container">
    <el-icon :size="80" color="#409eff"><Cpu /></el-icon>
    <h1 class="title">{{ $t('welcome.title') }}</h1>
    <p class="subtitle">{{ $t('welcome.subtitle') }}</p>
    <el-button type="primary" size="large" @click="start">{{ $t('welcome.button') }}</el-button>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { Cpu } from '@element-plus/icons-vue'

const router = useRouter()

const start = () => {
  router.push('/config')
}
</script>

<style scoped>
.welcome-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
  text-align: center;
}
.title {
  font-size: 3rem;
  font-weight: 500;
  margin: 24px 0 12px;
  color: var(--el-text-color-primary);
}
.subtitle {
  font-size: 1.2rem;
  color: var(--el-text-color-secondary);
  margin-bottom: 40px;
}
</style>