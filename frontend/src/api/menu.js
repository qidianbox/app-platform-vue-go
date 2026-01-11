import request from '@/utils/request'

/**
 * 获取APP菜单列表（树形结构）
 * @param {number|string} appId APP ID
 */
export function getAppMenus(appId) {
  return request({
    url: `/apps/${appId}/menus`,
    method: 'get'
  })
}

/**
 * 获取APP菜单列表（平铺结构）
 * @param {number|string} appId APP ID
 */
export function getAppMenuList(appId) {
  return request({
    url: `/apps/${appId}/menus/list`,
    method: 'get'
  })
}

/**
 * 创建菜单
 * @param {number|string} appId APP ID
 * @param {Object} data 菜单数据
 */
export function createMenu(appId, data) {
  return request({
    url: `/apps/${appId}/menus`,
    method: 'post',
    data
  })
}

/**
 * 获取菜单详情
 * @param {number|string} appId APP ID
 * @param {number} menuId 菜单ID
 */
export function getMenuDetail(appId, menuId) {
  return request({
    url: `/apps/${appId}/menus/${menuId}`,
    method: 'get'
  })
}

/**
 * 更新菜单
 * @param {number|string} appId APP ID
 * @param {number} menuId 菜单ID
 * @param {Object} data 菜单数据
 */
export function updateMenu(appId, menuId, data) {
  return request({
    url: `/apps/${appId}/menus/${menuId}`,
    method: 'put',
    data
  })
}

/**
 * 删除菜单
 * @param {number|string} appId APP ID
 * @param {number} menuId 菜单ID
 */
export function deleteMenu(appId, menuId) {
  return request({
    url: `/apps/${appId}/menus/${menuId}`,
    method: 'delete'
  })
}

/**
 * 批量更新菜单排序
 * @param {number|string} appId APP ID
 * @param {Array} sortData 排序数据
 */
export function updateMenuSort(appId, sortData) {
  return request({
    url: `/apps/${appId}/menus/sort`,
    method: 'put',
    data: sortData
  })
}
