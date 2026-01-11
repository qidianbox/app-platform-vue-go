package module

import (
	"strconv"

	"app-platform-backend/internal/model"
	"app-platform-backend/internal/pkg/database"
)

// getAppDatabaseID 根据路径参数获取 APP 的数据库 ID
// 参数可以是数字 ID 或者 app_id 字符串
func getAppDatabaseID(idParam string) (uint, error) {
	// 尝试作为数字 ID 解析
	if id, err := strconv.ParseUint(idParam, 10, 32); err == nil {
		// 验证这个 ID 是否存在
		var app model.App
		if err := database.GetDB().Select("id").First(&app, id).Error; err == nil {
			return uint(id), nil
		}
	}

	// 作为 app_id 字符串查询
	var app model.App
	if err := database.GetDB().Select("id").Where("app_id = ?", idParam).First(&app).Error; err != nil {
		return 0, err
	}

	return app.ID, nil
}
