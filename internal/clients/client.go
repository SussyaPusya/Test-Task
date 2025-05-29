package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"test_task/internal/config"
	"test_task/internal/dto"
)

type ExternalAPI struct {
	urls       *config.ExternalAPI
	httpClient *http.Client
}
type Target struct {
	CountryId string `json:"country_id"`
}

type ReqStruct struct {
	Age         int      `json:"age"`
	Gender      string   `json:"gender"`
	Nationality []Target `json:"country"`
}

func NewExtanlAPI(cfg *config.ExternalAPI) *ExternalAPI {

	return &ExternalAPI{urls: cfg, httpClient: &http.Client{}}
}

func (e *ExternalAPI) GetAge(ctx context.Context, person *dto.Person) (*dto.Person, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.urls.Age+fmt.Sprintf("/?name=%s", person.Name), nil)
	if err != nil {
		return nil, err
	}

	reqStr := ReqStruct{}
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return person, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&reqStr); err != nil {
		return person, err
	}
	fmt.Println(reqStr.Age)
	person.Age = reqStr.Age

	return person, nil

}

func (e *ExternalAPI) GetGender(ctx context.Context, person *dto.Person) (*dto.Person, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.urls.Gender+fmt.Sprintf("/?name=%s", person.Name), nil)
	if err != nil {
		return nil, err
	}

	reqStr := ReqStruct{}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return person, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&reqStr); err != nil {
		return person, err
	}

	person.Gender = reqStr.Gender

	return person, nil
}

func (e *ExternalAPI) GetNationaliti(ctx context.Context, person *dto.Person) (*dto.Person, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.urls.Nationality+fmt.Sprintf("/?name=%s", person.Name), nil)
	if err != nil {
		return nil, err
	}

	reqStr := ReqStruct{}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return person, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&reqStr); err != nil {
		return person, err
	}

	nati := reqStr.Nationality[0].CountryId

	person.Nationality = nati

	return person, nil
}

func (e *ExternalAPI) GetAll(ctx context.Context, person *dto.Person) (*dto.Person, error) {
	person, err := e.GetAge(ctx, person)
	if err != nil {
		return person, err
	}
	person, err = e.GetGender(ctx, person)
	if err != nil {
		return person, err
	}
	person, err = e.GetNationaliti(ctx, person)
	if err != nil {
		return person, err
	}

	return person, err

}
