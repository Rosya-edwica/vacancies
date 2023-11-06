package api

type GeekJobResponseFound struct {
	CountVacancies int `json:"documentsCount"`
}

type GeekJobResponse struct {
	Items []struct {
		Id string `json:"id"`
	} `json:"data"`
}
