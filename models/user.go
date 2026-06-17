package models

type User struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string `gorm:"unique;not null" json:"username"`
	Password    string `gorm:"not null" json:"password"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	AvatarURL   string `json:"avatar_url"`
}