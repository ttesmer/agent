package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"bytes"
)

type Agent struct {
	client *Client
	getUserMessage func() (string, bool)
}

type Client struct {
	Messages *Messages
}

type Messages struct {}

func NewAgent(client *Client, getUserMessage func() (string, bool)) *Agent {
	return &Agent{
		client: client,
		getUserMessage: getUserMessage,
	}
}

func NewClient() *Client {
	return &Client{}
}

func (a *Agent) runInference(ctx context.Context, conversation []Message) (Message, error) {
	message, err := a.client.Messages.New(ctx, MessageBody{
		Model: "minimax/minimax-m2.5",
		Messages: conversation,
		MaxTokens: int64(1024),
	})
	return message, err
}

func (m *Messages) New(ctx context.Context, msg MessageBody) (Message, error){
	API_URL := "https://openrouter.ai/api/v1/chat/completions"
	jsonBytes, _ := json.Marshal(msg)
	body := bytes.NewReader(jsonBytes)
	req, err := http.NewRequest("POST", API_URL, body)
	if err != nil {
		return Message{}, err
	}

	API_KEY := os.Getenv("OPENROUTER_API_KEY")
	if API_KEY == "" {
		return Message{}, fmt.Errorf("No API Key!\n")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()

	var result OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return Message{}, err
	}

	if len(result.Choices) == 0 {
			return Message{}, fmt.Errorf("no choices in response\n")
	}

	return result.Choices[0].Message, err
}

type Choice struct {
    Index        int     `json:"index"`
    FinishReason string  `json:"finish_reason"`
    Message      Message `json:"message"`
}

type OpenRouterResponse struct {
    Choices []Choice `json:"choices"`
}

type MessageBody struct {
	Model 		string    `json:"model"`
	MaxTokens int64     `json:"max_tokens"`
	Messages  []Message `json:"messages"`
	Tools 		[]Tool    `json:"tools,omitempty"`
}

type Tool struct {
	Type 		 string      `json:"type"` // "function"
	Function FunctionDef `json:"function"`
}

type FunctionDef struct {
	Name 				string `json:"name"`
	Description string `json:"description"`
	Parameters  any 	 `json:parameters"` // JSON Schema object
}

type ToolCall struct {
	ID 			 string 		  `json:"id"`
	Type 		 string 		  `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name 			string `json:"name"`
	Arguments string `json:"arguments"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"` // "user" or "assistant"
}

type Content struct {
	Type string // e.g. 'text'
	Text string
}

func NewUserMessage(input string) Message {
	return Message{
		Content: input,
		Role: "user",
	}
}

func (a *Agent) Run(ctx context.Context) error {
	conversation := []Message{}

	fmt.Println("Chat with The Agent (use 'ctrl-c' to quit)")
	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		userInput, ok := a.getUserMessage()
		if !ok {
			break
		}

		userMessage := NewUserMessage(userInput)
		conversation := append(conversation, userMessage)

		message, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}
		conversation = append(conversation, message)

		// this is probs the agentic stuff
		msgContent := []Content{{
			Type: "text",
			Text: message.Content,
		}}
		for _, content := range msgContent {
			switch content.Type {
			case "text":
				fmt.Printf("\u001b[93mThe Agent\u001b[0m: %s\n", content.Text)
			}
		}
	}

	return nil
}

func main() {
	client := NewClient()

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}
	agent := NewAgent(client, getUserMessage)
	err := agent.Run(context.TODO())
	if err != nil {
		fmt.Print("Error: %s\n", err.Error())
	}
}
