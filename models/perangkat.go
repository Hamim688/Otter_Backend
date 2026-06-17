package models

type Perangkat struct {
	ID              uint `gorm:"primaryKey;default:1"` // Selalu ID 1
	LampuKamar      bool `gorm:"default:false"`
	LampuTamu       bool `gorm:"default:false"`
	LampuKamarMandi bool `gorm:"default:false"`
	LampuDapur      bool `gorm:"default:false"`
	KipasKamar      bool `gorm:"default:false"`
	KecepatanKipas  int  `gorm:"default:255"`
	BuzzerAlrm      bool `gorm:"default:false"`
	LedMerahDapur   bool `gorm:"default:false"`
	KunciPintuRfid  bool `gorm:"default:true"`
}