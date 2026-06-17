// Package http - response.go - Describe Http Response
package http

type (
	DeckBuildCard struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Title    string `json:"title"`
		Number   string `json:"number"`
		Set      string `json:"set"`
		Type     string `json:"type"`
		Rarity   string `json:"rarity"`
		Quantity int    `json:"quantity"`
		Filepath string `json:"filepath"`
	}

	DeckBuildResponse struct {
		Size     int             `json:"size"`
		Strategy string          `json:"strategy"`
		Name     string          `json:"name"`
		Deck     []DeckBuildCard `json:"deck"`
	}
)
