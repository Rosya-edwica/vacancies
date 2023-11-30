package headhunter

import (
	"net/url"
	"strconv"
	apiJson "vacancies/api"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
)

type HeadHunter struct {
	PositionName string
	PositionId   int
	CityId       int
	CityEdwicaId int
	Cities       []models.City
	Currencies   []models.Currency
	Headers      map[string]string
}

/*
2000 - это максимальное количество вакансий для одного запроса.
Например, если вакансий python-разработчика всего 5 тысяч, то hh покажет нам максимум 2000 вакансий
В таком случае, чтобы собрать как можно больше вакансий, мы будем искать python-разработчиков по отдельным городам
Если всего вакансий python-разработчика меньше 5 тысяч, то нам не нужно делать лишние запросы по отдельным городам
и можно спокойно спарсить все вакансии без привязки к городу (то есть по всей России)
*/
const HeadhunterVacanciesLimitInOneCity = 2000

// Наполняем API-запрос к HH необходимыми параметрами. Так как это метод структуры HeadHunter, нам не нужно ничего передавать в функцию, т.к. необходимые данные лежат в структуре
func (api *HeadHunter) CreateQuery() (query string) {
	var params url.Values
	const (
		perPage      = "60"
		searchField  = "name"
		domain       = "https://api.hh.ru/vacancies?"
		idRussiaArea = "113"
	)
	// Если город не указан, значит ищем по всей России
	if api.CityId == 0 {
		params = url.Values{
			"search_field": {searchField},
			"per_page":     {perPage},
			"text":         {api.PositionName},
			"area":         {idRussiaArea},
		}
		return domain + params.Encode()
	}
	params = url.Values{
		"search_field": {searchField},
		"per_page":     {perPage},
		"text":         {api.PositionName},
		"area":         {strconv.Itoa(api.CityId)},
	}
	return domain + params.Encode()
}

// Подсчитываем количество вакансий, нам это нужно для того, чтобы понять искать по всей России или по отдельным городам
func (api *HeadHunter) CountVacanciesByQuery(url string) (count int) {
	resp, _ := apiJson.DecondeJsonResponse(url, "GET", api.Headers, &apiJson.HeadHunterResponseFound{})
	found := resp.(*apiJson.HeadHunterResponseFound)
	logger.Log.Printf("Нашлось %d вакансий для профессии '%s'", found.CountVacancies, api.PositionName)
	if found != nil {
		return found.CountVacancies
	}
	return 0
}

// Мы должны приводить всю зарплату к рублям
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
