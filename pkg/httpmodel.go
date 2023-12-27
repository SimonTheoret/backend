package main

import (
	"context"
	"fmt"

	requests "github.com/carlmjohnson/requests"
)

type HttpModel struct {
	base baseModeler //base struct
	dest string      // destination for the requests post
}

// Main implementation detail. This function takes in an array of bytes, such as
// the one resulting from calling json.Marshal() and returns the response from
// the model.
func (h *HttpModel) send(body []byte) (Json, error) {
	var result Json
	ctx := context.Background()

	err := requests.
		URL(h.dest).
		BodyJSON(&body).
		ToJSON(&result).
		Fetch(ctx)
	fmt.Println(result)
	if err != nil {
		return nil, fmt.Errorf("Failed to send or receive request: %w", err)
	}

	return result, nil
}
