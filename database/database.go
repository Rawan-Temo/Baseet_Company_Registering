package database

import (
	"log"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	// Connect with logger enabled for info level
	db, err := gorm.Open(sqlite.Open("apiGo2.db?_foreign_keys=on"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("❌ Failed to connect to DB:", err)
	}

	log.Println("✅ Connected to DB successfully")

	// Run migrations

	log.Println("🧩 Running migrations...")
	// db.Migrator().DropTable(&models.User{})
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	DB = db
	return db
}
