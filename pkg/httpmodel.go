package back

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/carlmjohnson/requests"
	"github.com/gin-gonic/gin"
)

type HttpModel struct {
	Modeler            //base struct
	dest    string     // destination for the requests
	out     OutputChan // Channel used to send ModelResponses
	in      InputChan  // Channel used to receive Frontend Queries
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

// Starts model, waiting for queries. Start is launched in its own goroutine
func (h *HttpModel) Start(rf *responseFormatter) {

	for {
		in := <-h.in

		t := in.QueryType
		// var err error
		// var res Json

		switch t {
		case GetLogs: // Case for getting the logs
			res, err := h.GetLogs(h.send)
			h.manageErrAndSend(rf, res, err)
		default: //Consider unknown and prediction queries as prediction
			res, err := h.Predict(&in, h.send)
			if err != nil {
				errBody, _ := json.Marshal(err)
				gin.DefaultErrorWriter.Write(errBody)
			}
			h.manageErrAndSend(rf, res, err)
		}
	}
}

// Helper function. It manages the errors and sends back the *modelResponse.
func (h *HttpModel) manageErrAndSend(rf *responseFormatter, res Json, err error) {
	var modelResponse *ModelResponse
	if err != nil {
		errBody, _ := json.Marshal(err)
		gin.DefaultErrorWriter.Write(errBody)
		jsonError := Json{"error": err.Error()}
		modelResponse, err = rf.FormatRawResponse(jsonError, h, nil)
		if err != nil {
			errBody, _ := json.Marshal(err)
			gin.DefaultErrorWriter.Write(errBody)
		}
	} else {
		modelResponse, err = rf.FormatRawResponse(res, h, nil)
		if err != nil {
			errBody, _ := json.Marshal(err)
			gin.DefaultErrorWriter.Write(errBody)
		}
	}
	h.out <- *modelResponse
}

// Returns the input channel
func (h *HttpModel) QueryChannel() InputChan {
	return h.in
}

// Returns the output channel
func (h *HttpModel) ResponseChannel() OutputChan {
	return h.out
}

// Builds a new HttpModel
func NewHttpModel(name string, id int, dest string) *HttpModel {
	base := NewBase(name, id)
	out := make(chan ModelResponse)
	in := make(chan FrontEndQuery)
	return &HttpModel{base, dest, out, in}
}

// DEPRECATED. This function achieves the same result as the send function.
// Benchmarks are needed to compare these two versions.
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
