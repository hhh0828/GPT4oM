package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Request struct {
	Model       string     `json:"model"`
	Messages    []Messages `json:"messages"`
	Temperature float32    `json:"temperature"`
}

func (r *Request) InputUserChat() *Request {
	fmt.Println("Hi Please input a prompt you want")
	var messageslice []Messages
	var content string
	reader := bufio.NewReader(os.Stdin)
	content, _ = reader.ReadString('\n')
	messageslice = append(messageslice, Messages{"user", content})
	r.Messages = messageslice

	return r
}

func (userreq *Request) MakingRequest() *http.Request {
	apikey := "sk-None-izmLhx0PUGxalUxl4RaRT3BlbkFJmSVSdLfv7ypsP6U036hH"
	url := "https://api.openai.com/v1/chat/completions"

	//part of creating go object regarding the user request.

	userreq.Model = "gpt-4o-mini-2024-07-18"
	userreq.Temperature = 0.7

	//요청을 따로 빼는게 좋음.. Makingrequest는 재사용우려가 있음...
	userreq.InputUserChat()

	//marshal the user request so the other lang can understand
	ur, _ := json.Marshal(userreq)
	//fmt.Println(string(ur))
	//create as http request.
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(ur))

	//header set
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apikey)
	return req
}

// when get a request from user, so my server can send request again to API server transforming the Json object that GPT API Server has set.
func (userreq *Request) Repackage() *http.Request {
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
	return req
}

func (userreq *Request) RetrieveSelf() string {

	var result strings.Builder
	for _, Message := range userreq.Messages {
		result.WriteString(Message.Content + "\n")
	}

	res := result.String() + "모델 정보1 : " + userreq.Model + "\n" + fmt.Sprintf("%f", userreq.Temperature)
	return res
}
