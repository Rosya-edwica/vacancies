package tools

import (
	"vacancies/pkg/models"
	"vacancies/pkg/telegram"

	apiJson "vacancies/api"
)

func UniqueNames(names []string) (unique []string) {
	allKeys := make(map[string]bool)
	for _, item := range names {
		if _, value := allKeys[item]; !value {
			if len(item) < 2 {
				continue
			}
			allKeys[item] = true
			unique = append(unique, item)
		}
	}
	return
}

func CheckErr(err error) {
	if err != nil {
		telegram.ErrorMessageMailing(err.Error())
		panic(err)
	}
}

func CollectCurrencies() (currencies []models.Currency) {
	resp, _ := apiJson.DecondeJsonResponse("https://api.hh.ru/dictionaries", nil, &apiJson.HeadHunterResponseCurrency{})
	data := resp.(*apiJson.HeadHunterResponseCurrency)
	for _, currency := range data.Items {
		currencies = append(currencies, models.Currency{
			Abbr: currency.Abbr,
			Name: currency.Name,
			Code: currency.Code,
			Rate: currency.Rate,
		})
	}
	return
}
