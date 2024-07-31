package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Request struct {
	Model       string     `json:"model"`
	Messages    []Messages `json:"messages"`
	Temperature float32    `json:"temperature"`
}

type Stream struct {
	StreamModel    string
	StreamMessages []Messages
	Stream         bool
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Usage   Usage     `json:"usage"`
	Choices []Choices `json:"choices"`
}

type Usage struct {
	PromptToken     int   `json:"prompt_tokens"`
	CompletionToken int   `json:"completion_tokens"`
	TotalToken      int64 `json:"total_tokens"`
}

type Choices struct {
	Messages     Messages `json:"message"`
	Logprobs     any      `json:"logprobs"`
	Finishreason string   `json:"finish_reason"`
	IDX          int      `json:"index"`
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

//get asks from User on the WEB site. handler is standbyfor .

type Making interface {
	MakingRequest()
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
	fmt.Println(string(ur))
	//create as http request.
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(ur))

	//header set
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apikey)
	return req
}

// Raw data from GPT API server
func GetResponse(req *http.Request) *http.Response {
	//client := new(http.Client)
	client2 := http.DefaultClient
	res, _ := client2.Do(req)
	//fmt.Println(res)
	//data, _ := io.ReadAll(res.Body)
	//datas := bytes.NewBuffer(data)
	//fmt.Println(datas)

	return res
}

// Transform the gpt data to readable base.
func TransformRes(rawres *http.Response, gres *Response) *Response {

	json.NewDecoder(rawres.Body).Decode(gres)

	//fmt.Println(gres)
	return gres
}

func PrintResponse(res *Response) {

	fmt.Println(res)

}

// get userinpput with json object.
// need to transform the data to go object so the g ocompliler can read.
type Userinput struct {
	Request string `json:"prompt"`
}

type GPToutput struct {
	Output string `json:"response"`
}

// part of transforming the data to go object.

func RequestHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", UserinputHandler)
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)
	return mux
}

// JS input data from user handler.
func UserinputHandler(w http.ResponseWriter, r *http.Request) {
	//get a chat data from Userinput.

	//need to change the data to the Go object with decoder.
	var Uinput Userinput

	json.NewDecoder(r.Body).Decode(&Uinput)
	//get data and transfer the data with the json code before transfering you must do check if it fit to json type API request for GPT api.
	var messageslice []Messages
	messageslice = append(messageslice, Messages{"user", Uinput.Request})
	fmt.Println(Uinput.Request)
	//Create reqeust that fit to Json.

	//transferring the inputdata that hadnled by upper code. and get response
	req := func() *http.Request {
		userreq := new(Request)
		userreq.Messages = messageslice
		apikey := "sk-None-izmLhx0PUGxalUxl4RaRT3BlbkFJmSVSdLfv7ypsP6U036hH"
		url := "https://api.openai.com/v1/chat/completions"

		//part of creating go object regarding the user request.

		userreq.Model = "gpt-4o-mini-2024-07-18"
		userreq.Temperature = 0.7

		//marshal the user request so the other lang can understand
		ur, _ := json.Marshal(userreq)
		fmt.Println(string(ur))
		//create as http request.
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(ur))

		//header set
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apikey)
		return req

	}
	rawres := GetResponse(req())
	transformedd := TransformRes(rawres, &Response{})
	//바디에 있는 데이터만 UI쪽으로 전달.
	outputt := new(GPToutput)
	outputt.Output = transformedd.Choices[0].Messages.Content
	w.Header().Set("Content-Type", "application/json")
	sdata, _ := json.Marshal(outputt)
	w.Write(sdata)
}

/*
func SavingResponse(r *Response) []Messages {

}
*/
func main() {

	// / 경로로 들어오는 모든 요청을 ./static 디렉토리의 index.html 파일로 라우팅

	http.ListenAndServe(":8080", RequestHandler())
	userreq := new(Request)
	req := userreq.MakingRequest()
	rawres := GetResponse(req)

	//fmt.Println(rawres)
	res := new(Response)
	transfromedres := TransformRes(rawres, res)

	PrintResponse(transfromedres)

}

///need to implement history feature so it works as gpt streaming server.
///Server console interface creating requried.
