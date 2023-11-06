package geekjob

import (
	"fmt"
	"sync"
	apiJson "vacancies/api"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
)

func (api *GeekJob) CollectAllVacanciesByQuery(position models.Position) (vacancies []models.Vacancy) {
	api.PositionId = position.Id
	api.PositionName = position.Name

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
	resp, _ := apiJson.DecondeJsonResponse(url, nil, &apiJson.GeekJobResponse{})
	data := resp.(*apiJson.GeekJobResponse)
	for _, vacancy := range data.Items {
		ids = append(ids, vacancy.Id)
	}
	return
}
