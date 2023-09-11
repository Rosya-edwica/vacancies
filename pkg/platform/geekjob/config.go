package geekjob

import (
	"net/url"
	"vacancies/pkg/models"
	"vacancies/pkg/tools"

	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

type GeekJob struct {
	PositionName string
	PositionId   int
	CityId       int
	CityEdwicaId int
	Cities       []models.City
	Currencies   []models.Currency
}

func (api *GeekJob) CreateQuery() (query string) {
	params := url.Values{
		"qs": {api.PositionName},
	}
	query = "https://geekjob.ru/json/find/vacancy?" + params.Encode()
	return
}

func (api *GeekJob) CountVacanciesByQuery(url string) (count int) {
	json, err := tools.GetJson(url, "geekjob")
	tools.CheckErr(err)
	count = int(gjson.Get(json, "documentsCount").Int())
	return
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
