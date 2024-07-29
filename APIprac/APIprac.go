package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	//URL
	url := "https://openapi.gooroomee.com/api/v1/room"

	//요청 구조체 생성
	type Payload struct {
		CallType         string `json:"callType"`
		LiveMode         bool   `json:"liveMode"`
		MaxJoinCount     int    `json:"maxJoinCount"`
		LiveMaxJoinCount int    `json:"liveMaxJoinCount"`
		LayoutType       int    `json:"layoutType"`
		SfuIncludeAll    bool   `json:"sfuIncludeAll"`
	}
	//요청 객체 생성
	Pl := new(Payload)
	//요청 입력
	fmt.Println("what is callTpye? : ")
	fmt.Scanln(&Pl.CallType)
	fmt.Println("what is a mode you want choose, live or dead")
	fmt.Scanln(&Pl.LiveMode)
	fmt.Println("how many people you want them to join maximumly")
	fmt.Scanln(&Pl.MaxJoinCount)
	fmt.Println("LiveMaxJoinCount : ?")
	fmt.Scanln(&Pl.LiveMaxJoinCount)
	fmt.Println("Layout type : ?")
	fmt.Scanln(&Pl.LayoutType)
	fmt.Println("sfuIncludeAll : ?")
	fmt.Scanln(&Pl.SfuIncludeAll)

	Pl.MaxJoinCount = 16
	// 요청 Marshal Go object to Json data
	sendingdata, _ := json.Marshal(Pl)

	//payload := strings.NewReader("callType=P2P&liveMode=false&maxJoinCount=4&liveMaxJoinCount=100&layoutType=4&sfuIncludeAll=true")
	// 요청 객체 생성
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(sendingdata))

	//요청 Header 추가
	req.Header.Add("accept", "application/json")
	//요청 Header 추가
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	//요청 Header Auth Token 입력 (API Key)
	req.Header.Add("X-GRM-AuthToken", "12056163501988613cf51b7b51cdd8140bb172761d02211a8b")

	//응답 객체 생성
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	//body, _ := io.ReadAll(res.Body)

	type Room struct {
		EndDate      string `json:"endDate"` // endData -> EndDate (필드 이름 대문자)
		MaxJoinCount int    `json:"maxJoinCount"`
		RoomID       string `json:"roomId"` // roomID -> RoomID
	}

	type Data struct {
		Room Room `json:"room"` // room -> Room
	}

	type Response struct {
		ResultCode  string `json:"resultCode"`
		Description string `json:"description"`
		Data        Data   `json:"data"`
	}

	Respon := &Response{}
	json.NewDecoder(res.Body).Decode(Respon)

	fmt.Println(Respon.Data.Room.MaxJoinCount)
	fmt.Println(Respon)
	//fmt.Println(string(body))

}
