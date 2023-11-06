package trudvsem

import (
	"fmt"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
)

func (api *TrudVsem) CollectAllVacanciesByQuery(position models.Position) (vacancies []models.Vacancy) {
	api.PositionId = position.Id
	api.PositionName = position.Name

	countVacancies := api.CountVacanciesByQuery(api.CreateQuery())
	if countVacancies == 0 {
		return
	}
	return api.FindVacanciesInRussia()
}

func (api *TrudVsem) FindVacanciesInRussia() (vacancies []models.Vacancy) {
	var pageNum = 0
	for {
		url := fmt.Sprintf("%s&offset=%d", api.CreateQuery(), pageNum)
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

func (api *TrudVsem) FindVacanciesInCurrentCity(city models.City) (vacancies []models.Vacancy) {
	return
}
