package domain

type Limits struct {
	MaxPromptChars    int
	MaxOutputTokens   int
	MaxFileSizeBytes  int
	MaxFileTextChars  int
	MaxHistoryChars   int
	MaxRequestChars   int
	MaxRequestsPerMin int
	MaxConcurrentLLM  int
}
