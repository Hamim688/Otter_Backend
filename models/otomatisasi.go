package models

type Otomatisasi struct {
	ID                   uint    `gorm:"primaryKey` // Selalu ID 1
	ModeAutoLampu        bool    `gorm:"default:true"`
	AutoLampuKamar       bool    `gorm:"default:true"`
	AutoLampuTamu        bool    `gorm:"default:true"`
	AutoLampuKamarMandi  bool    `gorm:"default:true"`
	AutoLampuDapur       bool    `gorm:"default:true"`
	ModeAutoKipas        bool    `gorm:"default:true"`
	BatasGelapLampu      int     `gorm:"default:30"`
	BatasPanasKamar      float64 `gorm:"default:29.0"`
}