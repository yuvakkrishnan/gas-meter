package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type SensorReading struct {
	Meter_id         string    `json:"Meter_id"`
	Reading          float64   `json:"Reading"`
	Unit             string    `json:"Unit"`
	Timestamp        time.Time `json:"Timestamp"`
	Status           string    `json:"Status"`
	Location         Location  `json:"Location"`
	Manufacturer     string    `json:"Manufacturer"`
	Serial_number    string    `json:"Serial_number"`
	Software_version string    `json:"Software_version"`
	Battery_level    int64     `json:"Battery_level"`
	Signal_strength  int64     `json:"Signal_strength"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

func main() {
	// Open the JSON file
	jsFile, err := os.Open("gas_meter.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer jsFile.Close()

	// Parse the JSON file
	var reading SensorReading
	err = json.NewDecoder(jsFile).Decode(&reading)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Store the message in AWS DynamoDB
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess, aws.NewConfig().WithRegion("us-east-1"))

	av, err := dynamodbattribute.MarshalMap(reading)
	if err != nil {
		fmt.Println("Error marshaling message:", err)
		return
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("SensorReadings"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Error storing message:", err)
		return
	}

	fmt.Println("Message stored successfully!")
}



brfore edited with fronend

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

		// Marshal the SensorReading objects into JSON and send the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(readings)
	})

	// Start the HTTP server
	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
#   g a s - m e t e r  
 