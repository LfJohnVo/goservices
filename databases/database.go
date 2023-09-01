package databases

import (
	"goservices/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func Connect() error {
	var err error

	switch config.GetEnvValue("DBConnection") {
	case "mysql":
		var DATABASE_URI string = config.GetEnvValue("DBUser") + ":" + config.GetEnvValue("DBPassword") + "@tcp(" + config.GetEnvValue("DBHost") + ":" + config.GetEnvValue("DBPort") + ")/" + config.GetEnvValue("DBName") + "?charset=utf8mb4&parseTime=True&loc=Local"
		Database, err = gorm.Open(mysql.Open(DATABASE_URI), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})
	case "postgres":
		var DATABASE_URI string = "user=" + config.GetEnvValue("DBUser") + " dbname=" + config.GetEnvValue("DBName") + " host=" + config.GetEnvValue("DBHost") + " port=" + config.GetEnvValue("DBPort") + " sslmode=disable"
		Database, err = gorm.Open(postgres.Open(DATABASE_URI), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})
	case "sqlite":
		// Database, err = gorm.Open(sqlite.Open(DATABASE_URI), &gorm.Config{
		// 	SkipDefaultTransaction: true,
		// 	PrepareStmt:            true,
		// })
	}

	if err != nil {
		panic(err)
	}

	Database.Logger = logger.Default.LogMode(logger.Info)

	// db.AutoMigrate(&models.User{}, &models.PasswordReset{})

	return nil
}
