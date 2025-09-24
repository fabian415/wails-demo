<template>
  <el-container style="height: 100%;">
    <el-main style="padding: 20px;">
      <div class="config-grid">
        <el-card v-for="(option, index) in configOptions" :key="index" class="config-card">
          <template #header>
            <div class="card-header">
              <el-checkbox v-model="option.enabled" :label="option.title" size="large" />
              <el-tooltip :content="option.description" placement="top">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <div class="card-body">
            <el-input v-model="option.path" :placeholder="$t('config.placeholder.path')" readonly>
              <template #append>
                <el-button @click="clearPath(option)" :disabled="!option.enabled"><el-icon><Close /></el-icon></el-button>
              </template>
            </el-input>
            <el-button type="primary" @click="selectPath(option)" :disabled="!option.enabled" class="select-button">
              {{ $t('config.button.select') }}
            </el-button>
          </div>
        </el-card>
      </div>
    </el-main>
  </el-container>
  <div class="footer-buttons">
    <el-button @click="prev">{{ $t('config.button.prev') }}</el-button>
    <div>
      <el-button @click="loadConfig">{{ $t('config.button.load') }}</el-button>
      <el-button type="success" @click="saveConfig">{{ $t('config.button.save') }}</el-button>
      <el-button type="primary" @click="next">{{ $t('config.button.next') }}</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { InfoFilled, Close } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { SaveConfig, LoadConfig, SelectFile, SelectDirectory } from '../../wailsjs/go/main/App'
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const router = useRouter()

const configOptions = ref([
  { id: 'ISO', type: 'file', enabled: false, path: '', title: t('config.options.iso.title'), description: t('config.options.iso.description'), patterns: ['*.*'] },
  { id: 'DEB', type: 'folder', enabled: false, path: '', title: t('config.options.deb.title'), description: t('config.options.deb.description') },
  { id: 'LOGO', type: 'file', enabled: false, path: '', title: t('config.options.logo.title'), description: t('config.options.logo.description'), patterns: ['*.png', '*.jpg', '*.svg'] },
  { id: 'DRIVER', type: 'folder', enabled: false, path: '', title: t('config.options.driver.title'), description: t('config.options.driver.description') },
  { id: 'PRE_SCRIPT', type: 'file', enabled: false, path: '', title: t('config.options.pre_script.title'), description: t('config.options.pre_script.description'), patterns: ['*.sh', '*.bat'] },
  { id: 'AUTOSTART', type: 'folder', enabled: false, path: '', title: t('config.options.autostart.title'), description: t('config.options.autostart.description') },
  { id: 'OTHER', type: 'folder', enabled: false, path: '', title: t('config.options.other.title'), description: t('config.options.other.description') },
  { id: 'BOOT_DEVICE', type: 'folder', enabled: false, path: '', title: t('config.options.boot_device.title'), description: t('config.options.boot_device.description') },
]);


const selectPath = async (option) => {
  try {
    let result;
    if (option.type === 'file') {
      const options = {
        title: t('config.dialog.select_file'),
        patterns: option.patterns,
      };
      result = await SelectFile(JSON.stringify(options));
    } else {
      result = await SelectDirectory(t('config.dialog.select_folder'));
    }
    if (result) {
      option.path = result;
    }
  } catch (error) {
    console.error("Error selecting path:", error);
    ElMessage.error(t('config.message.select_error'));
  }
};

const clearPath = (option) => {
  option.path = '';
};

const saveConfig = async () => {
  const configToSave = {};
  configOptions.value.forEach(opt => {
    if (opt.enabled && opt.path) {
      configToSave[opt.id] = opt.path;
    }
  });

  try {
    await SaveConfig(JSON.stringify(configToSave));
    ElMessage.success(t('config.message.save_success'));
  } catch (error) {
    console.error("Error saving config:", error);
    ElMessage.error(t('config.message.save_error'));
  }
};

const loadConfig = async () => {
  try {
    const configStr = await LoadConfig();
    if (configStr) {
      const savedConfig = JSON.parse(configStr);
      configOptions.value.forEach(opt => {
        if (savedConfig[opt.id]) {
          opt.path = savedConfig[opt.id];
          opt.enabled = true;
        } else {
          opt.path = '';
          opt.enabled = false;
        }
      });
      ElMessage.success(t('config.message.load_success'));
    }
  } catch (error) {
    console.error("Error loading config:", error);
    ElMessage.error(t('config.message.load_error'));
  }
};

const next = () => {
  router.push('/build');
};

const prev = () => {
  router.push('/');
};

onMounted(loadConfig);

</script>

<style scoped>
.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 20px;
}

.config-card {
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.select-button {
  width: 100%;
}

.footer-buttons {
  display: flex;
  justify-content: space-between;
  padding: 0 20px 20px;
  align-items: center;
}
</style>
