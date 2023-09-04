package headhunter

import (
	"fmt"
	"strings"
	"sync"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	"vacancies/pkg/tools"

	"github.com/tidwall/gjson"
)

func (api *HeadHunter) PutVacancyToArrayById(id string, wg *sync.WaitGroup, vacancies *[]models.Vacancy) {
	defer wg.Done()
	var vacancy models.Vacancy
	json, err := tools.GetJson("https://api.hh.ru/vacancies/"+id, "headhunter")
	if err != nil {
		logger.Log.Fatalf("Ошибка при подключении к странице %s.\nТекст ошибки: %s", err, "https://api.hh.ru/vacancies/"+id)
		return
	}

	salary := api.getSalary(json)
	vacancy.Platform = "hh"
	vacancy.CityId = api.getEdwicaCityID(int(gjson.Get(json, "area.id").Int()))
	vacancy.SalaryFrom = salary.From
	vacancy.ProfessionId = api.PositionId
	vacancy.SalaryTo = salary.To
	vacancy.Skills = getSkills(json)
	vacancy.Specializations = getSpecializations(json)
	vacancy.ProfAreas = getProfAreas(json)
	vacancy.Id = gjson.Get(json, "id").String()
	vacancy.Title = gjson.Get(json, "name").String()
	vacancy.Url = gjson.Get(json, "alternate_url").String()
	vacancy.Experience = gjson.Get(json, "experience.name").String()
	vacancy.DateUpdate = gjson.Get(json, "created_at").String()
	if vacancy.CityId == 0 {
		logger.Log.Printf("Ошибка: Вакансия %s не была добавлена в базу, т.к города '%s' в базе Эдвики нет", vacancy.Id, gjson.Get(json, "area.name").String())
	} else {
		*vacancies = append(*vacancies, vacancy)
	}
	fmt.Println(vacancy.Url, vacancy.CityId, gjson.Get(json, "area.name").String())
	return
}

func (api *HeadHunter) getSalary(json string) (salary models.Salary) {
	salary.Currency = gjson.Get(json, "salary.currency").String()
	salary.From = gjson.Get(json, "salary.from").Float()
	salary.To = gjson.Get(json, "salary.to").Float()

	switch salary.Currency {
	case "RUR":
		return salary
	case "":
		return models.Salary{}
	default:
		return api.convertSalaryToRUR(salary)
	}
}

func (api *HeadHunter) getEdwicaCityID(hhId int) (edwicaId int) {
	for _, city := range api.Cities {
		if city.HH_ID == hhId {
			return city.EDWICA_ID
		}
	}
	return 0
}

func getSkills(json string) (skills string) {
	var skillsList []string
	for _, item := range gjson.Get(json, "key_skills").Array() {
		skillsList = append(skillsList, item.Get("name").String())
	}
	languages := getLanguages(json)
	skillsList = append(skillsList, languages...)
	return strings.Join(skillsList, "|")
}

func getProfAreas(json string) (areas string) {
	var profAreas []string
	for _, item := range gjson.Get(json, "specializations").Array() {
		profAreas = append(profAreas, item.Get("profarea_name").String())
	}
	return strings.Join(tools.UniqueNames((profAreas)), "|")
}

func getSpecializations(json string) (specs string) {
	var specializations []string
	for _, item := range gjson.Get(json, "specializations").Array() {
		specializations = append(specializations, item.Get("name").String())
	}
	return strings.Join(tools.UniqueNames(specializations), "|")
}

func getLanguages(json string) (languages []string) {
	for _, item := range gjson.Get(json, "languages").Array() {
		lang := item.Get("name").String()
		level := item.Get("level.name").String()
		languages = append(languages, fmt.Sprintf("%s (%s)", lang, level))
	}
	return
}
