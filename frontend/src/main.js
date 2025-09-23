import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import './style.css'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import { createI18n } from 'vue-i18n'
import en from './locales/en.json'
import zhCN from './locales/zh-CN.json'
import zhTW from './locales/zh-TW.json'

const i18n = createI18n({
  locale: 'zh-TW', // set locale
  fallbackLocale: 'en', // set fallback locale
  messages: {
    en,
    'zh-CN': zhCN,
    'zh-TW': zhTW,
  },
})

const app = createApp(App)
app.use(router)
.use(ElementPlus)
.use(i18n)
.mount('#app')
