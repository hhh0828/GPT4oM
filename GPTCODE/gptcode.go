package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// 요청 Structure 구현
type GRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// Scan으로 Input 요청 구현
func (r *GRequest) ScanUserRequest() *GRequest {
	fmt.Println("Enter your prompt:")
	fmt.Scanln(&r.Prompt)

	fmt.Println("Enter max tokens (integer):")
	fmt.Scanln(&r.MaxTokens)

	fmt.Println("Enter temperature (float):")
	fmt.Scanln(&r.Temperature)

	return r
}

// 요청 Structure Unmarshal Go object >> Json
func (r GRequest) CreateAPIrequest(apiKey string) *http.Request {
	reqBody, _ := json.Marshal(r)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/engines/davinci/completions", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return req
}

// 응답 Structure 구현
type Response struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

// 응답 Structure Marshal 후 Go object로 전환
func (res *Response) Decoding(re *http.Response) {
	_ = json.NewDecoder(re.Body).Decode(res)
}

func main() {
	// API 키 설정 (보안 문제로 환경 변수에서 가져오는 것이 좋습니다)
	apiKey := "sk-None-izmLhx0PUGxalUxl4RaRT3BlbkFJmSVSdLfv7ypsP6U036hH"
	if apiKey == "" {
		fmt.Println("API key is not set")
		return
	}

	// 빈 객체 생성
	requestData := GRequest{}
	// req 사용자 입력 받고 API 리퀘스트 생성
	req := requestData.ScanUserRequest().CreateAPIrequest(apiKey)

	// Server로 부터 응답을 받는 client 객체 생성
	client := &http.Client{}
	// request를 요청해서 client.Do(request)로 response 반환
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	// 응답 데이터 구조체 생성 및 디코딩
	responseData := new(Response)
	responseData.Decoding(resp)

	// 결과 출력
	if len(responseData.Choices) > 0 {
		fmt.Println("Response:", responseData.Choices[0].Text)
	} else {
		fmt.Println("No response received.")
	}
}
