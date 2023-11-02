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
	Cities       []models.City
}

type TrudvsemResponse struct {
	Results struct {
		Vacancies []struct {
			Vacancy struct {
				Id             string `json:"id"`
				Name           string `json:"job-name"`
				SalaryFrom     int    `json:"salary_min"`
				SalaryTo       int    `json:"salary_max"`
				Url            string `json:"vac_url"`
				DateUpdate     string `json:"creation-date"`
				Specialisation struct {
					Name string `json:"specialisation"`
				} `json:"category"`
				Addressses struct {
					Address []struct {
						Location string `json:"location"`
					} `json:"address"`
				} `json:"addresses"`
			} `json:"vacancy"`
		} `json:"vacancies"`
	} `json:"results"`
}

func (api *TrudVsem) CreateQuery() (query string) {
	params := url.Values{
		"text": {api.PositionName},
	}
	query = "http://opendata.trudvsem.ru/api/v1/vacancies?" + params.Encode()
	return
}
