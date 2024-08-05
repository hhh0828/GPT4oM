package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool) // 활성 클라이언트 목록
	broadcast = make(chan string)              // 브로드캐스트 채널
	mu        sync.Mutex                       // 동시성 제어를 위한 뮤텍스
)

func handleMessages() {
	for {

		msg := <-broadcast

		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println("메시지 전송 오류:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()

	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {

	//웹소켓 업그레이더 웹소켓을 상시 업그레이드 해줌
	upgrader := new(websocket.Upgrader)
	upgrader.ReadBufferSize = 1024
	upgrader.WriteBufferSize = 1024
	//conn객체 생성

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	for {
		req := new(Messages)
		conn.ReadJSON(&req)
		broadcast <- req.Content
		break
	}

}
