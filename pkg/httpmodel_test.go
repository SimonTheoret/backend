package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Setup function. Returns a functionning http model, a basic query and a server
// for testing
func setUp(ts *httptest.Server) HttpModel {

	h := HttpModel{
		base: &base{
			state:     Ready,
			ModelName: "testing base",
			id:        0,
		},
		dest: ts.URL,
	}

	return h
}

func TestNewRequests(t *testing.T) {
	pred := Response{
		Response: Json{
			"val0": rand.Float32(),
			"val1": rand.Float32(),
			"val2": rand.Float32(),
			"val3": rand.Float32(),
		},
		id: 0}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(pred)
	}))
	defer ts.Close()

	h := setUp(ts)
	message := Query{
		Instruction: "Return a response, any response",
		Input: Json{
			"val0": rand.Float32(),
			"val1": rand.Float32(),
			"val2": rand.Float32(),
			"val3": rand.Float32(),
		},
		id: 0,
	}
	encoded, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	res, err := h.send(encoded)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(os.Stdout)
	assert.EqualValues(t, res, pred, fmt.Sprintf("%+v and %+v should be equal", res, pred))
}
