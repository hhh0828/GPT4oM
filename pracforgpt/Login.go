package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ID struct {
	IDX       int
	CreatedAt time.Time
	User
	ID       string
	Password string
}

type User struct {
	IDX     int
	Name    string
	Email   string
	Address string
}

// 0817 study
//go get -u gorm.io/gorm
// go get -u gorm.io/driver/mysql # MySQL 드라이버 예시

func ConnDB() {
	dsn := "host=localhost user=postgres password=root1234 port=8801 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	fmt.Println("PostgreSQL에 성공적으로 연결되었습니다!")
	newDBName := "testDB"
	sqlnewdb := fmt.Sprintf("CREATE DATABASE %s;", newDBName)
	db.Exec(sqlnewdb)
	db.AutoMigrate(&ID{})

	//err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

}
