package controllers

import (
	"Backend/config"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

// GetProfile mengambil profil admin (ID 1) dari PostgreSQL
func GetProfile(c *fiber.Ctx) error {
	var user models.User
	
	// Cari user dengan ID 1. Jika belum ada, gunakan default seeder
	if err := config.DB.First(&user, 1).Error; err != nil {
		// Fallback jika database seeder belum jalan
		user = models.User{
			ID:          1,
			Username:    "admin",
			Password:    "1234",
			DisplayName: "Mimah Dudim",
			Role:        "Administrator Rumah Pintar",
			AvatarURL:   "",
		}
		config.DB.Create(&user)
	}

	return c.JSON(user)
}

// UpdateProfile memperbarui profil admin di PostgreSQL
func UpdateProfile(c *fiber.Ctx) error {
	var user models.User

	if err := config.DB.First(&user, 1).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	type ProfileBody struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
		Role        string `json:"role"`
		AvatarURL   string `json:"avatar_url"`
	}

	var body ProfileBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format JSON salah!"})
	}

	user.Username = body.Username
	user.Password = body.Password
	user.DisplayName = body.DisplayName
	user.Role = body.Role
	user.AvatarURL = body.AvatarURL

	config.DB.Save(&user)

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Profil berhasil diperbarui di database!",
		"data":    user,
	})
}
