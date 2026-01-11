import request from '@/utils/request'

// ==================== 系统API管理 ====================

/**
 * 获取系统API列表
 * @param {Object} params 查询参数
 */
export function getSystemAPIs(params) {
  return request({
    url: '/system-apis',
    method: 'get',
    params
  })
}

/**
 * 获取API分类列表
 */
export function getAPICategories() {
  return request({
    url: '/system-apis/categories',
    method: 'get'
  })
}

/**
 * 获取API所属模块列表
 */
export function getAPIModules() {
  return request({
    url: '/system-apis/modules',
    method: 'get'
  })
}

// ==================== APP API授权管理 ====================

/**
 * 获取APP已授权的API列表
 * @param {number|string} appId APP ID
 */
export function getAppAPIPermissions(appId) {
  return request({
    url: `/apps/${appId}/api-permissions`,
    method: 'get'
  })
}

/**
 * 授权API给APP
 * @param {number|string} appId APP ID
 * @param {Object} data 授权数据
 */
export function grantAPIPermission(appId, data) {
  return request({
    url: `/apps/${appId}/api-permissions`,
    method: 'post',
    data
  })
}

/**
 * 撤销API授权
 * @param {number|string} appId APP ID
 * @param {string} apiCode API标识码
 */
export function revokeAPIPermission(appId, apiCode) {
  return request({
    url: `/apps/${appId}/api-permissions/${apiCode}`,
    method: 'delete'
  })
}

// ==================== APP API密钥管理 ====================

/**
 * 获取APP的API密钥列表
 * @param {number|string} appId APP ID
 */
export function getAppAPIKeys(appId) {
  return request({
    url: `/apps/${appId}/api-keys`,
    method: 'get'
  })
}

/**
 * 创建API密钥
 * @param {number|string} appId APP ID
 * @param {Object} data 密钥数据
 */
export function createAppAPIKey(appId, data) {
  return request({
    url: `/apps/${appId}/api-keys`,
    method: 'post',
    data
  })
}

/**
 * 更新API密钥状态
 * @param {number|string} appId APP ID
 * @param {number} keyId 密钥ID
 * @param {number} status 状态
 */
export function updateAppAPIKeyStatus(appId, keyId, status) {
  return request({
    url: `/apps/${appId}/api-keys/${keyId}/status`,
    method: 'put',
    data: { status }
  })
}

/**
 * 删除API密钥
 * @param {number|string} appId APP ID
 * @param {number} keyId 密钥ID
 */
export function deleteAppAPIKey(appId, keyId) {
  return request({
    url: `/apps/${appId}/api-keys/${keyId}`,
    method: 'delete'
  })
}

// ==================== API调用统计 ====================

/**
 * 获取APP的API调用统计
 * @param {number|string} appId APP ID
 */
export function getAppAPIStats(appId) {
  return request({
    url: `/apps/${appId}/api-stats`,
    method: 'get'
  })
}

/**
 * 获取APP的API调用日志
 * @param {number|string} appId APP ID
 * @param {Object} params 查询参数
 */
export function getAppAPICallLogs(appId, params) {
  return request({
    url: `/apps/${appId}/api-logs`,
    method: 'get',
    params
  })
}
