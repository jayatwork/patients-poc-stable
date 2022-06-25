package main

//single main package where all src live for now

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Patient struct {
	ID            int    `json:"id"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	StreetAddress string `json:"address"`
	State         string `json:"state"`
	City          string `json:"city"`
	Zip           int    `json:"zip"`
	Telephone     int    `json:"telephone"`
}

type Billing struct {
	primaryCC   int64
	billingAddr string
	owing       bool
	balance     float32
}

type PatientRecord struct {
	patient          Patient
	demographic      string
	medHistory       []string
	labResults       []string
	mentalHealth     string
	insuranceCarrier string
	billing          Billing
}

var PatientRecords []PatientRecord

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

// Create new Patient Record
func PatientsCreate(w http.ResponseWriter, r *http.Request) {
	var pr PatientRecord                        //initialize with struct
	json.NewDecoder(r.Body).Decode(&pr)         //passing JSON object to struct
	pr.patient.ID = rand.Intn(100)              //Generate random ID
	PatientRecords = append(PatientRecords, pr) //save user into all users in memory
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("user Created")
}

// Return a single patient record by patient.ID
func PatientsFind(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["ID"]

	// Loop over all of our existing PatientRecords
	for _, pr := range PatientRecords {
		id, err := strconv.Atoi(key)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		if pr.patient.ID == id {
			json.NewEncoder(w).Encode(pr.patient.ID)
		}
	}
}

// Delete a patients record entry
func PatientsDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	for index, pr := range PatientRecords {
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		if pr.patient.ID == id {
			PatientRecords = append(PatientRecords[:index], PatientRecords[index+1:]...)
			json.NewEncoder(w).Encode("Successfully deleted !!!")
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode("User Not Found")
}

func PatientsEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var pr PatientRecord
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pr); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	pr.patient.ID = id

	respondWithJSON(w, http.StatusOK, pr)
}

// Get ALL Patients
func PatientsAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(PatientRecords)
}

// Helper function for errors rendering
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Server initialization and Startup

func serverStart() {
	// Construct the routing site endpoints
	router := mux.NewRouter()
	router.HandleFunc("/", PatientsAll).Methods("GET")
	router.HandleFunc("/patients", PatientsAll).Methods("GET")
	router.HandleFunc("/patients/{id:[0-9]+}", PatientsFind).Methods("GET")
	router.HandleFunc("/patients", PatientsCreate).Methods("POST")
	router.HandleFunc("/patients/{id:[0-9]+}", PatientsEdit).Methods("PUT")
	router.HandleFunc("/patients/{id:[0-9]+}", PatientsDelete).Methods("DELETE")
	router.HandleFunc("/health", HealthCheckHandler)

	http.Handle("/", router)
	// TODO need more graceful startup and server shutdown

	// Start server
	port := ":8888"
	fmt.Println("\nListening and serving up content on port " + port)
	http.ListenAndServe(port, router)
}

func main() {

	// Initializing the first and second patient objects
	p1 := &PatientRecord{
		patient: Patient{
			ID:            12345,
			FirstName:     "Jane",
			LastName:      "Doe",
			StreetAddress: "1234 Some Patient Drive,  Some City USA",
			State:         "GA",
			Zip:           12345,
			Telephone:     7700000000,
		},
		demographic:      "Pacific Islander",
		medHistory:       []string{"Some medical history 1", "Some medical history 2", "Some medical history 3"},
		labResults:       []string{"Some lab results 1", "Some lab results 2", "Some lab results 3"},
		mentalHealth:     "Some mental health assessment",
		insuranceCarrier: "XYZ UnitedHealth",
		billing: Billing{
			primaryCC:   4444000011115555, //TODO to encode this cleartext field of CCnumber
			billingAddr: "1234 Some Client Drive, Some City USA, 12345",
			owing:       false,
			balance:     000.00,
		},
	}

	p2 := &PatientRecord{
		patient: Patient{
			ID:            67890,
			FirstName:     "John",
			LastName:      "Doe",
			StreetAddress: "5678 Some Patient Drive,  Some City USA",
			State:         "GA",
			Zip:           67890,
			Telephone:     4040000000,
		},
		demographic:      "African American",
		medHistory:       []string{"Some medical history 1", "Some medical history 2", "Some medical history 3"},
		labResults:       []string{"Some lab results 1", "Some lab results 2", "Some lab results 3"},
		mentalHealth:     "Some mental health assessment",
		insuranceCarrier: "XYZ GlobalLife",
		billing: Billing{
			primaryCC:   4444222266669999, //TODO to encode this cleartext field of CCnumber
			billingAddr: "5678 Some Client Drive, Some City USA, 67890",
			owing:       true,
			balance:     850.00,
		},
	}

	fmt.Println(p1, "\n", p2)         //Observe initial data structure
	fmt.Println(p1.billing.primaryCC) //Eventually check if known cleartext fields are obfuscated

	serverStart()

}
