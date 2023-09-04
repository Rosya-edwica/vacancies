package api

import (
	"vacancies/pkg/database"
	"vacancies/pkg/models"
	"sync"
)

type API interface {
	// Метод создания ссылки, в которой будут лежать данные
	CreateQuery() (query string)
	
	// Подсчет вакансий одной профессии по всей России для того, чтобы определить какой метод нам использовать FindVacanciesInRussia() или FindVacanciesInCurrentCity()
	CountVacanciesByQuery(url string) (count int)

	// Сбор всех вакансий с одного запроса по API (все вакансии профессии в городе)
	CollectAllVacanciesByQuery(position models.Position, db *database.DB) (vacancies []models.Vacancy)
	
	// Поиск вакансий по всей России без привязки к городу
	FindVacanciesInRussia() (vacancies []models.Vacancy)

	// Поиск вакансий по конкретному городу для популярных профессий, которых больше 2000 на платформе
	FindVacanciesInCurrentCity(city models.City) (vacancies []models.Vacancy)

	// Сбор всех вакансий с одной страницы запроса
	CollectVacanciesFromPage(url string) (vacancies []models.Vacancy)

	// Сбор одной конкретной вакансии
	PutVacancyToArrayById(id string, wg *sync.WaitGroup, vacancies *[]models.Vacancy)
}