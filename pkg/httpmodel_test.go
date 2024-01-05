package back

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

// Sets up the model for testing
func setUpHttpModel(ts *httptest.Server) HttpModel {

	h := HttpModel{
		dest: ts.URL,
	}

	return h
}

// Sets up the Query for testing
func setUpQuery(t *testing.T) message[queryType] {

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
func setUpResponse(t *testing.T) message[responseType] {

	return message[queryType]{
		content: Json{
			"val0": 1.0,
			"val1": 2.0,
			"val2": 3.0,
			"val3": 4.0,
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

// Sets up the test server. This server returns the content of the given
// message[responseType] res. It is implemented as if it was an actual (constant) model
// which always sends back the same input.
func setUpModelPrediction(t *testing.T, res message[responseType]) *httptest.Server {

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(res.content) // res.content is the constant output of the model
	}))
}

// Sets up the test server. This servers returns a fake log.
func setUpModelLogs(t *testing.T, res ModelResponse) *httptest.Server {

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(res.Response) // res.Response is the constant output of the model
	}))
}

// Tests the marshalling and unmarshalling of a query
func TestMarshalUnmarshalQuery(t *testing.T) {
	q := setUpQuery(t)
	mq, err := json.Marshal(q)
	if err != nil {
		err = fmt.Errorf("Could not marshal the query: %w", err)
		fmt.Println(err)
	}
	var uq FrontEndQuery
	err = json.Unmarshal(mq, &uq)
	if err != nil {
		err = fmt.Errorf("Could not unmarshal the query: %w", err)
		fmt.Println(err)
	}

	assert.EqualValues(t, q, uq, fmt.Sprintf("%+v and %+v should be equal", q, uq))
}

// Tests the marshalling and unmarshalling of a response
func TestMarshalUnmarshalResponse(t *testing.T) {
	r := setUpResponse(t)
	mr, err := json.Marshal(r)
	if err != nil {
		err = fmt.Errorf("Could not marshal the response: %w", err)
		fmt.Println(err)
	}
	var ur ModelResponse
	err = json.Unmarshal(mr, &ur)
	if err != nil {
		err = fmt.Errorf("Could not unmarshal the response: %w", err)
		fmt.Println(err)
	}
	assert.EqualValues(t, r, ur, fmt.Sprintf("%+v and %+v should be equal", r, ur))
}

// Tests if the prediction and the returned prediction are identical. It uses
// a mock server
func TestSend(t *testing.T) {

	pred := setUpResponse(t)
	ts := setUpModelPrediction(t, pred)
	defer ts.Close()

	h := setUpHttpModel(ts)

	message := setUpQuery(t)
	encoded, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	res, err := h.send(encoded)
	if err != nil {
		log.Fatal(err)
	}
	testPred := ModelResponse{Response: nil, ResponseType: Predictions, Id: 0} //mock the formating of a model
	testPred.Response = res                                                    // continues the mocking

	assert.EqualValues(t,
		pred,
		testPred,
		fmt.Sprintf("%+v and %+v should be equal", pred, testPred))
}

// Tests wheter the connection with the mock server is well established
func TestServerReceivesResponse(t *testing.T) {
	notEmpty := Json{"notempty": true}
	ts := setUpModelPrediction(t, ModelResponse{Response: notEmpty})
	defer ts.Close()
	h := setUpHttpModel(ts)
	message := "test"
	encoded, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	jsonRes, err := h.send(encoded)
	assert.NotEmpty(t, jsonRes, "Response of server should not be empty")
}

func TestPredict(t *testing.T) {

	q := setUpQuery(t)
	pred := setUpResponse(t)
	ts := setUpModelPrediction(t, pred)
	defer ts.Close()

	h := setUpHttpModel(ts)
	returnValues, err := h.Predict(&q, h.send)
	if err != nil {
		log.Fatal(err)
	}
	testPred := ModelResponse{Response: returnValues, ResponseType: Predictions, Id: 0}
	assert.EqualValues(t, pred, testPred)

}

// Tests if the returned modeler of setUpodel implements the Modeler interface
func TestInterface(t *testing.T) {
	h := setUpHttpModel(setUpModelPrediction(t, ModelResponse{}))
	_, ok := any(&h).(Sender) //Must be &h because an interface is always a pointer type
	assert.True(t, ok, "httpmodeler does NOT implement the modeler interface")
}
