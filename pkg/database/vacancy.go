package database

import (
	"fmt"
	"vacancies/pkg/logger"
	"vacancies/pkg/models"
	"vacancies/pkg/tools"
)

func (d *DB) SaveOneVacancy(v models.Vacancy) {
	if v.Title == "" { return }

	columns := "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	smt := fmt.Sprintf(`
		INSERT INTO 
		h_vacancy (id, url, name, city_id, position_id, prof_areas, specs, experience, salary_from, 
			salary_to, key_skills, vacancy_date, platform) 
		VALUES %s`, 
		columns)

	tx, _ := d.Connection.Begin()
	_, err := d.Connection.Exec(smt, v.Id, v.Url, v.Title, v.CityId, v.ProfessionId, v.ProfAreas, v.Specializations,
		 v.Experience, v.SalaryFrom, v.SalaryTo, v.Skills, v.DateUpdate, v.Platform)
	if err != nil {
		logger.Log.Printf("Ошибка: Вакансия %d не была добавлена в базу - %s", v.Id, err)
		err = tx.Commit()
		tools.CheckErr(err)
		return
	}
	err = tx.Commit()
	tools.CheckErr(err)
	logger.Log.Printf("Успех: Вакансия %s была добавлена в базу", v.Id)
}

func (d *DB) SaveManyVacancies(vacancies []models.Vacancy) {
	groups := groupVacancies(vacancies)
	for _, group := range groups {
		d.SaveVacancies(group)
	}
}

func (d *DB) SaveVacancies(vacancies []models.Vacancy) {
	if len(vacancies) == 0 { return }

	query, vals := createQueryForMultipleInsertVacanciesMYSQL(vacancies)
	tx, _ := d.Connection.Begin()
	_, err := d.Connection.Exec(query, vals...)
	tools.CheckErr(err)
	tx.Commit()
	logger.Log.Printf("Успех: Cохранили %d вакансий", len(vacancies))
}


func createQueryForMultipleInsertVacanciesMYSQL(vacancies []models.Vacancy) (query string, valArgs []interface{}) {
	query = `
		INSERT IGNORE INTO 
		h_vacancy (id, name, url, city_id, position_id, prof_areas, specs, experience, salary_from, 
			salary_to, key_skills, vacancy_date, platform) 
		VALUES `

	for _, v := range vacancies {
		query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
		valArgs = append(valArgs,  v.Id, v.Title, v.Url, v.CityId, v.ProfessionId, v.ProfAreas, 
			v.Specializations, v.Experience, v.SalaryFrom, v.SalaryTo, v.Skills, v.DateUpdate, v.Platform)
	}
	query = query[0:len(query)-1]
	return
}