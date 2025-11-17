package dto

type Scenario struct {
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ScenariosResponse struct {
	Scenarios []Scenario `json:"scenarios"`
}
