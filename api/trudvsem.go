package api

type TrudvsemResponse struct {
	Results struct {
		Vacancies []struct {
			Vacancy struct {
				Id             string `json:"id"`
				Name           string `json:"job-name"`
				SalaryFrom     int    `json:"salary_min"`
				SalaryTo       int    `json:"salary_max"`
				Url            string `json:"vac_url"`
				DateUpdate     string `json:"creation-date"`
				Specialisation struct {
					Name string `json:"specialisation"`
				} `json:"category"`
				Addressses struct {
					Address []struct {
						Location string `json:"location"`
					} `json:"address"`
				} `json:"addresses"`
			} `json:"vacancy"`
		} `json:"vacancies"`
	} `json:"results"`
}

type TrudvsemResponseFound struct {
	Meta struct {
		VacanciesCount int `json:"total"`
	} `json:"meta"`
}
