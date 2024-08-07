package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

type ImageRequest struct {
	Image  os.File `json:"image"`
	Model  string  `json:"model"`
	Prompt string  `json:"prompt"`
	N      int     `json:"n"`
	Size   string  `json:"size"` // need to get a size value from the Frontend.
}

type ImageResponse struct {
	Created string
	Data    []Url
}

type Url struct {
	Url string
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
	//fmt.Println(string(ur))
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
	mux.HandleFunc("/save", History) // 대화기록 저장
	//mux.Handlefunc("/image, Imagehandler) // 이미지 업로드
	mux.HandleFunc("/chat", UserinputHandler)
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/createthread", CreateThread)
	mux.HandleFunc("/whatismyip", IPreturn)
	mux.HandleFunc("/ws", HandleConnections)
	return mux
}

type AccessUserinfo struct {
	IPadd string `json:"ipresponse"`
}

func IPreturn(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	Accuserinfo := new(AccessUserinfo)
	Accuserinfo.IPadd = ip
	data, _ := json.Marshal(Accuserinfo)
	w.Write(data)
}

var messageslice []Messages
var ResCollector []Messages

// Saving history handler. return newfile and
func History(w http.ResponseWriter, r *http.Request) {
	// get the history of conversation.
	fmt.Println("start saving")
	newfile, _ := os.Create("ConversationHistory.txt")
	defer newfile.Close()

	for _, v := range ResCollector {
		if v.Role == "user" {

			newfile.WriteString(v.Content + "\n")
		} else if v.Role == "assistant" {
			newfile.WriteString("GPT : " + v.Content + "\n")
		}

	}
	w.Header().Set("Content-Disposition", "attachment; filename=ConversationHistory.txt")
	w.Header().Set("Content-Type", "text/plain")
	http.ServeFile(w, r, "ConversationHistory.txt")
}

type SendDatatoGO struct {
	IPaddr string
	Msg    string
	Gres   string
}

// JS input data from user handler.
func UserinputHandler(w http.ResponseWriter, r *http.Request) {
	//get a chat data from Userinput.

	// Decode the data, and check if there's a Image []Byte.
	// if there's no Image byte continue the code.
	// if there's Image data, redirect request to /Image Handfunc.
	// as Image Handler has a feature to hanlde a prompt with Image given.

	//need to change the data to the Go object with decoder.
	var Uinput Userinput

	json.NewDecoder(r.Body).Decode(&Uinput)
	//get data and transfer the data with the json code before transfering you must do check if it fit to json type API request for GPT api.

	messageslice = CachePreviousConver(messageslice, Messages{"user", r.RemoteAddr + "님의 채팅 : " + Uinput.Request})

	//fmt.Println(Uinput.Request)
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

	//Stack Conversation at max 1000 (You and GPT)
	CachePreviouosGres(transformedd, Messages{"user", r.RemoteAddr + "님의 채팅 : " + Uinput.Request})
	CachePreviouosGres(transformedd, transformedd.Choices[0].Messages)

	var uc SendDatatoGO = SendDatatoGO{
		IPaddr: r.RemoteAddr,
		Msg:    Uinput.Request,
		Gres:   "GPT : " + transformedd.Choices[0].Messages.Content,
	}
	userchat <- uc
	//Stack previous Conver at Max 7,
	messageslice = CachePreviousConver(messageslice, transformedd.Choices[0].Messages)

	//fmt.Println("[]Messages: GPT와의 대화 : ", ResCollector)

	//바디에 있는 데이터만 UI쪽으로 전달.
	/*
			a := "GPT : " + transformedd.Choices[0].Messages.Content
		broadcast <- a
		fmt.Println(a, "여기는 inputhandler쪽")
			outputt := new(GPToutput)
			outputt.Output = transformedd.Choices[0].Messages.Content
			w.Header().Set("Content-Type", "application/json")
			sdata, _ := json.Marshal(outputt)
			w.Write(sdata)
	*/
}

func SendingChan() {

	broadcast <- ResCollector[len(ResCollector)-1].Content
}

/*
func SavingResponse(r *Response) []Messages {

}
*/

func CachePreviouosGres(res *Response, content Messages) []Messages {
	ResCollector = append(ResCollector, content)
	if len(ResCollector) == 1000 {
		ResCollector = ResCollector[1:]
	}
	return ResCollector
}

// Messages slice Que - 012,
func CachePreviousConver(messagesslice []Messages, content Messages) []Messages {

	messagesslice = append(messagesslice, content)
	if len(messagesslice) == 7 {
		messagesslice = messagesslice[1:]
	}
	//fmt.Println(messagesslice)
	//fmt.Println(len(messagesslice))
	return messagesslice
}

type Tresponse struct {
	ID            string                 `json:"id"`
	Object        string                 `json:"object"`
	CreatedAt     int                    `json:"created_at"`
	Metadata      map[string]interface{} `json:"metadata"`
	ToolResources map[string]interface{} `json:"tool_resources"`
}

type Trequest struct {
}

var UpdateDetector bool

// Thread created - but never used yet.
func CreateThread(w http.ResponseWriter, r *http.Request) {

	treq := new(Trequest)
	apikey := "sk-None-izmLhx0PUGxalUxl4RaRT3BlbkFJmSVSdLfv7ypsP6U036hH"
	url := "https://api.openai.com/v1/threads"

	//part of creating go object regarding the user request.

	//marshal the user request so the other lang can understand
	ur, _ := json.Marshal(treq)
	fmt.Println(string(ur))
	//create as http request.
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(ur))

	//header set
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apikey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	rawres := GetResponse(req)
	fmt.Println(rawres)

	Tres := Tresponse{}
	json.NewDecoder(rawres.Body).Decode(&Tres)
	fmt.Println(Tres)

}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.RemoteAddr + "IndexHandler 로부터 Print 된 IP 주소")
	http.ServeFile(w, r, "./static/index.html")
}

func main() {

	// / 경로로 들어오는 모든 요청을 ./static 디렉토리의 index.html 파일로 라우팅
	go handleMessages()
	http.ListenAndServe("0.0.0.0:8080", RequestHandler())
	userreq := new(Request)
	req := userreq.MakingRequest()
	rawres := GetResponse(req)
	//part of Docker container response.
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	e.Logger.Fatal(e.Start(":" + httpPort))
	//fmt.Println(rawres)
	res := new(Response)
	transfromedres := TransformRes(rawres, res)

	//
	PrintResponse(transfromedres)
	r := http.Request{}
	r.FormFile("image")

}

//Making request API 입력 테스트 및 구조체 구현
//Json 엔코딩 디코딩 구현

//Mux 추가
//Mux 핸드러 Userinput Handler 추가

//질문간 맥락유지 기능 추가 완료
//큐 구조 슬라이스로 메시지 구성

//대화기록 세이브 기능 추가 완료
//History핸들러 기능 추가 및 파일 쓰기
//Frontend side ? 파일 브라우저 다운로드 구현
//Port 포워딩 및 빌드 후 시운용 성공

//

//HTML 수정 및 디자인 약간 수정 완료

//need to implement

///Server console interface creating requried.  ?? for temporary terminal view created.

//Thread, Session implement
//create conversation session on each ID or Client or IP address
//how could identify each user who access to website ?

//client access > send request > with IP address (network base)
//client access > ID / PW
//clent access > ??

//DB 서버 올리고
//PW 원웨이 함수로 해쉬처리
//ID / PW 저장 회원가입 받음

//implemtented.
//3 people uses possible at a same time? //livemode

//ID PW system implemet // SQL server.

//Image handler추가,.. 사용자 request 기반.

//웹소켓 연결 끊김시, GPT 응답이 한번 안들려오는 문제 개선해야함. 버그찾아야함 ....
//
