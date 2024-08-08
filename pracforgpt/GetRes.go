package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

//인터페이스 메서드 집합... Request들
//같은 기능을 가진녀석들 을 대표하는 이름같은것.

type Requests interface {
	Repackage() *http.Request
	RetrieveSelf() string
}

type Reque interface {
}

func Logwrite1(re Reque) {

}

func Printtest(r Requests) {

	println(r.Repackage())
}

// 인터페이스를 매개변수로 받는 함수
// 즉 인터페이스의 함수들을 구현하는 구조체들은 매개변수로 들어올 수 있음
// 덕타이핑을위함
func Logwrite(rs Requests) {

	/*
		Logs, _ := os.Create("userlog.txt")
		defer Logs.Close()

		Logs.WriteString(rs.RetrieveSelf())
	*/
	println(rs.RetrieveSelf())
}

// 인터페이스 >> 새로운 리퀘스트를 만드는 구조체들을 대표하는 이름.
type Reques interface {
	NewReq() *http.Request
}

// Request의 순기능 // http.Request를 반환함 // 인터페이스의 newreq()를 구현하고있음
func (re *Request) NewReq() *http.Request {

	data, _ := json.Marshal(re)
	a, _ := http.NewRequest("Post", "url", bytes.NewBufferString(string(data)))
	return a
}

// Repackaging 구현한 구조체를 매개변수로 받고 http.response를 구현함.
func GetRes(R Requests) *http.Response {
	//client := new(http.Client)
	client2 := http.DefaultClient
	//자신만의 방식으로 리패키징을 해서 req를 보냄..
	a := R.Repackage()
	res, _ := client2.Do(a)
	//fmt.Println(res)
	//data, _ := io.ReadAll(res.Body)
	//datas := bytes.NewBuffer(data)
	//fmt.Println(datas)

	return res
}

/*
//when my server get a request from user so this server  can send the data to API server. we need to Repackaging with Json format that GPT API has set.
func(userreq Request) Repackage() *http.Request {
	apikey := "sk-None-izmLhx0PUGxalUxl4RaRT3BlbkFJmSVSdLfv7ypsP6U036hH"
	url := "https://api.openai.com/v1/chat/completions"

	//part of creating go object regarding the user request.

	userreq.Model = "gpt-4o-mini-2024-07-18"
	userreq.Temperature = 0.7

	ur, _ := json.Marshal(userreq)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(ur))

	//header set
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apikey)

	var rs Requests
	rs.Repackage()
	return req

}


*/
