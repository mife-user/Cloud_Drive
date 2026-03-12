<template>
  <div class="trash-container">
    <Navbar />
    
    <main class="main-content">
      <div class="toolbar">
        <div class="toolbar-left">
          <h2>回收站</h2>
          <span class="file-count">{{ files.length }} 个文件</span>
        </div>
        <div class="toolbar-right">
          <el-alert
            v-if="files.length > 0"
            type="info"
            :closable="false"
            show-icon
          >
            文件将在回收站保存1小时后自动删除
          </el-alert>
        </div>
      </div>
      
      <div class="file-section">
        <div v-if="loading" class="loading-state">
          <el-icon class="is-loading" :size="48"><Loading /></el-icon>
          <p>加载中...</p>
        </div>
        
        <div v-else-if="files.length === 0" class="empty-state">
          <el-icon :size="64"><Delete /></el-icon>
          <p>回收站为空</p>
        </div>
        
        <div v-else class="file-grid">
          <div
            v-for="file in files"
            :key="file.ID"
            class="file-card deleted"
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
              </p>
            </div>
            <div class="file-actions">
              <el-button
                circle
                size="small"
                type="danger"
                @click="deleteForever(file.ID)"
              >
                <el-icon><DeleteFilled /></el-icon>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { Loading, Delete, Document, DeleteFilled } from '@element-plus/icons-vue'
import Navbar from '@/components/Navbar.vue'
import { getDeletedFiles, deleteFileForever } from '@/api/file'

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
    const res = await getDeletedFiles()
    files.value = res.files || []
  } catch (error) {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}

async function deleteForever(fileId) {
  try {
    await ElMessageBox.confirm('确定要永久删除这个文件吗？此操作不可恢复', '永久删除', {
      type: 'warning',
      confirmButtonText: '永久删除',
      confirmButtonClass: 'el-button--danger'
    })
    await deleteFileForever(fileId)
    ElMessage.success('文件已永久删除')
    loadFiles()
  } catch (error) {
    if (error !== 'cancel') {
      // 错误已处理
    }
  }
}

onMounted(() => {
  loadFiles()
})
</script>

<style scoped>
.trash-container {
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
  transition: all 0.3s ease;
  border: 1px solid transparent;
}

.file-card.deleted {
  opacity: 0.7;
}

.file-card:hover {
  opacity: 1;
  border-color: #f56c6c;
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
  filter: grayscale(50%);
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
