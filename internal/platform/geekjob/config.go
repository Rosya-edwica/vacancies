package geekjob

import (
	"net/url"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	"vacancies/tools"

	apiJson "vacancies/api"

	"github.com/gocolly/colly"
)

type GeekJob struct {
	PositionName string
	PositionId   int
	CityId       int
	CityEdwicaId int
	Cities       []models.City
	Currencies   []models.Currency
	Headers      map[string]string
}

func (api *GeekJob) CreateQuery() (query string) {
	params := url.Values{
		"qs": {api.PositionName},
	}
	query = "https://geekjob.ru/json/find/vacancy?" + params.Encode()
	return
}

func (api *GeekJob) CountVacanciesByQuery(url string) (count int) {
	resp, _ := apiJson.DecondeJsonResponse(url, "GET", nil, &apiJson.GeekJobResponseFound{})
	found := resp.(*apiJson.GeekJobResponseFound)
	logger.Log.Printf("Нашлось %d вакансий для профессии '%s'", found.CountVacancies, api.PositionName)
	return found.CountVacancies
}

func getHTMLBody(url string) (body *colly.HTMLElement) {
	c := colly.NewCollector()

	c.OnHTML("body", func(h *colly.HTMLElement) {
		body = h
	})
	err := c.Visit(url)
	tools.CheckErr(err)
	return
}
