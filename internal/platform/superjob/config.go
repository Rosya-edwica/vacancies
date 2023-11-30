package superjob

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"

	apiJson "vacancies/api"
)

type Superjob struct {
	PositionName string
	PositionId   int
	CityId       int
	CityEdwicaId int
	Cities       []models.City
	Currencies   []models.Currency
	Headers      map[string]string
}

const ClientId = "2915"

func (api *Superjob) CreateQuery() (query string) {
	params := url.Values{
		"count":             {"100"},
		"keywords[0][srws]": {"1"},              // Ищем в названии вакансии
		"keywords[0][skwc]": {"particular"},     // Ищем точную фразу
		"keywords[0][keys]": {api.PositionName}, // Фраза
	}
	if api.CityId != 0 {
		params.Add("town", strconv.Itoa(api.CityId))
	}

	return "https://api.superjob.ru/2.0/vacancies?" + params.Encode()
}

func (api *Superjob) UpdateAccessToken() (token string) {
	fmt.Println("Пытаемся обновить токен superjob")
	url := api.createQueryToUpdateToken()
	resp, statusCode := apiJson.DecondeJsonResponse(url, "POST", nil, &apiJson.SuperJobResponseAccessToken{})
	if statusCode != 200 {
		panic("Проблема с токеном Superjob")
	}
	token = resp.(*apiJson.SuperJobResponseAccessToken).Token
	if token != "" {
		fmt.Println("Успешно обновили токен!")
	} else {
		panic("Проблема с токеном Superjob!")
	}
	return
}

func (api *Superjob) CountVacanciesByQuery(url string) (count int) {
	resp, statusCode := apiJson.DecondeJsonResponse(url, "GET", api.Headers, &apiJson.SuperJobResponseFound{})
	if statusCode == 410 {
		newToken := api.UpdateAccessToken()
		api.Headers["Authorization"] = fmt.Sprintf("Bearer %s", newToken)
		return api.CountVacanciesByQuery(url)
	}
	found := resp.(*apiJson.SuperJobResponseFound)
	logger.Log.Printf("[StatusCode: %d] Нашлось %d вакансий для профессии '%s'", statusCode, found.CountVacancies, api.PositionName)
	if found != nil {
		return found.CountVacancies
	}
	return 0
}

func (api *Superjob) createQueryToUpdateToken() string {
	params := url.Values{
		"refresh_token": {strings.ReplaceAll(api.Headers["Authorization"], "Bearer ", "")},
		"client_id":     {ClientId},
		"client_secret": {api.Headers["X-Api-App-Id"]},
	}
	url := "https://api.superjob.ru/2.0/oauth2/refresh_token?" + params.Encode()
	return url
}
