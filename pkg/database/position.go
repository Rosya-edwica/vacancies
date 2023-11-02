package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	"vacancies/pkg/tools"
)

func (d *DB) GetPositionsByArea(areaList []string) (positions []models.Position) {
	areas := arrayToPostgresList(areaList)
	query := fmt.Sprintf(`
		SELECT position.id, position.name, position.other_names
		FROM position
		LEFT JOIN position_to_prof_area ON position_to_prof_area.position_id=position.id
		LEFT JOIN prof_area_to_specialty ON prof_area_to_specialty.id=position_to_prof_area.area_id
		LEFT JOIN professional_area ON professional_area.id=prof_area_to_specialty.prof_area_id
		WHERE LOWER(professional_area.name) IN %s
	`, strings.ToLower(areas))

	logger.Log.Printf("Парсим профессии из этих профобластей: %s", areas)
	return d.getPositions(query)
}

func (d *DB) GetAllPositions() (positions []models.Position) {
	query := `
		SELECT id, name, other_names
		FROM position`

	logger.Log.Println("Парсим абсолютно все профессии")
	return d.getPositions(query)
}

func (d *DB) GetNewPositions() (positions []models.Position) {
	query := "SELECT id, name, other_names FROM position WHERE id >= 17178"
	logger.Log.Println("Парсим новые профессии (придуманные GPT)")
	return d.getPositions(query)
}

func (d *DB) getPositions(query string) (positions []models.Position) {
	rows, err := d.Connection.Query(query)
	tools.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		var (
			name  string
			other sql.NullString
			id    int
		)
		err = rows.Scan(&id, &name, &other)
		tools.CheckErr(err)

		prof := models.Position{
			Id:         id,
			Name:       name,
			OtherNames: strings.Split(other.String, "|"),
		}
		positions = append(positions, prof)

	}
	return
}

func (d *DB) GetProfAreasFromFile(path string) (areas []string) {
	file, err := os.Open(path)
	if err != nil {
		panic("Создайте файл prof_areas.txt для парсинга по профобластям")
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		areas = append(areas, fileScanner.Text())
	}
	file.Close()
	return
}
