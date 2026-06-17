package models

type Notification struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Category  string `json:"category"`
	Priority  string `json:"priority"`
	IsRead    bool   `gorm:"default:false" json:"is_read"`
	Timestamp string `json:"timestamp"`
}