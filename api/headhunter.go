package api

type HeadHunterResponse struct {
	Items []struct {
		Id string `json:"id"`
	} `json:"items"`
}

type HeadHunterResponseFound struct {
	CountVacancies int `json:"found"`
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

type HeadHunterResponseCurrency struct {
	Items []struct {
		Code string  `json:"code"`
		Abbr string  `json:"abbr"`
		Name string  `json:"name"`
		Rate float64 `json:"rate"`
	} `json:"currency"`
}
