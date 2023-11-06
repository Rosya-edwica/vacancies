package headhunter

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	apiJson "vacancies/api"
	"vacancies/pkg/models"
)

// Goroutina для одновременного парсинга сразу нескольких вакансий
func (api *HeadHunter) PutVacancyToArrayById(id string, wg *sync.WaitGroup, vacancies *[]models.Vacancy) {
	defer wg.Done()
	vacancyUrl := fmt.Sprintf("http://api.hh.ru/vacancies/%s", id)

	// Декодируем json-ответ headhunter в виде подробной информации о вакансии
	resp, _ := apiJson.DecondeJsonResponse(vacancyUrl, api.Headers, &apiJson.HeadHunterVacancyResponse{})
	vacancyResp := resp.(*apiJson.HeadHunterVacancyResponse)
	vacancy := api.steamUpHeadHunterVacancyResponse(vacancyResp)
	*vacancies = append(*vacancies, vacancy)
	return
}

func (api *HeadHunter) convertSalary(salary models.Salary) models.Salary {
	switch salary.Currency {
	case "RUR":
		return salary
	case "":
		return models.Salary{}
	default:
		return api.convertSalaryToRUR(salary)
	}
}

// id городов Эдвики и HH отличаются, во время парсинга мы используем id hh, а при сохранении в БД нам нужен id Эдвики
func (api *HeadHunter) getEdwicaCityID(hhId int) (edwicaId int) {
	for _, city := range api.Cities {
		if city.HH_ID == hhId {
			return city.EDWICA_ID
		}
	}
	return 0
}

// Парсим json-структуру вакансии в модель Vacancy
func (api *HeadHunter) steamUpHeadHunterVacancyResponse(vacancyResp *apiJson.HeadHunterVacancyResponse) models.Vacancy {
	cityId, _ := strconv.Atoi(vacancyResp.City.Id)
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
		Skills:       strings.Join(steamUpSkills(vacancyResp), "|"),
		SalaryFrom:   salary.From,
		SalaryTo:     salary.To,
	}
}

// Объединяем навыки и языки в один список
func steamUpSkills(vacancyResp *apiJson.HeadHunterVacancyResponse) (skills []string) {
	var languages []string
	for _, skill := range vacancyResp.Skills {
		skills = append(skills, skill.Name)
	}
	for _, item := range vacancyResp.Languages {
		lang := fmt.Sprintf("%s (%s)", item.Name, item.Level.Name)
		languages = append(languages, lang)
	}
	skills = append(skills, languages...)
	return
}
