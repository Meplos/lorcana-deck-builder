package infra

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type AIAgent interface {
	Send(req DeckRequest) (DeckResponse, error)
}

type MistralAIAgent struct {
	Key     string
	AgentID string
	Client  *http.Client
	Model   string

	URL string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Object  string `json:"object"`
	Type    string `json:"type"`
}

type ChatRequest struct {
	Inputs  []Message `json:"inputs"`
	Stream  bool      `json:"stream"`
	AgentID string    `json:"agent_id"`
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Output struct {
	Object      string    `json:"object"`
	Type        string    `json:"message.output"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	AgentID     string    `json:"agent_id"`
	Model       string    `json:"model"`
	ID          string    `json:"id"`
	Role        string    `json:"role"`
	Content     string    `json:"content"`
}

type ChatResponse struct {
	Object         string     `json:"object"`
	ConversationID string     `json:"conversation_id"`
	Outputs        []Output   `json:"outputs"`
	Usage          TokenUsage `json:"usage"`
	Guardrails     any        `json:"guardrails"`
	// "guardrails":null
}

func NewAIAgent() AIAgent {
	agentID := os.Getenv("MISTRAL_AGENT_ID")
	key := os.Getenv("MISTRAL_API_KEY")
	if agentID == "" || key == "" {
		log.Fatal("missing MISTRAL_AGENT_ID or MISTRAL_API_KEY")
	}
	model := os.Getenv("MISTRAL_MODEL")
	if model == "" {
		model = "mistral-medium-latest"
	}
	url := os.Getenv("MISTRAL_URL")
	if url == "" {
		url = "https://api.mistral.ai/v1/conversations"
	}
	return &MistralAIAgent{
		Model:   model,
		AgentID: agentID,
		Client:  &http.Client{},
		Key:     key,
		URL:     url,
	}
}

func (ai *MistralAIAgent) Send(req DeckRequest) (DeckResponse, error) {
	rawContent, err := json.Marshal(req)
	if err != nil {
		return DeckResponse{}, err
	}

	body := ChatRequest{
		Stream:  false,
		AgentID: ai.AgentID,
		Inputs: []Message{
			{
				Role:    "user",
				Content: string(rawContent),
				Object:  "entry",
				Type:    "message.input",
			},
		},
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return DeckResponse{}, err
	}

	request, err := http.NewRequest(http.MethodPost, ai.URL, bytes.NewReader(raw))
	if err != nil {
		return DeckResponse{}, err
	}

	request.Header.Add("Authorization", "Bearer "+ai.Key)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := ai.Client.Do(request)
	if err != nil {
		return DeckResponse{}, err
	}

	log.Printf("Received status: %d", response.StatusCode)
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return DeckResponse{}, err
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Error : %s", responseBody)
		return DeckResponse{}, errors.New("error, occured")
	}

	var rawResponse ChatResponse
	if err := json.Unmarshal(responseBody, &rawResponse); err != nil {
		return DeckResponse{}, err
	}

	message := rawResponse.Outputs[0].Content
	var deck DeckResponse
	if err := json.Unmarshal([]byte(message), &deck); err != nil {
		return DeckResponse{}, nil
	}

	return deck, nil
}
