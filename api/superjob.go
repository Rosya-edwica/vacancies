package api

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

type SuperJobResponseFound struct {
	CountVacancies int `json:"total"`
}

type SuperJobResponseAccessToken struct {
	Token string `json:"access_token"`
}
