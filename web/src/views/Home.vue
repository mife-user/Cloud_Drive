<template>
  <div class="home-container">
    <Navbar />
    
    <main class="main-content">
      <div class="toolbar">
        <div class="toolbar-left">
          <h2>我的文件</h2>
          <span class="file-count">{{ files.length }} 个文件</span>
        </div>
        <div class="toolbar-right">
          <el-button type="primary" @click="uploadRef?.open()">
            <el-icon><Upload /></el-icon>
            上传文件
          </el-button>
        </div>
      </div>
      
      <div class="file-section">
        <div v-if="loading" class="loading-state">
          <el-icon class="is-loading" :size="48"><Loading /></el-icon>
          <p>加载中...</p>
        </div>
        
        <div v-else-if="files.length === 0" class="empty-state">
          <el-icon :size="64"><FolderOpened /></el-icon>
          <p>暂无文件</p>
          <el-button type="primary" @click="uploadRef?.open()">
            上传第一个文件
          </el-button>
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
                <span>{{ file.Permissions === 'public' ? '公开' : '私有' }}</span>
              </p>
            </div>
            <div class="file-actions" @click.stop>
              <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, file)">
                <el-button circle size="small">
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="download">
                      <el-icon><Download /></el-icon>
                      下载
                    </el-dropdown-item>
                    <el-dropdown-item command="share">
                      <el-icon><Share /></el-icon>
                      分享
                    </el-dropdown-item>
                    <el-dropdown-item command="favorite">
                      <el-icon><Star /></el-icon>
                      收藏
                    </el-dropdown-item>
                    <el-dropdown-item command="permission">
                      <el-icon><Lock /></el-icon>
                      修改权限
                    </el-dropdown-item>
                    <el-dropdown-item divided command="delete">
                      <el-icon><Delete /></el-icon>
                      删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </div>
      </div>
    </main>
    
    <FileUpload ref="uploadRef" @success="loadFiles" />
    <ShareDialog ref="shareRef" />
    
    <el-dialog v-model="permissionDialog.visible" title="修改权限" width="400px">
      <el-form>
        <el-form-item label="文件权限">
          <el-radio-group v-model="permissionDialog.permission">
            <el-radio value="public">公开</el-radio>
            <el-radio value="private">私有</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="permissionDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="updatePermission">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, Loading, FolderOpened, Document, MoreFilled, Download, Share, Star, Lock, Delete } from '@element-plus/icons-vue'
import Navbar from '@/components/Navbar.vue'
import FileUpload from '@/components/FileUpload.vue'
import ShareDialog from '@/components/ShareDialog.vue'
import { getFiles, deleteFile, updatePermissions, addFavorite } from '@/api/file'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const loading = ref(false)
const files = ref([])
const uploadRef = ref(null)
const shareRef = ref(null)

const permissionDialog = reactive({
  visible: false,
  fileId: null,
  permission: 'private'
})

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
    const res = await getFiles(authStore.userName)
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

async function handleCommand(command, file) {
  switch (command) {
    case 'download':
      window.open(getFileUrl(file.ID), '_blank')
      break
    case 'share':
      shareRef.value?.open(file.ID)
      break
    case 'favorite':
      try {
        await addFavorite(file.ID)
        ElMessage.success('收藏成功')
      } catch (error) {
        // 错误已处理
      }
      break
    case 'permission':
      permissionDialog.fileId = file.ID
      permissionDialog.permission = file.Permissions === 'public' ? 'public' : 'private'
      permissionDialog.visible = true
      break
    case 'delete':
      try {
        await ElMessageBox.confirm('确定要删除这个文件吗？', '删除确认', {
          type: 'warning'
        })
        await deleteFile(file.ID)
        ElMessage.success('文件已移入回收站')
        loadFiles()
      } catch (error) {
        if (error !== 'cancel') {
          // 错误已处理
        }
      }
      break
  }
}

async function updatePermission() {
  try {
    await updatePermissions(permissionDialog.fileId, permissionDialog.permission)
    ElMessage.success('权限更新成功')
    permissionDialog.visible = false
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
.home-container {
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
