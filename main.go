package main

//single main package where all src live for now

import (
	"encoding/json"
	"fmt"
	"io"

	//	"io/utils"
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

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
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	respondWithJSON(w, http.StatusOK, pr)
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")

	// Declare the static file directory and point it to the
	// directory we just made
	staticFileDirectory := http.Dir("./assets/")
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for
	// "./assets/assets/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	return r
}

func getPatientsHandler(w http.ResponseWriter, r *http.Request) {
	//Transform our patientrecord to json
	patientRecBytes, err := json.Marshal(PatientRecords)

	// If there is an error, print it to the console, and return a server
	// error response to the user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// If all goes well, write the JSON list of birds to the response
	w.Write(patientRecBytes)
}

func createPatientsHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new instance of Patient Record
	pat := PatientRecord{}

	// We send all our data as HTML form data
	// the `ParseForm` method of the request, parses the
	// form values
	err := r.ParseForm()

	// In case of any error, we respond with an error to the user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the information about the patient from the form info
	pat.patient.FirstName = r.Form.Get("species")
	pat.patient.LastName = r.Form.Get("description")

	// Append our existing list of birds with a new entry
	PatientRecords = append(PatientRecords, pat)

	//Finally, we redirect the user to the original HTMl page
	// (located at `/assets/`), using the http libraries `Redirect` method
	http.Redirect(w, r, "/assets/", http.StatusFound)
}

// Get ALL Patients 1st iteration
func PatientsAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//json.NewEncoder(w).Encode(PatientRecords)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&PatientRecords)
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Initializing the first and second patient objects
	p1 := PatientRecord{
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

	p2 := PatientRecord{
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

	json.NewEncoder(w).Encode(p1.patient.FirstName)
	json.NewEncoder(w).Encode(p2.patient.LastName)

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

// Print the entire PatientRecords collection currently in memory
func PrintAllStruct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var printSample PatientRecord
	printSample = PatientRecord{
		patient: Patient{
			ID:            0000000,
			FirstName:     "Hello from PrintAllStruct (see fields)",
			LastName:      "Doe",
			StreetAddress: "1234 Some Printer stuff,  Some City USA",
			State:         "GA",
			Zip:           12345,
			Telephone:     7700000000,
		},
		demographic:      "We're all the same",
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

	io.WriteString(w, printSample.patient.FirstName)
}

// Server initialization and Startup

func serverStart() {
	// Construct the routing site endpoints
	router := mux.NewRouter()
	router.HandleFunc("/", getPatientsHandler).Methods("GET")
	router.HandleFunc("/patients", getPatientsHandler).Methods("GET")
	router.HandleFunc("/patients/{id:[0-9]+}", PatientsFind).Methods("GET")
	router.HandleFunc("/create", getPatientsHandler).Methods("GET")
	router.HandleFunc("/create", createPatientsHandler).Methods("POST")
	router.HandleFunc("/patients/edit/{id:[0-9]+}", PatientsEdit).Methods("PUT")
	router.HandleFunc("/patients/delete/{id:[0-9]+}", PatientsDelete).Methods("DELETE")
	router.HandleFunc("/print", PrintAllStruct).Methods("GET")
	router.HandleFunc("/health", HealthCheckHandler)
	router.HandleFunc("/hello", handler).Methods("GET")

	http.Handle("/", router)
	// TODO need more graceful startup and server shutdown

	// Start server
	port := ":8888"
	fmt.Println("\nListening and serving up content on port " + port)
	http.ListenAndServe(port, router)
}

func main() {

	// Initializing the first and second patient objects
	p1 := PatientRecord{
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

	p2 := PatientRecord{
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
