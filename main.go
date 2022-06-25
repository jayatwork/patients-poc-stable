package main

//single main package where all src live for now

import (
	"fmt"
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

//TODO create all other fields to support data organization

func main() {
	// Initializing the first patient object
	p := &PatientRecord{
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
			owing:       true,
			balance:     850.00,
		},
	}
	fmt.Println(p)                   //Observe initial data structure
	fmt.Println(p.billing.primaryCC) //Eventually check if known cleartext fields are obfuscated

}
