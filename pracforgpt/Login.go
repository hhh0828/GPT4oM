package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:100"`
	Email     string `gorm:"uniqueIndex;size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 0817 study
//go get -u gorm.io/gorm
// go get -u gorm.io/driver/mysql # MySQL 드라이버 예시

func InitDB() {
	dsn := "host=localhost user=postgres password=root1234 port=8801 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	fmt.Println("PostgreSQL에 성공적으로 연결되었습니다!")
	newDBName := "testDB"
	sqlnewdb := fmt.Sprintf("CREATE DATABASE %s;", newDBName)
	db.Exec(sqlnewdb)
	db.AutoMigrate(&User{})

	//err = db.AutoMigrate(&User{})

}

func ConnDB() (db *gorm.DB, er error) {
	dsn := "host=localhost user=postgres password=root1234 dbname=testdb port=8801 sslmode=disable"
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	fmt.Println("PostgreSQL에 성공적으로 연결되었습니다!")
	return db, er

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

}
