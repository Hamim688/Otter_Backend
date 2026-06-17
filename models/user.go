package models

type User struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Username    string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	DisplayName string
	Role        string
	AvatarURL   string
}