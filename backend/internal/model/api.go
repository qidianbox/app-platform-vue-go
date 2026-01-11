package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemAPI 系统API定义模型
type SystemAPI struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	Code           string         `gorm:"size:100;not null;uniqueIndex" json:"code"`
	Path           string         `gorm:"size:255;not null" json:"path"`
	Method         string         `gorm:"size:20;not null" json:"method"`
	ModuleCode     string         `gorm:"size:100;index" json:"module_code"`
	Category       string         `gorm:"size:100;index" json:"category"`
	Description    string         `gorm:"type:text" json:"description"`
	RequestParams  string         `gorm:"type:json" json:"request_params"`
	ResponseParams string         `gorm:"type:json" json:"response_params"`
	Version        string         `gorm:"size:20;default:v1" json:"version"`
	Status         int8           `gorm:"default:1" json:"status"`
	IsPublic       int8           `gorm:"default:0" json:"is_public"`
	RateLimit      int            `gorm:"default:0" json:"rate_limit"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SystemAPI) TableName() string {
	return "system_apis"
}

// AppAPIPermission APP API授权模型
type AppAPIPermission struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	AppID     uint       `gorm:"not null;index" json:"app_id"`
	APIID     uint       `gorm:"column:api_id;not null;index" json:"api_id"`
	APICode   string     `gorm:"size:100;not null" json:"api_code"`
	Status    int8       `gorm:"default:1" json:"status"`
	RateLimit int        `gorm:"default:0" json:"rate_limit"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (AppAPIPermission) TableName() string {
	return "app_api_permissions"
}

// AppAPIKey APP API密钥模型
type AppAPIKey struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	AppID       uint           `gorm:"not null;index" json:"app_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	APIKey      string         `gorm:"size:64;not null;uniqueIndex" json:"api_key"`
	APISecret   string         `gorm:"size:128;not null" json:"-"` // 不返回给前端
	Status      int8           `gorm:"default:1" json:"status"`
	Permissions string         `gorm:"type:json" json:"permissions"`
	IPWhitelist string         `gorm:"type:text" json:"ip_whitelist"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AppAPIKey) TableName() string {
	return "app_api_keys"
}

// APICallLog API调用日志模型
type APICallLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AppID         uint      `gorm:"not null;index" json:"app_id"`
	APIID         uint      `gorm:"column:api_id" json:"api_id"`
	APICode       string    `gorm:"size:100;index" json:"api_code"`
	APIKeyID      uint      `gorm:"column:api_key_id" json:"api_key_id"`
	RequestMethod string    `gorm:"size:20;not null" json:"request_method"`
	RequestPath   string    `gorm:"size:500;not null" json:"request_path"`
	RequestParams string    `gorm:"type:text" json:"request_params"`
	RequestBody   string    `gorm:"type:text" json:"request_body"`
	ResponseCode  int       `gorm:"default:0" json:"response_code"`
	ResponseBody  string    `gorm:"type:text" json:"response_body"`
	ClientIP      string    `gorm:"size:50" json:"client_ip"`
	UserAgent     string    `gorm:"size:500" json:"user_agent"`
	Duration      int       `gorm:"default:0" json:"duration"`
	Status        int8      `gorm:"default:1;index" json:"status"`
	ErrorMessage  string    `gorm:"type:text" json:"error_message"`
	CreatedAt     time.Time `gorm:"index" json:"created_at"`
}

func (APICallLog) TableName() string {
	return "api_call_logs"
}

// APICallStats API调用统计模型
type APICallStats struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	AppID        uint      `gorm:"not null" json:"app_id"`
	APICode      string    `gorm:"size:100;not null" json:"api_code"`
	StatHour     time.Time `gorm:"not null" json:"stat_hour"`
	TotalCalls   int       `gorm:"default:0" json:"total_calls"`
	SuccessCalls int       `gorm:"default:0" json:"success_calls"`
	FailCalls    int       `gorm:"default:0" json:"fail_calls"`
	AvgDuration  int       `gorm:"default:0" json:"avg_duration"`
	MaxDuration  int       `gorm:"default:0" json:"max_duration"`
	MinDuration  int       `gorm:"default:0" json:"min_duration"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (APICallStats) TableName() string {
	return "api_call_stats"
}
