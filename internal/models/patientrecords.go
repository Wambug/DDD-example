package models

import "time"

// patient record model
type (
	Patientrecords struct {
		Recordid     int
		Patienid     int
		Doctorid     int
		Date         time.Time
		Diagnosis    string
		Disease      string
		Prescription string
		Weight       string
		Nurseid      int
	}

	ListPatientRecords struct {
		Limit  int
		Offset int
	}
)

// Patientrecordsrepository represent the Patientrecords repository contract
type Patientrecordsrepository interface {
	Create(patientrecords Patientrecords) (Patientrecords, error)
	Find(id int) (Patientrecords, error)
	FindAll(ListPatientRecords) ([]Patientrecords, error)
	Count() (int, error)
	FindAllByDoctor(id int) ([]Patientrecords, error)
	FindAllByPatient(id int) ([]Patientrecords, error)
	Delete(id int) error
	Update(record Patientrecords) (Patientrecords, error)
}
