<template>
  <header class="navbar">
    <div class="navbar-brand">
      <el-icon class="brand-icon"><FolderOpened /></el-icon>
      <span class="brand-text">云盘</span>
    </div>
    
    <nav class="navbar-menu">
      <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }">
        <el-icon><Document /></el-icon>
        <span>我的文件</span>
      </router-link>
      <router-link to="/favorites" class="nav-item" :class="{ active: $route.path === '/favorites' }">
        <el-icon><Star /></el-icon>
        <span>收藏</span>
      </router-link>
      <router-link to="/trash" class="nav-item" :class="{ active: $route.path === '/trash' }">
        <el-icon><Delete /></el-icon>
        <span>回收站</span>
      </router-link>
    </nav>
    
    <div class="navbar-user">
      <el-dropdown trigger="click" @command="handleCommand">
        <div class="user-info">
          <el-avatar :size="36" :src="avatarUrl">
            {{ userName?.charAt(0)?.toUpperCase() }}
          </el-avatar>
          <span class="user-name">{{ userName }}</span>
          <el-icon><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="avatar">
              <el-icon><User /></el-icon>
              更换头像
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      
      <input
        ref="avatarInput"
        type="file"
        accept="image/*"
        style="display: none"
        @change="handleAvatarChange"
      />
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { FolderOpened, Document, Star, Delete, ArrowDown, User, SwitchButton } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { updateHeader } from '@/api/user'

const router = useRouter()
const authStore = useAuthStore()

const avatarInput = ref(null)

const userName = computed(() => authStore.userName)
const avatarUrl = computed(() => {
  if (authStore.user?.header_path) {
    return `/api${authStore.user.header_path}`
  }
  return ''
})

function handleCommand(command) {
  if (command === 'logout') {
    authStore.logout()
    router.push('/login')
  } else if (command === 'avatar') {
    avatarInput.value?.click()
  }
}

async function handleAvatarChange(event) {
  const file = event.target.files[0]
  if (!file) return
  
  const formData = new FormData()
  formData.append('header', file)
  
  try {
    await updateHeader(formData)
    ElMessage.success('头像更新成功')
    // 重新加载用户信息
    window.location.reload()
  } catch (error) {
    // 错误已在拦截器中处理
  }
  
  event.target.value = ''
}
</script>

<style scoped>
.navbar {
  height: 64px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-icon {
  font-size: 28px;
  color: #667eea;
}

.brand-text {
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.navbar-menu {
  display: flex;
  gap: 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 8px;
  color: #666;
  text-decoration: none;
  font-size: 14px;
  transition: all 0.3s ease;
}

.nav-item:hover {
  background: #f5f7fa;
  color: #667eea;
}

.nav-item.active {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  color: #667eea;
  font-weight: 500;
}

.navbar-user {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 24px;
  transition: all 0.3s ease;
}

.user-info:hover {
  background: #f5f7fa;
}

.user-name {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}
</style>
