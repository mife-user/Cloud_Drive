<template>
  <div class="share-container">
    <div class="share-card">
      <div class="share-header">
        <el-icon class="share-icon"><Share /></el-icon>
        <h1>文件分享</h1>
        <p>请输入访问密钥查看分享的文件</p>
      </div>
      
      <div v-if="!fileData" class="share-form">
        <el-input
          v-model="accessKey"
          placeholder="请输入访问密钥"
          size="large"
          @keyup.enter="accessShareFile"
        >
          <template #prefix>
            <el-icon><Key /></el-icon>
          </template>
        </el-input>
        
        <el-button
          type="primary"
          size="large"
          :loading="loading"
          @click="accessShareFile"
        >
          访问文件
        </el-button>
      </div>
      
      <div v-else class="share-result">
        <div class="file-preview-large">
          <img
            v-if="isImage(fileData.FileName)"
            :src="fileData.url"
            alt="preview"
          />
          <el-icon v-else :size="80"><Document /></el-icon>
        </div>
        
        <div class="file-details">
          <h3>{{ fileData.FileName }}</h3>
          <p>文件大小: {{ formatSize(fileData.Size) }}</p>
          <p>分享者: {{ fileData.Owner }}</p>
        </div>
        
        <div class="file-actions">
          <el-button type="primary" size="large" @click="downloadFile">
            <el-icon><Download /></el-icon>
            下载文件
          </el-button>
          
          <el-button
            v-if="authStore.isLoggedIn"
            size="large"
            @click="addToFavorite"
          >
            <el-icon><Star /></el-icon>
            收藏到我的云盘
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Share, Key, Document, Download, Star } from '@element-plus/icons-vue'
import { accessShare } from '@/api/file'
import { addFavorite } from '@/api/file'
import { useAuthStore } from '@/stores/auth'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const accessKey = ref('')
const loading = ref(false)
const fileData = ref(null)
const shareId = ref('')

function isImage(filename) {
  const ext = filename?.toLowerCase().split('.').pop()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg'].includes(ext)
}

function formatSize(bytes) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  return (bytes / (1024 * 1024 * 1024)).toFixed(1) + ' GB'
}

async function accessShareFile() {
  if (!accessKey.value.trim()) {
    ElMessage.warning('请输入访问密钥')
    return
  }
  
  loading.value = true
  try {
    // 直接访问分享API获取文件
    const response = await axios.get(`/api/file/share/${shareId.value}`, {
      params: { access_key: accessKey.value }
    })
    
    fileData.value = {
      ...response.data,
      url: `/api/file/share/${shareId.value}?access_key=${accessKey.value}`
    }
  } catch (error) {
    ElMessage.error('访问密钥错误或分享已过期')
  } finally {
    loading.value = false
  }
}

function downloadFile() {
  window.open(fileData.value.url, '_blank')
}

async function addToFavorite() {
  try {
    await addFavorite(fileData.value.ID, accessKey.value)
    ElMessage.success('收藏成功')
  } catch (error) {
    // 错误已处理
  }
}

onMounted(() => {
  shareId.value = route.params.shareId
  if (!shareId.value) {
    ElMessage.error('无效的分享链接')
    router.push('/login')
  }
})
</script>

<style scoped>
.share-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.share-card {
  width: 100%;
  max-width: 500px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(10px);
}

.share-header {
  text-align: center;
  margin-bottom: 32px;
}

.share-icon {
  font-size: 48px;
  color: #667eea;
  margin-bottom: 16px;
}

.share-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
}

.share-header p {
  color: #666;
  font-size: 14px;
  margin: 0;
}

.share-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.share-result {
  text-align: center;
}

.file-preview-large {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  border-radius: 12px;
  margin-bottom: 24px;
  overflow: hidden;
}

.file-preview-large img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.file-preview-large .el-icon {
  color: #c0c4cc;
}

.file-details {
  margin-bottom: 24px;
}

.file-details h3 {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 12px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-details p {
  font-size: 14px;
  color: #666;
  margin: 4px 0;
}

.file-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
}
</style>
