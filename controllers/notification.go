package controllers

import (
	"Backend/config"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

// Ambil semua notifikasi, urutkan dari yang paling baru
func GetNotifications(c *fiber.Ctx) error {
	var list []models.Notification
	if err := config.DB.Order("timestamp desc").Find(&list).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data notifikasi"})
	}
	return c.JSON(list)
}

// Tandai satu notifikasi telah dibaca
func MarkAsRead(c *fiber.Ctx) error {
	id := c.Params("id")
	var notif models.Notification

	if err := config.DB.First(&notif, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notifikasi tidak ditemukan"})
	}

	notif.IsRead = true
	config.DB.Save(&notif)

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Notifikasi ditandai dibaca",
		"data":    notif,
	})
}

// Tandai semua notifikasi telah dibaca
func MarkAllAsRead(c *fiber.Ctx) error {
	if err := config.DB.Model(&models.Notification{}).Where("is_read = ?", false).Update("is_read", true).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menandai semua notifikasi"})
	}

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Semua notifikasi ditandai dibaca",
	})
}

// Hapus satu notifikasi
func DeleteNotification(c *fiber.Ctx) error {
	id := c.Params("id")
	var notif models.Notification

	if err := config.DB.First(&notif, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notifikasi tidak ditemukan"})
	}

	config.DB.Delete(&notif)

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Notifikasi berhasil dihapus",
	})
}

// Hapus seluruh notifikasi (Clear All)
func ClearAllNotifications(c *fiber.Ctx) error {
	if err := config.DB.Exec("DELETE FROM notifications").Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membersihkan notifikasi"})
	}

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Seluruh riwayat notifikasi dibersihkan",
	})
}
