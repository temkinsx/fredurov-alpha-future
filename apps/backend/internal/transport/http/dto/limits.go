package dto

type LimitsResponse struct {
	MaxFileSizeBytes int `json:"max_file_size_bytes"`
	MaxFileTextChars int `json:"max_file_text_chars"`
	MaxHistoryChars  int `json:"max_history_chars"`
	MaxPromptChars   int `json:"max_prompt_chars"`
}
