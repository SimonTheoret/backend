package back

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
)

type HttpModel struct {
	dest string // destination for the requests
}

// This function takes in an array of bytes, such as the one resulting from
// calling json.Marshal() and returns the response from the model. It executes a
// single http request to a model.
func (h *HttpModel) send(body []byte, qt queryType, options ...any) (Json, error) {
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
