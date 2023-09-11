package geekjob

import (
	"fmt"
	"sync"
	"vacancies/pkg/database"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	"vacancies/pkg/tools"

	"github.com/tidwall/gjson"
)

func (api *GeekJob) CollectAllVacanciesByQuery(position models.Position, db *database.DB) (vacancies []models.Vacancy) {
	api.PositionId = position.Id
	api.PositionName = position.Name

	if len(api.Cities) == 0 {
		api.Cities = db.GetCities()
	}
	if len(api.Currencies) == 0 {
		api.Currencies = database.GetCurrencies()
	}
	return api.FindVacanciesInRussia()
}

func (api *GeekJob) FindVacanciesInRussia() (vacancies []models.Vacancy) {
	logger.Log.Println("Ищем вакансии по всей России")
	return api.FindVacanciesInCurrentCity(models.City{})
}
func (api *GeekJob) FindVacanciesInCurrentCity(city models.City) (vacancies []models.Vacancy) {
	pageNum := 1
	for {
		url := fmt.Sprintf("%s&page=%d", api.CreateQuery(), pageNum)
		pageVacancies := api.CollectVacanciesFromPage(url)
		if len(pageVacancies) == 0 {
			break
		}
		pageNum++
		logger.Log.Printf("Количество вакансий - %d на %d странице", len(pageVacancies), pageNum)
		vacancies = append(vacancies, pageVacancies...)
	}
	return
}
func (api *GeekJob) CollectVacanciesFromPage(url string) (vacancies []models.Vacancy) {
	vacanciesId := api.collectVacanciesIdFromPage(url)
	if len(vacanciesId) == 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(vacanciesId))
	for _, id := range vacanciesId {
		go api.PutVacancyToArrayById(id, &wg, &vacancies)
	}
	wg.Wait()
	return
}

func (api *GeekJob) collectVacanciesIdFromPage(url string) (ids []string) {
	json, err := tools.GetJson(url, "geekjob")
	tools.CheckErr(err)

	for _, item := range gjson.Get(json, "data").Array() {
		id := item.Get("id").String()
		ids = append(ids, id)
	}
	return
}
