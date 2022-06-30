package main

//single main package where all src live for now

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

var Patients []Patient // Globally scoped Patient collection to be consulted throughout

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println("Endpoint Hit: /health")
	io.WriteString(w, `{"alive": true}`)
}

// Create new Patient Record
func patientsCreate(w http.ResponseWriter, r *http.Request) {
	var pat Patient //initialize with struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &pat)
	Patients = append(Patients, pat)
	json.NewEncoder(w).Encode(Patients)
	log.Println("Endpoint Hit: /create")
}

// Return a single patient record by patient.ID
func PatientsFind(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["ID"]

	// Loop over all of our existing PatientRecords
	for _, pat := range Patients {
		id, err := strconv.Atoi(key)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		if pat.ID == id {
			json.NewEncoder(w).Encode(pat.ID)
		}
	}
	log.Println("Endpoint Hit: /find")
}

// Delete a patients record entry
func PatientsDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	for index, pat := range Patients {
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		if pat.ID == id {
			Patients = append(Patients[:index], Patients[index+1:]...)
			json.NewEncoder(w).Encode("Successfully deleted !!!")
			return
		}
	}
	log.Println("Endpoint Hit: /delete")
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
	var pat Patient
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pat); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	pat.ID = id

	respondWithJSON(w, http.StatusOK, pat)
}

// Get ALL Patients
func PatientsAll(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: / or /patients")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(Patients)
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

//jsonHandler returns http response in JSON format.
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Construct an initiat Patient object to render
	pat := Patient{
		ID:            67890,
		FirstName:     "John",
		LastName:      "Doe",
		StreetAddress: "5678 Some Patient Drive,  Some City USA",
		State:         "GA",
		Zip:           67890,
		Telephone:     4040000000,
	}
	log.Println("Endpoint Hit: /json")
	Patients = append(Patients, pat) //save initial patient in memory
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Patients)
	json.NewEncoder(w).Encode(pat)
}

// Server initialization and Startup

func serverStart() {
	// Construct the routing site endpoints
	router := mux.NewRouter()
	router.HandleFunc("/json", jsonHandler)
	router.HandleFunc("/", PatientsAll)
	router.HandleFunc("/patients", PatientsAll)
	router.HandleFunc("/find/{id}", PatientsFind)
	router.HandleFunc("/create", patientsCreate).Methods("POST")
	router.HandleFunc("/edit/{id}", PatientsEdit)
	router.HandleFunc("/delete/{id}", PatientsDelete).Methods("DELETE")
	router.HandleFunc("/health", HealthCheckHandler)

	// TODO need more graceful startup and server shutdown

	// Start server
	port := ":8888"
	fmt.Println("\nListening and serving up content on port " + port)
	http.ListenAndServe(port, router)
}

func main() {
	serverStart()
}
