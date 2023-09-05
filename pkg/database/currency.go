package database

import (
	"fmt"
	"vacancies/pkg/models"
	"vacancies/pkg/tools"

	"github.com/tidwall/gjson"
)


const dictionaryUrl = "https://api.hh.ru/dictionaries"

func GetCurrencies() (currencies []models.Currency) {
	json, err := tools.GetJson(dictionaryUrl, "headhunter")
	if err != nil {
		fmt.Printf("Не удалось обновить валюту. Текст сообщения: %s", err)
	}
	for _, item := range gjson.Get(json, "currency").Array() {
		currencies = append(currencies, models.Currency{
			Code: item.Get("code").String(),
			Abbr: item.Get("abbr").String(),
			Name: item.Get("name").String(),
			Rate: item.Get("rate").Float(),
		})
	}
	return
}