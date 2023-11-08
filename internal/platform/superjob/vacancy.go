package superjob

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"vacancies/pkg/models"

	apiJson "vacancies/api"
)

func (api *Superjob) PutVacancyToArrayById(id string, wg *sync.WaitGroup, vacancies *[]models.Vacancy) {
	return
}

func (api *Superjob) CollectVacanciesFromPage(url string) (vacancies []models.Vacancy) {
	resp, statusCode := apiJson.DecondeJsonResponse(url, api.Headers, &apiJson.SuperJobResponse{})
	fmt.Println(statusCode, url)
	vacanciesResp := resp.(*apiJson.SuperJobResponse)
	for _, item := range vacanciesResp.Items {
		var vacancy models.Vacancy
		var profAreas, specs []string

		if api.CityEdwicaId != 0 {
			vacancy.CityId = api.CityEdwicaId
		} else {
			vacancy.CityId = api.getEdwicaCityID(item.City.Id)
		}

		for _, area := range item.ProfAreas {
			profAreas = append(profAreas, area.Name)
			for _, spec := range area.Specializations {
				specs = append(specs, spec.Name)
			}
		}

		salary := api.convertSalaryToRUR(models.Salary{
			From:     float64(item.SalaryFrom),
			To:       float64(item.SalaryTo),
			Currency: item.Currency,
		})

		vacancy = models.Vacancy{
			Platform:        "superjob",
			ProfessionId:    api.PositionId,
			Id:              strconv.Itoa(item.Id),
			Title:           item.Name,
			Url:             item.Url,
			DateUpdate:      convertDateUpdate(item.PublishedAt),
			Experience:      convertExperienceId(item.Experience.Id),
			SalaryFrom:      salary.From,
			SalaryTo:        salary.To,
			ProfAreas:       strings.Join(profAreas, "|"),
			Specializations: strings.Join(specs, "|"),
		}
		vacancies = append(vacancies, vacancy)

	}
	return
}

func convertDateUpdate(timestamp int64) (date string) {
	dateUpdate := time.Unix(timestamp, 0)
	date = dateUpdate.Format(time.RFC3339Nano)
	return
}

func convertExperienceId(id int) (experience string) {
	switch id {
	case 2:
		experience = "От 1 года до 3 лет"
	case 3:
		experience = "От 3 до 6 лет"
	case 4:
		experience = "От 6 лет"
	default:
		experience = "Нет опыта"
	}
	return
}

func (api *Superjob) convertSalaryToRUR(salary models.Salary) (convertedSalary models.Salary) {
	for _, cur := range api.Currencies {
		if cur.Abbr == salary.Currency {
			salary.To = salary.To / cur.Rate
			salary.From = salary.From / cur.Rate
			salary.Currency = "rub"
			return salary
		}
	}

	return
}

func (api *Superjob) getEdwicaCityID(superjobId int) (id int) {
	if superjobId != 0 {
		return superjobId
	}
	for _, city := range api.Cities {
		if city.SUPERJOB_ID == superjobId {
			return city.EDWICA_ID
		}
	}
	fmt.Println("Не удалось подобрать город")
	return
}
