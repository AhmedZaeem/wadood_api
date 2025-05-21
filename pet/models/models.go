package models

type Pet struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Type    string `json:"type"`
	OwnerID int    `json:"owner_id"`
	Gender  string `json:"gender"`
}
