package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CompletionReq struct {
	Model       string       `json:"model"`
	Messages    []ReqMessage `json:"messages"`
	MaxTokens   *int         `json:"max_tokens,omitempty"`
	Temperature *float64     `json:"temperature,omitempty"`
}

type ReqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type Message struct {
	Role    string  `json:"role"`
	Content string  `json:"content"`
	Refusal *string `json:"refusal"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type CompletionResp struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

func (resp CompletionResp) GetText() string {
	return resp.Choices[0].Message.Content
}

func getToken() (token string, err error) {
	token, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		err = errors.New("OPENAI_API_KEY is not setted")
	}
	return
}

func FetchCompletions(requestBody CompletionReq) (response CompletionResp, err error) {
	token, err := getToken()
	if err != nil {
		return CompletionResp{}, fmt.Errorf("error retrieving token: %w", err)
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return CompletionResp{}, fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return CompletionResp{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CompletionResp{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return CompletionResp{}, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		indentedRespBody, _ := json.MarshalIndent(json.RawMessage(respBody), "", "  ")
		return CompletionResp{}, fmt.Errorf("received status code %d, request: %s, response: %s", resp.StatusCode, string(indentedRespBody), string(respBody))
	}

	var chatgptResponse CompletionResp
	if err := json.Unmarshal(respBody, &chatgptResponse); err != nil {
		return CompletionResp{}, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return chatgptResponse, nil
}
