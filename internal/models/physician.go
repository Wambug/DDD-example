package models

import "time"

//Physcian struct
type Physician struct {
	Physicianid         int
	Username            string
	Full_name           string
	Email               string
	Contact             string
	Hashed_password     string
	Password_changed_at time.Time
	Created_at          time.Time
	//verfied string
}

//update Physcian
type UpdatePhysician struct {
	Username            string
	Full_name           string
	Email               string
	Contact             string
	Hashed_password     string
	Password_changed_at time.Time
}

//Physicianrepository represent the Physician repository contract
type Physicianrepository interface {
	Create(physician Physician) (Physician, error)
	Find(id int) (Physician, error)
	FindAll() ([]Physician, error)
	Delete(id int) error
	Update(physician UpdatePhysician, id int) (Physician, error)
}
