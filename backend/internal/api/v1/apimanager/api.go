package apimanager

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"time"

	"app-platform-backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB(database *gorm.DB) {
	db = database
}

// getAppDatabaseID 获取APP的数据库ID
func getAppDatabaseID(appIDParam string) (uint, error) {
	if id, err := strconv.ParseUint(appIDParam, 10, 64); err == nil {
		return uint(id), nil
	}
	var app model.App
	if err := db.Where("app_id = ?", appIDParam).First(&app).Error; err != nil {
		return 0, err
	}
	return app.ID, nil
}

// generateAPIKey 生成API Key
func generateAPIKey() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return "ak_" + hex.EncodeToString(bytes)
}

// generateAPISecret 生成API Secret
func generateAPISecret() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ==================== 系统API管理 ====================

// GetSystemAPIs 获取系统API列表
func GetSystemAPIs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	moduleCode := c.Query("module_code")
	category := c.Query("category")
	keyword := c.Query("keyword")

	query := db.Model(&model.SystemAPI{})

	if moduleCode != "" {
		query = query.Where("module_code = ?", moduleCode)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR path LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var apis []model.SystemAPI
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("module_code, id").Find(&apis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取API列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": apis,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetSystemAPICategories 获取API分类列表
func GetSystemAPICategories(c *gin.Context) {
	var categories []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}

	db.Model(&model.SystemAPI{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Order("category").
		Scan(&categories)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": categories,
	})
}

// GetSystemAPIModules 获取API所属模块列表
func GetSystemAPIModules(c *gin.Context) {
	var modules []struct {
		ModuleCode string `json:"module_code"`
		Count      int64  `json:"count"`
	}

	db.Model(&model.SystemAPI{}).
		Select("module_code, COUNT(*) as count").
		Group("module_code").
		Order("module_code").
		Scan(&modules)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": modules,
	})
}

// ==================== APP API授权管理 ====================

// GetAppAPIPermissions 获取APP已授权的API列表
func GetAppAPIPermissions(c *gin.Context) {
	appIDParam := c.Param("id")
	log.Printf("[DEBUG] GetAppAPIPermissions - appIDParam: %s", appIDParam)

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	var permissions []model.AppAPIPermission
	if err := db.Where("app_id = ?", appID).Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取授权列表失败"})
		return
	}

	// 获取API详情
	var apiCodes []string
	for _, p := range permissions {
		apiCodes = append(apiCodes, p.APICode)
	}

	var apis []model.SystemAPI
	if len(apiCodes) > 0 {
		db.Where("code IN ?", apiCodes).Find(&apis)
	}

	// 合并数据
	apiMap := make(map[string]model.SystemAPI)
	for _, api := range apis {
		apiMap[api.Code] = api
	}

	type PermissionWithAPI struct {
		model.AppAPIPermission
		APIName     string `json:"api_name"`
		APIPath     string `json:"api_path"`
		APIMethod   string `json:"api_method"`
		ModuleCode  string `json:"module_code"`
		Description string `json:"description"`
	}

	var result []PermissionWithAPI
	for _, p := range permissions {
		item := PermissionWithAPI{AppAPIPermission: p}
		if api, ok := apiMap[p.APICode]; ok {
			item.APIName = api.Name
			item.APIPath = api.Path
			item.APIMethod = api.Method
			item.ModuleCode = api.ModuleCode
			item.Description = api.Description
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
}

// GrantAPIPermission 授权API给APP
func GrantAPIPermission(c *gin.Context) {
	appIDParam := c.Param("id")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	var req struct {
		APICodes  []string `json:"api_codes"`
		RateLimit int      `json:"rate_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 获取API信息
	var apis []model.SystemAPI
	if err := db.Where("code IN ?", req.APICodes).Find(&apis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取API信息失败"})
		return
	}

	tx := db.Begin()
	for _, api := range apis {
		permission := model.AppAPIPermission{
			AppID:     appID,
			APIID:     api.ID,
			APICode:   api.Code,
			Status:    1,
			RateLimit: req.RateLimit,
		}

		// 使用 ON DUPLICATE KEY UPDATE
		if err := tx.Where("app_id = ? AND api_id = ?", appID, api.ID).
			Assign(permission).
			FirstOrCreate(&permission).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "授权失败"})
			return
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "授权成功",
	})
}

// RevokeAPIPermission 撤销API授权
func RevokeAPIPermission(c *gin.Context) {
	appIDParam := c.Param("id")
	apiCode := c.Param("apiCode")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	if err := db.Where("app_id = ? AND api_code = ?", appID, apiCode).
		Delete(&model.AppAPIPermission{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "撤销授权失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "撤销成功",
	})
}

// ==================== APP API密钥管理 ====================

// GetAppAPIKeys 获取APP的API密钥列表
func GetAppAPIKeys(c *gin.Context) {
	appIDParam := c.Param("id")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	var keys []model.AppAPIKey
	if err := db.Where("app_id = ?", appID).Find(&keys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取密钥列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": keys,
	})
}

// CreateAppAPIKey 创建API密钥
func CreateAppAPIKey(c *gin.Context) {
	appIDParam := c.Param("id")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		IPWhitelist string `json:"ip_whitelist"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: 密钥名称不能为空"})
		return
	}

	apiKey := generateAPIKey()
	apiSecret := generateAPISecret()

	key := model.AppAPIKey{
		AppID:       appID,
		Name:        req.Name,
		APIKey:      apiKey,
		APISecret:   apiSecret,
		Status:      1,
		IPWhitelist: req.IPWhitelist,
	}

	if err := db.Create(&key).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建密钥失败"})
		return
	}

	// 返回时包含 secret（仅创建时返回一次）
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":         key.ID,
			"name":       key.Name,
			"api_key":    apiKey,
			"api_secret": apiSecret,
			"status":     key.Status,
			"created_at": key.CreatedAt,
		},
	})
}

// DeleteAppAPIKey 删除API密钥
func DeleteAppAPIKey(c *gin.Context) {
	appIDParam := c.Param("id")
	keyIDParam := c.Param("keyId")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	keyID, err := strconv.ParseUint(keyIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的密钥ID"})
		return
	}

	if err := db.Where("id = ? AND app_id = ?", keyID, appID).
		Delete(&model.AppAPIKey{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除密钥失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// UpdateAppAPIKeyStatus 更新API密钥状态
func UpdateAppAPIKeyStatus(c *gin.Context) {
	appIDParam := c.Param("id")
	keyIDParam := c.Param("keyId")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	keyID, err := strconv.ParseUint(keyIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的密钥ID"})
		return
	}

	var req struct {
		Status int8 `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := db.Model(&model.AppAPIKey{}).
		Where("id = ? AND app_id = ?", keyID, appID).
		Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// ==================== API调用统计 ====================

// GetAppAPIStats 获取APP的API调用统计
func GetAppAPIStats(c *gin.Context) {
	appIDParam := c.Param("id")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	// 获取最近7天的统计
	startTime := time.Now().AddDate(0, 0, -7)

	var stats []model.APICallStats
	if err := db.Where("app_id = ? AND stat_hour >= ?", appID, startTime).
		Order("stat_hour DESC").
		Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取统计失败"})
		return
	}

	// 汇总统计
	var totalCalls, successCalls, failCalls int
	for _, s := range stats {
		totalCalls += s.TotalCalls
		successCalls += s.SuccessCalls
		failCalls += s.FailCalls
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"summary": gin.H{
				"total_calls":   totalCalls,
				"success_calls": successCalls,
				"fail_calls":    failCalls,
				"success_rate":  float64(successCalls) / float64(max(totalCalls, 1)) * 100,
			},
			"details": stats,
		},
	})
}

// GetAppAPICallLogs 获取APP的API调用日志
func GetAppAPICallLogs(c *gin.Context) {
	appIDParam := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	apiCode := c.Query("api_code")
	status := c.Query("status")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	query := db.Model(&model.APICallLog{}).Where("app_id = ?", appID)

	if apiCode != "" {
		query = query.Where("api_code = ?", apiCode)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var logs []model.APICallLog
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取日志失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": logs,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
