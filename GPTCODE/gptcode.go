package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Request struct {
	Model       string   `json:"model"`
	Messages    Messages `json:"messages"`
	Temperature float32  `json:"temperature"`
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created time.Time `json:"created"`
	Model   string    `json:"model"`
	Usage   Usage     `json:"usage"`
	Choices []Choices `json:"choices"`
}

type Usage struct {
	PromptToken     int   `json:"prompt_token"`
	CompletionToken int   `json:"completion_token"`
	TotalToken      int64 `json:"total_token"`
}

type Choices struct {
	Messages     Messages `json:"message"` // Corrected JSON tag from "messages"
	Logprobs     any      `json:"logprobs"`
	Finishreason string   `json:"finish_reason"` // Corrected JSON tag from "finishreason"
	Index        int      `json:"index"`         // Corrected field name from IDX to Index
}

func (r *Request) InputUserChat() *Request {
	fmt.Println("Hi! Please input a prompt you want:")
	fmt.Scanln(&r.Messages.Content)
	return r
}

func MakingRequest() *http.Request {
	apikey := "YOUR_API_KEY" // Replace with your actual API key
	url := "https://api.openai.com/v1/chat/completions"

	userreq := &Request{
		Model:       "gpt-4o-mini",
		Temperature: 0.7,
		Messages: Messages{
			Role: "user",
		},
	}
	userreq.InputUserChat()

	ur, err := json.Marshal(userreq)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return nil
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(ur))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apikey)
	return req
}

func GetResponse(req *http.Request) *http.Response {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error during HTTP request:", err)
		return nil
	}

	return res
}

func (gres *Response) TransformRes(res *http.Response) *Response {
	if res != nil && res.Body != nil {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return nil
		}

		err = json.Unmarshal(body, gres)
		if err != nil {
			fmt.Println("Error decoding response JSON:", err)
			fmt.Println("Response body:", string(body)) // Print raw response body for debugging
			return nil
		}
	} else {
		fmt.Println("No response received or response body is nil")
	}
	return gres
}

func PrintResponse(res *Response) {
	if res != nil {
		fmt.Printf("Response ID: %s\nModel: %s\nCreated: %s\n", res.ID, res.Model, res.Created)
		fmt.Printf("Choices: %+v\n", res.Choices)
		fmt.Printf("Usage: %+v\n", res.Usage)
	} else {
		fmt.Println("No response to print")
	}
}

func main() {
	req := MakingRequest()
	if req != nil {
		rawres := GetResponse(req)
		if rawres != nil {
			res := new(Response)
			transformedRes := res.TransformRes(rawres)
			PrintResponse(transformedRes)
		} else {
			fmt.Println("Failed to get a valid response from the server.")
		}
	} else {
		fmt.Println("Failed to create a valid request.")
	}
}
