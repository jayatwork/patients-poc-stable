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

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Patient struct {
	Id            int64  `json:"id"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	Dob           string `json:"dob"`
	StreetAddress string `json:"address"`
	State         string `json:"state"`
	City          string `json:"city"`
	Zip           string `json:"zip"`
	Email         string `json:"email"`
	Telephone     string `json:"telephone"`
	Appointment   string `json:"appointment"`
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// Create new Patient Record
func patientsRegister(w http.ResponseWriter, r *http.Request) {

	var pat Patient //initialize with struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &pat)
	Patients = append(Patients, pat)
	json.NewEncoder(w).Encode(Patients)
	log.Println("Endpoint Hit: /register")
}

// Return a single patient record by patient.ID
func PatientsFind(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Loop over all of our existing PatientRecords
	for _, pat := range Patients {
		id, err := strconv.ParseInt(params["id"], 10, 64)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		if pat.Id == id {
			json.NewEncoder(w).Encode(pat)
		}
	}
	log.Println("Endpoint Hit: /find/<some_id>")
}

// Delete a patients record entry
func PatientsDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	for index, pat := range Patients {
		log.Println(strconv.ParseInt(params["id"], 10, 64))
		id, err := strconv.ParseInt(params["id"], 10, 64)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		if pat.Id == id {
			fmt.Println(pat)
			fmt.Println(Patients[:index], Patients[index+1:])
			Patients = append(Patients[:index], Patients[index+1:]...)
			json.NewEncoder(w).Encode("Successfully deleted - No longer Registered !!!")
			return
		}
	}
	log.Println("Endpoint Hit: /delete")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode("User Not Found")
}

func PatientsEdit(w http.ResponseWriter, r *http.Request) {

	// Get the ID from the url
	params := mux.Vars(r)
	var updatedEvent Patient
	log.Println(strconv.ParseInt(params["id"], 10, 64))
	eventID, err := strconv.ParseInt(params["id"], 10, 64)
	// Convert r.Body into a readable format
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Enter required fields FirstName Lastname Appointment time revised")
	}

	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range Patients {
		if singleEvent.Id == eventID {
			singleEvent.FirstName = updatedEvent.FirstName
			singleEvent.LastName = updatedEvent.LastName
			singleEvent.Appointment = updatedEvent.Appointment
			Patients[i] = singleEvent
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
	log.Println("Endpoint Hit: /edit")
}

// Get ALL Patients
func PatientsAll(w http.ResponseWriter, r *http.Request) {

	log.Println("Endpoint Hit: / or /patients")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
// func jsonHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	//Construct an initial Patient object to render
// 	pat := Patient{
// 		ID:            67890,
// 		FirstName:     "John",
// 		LastName:      "Doe",
// 		StreetAddress: "5678 Some Patient Drive,  Some City USA",
// 		State:         "GA",
// 		Zip:           67890,
// 		Telephone:     4040000000,
// 	}
// 	log.Println("Endpoint Hit: /json")
// 	Patients = append(Patients, pat) //save initial patient in memory
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(Patients)
// 	json.NewEncoder(w).Encode(pat)
// }

// Server initialization and Startup
func serverStart() {
	// Construct the routing site endpoints
	router := mux.NewRouter()
	//	router.HandleFunc("/json", jsonHandler)
	router.HandleFunc("/", PatientsAll)
	router.HandleFunc("/patients", PatientsAll)
	router.HandleFunc("/find/{id:[0-9]+}", PatientsFind).Methods("GET")
	router.HandleFunc("/register", patientsRegister).Methods("POST")
	router.HandleFunc("/edit/{id:[0-9]+}", PatientsEdit).Methods("PUT")
	router.HandleFunc("/delete/{id:[0-9]+}", PatientsDelete).Methods("DELETE")
	router.HandleFunc("/health", HealthCheckHandler)

	// TODO need more graceful startup and server shutdown

	// Start server
	port := ":8888"
	fmt.Println("\nListening and serving up content on port " + port)
	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED"), "*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router))
}

func main() {
	serverStart()
}
