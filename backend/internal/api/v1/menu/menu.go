package menu

import (
	"log"
	"net/http"
	"strconv"

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
	// 尝试解析为数字
	if id, err := strconv.ParseUint(appIDParam, 10, 64); err == nil {
		return uint(id), nil
	}
	// 如果不是数字，按 app_id 字符串查询
	var app model.App
	if err := db.Where("app_id = ?", appIDParam).First(&app).Error; err != nil {
		return 0, err
	}
	return app.ID, nil
}

// GetMenus 获取APP菜单列表（树形结构）
func GetMenus(c *gin.Context) {
	appIDParam := c.Param("id")
	log.Printf("[DEBUG] GetMenus - appIDParam: %s", appIDParam)

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	var menus []model.AppMenu
	if err := db.Where("app_id = ?", appID).Order("sort_order ASC, id ASC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取菜单失败"})
		return
	}

	// 构建树形结构
	tree := model.BuildMenuTree(menus, 0)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": tree,
	})
}

// GetMenuList 获取APP菜单列表（平铺结构）
func GetMenuList(c *gin.Context) {
	appIDParam := c.Param("id")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	var menus []model.AppMenu
	if err := db.Where("app_id = ?", appID).Order("sort_order ASC, id ASC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": menus,
	})
}

// CreateMenu 创建菜单
func CreateMenu(c *gin.Context) {
	appIDParam := c.Param("id")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	var menu model.AppMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	menu.AppID = appID

	// 检查 code 是否重复
	var count int64
	db.Model(&model.AppMenu{}).Where("app_id = ? AND code = ?", appID, menu.Code).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "菜单标识已存在"})
		return
	}

	if err := db.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    menu,
	})
}

// UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context) {
	appIDParam := c.Param("id")
	menuIDParam := c.Param("menuId")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	menuID, err := strconv.ParseUint(menuIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的菜单ID"})
		return
	}

	var menu model.AppMenu
	if err := db.Where("id = ? AND app_id = ?", menuID, appID).First(&menu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单不存在"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 不允许修改 app_id 和 id
	delete(updateData, "app_id")
	delete(updateData, "id")

	if err := db.Model(&menu).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context) {
	appIDParam := c.Param("id")
	menuIDParam := c.Param("menuId")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	menuID, err := strconv.ParseUint(menuIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的菜单ID"})
		return
	}

	// 检查是否有子菜单
	var childCount int64
	db.Model(&model.AppMenu{}).Where("parent_id = ?", menuID).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请先删除子菜单"})
		return
	}

	if err := db.Where("id = ? AND app_id = ?", menuID, appID).Delete(&model.AppMenu{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// GetMenuDetail 获取菜单详情
func GetMenuDetail(c *gin.Context) {
	appIDParam := c.Param("id")
	menuIDParam := c.Param("menuId")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	menuID, err := strconv.ParseUint(menuIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的菜单ID"})
		return
	}

	var menu model.AppMenu
	if err := db.Where("id = ? AND app_id = ?", menuID, appID).First(&menu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": menu,
	})
}

// UpdateMenuSort 批量更新菜单排序
func UpdateMenuSort(c *gin.Context) {
	appIDParam := c.Param("id")

	appID, err := getAppDatabaseID(appIDParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		return
	}

	var sortData []struct {
		ID        uint `json:"id"`
		SortOrder int  `json:"sort_order"`
		ParentID  uint `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&sortData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	tx := db.Begin()
	for _, item := range sortData {
		if err := tx.Model(&model.AppMenu{}).
			Where("id = ? AND app_id = ?", item.ID, appID).
			Updates(map[string]interface{}{
				"sort_order": item.SortOrder,
				"parent_id":  item.ParentID,
			}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新排序失败"})
			return
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "排序更新成功",
	})
}
