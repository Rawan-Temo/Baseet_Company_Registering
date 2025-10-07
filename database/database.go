package database

import (
	"log"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/models"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
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
	db.Migrator().DropTable(&models.User{})
	if err := db.AutoMigrate(&models.User{} ,&models.Company{}); err != nil {
		log.Fatal("❌ Migration failed:", err)
	}
	log.Println("✅ Migration completed")
	if err := utils.CreateDefaultAdmin(db); err != nil {
		log.Fatal("❌ Creating default admin user failed:", err)
	}
	log.Println("✅ Default admin user ensured")



	DB = db
	return db
}
