package trudvsem

import (
	"net/url"
	"vacancies/pkg/models"
)

type TrudVsem struct {
	PositionName string
	PositionId   int
	CityId       int
	CityEdwicaId int
	Cities []models.City
}

func (api *TrudVsem) CreateQuery() (query string) {
	params := url.Values{
		"text": {api.PositionName},
	}
	query = "http://opendata.trudvsem.ru/api/v1/vacancies?" + params.Encode()
	return
}

