package dto

import "time"

// Person структура человека
// @Description Структура, содержащая информацию о человеке.
type Person struct {
	ID          string `json:"id" example:"1"`
	Name        string `json:"name" example:"Dmitriy"`
	Surname     string `json:"surname"  example:"Ushakov"`
	Patronymic  string `json:"patronymic,omitempty" example:"Vasilevich"`
	Gender      string `json:"gender" example:"male"`
	Age         int    `json:"age" example:"32"`
	Nationality string `json:"nationality" example:"RU"`
	CreatedAt   time.Time
}

// PersonFilter структура для фильтрации
// @Description Параметры фильтрации людей при получении списка.
type PersonFilter struct {
	Name     string `json:"name,omitempty" query:"name"`
	Surname  string `json:"surname,omitempty" query:"surname"`
	Patronym string `json:"patronymic,omitempty" query:"patronymic"`
	Gender   string `json:"gender,omitempty" query:"gender"`
	Age      string `json:"age,omitempty" query:"age"`
	Country  string `json:"country,omitempty" query:"country"`
}
