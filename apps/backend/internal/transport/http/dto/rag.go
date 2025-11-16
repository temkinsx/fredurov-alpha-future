package dto

type RAGSearchRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k"`
}

type RAGChunk struct {
	Text       string  `json:"text"`
	DocumentID string  `json:"document_id"`
	SourceName string  `json:"source_name"`
	Score      float64 `json:"score"`
}

type RAGSearchResponse struct {
	Chunks []RAGChunk `json:"chunks"`
}
