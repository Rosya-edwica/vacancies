package api

import (
	"encoding/json"

	"net/http"
	"sync"
	"time"
	"vacancies/pkg/models"
)

type API interface {
	// Метод создания ссылки, в которой будут лежать данные
	CreateQuery() (query string)

	// Подсчет вакансий одной профессии по всей России для того, чтобы определить какой метод нам использовать FindVacanciesInRussia() или FindVacanciesInCurrentCity()
	CountVacanciesByQuery(url string) (count int)

	// Сбор всех вакансий с одного запроса по API (все вакансии профессии в городе)
	CollectAllVacanciesByQuery(position models.Position) (vacancies []models.Vacancy)

	// Поиск вакансий по всей России без привязки к городу
	FindVacanciesInRussia() (vacancies []models.Vacancy)

	// Поиск вакансий по конкретному городу для популярных профессий, которых больше 2000 на платформе
	FindVacanciesInCurrentCity(city models.City) (vacancies []models.Vacancy)

	// Сбор всех вакансий с одной страницы запроса
	CollectVacanciesFromPage(url string) (vacancies []models.Vacancy)

	// Сбор одной конкретной вакансии
	PutVacancyToArrayById(id string, wg *sync.WaitGroup, vacancies *[]models.Vacancy)
}

type AnyJsonResponseStruct interface{}

// Парсим Json в структуру, которую передают в dataStruct
func DecondeJsonResponse(url string, headers map[string]string, dataStruct interface{}) (data interface{}, statusCode int) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(dataStruct)
	if err != nil {
		panic(err)
	}
	return dataStruct, resp.StatusCode
}
