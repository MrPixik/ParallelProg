package static

const (
	ResourceURL       = "https://plato.stanford.edu/search/search"
	Query             = "elbow+room"
	PageToParseAmount = 5
	ApiKey            = "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"
	AiURL             = "https://openrouter.ai/api/v1/completions"
	Model             = "deepseek/deepseek-chat-v3-0324:free"
	PromptSample      = "Переведи на русский и сделай краткую выжимку текста:\n"
	MaxTokens         = 300
	FilenameAsync     = "lab3/async/summaries.txt"
	FilenameSync      = "lab3/syncr/summaries.txt"
)
