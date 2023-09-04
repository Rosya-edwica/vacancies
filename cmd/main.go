package main

import (
	"fmt"
	"time"

	"vacancies/pkg/database"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	"vacancies/pkg/telegram"
	"vacancies/pkg/tools"

)


func main() {
	startTime := time.Now().Unix()
	initSettings()
	db := initDatabase()
	defer db.Close()
	
	positions := getPositions(db)
	positionsCount := len(positions)
	for i, position := range positions {
		logger.Log.Printf("Профессия № %d:%s (Осталось: %d)", position.Id, position.Name, positionsCount - i+1)
		findPositionVacancies(position, db)
	}
	telegram.SuccessMessageMailing(fmt.Sprintf("Программа завершилась за %d секунд.", time.Now().Unix() - startTime))

}



func findPositionVacancies(position models.Position, db *database.DB) {
	professionNames := position.OtherNames
	professionNames = append(professionNames, position.Name)
	for _, name := range tools.UniqueNames(professionNames) {
		position.Name = name
		vacancies := API.CollectAllVacanciesByQuery(position, db)
		db.SaveManyVacancies(vacancies)
	}
}