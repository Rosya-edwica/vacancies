package headhunter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	"vacancies/pkg/tools"

	"github.com/tidwall/gjson"
)

type HeadHunter struct {
	PositionName string
	PositionId   int
	CityId       int
	CityEdwicaId int
	Cities       []models.City
	Currencies   []models.Currency
}

type HeadHunterResponse struct {
	Items []struct {
		Id string `json:"id"`
	} `json:"items"`
}

type HeadHunterVacancyResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Url       string `json:"alternate_url"`
	Salary    struct {
		From     int    `json:"from"`
		To       int    `json:"to"`
		Currency string `json:"currency"`
	} `json:"salary"`
	Experience struct {
		Name string `json:"name"`
	} `json:"experience"`
	City struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"area"`
	Skills []struct {
		Name string `json:"name"`
	} `json:"key_skills"`
	Languages []struct {
		Name  string `json:"name"`
		Level struct {
			Name string `json:"name"`
		} `json:"level"`
	} `json:"languages"`
}

const (
	TOKEN        = "QQAVSIBVU4B0JCR296THKB22JP05A92H329U49TDD9CRIS8DT9BRPPT7M9OLQ6HD"
	per_page     = "60"
	search_field = "name"
)

func (api *HeadHunter) CreateQuery() (query string) {
	var params url.Values
	if api.CityId == 0 {
		params = url.Values{
			"search_field": {search_field},
			"per_page":     {per_page},
			"text":         {api.PositionName},
		}
	} else {
		params = url.Values{
			"search_field": {search_field},
			"per_page":     {per_page},
			"text":         {api.PositionName},
			"area":         {strconv.Itoa(api.CityId)},
		}
	}
	query = "https://api.hh.ru/vacancies?" + params.Encode()
	return
}

func (api *HeadHunter) CountVacanciesByQuery(url string) (count int) {
	json, err := tools.GetJson(url, "headhunter")
	tools.CheckErr(err)
	count = int(gjson.Get(json, "found").Int())
	logger.Log.Printf("Нашлось %d вакансий для профессии '%s'", count, api.PositionName)
	return
}

func (api *HeadHunter) convertSalaryToRUR(salary models.Salary) models.Salary {
	for _, cur := range api.Currencies {
		fmt.Println(cur.Code)
		if cur.Code == salary.Currency {
			salary.To = salary.To / cur.Rate
			salary.From = salary.From / cur.Rate
			salary.Currency = "RUR"
			return salary
		}
	}
	return salary
}

func TestParse() {
	var v HeadHunterVacancyResponse

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get("http://api.hh.ru/vacancies/88408515")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}
