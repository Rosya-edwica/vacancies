package database

import (
	"vacancies/pkg/models"
)

func (d *DB) GetCities() (cities []models.City) {
	query := `
		SELECT id_edwica, id_superjob, id_hh, name
		FROM h_city 
		ORDER BY id_hh ASC`
	rows, err := d.Connection.Query(query)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var name string
		var edwicaId, superjobId, hhId int

		err = rows.Scan(&edwicaId, &superjobId, &hhId, &name)
		checkErr(err)
		cities = append(cities, models.City{
			Name:        name,
			HH_ID:       hhId,
			EDWICA_ID:   edwicaId,
			SUPERJOB_ID: superjobId,
		})
	}
	return
}
