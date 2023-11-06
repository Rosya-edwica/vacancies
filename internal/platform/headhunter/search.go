package headhunter

import (
	"fmt"
	"sync"
	apiJson "vacancies/api"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
)

// Собираем все вакансии одной профессии. Этот метод определяет искать вакансии с привязкой к городу или нет
func (api *HeadHunter) CollectAllVacanciesByQuery(position models.Position) (vacancies []models.Vacancy) {
	api.PositionId = position.Id
	api.PositionName = position.Name

	countVacancies := api.CountVacanciesByQuery(api.CreateQuery())
	if countVacancies == 0 {
		return
	}

	if countVacancies < HeadhunterVacanciesLimitInOneCity {
		vacancies = api.FindVacanciesInRussia()
		return
	}
	for _, city := range api.Cities {
		if city.HH_ID == 0 {
			continue
		}
		cityVacancies := api.FindVacanciesInCurrentCity(city)
		vacancies = append(vacancies, cityVacancies...)
	}
	return
}

// Собираем вакансии без привязки к городу с помощью пустой структуры города
func (api *HeadHunter) FindVacanciesInRussia() (vacancies []models.Vacancy) {
	logger.Log.Println("Поиск без привязки к городу")
	return api.FindVacanciesInCurrentCity(models.City{})
}

// Собираем вакансии в конкретном городе
func (api *HeadHunter) FindVacanciesInCurrentCity(city models.City) (vacancies []models.Vacancy) {
	if city.Name != "" {
		logger.Log.Printf("Поиск в городе:%s", city.Name)
	}
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
		logger.Log.Printf("Страница № %d: количество вакансий: %d ", pageNum, len(pageVacancies))
		vacancies = append(vacancies, pageVacancies...)
	}
	return
}

// Переходим по всем ссылкам на вакансии и парсим их из списка вакансий
func (api *HeadHunter) CollectVacanciesFromPage(url string) (vacancies []models.Vacancy) {
	// Декодируем json-ответ headhunter в виде списка вакансий
	respData, _ := apiJson.DecondeJsonResponse(url, api.Headers, &apiJson.HeadHunterResponse{})
	resp := respData.(*apiJson.HeadHunterResponse)

	// Формируем список из id вакансий. Чтобы передать их горутинам
	var listVacaciesId []string
	for _, item := range resp.Items {
		listVacaciesId = append(listVacaciesId, item.Id)
	}

	// Создаем горутины
	var wg sync.WaitGroup
	wg.Add(len(listVacaciesId))
	for _, item := range listVacaciesId {
		go api.PutVacancyToArrayById(item, &wg, &vacancies)
	}
	wg.Wait()

	return
}
