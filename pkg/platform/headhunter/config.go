package headhunter

import (
	"fmt"
	"net/url"
	"strconv"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	"vacancies/pkg/tools"

	"github.com/tidwall/gjson"
)

type HeadHunter struct {
	PositionName string
	PositionId   int
	CityId       int
	CityEdwicaId int
	Cities 	[]models.City
	Currencies []models.Currency
}

const (
	TOKEN = "QQAVSIBVU4B0JCR296THKB22JP05A92H329U49TDD9CRIS8DT9BRPPT7M9OLQ6HD"
	dictionaryUrl = "https://api.hh.ru/dictionaries"
	per_page = "60"
	search_field = "name"
)


func (api *HeadHunter) CreateQuery() (query string) {
	var params url.Values
	if api.CityId == 0 {
		params = url.Values{
			"search_field": {search_field},
			"per_page":     {per_page},
			"text":         {api.PositionName},
		}
	} else {
		params = url.Values{
			"search_field": {search_field},
			"per_page":     {per_page},
			"text":         {api.PositionName},
			"area":         {strconv.Itoa(api.CityId)},
		}
	}
	query = "https://api.hh.ru/vacancies?" + params.Encode()
	return
}

func (api *HeadHunter) CountVacanciesByQuery(url string) (count int) {
	json, err := tools.GetJson(url,"headhunter")
	tools.CheckErr(err)
	count = int(gjson.Get(json, "found").Int())
	logger.Log.Printf("Нашлось %d вакансий для профессии '%s'", count, api.PositionName)
	return
}


func (api *HeadHunter) convertSalaryToRUR(salary models.Salary) models.Salary {
	for _, cur := range api.Currencies {
		if cur.Code == salary.Currency {
			salary.To = salary.To / cur.Rate
			salary.From = salary.From / cur.Rate
			salary.Currency = "RUR"
			return salary
		}
	}
	return salary
}

func getCurrencies() (currencies []models.Currency) {
	json, err := tools.GetJson(dictionaryUrl, "headhunter")
	if err != nil {
		fmt.Printf("Не удалось обновить валюту. Текст сообщения: %s", err)
	}
	for _, item := range gjson.Get(json, "currency").Array() {
		currencies = append(currencies, models.Currency{
			Code: item.Get("code").String(),
			Abbr: item.Get("abbr").String(),
			Name: item.Get("name").String(),
			Rate: item.Get("rate").Float(),
		})
	}
	return
}