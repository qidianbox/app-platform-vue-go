<template>
  <div class="menu-management">
    <!-- 页面标题和操作栏 -->
    <div class="page-header">
      <div class="header-left">
        <h2>菜单管理</h2>
        <span class="subtitle">管理应用内的菜单结构和导航</span>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="handleAdd(null)">
          <el-icon><Plus /></el-icon>
          新增菜单
        </el-button>
      </div>
    </div>

    <!-- 菜单树表格 -->
    <el-card class="menu-card">
      <el-table
        :data="menuTree"
        row-key="id"
        border
        default-expand-all
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        v-loading="loading"
      >
        <el-table-column prop="name" label="菜单名称" min-width="200">
          <template #default="{ row }">
            <div class="menu-name">
              <el-icon v-if="row.icon" class="menu-icon">
                <component :is="row.icon" />
              </el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="code" label="菜单标识" width="150" />
        <el-table-column prop="path" label="路由路径" width="180" />
        <el-table-column prop="menu_type" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getMenuTypeTag(row.menu_type)">
              {{ getMenuTypeName(row.menu_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" align="center" />
        <el-table-column prop="visible" label="可见" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.visible"
              :active-value="1"
              :inactive-value="0"
              @change="handleVisibleChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleAdd(row)">
              添加子菜单
            </el-button>
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="上级菜单">
          <el-tree-select
            v-model="formData.parent_id"
            :data="menuTreeOptions"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="选择上级菜单（不选则为顶级菜单）"
            clearable
            check-strictly
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="菜单类型" prop="menu_type">
          <el-radio-group v-model="formData.menu_type">
            <el-radio :value="1">目录</el-radio>
            <el-radio :value="2">菜单</el-radio>
            <el-radio :value="3">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="菜单标识" prop="code">
          <el-input v-model="formData.code" placeholder="请输入菜单标识（英文）" />
        </el-form-item>
        <el-form-item label="图标" v-if="formData.menu_type !== 3">
          <el-input v-model="formData.icon" placeholder="请输入图标名称" />
        </el-form-item>
        <el-form-item label="路由路径" v-if="formData.menu_type !== 3">
          <el-input v-model="formData.path" placeholder="请输入路由路径" />
        </el-form-item>
        <el-form-item label="组件路径" v-if="formData.menu_type === 2">
          <el-input v-model="formData.component" placeholder="请输入组件路径" />
        </el-form-item>
        <el-form-item label="权限标识" v-if="formData.menu_type === 3">
          <el-input v-model="formData.permission" placeholder="请输入权限标识" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="formData.sort_order" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="是否可见">
          <el-switch v-model="formData.visible" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
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
  getAppMenus,
  createMenu,
  updateMenu,
  deleteMenu
} from '@/api/menu'

const route = useRoute()
const appId = computed(() => route.params.id)

// 状态
const loading = ref(false)
const submitting = ref(false)
const menuTree = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('新增菜单')
const editingId = ref(null)

// 表单
const formRef = ref(null)
const formData = ref({
  parent_id: 0,
  menu_type: 2,
  name: '',
  code: '',
  icon: '',
  path: '',
  component: '',
  permission: '',
  sort_order: 0,
  visible: 1,
  remark: ''
})

const formRules = {
  name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  code: [
    { required: true, message: '请输入菜单标识', trigger: 'blur' },
    { pattern: /^[a-zA-Z_][a-zA-Z0-9_]*$/, message: '只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  menu_type: [{ required: true, message: '请选择菜单类型', trigger: 'change' }]
}

// 菜单树选项（用于选择上级菜单）
const menuTreeOptions = computed(() => {
  const options = [{ id: 0, name: '顶级菜单', children: [] }]
  if (menuTree.value.length > 0) {
    options[0].children = menuTree.value
  }
  return options
})

// 获取菜单类型名称
const getMenuTypeName = (type) => {
  const types = { 1: '目录', 2: '菜单', 3: '按钮' }
  return types[type] || '未知'
}

// 获取菜单类型标签样式
const getMenuTypeTag = (type) => {
  const tags = { 1: 'warning', 2: 'success', 3: 'info' }
  return tags[type] || 'info'
}

// 加载菜单列表
const loadMenus = async () => {
  loading.value = true
  try {
    const res = await getAppMenus(appId.value)
    menuTree.value = Array.isArray(res) ? res : (res?.data || [])
  } catch (error) {
    console.error('加载菜单失败:', error)
    ElMessage.error('加载菜单失败')
  } finally {
    loading.value = false
  }
}

// 新增菜单
const handleAdd = (parent) => {
  editingId.value = null
  dialogTitle.value = '新增菜单'
  formData.value = {
    parent_id: parent ? parent.id : 0,
    menu_type: 2,
    name: '',
    code: '',
    icon: '',
    path: '',
    component: '',
    permission: '',
    sort_order: 0,
    visible: 1,
    remark: ''
  }
  dialogVisible.value = true
}

// 编辑菜单
const handleEdit = (row) => {
  editingId.value = row.id
  dialogTitle.value = '编辑菜单'
  formData.value = { ...row }
  dialogVisible.value = true
}

// 删除菜单
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除菜单"${row.name}"吗？`,
      '删除确认',
      { type: 'warning' }
    )
    await deleteMenu(appId.value, row.id)
    ElMessage.success('删除成功')
    loadMenus()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    if (editingId.value) {
      await updateMenu(appId.value, editingId.value, formData.value)
      ElMessage.success('更新成功')
    } else {
      await createMenu(appId.value, formData.value)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    loadMenus()
  } catch (error) {
    if (error !== false) {
      ElMessage.error(error.message || '操作失败')
    }
  } finally {
    submitting.value = false
  }
}

// 切换可见状态
const handleVisibleChange = async (row) => {
  try {
    await updateMenu(appId.value, row.id, { visible: row.visible })
    ElMessage.success('更新成功')
  } catch (error) {
    row.visible = row.visible === 1 ? 0 : 1
    ElMessage.error('更新失败')
  }
}

// 切换状态
const handleStatusChange = async (row) => {
  try {
    await updateMenu(appId.value, row.id, { status: row.status })
    ElMessage.success('更新成功')
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('更新失败')
  }
}

onMounted(() => {
  loadMenus()
})
</script>

<style scoped>
.menu-management {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.header-left .subtitle {
  font-size: 14px;
  color: #909399;
  margin-left: 12px;
}

.menu-card {
  background: #fff;
}

.menu-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.menu-icon {
  color: #409eff;
}
</style>
