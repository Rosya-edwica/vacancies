package superjob

import (
	"net/url"
	"os"
	"strconv"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	"vacancies/pkg/tools"

	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
)

type Superjob struct {
	PositionName string
	PositionId   int
	CityId       int
	CityEdwicaId int
	Cities       []models.City
	Currencies   []models.Currency
}

type SuperJobResponse struct {
	Items []struct {
		Id          int    `json:"id"`
		SalaryFrom  int    `json:"payment_from"`
		SalaryTo    int    `json:"payment_to"`
		Currency    string `json:"currency"`
		PublishedAt int64  `json:"date_published"`
		Name        string `json:"profession"`
		Url         string `json:"link"`
		Experience  struct {
			Id int `json:"id"`
		} `json:"experience"`

		City struct {
			Id   int    `json:"id"`
			Name string `json:"title"`
		} `json:"town"`

		ProfAreas []struct {
			Name            string `json:"title"`
			Specializations []struct {
				Name string `json:"title"`
			} `json:"positions"`
		} `json:"catalogues"`
	} `json:"objects"`
}

func (api *Superjob) CreateQuery() (query string) {
	params := url.Values{
		"count":             {"100"},
		"keywords[0][srws]": {"1"},              // Ищем в названии вакансии
		"keywords[0][skwc]": {"particular"},     // Ищем точную фразу
		"keywords[0][keys]": {api.PositionName}, // Фраза
	}
	if api.CityId != 0 {
		params.Add("town", strconv.Itoa(api.CityId))
	}

	return "https://api.superjob.ru/2.0/vacancies?" + params.Encode()
}

// FIXME: НЕ работает. Не получается обновить .env
func (api *Superjob) UpdateAccessToken() (token string) {
	params := url.Values{
		"refresh_token": {tools.SUPERJOB_TOKEN},
		"client_id":     {tools.SUPERJOB_ID},
		"client_secret": {tools.SUPERJOB_SECRET},
	}
	json, err := tools.GetJson("https://api.superjob.ru/2.0/oauth2/refresh_token?"+params.Encode(), "superjob")
	tools.CheckErr(err)
	token = gjson.Get(json, "access_token").String()
	tools.SUPERJOB_TOKEN = token

	err = godotenv.Load(".env")
	tools.CheckErr(err)
	err = os.Setenv("SUPERJOB_TOKEN", tools.SUPERJOB_TOKEN)
	tools.CheckErr(err)
	return
}

func (api *Superjob) CountVacanciesByQuery(url string) (count int) {
	json, err := tools.GetJson(url, "superjob")
	tools.CheckErr(err)
	count = int(gjson.Get(json, "total").Int())
	logger.Log.Printf("Нашлось %d вакансий для профессии '%s'", count, api.PositionName)
	return
}
