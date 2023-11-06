package models

type Currency struct {
	Code string
	Abbr string
	Name string
	Rate float64
}

type Salary struct {
	From float64
	To float64
	Currency string
}