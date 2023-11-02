package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"vacancies/pkg/api"
	"vacancies/pkg/database"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"

	"vacancies/pkg/platform/geekjob"
	"vacancies/pkg/platform/headhunter"
	"vacancies/pkg/platform/superjob"
	"vacancies/pkg/platform/trudvsem"
	"vacancies/pkg/telegram"
	"vacancies/pkg/tools"

	"github.com/joho/godotenv"
)

var API api.API
var PARSING_BY_AREA bool
var PLATFORM string
var ERROR_MESSAGE = strings.Join([]string{
	"\nОшибка: Запусти программу с дополнительными параметрами:",
	"1. 'headhunter' OR 'superjob' OR 'trudvsem' OR 'geekjob'  - для выбора платформы парсинга",
	"2. 'area'- для парсинга по определенным профобластям, которые перечислены в файле prof_areas.txt (необязательно)",
}, "\n")

func main() {
	startTime := time.Now().Unix()
	initSettings()
	db := initDatabase()
	defer db.Close()

	positions := getPositions(db)
	positionsCount := len(positions)
	for i, position := range positions {
		logger.Log.Printf("[%s] Профессия № %d: %s (Осталось: %d)", PLATFORM, position.Id, position.Name, positionsCount-i+1)
		log.Printf("[%s] Профессия № %d: %s (Осталось: %d)", PLATFORM, position.Id, position.Name, positionsCount-i+1)
		findPositionVacancies(position, db)
	}
	telegram.SuccessMessageMailing(fmt.Sprintf("Программа завершилась за %d секунд.", time.Now().Unix()-startTime))

}

func findPositionVacancies(position models.Position, db *database.DB) {
	professionNames := position.OtherNames
	professionNames = append(professionNames, position.Name)
	for _, name := range tools.UniqueNames(professionNames) {
		position.Name = name
		vacancies := API.CollectAllVacanciesByQuery(position, db)
		logger.Log.Printf("[%s] Количество вакансий для %s:%d\n", PLATFORM, position.Name, len(vacancies))
		db.SaveManyVacancies(vacancies)
	}
}

func initDatabase() (db *database.DB) {
	err := godotenv.Load(".env")
	if err != nil {
		tools.CheckErr(errors.New("Создай файл с переменными окружениями"))
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
	PLATFORM = sysArgs[1]
	switch PLATFORM {
	case "trudvsem":
		API = &trudvsem.TrudVsem{}
	case "superjob":
		API = &superjob.Superjob{}
	case "headhunter":
		API = &headhunter.HeadHunter{}
	case "geekjob":
		API = &geekjob.GeekJob{}
	default:
		panic(ERROR_MESSAGE)
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
		positions = db.GetNewPositions()
	}
	return
}
