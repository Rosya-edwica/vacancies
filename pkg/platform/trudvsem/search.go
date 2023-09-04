package trudvsem

import (
	"fmt"
	"vacancies/pkg/database"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	"vacancies/pkg/tools"

	"github.com/tidwall/gjson"
)


func (api *TrudVsem) CollectAllVacanciesByQuery(position models.Position, db *database.DB) (vacancies []models.Vacancy) {
	api.PositionId = position.Id
	api.PositionName = position.Name
	if len(api.Cities) == 0 {
		api.Cities = db.GetCities()
	}
	return api.FindVacanciesInRussia()
}

func (api *TrudVsem) CountVacanciesByQuery(url string) (count int) {
	json, err := tools.GetJson(url, "trudvsem")
	tools.CheckErr(err)
	count = int(gjson.Get(json, "meta.total").Int())
	logger.Log.Printf("Нашлось %d вакансий для профессии '%s'", count, api.PositionName)
	return
}

func (api *TrudVsem) FindVacanciesInRussia() (vacancies []models.Vacancy) {
	var pageNum = 0
	for {
		url := fmt.Sprintf("%s&offset=%d", api.CreateQuery(), pageNum)
		pageVacancies := api.CollectVacanciesFromPage(url)
		if len(pageVacancies) == 0 { break }
		pageNum++
		logger.Log.Printf("Количество вакансий - %d на %d странице", len(pageVacancies), pageNum)
		vacancies = append(vacancies, pageVacancies...)
	}
	return
}

func (api *TrudVsem) FindVacanciesInCurrentCity(city models.City) (vacancies []models.Vacancy) {
	// return api.FindVacanciesInRussia(db)
	return
}
