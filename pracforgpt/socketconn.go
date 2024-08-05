package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func Sockethandler(w http.ResponseWriter, r *http.Request) {

	//웹소켓 업그레이더 웹소켓을 상시 업그레이드 해줌
	upgrader := new(websocket.Upgrader)
	upgrader.ReadBufferSize = 1024
	upgrader.WriteBufferSize = 1024
	//conn객체 생성

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		req := new(Request)
		conn.ReadJSON(&req)
		if err != nil {
			log.Println(err)
			return
		}
		println(req.Messages)
		usertemp := ResCollector[len(ResCollector)-1]
		if err := conn.WriteJSON(usertemp); err != nil {
			log.Println(err)
			return
		}
		temp := ResCollector[len(ResCollector)]

		if err := conn.WriteJSON(temp); err != nil {
			log.Println(err)
			return
		}

	}

}
