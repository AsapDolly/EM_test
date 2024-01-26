package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/AsapDolly/EM_test/internal/entity"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// PostgreConnect хранит соединение с базой данных.
type PostgreConnect struct {
	DBConnect *sql.DB
}

// URLRow используется для чтения данных из базы данных.
type URLRow struct {
	UserID      string
	ShortURL    string
	OriginalURL string
}

// GetNewConnection - конструктор PostgreConnect.
func GetNewConnection(db *sql.DB, dbConf string, migrationAddress string) PostgreConnect {

	dbConn := PostgreConnect{DBConnect: db}

	migration, err := migrate.New(migrationAddress, dbConf)
	if err != nil {
		log.Print(err)
	}

	if err = migration.Up(); errors.Is(err, migrate.ErrNoChange) {
		log.Print(err)
	}

	return dbConn
}

// GetPersons читает данные из базы данных.
func (s PostgreConnect) GetPersons(ctx context.Context, params map[string][]string) (res map[int]entity.Person, err error) {
	res = make(map[int]entity.Person)
	limit := ""
	query := "WITH PersonData AS (SELECT person_id, name, surname, patronymic, age, gender FROM Persons WHERE isdeleted = FALSE"

	//собираем параметры в массив с записями типа name = 'John'
	var whereClauses []string
	for key, values := range params {
		for _, value := range values {
			if strings.EqualFold(key, "limit") {
				limit = value
				continue
			}
			whereClauses = append(whereClauses, key+"='"+value+"'")
		}
	}

	//если указаны параметры, то добавляем к запросу
	if len(whereClauses) > 0 {
		query = fmt.Sprintf("%s AND  %s", query, strings.Join(whereClauses, " AND "))
	}

	if limit != "" {
		query = fmt.Sprintf("%s LIMIT  %s", query, limit)
	}

	query = fmt.Sprintf("%s) SELECT pd.person_id, pd.name, pd.surname, pd.patronymic, pd.age, pd.gender, n.nationality_name, pn.probability FROM PersonData pd LEFT JOIN PersonNationalities pn ON pd.person_id = pn.person_id LEFT JOIN Nationalities n ON pn.nationality_id = n.nationality_id;", query)

	nationalitiesRows, err := s.DBConnect.QueryContext(ctx, query)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer nationalitiesRows.Close()

	for nationalitiesRows.Next() {
		var p entity.Person
		var n entity.Nationality

		err = nationalitiesRows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &n.CountryID, &n.Probability)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		if _, exist := res[p.ID]; exist {
			p.Nationality = append(res[p.ID].Nationality, n)
		} else {
			p.Nationality = append(p.Nationality, n)
		}

		res[p.ID] = p
	}

	err = nationalitiesRows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return res, nil
}

// WritePersonData сохраняет данные в БД.
func (s PostgreConnect) WritePersonData(ctx context.Context, person entity.Person) error {
	tx, err := s.DBConnect.BeginTx(ctx, nil)
	if err != nil {
		log.Print(err)
		return err
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRow("INSERT INTO Persons (name, surname, patronymic, age, gender) VALUES ($1, $2, $3, $4, $5) RETURNING person_id;", person.Name, person.Surname, person.Patronymic, person.Age, person.Gender).Scan(&userID)
	if err != nil {
		err = fmt.Errorf("ошибка вставки в таблицу Persons: %w", err)
		log.Print(err)
		return err
	}

	err = insertPersonNationalities(tx, userID, person)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}

// UpdateData обновляет данные из БД.
func (s PostgreConnect) UpdateData(ctx context.Context, person entity.Person) error {
	tx, err := s.DBConnect.BeginTx(ctx, nil)
	if err != nil {
		log.Print(err)
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE Persons SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5 WHERE person_id = $6;", person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.ID)
	if err != nil {
		err = fmt.Errorf("ошибка обновления таблицы Persons: %w", err)
		log.Print(err)
		return err
	}

	_, err = tx.Exec("DELETE FROM PersonNationalities WHERE person_id = $1;", person.ID)
	if err != nil {
		err = fmt.Errorf("ошибка обновления таблицы Persons: %w", err)
		log.Print(err)
		return err
	}

	err = insertPersonNationalities(tx, person.ID, person)
	if err != nil {
		log.Print(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}

// DeleteData удаляет данные из БД.
func (s PostgreConnect) DeleteData(ctx context.Context, personID int) error {
	tx, err := s.DBConnect.BeginTx(ctx, nil)
	if err != nil {
		log.Print(err)
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE Persons SET isdeleted = TRUE WHERE person_id = $1;", personID)
	if err != nil {
		err = fmt.Errorf("ошибка удаление из таблицы Persons: %w", err)
		log.Print(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func insertPersonNationalities(tx *sql.Tx, userID int, person entity.Person) error {

	sqlInsertPersonNationalities, err := tx.Prepare("INSERT INTO PersonNationalities (person_id, nationality_id, probability) VALUES ($1, $2, $3);")
	if err != nil {
		log.Print(err)
		return err
	}

	defer sqlInsertPersonNationalities.Close()

	for _, values := range person.Nationality {

		var nationalityID int

		err = tx.QueryRow("WITH inserted_row AS (INSERT INTO Nationalities (nationality_name) VALUES ($1) ON CONFLICT (nationality_name) DO NOTHING RETURNING nationality_id) SELECT nationality_id FROM inserted_row UNION SELECT nationality_id FROM Nationalities WHERE nationality_name = $1;\n", values.CountryID).Scan(&nationalityID)
		if err != nil {
			err = fmt.Errorf("ошибка вставки в таблицу Nationalities: %w", err)
			log.Print(err)
			return err
		}

		_, err = sqlInsertPersonNationalities.Exec(userID, nationalityID, values.Probability)
		if err != nil {
			err = fmt.Errorf("ошибка вставки в таблицу PersonNationalities: %w", err)
			log.Print(err)
			return err
		}

	}

	return nil
}
