<template>
  <el-container style="height: 100%;">
    <el-aside width="250px" style="background-color: var(--el-bg-color-overlay); border-right: 1px solid var(--el-border-color-light);">
      <el-menu
        default-active="1"
        @select="(index) => activeMenu = index"
        style="height: 100%; border-right: none; background-color: var(--el-bg-color-overlay);"
      >
        <el-menu-item index="1">
          <el-icon><Files /></el-icon>
          <span>{{ $t('config.menu.install') }}</span>
        </el-menu-item>
        <el-menu-item index="2">
          <el-icon><Picture /></el-icon>
          <span>{{ $t('config.menu.logo') }}</span>
        </el-menu-item>
        <el-menu-item index="3">
          <el-icon><Box /></el-icon>
          <span>{{ $t('config.menu.unpack') }}</span>
        </el-menu-item>
        <el-menu-item index="4">
          <el-icon><Document /></el-icon>
          <span>{{ $t('config.menu.script') }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-main style="padding: 20px;">
      <div v-if="activeMenu === '1'">
        <h3>{{ $t('config.title.install') }}</h3>
        <p>{{ $t('config.description.install') }}</p>
      </div>
      <div v-if="activeMenu === '2'">
        <h3>{{ $t('config.title.logo') }}</h3>
        <p>{{ $t('config.description.logo') }}</p>
      </div>
      <div v-if="activeMenu === '3'">
        <h3>{{ $t('config.title.unpack') }}</h3>
        <p>{{ $t('config.description.unpack') }}</p>
      </div>
      <div v-if="activeMenu === '4'">
        <h3>{{ $t('config.title.script') }}</h3>
        <p>{{ $t('config.description.script') }}</p>
      </div>
    </el-main>
  </el-container>
  <div class="footer-buttons">
    <el-button @click="prev">{{ $t('config.button.prev') }}</el-button>
    <el-button type="primary" @click="next">{{ $t('config.button.next') }}</el-button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Picture, Files, Box, Document } from '@element-plus/icons-vue'

const router = useRouter()
const activeMenu = ref('1')

const next = () => {
  router.push('/build')
}

const prev = () => {
  router.push('/')
}
</script>

<style scoped>
/* Scoped styles for Configuration.vue */
.el-main h3 {
  margin-top: 0;
  color: var(--el-text-color-primary);
}
.el-main p {
  color: var(--el-text-color-regular);
}
</style>