package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// 요청 Structure 구현
type GRequest struct {
	content string
	prompt  string
	token   string
	body    string
}

// Scan으로 Input 요청 구현
func (r *GRequest) ScanUserRequest() *GRequest {
	p, _ := fmt.Scanln()
	b, _ := fmt.Scanln()
	c, _ := fmt.Scanln()
	t, _ := fmt.Scanln()
	r.prompt = string(p)
	r.body = string(b)
	r.content = string(c)
	r.token = string(t)
	return r
}

// 요청 Structure Unmarshal Go object >> Json
func (r GRequest) CreateAPIrequest() *http.Request {
	greq := new(GRequest)
	Inreq := greq.ScanUserRequest()
	reqs, _ := json.Marshal(Inreq)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/engines/davinci/completions", bytes.NewBuffer(reqs))
	return req
}

//응답 Structure 구현

type Response struct {
	respon string
}

func (res *Response) Decoding(re *http.Response) {

	json.NewDecoder(re.Body).Decode(res)
}

//응답 Structure Marshal 후 Go object로 전환.

//답변 Get.

func main() {
	//빈 객체 생성
	GRequest := GRequest{}
	//req 사용자 입력 받고 API리퀘스트 생성
	req := GRequest.ScanUserRequest().CreateAPIrequest()
	//Server로 부터 응답을 받는 client 객체 생성
	client := &http.Client{}
	//request를 요청해서 client.Do(request)로 response 반환
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	res := new(Response)
	//to Go object
	res.Decoding(resp)
	//json.NewDecoder(resp.Body).Decode(res)
	fmt.Println(resp.Body)
}
