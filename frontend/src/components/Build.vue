<template>
  <div class="build-container">
    <div style="width: 70%;">
      <h2>正在建立映像...</h2>
      <el-progress :percentage="progress" :text-inside="true" :stroke-width="24" status="success" class="progress-bar"/>
      <p style="text-align: center; color: var(--el-text-color-secondary);">{{ statusText }}</p>
    </div>

    <el-input
      type="textarea"
      :rows="12"
      placeholder=""
      v-model="logs"
      readonly
      class="log-window"
    />

    <el-dialog v-model="showResultDialog" :title="resultTitle" width="30%" center>
      <span>{{ resultMessage }}</span>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="startOver">完成</el-button>
        </span>
      </template>
    </el-dialog>

  </div>
  <div class="footer-buttons">
      <el-button @click="prev" :disabled="progress > 0 && progress < 100">上一步</el-button>
      <el-button type="primary" @click="startOver" v-if="progress === 100">重新開始</el-button>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const progress = ref(0)
const logs = ref('')
const showResultDialog = ref(false)
const resultTitle = ref('')
const resultMessage = ref('')

const statusText = computed(() => {
  if (progress.value < 100) {
    return `目前步驟: ${Math.floor(progress.value / 10) + 1} / 10`
  } else {
    return "建立完成！"
  }
})

onMounted(() => {
  const interval = setInterval(() => {
    if (progress.value < 100) {
      progress.value += 10
      logs.value += `[INFO] ${new Date().toLocaleTimeString()}: Build step ${progress.value / 10} completed.\n`
    } else {
      clearInterval(interval)
      logs.value += `[SUCCESS] ${new Date().toLocaleTimeString()}: Build completed successfully!\n`
      resultTitle.value = "成功"
      resultMessage.value = "您的客製化作業系統映像已成功建立。"
      showResultDialog.value = true
    }
  }, 300)
})

const startOver = () => {
  showResultDialog.value = false
  router.push('/')
}

const prev = () => {
  router.push('/config')
}
</script>

<style scoped>
.build-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
}
.progress-bar {
  margin: 20px 0;
}
.log-window {
  width: 70%;
  margin-top: 20px;
  font-family: monospace;
}
</style>