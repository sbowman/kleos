package main

import (
	"time"

	"github.com/sbowman/kleos"
)

func main() {
	writer := kleos.NewLogstashWriter(":9999", time.Second)
	if err := writer.Dial(); err != nil {
		panic(err)
	}
	defer writer.Close()

	elk := kleos.NewJSONOutput("testing", writer)
	elk.Timestamp = "@timestamp"
	kleos.SetOutput(elk)

	// tw := kleos.NewTextOutput(os.Stdout)
	// kleos.SetOutput(tw)

	kleos.WithFields(kleos.Fields{
		"name": "James T. Kirk",
		"rank": "Captain",
		"assignment": "U.S.S. Enterprise",
		"mission": 5,
	}).Info("Recording new mission")

	kleos.WithFields(kleos.Fields{
		"name": "Spock",
		"rank": "Commander",
		"assignment": "U.S.S. Enterprise",
		"mission": 5,
	}).Info("Added science officer")
}
