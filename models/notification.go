package models

type Notification struct {
	ID        string `gorm:"primaryKey"`
	Title     string
	Message   string
	Category  string
	Priority  string
	IsRead    bool   `gorm:"default:false"`
	Timestamp string 
}