<template>
  <el-dialog
    v-model="visible"
    title="上传文件"
    width="500px"
    :close-on-click-modal="false"
  >
    <el-upload
      ref="uploadRef"
      class="upload-area"
      drag
      :auto-upload="false"
      :on-change="handleFileChange"
      :file-list="fileList"
      :limit="10"
      multiple
    >
      <el-icon class="upload-icon"><UploadFilled /></el-icon>
      <div class="upload-text">
        拖拽文件到此处，或<em>点击上传</em>
      </div>
      <template #tip>
        <div class="upload-tip">
          支持上传 jpg, png, gif, pdf, doc, docx, xls, xlsx, ppt, pptx, txt, zip, rar 格式文件
        </div>
      </template>
    </el-upload>
    
    <el-form class="permission-form">
      <el-form-item label="文件权限">
        <el-radio-group v-model="permissions">
          <el-radio value="public">公开</el-radio>
          <el-radio value="private">私有</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>
    
    <template #footer>
      <el-button @click="close">取消</el-button>
      <el-button type="primary" :loading="uploading" @click="handleUpload">
        上传
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import { uploadFile } from '@/api/file'

const emit = defineEmits(['success'])

const visible = ref(false)
const uploadRef = ref(null)
const fileList = ref([])
const permissions = ref('private')
const uploading = ref(false)

function open() {
  visible.value = true
  fileList.value = []
  permissions.value = 'private'
}

function close() {
  visible.value = false
  fileList.value = []
}

function handleFileChange(file, list) {
  fileList.value = list
}

async function handleUpload() {
  if (fileList.value.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }
  
  uploading.value = true
  
  try {
    for (const file of fileList.value) {
      const formData = new FormData()
      formData.append('files', file.raw)
      formData.append('permissions', permissions.value)
      
      await uploadFile(formData, (progressEvent) => {
        const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        file.percentage = percent
      })
    }
    
    ElMessage.success('上传成功')
    emit('success')
    close()
  } catch (error) {
    // 错误已在拦截器中处理
  } finally {
    uploading.value = false
  }
}

defineExpose({ open, close })
</script>

<style scoped>
.upload-area {
  width: 100%;
}

.upload-area :deep(.el-upload-dragger) {
  width: 100%;
  height: 180px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  border: 2px dashed #d9d9d9;
  transition: all 0.3s ease;
}

.upload-area :deep(.el-upload-dragger:hover) {
  border-color: #667eea;
}

.upload-icon {
  font-size: 48px;
  color: #c0c4cc;
  margin-bottom: 16px;
}

.upload-text {
  color: #666;
  font-size: 14px;
}

.upload-text em {
  color: #667eea;
  font-style: normal;
}

.upload-tip {
  text-align: center;
  color: #999;
  font-size: 12px;
  margin-top: 12px;
}

.permission-form {
  margin-top: 20px;
}
</style>
