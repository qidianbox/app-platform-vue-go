/**
 * 前端调试工具
 */

// 调试模式开关
export const DEBUG_MODE = import.meta.env.DEV || localStorage.getItem('debug_mode') === 'true'

// 日志级别
export const LogLevel = {
  DEBUG: 0,
  INFO: 1,
  WARN: 2,
  ERROR: 3
}

// 当前日志级别
let currentLogLevel = DEBUG_MODE ? LogLevel.DEBUG : LogLevel.INFO

/**
 * 设置日志级别
 */
export function setLogLevel(level) {
  currentLogLevel = level
}

/**
 * 调试日志
 */
export function debug(...args) {
  if (currentLogLevel <= LogLevel.DEBUG) {
    console.log('[DEBUG]', new Date().toISOString(), ...args)
  }
}

/**
 * 信息日志
 */
export function info(...args) {
  if (currentLogLevel <= LogLevel.INFO) {
    console.log('[INFO]', new Date().toISOString(), ...args)
  }
}

/**
 * 警告日志
 */
export function warn(...args) {
  if (currentLogLevel <= LogLevel.WARN) {
    console.warn('[WARN]', new Date().toISOString(), ...args)
  }
}

/**
 * 错误日志
 */
export function error(...args) {
  if (currentLogLevel <= LogLevel.ERROR) {
    console.error('[ERROR]', new Date().toISOString(), ...args)
  }
}

/**
 * API 请求日志
 */
export function logRequest(method, url, params, data) {
  if (DEBUG_MODE) {
    console.group(`%c[API REQUEST] ${method} ${url}`, 'color: #4CAF50; font-weight: bold')
    if (params) console.log('Params:', params)
    if (data) console.log('Data:', data)
    console.groupEnd()
  }
}

/**
 * API 响应日志
 */
export function logResponse(method, url, response, duration) {
  if (DEBUG_MODE) {
    const statusColor = response.status >= 200 && response.status < 300 ? '#4CAF50' : '#F44336'
    console.group(`%c[API RESPONSE] ${method} ${url} (${duration}ms)`, `color: ${statusColor}; font-weight: bold`)
    console.log('Status:', response.status)
    console.log('Headers:', response.headers)
    console.log('Data:', response.data)
    
    // 显示调试头
    const debugHeaders = {}
    for (const [key, value] of Object.entries(response.headers)) {
      if (key.toLowerCase().startsWith('x-debug-')) {
        debugHeaders[key] = value
      }
    }
    if (Object.keys(debugHeaders).length > 0) {
      console.log('%cDebug Headers:', 'color: #FF9800; font-weight: bold', debugHeaders)
    }
    
    console.groupEnd()
  }
}

/**
 * API 错误日志
 */
export function logError(method, url, error, duration) {
  if (DEBUG_MODE) {
    console.group(`%c[API ERROR] ${method} ${url} (${duration}ms)`, 'color: #F44336; font-weight: bold')
    console.error('Error:', error)
    if (error.response) {
      console.log('Status:', error.response.status)
      console.log('Headers:', error.response.headers)
      console.log('Data:', error.response.data)
      
      // 显示调试头
      const debugHeaders = {}
      for (const [key, value] of Object.entries(error.response.headers)) {
        if (key.toLowerCase().startsWith('x-debug-')) {
          debugHeaders[key] = value
        }
      }
      if (Object.keys(debugHeaders).length > 0) {
        console.log('%cDebug Headers:', 'color: #FF9800; font-weight: bold', debugHeaders)
      }
    }
    console.groupEnd()
  }
}

/**
 * 组件生命周期日志
 */
export function logLifecycle(componentName, lifecycle, data) {
  if (DEBUG_MODE) {
    console.log(`%c[${componentName}] ${lifecycle}`, 'color: #2196F3; font-weight: bold', data || '')
  }
}

/**
 * 数据变化日志
 */
export function logDataChange(name, oldValue, newValue) {
  if (DEBUG_MODE) {
    console.log(`%c[DATA CHANGE] ${name}`, 'color: #9C27B0; font-weight: bold')
    console.log('Old:', oldValue)
    console.log('New:', newValue)
  }
}

/**
 * 性能监控
 */
export function measurePerformance(name, fn) {
  if (DEBUG_MODE) {
    const start = performance.now()
    const result = fn()
    const end = performance.now()
    console.log(`%c[PERFORMANCE] ${name}: ${(end - start).toFixed(2)}ms`, 'color: #FF5722; font-weight: bold')
    return result
  }
  return fn()
}

/**
 * 启用调试模式
 */
export function enableDebugMode() {
  localStorage.setItem('debug_mode', 'true')
  window.location.reload()
}

/**
 * 禁用调试模式
 */
export function disableDebugMode() {
  localStorage.removeItem('debug_mode')
  window.location.reload()
}

// 在开发环境下，将调试工具挂载到 window 对象
if (DEBUG_MODE) {
  window.debug = {
    enable: enableDebugMode,
    disable: disableDebugMode,
    setLogLevel,
    LogLevel
  }
  console.log('%c🔧 调试模式已启用', 'color: #4CAF50; font-size: 14px; font-weight: bold')
  console.log('%c使用 window.debug 访问调试工具', 'color: #2196F3; font-size: 12px')
}

export default {
  DEBUG_MODE,
  LogLevel,
  setLogLevel,
  debug,
  info,
  warn,
  error,
  logRequest,
  logResponse,
  logError,
  logLifecycle,
  logDataChange,
  measurePerformance,
  enableDebugMode,
  disableDebugMode
}
