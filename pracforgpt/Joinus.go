package main

import (
	"log"
	"net/http"
)

func Joinus(w http.ResponseWriter, r *http.Request) {
	db, err := ConnDB()
	if err != nil {
		log.Fatal("err occured ", err)
	}
	db.AutoMigrate(&User{})

	joinreq := new(User)

	joinreq.Email = "hhhcjswo@naver.com1"
	joinreq.Name = "Hyunho.Hong2"

	// Password 를 해쉬형태로... 추가하는거 해야함.
	db.Create(&joinreq)

	db.Exec("SELECT * FROM Users;")
}
