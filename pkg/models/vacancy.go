package models

type Vacancy struct {
	Id              string
	Title           string
	ProfessionId    int
	CityId          int
	DateUpdate      string
	Url             string
	Experience      string
	SalaryFrom      float64
	SalaryTo        float64
	ProfAreas       string // Slice joined by |
	Skills          string
	Specializations string
	Platform 		string
}
