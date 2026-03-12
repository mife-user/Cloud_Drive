<template>
  <el-dialog
    v-model="visible"
    title="分享文件"
    width="450px"
  >
    <div v-if="shareInfo" class="share-content">
      <el-alert
        type="success"
        :closable="false"
        show-icon
        style="margin-bottom: 20px;"
      >
        文件分享成功
      </el-alert>
      
      <el-form label-width="80px">
        <el-form-item label="分享链接">
          <el-input
            v-model="shareInfo.share_url"
            readonly
          >
            <template #append>
              <el-button @click="copyToClipboard(shareInfo.share_url)">
                复制
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item label="访问密钥">
          <el-input
            v-model="shareInfo.access_key"
            readonly
          >
            <template #append>
              <el-button @click="copyToClipboard(shareInfo.access_key)">
                复制
              </el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      
      <el-alert
        type="info"
        :closable="false"
        style="margin-top: 16px;"
      >
        请将链接和密钥发送给需要访问的人
      </el-alert>
    </div>
    
    <template #footer>
      <el-button type="primary" @click="close">完成</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { createShare } from '@/api/file'

const visible = ref(false)
const shareInfo = ref(null)

async function open(fileId) {
  try {
    const res = await createShare(fileId)
    shareInfo.value = {
      share_url: res.share_url,
      access_key: res.access_key
    }
    visible.value = true
  } catch (error) {
    // 错误已在拦截器中处理
  }
}

function close() {
  visible.value = false
  shareInfo.value = null
}

function copyToClipboard(text) {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

defineExpose({ open, close })
</script>

<style scoped>
.share-content {
  padding: 10px 0;
}
</style>
