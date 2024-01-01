package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/carlmjohnson/requests"
)

type HttpModel struct {
	basicModeler        //base struct
	dest         string // destination for the requests
}

// This function takes in an array of bytes, such as the one resulting from
// calling json.Marshal() and returns the response from the model. It executes a
// single http request to a model.
func (h *HttpModel) send(body []byte) (Json, error) {
	// var buf bytes.Buffer
	var jsonBody Json
	err := requests.
		URL(h.dest).                // gives the destination url
		BodyBytes(body).            // writes the body of the request as body, a slice of bytes
		ToJSON(&jsonBody).          // convert the response to json and writes it into jsonBody
		Fetch(context.Background()) //Does the
	if err != nil {
		return nil, fmt.Errorf("Error during send method: %w", err)
	}
	return jsonBody, nil
}

// DEPRECATED. This function achieves the same result as the send function.
// Benchmarks are needed to compare these two versions. This version only uses
// the std library.
func (h *HttpModel) sendold(body []byte) (Json, error) {

	request, err := http.NewRequest("POST", h.dest, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error during POST, : %w", err)
	}
	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	fmt.Println("response Body:", string(resBody))
	var content Json
	json.Unmarshal(resBody, &content)
	return content, nil
}
