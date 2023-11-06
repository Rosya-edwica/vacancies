package trudvsem

import (
	"regexp"
	"strings"
	"sync"
	apiJson "vacancies/api"
	"vacancies/pkg/models"
)

func (api *TrudVsem) CollectVacanciesFromPage(url string) (vacancies []models.Vacancy) {
	resp, _ := apiJson.DecondeJsonResponse(url, api.Headers, &apiJson.TrudvsemResponse{})
	trudResp := resp.(*apiJson.TrudvsemResponse)
	for _, item := range trudResp.Results.Vacancies {
		var vacancy models.Vacancy

		vacancy.Platform = "trudvsem"
		vacancy.ProfessionId = api.PositionId
		vacancy.Id = item.Vacancy.Id
		vacancy.Title = item.Vacancy.Name
		vacancy.Url = item.Vacancy.Url
		vacancy.DateUpdate = item.Vacancy.DateUpdate
		vacancy.SalaryFrom = float64(item.Vacancy.SalaryFrom)
		vacancy.SalaryTo = float64(item.Vacancy.SalaryTo)
		vacancy.Specializations = item.Vacancy.Specialisation.Name
		vacancy.CityId = api.parseCity(item.Vacancy.Addressses.Address[0].Location)

		vacancies = append(vacancies, vacancy)
	}
	return
}

func (api *TrudVsem) PutVacancyToArrayById(id string, wg *sync.WaitGroup, vacancies *[]models.Vacancy) {
	return
}

func (api *TrudVsem) parseCity(cityName string) (cityId int) {
	re := regexp.MustCompile(`г .*?,|г .*? `)
	reSub := regexp.MustCompile(`г |г\.|,`)

	city := re.FindString(cityName + " ")
	if len(city) == 0 {
		return
	}

	city = reSub.ReplaceAllString(city, "")
	city = strings.TrimSpace(city)
	for _, item := range api.Cities {
		if strings.ToLower(item.Name) == strings.ToLower(city) {
			return item.EDWICA_ID
		}
	}
	return
}
