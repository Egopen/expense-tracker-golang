package initializer

import models "main/Models"

func Sync() {
	DB.AutoMigrate(&models.Users{}, &models.Categories{}, &models.Expenses{})
}
