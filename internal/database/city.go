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
		var edwica_id, superjob_id, hh_id int

		err = rows.Scan(&edwica_id, &superjob_id, &hh_id, &name)
		checkErr(err)
		cities = append(cities, models.City{
			Name:        name,
			HH_ID:       hh_id,
			EDWICA_ID:   edwica_id,
			SUPERJOB_ID: superjob_id,
		})
	}
	return
}
