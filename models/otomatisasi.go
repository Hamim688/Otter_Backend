package models

type Otomatisasi struct {
	ID                   uint    `gorm:"primaryKey;default:1" json:"id"` // Selalu ID 1
	ModeAutoLampu        bool    `gorm:"default:true" json:"mode_auto_lampu"`
	AutoLampuKamar       bool    `gorm:"default:true" json:"auto_lampu_kamar"`
	AutoLampuTamu        bool    `gorm:"default:true" json:"auto_lampu_tamu"`
	AutoLampuKamarMandi  bool    `gorm:"default:true" json:"auto_lampu_kamar_mandi"`
	AutoLampuDapur       bool    `gorm:"default:true" json:"auto_lampu_dapur"`
	ModeAutoKipas        bool    `gorm:"default:true" json:"mode_auto_kipas"`
	BatasGelapLampu      int     `gorm:"default:30" json:"batas_gelap_lampu"`
	BatasPanasKamar      float64 `gorm:"default:29.0" json:"batas_panas_kamar"`
}