package dto

import "time"

type Person struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
	CreatedAt   time.Time
}

type PersonFilter struct {
	Name     string
	Surname  string
	Patronym string
	Gender   string
	Country  string
	Age      string
}
