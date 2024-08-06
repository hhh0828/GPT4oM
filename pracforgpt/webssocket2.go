package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	//wg sync.WaitGroup
	clients   = make(map[*websocket.Conn]bool) // 활성 클라이언트 목록
	broadcast = make(chan string)              // 브로드캐스트 채널

	userchat = make(chan SendDatatoGO)
)

type Jrequest struct {
	Msg string `json:"data"`
}

func handleMessages() {

	for {

		usermsg := <-userchat

		//fmt.Println("chan working", msg)

		for client := range clients {
			//fmt.Println(client)

			//a := client.RemoteAddr().String()
			//a := client.NetConn().RemoteAddr().String()
			//a1 := client.LocalAddr().String()
			clientAddr := client.NetConn().RemoteAddr().String()

			//fmt.Println("thisis checking for ip", clientAddr)
			//fmt.Println("thisis from usermsg.IPaddr", usermsg.IPaddr)
			if !isSameIP(clientAddr, usermsg.IPaddr) {
				err1 := client.WriteMessage(websocket.TextMessage, []byte(usermsg.IPaddr+" 님의 채팅 : "+usermsg.Msg))
				if err1 != nil {
					fmt.Println("메시지 전송 오류:", err1, client.NetConn().RemoteAddr().String())
					client.Close()
					delete(clients, client)
				}
			}
			err := client.WriteMessage(websocket.TextMessage, []byte(usermsg.Gres))
			if err != nil {
				fmt.Println("메시지 전송 오류:", err, client.NetConn().RemoteAddr().String())
				client.Close()
				delete(clients, client)
			}
		}

	}

}
func isSameIP(clientAddr, usermsgAddr string) bool {
	clientIP, _, err := net.SplitHostPort(clientAddr)
	if err != nil {
		return false
	}
	usermsgIP, _, err := net.SplitHostPort(usermsgAddr)
	if err != nil {
		return false
	}

	return clientIP == usermsgIP
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {

	//웹소켓 업그레이더 웹소켓을 상시 업그레이드 해줌
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// 모든 출처에서 연결을 허용
			return true
		},
	}
	upgrader.ReadBufferSize = 1024
	upgrader.WriteBufferSize = 1024

	//conn객체 생성

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Conn객체 만들어짐")
	//fmt.Println("아ㅣ거는 테스트이다 ㅇtest")
	defer conn.Close()

	clients[conn] = true

	for {
		req := new(Messages)
		conn.ReadJSON(&req)
		//fmt.Println(req)
		broadcast <- req.Content

		break
	}

}
