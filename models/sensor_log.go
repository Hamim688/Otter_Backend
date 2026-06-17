package models

import "time"

type SensorLog struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CahayaAtap      int       `gorm:"column:cahaya_atap" json:"cahaya_atap"`
	DapurSuhu       float64   `gorm:"column:dapur_suhu" json:"dapur_suhu"`
	DapurKelembapan float64   `gorm:"column:dapur_kelembapan" json:"dapur_kelembapan"`
	DapurFlame      int       `gorm:"column:dapur_flame" json:"dapur_flame"`
	KamarSuhu       float64   `gorm:"column:kamar_suhu" json:"kamar_suhu"`
	KamarKelembapan float64   `gorm:"column:kamar_kelembapan" json:"kamar_kelembapan"`
	TamuGerak       bool      `gorm:"column:tamu_gerak" json:"tamu_gerak"`
	CreatedAt       time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}
