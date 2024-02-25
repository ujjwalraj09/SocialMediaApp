package initializers

import (
	"go/src/ujjwal/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MigrateDatabase() {
	dsn := os.Getenv("DB")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Auto Migrate the models
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Like{}, &models.PostImage{})
	if err != nil {
		panic("failed to migrate database")
	}

	DB = db
}
