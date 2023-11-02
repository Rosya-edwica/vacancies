package superjob

import (
	"fmt"
	// "strings"
	"sync"
	"time"
	// "vacancies/pkg/logger"
	"vacancies/pkg/models"
	// "vacancies/pkg/tools"

	"github.com/tidwall/gjson"
)

// func (api *Superjob) CollectVacanciesFromPage(url string) (vacancies []models.Vacancy) {
// 	json, err := tools.GetJson(url, "superjob")
// 	if err != nil {
// 		if err.Error() == "Limit is over" {
// 			time.Sleep(time.Second * 3)
// 			json, err = tools.GetJson(url, "superjob")
// 		} else {
// 			logger.Log.Printf("ОШИБКА: %s", err.Error())
// 			return
// 		}
// 	}
// 	for _, item := range gjson.Get(json, "objects").Array() {
// 		var vacancy models.Vacancy
// 		salary := api.getSalary(item)

// 		if api.CityEdwicaId != 0 {
// 			vacancy.CityId = api.CityEdwicaId
// 			// } else {
// 			// vacancy.CityId = api.getEdwicaCityID(item)
// 		}
// 		vacancy.Id = item.Get("id").String()
// 		vacancy.Platform = "superjob"
// 		vacancy.ProfessionId = api.PositionId
// 		vacancy.Title = item.Get("profession").String()
// 		vacancy.SalaryFrom = salary.From
// 		vacancy.SalaryTo = salary.To
// 		vacancy.Url = item.Get("link").String()
// 		vacancy.ProfAreas = strings.Join(getProfAreas(item), "|")
// 		vacancy.Specializations = strings.Join(getSpecs(item), "|")
// 		vacancy.Experience = convertExperienceId(int(item.Get("experience.id").Int()))
// 		vacancy.DateUpdate = convertDateUpdate(item.Get("date_published").Int())
// 		if vacancy.CityId == 0 {
// 			logger.Log.Printf("Ошибка: Вакансия %s не была добавлена в базу, т.к города '%s' в базе Эдвики нет", vacancy.Id, gjson.Get(json, "address").String())
// 		} else {
// 			vacancies = append(vacancies, vacancy)
// 		}
// 	}
// 	return
// }

func (api *Superjob) PutVacancyToArrayById(id string, wg *sync.WaitGroup, vacancies *[]models.Vacancy) {
	return
}

func getProfAreas(vacancyJson gjson.Result) (profAreas []string) {

	for _, item := range vacancyJson.Get("catalogues").Array() {
		profAreas = append(profAreas, item.Get("title").String())
	}
	return
}
func getSpecs(vacancyJson gjson.Result) (specs []string) {
	for _, profArea := range vacancyJson.Get("catalogues").Array() {
		for _, item := range profArea.Get("positions").Array() {
			specs = append(specs, item.Get("title").String())
		}
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

func (api *Superjob) getSalary(jsonVacancy gjson.Result) (salary models.Salary) {
	salary.From = jsonVacancy.Get("payment_from").Float()
	salary.To = jsonVacancy.Get("payment_to").Float()
	salary.Currency = jsonVacancy.Get("currency").String()

	switch salary.Currency {
	case "rub":
		return salary
	case "":
		return models.Salary{}
	default:
		return api.convertSalaryToRUR(salary)
	}

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
	for _, city := range api.Cities {
		if city.SUPERJOB_ID == superjobId {
			return city.EDWICA_ID
		}
	}
	fmt.Println("Не удалось подобрать город")
	return
}
