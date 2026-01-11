<template>
  <div class="app-config-layout">
    <!-- 顶部导航栏 -->
    <header class="app-header">
      <div class="header-left">
        <el-button type="text" @click="goBack" class="back-btn">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <div class="app-info" v-if="appInfo">
          <span class="app-name">{{ appInfo.name }}</span>
          <el-tag size="small" type="info">{{ appInfo.app_id }}</el-tag>
        </div>
      </div>
      <div class="header-center">
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          @select="handleMenuSelect"
          class="top-menu"
        >
          <el-menu-item index="config">
            <el-icon><Setting /></el-icon>
            概览
          </el-menu-item>
          <el-menu-item index="menus">
            <el-icon><Menu /></el-icon>
            菜单管理
          </el-menu-item>
          <el-menu-item index="apis">
            <el-icon><Connection /></el-icon>
            API管理
          </el-menu-item>
        </el-menu>
      </div>
      <div class="header-right">
        <el-dropdown>
          <span class="user-info">
            <el-avatar :size="32" icon="UserFilled" />
            <span class="username">{{ username }}</span>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>

    <!-- 主内容区域 -->
    <main class="app-main">
      <router-view v-if="appInfo" />
      <div v-else class="loading-container">
        <el-skeleton :rows="10" animated />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Setting, Menu, Connection } from '@element-plus/icons-vue'
import { getAppDetail } from '@/api/app'

const route = useRoute()
const router = useRouter()

const appId = computed(() => route.params.appId)
const appInfo = ref(null)
const username = ref(localStorage.getItem('username') || 'Admin')

// 当前激活的菜单
const activeMenu = computed(() => {
  const path = route.path
  if (path.includes('/menus')) return 'menus'
  if (path.includes('/apis')) return 'apis'
  return 'config'
})

// 加载APP信息
const loadAppInfo = async () => {
  try {
    const res = await getAppDetail(appId.value)
    appInfo.value = res || res?.data
  } catch (error) {
    console.error('加载APP信息失败:', error)
    ElMessage.error('加载APP信息失败')
  }
}

// 返回APP列表
const goBack = () => {
  router.push('/apps')
}

// 菜单选择
const handleMenuSelect = (index) => {
  router.push(`/apps/${appId.value}/${index}`)
}

// 退出登录
const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  router.push('/login')
}

// 监听路由变化
watch(appId, (newVal) => {
  if (newVal) {
    loadAppInfo()
  }
}, { immediate: true })

onMounted(() => {
  loadAppInfo()
})
</script>

<style scoped>
.app-config-layout {
  min-height: 100vh;
  background: #f5f7fa;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding: 0 20px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.back-btn {
  font-size: 14px;
  color: #606266;
}

.back-btn:hover {
  color: #409eff;
}

.app-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.app-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.top-menu {
  border-bottom: none;
}

.top-menu .el-menu-item {
  height: 60px;
  line-height: 60px;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  font-size: 14px;
  color: #606266;
}

.app-main {
  padding: 20px;
  min-height: calc(100vh - 60px);
}

.loading-container {
  padding: 40px;
  background: #fff;
  border-radius: 4px;
}
</style>
