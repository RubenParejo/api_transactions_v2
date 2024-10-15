package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"api_transactions_v2/pkg/model"
)

const url = "http://localhost:8080/transactions"

func LaunchScript() {
	data, err := getMockedData()
	if err != nil {
		log.Printf("LaunchScript - Error getting mocked data")
		return
	}

	start := time.Now()
	callApiPost(data)
	duration := time.Since(start)
	fmt.Printf("Single POST duration: %s\n", duration)

	start = time.Now()
	callApiGet(data.Id)
	duration = time.Since(start)
	fmt.Printf("Single GET duration: %s\n", duration)

	const numCalls = 1000
	const numGoroutines = 5
	successCount := 0
	errorCount := 0
	var wg sync.WaitGroup

	results := make(chan bool, numCalls)
	start = time.Now()

	makeAPICalls := func() {
		defer wg.Done()
		for i := 0; i < numCalls/numGoroutines; i++ {
			if err := callApiPost(data); err != nil {
				results <- false
			} else {
				results <- true
			}

			if err := callApiGet(data.Id); err != nil {
				results <- false
			} else {
				results <- true
			}
		}
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go makeAPICalls()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result {
			successCount++
		} else {
			errorCount++
		}
	}
	duration = time.Since(start)

	fmt.Printf("Total duration for concurrent calls: %s\n", duration)
	fmt.Printf("Successful executions: %d\n", successCount)
	fmt.Printf("Error executions: %d\n", errorCount)
}

func getMockedData() (model.Data, error) {
	layout := "2006-01-02T15:04:05Z"
	locationDateTime, err := time.Parse(layout, "2024-10-20T08:38:34Z")
	if err != nil {
		log.Printf("getMockedData - Error parsing dates: %s\n", err.Error())
		return model.Data{}, err
	}
	data := model.Data{
		Id:               "8834HR43F9FNF3F8J98",
		LocationDateTime: locationDateTime,
		Location:         "North Highway",
		TotalAmount:      10.5,
		Currency:         "EUR",
		Vehicle: model.Vehicle{
			VRM:     "1234BCD",
			Country: "ES",
			Make:    "SEAT",
		},
		Driver: model.Driver{
			FirstName: "Jose",
			LastName:  "Garcia",
			Address1:  "Apple Street",
			Address2:  "",
			PostCode:  "1234",
			City:      "Madrid",
			Region:    "",
			Country:   "ES",
			Phone:     "111-222-333",
			Email:     "josegarcia@abc.es",
		},
	}
	return data, nil
}

func callApiPost(data model.Data) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling data: %s\n", err.Error())
		return err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error making POST request: %s\n", err.Error())
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		log.Println("Data posted successfully")
		return nil
	}
	log.Printf("Failed to post data: %s\n", response.Status)
	return fmt.Errorf("failed to post data: %s", response.Status)
}

func callApiGet(id string) error {
	url_param := url + "?id=" + id

	response, err := http.Get(url_param)
	if err != nil {
		log.Printf("Error making GET request: %s\n", err.Error())
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var data model.Data
		if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
			log.Printf("Error decoding response: %s\n", err.Error())
			return err
		}

		log.Printf("Data retrieved successfully: %+v\n", data)
		return nil
	}
	log.Printf("Failed to get data: %s\n", response.Status)
	return fmt.Errorf("failed to get data: %s", response.Status)
}
