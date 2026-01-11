package model

import (
	"time"

	"gorm.io/gorm"
)

// AppMenu APP菜单模型
type AppMenu struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	AppID      uint           `gorm:"not null;index" json:"app_id"`
	ParentID   uint           `gorm:"default:0" json:"parent_id"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	Code       string         `gorm:"size:100;not null" json:"code"`
	Icon       string         `gorm:"size:100" json:"icon"`
	Path       string         `gorm:"size:255" json:"path"`
	Component  string         `gorm:"size:255" json:"component"`
	MenuType   int8           `gorm:"default:1" json:"menu_type"` // 1-目录 2-菜单 3-按钮
	Visible    int8           `gorm:"default:1" json:"visible"`   // 0-隐藏 1-显示
	Status     int8           `gorm:"default:1" json:"status"`    // 0-禁用 1-启用
	SortOrder  int            `gorm:"default:0" json:"sort_order"`
	Permission string         `gorm:"size:255" json:"permission"`
	Remark     string         `gorm:"size:500" json:"remark"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Children   []*AppMenu     `gorm:"-" json:"children,omitempty"`
}

func (AppMenu) TableName() string {
	return "app_menus"
}

// MenuTree 构建菜单树
func BuildMenuTree(menus []AppMenu, parentID uint) []*AppMenu {
	var tree []*AppMenu
	for i := range menus {
		if menus[i].ParentID == parentID {
			menu := &menus[i]
			menu.Children = BuildMenuTree(menus, menu.ID)
			tree = append(tree, menu)
		}
	}
	return tree
}
