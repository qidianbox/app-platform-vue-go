<template>
  <div class="api-management">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>API 管理</h2>
      <span class="subtitle">管理应用的 API 授权和密钥</span>
    </div>

    <!-- 标签页 -->
    <el-tabs v-model="activeTab" class="api-tabs">
      <!-- API 授权 -->
      <el-tab-pane label="API 授权" name="permissions">
        <div class="tab-header">
          <el-button type="primary" @click="showGrantDialog">
            <el-icon><Plus /></el-icon>
            授权 API
          </el-button>
        </div>
        
        <el-table :data="permissions" v-loading="loadingPermissions" border>
          <el-table-column prop="api_name" label="API 名称" min-width="150" />
          <el-table-column prop="api_code" label="API 标识" width="150" />
          <el-table-column prop="api_method" label="方法" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="getMethodType(row.api_method)" size="small">
                {{ row.api_method }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="api_path" label="路径" min-width="200" />
          <el-table-column prop="module_code" label="所属模块" width="120" />
          <el-table-column prop="rate_limit" label="限流" width="100" align="center">
            <template #default="{ row }">
              {{ row.rate_limit || '不限制' }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button type="danger" link size="small" @click="handleRevoke(row)">
                撤销
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- API 密钥 -->
      <el-tab-pane label="API 密钥" name="keys">
        <div class="tab-header">
          <el-button type="primary" @click="showCreateKeyDialog">
            <el-icon><Plus /></el-icon>
            创建密钥
          </el-button>
        </div>

        <el-table :data="apiKeys" v-loading="loadingKeys" border>
          <el-table-column prop="name" label="密钥名称" min-width="150" />
          <el-table-column prop="api_key" label="API Key" min-width="250">
            <template #default="{ row }">
              <div class="key-cell">
                <code>{{ row.api_key }}</code>
                <el-button type="primary" link size="small" @click="copyToClipboard(row.api_key)">
                  复制
                </el-button>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-switch
                v-model="row.status"
                :active-value="1"
                :inactive-value="0"
                @change="handleKeyStatusChange(row)"
              />
            </template>
          </el-table-column>
          <el-table-column prop="last_used_at" label="最后使用" width="180">
            <template #default="{ row }">
              {{ row.last_used_at || '从未使用' }}
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180" />
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button type="danger" link size="small" @click="handleDeleteKey(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- 调用统计 -->
      <el-tab-pane label="调用统计" name="stats">
        <div class="stats-summary">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-card class="stat-card">
                <div class="stat-value">{{ stats.summary?.total_calls || 0 }}</div>
                <div class="stat-label">总调用次数</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card success">
                <div class="stat-value">{{ stats.summary?.success_calls || 0 }}</div>
                <div class="stat-label">成功次数</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card danger">
                <div class="stat-value">{{ stats.summary?.fail_calls || 0 }}</div>
                <div class="stat-label">失败次数</div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card info">
                <div class="stat-value">{{ (stats.summary?.success_rate || 0).toFixed(1) }}%</div>
                <div class="stat-label">成功率</div>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <el-card class="logs-card">
          <template #header>
            <span>调用日志</span>
          </template>
          <el-table :data="callLogs" v-loading="loadingLogs" border>
            <el-table-column prop="api_code" label="API" width="150" />
            <el-table-column prop="request_method" label="方法" width="80" align="center" />
            <el-table-column prop="request_path" label="路径" min-width="200" />
            <el-table-column prop="response_code" label="状态码" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="row.response_code < 400 ? 'success' : 'danger'" size="small">
                  {{ row.response_code }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="duration" label="耗时" width="100" align="center">
              <template #default="{ row }">
                {{ row.duration }}ms
              </template>
            </el-table-column>
            <el-table-column prop="client_ip" label="客户端IP" width="130" />
            <el-table-column prop="created_at" label="时间" width="180" />
          </el-table>
          <div class="pagination">
            <el-pagination
              v-model:current-page="logsPage"
              v-model:page-size="logsPageSize"
              :total="logsTotal"
              layout="total, prev, pager, next"
              @current-change="loadCallLogs"
            />
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 授权 API 对话框 -->
    <el-dialog v-model="grantDialogVisible" title="授权 API" width="800px" destroy-on-close>
      <div class="api-filter">
        <el-input
          v-model="apiKeyword"
          placeholder="搜索 API"
          style="width: 200px"
          clearable
          @input="filterAPIs"
        />
        <el-select v-model="apiModuleFilter" placeholder="选择模块" clearable style="width: 150px; margin-left: 10px" @change="filterAPIs">
          <el-option v-for="m in apiModules" :key="m.module_code" :label="m.module_code" :value="m.module_code" />
        </el-select>
      </div>
      <el-table
        ref="apiTableRef"
        :data="filteredSystemAPIs"
        v-loading="loadingSystemAPIs"
        border
        max-height="400"
        @selection-change="handleAPISelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="API 名称" min-width="150" />
        <el-table-column prop="code" label="标识" width="150" />
        <el-table-column prop="method" label="方法" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="getMethodType(row.method)" size="small">{{ row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路径" min-width="200" />
        <el-table-column prop="module_code" label="模块" width="120" />
      </el-table>
      <template #footer>
        <el-button @click="grantDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleGrant" :loading="granting" :disabled="selectedAPIs.length === 0">
          授权 ({{ selectedAPIs.length }})
        </el-button>
      </template>
    </el-dialog>

    <!-- 创建密钥对话框 -->
    <el-dialog v-model="keyDialogVisible" title="创建 API 密钥" width="500px" destroy-on-close>
      <el-form ref="keyFormRef" :model="keyFormData" :rules="keyFormRules" label-width="100px">
        <el-form-item label="密钥名称" prop="name">
          <el-input v-model="keyFormData.name" placeholder="请输入密钥名称" />
        </el-form-item>
        <el-form-item label="IP 白名单">
          <el-input
            v-model="keyFormData.ip_whitelist"
            type="textarea"
            :rows="3"
            placeholder="多个 IP 用逗号分隔，留空表示不限制"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="keyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateKey" :loading="creatingKey">创建</el-button>
      </template>
    </el-dialog>

    <!-- 显示新创建的密钥 -->
    <el-dialog v-model="newKeyDialogVisible" title="密钥创建成功" width="600px">
      <el-alert type="warning" :closable="false" show-icon style="margin-bottom: 20px">
        请妥善保存以下密钥信息，API Secret 仅显示一次！
      </el-alert>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="密钥名称">{{ newKeyData.name }}</el-descriptions-item>
        <el-descriptions-item label="API Key">
          <code>{{ newKeyData.api_key }}</code>
          <el-button type="primary" link size="small" @click="copyToClipboard(newKeyData.api_key)">复制</el-button>
        </el-descriptions-item>
        <el-descriptions-item label="API Secret">
          <code>{{ newKeyData.api_secret }}</code>
          <el-button type="primary" link size="small" @click="copyToClipboard(newKeyData.api_secret)">复制</el-button>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button type="primary" @click="newKeyDialogVisible = false">我已保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getSystemAPIs,
  getAPIModules,
  getAppAPIPermissions,
  grantAPIPermission,
  revokeAPIPermission,
  getAppAPIKeys,
  createAppAPIKey,
  updateAppAPIKeyStatus,
  deleteAppAPIKey,
  getAppAPIStats,
  getAppAPICallLogs
} from '@/api/apiManager'

const route = useRoute()
const appId = computed(() => route.params.appId)

// 标签页
const activeTab = ref('permissions')

// API 授权
const permissions = ref([])
const loadingPermissions = ref(false)

// API 密钥
const apiKeys = ref([])
const loadingKeys = ref(false)

// 调用统计
const stats = ref({ summary: {} })
const callLogs = ref([])
const loadingLogs = ref(false)
const logsPage = ref(1)
const logsPageSize = ref(10)
const logsTotal = ref(0)

// 授权对话框
const grantDialogVisible = ref(false)
const systemAPIs = ref([])
const filteredSystemAPIs = ref([])
const loadingSystemAPIs = ref(false)
const selectedAPIs = ref([])
const granting = ref(false)
const apiKeyword = ref('')
const apiModuleFilter = ref('')
const apiModules = ref([])

// 密钥对话框
const keyDialogVisible = ref(false)
const keyFormRef = ref(null)
const keyFormData = ref({ name: '', ip_whitelist: '' })
const keyFormRules = {
  name: [{ required: true, message: '请输入密钥名称', trigger: 'blur' }]
}
const creatingKey = ref(false)

// 新密钥显示
const newKeyDialogVisible = ref(false)
const newKeyData = ref({})

// 获取方法类型样式
const getMethodType = (method) => {
  const types = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger'
  }
  return types[method] || 'info'
}

// 加载 API 授权列表
const loadPermissions = async () => {
  loadingPermissions.value = true
  try {
    const res = await getAppAPIPermissions(appId.value)
    permissions.value = Array.isArray(res) ? res : (res?.data || [])
  } catch (error) {
    console.error('加载授权列表失败:', error)
  } finally {
    loadingPermissions.value = false
  }
}

// 加载 API 密钥列表
const loadAPIKeys = async () => {
  loadingKeys.value = true
  try {
    const res = await getAppAPIKeys(appId.value)
    apiKeys.value = Array.isArray(res) ? res : (res?.data || [])
  } catch (error) {
    console.error('加载密钥列表失败:', error)
  } finally {
    loadingKeys.value = false
  }
}

// 加载调用统计
const loadStats = async () => {
  try {
    const res = await getAppAPIStats(appId.value)
    stats.value = res || { summary: {} }
  } catch (error) {
    console.error('加载统计失败:', error)
  }
}

// 加载调用日志
const loadCallLogs = async () => {
  loadingLogs.value = true
  try {
    const res = await getAppAPICallLogs(appId.value, {
      page: logsPage.value,
      page_size: logsPageSize.value
    })
    callLogs.value = res?.data || []
    logsTotal.value = res?.pagination?.total || 0
  } catch (error) {
    console.error('加载日志失败:', error)
  } finally {
    loadingLogs.value = false
  }
}

// 加载系统 API 列表
const loadSystemAPIs = async () => {
  loadingSystemAPIs.value = true
  try {
    const [apisRes, modulesRes] = await Promise.all([
      getSystemAPIs({ page_size: 100 }),
      getAPIModules()
    ])
    // axios 响应拦截器已解包，apisRes 可能直接是数组或 { data: [...] }
    systemAPIs.value = Array.isArray(apisRes) ? apisRes : (apisRes?.data || apisRes?.list || [])
    filteredSystemAPIs.value = systemAPIs.value
    apiModules.value = modulesRes || []
  } catch (error) {
    console.error('加载系统 API 失败:', error)
  } finally {
    loadingSystemAPIs.value = false
  }
}

// 过滤 API
const filterAPIs = () => {
  let result = systemAPIs.value
  if (apiKeyword.value) {
    const keyword = apiKeyword.value.toLowerCase()
    result = result.filter(api =>
      api.name.toLowerCase().includes(keyword) ||
      api.code.toLowerCase().includes(keyword) ||
      api.path.toLowerCase().includes(keyword)
    )
  }
  if (apiModuleFilter.value) {
    result = result.filter(api => api.module_code === apiModuleFilter.value)
  }
  filteredSystemAPIs.value = result
}

// 显示授权对话框
const showGrantDialog = async () => {
  grantDialogVisible.value = true
  selectedAPIs.value = []
  apiKeyword.value = ''
  apiModuleFilter.value = ''
  await loadSystemAPIs()
}

// API 选择变化
const handleAPISelectionChange = (selection) => {
  selectedAPIs.value = selection
}

// 授权 API
const handleGrant = async () => {
  granting.value = true
  try {
    const apiCodes = selectedAPIs.value.map(api => api.code)
    await grantAPIPermission(appId.value, { api_codes: apiCodes })
    ElMessage.success('授权成功')
    grantDialogVisible.value = false
    loadPermissions()
  } catch (error) {
    ElMessage.error(error.message || '授权失败')
  } finally {
    granting.value = false
  }
}

// 撤销授权
const handleRevoke = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要撤销 API "${row.api_name}" 的授权吗？`, '撤销确认', { type: 'warning' })
    await revokeAPIPermission(appId.value, row.api_code)
    ElMessage.success('撤销成功')
    loadPermissions()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '撤销失败')
    }
  }
}

// 显示创建密钥对话框
const showCreateKeyDialog = () => {
  keyFormData.value = { name: '', ip_whitelist: '' }
  keyDialogVisible.value = true
}

// 创建密钥
const handleCreateKey = async () => {
  try {
    await keyFormRef.value.validate()
    creatingKey.value = true
    const res = await createAppAPIKey(appId.value, keyFormData.value)
    keyDialogVisible.value = false
    newKeyData.value = res?.data || res
    newKeyDialogVisible.value = true
    loadAPIKeys()
  } catch (error) {
    if (error !== false) {
      ElMessage.error(error.message || '创建失败')
    }
  } finally {
    creatingKey.value = false
  }
}

// 更新密钥状态
const handleKeyStatusChange = async (row) => {
  try {
    await updateAppAPIKeyStatus(appId.value, row.id, row.status)
    ElMessage.success('更新成功')
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('更新失败')
  }
}

// 删除密钥
const handleDeleteKey = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除密钥 "${row.name}" 吗？`, '删除确认', { type: 'warning' })
    await deleteAppAPIKey(appId.value, row.id)
    ElMessage.success('删除成功')
    loadAPIKeys()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 复制到剪贴板
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

onMounted(() => {
  loadPermissions()
  loadAPIKeys()
  loadStats()
  loadCallLogs()
})
</script>

<style scoped>
.api-management {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.page-header .subtitle {
  font-size: 14px;
  color: #909399;
  margin-left: 12px;
}

.api-tabs {
  background: #fff;
  padding: 20px;
  border-radius: 4px;
}

.tab-header {
  margin-bottom: 20px;
}

.key-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.key-cell code {
  background: #f5f7fa;
  padding: 4px 8px;
  border-radius: 4px;
  font-family: monospace;
}

.stats-summary {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
  padding: 20px;
}

.stat-card .stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
}

.stat-card .stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 8px;
}

.stat-card.success .stat-value {
  color: #67c23a;
}

.stat-card.danger .stat-value {
  color: #f56c6c;
}

.stat-card.info .stat-value {
  color: #409eff;
}

.logs-card {
  margin-top: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.api-filter {
  margin-bottom: 15px;
}
</style>
