package headhunter

import (
	"fmt"
	"sync"
	"vacancies/pkg/database"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	"vacancies/pkg/tools"

	"github.com/tidwall/gjson"
)

func (api *HeadHunter) CollectAllVacanciesByQuery(position models.Position, db *database.DB) (vacancies []models.Vacancy) {
	api.PositionId = position.Id
	api.PositionName = position.Name

	if len(api.Cities) == 0 {
		api.Cities = db.GetCities()
	}
	if len(api.Currencies) == 0 {
		api.Currencies = getCurrencies()
	}

	queryForCounting := api.CreateQuery()
	countVacancies := api.CountVacanciesByQuery(queryForCounting)

	if countVacancies < 2000 {
		vacancies = api.FindVacanciesInRussia()
	} else {
		for _, city := range api.Cities {
			if city.HH_ID == 0 { continue }
			logger.Log.Printf("Ищем вакансии в городе:%s", city.Name)
			cityVacancies := api.FindVacanciesInCurrentCity(city)
			vacancies = append(vacancies, cityVacancies...)
		}
	}
	return
}

func (api *HeadHunter) FindVacanciesInRussia() (vacancies []models.Vacancy) {
	logger.Log.Println("Ищем вакансии по всей России")
	return api.FindVacanciesInCurrentCity(models.City{})
}

func (api *HeadHunter) FindVacanciesInCurrentCity(city models.City) (vacancies []models.Vacancy) {
	api.CityEdwicaId = city.EDWICA_ID
	api.CityId = city.HH_ID

	var pageNum = 0
	for {
		url := fmt.Sprintf("%s&page=%d", api.CreateQuery(), pageNum)
		pageVacancies := api.CollectVacanciesFromPage(url)
		if len(pageVacancies) == 0 { break }
		pageNum++
		logger.Log.Printf("Количество вакансий - %d на %d странице", len(pageVacancies), pageNum)
		vacancies = append(vacancies, pageVacancies...)
	}
	return
}

func (api *HeadHunter) CollectVacanciesFromPage(url string) (vacancies []models.Vacancy) {
	json, err := tools.GetJson(url, "headhunter")
	if err != nil {
		logger.Log.Printf("Не удалось подключиться к странице %s.\nТекст ошибки: %s", err, url)
		return
	}
	var wg sync.WaitGroup
	items := gjson.Get(json, "items").Array()
	wg.Add(len(items))
	for _, item := range items {
		vacancyId := item.Get("id").String()
		go api.PutVacancyToArrayById(vacancyId, &wg, &vacancies)
	}
	wg.Wait()
	fmt.Println(len(items), "=", len(vacancies))
	return
}