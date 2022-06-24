package main

//single main package where all src live for now

import (
	"log"
)

type Patient struct {
	ID            int    `json:"id"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	StreetAddress string `json:"address"`
	State         string `json:"state"`
	ZIP           int    `json:"zip"`
	Telephone     int    `json:"telephone"`
}

type Billing struct {
}

type PatientRecord struct {
	patient          Patient
	demographic      string
	medhistory       []string
	labresults       []string
	mentalhealth     string
	insurancecarrier string
	billing          Billing
}

//TODO create all other fields to support data organization

func main() {
	// Initializing the first patient object
	dummyp := PatientRecord{}
	//TODO complete the populating all fields of nested structures

	if err != nil {
		log.Fatalf("Unable to render the initial patient record %v", err)
	}

}
