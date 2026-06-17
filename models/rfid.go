package models

type RfidCard struct {
	UID         string `gorm:"primaryKey;column:uid" json:"uid"`
	NamaPemilik string `gorm:"column:nama_pemilik" json:"nama_pemilik"`
	Status      string `gorm:"column:status;default:'menunggu'" json:"status"` 
}