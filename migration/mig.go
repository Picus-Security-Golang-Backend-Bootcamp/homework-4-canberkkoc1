package migration

import (
	"ck/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var err error

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "<password>"
	dbname   = "<your_DB>"
)

func InitialMigration() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	DB, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&models.Books{})

}
