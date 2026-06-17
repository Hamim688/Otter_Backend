package controllers

import (
	"Backend/config"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

// Ambil data profil administrator
func GetProfile(c *fiber.Ctx) error {
	var user models.User
	
	// Cari user ID 1. Jika belum ada, cari user pertama di database
	if err := config.DB.First(&user, 1).Error; err != nil {
		if err := config.DB.First(&user).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Profil pengguna belum di-seed!"})
		}
	}

	return c.JSON(user)
}

// Update data profil administrator
func UpdateProfile(c *fiber.Ctx) error {
	var user models.User

	// Ambil data user ID 1
	if err := config.DB.First(&user, 1).Error; err != nil {
		if err := config.DB.First(&user).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Profil tidak ditemukan!"})
		}
	}

	type ProfileUpdateBody struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
		Role        string `json:"role"`
		AvatarURL   string `json:"avatar_url"`
	}

	var body ProfileUpdateBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format JSON salah!"})
	}

	if body.Username != "" {
		user.Username = body.Username
	}
	if body.Password != "" {
		user.Password = body.Password
	}
	if body.DisplayName != "" {
		user.DisplayName = body.DisplayName
	}
	if body.Role != "" {
		user.Role = body.Role
	}
	// Always set avatar URL if provided (or can be empty string)
	user.AvatarURL = body.AvatarURL

	config.DB.Save(&user)

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Profil berhasil diperbarui di database!",
		"data":    user,
	})
}
