package http

type (
	SaveCardRequest struct {
		ID       string `json:"id"`
		Quantity int    `json:"quantity"`
	}

	SaveDeckRequest struct {
		Size     int               `json:"size"`
		Strategy string            `json:"strategy"`
		Name     string            `json:"name"`
		Deck     []SaveCardRequest `json:"deck"`
	}
)
