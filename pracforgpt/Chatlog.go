package main

import "log"

type Userchat struct {
	ID             uint `gorm:"primaryKey;autoIncrement"`
	RequestContent string
	UserID         uint `gorm:"not null"`
	User           User `gorm:"foreignKey:UserID"`
}

func LogUserChat(ur *Userchat) {
	db, err := ConnDB()
	if err != nil {
		log.Fatal("error occured", err)
	}

	db.AutoMigrate(&Userchat{})
	db.Create(&ur)
}
