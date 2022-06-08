package controllers

import (
	"testing"

	"github.com/patienttracker/internal/models"
	"github.com/patienttracker/internal/utils"
	"github.com/stretchr/testify/require"
)

func CreateAppointment() models.Appointment {
	time := utils.Randate()
	patient := RandPatient()
	patient1, _ := controllers.Patient.Create(patient)
	physcian := RandDoctor()
	doc, _ := controllers.Doctors.Create(physcian)
	appointment, _ := controllers.Appointment.Create(models.Appointment{
		Patientid:       patient1.Patientid,
		Doctorid:        doc.Physicianid,
		Appointmentdate: time,
	})
	return appointment
}

func TestCreateNewAppointment(t *testing.T) {
	time := utils.Randate()
	patient := RandPatient()
	patient1, _ := controllers.Patient.Create(patient)
	physcian := RandDoctor()
	doc, _ := controllers.Doctors.Create(physcian)
	appointment, err := controllers.Appointment.Create(models.Appointment{
		Patientid:       patient1.Patientid,
		Doctorid:        doc.Physicianid,
		Appointmentdate: time,
	})
	require.NoError(t, err)
	require.Equal(t, appointment.Patientid, doc.Physicianid)
}

func TestFindAppointment(t *testing.T) {
	appointment := CreateAppointment()
	schedule, err := controllers.Appointment.Find(appointment.Appointmentid)
	require.NoError(t, err)
	require.NotEmpty(t, appointment)
	require.Equal(t, appointment.Appointmentdate, schedule.Appointmentdate)
}

func TestListAppointments(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateAppointment()

	}
	appointment, err := controllers.Appointment.FindAll()
	require.NoError(t, err)
	for _, v := range appointment {
		require.NotNil(t, v)
		require.NotEmpty(t, v)
	}

}

func TestDeleteAppointments(t *testing.T) {
	appointment := CreateAppointment()
	err := controllers.Appointment.Delete(appointment.Appointmentid)
	require.NoError(t, err)
	schedule, err := controllers.Appointment.Find(appointment.Appointmentid)
	require.Error(t, err)
	require.Empty(t, schedule)
}

func TestUpdateAppointment(t *testing.T) {
	appointment := CreateAppointment()
	time := utils.Randate()
	updatedtime, err := controllers.Appointment.Update(time, appointment.Appointmentid)
	require.NoError(t, err)
	require.NotEqual(t, appointment.Appointmentdate, updatedtime)
}