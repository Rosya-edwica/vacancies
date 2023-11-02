package headhunter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	// "sync"
	"time"
	"vacancies/pkg/database"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	// "vacancies/pkg/tools"
	// "github.com/tidwall/gjson"
)

func (api *HeadHunter) CollectAllVacanciesByQuery(position models.Position, db *database.DB) (vacancies []models.Vacancy) {
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
			if city.HH_ID == 0 {
				continue
			}
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
		if len(pageVacancies) == 0 {
			break
		}
		pageNum++
		logger.Log.Printf("Количество вакансий - %d на %d странице", len(pageVacancies), pageNum)
		vacancies = append(vacancies, pageVacancies...)
	}
	return
}

// func (api *HeadHunter) CollectVacanciesFromPage(url string) (vacancies []models.Vacancy) {
// 	json, err := tools.GetJson(url, "headhunter")
// 	if err != nil {
// 		logger.Log.Printf("Не удалось подключиться к странице %s.\nТекст ошибки: %s", err, url)
// 		return
// 	}
// 	var wg sync.WaitGroup
// 	items := gjson.Get(json, "items").Array()
// 	wg.Add(len(items))
// 	for _, item := range items {
// 		vacancyId := item.Get("id").String()
// 		go api.PutVacancyToArrayById(vacancyId, &wg, &vacancies)
// 	}
// 	wg.Wait()
// 	return
// }

func (api *HeadHunter) CollectVacanciesFromPage(url string) (vacancies []models.Vacancy) {
	hhResp := api.decodeHeadHunterResponse(url)
	for _, item := range hhResp.Items {
		vacancyUrl := fmt.Sprintf("http://api.hh.ru/vacancies/%s", item.Id)
		vacancyResp := api.decodeHeadhunterVacancyResponse(vacancyUrl)
		vacancy := api.steamUpHeadHunterVacancyResponse(vacancyResp)
		vacancies = append(vacancies, vacancy)
	}
	return
}

func (api *HeadHunter) decodeHeadHunterResponse(url string) HeadHunterResponse {
	var hhResp HeadHunterResponse

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&hhResp)
	if err != nil {
		panic(err)
	}
	return hhResp
}

func (api *HeadHunter) decodeHeadhunterVacancyResponse(url string) HeadHunterVacancyResponse {
	var vacancyResp HeadHunterVacancyResponse

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPad; CPU OS 7_2_1 like Mac OS X; en-US) AppleWebKit/533.14.6 (KHTML, like Gecko) Version/3.0.5 Mobile/8B116 Safari/6533.14.6")
	req.Header.Set("Authorization", "Bearer QQAVSIBVU4B0JCR296THKB22JP05A92H329U49TDD9CRIS8DT9BRPPT7M9OLQ6HD")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&vacancyResp)
	if err != nil {
		panic(err)
	}
	// fmt.Println(resp.StatusCode, url)
	return vacancyResp
}

func (api *HeadHunter) steamUpHeadHunterVacancyResponse(vacancyResp HeadHunterVacancyResponse) models.Vacancy {

	cityId, _ := strconv.Atoi(vacancyResp.City.Id)
	var skills, languages []string
	for _, skill := range vacancyResp.Skills {
		skills = append(skills, skill.Name)
	}
	for _, item := range vacancyResp.Languages {
		lang := fmt.Sprintf("%s (%s)", item.Name, item.Level.Name)
		languages = append(languages, lang)
	}
	skills = append(skills, languages...)

	salary := api.convertSalary(models.Salary{
		From:     float64(vacancyResp.Salary.From),
		To:       float64(vacancyResp.Salary.To),
		Currency: vacancyResp.Salary.Currency,
	})
	return models.Vacancy{
		Platform:     "hh",
		Id:           vacancyResp.Id,
		Title:        vacancyResp.Name,
		Url:          vacancyResp.Url,
		DateUpdate:   vacancyResp.CreatedAt,
		Experience:   vacancyResp.Experience.Name,
		CityId:       api.getEdwicaCityID(cityId),
		ProfessionId: api.PositionId,
		Skills:       strings.Join(skills, "|"),
		SalaryFrom:   salary.From,
		SalaryTo:     salary.To,
	}
}
