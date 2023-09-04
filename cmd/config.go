package main

import (
	"os"
	"strings"
	"vacancies/pkg/api"
	"vacancies/pkg/database"
	"vacancies/pkg/models"
	"vacancies/pkg/platform/headhunter"
	"vacancies/pkg/platform/trudvsem"

	"github.com/joho/godotenv"
)

var API api.API
var PARSING_BY_AREA bool

var ERROR_MESSAGE = strings.Join([]string{
	"\nОшибка: Запусти программу с дополнительными параметрами:",
	"1. 'headhunter' OR 'superjob' OR 'trudvsem'  - для выбора платформы парсинга",
	"2. 'area'- для парсинга по определенным профобластям, которые перечислены в файле prof_areas.txt (необязательно)",
	}, "\n")


func initDatabase() (db *database.DB) {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Создай файл с переменными окружениями .env!")
	}
	db = &database.DB{
		Host: os.Getenv("MYSQL_HOST"),
		User: os.Getenv("MYSQL_USER"),
		Name: os.Getenv("MYSQL_DATABASE"),
		Pass: os.Getenv("MYSQL_PASSWORD"),
		Port: os.Getenv("MYSQL_PORT"),
	}
	db.Connection = db.Connect()
	return
}

func initSettings() {
	sysArgs := os.Args
	if len(sysArgs) < 2 {
		panic(ERROR_MESSAGE)
	}
	switch sysArgs[1] {
		case "trudvsem" : API = &trudvsem.TrudVsem{}
		// case "superjob" : PARSING_PLATFORM = "superjob"
		case "headhunter" : API = &headhunter.HeadHunter{}
		default: panic(ERROR_MESSAGE)
	}
	if len(sysArgs) == 3 && sysArgs[2] == "area" {
		PARSING_BY_AREA = true
	}

}

func getPositions(db *database.DB) (positions []models.Position) {
	if PARSING_BY_AREA {
		areas := db.GetProfAreasFromFile("prof_areas.txt")
		positions = db.GetPositionsByArea(areas)
	} else {
		positions = db.GetAllPositions()
	}
	return
}
