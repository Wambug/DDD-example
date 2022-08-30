package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/patienttracker/internal/services"
	"github.com/patienttracker/pkg/logger"
)

// TODO: Error handling and logs
// TODO: Enum type for Bloodgroup i.e: A,B,AB,O
// TODO: Salt password
// TODO: Password updated at field
// TODO: Mock API calls
// TODO: Work on cancel appointments and delete appointments
// TODO: Work on Update structs on api calls
// TODO: Department Templates for admin not api calls
// TODO: Access token
// TODO: fix update fields
const version = "1.0.0"

type Server struct {
	Router   *mux.Router
	Services services.Service
	Log      *logger.Logger
}

func NewServer() *Server {
	var wait time.Duration
	mux := mux.NewRouter()
	conn := SetupDb("postgresql://postgres:secret@localhost:5432/patient_tracker?sslmode=disable")
	services := services.NewService(conn)
	logger := logger.New()
	server := Server{
		Router:   mux,
		Log:      logger,
		Services: services,
	}
	server.Routes()
	srve := http.Server{
		Addr:         "localhost:9000",
		Handler:      mux,
		ErrorLog:     log.New(logger, "", 0),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	fmt.Println("serving at port :9000")
	//srve.ListenAndServe()
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srve.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srve.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
	return &server
}

func SetupDb(conn string) *sql.DB {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	db.Ping()
	db.SetMaxOpenConns(65)
	db.SetMaxIdleConns(65)
	db.SetConnMaxLifetime(time.Hour)
	return db
}

func (server *Server) Routes() {
	http.Handle("/", server.Router)
	//server.Router.Use(server.contentTypeMiddleware)
	server.Router.HandleFunc("/v1/healthcheck", server.Healthcheck).Methods("GET")
	server.Router.HandleFunc("/v1/department", server.createdepartment).Methods("POST")
	server.Router.HandleFunc("/v1/department", server.finddepartment).Methods("GET")
	server.Router.HandleFunc("/v1/department/{id:[0-9]+}", server.deletedepartment).Methods("DELETE")
	//queryparams: ->page_id && page_size
	server.Router.HandleFunc("/v1/departments", server.findalldepartment).Methods("GET")
	server.Router.HandleFunc("/v1/department/{id:[0-9]+}", server.updatedepartment).Methods("POST")
	//queryparams: ->page_id && page_size
	server.Router.HandleFunc("/v1/{departmentname}", server.findalldoctorsbydepartment).Methods("GET")

	server.Router.HandleFunc("/v1/doctor", server.createdoctor).Methods("POST")
	server.Router.HandleFunc("/v1/doctor", server.finddoctor).Methods("GET")
	server.Router.HandleFunc("/v1/doctor/{id:[0-9]+}", server.deletedoctor).Methods("DELETE")
	//queryparams: ->page_id && page_size
	server.Router.HandleFunc("/v1/doctors", server.findalldoctors).Methods("GET")
	server.Router.HandleFunc("/v1/doctor/{id:[0-9]+}", server.updatedoctor).Methods("POST")
	server.Router.HandleFunc("/v1/doctor/{id:[0-9]+}/schedules", server.findallschedulesbydoctor).Methods("GET")
	server.Router.HandleFunc("/v1/doctor/{id:[0-9]+}/appoinmtents", server.findallappointmentsbydoctor).Methods("GET")
	server.Router.HandleFunc("/v1/doctor/{id:[0-9]+}/records", server.findallrecordsbydoctor).Methods("GET")

	server.Router.HandleFunc("/v1/patient", server.createpatient).Methods("POST")
	server.Router.HandleFunc("/v1/patient", server.findpatient).Methods("GET")
	server.Router.HandleFunc("/v1/patient/{id:[0-9]+}", server.deletepatient).Methods("DELETE")
	server.Router.HandleFunc("/v1/patient", server.findallpatients).Methods("GET")
	server.Router.HandleFunc("/v1/patient/{id:[0-9]+}", server.updatepatient).Methods("POST")
	server.Router.HandleFunc("/v1/patient/{id:[0-9]+}/appoinmtents", server.findallappointmentsbypatient).Methods("GET")
	server.Router.HandleFunc("/v1/patient/{id:[0-9]+}/records", server.findallrecordsbypatient).Methods("GET")

	server.Router.HandleFunc("/v1/schedule", server.createschedule).Methods("POST")
	server.Router.HandleFunc("/v1/schedule", server.findschedule).Methods("GET")
	server.Router.HandleFunc("/v1/schedule/{id:[0-9]+}", server.deleteschedule).Methods("DELETE")
	server.Router.HandleFunc("/v1/schedules", server.findallschedules).Methods("GET")
	server.Router.HandleFunc("/v1/schedule/{id:[0-9]+}", server.updateschedule).Methods("POST")

	server.Router.HandleFunc("/v1/appointment/patient/{id:[0-9]+}", server.createappointmentbypatient).Methods("POST")
	server.Router.HandleFunc("/v1/appointment/doctor/{id:[0-9]+}", server.createappointmentbydoctor).Methods("POST")
	server.Router.HandleFunc("/v1/appointment", server.findappointment).Methods("GET")
	server.Router.HandleFunc("/v1/appointment/{id:[0-9]+}", server.deleteappointment).Methods("DELETE")
	server.Router.HandleFunc("/v1/appointments", server.findallappointments).Methods("GET")
	server.Router.HandleFunc("/v1/appointment/{doctorid:[0-9]+}/{id:[0-9]+}", server.updateappointmentbyDoctor).Methods("POST")
	server.Router.HandleFunc("/v1/appointment/{patientid:[0-9]+}/{id:[0-9]+}", server.updateappointmentbyPatient).Methods("POST")

	server.Router.HandleFunc("/v1/record", server.createpatientrecord).Methods("POST")
	server.Router.HandleFunc("/v1/record", server.findpatientrecord).Methods("GET")
	server.Router.HandleFunc("/v1/record/{id:[0-9]+}", server.deletepatientrecord).Methods("DELETE")
	server.Router.HandleFunc("/v1/records", server.findallpatientrecords).Methods("GET")
	server.Router.HandleFunc("/v1/record/{id:[0-9]+}", server.updatepatientrecords).Methods("POST")
	err := server.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

}
func (server *Server) Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "status: available\n")
	fmt.Fprintf(w, "version: %s\n", version)
	fmt.Fprintf(w, "Environment: Production")
}
