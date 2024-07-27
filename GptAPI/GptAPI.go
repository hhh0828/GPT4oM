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
func (r GRequest) Marshal() *http.Request {
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

//API 요청 request 구현

//응답 Structure Marshal 후 Go object로 전환.

//답변 Get.

func main() {

}
