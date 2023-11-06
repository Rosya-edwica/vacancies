package trudvsem

import (
	"net/url"
	apiJson "vacancies/api"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
)

type TrudVsem struct {
	PositionName string
	PositionId   int
	CityId       int
	CityEdwicaId int
	Cities       []models.City
	Headers      map[string]string
}

func (api *TrudVsem) CreateQuery() (query string) {
	params := url.Values{
		"text": {api.PositionName},
	}
	query = "http://opendata.trudvsem.ru/api/v1/vacancies?" + params.Encode()
	return
}

func (api *TrudVsem) CountVacanciesByQuery(url string) (count int) {
	resp, _ := apiJson.DecondeJsonResponse(url, nil, &apiJson.TrudvsemResponseFound{})
	found := resp.(*apiJson.TrudvsemResponseFound)
	logger.Log.Printf("Нашлось %d вакансий для профессии '%s'", found.Meta.VacanciesCount, api.PositionName)
	if found != nil {
		return found.Meta.VacanciesCount
	}
	return 0
}
