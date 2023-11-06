package superjob

import (
	"fmt"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
)

func (api *Superjob) CollectAllVacanciesByQuery(position models.Position) (vacancies []models.Vacancy) {
	api.PositionId = position.Id
	api.PositionName = position.Name

	queryForCounting := api.CreateQuery()
	countVacancies := api.CountVacanciesByQuery(queryForCounting)
	if countVacancies == 0 {
		return
	}

	if countVacancies < 2000 {
		vacancies = api.FindVacanciesInRussia()
	} else {
		for _, city := range api.Cities {
			if city.SUPERJOB_ID == 0 {
				continue
			}
			logger.Log.Printf("Ищем вакансии в городе:%s", city.Name)
			cityVacancies := api.FindVacanciesInCurrentCity(city)
			vacancies = append(vacancies, cityVacancies...)
		}
	}
	return
}

func (api *Superjob) FindVacanciesInRussia() (vacancies []models.Vacancy) {
	logger.Log.Println("Ищем вакансии по всей России")
	return api.FindVacanciesInCurrentCity(models.City{})
}

func (api *Superjob) FindVacanciesInCurrentCity(city models.City) (vacancies []models.Vacancy) {
	api.CityEdwicaId = city.EDWICA_ID
	api.CityId = city.SUPERJOB_ID

	var pageNum = 0
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
