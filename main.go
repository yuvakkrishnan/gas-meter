package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type SensorReading struct {
	Meter_id         string    `json:"Meter_id"`
	Reading          float64   `json:"reading"`
	Unit             string    `json:"unit"`
	Timestamp        time.Time `json:"timestamp"`
	Status           string    `json:"status"`
	Location         Location  `json:"location"`
	Manufacturer     string    `json:"manufacturer"`
	Serial_number    string    `json:"serial_number"`
	Software_version string    `json:"software_version"`
	Battery_level    int64     `json:"battery_level"`
	Signal_strength  int64     `json:"signal_strength"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

func main() {
	// Set up a HTTP handler
	http.HandleFunc("/readings", func(w http.ResponseWriter, r *http.Request) {
		// Connect to DynamoDB
		sess := session.Must(session.NewSession())
		svc := dynamodb.New(sess, aws.NewConfig().WithRegion("us-east-1"))

		// Query DynamoDB for all items in the SensorReadings table
		input := &dynamodb.ScanInput{
			TableName: aws.String("SensorReadings"),
		}
		result, err := svc.Scan(input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal the DynamoDB items into SensorReading objects
		readings := make([]SensorReading, 0)
		for _, item := range result.Items {
			reading := SensorReading{}
			err = dynamodbattribute.UnmarshalMap(item, &reading)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			readings = append(readings, reading)
		}

		// Render the HTML template with the SensorReading data and send the response
		tmpl := template.Must(template.ParseFiles("template.html"))
		err = tmpl.Execute(w, readings)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Start the HTTP server
	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
