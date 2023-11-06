package geekjob

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"vacancies/pkg/models"
	"vacancies/tools"

	"github.com/gocolly/colly"
)

func (api *GeekJob) PutVacancyToArrayById(id string, wg *sync.WaitGroup, vacancies *[]models.Vacancy) {
	var vacancy models.Vacancy
	defer wg.Done()
	url := "https://geekjob.ru/vacancy/" + id
	html := getHTMLBody(url)
	vacancy.Title = getTitle(html)
	if vacancy.Title == "" {
		fmt.Println("Archived: ", url)
		return
	}
	vacancy.Url = url
	vacancy.Id = id
	vacancy.Platform = "geekjob"
	vacancy.ProfessionId = api.PositionId

	vacancy.CityId = api.getCityId(html)
	vacancy.Specializations = getSpecialization(html)
	salary := api.getSalary(html)
	vacancy.SalaryFrom = salary.From
	vacancy.SalaryTo = salary.To
	vacancy.Experience = getExperience(html)
	vacancy.DateUpdate = getDateUpdate(html)
	vacancy.Skills = strings.Join(getSkills(html), "|")
	*vacancies = append(*vacancies, vacancy)
	return
}

// Подгоняем к единому формату в БД
func getExperience(html *colly.HTMLElement) string {
	blockExperience := strings.Split(html.ChildText("span.jobformat"), "\n")
	var text string
	if len(blockExperience) > 1 {
		text = blockExperience[1]
	} else {
		text = blockExperience[0]
	}
	if strings.Contains(text, "любой") || strings.Contains(text, "менее 1 года") {
		return "Нет опыта"
	} else if strings.Contains(text, "от 1 года до 3х") {
		return "От 1 года до 3 лет"
	} else if strings.Contains(text, "от 3 до 5 лет") {
		return "От 3 до 6 лет"
	} else if strings.Contains(text, "более 5 лет") {
		return "Более 6 лет"
	}
	return "Нет опыта"
}

func getSpecialization(html *colly.HTMLElement) string {
	tags := strings.Split(html.ChildText("div.tags"), "•")
	return tags[0]
}
func getSkills(html *colly.HTMLElement) []string {
	tags := strings.Split(html.ChildText("div.tags"), "•")[1:]
	return removeAreasFromSkills(tags)
}
func getTitle(html *colly.HTMLElement) (title string) {
	title = html.ChildText("h1")
	return
}

func (api *GeekJob) getCityId(html *colly.HTMLElement) (id int) {
	cityName := html.ChildText("div.location")
	if cityName == "" {
		return 0
	}
	cityName = strings.Split(cityName, ",")[0]
	for _, item := range api.Cities {
		if strings.ToLower(item.Name) == strings.ToLower(cityName) {
			return item.EDWICA_ID
		}
	}
	return
}

func (api *GeekJob) getSalary(html *colly.HTMLElement) (salary models.Salary) {
	var salaryText string
	reDigits := regexp.MustCompile(`\d+`)
	reCurrency := regexp.MustCompile(`$|€|₽`)
	html.ForEach("span.salary", func(i int, h *colly.HTMLElement) {
		if i == 0 {
			salaryText = strings.ReplaceAll(h.Text, " ", "")
		}
	})
	currency := reCurrency.FindString(salaryText)
	if strings.Contains(salaryText, "от") && strings.Contains(salaryText, "до") {
		digits := reDigits.FindAllString(salaryText, 2)
		from, err := strconv.ParseFloat(digits[0], 64)
		tools.CheckErr(err)
		to, err := strconv.ParseFloat(digits[1], 64)
		tools.CheckErr(err)
		salary = models.Salary{
			From:     from,
			To:       to,
			Currency: currency,
		}

	} else if strings.Contains(salaryText, "от") {
		digit := reDigits.FindString(salaryText)
		from, err := strconv.ParseFloat(digit, 64)
		tools.CheckErr(err)
		salary = models.Salary{
			From:     from,
			Currency: currency,
		}
	} else {
		digit := reDigits.FindString(salaryText)
		to, _ := strconv.ParseFloat(digit, 64)
		salary = models.Salary{
			To:       to,
			Currency: currency,
		}
	}
	return api.convertSalaryToRUR(salary)
}

func (api *GeekJob) convertSalaryToRUR(salary models.Salary) models.Salary {
	if salary.Currency == "₽" {
		return models.Salary{
			From:     salary.From,
			To:       salary.To,
			Currency: "RUR",
		}
	}
	for _, currency := range api.Currencies {
		if currency.Code == "EUR" && salary.Currency == "€" {
			return models.Salary{
				From:     salary.From / currency.Rate,
				To:       salary.To / currency.Rate,
				Currency: currency.Code,
			}
		} else if currency.Code == "USR" && salary.Currency == "$" {
			return models.Salary{
				From:     salary.From / currency.Rate,
				To:       salary.To / currency.Rate,
				Currency: currency.Code,
			}
		}
	}
	return salary
}

func getDateUpdate(html *colly.HTMLElement) string {
	date := html.ChildText("div.time")
	return date + " 2023"
}

func removeAreasFromSkills(skills []string) (updated []string) {
	for _, item := range skills {
		if !checkSkillIsArea(item) {
			updated = append(updated, item)
		}
	}
	return
}

func checkSkillIsArea(skill string) bool {
	areas := []string{
		"Торговля и общепит",
		"СМИ, Медиа и индустрия развлечений",
		"Образование",
		"Заказная разработка",
		"Производство",
		"Промышленность",
		"Логистика и транспорт",
		"Медицина и фармацевтика",
		"Телекоммуникации",
		"Строительство и недвижимость",
		"Банковская и страховая сфера",
		"Наука",
		"Сельское хозяйство",
		"Консалтинг, профессиональные услуги",
		"Культура и искусство",
		"Государственные проекты",
	}
	for _, area := range areas {
		if strings.TrimSpace(area) == strings.TrimSpace(skill) {
			return true
		}
	}
	return false
}
