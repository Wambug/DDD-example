package controllers

import (
	"context"
	"database/sql"
	"log"

	"github.com/patienttracker/internal/models"
)

type Nurse struct {
	db *sql.DB
}

func (n Nurse) Create(nurse models.Nurse) (models.Nurse, error) {
	sqlStatement := `
  INSERT INTO nurse (username,full_name,email,hashed_password) 
  VALUES($1,$2,$3,$4)
  RETURNING *
  `
	err := n.db.QueryRow(sqlStatement, nurse.Username, nurse.Full_name,
		nurse.Email, nurse.Hashed_password).Scan(
		&nurse.Id,
		&nurse.Username,
		&nurse.Full_name,
		&nurse.Email,
		&nurse.Hashed_password,
		&nurse.Password_changed_at,
		&nurse.Created_at,
	)
	if err != nil {
		log.Fatal(err)
	}
	return nurse, nil

}

func (n Nurse) Find(id int) (models.Nurse, error) {
	sqlStatement := `
  SELECT * FROM nurse
  WHERE nurse.id = $1
  `
	var nurse models.Nurse
	err := n.db.QueryRowContext(context.Background(), sqlStatement, id).Scan(
		&nurse.Id,
		&nurse.Username,
		&nurse.Full_name,
		&nurse.Email,
		&nurse.Hashed_password,
		&nurse.Password_changed_at,
		&nurse.Created_at)
	return nurse, err
}

func (n Nurse) FindbyEmail(email string) (models.Nurse, error) {
	sqlStatement := `
  SELECT * FROM nurse
  WHERE nurse.email = $1
  `
	var nurse models.Nurse
	err := n.db.QueryRowContext(context.Background(), sqlStatement, email).Scan(
		&nurse.Id,
		&nurse.Username,
		&nurse.Full_name,
		&nurse.Email,
		&nurse.Hashed_password,
		&nurse.Password_changed_at,
		&nurse.Created_at)
	return nurse, err
}

func (n Nurse) Count() (int, error) {

	counter := 0
	rows, err := n.db.Query("SELECT * FROM nurse")
	if err != nil {
		return counter, err
	}
	defer rows.Close()

	for rows.Next() {
		counter++
	}
	return counter, nil
}

func (p Nurse) FindAll(args models.ListNurses) ([]models.Nurse, error) {
	sqlStatement := `
 SELECT id, username,full_name,email FROM nurse
 ORDER BY id
 LIMIT $1
 OFFSET $2
  `
	rows, err := p.db.QueryContext(context.Background(), sqlStatement, args.Limit, args.Offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var items []models.Nurse
	for rows.Next() {
		var i models.Nurse
		if err := rows.Scan(
			&i.Id,
			&i.Username,
			&i.Full_name,
			&i.Email,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (n Nurse) Delete(id int) error {
	sqlStatement := `DELETE FROM nurse
  WHERE id  = $1
  `
	_, err := n.db.Exec(sqlStatement, id)
	return err
}

func (p Nurse) Update(nurse models.Nurse) (models.Nurse, error) {
	sqlStatement := `UPDATE nurse
SET username = $2, full_name = $3, email = $4,hashed_password=$5,password_changed_at=$6
WHERE id = $1
RETURNING id,full_name,username,email;
  `
	var nur models.Nurse
	err := p.db.QueryRow(sqlStatement, nurse.Id, nurse.Username, nurse.Full_name, nurse.Email, nurse.Hashed_password, nurse.Password_changed_at).Scan(
		&nur.Id,
		&nur.Full_name,
		&nur.Username,
		&nur.Email,
	)
	if err != nil {
		return nur, err
	}
	return nur, nil
}
