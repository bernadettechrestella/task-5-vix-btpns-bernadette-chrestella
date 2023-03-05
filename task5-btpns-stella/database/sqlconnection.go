package database

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"task5-btpns-stella/helpers"
	"time"
)

func ConnectDB() *gorm.DB {
	dsn := os.Getenv("DB_CONNECTION_STRING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("DB connection failed!")
	}

	// get sql.Db for set database connection pool
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("DB connection failed!")
	}

	// set3dsxzC database conneciton pool
	sqlDb.SetMaxIdleConns(5)
	sqlDb.SetMaxOpenConns(20)
	sqlDb.SetConnMaxLifetime(60 * time.Minute)
	sqlDb.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func GetInstanceDbMock(dbMock *sql.DB) *gorm.DB {
	// d := os.Getenv("DB_CONN_STRING")

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dbMock,
	}))
	helpers.LogError("error connect to database", err, helpers.ErrFatalCallback)

	// get sql.Db for set database connection pool
	sqlDb, err := db.DB()
	helpers.LogError("error get instance for sql.DB", err, helpers.ErrErrorsCallback)

	// set database conneciton pool
	sqlDb.SetMaxIdleConns(5)
	sqlDb.SetMaxOpenConns(20)
	sqlDb.SetConnMaxLifetime(60 * time.Minute)
	sqlDb.SetConnMaxIdleTime(10 * time.Minute)
	return db
}
