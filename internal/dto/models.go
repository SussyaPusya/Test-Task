package dto

import "time"

type Person struct {
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Gender      string  `json:"gender"`
	Age         int     `json:"age"`
	Nationality string  `json:"nationality"`
	CreatedAt   time.Time
}
