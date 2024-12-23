package models

type Transaction struct {
	ID          string `json:"id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	Date        string `json:"date"`
	UserID      string `json:"userId"`
	Category    string `json:"category"`
}
