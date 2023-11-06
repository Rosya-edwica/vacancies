package database

import (
	"database/sql"
	"fmt"
	"strings"
	"vacancies/pkg/models"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Connection *sql.DB
	Name       string
	Host       string
	Port       string
	User       string
	Pass       string
}

func (d *DB) Connect() (connection *sql.DB) {
	connection, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.User, d.Pass, d.Host, d.Port, d.Name))
	checkErr(err)
	return
}

func (d *DB) Close() {
	d.Connection.Close()
}

func arrayToPostgresList(items []string) (result string) {
	var updatedList []string
	for _, i := range items {
		updatedList = append(updatedList, fmt.Sprintf("'%s'", i))
	}
	result = "(" + strings.Join(updatedList, ",") + ")"
	return
}

func groupVacancies(vacancies []models.Vacancy) (groups [][]models.Vacancy) {
	LIMIT := 2000
	for i := 0; i < len(vacancies); i += LIMIT {
		group := vacancies[i:]
		if len(group) >= 2000 {
			group = group[:LIMIT]
		}
		groups = append(groups, group)
	}
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
