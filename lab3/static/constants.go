package static

const (
	ResourceURL       = "https://plato.stanford.edu/search/search"
	Query             = "elbow+room"
	PageToParseAmount = 2
	ApiKey            = "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"
	AiURL             = "https://openrouter.ai/api/v1/completions"
	Model             = "deepseek/deepseek-chat-v3-0324:free"
	PromptSample      = "Переведи на русский и сделай краткую выжимку текста:\n"
	MaxTokens         = 1000
	FilenameAsync     = "lab3/async/summaries.txt"
	FilenameSync      = "lab3/syncr/summaries.txt"
)
