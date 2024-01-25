package back

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

// This module tests any modeler implementation.

// Tests if the returned modelers of setUpodel implements the sender interface
func TestInterfaceHTTPModeler(t *testing.T) {
	h, _ := setUpModeler(
		t,
		"httpmodeler",
	) // setups HttpModeler

	_, ok := any(h).(Sender) // Must be *modeler because an *modeler implements sender
	assert.True(t, ok, "httpmodeler does NOT implement the Sender interface")
}

func setUpHttpModeler(t *testing.T) (HttpModel, *httptest.Server) {
	ts := setUpModelPrediction(t, setUpResponse())
    return HttpModel{Dest: ts.URL}, ts
}

// Sets up the model for testing.
// Modeler type the string type of the desired modeler
func setUpModeler(t *testing.T, modelerType string, options ...any) (Modeler, any) {
	if strings.ToLower(modelerType) == "httpmodeler" {
        h, _ := setUpHttpModeler(t)
        base := NewBase("modeler test", &h)
		return &base, nil
	} else {
		return nil, nil
	}
}

// Sets up the Query for testing
func setUpQuery() message[queryType] {
	return message[queryType]{
		content: Json{
			"val0": 0.0,
			"val1": 1.0,
			"val2": 2.0,
			"val3": 3.0,
		},
		messageType:     Predict,
		id:              Id(xid.New()),
		creationTime:    time.Now(),
		receiverId:      Id(xid.New()),
		sender:          "yes yes sender yes yes",
		queryOptions:    nil,
		responseOptions: nil,
	}
}

// Sets up the Prediction for testing
func setUpResponse() message[responseType] {
	return message[responseType]{
		content: Json{
			"val0": 1.0,
			"val1": 2.0,
			"val2": 3.0,
			"val3": 4.0,
		},
		messageType:     Predictions,
		id:              Id(xid.New()),
		creationTime:    time.Now(),
		receiverId:      Id(xid.New()),
		sender:          "yes yes sender yes yes",
		queryOptions:    nil,
		responseOptions: nil,
	}
}

// Sets up the test server. This server returns the content of the given
// message[responseType] res. It is implemented as if it was an actual (constant) model
// which always sends back the same input.
func setUpModelPrediction(t *testing.T, res message[responseType]) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(res.content) // res.content is the constant output of the model
	}))
}

// Tests if the prediction and the returned prediction are identical. It uses
// a mock server
func TestSend(t *testing.T) {
	pred := setUpResponse()
	ts := setUpModelPrediction(t, pred)
	defer ts.Close()

	h, _ := setUpModeler(t, "httpmodeler")

	message := setUpQuery()
	encoded, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	res, err := h.send(encoded, UnknownQuery)
	if err != nil {
		log.Fatal(err)
	}
	testPred := setUpResponse()
	testPred.content = res // changes the content to what the model sent back

	assert.EqualValues(t,
		pred,
		testPred,
		fmt.Sprintf("%+v and %+v should be equal", pred, testPred))
}

// Tests predict for a HttpModeler
func TestPredict(t *testing.T) {
	q := setUpQuery()
	pred := setUpResponse()
	ts := setUpModelPrediction(t, pred)
	defer ts.Close()

	h, _ := setUpModeler(t, "httpmodeler")
	returnValues, err := h.Predict(&q, nil)
	if err != nil {
		log.Fatal(err)
	}
	testPred := setUpResponse()
    testPred.content = returnValues
	assert.EqualValues(t, pred, testPred)
}
