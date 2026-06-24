package controllers

import (
	"Backend/config"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

// Ambil semua kartu RFID terdaftar
func GetRfidCards(c *fiber.Ctx) error {
	var cards []models.RfidCard
	if err := config.DB.Find(&cards).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data RFID"})
	}
	return c.JSON(cards)
}

// Tambah kartu RFID baru secara manual
func CreateRfidCard(c *fiber.Ctx) error {
	var card models.RfidCard
	if err := c.BodyParser(&card); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format JSON salah!"})
	}
	if card.UID == "" || card.NamaPemilik == "" {
		return c.Status(400).JSON(fiber.Map{"error": "UID dan Nama Pemilik wajib diisi!"})
	}

	var existing models.RfidCard
	if err := config.DB.First(&existing, "uid = ?", card.UID).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "UID RFID sudah terdaftar!"})
	}

	card.Status = "aktif"
	if err := config.DB.Create(&card).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan kartu RFID baru"})
	}

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Kartu RFID berhasil ditambahkan secara manual!",
		"data":    card,
	})
}

// Setujui pendaftaran kartu RFID baru (Ubah Unknown -> Nama Pemilik asli, Status pending -> aktif)
func ApproveRfidCard(c *fiber.Ctx) error {
	uid := c.Params("uid")
	var card models.RfidCard

	if err := config.DB.First(&card, "uid = ?", uid).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kartu RFID tidak ditemukan"})
	}

	type ApproveBody struct {
		NamaPemilik string `json:"nama_pemilik"`
	}
	var body ApproveBody
	if err := c.BodyParser(&body); err != nil || body.NamaPemilik == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Nama pemilik wajib diisi!"})
	}

	card.NamaPemilik = body.NamaPemilik
	card.Status = "aktif"
	config.DB.Save(&card)



	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Pendaftaran kartu RFID disetujui!",
		"data":    card,
	})
}

// Aktifkan atau Nonaktifkan kartu RFID
func UpdateRfidStatus(c *fiber.Ctx) error {
	uid := c.Params("uid")
	var card models.RfidCard

	if err := config.DB.First(&card, "uid = ?", uid).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kartu RFID tidak ditemukan"})
	}

	type StatusBody struct {
		Status string `json:"status"` // 'aktif' atau 'nonaktif'
	}
	var body StatusBody
	if err := c.BodyParser(&body); err != nil || (body.Status != "aktif" && body.Status != "nonaktif") {
		return c.Status(400).JSON(fiber.Map{"error": "Status harus 'aktif' atau 'nonaktif'!"})
	}

	card.Status = body.Status
	config.DB.Save(&card)

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Status kartu RFID berhasil diupdate!",
		"data":    card,
	})
}

// Tolak pendaftaran kartu pending atau hapus kartu terdaftar
func DeleteRfidCard(c *fiber.Ctx) error {
	uid := c.Params("uid")
	var card models.RfidCard

	if err := config.DB.First(&card, "uid = ?", uid).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kartu RFID tidak ditemukan"})
	}

	config.DB.Delete(&card)

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Kartu RFID berhasil dihapus/ditolak!",
	})
}
