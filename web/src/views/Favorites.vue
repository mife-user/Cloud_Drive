<template>
  <div class="favorites-container">
    <Navbar />
    
    <main class="main-content">
      <div class="toolbar">
        <div class="toolbar-left">
          <h2>我的收藏</h2>
          <span class="file-count">{{ files.length }} 个文件</span>
        </div>
      </div>
      
      <div class="file-section">
        <div v-if="loading" class="loading-state">
          <el-icon class="is-loading" :size="48"><Loading /></el-icon>
          <p>加载中...</p>
        </div>
        
        <div v-else-if="files.length === 0" class="empty-state">
          <el-icon :size="64"><Star /></el-icon>
          <p>暂无收藏文件</p>
          <router-link to="/">
            <el-button type="primary">去收藏文件</el-button>
          </router-link>
        </div>
        
        <div v-else class="file-grid">
          <div
            v-for="file in files"
            :key="file.ID"
            class="file-card"
            @click="handleFileClick(file)"
          >
            <div class="file-preview">
              <img
                v-if="isImage(file.FileName)"
                :src="getFileUrl(file.ID)"
                alt="preview"
              />
              <el-icon v-else :size="48"><Document /></el-icon>
            </div>
            <div class="file-info">
              <p class="file-name" :title="file.FileName">{{ file.FileName }}</p>
              <p class="file-meta">
                <span>{{ formatSize(file.Size) }}</span>
                <span>{{ file.Owner }}</span>
              </p>
            </div>
            <div class="file-actions" @click.stop>
              <el-button
                circle
                size="small"
                type="warning"
                @click="removeFromFavorite(file.ID)"
              >
                <el-icon><StarFilled /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Loading, Star, Document, StarFilled } from '@element-plus/icons-vue'
import Navbar from '@/components/Navbar.vue'
import { getFavorites, removeFavorite } from '@/api/file'

const loading = ref(false)
const files = ref([])

function getFileUrl(fileId) {
  return `/api/file/view/${fileId}`
}

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

async function loadFiles() {
  loading.value = true
  try {
    const res = await getFavorites()
    files.value = res.files || []
  } catch (error) {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}

function handleFileClick(file) {
  window.open(getFileUrl(file.ID), '_blank')
}

async function removeFromFavorite(fileId) {
  try {
    await removeFavorite(fileId)
    ElMessage.success('已取消收藏')
    loadFiles()
  } catch (error) {
    // 错误已处理
  }
}

onMounted(() => {
  loadFiles()
})
</script>

<style scoped>
.favorites-container {
  min-height: 100vh;
  background: #f5f7fa;
}

.main-content {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.toolbar-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.toolbar-left h2 {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.file-count {
  font-size: 14px;
  color: #999;
}

.file-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  min-height: 400px;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  color: #999;
}

.empty-state p {
  margin: 16px 0;
  font-size: 16px;
}

.file-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.file-card {
  position: relative;
  background: #f9fafb;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid transparent;
}

.file-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
  border-color: #667eea;
}

.file-preview {
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ed 100%);
}

.file-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.file-preview .el-icon {
  color: #c0c4cc;
}

.file-info {
  padding: 12px;
}

.file-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin: 0 0 6px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
  margin: 0;
}

.file-actions {
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.file-card:hover .file-actions {
  opacity: 1;
}
</style>
