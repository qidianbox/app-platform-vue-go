package user

import (
	"log"
	"time"

	"app-platform-backend/internal/response"
	"app-platform-backend/internal/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB(database *gorm.DB) {
	db = database
	log.Println("[UserAPI] Database connection initialized")
}

// AppPlatformUser 本地用户表结构（匹配实际数据库结构）
type AppPlatformUser struct {
	ID        uint       `gorm:"column:id;primaryKey" json:"id"`
	Username  string     `gorm:"column:username" json:"username"`
	Password  string     `gorm:"column:password" json:"-"`
	Email     *string    `gorm:"column:email" json:"email"`
	Phone     *string    `gorm:"column:phone" json:"phone"`
	Status    int        `gorm:"column:status;default:1" json:"status"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName 指定表名
func (AppPlatformUser) TableName() string {
	return "users"
}

// UserResponse 用户响应结构（适配前端）
type UserResponse struct {
	ID          uint       `json:"id"`
	Nickname    string     `json:"nickname"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	Status      int        `json:"status"`
	Role        string     `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

// ListRequest 用户列表请求参数
type ListRequest struct {
	AppID  string `form:"app_id"`
	Page   int    `form:"page" binding:"min=1"`
	Size   int    `form:"size" binding:"min=1,max=100"`
	Status *int   `form:"status"`
	Search string `form:"search"`
}

// List 获取用户列表
func List(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[UserAPI] List - Invalid request: %v", err)
		response.ParamError(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证分页参数
	req.Page, req.Size = validator.ValidatePagination(req.Page, req.Size)

	log.Printf("[UserAPI] List - Request: page=%d, size=%d, search=%s, app_id=%s", req.Page, req.Size, req.Search, req.AppID)

	// 查询本地users表
	query := db.Model(&AppPlatformUser{}).Where("deleted_at IS NULL")

	// 搜索（按用户名或邮箱）
	if req.Search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}

	// 状态过滤
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Printf("[UserAPI] List - Count error: %v", err)
		response.DBError(c, err)
		return
	}

	// 分页查询
	var users []AppPlatformUser
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Order("created_at DESC").Find(&users).Error; err != nil {
		log.Printf("[UserAPI] List - Query error: %v", err)
		response.DBError(c, err)
		return
	}

	// 转换为前端期望的格式
	var responseList []UserResponse
	for _, u := range users {
		email := ""
		if u.Email != nil {
			email = *u.Email
		}
		phone := ""
		if u.Phone != nil {
			phone = *u.Phone
		}
		responseList = append(responseList, UserResponse{
			ID:          u.ID,
			Nickname:    u.Username,
			Phone:       phone,
			Email:       email,
			Status:      u.Status,
			Role:        "user",
			CreatedAt:   u.CreatedAt,
			LastLoginAt: nil,
		})
	}

	log.Printf("[UserAPI] List - Found %d users, total %d", len(users), total)

	response.PageSuccess(c, responseList, total, req.Page, req.Size)
}

// Detail 获取用户详情
func Detail(c *gin.Context) {
	idStr := c.Param("id")

	// 验证ID
	id, err := validator.ValidateID(idStr)
	if err != nil {
		log.Printf("[UserAPI] Detail - Invalid user ID: %s", idStr)
		response.ParamError(c, "无效的用户ID")
		return
	}

	log.Printf("[UserAPI] Detail - Getting user %d", id)

	var user AppPlatformUser
	if err := db.Where("deleted_at IS NULL").First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[UserAPI] Detail - User %d not found", id)
			response.NotFound(c, "用户不存在")
			return
		}
		log.Printf("[UserAPI] Detail - Query error: %v", err)
		response.DBError(c, err)
		return
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	phone := ""
	if user.Phone != nil {
		phone = *user.Phone
	}

	userResponse := UserResponse{
		ID:          user.ID,
		Nickname:    user.Username,
		Phone:       phone,
		Email:       email,
		Status:      user.Status,
		Role:        "user",
		CreatedAt:   user.CreatedAt,
		LastLoginAt: nil,
	}

	response.Success(c, userResponse)
}

// UpdateStatusRequest 更新用户状态请求参数
type UpdateStatusRequest struct {
	Status int `json:"status" binding:"oneof=0 1"`
}

// UpdateStatus 更新用户状态
func UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")

	// 验证ID
	id, err := validator.ValidateID(idStr)
	if err != nil {
		log.Printf("[UserAPI] UpdateStatus - Invalid user ID: %s", idStr)
		response.ParamError(c, "无效的用户ID")
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UserAPI] UpdateStatus - Invalid request: %v", err)
		response.ParamError(c, "请求参数错误: "+err.Error())
		return
	}

	log.Printf("[UserAPI] UpdateStatus - User %d status to %d", id, req.Status)

	// 更新用户状态
	if err := db.Model(&AppPlatformUser{}).Where("id = ? AND deleted_at IS NULL", id).Update("status", req.Status).Error; err != nil {
		log.Printf("[UserAPI] UpdateStatus - Update error: %v", err)
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "用户状态更新成功")
}

// Stats 用户统计
func Stats(c *gin.Context) {
	log.Printf("[UserAPI] Stats - Getting user statistics")

	// 总用户数
	var total int64
	if err := db.Model(&AppPlatformUser{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		log.Printf("[UserAPI] Stats - Count error: %v", err)
		total = 0
	}

	// 活跃用户数（状态为1的用户）
	var active int64
	if err := db.Model(&AppPlatformUser{}).Where("deleted_at IS NULL AND status = 1").Count(&active).Error; err != nil {
		log.Printf("[UserAPI] Stats - Active count error: %v", err)
		active = 0
	}

	// 今日新增
	var todayNew int64
	today := time.Now().Format("2006-01-02")
	if err := db.Model(&AppPlatformUser{}).Where("deleted_at IS NULL AND DATE(created_at) = ?", today).Count(&todayNew).Error; err != nil {
		log.Printf("[UserAPI] Stats - Today new count error: %v", err)
		todayNew = 0
	}

	// 禁用用户数量
	var disabled int64
	if err := db.Model(&AppPlatformUser{}).Where("deleted_at IS NULL AND status = 0").Count(&disabled).Error; err != nil {
		log.Printf("[UserAPI] Stats - Disabled count error: %v", err)
		disabled = 0
	}

	log.Printf("[UserAPI] Stats - total=%d, active=%d, todayNew=%d, disabled=%d", total, active, todayNew, disabled)

	response.Success(c, gin.H{
		"total":     total,
		"active":    active,
		"today_new": todayNew,
		"disabled":  disabled,
		"normal":    active,
		"admin":     0,
	})
}
