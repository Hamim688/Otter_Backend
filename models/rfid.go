package models

type RfidCard struct {
	UID         string `gorm:"primaryKey;column:uid"`
	NamaPemilik string `gorm:"column:nama_pemilik"`
	Status      string `gorm:"column:status;default:'menunggu'"` 
}