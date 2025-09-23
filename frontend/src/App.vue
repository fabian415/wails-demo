<template>
  <el-container style="height: 100%;">
    <el-header style="display: flex; justify-content: space-between; align-items: center; background-color: var(--el-bg-color-overlay); color: var(--el-text-color-primary);">
      <div style="display: flex; align-items: center;">
        <img src="./assets/images/logo-universal.png" alt="logo" style="height: 32px; margin-right: 12px;"/>
        <span style="font-size: 20px; font-weight: bold;">OSBuilder</span>
      </div>
      <div>
        <el-switch
          v-model="isDark"
          inline-prompt
          active-text="Dark"
          inactive-text="Light"
          @change="toggleDark"
        />
        <el-dropdown @command="handleLanguageChange" style="margin-left: 20px;">
          <span class="el-dropdown-link" style="color: var(--el-text-color-primary); cursor: pointer;">
            {{ currentLanguage }}<el-icon class="el-icon--right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="en">English</el-dropdown-item>
              <el-dropdown-item command="zh-CN">简体中文</el-dropdown-item>
              <el-dropdown-item command="zh-TW">繁體中文</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    <el-main class="main-container">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { computed, ref } from 'vue'
import { ArrowDown } from '@element-plus/icons-vue'
import { useDark, useToggle } from '@vueuse/core'

const isDark = useDark()
const toggleDark = useToggle(isDark)

const { locale } = useI18n()

const languageMap = {
  'en': 'English',
  'zh-CN': '简体中文',
  'zh-TW': '繁體中文'
}

const currentLanguage = computed(() => languageMap[locale.value])

const handleLanguageChange = (lang) => {
  locale.value = lang
}
</script>

<style>
/* Global styles are in style.css */
.el-header {
  --el-header-padding: 0 20px;
}
</style>