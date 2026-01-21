package database

import (
	"log"

	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	general_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/general"
	media_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/media"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	// Connect with logger enabled for info level
	dsn := "host=localhost user=postgres  password=rawan445153 dbname=testdb port=5432 sslmode=disable "
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("❌ Failed to connect to DB:", err)
	}

	log.Println("✅ Connected to DB successfully")

	// Run migrations

	log.Println("🧩 Running migrations...")
	if err := db.AutoMigrate(&auth_models.User{},
		&company_models.Company{},
		&general_models.CompanyType{},
		&general_models.Office{},
		&auth_models.License{},
		&company_models.People{},
		&company_models.CompanyActivity{},
		&company_models.CompanyActivity{},
		&media_models.MediaType{}); err != nil {
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
