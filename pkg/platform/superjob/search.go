package superjob

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"vacancies/pkg/database"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
)

func (api *Superjob) CollectAllVacanciesByQuery(position models.Position, db *database.DB) (vacancies []models.Vacancy) {
	api.PositionId = position.Id
	api.PositionName = position.Name

	if len(api.Cities) == 0 {
		api.Cities = db.GetCities()
	}
	if len(api.Currencies) == 0 {
		api.Currencies = database.GetCurrencies()
	}

	queryForCounting := api.CreateQuery()
	countVacancies := api.CountVacanciesByQuery(queryForCounting)
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

func (api *Superjob) CollectVacanciesFromPage(url string) (vacancies []models.Vacancy) {
	vacanciesResp := api.decodeSuperjobResponse(url)
	for _, item := range vacanciesResp.Items {
		var vacancy models.Vacancy
		var profAreas, specs []string
		if api.CityEdwicaId != 0 {
			vacancy.CityId = api.CityEdwicaId
		} else {
			vacancy.CityId = api.getEdwicaCityID(item.City.Id)
		}

		salary := api.convertSalaryToRUR(models.Salary{
			From:     float64(item.SalaryFrom),
			To:       float64(item.SalaryTo),
			Currency: item.Currency,
		})
		for _, area := range item.ProfAreas {
			profAreas = append(profAreas, area.Name)
			for _, spec := range area.Specializations {
				specs = append(specs, spec.Name)
			}
		}

		vacancy = models.Vacancy{
			Platform:        "superjob",
			ProfessionId:    api.PositionId,
			Id:              string(item.Id),
			Title:           item.Name,
			Url:             item.Url,
			DateUpdate:      convertDateUpdate(item.PublishedAt),
			Experience:      convertExperienceId(item.Experience.Id),
			SalaryFrom:      salary.From,
			SalaryTo:        salary.To,
			ProfAreas:       strings.Join(profAreas, "|"),
			Specializations: strings.Join(specs, "|"),
		}

	}
	return
}

func (api *Superjob) decodeSuperjobResponse(url string) SuperJobResponse {
	var superjobResp SuperJobResponse
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/117.0)")
	req.Header.Add("Authorization", "Bearer v3.r.135952025.d48b84420bb8608422ec51987aea9ed6d5a42598.db29fd855ce56de6a98d41dfb1eae04523fd9ae9")
	req.Header.Add("X-Api-App-Id", "v3.r.135952025.dd781d8411025c15dc44627a14945a740b16fdcb.1dc93ca9bece339431c472b92760366c2f6c52a5")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&superjobResp)
	if err != nil {
		panic(err)
	}
	return superjobResp
}
