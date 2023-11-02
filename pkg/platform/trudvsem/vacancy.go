package trudvsem

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
	"vacancies/pkg/models"
)

// func (api *TrudVsem) CollectVacanciesFromPage(url string) (vacancies []models.Vacancy) {
// 	json, err := tools.GetJson(url, "trudvsem")
// 	tools.CheckErr(err)
// 	for _, item := range gjson.Get(json, "results.vacancies").Array() {
// 		var vacancy models.Vacancy
// 		vacancy.Id = item.Get("vacancy.id").String()
// 		vacancy.Platform = "trudvsem"
// 		vacancy.Title = item.Get("vacancy.job-name").String()
// 		vacancy.ProfessionId = api.PositionId
// 		vacancy.SalaryFrom = item.Get("vacancy.salary_min").Float()
// 		vacancy.SalaryTo = item.Get("vacancy.salary_max").Float()
// 		vacancy.Url = item.Get("vacancy.vac_url").String()
// 		vacancy.ProfAreas = item.Get("vacancy.category.specialisation").String()
// 		vacancy.DateUpdate = item.Get("vacancy.creation-date").String()
// 		city := item.Get("vacancy.addresses.address.0.location").String()
// 		vacancy.CityId = api.parseCity(city)
// 		vacancies = append(vacancies, vacancy)
// 	}
// 	return
// }

func (api *TrudVsem) CollectVacanciesFromPage(url string) (vacancies []models.Vacancy) {
	trudvsemResp := api.decodeTrudvsemResponse(url)
	for _, item := range trudvsemResp.Results.Vacancies {
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

func (api *TrudVsem) decodeTrudvsemResponse(url string) TrudvsemResponse {
	var trudvsemResp TrudvsemResponse

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&trudvsemResp)
	if err != nil {
		panic(err)
	}
	return trudvsemResp
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
