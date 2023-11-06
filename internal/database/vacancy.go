package database

import (
	"fmt"
	"strings"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
)

func (d *DB) SaveOneVacancy(v models.Vacancy) {
	if v.Title == "" {
		return
	}

	columns := "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	smt := fmt.Sprintf(`
		INSERT INTO 
		h_vacancy (vacancy_id, url, name, city_id, position_id, prof_areas, specs, experience, salary_from, 
			salary_to, key_skills, vacancy_date, platform) 
		VALUES %s`,
		columns)

	tx, _ := d.Connection.Begin()
	_, err := d.Connection.Exec(smt, v.Id, v.Url, v.Title, v.CityId, v.ProfessionId, v.ProfAreas, v.Specializations,
		v.Experience, v.SalaryFrom, v.SalaryTo, v.Skills, v.DateUpdate, v.Platform)
	if err != nil {
		logger.Log.Printf("Ошибка: Вакансия %s не была добавлена в базу - %s", v.Id, err)
		err = tx.Commit()
		checkErr(err)
		return
	}
	err = tx.Commit()
	checkErr(err)
	logger.Log.Printf("Успешно сохранили вакансию %s", v.Id)
}

func (d *DB) SaveManyVacancies(vacancies []models.Vacancy) {
	groups := groupVacancies(vacancies)
	for _, group := range groups {
		d.SaveVacancies(group)
	}
}

func (d *DB) SaveVacancies(vacancies []models.Vacancy) {
	if len(vacancies) == 0 {
		return
	}

	query, vals := createQueryForMultipleInsertVacanciesMYSQL(vacancies)
	if vals == nil {
		return
	}
	tx, _ := d.Connection.Begin()
	_, err := d.Connection.Exec(query, vals...)
	checkErr(err)
	tx.Commit()
	logger.Log.Printf("Успешно сохранили %d вакансий", len(vacancies))
}

func createQueryForMultipleInsertVacanciesMYSQL(vacancies []models.Vacancy) (query string, valArgs []interface{}) {
	query = `
		INSERT IGNORE INTO 
		h_vacancy (vacancy_id, name, url, city_id, position_id, prof_areas, specs, experience, salary_from, 
			salary_to, key_skills, vacancy_date, platform) 
		VALUES `

	for _, v := range vacancies {
		if v.Title == "" || v.Url == "" || v.Id == "" {
			continue
		}
		query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
		valArgs = append(valArgs, v.Id, v.Title, v.Url, v.CityId, v.ProfessionId, v.ProfAreas,
			v.Specializations, v.Experience, v.SalaryFrom, v.SalaryTo, v.Skills, v.DateUpdate, v.Platform)
	}
	query = query[0 : len(query)-1]

	// Проверяем, что мы не вернем пустой запрос
	if strings.HasSuffix(query, "VALUES ") || len(valArgs) == 0 {
		return "", nil
	}
	return
}
