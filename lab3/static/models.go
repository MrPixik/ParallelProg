package static

type Article struct {
	Link    string
	Content string `json:"-"`
	Summary string
	Err     error `json:"-"`
}

type ReqOpenRouter struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}
